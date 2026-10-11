#!/bin/bash

# NGO Donation Tracker - Production Deployment Script
# This script handles the deployment of the application to production

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_NAME="ngo-donation-tracker"
COMPOSE_FILE="docker-compose.prod.yml"
ENV_FILE=".env.production"
BACKUP_DIR="./backups"
LOG_FILE="./deploy.log"

# Functions
log() {
    echo -e "${GREEN}[$(date +'%Y-%m-%d %H:%M:%S')]${NC} $1" | tee -a "$LOG_FILE"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1" | tee -a "$LOG_FILE"
    exit 1
}

warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1" | tee -a "$LOG_FILE"
}

info() {
    echo -e "${BLUE}[INFO]${NC} $1" | tee -a "$LOG_FILE"
}

check_requirements() {
    log "Checking deployment requirements..."
    
    # Check if Docker is installed
    if ! command -v docker &> /dev/null; then
        error "Docker is not installed. Please install Docker first."
    fi
    
    # Check if Docker Compose is installed
    if ! command -v docker-compose &> /dev/null; then
        error "Docker Compose is not installed. Please install Docker Compose first."
    fi
    
    # Check if environment file exists
    if [ ! -f "$ENV_FILE" ]; then
        error "Environment file $ENV_FILE not found. Please copy .env.production.template to $ENV_FILE and configure it."
    fi
    
    # Check if SSL certificates exist (in production)
    if [ "$ENVIRONMENT" = "production" ] && [ ! -f "./nginx/ssl/cert.pem" ]; then
        warning "SSL certificates not found. HTTPS will not be available."
    fi
    
    log "Requirements check completed successfully"
}

backup_database() {
    log "Creating database backup..."
    
    # Create backup directory if it doesn't exist
    mkdir -p "$BACKUP_DIR"
    
    # Generate backup filename with timestamp
    BACKUP_FILE="$BACKUP_DIR/ngo_tracker_$(date +%Y%m%d_%H%M%S).sql"
    
    # Check if database container is running
    if docker-compose -f "$COMPOSE_FILE" ps postgres | grep -q "Up"; then
        # Create database backup
        docker-compose -f "$COMPOSE_FILE" exec -T postgres pg_dump -U ngo_user ngo_tracker > "$BACKUP_FILE"
        
        if [ $? -eq 0 ]; then
            log "Database backup created: $BACKUP_FILE"
            
            # Compress backup
            gzip "$BACKUP_FILE"
            log "Database backup compressed: $BACKUP_FILE.gz"
        else
            error "Failed to create database backup"
        fi
    else
        warning "Database container is not running. Skipping backup."
    fi
}

build_images() {
    log "Building Docker images..."
    
    # Build backend image
    info "Building backend image..."
    docker-compose -f "$COMPOSE_FILE" build backend
    
    if [ $? -ne 0 ]; then
        error "Failed to build backend image"
    fi
    
    # Build frontend image
    info "Building frontend image..."
    docker-compose -f "$COMPOSE_FILE" build frontend
    
    if [ $? -ne 0 ]; then
        error "Failed to build frontend image"
    fi
    
    log "Docker images built successfully"
}

deploy_services() {
    log "Deploying services..."
    
    # Pull external images
    info "Pulling external images..."
    docker-compose -f "$COMPOSE_FILE" pull postgres redis nginx prometheus grafana elasticsearch kibana
    
    # Start core services first
    info "Starting core services..."
    docker-compose -f "$COMPOSE_FILE" up -d postgres redis
    
    # Wait for database to be ready
    info "Waiting for database to be ready..."
    timeout=60
    counter=0
    while ! docker-compose -f "$COMPOSE_FILE" exec -T postgres pg_isready -U ngo_user -d ngo_tracker; do
        sleep 2
        counter=$((counter + 2))
        if [ $counter -ge $timeout ]; then
            error "Database failed to start within $timeout seconds"
        fi
    done
    
    # Run database migrations
    info "Running database migrations..."
    # Add migration commands here when available
    
    # Start application services
    info "Starting application services..."
    docker-compose -f "$COMPOSE_FILE" up -d backend frontend nginx
    
    # Start monitoring services (if profile is enabled)
    if [ "$ENABLE_MONITORING" = "true" ]; then
        info "Starting monitoring services..."
        docker-compose -f "$COMPOSE_FILE" --profile monitoring up -d
    fi
    
    # Start logging services (if profile is enabled)
    if [ "$ENABLE_LOGGING" = "true" ]; then
        info "Starting logging services..."
        docker-compose -f "$COMPOSE_FILE" --profile logging up -d
    fi
    
    log "Services deployed successfully"
}

health_check() {
    log "Performing health checks..."
    
    # Wait for services to be healthy
    info "Waiting for services to be healthy..."
    
    # Check backend health
    timeout=120
    counter=0
    while ! curl -sf http://localhost:8080/health > /dev/null 2>&1; do
        sleep 5
        counter=$((counter + 5))
        if [ $counter -ge $timeout ]; then
            error "Backend health check failed after $timeout seconds"
        fi
    done
    info "Backend is healthy"
    
    # Check frontend health
    counter=0
    while ! curl -sf http://localhost:3000/health > /dev/null 2>&1; do
        sleep 5
        counter=$((counter + 5))
        if [ $counter -ge $timeout ]; then
            error "Frontend health check failed after $timeout seconds"
        fi
    done
    info "Frontend is healthy"
    
    # Check nginx health
    counter=0
    while ! curl -sf http://localhost:80/health > /dev/null 2>&1; do
        sleep 5
        counter=$((counter + 5))
        if [ $counter -ge $timeout ]; then
            error "Nginx health check failed after $timeout seconds"
        fi
    done
    info "Nginx is healthy"
    
    log "All health checks passed"
}

cleanup_old_images() {
    log "Cleaning up old Docker images..."
    
    # Remove unused images
    docker image prune -f
    
    # Remove unused volumes
    docker volume prune -f
    
    log "Cleanup completed"
}

show_status() {
    log "Deployment status:"
    echo
    docker-compose -f "$COMPOSE_FILE" ps
    echo
    
    info "Application URLs:"
    echo "  Frontend: http://localhost:80 (or https://localhost:443 if SSL is configured)"
    echo "  API: http://localhost:8080"
    echo "  Metrics: http://localhost:9090"
    
    if [ "$ENABLE_MONITORING" = "true" ]; then
        echo "  Grafana: http://localhost:3001"
        echo "  Prometheus: http://localhost:9091"
    fi
    
    if [ "$ENABLE_LOGGING" = "true" ]; then
        echo "  Kibana: http://localhost:5601"
        echo "  Elasticsearch: http://localhost:9200"
    fi
    echo
}

rollback() {
    log "Rolling back deployment..."
    
    # Stop current services
    docker-compose -f "$COMPOSE_FILE" down
    
    # Restore from latest backup if available
    LATEST_BACKUP=$(ls -t "$BACKUP_DIR"/*.sql.gz 2>/dev/null | head -n1)
    if [ -n "$LATEST_BACKUP" ]; then
        warning "Restoring from backup: $LATEST_BACKUP"
        # Add restoration commands here
    fi
    
    # Restart services with previous images
    # This would require tagging images with versions
    warning "Rollback functionality requires version tagging implementation"
    
    log "Rollback completed"
}

main() {
    case "${1:-deploy}" in
        "deploy")
            log "Starting deployment of $PROJECT_NAME..."
            check_requirements
            backup_database
            build_images
            deploy_services
            health_check
            cleanup_old_images
            show_status
            log "Deployment completed successfully!"
            ;;
        "update")
            log "Updating $PROJECT_NAME..."
            backup_database
            build_images
            docker-compose -f "$COMPOSE_FILE" up -d --force-recreate
            health_check
            cleanup_old_images
            show_status
            log "Update completed successfully!"
            ;;
        "rollback")
            rollback
            ;;
        "status")
            show_status
            ;;
        "stop")
            log "Stopping $PROJECT_NAME..."
            docker-compose -f "$COMPOSE_FILE" down
            log "Services stopped"
            ;;
        "restart")
            log "Restarting $PROJECT_NAME..."
            docker-compose -f "$COMPOSE_FILE" restart
            health_check
            show_status
            log "Services restarted"
            ;;
        "logs")
            docker-compose -f "$COMPOSE_FILE" logs -f ${2:-}
            ;;
        "backup")
            backup_database
            ;;
        "cleanup")
            cleanup_old_images
            ;;
        *)
            echo "Usage: $0 {deploy|update|rollback|status|stop|restart|logs|backup|cleanup}"
            echo
            echo "Commands:"
            echo "  deploy   - Full deployment (default)"
            echo "  update   - Update existing deployment"
            echo "  rollback - Rollback to previous version"
            echo "  status   - Show deployment status"
            echo "  stop     - Stop all services"
            echo "  restart  - Restart all services"
            echo "  logs     - Show logs (optionally specify service name)"
            echo "  backup   - Create database backup"
            echo "  cleanup  - Clean up unused Docker resources"
            exit 1
            ;;
    esac
}

# Load environment variables
if [ -f "$ENV_FILE" ]; then
    source "$ENV_FILE"
fi

# Run main function
main "$@"
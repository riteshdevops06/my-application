# NGO Donation Transparency Tracker

A comprehensive platform for transparent NGO donation tracking with milestone-based progress monitoring, real-time updates, and advanced analytics.

## 🚀 Features

### Core Features
- **Multi-role Authentication**: Donors, NGO Admins, Platform Admins
- **Campaign Management**: Create, manage, and track donation campaigns
- **Milestone Tracking**: Progress-based milestone system with evidence uploads
- **Expense Management**: Transparent expense tracking and approval workflow
- **Real-time Updates**: Live donation feeds and campaign progress
- **File Uploads**: Secure document and image management
- **Advanced Analytics**: Charts, reports, and performance metrics

### Security Features
- **JWT Authentication** with refresh tokens
- **Role-based Access Control** (RBAC)
- **Rate Limiting** and DDoS protection
- **Input Validation** and sanitization
- **HTTPS/TLS** encryption
- **Security Headers** (CSP, HSTS, etc.)
- **File Upload Security** with virus scanning support

### Production Features
- **Docker Containerization** for easy deployment
- **Horizontal Scaling** support
- **Database Migrations** and backups
- **Monitoring & Alerting** with Prometheus/Grafana
- **Centralized Logging** with ELK stack
- **Health Checks** and graceful shutdowns
- **Performance Optimization** and caching

## 🏗️ Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   React SPA     │    │   Go Backend    │    │   PostgreSQL    │
│   (Frontend)    │◄──►│     (API)       │◄──►│   (Database)    │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         │              ┌─────────────────┐             │
         │              │     Redis       │             │
         └──────────────►│    (Cache)      │◄────────────┘
                        └─────────────────┘
                                 │
         ┌─────────────────────────────────────────────┐
         │              Monitoring Stack                │
         │  ┌─────────────┐  ┌─────────────┐          │
         │  │ Prometheus  │  │   Grafana   │          │
         │  └─────────────┘  └─────────────┘          │
         └─────────────────────────────────────────────┘
```

## 📋 Prerequisites

- **Docker** 20.10+ and **Docker Compose** 2.0+
- **Node.js** 18+ (for local development)
- **Go** 1.21+ (for local development)
- **PostgreSQL** 15+ (for local development)
- **Redis** 7+ (for caching)

## 🚀 Quick Start

### 1. Clone the Repository

```bash
git clone <repository-url>
cd ngo-donation-tracker
```

### 2. Environment Setup

```bash
# Copy environment template
cp .env.production.template .env.production

# Edit the configuration
nano .env.production
```

**Important**: Update all passwords, secrets, and API keys in the environment file.

### 3. SSL Certificates (Production)

For HTTPS in production, add SSL certificates:

```bash
mkdir -p nginx/ssl
# Copy your SSL certificate files
cp /path/to/your/cert.pem nginx/ssl/cert.pem
cp /path/to/your/key.pem nginx/ssl/key.pem
```

### 4. Deploy

```bash
# Make deployment script executable (Linux/Mac)
chmod +x deploy.sh

# Deploy the application
./deploy.sh deploy
```

### 5. Access the Application

- **Frontend**: http://localhost (or https://localhost with SSL)
- **API**: http://localhost:8080
- **Monitoring**: http://localhost:3001 (Grafana)
- **Metrics**: http://localhost:9091 (Prometheus)

## 🛠️ Development Setup

### Backend Development

```bash
cd backend

# Install dependencies
go mod download

# Set up development database
createdb ngo_tracker_dev

# Run migrations (when available)
# go run ./cmd/migrate up

# Start development server
go run ./cmd/api
```

### Frontend Development

```bash
cd frontend

# Install dependencies
npm install

# Start development server
npm run dev
```

## 🐳 Docker Commands

### Production Deployment

```bash
# Full deployment
./deploy.sh deploy

# Update deployment
./deploy.sh update

# Show status
./deploy.sh status

# View logs
./deploy.sh logs [service-name]

# Stop services
./deploy.sh stop

# Restart services
./deploy.sh restart
```

### Development with Docker

```bash
# Start development environment
docker-compose up -d

# View logs
docker-compose logs -f

# Stop services
docker-compose down

# Rebuild services
docker-compose build --no-cache
```

## 🔧 Configuration

### Environment Variables

Key configuration options in `.env.production`:

#### Database Configuration
```env
DB_HOST=postgres
DB_PASSWORD=your_secure_password
DB_NAME=ngo_tracker
```

#### Security Configuration
```env
JWT_SECRET=your_jwt_secret_32_characters_min
ENCRYPTION_KEY=your_encryption_key_32_chars
BCRYPT_COST=12
```

#### Email Configuration
```env
SMTP_HOST=your_smtp_server
SMTP_USER=your_smtp_username
SMTP_PASSWORD=your_smtp_password
FROM_ADDRESS=noreply@yourdomain.com
```

#### File Upload Configuration
```env
UPLOAD_MAX_FILE_SIZE=10485760
USE_S3=false
S3_BUCKET=your_s3_bucket
```

### Rate Limiting

Configure rate limits for different endpoints:

```env
RATE_LIMIT_GLOBAL=1000      # 1000 requests per hour
RATE_LIMIT_AUTH=10          # 10 auth requests per minute  
RATE_LIMIT_API=100          # 100 API requests per minute
RATE_LIMIT_UPLOAD=5         # 5 upload requests per minute
```

## 📊 Monitoring & Observability

### Prometheus Metrics

The application exposes metrics at `/metrics`:

- HTTP request duration and count
- Database connection pool stats
- Cache hit/miss rates
- Custom business metrics (donations, campaigns, etc.)

### Grafana Dashboards

Access Grafana at `http://localhost:3001`:

- **Application Dashboard**: Request rates, response times, errors
- **Infrastructure Dashboard**: CPU, memory, disk usage
- **Business Dashboard**: Donations, campaigns, user activity

### Logging

Structured JSON logging with configurable levels:

- **Console**: Development logging
- **File**: Production logging with rotation
- **Elasticsearch**: Centralized logging (with ELK stack)

### Health Checks

Health check endpoints:

- **Application**: `GET /health`
- **Database**: `GET /health/db`
- **Cache**: `GET /health/cache`

## 🔒 Security

### Security Features

1. **Authentication & Authorization**
   - JWT tokens with configurable expiration
   - Role-based access control
   - Password strength requirements

2. **Input Validation**
   - Server-side validation for all inputs
   - SQL injection prevention
   - XSS protection

3. **Transport Security**
   - HTTPS/TLS encryption
   - Security headers (HSTS, CSP, etc.)
   - Secure cookie settings

4. **Rate Limiting**
   - Per-IP rate limiting
   - Endpoint-specific limits
   - Redis-backed distributed limiting

5. **File Upload Security**
   - File type validation
   - Size limits
   - Virus scanning support
   - Secure file storage

### Security Configuration

```nginx
# Example security headers in nginx.conf
add_header X-Frame-Options "DENY" always;
add_header X-Content-Type-Options "nosniff" always;
add_header X-XSS-Protection "1; mode=block" always;
add_header Strict-Transport-Security "max-age=31536000" always;
```

## 🧪 Testing

### Backend Testing

```bash
cd backend

# Run unit tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run integration tests
go test -tags=integration ./tests/...

# Load testing
go test -bench=. ./tests/...
```

### Frontend Testing

```bash
cd frontend

# Run unit tests
npm test

# Run tests with coverage
npm run test:coverage

# Run e2e tests
npm run test:e2e
```

## 📈 Performance Optimization

### Backend Optimizations

- **Database Connection Pooling**: Configurable pool sizes
- **Redis Caching**: Session and data caching
- **Compression**: Gzip compression for responses
- **Static Asset Optimization**: CDN-ready file serving

### Frontend Optimizations

- **Code Splitting**: Lazy loading of routes and components
- **Bundle Optimization**: Tree shaking and minification
- **Image Optimization**: Automatic image compression
- **Caching**: Browser and CDN caching strategies

### Database Optimizations

- **Indexing**: Strategic database indexes
- **Query Optimization**: Efficient queries with proper joins
- **Connection Pooling**: Optimized connection management
- **Read Replicas**: Read scaling support

## 🔄 Backup & Recovery

### Automated Backups

```bash
# Create backup
./deploy.sh backup

# Scheduled backups (add to crontab)
0 2 * * * /path/to/ngo-donation-tracker/deploy.sh backup
```

### Manual Backup

```bash
# Database backup
docker-compose exec postgres pg_dump -U ngo_user ngo_tracker > backup.sql

# File uploads backup
tar -czf uploads_backup.tar.gz uploads/
```

### Recovery

```bash
# Restore database
docker-compose exec -T postgres psql -U ngo_user ngo_tracker < backup.sql

# Restore files
tar -xzf uploads_backup.tar.gz
```

## 🚀 Deployment Options

### Docker Compose (Recommended)

- **Development**: `docker-compose.yml`
- **Production**: `docker-compose.prod.yml`
- **Monitoring**: Use `--profile monitoring`
- **Logging**: Use `--profile logging`

### Kubernetes (Advanced)

For Kubernetes deployment, see the `k8s/` directory for manifests.

### Cloud Deployment

#### AWS
- Use ECS or EKS for container orchestration
- RDS for PostgreSQL
- ElastiCache for Redis
- S3 for file storage

#### Google Cloud
- Use GKE for Kubernetes
- Cloud SQL for PostgreSQL
- Cloud Memorystore for Redis
- Cloud Storage for files

#### Azure
- Use AKS for Kubernetes
- Azure Database for PostgreSQL
- Azure Cache for Redis
- Azure Blob Storage for files

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Run the test suite
6. Submit a pull request

### Development Guidelines

- Follow Go and React best practices
- Write comprehensive tests
- Document new features
- Follow the existing code style
- Update documentation

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🆘 Support

### Documentation

- **API Documentation**: Available at `/api/docs` when running
- **Architecture Guide**: See `docs/architecture.md`
- **Deployment Guide**: See `docs/deployment.md`

### Getting Help

1. Check the documentation
2. Search existing issues
3. Create a new issue with:
   - Environment details
   - Steps to reproduce
   - Expected vs actual behavior
   - Relevant logs

### Common Issues

#### Database Connection Issues
```bash
# Check database status
docker-compose exec postgres pg_isready -U ngo_user

# View database logs
docker-compose logs postgres
```

#### SSL/HTTPS Issues
```bash
# Verify certificate files
ls -la nginx/ssl/

# Check nginx configuration
docker-compose exec nginx nginx -t
```

#### Performance Issues
```bash
# Check resource usage
docker stats

# View application metrics
curl http://localhost:9090/metrics
```

## 🗺️ Roadmap

### Phase 1 ✅
- Basic authentication and authorization
- Campaign and donation management
- File upload system
- Basic monitoring

### Phase 2 ✅
- Advanced UI components
- Real-time notifications
- Search and filtering
- Security hardening

### Phase 3 (Future)
- Mobile application
- Payment gateway integration
- Advanced analytics
- Multi-language support
- API rate limiting improvements
- Blockchain integration for transparency

---

**Built with ❤️ for NGO transparency and accountability**
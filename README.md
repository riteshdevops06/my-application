# NGO Donation Transparency Tracker

A full-stack platform for transparent donation tracking with three user roles:
- **Donors** - Browse campaigns, make donations, track giving history
- **NGO Admins** - Manage campaigns, milestones, expenses, verification
- **Platform Admins** - Oversee NGOs, campaigns, platform settings, reports

---

## 🏗️ Tech Stack

### Frontend
- **React 19** + **Vite 8** - Modern React with fast build tool
- **React Router 7** - Client-side routing with role-based protection
- **TailwindCSS 4** - Utility-first CSS framework
- **Chart.js** - Data visualization for dashboards
- **Axios** - HTTP client with interceptors

### Backend
- **Go 1.23** - High-performance backend
- **Gorilla Mux** - HTTP router
- **PostgreSQL** - Primary database
- **lib/pq** - PostgreSQL driver
- **golang-jwt/jwt** - JWT authentication
- **bcrypt** - Password hashing
- **zerolog** - Structured logging

### DevOps
- **Docker & Docker Compose** - Containerization
- **Nginx** - Reverse proxy for frontend

---

## 📁 Project Structure

```
ngo-donation-tracker/
├── backend/                 # Go backend
│   ├── cmd/api/            # Main entry point
│   ├── config/             # Configuration management
│   ├── database/           # Database connection & migrations
│   ├── handlers/           # HTTP handlers
│   ├── middleware/         # Auth, logging, CORS, recovery
│   ├── models/             # Data models (to be created)
│   ├── repositories/       # Data access layer (to be created)
│   ├── routes/             # Route definitions
│   ├── services/           # Business logic (to be created)
│   ├── utils/              # Helpers (JWT, password, pagination)
│   └── uploads/            # File upload handling
├── frontend/               # React frontend
│   ├── src/
│   │   ├── components/     # Reusable UI components
│   │   ├── context/        # React context (Auth, Notifications)
│   │   ├── hooks/          # Custom React hooks
│   │   ├── layouts/        # Layout components (Public, Donor, NGO, Admin)
│   │   ├── pages/          # Page components by role
│   │   ├── routes/         # Route configuration
│   │   ├── services/       # API service layer
│   │   ├── utils/          # Utility functions
│   │   ├── App.jsx         # Main app with routing
│   │   └── main.jsx        # Entry point
│   └── public/             # Static assets
├── database/
│   ├── migrations/         # SQL migration files
│   └── seeds/              # Seed data scripts
├── docs/                   # Documentation
├── docker-compose.yml      # Local development stack
└── README.md               # This file
```

---

## 🚀 Quick Start

### Prerequisites
- **Docker & Docker Compose** (recommended)
- OR: **Go 1.23+**, **Node.js 20+**, **PostgreSQL 16+**

### Option 1: Docker Compose (Recommended)

```bash
# Clone and navigate
cd ngo-donation-tracker

# Start all services
docker-compose up -d

# View logs
docker-compose logs -f

# Stop services
docker-compose down
```

Services will be available at:
- **Frontend**: http://localhost:5173
- **Backend API**: http://localhost:8080/api/v1
- **Database**: localhost:5432

### Option 2: Manual Setup

#### Backend
```bash
cd ngo-donation-tracker/backend

# Copy env file
cp .env.example .env

# Edit .env with your database credentials
# DATABASE_URL=postgres://user:pass@localhost:5432/ngo_tracker?sslmode=disable

# Install dependencies
go mod tidy

# Run database migrations (requires running PostgreSQL)
go run cmd/api/main.go
```

#### Frontend
```bash
cd ngo-donation-tracker/frontend

# Install dependencies
npm install

# Start dev server
npm run dev
```

#### Database Seeding
```bash
cd database/seeds
go run seed.go
```

---

## 🔐 Default Test Accounts

After seeding the database:

| Role | Email | Password |
|------|-------|----------|
| Platform Admin | admin@ngotracker.org | password123 |
| Donor | donor@example.com | password123 |
| NGO Admin | ngo@example.com | password123 |

---

## 📡 API Endpoints

### Authentication
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/auth/login` | User login |
| POST | `/api/v1/auth/register` | User registration |
| POST | `/api/v1/auth/refresh` | Refresh access token |
| POST | `/api/v1/auth/forgot-password` | Request password reset |
| POST | `/api/v1/auth/reset-password` | Reset password |

### Public Campaigns
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/campaigns` | List campaigns (with filters) |
| GET | `/api/v1/campaigns/{id}` | Get campaign details |
| GET | `/api/v1/campaigns/{id}/milestones` | Get campaign milestones |

### Donor (Requires `donor` role)
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/donor/dashboard` | Donor dashboard stats |
| GET | `/api/v1/donor/donations` | List donor's donations |
| GET | `/api/v1/donor/donations/{id}` | Get donation details |
| POST | `/api/v1/donor/donations` | Create new donation |

### NGO Admin (Requires `ngo_admin` role)
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/ngo/dashboard` | NGO dashboard stats |
| GET | `/api/v1/ngo/profile` | Get NGO profile |
| PUT | `/api/v1/ngo/profile` | Update NGO profile |
| GET | `/api/v1/ngo/campaigns` | List NGO's campaigns |
| POST | `/api/v1/ngo/campaigns` | Create campaign |
| GET | `/api/v1/ngo/campaigns/{id}` | Get campaign |
| PUT | `/api/v1/ngo/campaigns/{id}` | Update campaign |
| DELETE | `/api/v1/ngo/campaigns/{id}` | Delete campaign |
| POST | `/api/v1/ngo/campaigns/{id}/publish` | Publish campaign |
| GET | `/api/v1/ngo/campaigns/{campaign_id}/milestones` | List milestones |
| POST | `/api/v1/ngo/campaigns/{campaign_id}/milestones` | Create milestone |
| GET | `/api/v1/ngo/campaigns/{campaign_id}/expenses` | List expenses |
| POST | `/api/v1/ngo/campaigns/{campaign_id}/expenses` | Create expense |
| GET | `/api/v1/ngo/verification` | Get verification status |
| POST | `/api/v1/ngo/verification` | Submit verification |

### Platform Admin (Requires `platform_admin` role)
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/admin/dashboard` | Platform dashboard |
| GET | `/api/v1/admin/ngos` | List all NGOs |
| GET | `/api/v1/admin/ngos/{id}` | Get NGO details |
| POST | `/api/v1/admin/ngos/{id}/verify` | Verify NGO |
| POST | `/api/v1/admin/ngos/{id}/reject` | Reject NGO |
| GET | `/api/v1/admin/campaigns` | List all campaigns |
| GET | `/api/v1/admin/milestones` | List all milestones |
| POST | `/api/v1/admin/milestones/{id}/approve` | Approve milestone |
| POST | `/api/v1/admin/milestones/{id}/reject` | Reject milestone |
| GET | `/api/v1/admin/reports/donations` | Donation reports |
| GET | `/api/v1/admin/reports/ngos` | NGO reports |
| GET | `/api/v1/admin/reports/campaigns` | Campaign reports |
| GET | `/api/v1/admin/settings` | Get platform settings |
| PUT | `/api/v1/admin/settings` | Update platform settings |

---

## 🔧 Development

### Running Tests
```bash
# Backend tests
cd backend && go test ./...

# Frontend tests
cd frontend && npm test
```

### Code Quality
```bash
# Backend linting
cd backend && go vet ./...

# Frontend linting
cd frontend && npm run lint
```

### Database Migrations
Migrations run automatically on backend startup. To create new migrations, add SQL files to `database/migrations/`.

---

## 🌐 Environment Variables

### Backend (.env)
| Variable | Description | Default |
|----------|-------------|---------|
| PORT | Server port | 8080 |
| DATABASE_URL | PostgreSQL connection string | Required |
| JWT_SECRET | JWT signing secret | Required |
| JWT_EXPIRY_HOURS | Token expiry in hours | 24 |
| ENVIRONMENT | Environment (development/production) | development |

### Frontend
| Variable | Description | Default |
|----------|-------------|---------|
| VITE_API_URL | Backend API base URL | http://localhost:8080/api/v1 |

---

## 📄 License

MIT License - see LICENSE file for details.
# WebSocket Chat Server

A real-time chat server built with Go, JS, WebSockets, and PostgreSQL.

## Features

- Real-time messaging via WebSockets
- JWT-based authentication
- PostgreSQL database as a data storage
- Docker containerization

## Prerequisites

- Docker and Docker Compose
- Go 1.23+ (for local development)

## Quick Start with Docker

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd websocket-chat-server
   ```

2. **Start the services**
   ```bash
   make dev
   ```
   or manually:
   ```bash
   docker-compose up -d
   ```

3. **Access the application**
   - Web interface: http://localhost:9080/chat/
   - WebSocket endpoint: ws://localhost:9080/chat/ws/

## Environment Variables

Create a `.env` file in the root directory:

```env
ADDR=:9080
ALLOWED_ORIGINS=http://localhost:9080
JWT_KEY=your-secret-jwt-key-here
```

## Database

The PostgreSQL database is automatically initialized with the following tables:
- `users` - User accounts and authentication
- `sessions` - JWT token management

### Database Commands

```bash
# Connect to database
make db-connect

# Run migrations
make db-migrate
```

## Development

### Local Development

```bash
# Build the application
make build

# Run locally
make run

# Run tests
make test
```

### Docker Commands

```bash
# Build containers
make docker-build

# Start services
make docker-up

# Stop services
make docker-down

# View logs
make docker-logs

# Clean up (removes volumes)
make docker-clean
```

## Architecture

The application follows SOLID principles

### Project Structure

```
├── cmd/
│   └── main.go              # Application entry point
├── db/
│   ├── migrations/          # database sql migrations
├── internal/
│   ├── auth/                # Authentication logic
│   ├── controller/          # WebSocket hub and client management
│   ├── handler/             # HTTP and WebSocket handlers
│   ├── middleware/          # HTTP middleware
│   └── models/              # Data models
├── src/                     # Static files
├── utils/                   # Utility functions
├── docker-compose.yml       # Docker services configuration
└── Dockerfile              # Application container
```

## API Endpoints

- `GET /chat/` - Web interface
- `GET /chat/ws/` - WebSocket connection (requires JWT authentication)

## Security

- JWT tokens for authentication
- Password hashing (bcrypt)
- Input validation

## Contributing

Not being considered at this time
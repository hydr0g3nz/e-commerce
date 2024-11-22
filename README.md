# E-commerce Backend API 

A modern e-commerce backend built with Go 1.23, following hexagonal architecture principles. This API provides comprehensive endpoints for managing products, categories, orders, and user authentication.

## 🛠 Tech Stack

- **Language:** Go 1.23
- **Framework:** Fiber
- **Database:** MongoDB
- **Architecture:** Hexagonal (Ports & Adapters)
- **Authentication:** JWT (Access & Refresh Tokens)
- **Containerization:** Docker
- **Message Queue:** RabbitMQ (for order processing)

## ✨ Features

- 🔐 JWT-based Authentication with Refresh Tokens
- 📦 Product Management with Variations
- 🗂 Category Management with Hierarchical Structure
- 🛒 Order Processing with Stock Management
- 🖼 Image Upload for Products
- 👥 Role-based Access Control (Admin/User)
- 🔄 Automatic Stock Management System

## 🏗 Architecture

The project follows hexagonal architecture (also known as ports and adapters) with the following structure:

```
internal/
├── adapters/
│   ├── handler/      # HTTP handlers
│   ├── middleware/   # HTTP middleware
│   ├── model/        # Database models
│   └── repository/   # Database implementations
├── core/
│   ├── domain/       # Business entities
│   ├── ports/        # Interfaces
│   └── services/     # Business logic
└── config/           # Application configuration
```

## 🔑 API Endpoints

### Authentication
```
POST /api/v1/auth/register    # Register new user
POST /api/v1/auth/login       # Login user
POST /api/v1/auth/refresh     # Refresh access token
```

### Categories
```
GET    /api/v1/category          # Get all categories
GET    /api/v1/category/:id      # Get category by ID
POST   /api/v1/category          # Create category (Admin)
PUT    /api/v1/category          # Update category (Admin)
DELETE /api/v1/category/:id      # Delete category (Admin)
POST   /api/v1/category/product  # Add product to category (Admin)
DELETE /api/v1/category/:cat_id/product/:prod_id  # Remove product from category (Admin)
```

### Products
```
GET    /api/v1/product                           # Get all products
GET    /api/v1/product/:id                       # Get product by ID
POST   /api/v1/product                           # Create product (Admin)
PUT    /api/v1/product                           # Update product (Admin)
DELETE /api/v1/product/:prod_id                  # Delete product (Admin)
POST   /api/v1/product/variant/:prod_id          # Add product variation (Admin)
DELETE /api/v1/product/:prod_id/variant/:var_id  # Remove product variation (Admin)
POST   /api/v1/product/image                     # Upload product image (Admin)
DELETE /api/v1/product/image/:filename           # Delete product image (Admin)
```

### Orders
```
POST   /api/v1/order  # Create order (Authenticated)
```

## 🚀 Getting Started

### Prerequisites
- Go 1.23 or higher
- MongoDB
- RabbitMQ
- Docker & Docker Compose

### Configuration
Create a `config.yml` file in the root directory:

```yaml
server:
  port: 8080
  path: "/api"

mongo:
  uri: "mongodb://localhost:27017"
  database: "e-commerce"

key:
  accessToken: "your-access-token-secret"
  refreshToken: "your-refresh-token-secret"

amqp:
  url: "amqp://guest:guest@localhost:5672/"

upload:
  serverPath: "./uploads"
```

### Running with Docker

1. Build the Docker image:
```bash
docker build -t ecommerce-api .
```

2. Run with Docker Compose:
```bash
docker-compose up -d
```

### Running Locally

1. Install dependencies:
```bash
go mod download
```

2. Run the application:
```bash
go run main.go
```

## 🔒 Security Features

- JWT-based authentication
- Role-based access control
- Password hashing with bcrypt
- Refresh token rotation
- Request rate limiting
- CORS protection
- Panic recovery middleware

## 💡 Key Implementation Details

### Product Variations
Products support multiple variations with:
- SKU tracking
- Stock management
- Size and color options
- Individual pricing
- Sale price support

### Order Processing
- Atomic stock updates
- RabbitMQ for async processing
- Automatic stock reservation
- Failed transaction handling

### Image Management
- Support for multiple product images
- Secure file upload
- Automatic file type validation
- Server-side image storage

## 📝 License

MIT License

## 🤝 Contributing

Contributions, issues, and feature requests are welcome!
# Digital Library Analytics Dashboard

A full-stack application for managing a digital library and visualizing analytics.

## Project Structure

- `backend/`: Golang (Fiber) REST API
- `frontend/`: Next.js (TypeScript, Tailwind) Web Application

## Features

- User authentication with JWT
- User registration system
- Book management (CRUD operations)
- Lending record management
- Search and filter functionality for books
- API Documentation with Swagger/OpenAPI and Redoc

## Quick Start with Docker

The easiest way to run the application is using Docker Compose. This will set up all necessary services including the database with sample data.

### Prerequisites

- [Docker](https://docs.docker.com/get-docker/)
- [Docker Compose](https://docs.docker.com/compose/install/)

### Running with Docker

1. **Clone the repository** (if you haven't already):
   ```bash
   git clone https://github.com/agussuwerdo/digital-library.git
   cd digital-library
   ```

2. **Build and start the containers**:
   ```bash
   docker-compose up --build
   ```

   This will:
   - Build the frontend and backend containers
   - Start a PostgreSQL database
   - Initialize the database with schema and sample data
   - Start all services

3. **Access the application**:
   - Frontend: [http://localhost:3000](http://localhost:3000)
   - Backend API: [http://localhost:3001](http://localhost:3001)
   - Swagger Documentation: [http://localhost:3001/api/apidocs/](http://localhost:3001/api/apidocs/)

### Docker Services

The Docker Compose setup includes three services:

1. **Frontend** (`frontend`):
   - Next.js application
   - Port: 3000
   - Environment: Development

2. **Backend** (`backend`):
   - Go Fiber API
   - Port: 3001
   - Connected to PostgreSQL database

3. **Database** (`db`):
   - PostgreSQL 15
   - Port: 5432
   - Persistent volume for data storage
   - Pre-populated with sample data

### Sample Data

The database is automatically populated with:

- **Users**:
  - Admin: 
    - Username: `admin`
    - Email: `admin@example.com`
    - Password: `123`
    - Role: `admin`
  - Regular User:
    - Username: `user`
    - Email: `user@example.com`
    - Password: `123`
    - Role: `user`

- **Books**: 10 sample books across different categories
- **Lending Records**: 5 sample lending records

### User Roles and Permissions

The application has two user roles with different access levels:

1. **Admin Role**:
   - Can access all features and pages
   - Can manage books (add, edit, delete)
   - Can manage lending records
   - Can access the admin dashboard
   - Can view API documentation

2. **User Role**:
   - Can only borrow books
   - Can view their own lending history
   - Can view available books

The navigation menu automatically adjusts based on the user's role:
- Admin users see all navigation items
- Regular users only see the "Lending" section

### Docker Commands

- **Start services**:
  ```bash
  docker-compose up
  ```

- **Start services in detached mode**:
  ```bash
  docker-compose up -d
  ```

- **Stop services**:
  ```bash
  docker-compose down
  ```

- **Stop services and remove volumes**:
  ```bash
  docker-compose down -v
  ```

- **View logs**:
  ```bash
  docker-compose logs -f
  ```

- **Rebuild and restart**:
  ```bash
  docker-compose up --build
  ```

### Environment Variables

The following environment variables are configured in the Docker Compose file:

- **Backend**:
  - `DATABASE_URL`: PostgreSQL connection string
  - `JWT_SECRET`: Secret key for JWT token generation

- **Frontend**:
  - `NEXT_PUBLIC_API_URL`: Backend API URL

- **Database**:
  - `POSTGRES_USER`: Database username
  - `POSTGRES_PASSWORD`: Database password
  - `POSTGRES_DB`: Database name

### Troubleshooting

1. **Port conflicts**:
   If you encounter port conflicts, you can modify the port mappings in `docker-compose.yml`.

2. **Database initialization**:
   If you need to reset the database:
   ```bash
   docker-compose down -v
   docker-compose up --build
   ```

3. **Container logs**:
   To view logs for a specific service:
   ```bash
   docker-compose logs -f [service_name]
   ```

4. **Container shell**:
   To access a container's shell:
   ```bash
   docker-compose exec [service_name] sh
   ```

## Database Schema (PostgreSQL)

We use three main tables:

1. **`books`**: Stores information about each book.
2. **`lending_records`**: Tracks borrowing history.
3. **`users`**: Stores user authentication information.

### `books` Table

```sql
CREATE TABLE books (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    author VARCHAR(255) NOT NULL,
    isbn VARCHAR(20) UNIQUE NOT NULL,
    quantity INTEGER NOT NULL DEFAULT 0 CHECK (quantity >= 0),
    category VARCHAR(100),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
```

- `id`: Unique identifier for the book.
- `title`: Title of the book.
- `author`: Author of the book.
- `isbn`: International Standard Book Number (unique).
- `quantity`: Number of available copies.
- `category`: Genre or category of the book.
- `created_at`, `updated_at`: Timestamps for record creation and modification.

### `lending_records` Table

```sql
CREATE TABLE lending_records (
    id SERIAL PRIMARY KEY,
    book_id INTEGER NOT NULL REFERENCES books(id) ON DELETE CASCADE,
    borrower_name VARCHAR(255) NOT NULL,
    borrow_date DATE NOT NULL,
    return_date DATE NULL, -- Null means the book hasn't been returned yet
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
```

- `id`: Unique identifier for the lending record.
- `book_id`: Foreign key referencing the `books` table.
- `borrower_name`: Name of the person borrowing the book.
- `borrow_date`: Date the book was borrowed.
- `return_date`: Date the book was returned (NULL if still borrowed).
- `created_at`, `updated_at`: Timestamps for record creation and modification.

### `users` Table

```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
```

- `id`: Unique identifier for the user.
- `username`: Unique username for authentication.
- `password_hash`: Securely hashed password.
- `email`: User's email address (unique).
- `created_at`, `updated_at`: Timestamps for record creation and modification.

### Trigger for `updated_at`

```sql
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
   NEW.updated_at = now(); 
   RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_books_updated_at
BEFORE UPDATE ON books
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_lending_records_updated_at
BEFORE UPDATE ON lending_records
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_users_updated_at
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();
```

## API Documentation

The application provides comprehensive API documentation using Swagger/OpenAPI specification and serves it through two different interfaces:

1. **Swagger UI** (Backend): Available at `/api/apidocs/*`
   - Available at `http://localhost:3001/api/apidocs/` in development

2. **Redoc UI** (Frontend): Available at `/apidocs`
   - Available at `http://localhost:3000/apidocs` in development

### Generating API Documentation

The API documentation is automatically generated from annotations in the Go code. To update the documentation:

1. Install Swagger tools (if not already installed):
   ```bash
   cd backend
   go get github.com/swaggo/swag/cmd/swag@v1.8.12
   go get github.com/gofiber/swagger@v0.1.9
   go get github.com/swaggo/fiber-swagger@v1.3.0
   ```

2. Add Swagger annotations to your handlers:
   ```go
   // @Summary Get all books
   // @Description Get all books with optional search and filtering
   // @Tags books
   // @Accept json
   // @Produce json
   // @Success 200 {array} models.Book
   // @Router /books [get]
   func GetBooks(c *fiber.Ctx) error {
       // ... handler implementation
   }
   ```

3. Generate Swagger documentation:
   ```bash
   cd backend
   swag init
   ```

4. The documentation will be generated in the `backend/docs` directory.

### Accessing API Documentation

#### Development Environment:
- Swagger UI: `http://localhost:3001/swagger/`
- Raw OpenAPI Spec: `http://localhost:3000/api/apidocs`
- Redoc UI: `http://localhost:3000/apidocs`

#### Production Environment:
- Redoc UI: `https://digital-library-frontend.werdev.my.id/apidocs`

### Customizing API Documentation

1. **Backend (Swagger)**:
   - Update the Swagger configuration in `backend/main.go`:
     ```go
     // @title Digital Library API
     // @version 1.0
     // @description This is the API documentation for the Digital Library application
     // @host your-api-domain
     // @BasePath /api
     ```

2. **Frontend (Redoc)**:
   - Customize the theme in `frontend/app/apidocs/page.tsx`:
     ```typescript
     <RedocStandalone
       spec={spec}
       options={{
         theme: {
           colors: {
             primary: { main: '#3B82F6' }
           },
           // ... other theme options
         }
       }}
     />
     ```

## Setup Instructions

Follow these steps to set up and run the project locally.

### Prerequisites

- **Node.js** (v18 or later recommended) and npm/yarn
- **Go** (v1.18 or later recommended)
- **PostgreSQL** (running instance)
- **Git**

### 1. Clone the Repository

```bash
git clone https://github.com/agussuwerdo/digital-library.git
cd digital-library
```

### 2. Backend Setup (Golang API)

a. **Navigate to backend directory:**
   ```bash
   cd backend
   ```

b. **Create `.env` file:**
   Create a `.env` file in the `backend` directory and add the following environment variables:
   ```dotenv
   # Example .env file for backend
   DATABASE_URL="postgresql://YOUR_USER:YOUR_PASSWORD@YOUR_HOST:YOUR_PORT/YOUR_DB_NAME?sslmode=disable"
   JWT_SECRET="your_very_strong_and_secret_jwt_key_here_minimum_32_chars"
   ```
   *Example `DATABASE_URL` for a typical local setup: `DATABASE_URL="postgresql://postgres:postgres@localhost:5432/digital_library?sslmode=disable"`*

c. **Install dependencies:**
   ```bash
   go mod tidy
   ```

d. **Database Migration:**
   Create the database and execute the SQL commands for creating the tables (`books`, `lending_records`, `users`) and triggers as shown in the Database Schema section.

e. **Run the backend server:**
   ```bash
   go run main.go
   ```
   The backend API should now be running on `http://localhost:3000`.

### 3. Frontend Setup (Next.js App)

a. **Navigate to frontend directory:**
   From the project root directory (`digital-library`):
   ```bash
   cd frontend
   ```

b. **Install dependencies:**
   ```bash
   npm install
   # or if you prefer yarn: yarn install
   ```

c. **Environment Variables:**
   Create a `.env.local` file in the `frontend` directory:
   ```dotenv
   NEXT_PUBLIC_API_URL=http://localhost:3000/api 
   ```

d. **Run the frontend development server:**
   ```bash
   npm run dev
   # or yarn dev
   ```
   The frontend application should now be running on `http://localhost:3001`.

### 4. Accessing the Application

1. Open your browser and navigate to the frontend URL (e.g., `http://localhost:3001`).
2. You can either:
   - Register a new account by clicking the "Register" button and filling out the registration form
   - Log in with your existing credentials if you already have an account

The application provides a complete library management system with features for managing books, tracking lending records, and viewing analytics.

## Deployment

The application is configured for deployment on Vercel:
- Frontend: Deployed directly through Vercel's Next.js integration
- Backend: Deployed as Vercel Serverless Functions
- Database: Hosted on a PostgreSQL provider of your choice

For detailed deployment instructions, refer to the deployment documentation in the respective frontend and backend directories. 
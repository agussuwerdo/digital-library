# Digital Library Analytics Dashboard

A full-stack application for managing a digital library and visualizing analytics.

## Project Structure

- `backend/`: Golang (Fiber) REST API
- `frontend/`: Next.js (TypeScript, Tailwind) Web Application

## Database Schema (PostgreSQL)

We use two main tables:

1.  **`books`**: Stores information about each book.
2.  **`lending_records`**: Tracks borrowing history.

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
- `book_id`: Foreign key referencing the `books` table. If a book is deleted, corresponding lending records are also deleted (`ON DELETE CASCADE`).
- `borrower_name`: Name of the person borrowing the book.
- `borrow_date`: Date the book was borrowed.
- `return_date`: Date the book was returned (NULL if still borrowed).
- `created_at`, `updated_at`: Timestamps for record creation and modification.

### Trigger for `updated_at`

It's recommended to add a trigger function to automatically update the `updated_at` timestamp whenever a row is updated in either table.

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
```

## API Documentation

*(TODO: Integrate Swagger/OpenAPI documentation generation for the backend API. This typically involves adding annotations to the Go handlers and using a tool like `swag` (`github.com/swaggo/swag`) to generate the spec.)*

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
   Create a `.env` file in the `backend` directory and add the following environment variables. Replace the placeholders with your actual database credentials and a secure JWT secret.
   ```dotenv
   # Example .env file for backend
   DATABASE_URL="postgresql://YOUR_USER:YOUR_PASSWORD@YOUR_HOST:YOUR_PORT/YOUR_DB_NAME?sslmode=disable"
   JWT_SECRET="your_very_strong_and_secret_jwt_key_here_minimum_32_chars"
   ```
   *Example `DATABASE_URL` for a typical local setup: `DATABASE_URL="postgresql://postgres:postgres@localhost:5432/digital_library?sslmode=disable"`*

c. **Install dependencies:**
   If you haven't already during development:
   ```bash
   go mod tidy
   ```

d. **Database Migration:**
   Connect to your PostgreSQL instance using a tool like `psql` or a GUI client (e.g., DBeaver, pgAdmin). Create the database specified in your `DATABASE_URL` (e.g., `digital_library`). Then, execute the SQL commands found in the [Database Schema](#database-schema-postgresql) section of this README within that database to create the necessary tables (`books`, `lending_records`) and the `update_updated_at_column` trigger function.

e. **Run the backend server:**
   ```bash
   go run main.go
   ```
   The backend API should now be running, typically on `http://localhost:3000`.

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

c. **Environment Variables (Optional but Recommended):**
   The frontend defaults to connecting to `http://localhost:3000/api`. If your backend is running elsewhere, create a `.env.local` file in the `frontend` directory:
   ```dotenv
   # Example .env.local for frontend
   NEXT_PUBLIC_API_URL=http://localhost:3000/api 
   ```

d. **Run the frontend development server:**
   ```bash
   npm run dev
   # or yarn dev
   ```
   The frontend application should now be running, typically on `http://localhost:3001` (Next.js often defaults here if 3000 is taken).

### 4. Accessing the Application

- Open your browser and navigate to the frontend URL (e.g., `http://localhost:3001`).
- Use the hardcoded login credentials (defined in `backend/handlers/auth_handler.go`):
    - **Username:** `librarian`
    - **Password:** `password123`

*(End of file)* 
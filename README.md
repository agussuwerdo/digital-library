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

## Setup Instructions

*(TODO: Add setup instructions for backend, frontend, and database)* 
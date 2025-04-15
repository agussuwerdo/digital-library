-- Create users table if not exists
DO $$ 
BEGIN
    IF NOT EXISTS (SELECT FROM pg_tables WHERE tablename = 'users') THEN
        CREATE TABLE users (
            id SERIAL PRIMARY KEY,
            username VARCHAR(255) UNIQUE NOT NULL,
            password_hash VARCHAR(255) NOT NULL,
            email VARCHAR(255) UNIQUE NOT NULL,
            role VARCHAR(50) NOT NULL DEFAULT 'user',
            created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
        );
    END IF;
END $$;

-- Create books table if not exists
DO $$ 
BEGIN
    IF NOT EXISTS (SELECT FROM pg_tables WHERE tablename = 'books') THEN
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
    END IF;
END $$;

-- Create lending_records table if not exists
DO $$ 
BEGIN
    IF NOT EXISTS (SELECT FROM pg_tables WHERE tablename = 'lending_records') THEN
        CREATE TABLE lending_records (
            id SERIAL PRIMARY KEY,
            book_id INTEGER NOT NULL REFERENCES books(id) ON DELETE CASCADE,
            borrower_name VARCHAR(255) NOT NULL,
            borrow_date DATE NOT NULL,
            return_date DATE NULL,
            created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
        );
    END IF;
END $$;

-- Create trigger function if not exists
DO $$ 
BEGIN
    IF NOT EXISTS (SELECT FROM pg_proc WHERE proname = 'update_updated_at_column') THEN
        CREATE OR REPLACE FUNCTION update_updated_at_column()
        RETURNS TRIGGER AS $$
        BEGIN
            NEW.updated_at = now(); 
            RETURN NEW;
        END;
        $$ language 'plpgsql';
    END IF;
END $$;

-- Create triggers if they don't exist
DO $$ 
BEGIN
    IF NOT EXISTS (SELECT FROM pg_trigger WHERE tgname = 'update_users_updated_at') THEN
        CREATE TRIGGER update_users_updated_at
        BEFORE UPDATE ON users
        FOR EACH ROW
        EXECUTE FUNCTION update_updated_at_column();
    END IF;

    IF NOT EXISTS (SELECT FROM pg_trigger WHERE tgname = 'update_books_updated_at') THEN
        CREATE TRIGGER update_books_updated_at
        BEFORE UPDATE ON books
        FOR EACH ROW
        EXECUTE FUNCTION update_updated_at_column();
    END IF;

    IF NOT EXISTS (SELECT FROM pg_trigger WHERE tgname = 'update_lending_records_updated_at') THEN
        CREATE TRIGGER update_lending_records_updated_at
        BEFORE UPDATE ON lending_records
        FOR EACH ROW
        EXECUTE FUNCTION update_updated_at_column();
    END IF;
END $$; 
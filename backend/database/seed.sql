-- Insert sample users (password: 123)
INSERT INTO users (username, password_hash, email, role) VALUES
('admin', '$2a$10$B3wviry8gDQIUHwpCQKo7u4OguHWEO3TCk8i11S5UOA/4Y./uOi.a', 'admin@example.com', 'admin'),
('user', '$2a$10$B3wviry8gDQIUHwpCQKo7u4OguHWEO3TCk8i11S5UOA/4Y./uOi.a', 'user@example.com', 'user');

-- Insert sample books
INSERT INTO books (title, author, isbn, quantity, category) VALUES
('The Great Gatsby', 'F. Scott Fitzgerald', '9780743273565', 5, 'Classic'),
('To Kill a Mockingbird', 'Harper Lee', '9780446310789', 3, 'Classic'),
('1984', 'George Orwell', '9780451524935', 4, 'Science Fiction'),
('Pride and Prejudice', 'Jane Austen', '9780141439518', 2, 'Classic'),
('The Hobbit', 'J.R.R. Tolkien', '9780547928227', 6, 'Fantasy'),
('The Catcher in the Rye', 'J.D. Salinger', '9780316769488', 3, 'Classic'),
('Brave New World', 'Aldous Huxley', '9780060850524', 4, 'Science Fiction'),
('The Lord of the Rings', 'J.R.R. Tolkien', '9780618640157', 5, 'Fantasy'),
('Crime and Punishment', 'Fyodor Dostoevsky', '9780143107637', 2, 'Classic'),
('The Alchemist', 'Paulo Coelho', '9780062315007', 4, 'Fiction');

-- Insert sample lending records
INSERT INTO lending_records (book_id, borrower_name, borrow_date, return_date) VALUES
(1, 'John Doe', '2024-01-15', '2024-02-15'),
(2, 'Jane Smith', '2024-02-01', NULL),
(3, 'Bob Johnson', '2024-01-20', '2024-02-20'),
(4, 'Alice Brown', '2024-02-10', NULL),
(5, 'Charlie Wilson', '2024-01-25', '2024-02-25'); 
// Matches backend/models/models.go
export interface Book {
  id: number;
  title: string;
  author: string;
  isbn: string;
  quantity: number;
  category: string;
  created_at: string; // Use string for dates, can parse if needed
  updated_at: string;
}

// User type definition
export interface User {
  id: number;
  username: string;
  email: string;
  role: 'admin' | 'user';
  created_at: string;
  updated_at: string;
}

// Matches backend/handlers/lending_handler.go -> LendingRecordDetail
export interface LendingRecordDetail {
  id: number;
  book_id: number;
  borrower: string;
  borrow_date: string; 
  return_date?: string | null; // Optional/nullable
  created_at: string;
  updated_at: string;
  book_title: string; // Joined data
  book_author: string; // Joined data
}

// Matches backend/handlers/analytics_handler.go -> BorrowCount
export interface BorrowCount {
  book_id: number;
  book_title: string;
  borrows: number;
}

// Matches backend/handlers/analytics_handler.go -> MonthlyTrend
export interface MonthlyTrend {
  month: string; // YYYY-MM
  count: number;
}

// Matches backend/handlers/analytics_handler.go -> CategoryDistribution
export interface CategoryDistribution {
  category: string;
  count: number;
} 
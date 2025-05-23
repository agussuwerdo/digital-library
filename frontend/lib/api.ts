import { Book, LendingRecordDetail, BorrowCount, MonthlyTrend, CategoryDistribution } from '@/lib/types';

// Base URL for the backend API
const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 
  (process.env.NODE_ENV === 'production' 
    ? 'https://digital-library-backend.werdev.my.id/api' 
    : 'http://localhost:3001/api');

// Helper function to get JWT token from storage
const getToken = (): string | null => {
  // Check if running in browser environment
  if (typeof window !== 'undefined') {
    return localStorage.getItem('authToken');
  }
  return null;
};

// Helper function for making API requests
// Using unknown for the generic default and casting where needed
const apiRequest = async <T = unknown>(endpoint: string, options: RequestInit = {}): Promise<T> => {
  const token = getToken();
  // Use Record<string, string> for headers to allow arbitrary keys
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    'Accept': 'application/json',
    ...(options.headers as Record<string, string>),
  };

  if (token) {
    headers['Authorization'] = `Bearer ${token}`;
  }

  let response: Response;
  try {
    response = await fetch(`${API_BASE_URL}${endpoint}`, {
      ...options,
      headers,
      credentials: 'include',
      mode: 'cors',
    });
  } catch (networkError: unknown) {
    console.error('Network Error:', networkError);
    throw { status: 0, message: 'Network error, could not connect to API' };
  }

  let data: unknown = null; // Use unknown for data before parsing
  try {
    const contentType = response.headers.get('content-type');
    if (response.status === 204) {
      // No content, data remains null
    } else if (contentType && contentType.includes('application/json')) {
      data = await response.json();
    } else {
      // Handle unexpected content types (e.g., plain text error)
      const textData = await response.text();
      data = { error: textData || response.statusText };
    }
  } catch (parsingError: unknown) { // Catch JSON parsing errors
    console.error('API Response Parsing Error:', parsingError);
    throw { status: response.status, message: 'Failed to parse API response' };
  }

  if (!response.ok) {
    // Try to extract a meaningful error message from the parsed data
    const errorMessage = (data as { error?: string })?.error || response.statusText;
    console.error(`API Error ${response.status}:`, errorMessage, data);
    throw { status: response.status, message: errorMessage, data };
  }

  return data as T; // Cast to expected type T only on success
};

// --- Auth API --- 

interface LoginResponse {
  token: string;
  user: {
    id: number;
    username: string;
    email: string;
    role: 'admin' | 'user';
    created_at: string;
    updated_at: string;
  };
}

interface RegisterResponse {
  message: string;
}

export const login = async (credentials: { username: string; password: string }): Promise<LoginResponse> => {
  return apiRequest<LoginResponse>('/login', {
    method: 'POST',
    body: JSON.stringify(credentials),
  });
};

export const register = async (userData: { username: string; password: string; email: string }): Promise<RegisterResponse> => {
  const response = await fetch(`${API_BASE_URL}/register`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(userData),
  });

  if (!response.ok) {
    const error = await response.json();
    throw new Error(error.error || 'Registration failed');
  }

  return response.json();
};

// --- Book API (Placeholders) --- 

// Define an input type for book creation/update
// (Mirroring backend/models/Book but omitting server-generated fields)
type BookInput = Omit<Book, 'id' | 'created_at' | 'updated_at'>;

export const getBooks = async (params?: Record<string, string>): Promise<Book[]> => {
  const queryString = params ? new URLSearchParams(params).toString() : '';
  return apiRequest<Book[]>(`/books${queryString ? `?${queryString}` : ''}`); 
};

export const getBook = async (id: string | number): Promise<Book> => {
  return apiRequest<Book>(`/books/${id}`);
};

export const createBook = async (bookData: BookInput): Promise<Book> => {
  return apiRequest<Book>('/books', {
    method: 'POST',
    body: JSON.stringify(bookData),
  });
};

export const updateBook = async (id: string | number, bookData: Partial<BookInput>): Promise<Book> => {
  // Use Partial<BookInput> to allow updating only some fields
  return apiRequest<Book>(`/books/${id}`, {
    method: 'PUT',
    body: JSON.stringify(bookData),
  });
};

interface DeleteResponse {
  message: string;
  id: number;
}

export const deleteBook = async (id: string | number): Promise<DeleteResponse> => {
  return apiRequest<DeleteResponse>(`/books/${id}`, {
    method: 'DELETE',
  });
};

// --- Lending API (Placeholders) --- 

export const getLendingRecords = async (params?: Record<string, string>): Promise<LendingRecordDetail[]> => {
  const queryString = params ? new URLSearchParams(params).toString() : '';
  return apiRequest<LendingRecordDetail[]>(`/lending${queryString ? `?${queryString}` : ''}`);
};

interface LendBookPayload {
  book_id: number;
  borrower: string;
}

// The backend returns the created LendingRecord on successful lend
export const lendBook = async (payload: LendBookPayload): Promise<LendingRecordDetail> => {
  return apiRequest<LendingRecordDetail>('/lending/lend', {
    method: 'POST',
    body: JSON.stringify(payload),
  });
};

interface ReturnResponse {
  message: string;
}

export const returnBook = async (lendingRecordId: string | number): Promise<ReturnResponse> => {
  return apiRequest<ReturnResponse>(`/lending/return/${lendingRecordId}`, {
    method: 'POST',
  });
};

interface DeleteLendingResponse {
  message: string;
}

export const deleteLendingRecord = async (lendingRecordId: string | number): Promise<DeleteLendingResponse> => {
  return apiRequest<DeleteLendingResponse>(`/lending/${lendingRecordId}`, {
    method: 'DELETE',
  });
};

// --- Analytics API --- 

export const getMostBorrowed = async (username?: string, role?: string): Promise<BorrowCount[]> => {
  const queryParams = new URLSearchParams();
  if (username) queryParams.append('username', username);
  if (role) queryParams.append('role', role);
  const queryString = queryParams.toString();
  return apiRequest<BorrowCount[]>(`/analytics/most-borrowed${queryString ? `?${queryString}` : ''}`);
};

export const getMonthlyTrends = async (username?: string, role?: string): Promise<MonthlyTrend[]> => {
  const queryParams = new URLSearchParams();
  if (username) queryParams.append('username', username);
  if (role) queryParams.append('role', role);
  const queryString = queryParams.toString();
  return apiRequest<MonthlyTrend[]>(`/analytics/monthly-trends${queryString ? `?${queryString}` : ''}`);
};

export const getCategoryDistribution = async (username?: string, role?: string): Promise<CategoryDistribution[]> => {
  const queryParams = new URLSearchParams();
  if (username) queryParams.append('username', username);
  if (role) queryParams.append('role', role);
  const queryString = queryParams.toString();
  return apiRequest<CategoryDistribution[]>(`/analytics/category-distribution${queryString ? `?${queryString}` : ''}`);
}; 
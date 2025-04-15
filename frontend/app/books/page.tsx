'use client';

// FORCE APPLY MARKER
import { useState, useEffect, useCallback, Suspense } from 'react';
import * as api from '@/lib/api';
import { Book } from '@/lib/types';
import { useAuth } from '@/context/AuthContext'; // To ensure user is authenticated
import BookFormModal from '@/components/BookFormModal'; // Import the modal
import SearchFilter from '@/components/SearchFilter';
import { useSearchParams } from 'next/navigation';

function BooksContent() {
  const [books, setBooks] = useState<Book[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const { isAuthenticated, isLoading: authLoading, user } = useAuth(); // Get auth state and user info

  // State for modal
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [selectedBook, setSelectedBook] = useState<Book | null>(null); // Add state for selected book
  const searchParams = useSearchParams();

  const isAdmin = user?.role === 'admin'; // Check if user is admin

  const fetchBooks = useCallback(async (params: Record<string, string> = {}) => {
    if (!isAuthenticated) return;
  
    try {
      setLoading(true);
      const data = await api.getBooks(params);
      setBooks(data);
      setError(null);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'An error occurred');
    } finally {
      setLoading(false);
    }
  }, [isAuthenticated]);
  
  // Handle search from URL parameters
  useEffect(() => {
    if (!authLoading && isAuthenticated) {
      const params: Record<string, string> = {};
      searchParams.forEach((value, key) => {
        params[key] = value;
      });
      fetchBooks(params);
    } else if (!authLoading && !isAuthenticated) {
      // If auth is done loading and user is not authenticated, stop loading state
      // AuthProvider should handle the redirect already
      setLoading(false);
    }
  }, [authLoading, isAuthenticated, searchParams, fetchBooks]);

  // --- Modal Handlers --- 
  const handleOpenModal = (book: Book | null = null) => {
    setSelectedBook(book); // Set for editing or null for adding
    setIsModalOpen(true);
  };

  const handleCloseModal = () => {
    setIsModalOpen(false);
    setSelectedBook(null); // Clear selection on close
  };

  const handleSave = () => {
    handleCloseModal(); // Close modal after save
    fetchBooks(); // Refresh book list
  };
  // ---------------------

  // TODO: Implement delete functionality
  const handleDelete = async (id: number) => {
    if (window.confirm('Are you sure you want to delete this book? This action cannot be undone.')) {
        try {
          await api.deleteBook(id);
          fetchBooks(); 
        } catch (err: unknown) { // Use unknown type
          console.error("Failed to delete book:", err);
          // Extract error message safely
          let message = "Failed to delete book";
          if (typeof err === 'object' && err !== null && 'message' in err) {
            message = String((err as { message: unknown }).message);
          }
          setError(message); // Show error to user
        }
    }
  };

  const bookFilters = [
    {
      name: 'available',
      label: 'Availability',
      options: [
        { value: 'true', label: 'Available' },
        { value: 'false', label: 'Borrowed' },
      ],
    },
  ];

  // Display loading state
  if (authLoading || loading) {
    return <div className="p-8 text-center">Loading...</div>;
  }

  // If not authenticated after loading, AuthProvider handles redirect, show nothing here
  if (!isAuthenticated) {
    return null; 
  }

  // Display error state
  if (error) {
    return <div className="p-8 text-center text-red-500">Error loading books: {error}</div>;
  }

  // Display Book Management UI
  return (
    <div className="container mx-auto px-4 py-8">
      <h1 className="text-3xl font-bold mb-8">Books</h1>
      
      <SearchFilter
        filters={bookFilters}
        placeholder="Search by title or author..."
        className="mb-8"
      />

      <div className="flex justify-between items-center mb-6">
        <h1 className="text-3xl font-bold">Book Management</h1>
        {/* Only show Add Book button for admin users */}
        {isAdmin && (
          <button 
            onClick={() => handleOpenModal(null)}
            className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded"
          >
            Add New Book
          </button>
        )}
      </div>
      
      {/* Books Table */}
      <div className="shadow overflow-hidden border-b border-gray-200 sm:rounded-lg">
        <div className="overflow-x-auto">
          <table className="min-w-full divide-y divide-gray-200">
            <thead className="bg-gray-50">
              <tr>
                <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Title</th>
                <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Author</th>
                <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">ISBN</th>
                <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Category</th>
                <th scope="col" className="px-6 py-3 text-center text-xs font-medium text-gray-500 uppercase tracking-wider">Quantity</th>
                {/* Only show Actions column for admin users */}
                {isAdmin && (
                  <th scope="col" className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">Actions</th>
                )}
              </tr>
            </thead>
            <tbody className="bg-white divide-y divide-gray-200">
              {books.length === 0 ? (
                <tr>
                  <td colSpan={isAdmin ? 6 : 5} className="px-6 py-4 whitespace-nowrap text-sm text-gray-500 text-center">
                    No books found. {isAdmin ? 'Add one to get started!' : ''}
                  </td>
                </tr>
              ) : (
                books.map((book) => (
                  <tr key={book.id}>
                    <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">{book.title}</td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{book.author}</td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{book.isbn}</td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{book.category || 'N/A'}</td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500 text-center">{book.quantity}</td>
                    {/* Only show action buttons for admin users */}
                    {isAdmin && (
                      <td className="px-6 py-4 whitespace-nowrap text-right text-sm font-medium space-x-2">
                        <button 
                          onClick={() => handleOpenModal(book)}
                          className="text-indigo-600 hover:text-indigo-900"
                        >
                          Edit
                        </button>
                        <button 
                          onClick={() => handleDelete(book.id)} 
                          className="text-red-600 hover:text-red-900"
                        >
                          Delete
                        </button>
                      </td>
                    )}
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>
      </div>

      {/* Only render the modal for admin users */}
      {isAdmin && (
        <BookFormModal 
          isOpen={isModalOpen} 
          onClose={handleCloseModal} 
          bookToEdit={selectedBook} 
          onSave={handleSave} 
        />
      )}
    </div>
  );
}

export default function BooksPage() {
  return (
    <Suspense fallback={<div className="p-8 text-center">Loading...</div>}>
      <BooksContent />
    </Suspense>
  );
}

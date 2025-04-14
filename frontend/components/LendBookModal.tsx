'use client';

import { useState, useEffect, FormEvent } from 'react';
import * as api from '@/lib/api';
import { Book } from '@/lib/types';

interface LendBookModalProps {
  isOpen: boolean;
  onClose: () => void;
  onSave: () => void; // Callback to refresh lending records
}

export default function LendBookModal({ isOpen, onClose, onSave }: LendBookModalProps) {
  const [books, setBooks] = useState<Book[]>([]);
  const [selectedBookId, setSelectedBookId] = useState<string>(''); // Store book ID as string
  const [borrower, setBorrower] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [booksLoading, setBooksLoading] = useState(false);

  // Fetch available books when the modal opens
  useEffect(() => {
    if (isOpen) {
      setBooksLoading(true);
      setError(null);
      setBorrower(''); // Reset fields
      setSelectedBookId('');
      api.getBooks()
        .then(data => {
          // Filter books with quantity > 0 ? Maybe backend should handle this better
          // Or display quantity and let user know if unavailable
          setBooks(data.filter(b => b.quantity > 0)); // Only show books in stock
        })
        .catch(err => {
          console.error("Failed to fetch books for lending modal:", err);
          setError("Could not load available books.");
        })
        .finally(() => {
          setBooksLoading(false);
        });
    }
  }, [isOpen]);

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    if (!selectedBookId || !borrower) {
      setError("Please select a book and enter borrower name.");
      return;
    }
    setLoading(true);
    setError(null);

    try {
      await api.lendBook({ book_id: parseInt(selectedBookId, 10), borrower });
      onSave(); // Refresh lending list
      onClose(); // Close modal
    } catch (err: unknown) {
      console.error("Failed to lend book:", err);
      let message = "Failed to lend book";
      if (typeof err === 'object' && err !== null && 'message' in err) {
        message = String((err as { message: unknown }).message);
      }
      // Specific check for 'Book is currently out of stock' might be useful
      if (message.includes("out of stock")) {
          setError("Selected book is out of stock. Please refresh or choose another.");
      } else {
          setError(message);
      }
    } finally {
      setLoading(false);
    }
  };

  if (!isOpen) {
    return null;
  }

  return (
    <div className="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50 flex items-center justify-center">
      <div className="relative mx-auto p-5 border w-full max-w-md shadow-lg rounded-md bg-white">
        <h3 className="text-lg font-medium leading-6 text-gray-900 mb-4">
          Lend a Book
        </h3>
        <form onSubmit={handleSubmit}>
          <div className="mb-4">
            <label htmlFor="book" className="block text-sm font-medium text-gray-700">Book</label>
            {booksLoading ? (
              <p className="text-sm text-gray-500">Loading books...</p>
            ) : (
              <select
                id="book"
                name="book"
                value={selectedBookId}
                onChange={(e) => setSelectedBookId(e.target.value)}
                required
                className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 disabled:bg-gray-100"
                disabled={loading}
              >
                <option value="" disabled>Select a book</option>
                {books.map(book => (
                  // Disable option if quantity is 0, though we filter initially
                  <option key={book.id} value={book.id} disabled={book.quantity <= 0}>
                    {book.title} by {book.author} (Qty: {book.quantity})
                  </option>
                ))}
              </select>
            )}
             {books.length === 0 && !booksLoading && <p className="text-sm text-red-500 mt-1">No books currently available to lend.</p>}
          </div>

          <div className="mb-4">
            <label htmlFor="borrower" className="block text-sm font-medium text-gray-700">Borrower Name</label>
            <input 
              type="text" 
              name="borrower" 
              id="borrower" 
              value={borrower} 
              onChange={(e) => setBorrower(e.target.value)} 
              required 
              className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 disabled:bg-gray-100"
              disabled={loading}
            />
          </div>

          {error && (
            <div className="mb-4 p-2 text-sm text-red-700 bg-red-100 rounded-md border border-red-300">
              Error: {error}
            </div>
          )}

          <div className="items-center px-4 py-3 flex justify-end space-x-2">
            <button
              type="button"
              onClick={onClose}
              className="px-4 py-2 bg-gray-200 text-gray-800 rounded hover:bg-gray-300 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-gray-500"
              disabled={loading}
            >
              Cancel
            </button>
            <button
              type="submit"
              className="px-4 py-2 bg-green-600 text-white rounded hover:bg-green-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-green-500 disabled:opacity-50"
              disabled={loading || booksLoading || books.length === 0}
            >
              {loading ? 'Lending...' : 'Lend Book'}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
} 
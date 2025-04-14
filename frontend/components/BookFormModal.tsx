'use client';

import { useState, useEffect, FormEvent } from 'react';
import { Book } from '@/lib/types';
import * as api from '@/lib/api';

interface BookFormModalProps {
  isOpen: boolean;
  onClose: () => void;
  bookToEdit: Book | null; // Null when adding a new book
  onSave: () => void; // Callback to refresh the book list
}

// Define the structure for form data (excluding server-generated fields)
type BookFormData = Omit<Book, 'id' | 'created_at' | 'updated_at'>;

const initialFormData: BookFormData = {
  title: '',
  author: '',
  isbn: '',
  quantity: 0,
  category: '',
};

export default function BookFormModal({ isOpen, onClose, bookToEdit, onSave }: BookFormModalProps) {
  const [formData, setFormData] = useState<BookFormData>(initialFormData);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  // Reset form when bookToEdit changes (or modal opens for adding)
  useEffect(() => {
    if (isOpen) {
      if (bookToEdit) {
        setFormData({
          title: bookToEdit.title,
          author: bookToEdit.author,
          isbn: bookToEdit.isbn,
          quantity: bookToEdit.quantity,
          category: bookToEdit.category || '',
        });
      } else {
        setFormData(initialFormData); // Reset for new book
      }
      setError(null); // Clear previous errors
    }
  }, [isOpen, bookToEdit]);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>) => {
    const { name, value } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: name === 'quantity' ? parseInt(value, 10) || 0 : value, // Ensure quantity is number
    }));
  };

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError(null);

    try {
      if (bookToEdit) {
        // Update existing book
        await api.updateBook(bookToEdit.id, formData);
      } else {
        // Create new book
        await api.createBook(formData);
      }
      onSave(); // Refresh the book list in the parent component
      onClose(); // Close the modal
    } catch (err: unknown) {
      console.error("Failed to save book:", err);
      let message = "Failed to save book";
      if (typeof err === 'object' && err !== null && 'message' in err) {
        message = String((err as { message: unknown }).message);
      }
      setError(message);
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
          {bookToEdit ? 'Edit Book' : 'Add New Book'}
        </h3>
        <form onSubmit={handleSubmit}>
          {/* Form Fields */} 
          <div className="mb-4">
            <label htmlFor="title" className="block text-sm font-medium text-gray-700">Title</label>
            <input type="text" name="title" id="title" value={formData.title} onChange={handleChange} required className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500" />
          </div>
          <div className="mb-4">
            <label htmlFor="author" className="block text-sm font-medium text-gray-700">Author</label>
            <input type="text" name="author" id="author" value={formData.author} onChange={handleChange} required className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500" />
          </div>
          <div className="mb-4">
            <label htmlFor="isbn" className="block text-sm font-medium text-gray-700">ISBN</label>
            <input type="text" name="isbn" id="isbn" value={formData.isbn} onChange={handleChange} required className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500" />
          </div>
          <div className="mb-4">
            <label htmlFor="quantity" className="block text-sm font-medium text-gray-700">Quantity</label>
            <input type="number" name="quantity" id="quantity" value={formData.quantity} onChange={handleChange} required min="0" className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500" />
          </div>
          <div className="mb-4">
            <label htmlFor="category" className="block text-sm font-medium text-gray-700">Category</label>
            <input type="text" name="category" id="category" value={formData.category} onChange={handleChange} className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500" />
          </div>

          {error && (
            <div className="mb-4 p-2 text-sm text-red-700 bg-red-100 rounded-md border border-red-300">
              Error: {error}
            </div>
          )}

          {/* Action Buttons */} 
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
              className="px-4 py-2 bg-indigo-600 text-white rounded hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:opacity-50"
              disabled={loading}
            >
              {loading ? 'Saving...' : (bookToEdit ? 'Update Book' : 'Add Book')}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
} 
"use client";

import { useState, useEffect, useCallback, Suspense } from "react";
import * as api from "@/lib/api";
import { LendingRecordDetail } from "@/lib/types";
import { useAuth } from "@/context/AuthContext";
import { format } from "date-fns"; // For formatting dates
import LendBookModal from "@/components/LendBookModal"; // Import the modal
import SearchFilter from "@/components/SearchFilter";
import { useSearchParams } from "next/navigation";

function LendingContent() {
  const [records, setRecords] = useState<LendingRecordDetail[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const { isAuthenticated, isLoading: authLoading, user } = useAuth();

  // State for Lend Book Modal
  const [isLendModalOpen, setIsLendModalOpen] = useState(false);
  const searchParams = useSearchParams();

  const fetchLendingRecords = useCallback(
    async (params: Record<string, string> = {}) => {
      if (!isAuthenticated) return;

      try {
        setLoading(true);
        const data = await api.getLendingRecords(params);
        // If user is not admin, filter to show only their records
        if (user?.role !== "admin") {
          setRecords(
            data.filter((record) => record.borrower === user?.username)
          );
        } else {
          setRecords(data);
        }
        setError(null);
      } catch (err) {
        setError(err instanceof Error ? err.message : "An error occurred");
      } finally {
        setLoading(false);
      }
    },
    [isAuthenticated, user]
  );

  // Handle search from URL parameters
  useEffect(() => {
    if (!authLoading && isAuthenticated) {
      const params: Record<string, string> = {};
      searchParams.forEach((value, key) => {
        params[key] = value;
      });
      fetchLendingRecords(params);
    } else if (!authLoading && !isAuthenticated) {
      setLoading(false);
    }
  }, [authLoading, isAuthenticated, searchParams, fetchLendingRecords]);

  const handleReturn = async (recordId: number) => {
    if (window.confirm("Mark this book as returned?")) {
      try {
        await api.returnBook(recordId);
        fetchLendingRecords(); // Refresh list
      } catch (err: unknown) {
        console.error("Failed to return book:", err);
        let message = "Failed to return book";
        if (typeof err === "object" && err !== null && "message" in err) {
          message = String((err as { message: unknown }).message);
        }
        setError(message);
      }
    }
  };

  const handleDelete = async (recordId: number) => {
    if (user?.role !== "admin") return; // Only admins can delete

    if (
      window.confirm(
        "Are you sure you want to delete this lending record? This may affect analytics."
      )
    ) {
      try {
        await api.deleteLendingRecord(recordId);
        fetchLendingRecords(); // Refresh list
      } catch (err: unknown) {
        console.error("Failed to delete lending record:", err);
        let message = "Failed to delete record";
        if (typeof err === "object" && err !== null && "message" in err) {
          message = String((err as { message: unknown }).message);
        }
        setError(message);
      }
    }
  };

  // Open lend modal handler
  const handleOpenLendModal = () => {
    setIsLendModalOpen(true);
  };

  // Close lend modal handler
  const handleCloseLendModal = () => {
    setIsLendModalOpen(false);
  };

  // Callback for successful save from lend modal
  const handleLendSave = () => {
    handleCloseLendModal();
    fetchLendingRecords(); // Refresh the list
  };

  const lendingFilters = [
    {
      name: "status",
      label: "Status",
      options: [
        { value: "active", label: "Active" },
        { value: "returned", label: "Returned" },
      ],
    },
  ];

  if (authLoading || loading) {
    return <div className="p-8 text-center">Loading...</div>;
  }
  if (!isAuthenticated) {
    return null;
  }
  if (error) {
    return <div className="p-8 text-center text-red-500">Error: {error}</div>;
  }

  const formatDate = (dateString: string | null | undefined) => {
    if (!dateString) return "N/A";
    try {
      // Assuming date string is in a format date-fns can parse, like ISO 8601 or YYYY-MM-DD
      return format(new Date(dateString), "yyyy-MM-dd");
    } catch {
      return "Invalid Date";
    }
  };

  return (
    <div className="container mx-auto px-4 py-8">
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-3xl font-bold">Lending Management</h1>
        <button
          onClick={handleOpenLendModal}
          className="bg-green-500 hover:bg-green-700 text-white font-bold py-2 px-4 rounded"
        >
          {user?.role === 'admin' ? 'Lend New Book' : 'Borrow New Book'}
        </button>
      </div>

      <SearchFilter
        filters={lendingFilters}
        placeholder="Search by book title or borrower..."
        className="mb-8"
      />

      {/* Lending Records Table */}
      <div className="shadow overflow-hidden border-b border-gray-200 sm:rounded-lg">
        <div className="overflow-x-auto">
          <table className="min-w-full divide-y divide-gray-200">
            <thead className="bg-gray-50">
              <tr>
                <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Book Title</th>
                <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Author</th>
                <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Borrower</th>
                <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Borrow Date</th>
                <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Return Date</th>
                <th scope="col" className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">Actions</th>
              </tr>
            </thead>
            <tbody className="bg-white divide-y divide-gray-200">
              {records.length === 0 ? (
                <tr>
                  <td colSpan={6} className="px-6 py-4 whitespace-nowrap text-sm text-gray-500 text-center">No lending records found.</td>
                </tr>
              ) : (
                records.map((record) => (
                  <tr key={record.id}>
                    <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">{record.book_title}</td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{record.book_author}</td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{record.borrower}</td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{formatDate(record.borrow_date)}</td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                      {record.return_date ? formatDate(record.return_date) : 'Not Returned'}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-right text-sm font-medium space-x-2">
                      {!record.return_date && (
                        <button
                          onClick={() => handleReturn(record.id)}
                          className="text-green-600 hover:text-green-900"
                        >
                          Return
                        </button>
                      )}
                      {user?.role === "admin" && (
                        <button
                          onClick={() => handleDelete(record.id)}
                          className="text-red-600 hover:text-red-900"
                        >
                          Delete
                        </button>
                      )}
                    </td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>
      </div>

      {/* Render the Lend Book Modal */}
      <LendBookModal
        isOpen={isLendModalOpen}
        onClose={handleCloseLendModal}
        onSave={handleLendSave}
      />
    </div>
  );
}

export default function LendingPage() {
  return (
    <Suspense fallback={<div className="p-8 text-center">Loading...</div>}>
      <LendingContent />
    </Suspense>
  );
}

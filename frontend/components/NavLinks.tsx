'use client';

import Link from 'next/link';
import { useAuth } from '@/context/AuthContext';

export default function NavLinks() {
  const { isAuthenticated, logout } = useAuth();

  return (
    <div className="space-x-4">
      {isAuthenticated ? (
        <>
          <Link href="/dashboard" className="hover:text-gray-300">Dashboard</Link>
          <Link href="/books" className="hover:text-gray-300">Books</Link>
          <Link href="/lending" className="hover:text-gray-300">Lending</Link>
          <button onClick={logout} className="hover:text-gray-300 bg-transparent border-none text-white cursor-pointer">
            Logout
          </button>
        </>
      ) : (
        <Link href="/login" className="hover:text-gray-300">Login</Link>
      )}
    </div>
  );
} 
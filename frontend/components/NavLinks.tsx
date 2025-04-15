'use client';

import Link from 'next/link';
import { useAuth } from '@/context/AuthContext';
import { usePathname } from 'next/navigation';
import { useState, useRef, useEffect } from 'react';

const adminLinks = [
  { name: 'Dashboard', href: '/dashboard' },
  { name: 'Books', href: '/books' },
  { name: 'Lending', href: '/lending' },
  { name: 'API Docs', href: '/apidocs' },
];

const userLinks = [
  { name: 'Lending', href: '/lending' },
];

export default function NavLinks() {
  const { isAuthenticated, user, logout } = useAuth();
  const pathname = usePathname();
  const [showProfile, setShowProfile] = useState(false);
  const profileRef = useRef<HTMLDivElement>(null);

  // Reset profile state when authentication status changes
  useEffect(() => {
    setShowProfile(false);
  }, [isAuthenticated]);

  // Handle click outside to close the profile panel
  useEffect(() => {
    function handleClickOutside(event: MouseEvent) {
      if (profileRef.current && !profileRef.current.contains(event.target as Node)) {
        setShowProfile(false);
      }
    }

    document.addEventListener('mousedown', handleClickOutside);
    return () => {
      document.removeEventListener('mousedown', handleClickOutside);
    };
  }, []);

  if (!isAuthenticated) {
    return null;
  }

  const links = user?.role === 'admin' ? adminLinks : userLinks;

  return (
    <div className="flex items-center gap-6">
      {links.map((link) => {
        const isActive = pathname === link.href;
        
        return (
          <Link
            key={link.name}
            href={link.href}
            className={`text-sm font-medium transition-colors hover:text-gray-300 ${
              isActive ? 'text-white' : 'text-gray-300'
            }`}
          >
            {link.name}
          </Link>
        );
      })}
      
      <div className="relative" ref={profileRef}>
        <button 
          onClick={() => setShowProfile(!showProfile)}
          className="flex items-center gap-2 text-sm font-medium hover:text-gray-300 text-gray-300 cursor-pointer"
        >
          <span className="w-8 h-8 rounded-full bg-blue-600 flex items-center justify-center">
            {user?.username.charAt(0).toUpperCase()}
          </span>
          <span>{user?.username}</span>
        </button>
        
        {showProfile && (
          <div className="absolute right-0 mt-2 w-48 bg-white rounded-md shadow-lg py-1 z-50">
            <div className="px-4 py-2 border-b">
              <p className="text-sm font-medium text-gray-900">{user?.username}</p>
              <p className="text-xs text-gray-500">{user?.email}</p>
              <p className="text-xs text-gray-500 capitalize">{user?.role}</p>
            </div>
            <button
              onClick={logout}
              className="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 cursor-pointer"
            >
              Logout
            </button>
          </div>
        )}
      </div>
    </div>
  );
} 
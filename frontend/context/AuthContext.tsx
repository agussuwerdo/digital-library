'use client';

import React, { createContext, useState, useContext, useEffect, ReactNode } from 'react';
import { useRouter, usePathname } from 'next/navigation';

interface AuthContextType {
  isAuthenticated: boolean;
  isLoading: boolean;
  login: (token: string) => void;
  logout: () => void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider = ({ children }: { children: ReactNode }) => {
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const [isLoading, setIsLoading] = useState(true); // Start loading until check is done
  const router = useRouter();
  const pathname = usePathname();

  useEffect(() => {
    // Check for token on initial load
    const token = localStorage.getItem('authToken');
    if (token) {
      // TODO: Optionally validate token with backend endpoint here
      setIsAuthenticated(true);
    } else {
      setIsAuthenticated(false);
    }
    setIsLoading(false);
  }, []);

  useEffect(() => {
    // Redirect unauthenticated users from protected pages
    if (!isLoading && !isAuthenticated) {
      const publicPaths = ['/login', '/', '/register'];
      if (!publicPaths.includes(pathname)) {
        router.push('/login');
      }
    }
    // Optional: Redirect authenticated users away from login page
    // if (!isLoading && isAuthenticated && pathname === '/login') {
    //   router.push('/dashboard');
    // }
  }, [isLoading, isAuthenticated, pathname, router]);

  const login = (token: string) => {
    localStorage.setItem('authToken', token);
    setIsAuthenticated(true);
    router.push('/dashboard'); // Redirect after login
  };

  const logout = () => {
    localStorage.removeItem('authToken');
    setIsAuthenticated(false);
    router.push('/login'); // Redirect to login after logout
  };

  return (
    <AuthContext.Provider value={{ isAuthenticated, isLoading, login, logout }}>
      {children}
    </AuthContext.Provider>
  );
};

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
}; 
'use client';

import { useState, useEffect, useRef } from 'react';
import { useSearchParams, useRouter } from 'next/navigation';

interface SearchFilterProps {
  filters: {
    name: string;
    label: string;
    options: { value: string; label: string }[];
  }[];
  placeholder?: string;
  className?: string;
}

export default function SearchFilter({ filters, placeholder = 'Search...', className = '' }: SearchFilterProps) {
  const searchParams = useSearchParams();
  const router = useRouter();
  const searchInputRef = useRef<HTMLInputElement>(null);
  const [searchTerm, setSearchTerm] = useState(searchParams.get('search') || '');
  const [selectedFilters, setSelectedFilters] = useState<Record<string, string>>(() => {
    const initialFilters: Record<string, string> = {};
    filters.forEach(filter => {
      const value = searchParams.get(filter.name);
      if (value) initialFilters[filter.name] = value;
    });
    return initialFilters;
  });
  const debounceTimeout = useRef<NodeJS.Timeout | null>(null);

  const updateUrlParams = (term: string, filters: Record<string, string>) => {
    const params = new URLSearchParams();
    if (term) params.set('search', term);
    Object.entries(filters).forEach(([key, value]) => {
      if (value) params.set(key, value);
    });
    router.replace(`?${params.toString()}`);
  };

  const handleSearchChange = (term: string) => {
    setSearchTerm(term);
    if (debounceTimeout.current) {
      clearTimeout(debounceTimeout.current);
    }
    debounceTimeout.current = setTimeout(() => {
      updateUrlParams(term, selectedFilters);
    }, 1500);
  };

  const handleFilterChange = (name: string, value: string) => {
    const newFilters = {
      ...selectedFilters,
      [name]: value,
    };
    setSelectedFilters(newFilters);
    updateUrlParams(searchTerm, newFilters);
  };

  // Cleanup timeout on unmount
  useEffect(() => {
    return () => {
      if (debounceTimeout.current) {
        clearTimeout(debounceTimeout.current);
      }
    };
  }, []);

  // Maintain input focus after URL updates
  useEffect(() => {
    if (searchInputRef.current) {
      searchInputRef.current.focus();
    }
  }, [searchParams]);

  return (
    <div className={`space-y-4 ${className}`}>
      <div className="relative">
        <input
          ref={searchInputRef}
          type="text"
          placeholder={placeholder}
          value={searchTerm}
          onChange={(e) => handleSearchChange(e.target.value)}
          className="w-full px-4 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
        />
        {searchTerm && (
          <button
            onClick={() => handleSearchChange('')}
            className="absolute right-3 top-1/2 transform -translate-y-1/2 text-gray-400 hover:text-gray-600"
          >
            ×
          </button>
        )}
      </div>

      <div className="flex flex-wrap gap-4">
        {filters.map((filter) => (
          <div key={filter.name} className="flex-1 min-w-[200px]">
            <label className="block text-sm font-medium text-gray-700 mb-1">
              {filter.label}
            </label>
            <select
              value={selectedFilters[filter.name] || ''}
              onChange={(e) => handleFilterChange(filter.name, e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            >
              <option value="">All</option>
              {filter.options.map((option) => (
                <option key={option.value} value={option.value}>
                  {option.label}
                </option>
              ))}
            </select>
          </div>
        ))}
      </div>
    </div>
  );
} 
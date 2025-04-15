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

  const updateUrlParams = (term: string, filters: Record<string, string>) => {
    const params = new URLSearchParams();
    if (term) params.set('search', term);
    Object.entries(filters).forEach(([key, value]) => {
      if (value) params.set(key, value);
    });
    router.replace(`?${params.toString()}`);
  };

  const handleSearch = () => {
    updateUrlParams(searchTerm, selectedFilters);
  };

  const handleKeyPress = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === 'Enter') {
      handleSearch();
    }
  };

  const handleFilterChange = (name: string, value: string) => {
    const newFilters = {
      ...selectedFilters,
      [name]: value,
    };
    setSelectedFilters(newFilters);
    updateUrlParams(searchTerm, newFilters);
  };

  // Maintain input focus after URL updates
  useEffect(() => {
    if (searchInputRef.current) {
      searchInputRef.current.focus();
    }
  }, [searchParams]);

  return (
    <div className={`space-y-4 ${className}`}>
      <div className="relative flex">
        <input
          ref={searchInputRef}
          type="text"
          placeholder={placeholder}
          value={searchTerm}
          onChange={(e) => setSearchTerm(e.target.value)}
          onKeyPress={handleKeyPress}
          className="flex-1 px-4 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
        />
        <button
          onClick={handleSearch}
          className="ml-2 px-4 py-2 bg-blue-500 text-white rounded-md hover:bg-blue-600 focus:outline-none focus:ring-2 focus:ring-blue-500"
        >
          Search
        </button>
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
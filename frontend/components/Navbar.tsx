'use client';

import Link from 'next/link';
import NavLinks from './NavLinks';

export default function Navbar() {
  return (
    <nav className="bg-blue-900 text-white p-4">
      <div className="container mx-auto flex justify-between items-center">
        <Link href="/" className="text-xl font-bold">Digital Library</Link>
        <NavLinks />
      </div>
    </nav>
  );
} 
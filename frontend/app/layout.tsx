import type { Metadata } from "next";
import { Inter } from "next/font/google";
import "./globals.css";
import Link from "next/link";

const inter = Inter({ subsets: ["latin"] });

export const metadata: Metadata = {
  title: "Digital Library",
  description: "Digital Library Analytics Dashboard",
};

function Navbar() {
  return (
    <nav className="bg-gray-800 text-white p-4 shadow-md">
      <div className="container mx-auto flex justify-between items-center">
        <Link href="/" className="text-xl font-bold">
          Digital Library
        </Link>
        <div className="space-x-4">
          <Link href="/dashboard" className="hover:text-gray-300">Dashboard</Link>
          <Link href="/books" className="hover:text-gray-300">Books</Link>
          <Link href="/lending" className="hover:text-gray-300">Lending</Link>
          <Link href="/login" className="hover:text-gray-300">Login</Link>
        </div>
      </div>
    </nav>
  );
}

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body className={inter.className}>
        <Navbar />
        <main>{children}</main>
      </body>
    </html>
  );
}

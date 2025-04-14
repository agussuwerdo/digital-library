import Link from "next/link";

export default function Home() {
  return (
    <main className="flex min-h-screen flex-col items-center justify-center p-24">
      <h1 className="text-4xl font-bold mb-8">Digital Library</h1>
      <p className="mb-4">Welcome to the Digital Library App.</p>
      <div className="flex space-x-4">
        <Link href="/dashboard" className="text-blue-500 hover:underline">
          Go to Dashboard
        </Link>
      </div>
    </main>
  );
}

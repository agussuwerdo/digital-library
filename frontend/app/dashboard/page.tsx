export default function DashboardPage() {
  return (
    <div className="p-8">
      <h1 className="text-3xl font-bold mb-6">Dashboard</h1>
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {/* Placeholder for charts */}
        <div className="border p-4 rounded shadow">Most Borrowed Books Chart</div>
        <div className="border p-4 rounded shadow">Monthly Lending Trends Chart</div>
        <div className="border p-4 rounded shadow">Category Distribution Chart</div>
      </div>
      {/* TODO: Implement Chart components and fetch data */}
    </div>
  );
}

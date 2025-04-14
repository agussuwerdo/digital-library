'use client';

import { useState, useEffect } from 'react';
import * as api from '@/lib/api';
import { useAuth } from '@/context/AuthContext';
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  BarElement,
  LineElement,
  PointElement,
  ArcElement,
  Title,
  Tooltip,
  Legend,
} from 'chart.js';
import { Bar, Line, Pie } from 'react-chartjs-2';
import { BorrowCount, MonthlyTrend, CategoryDistribution } from '@/lib/types';

// Register Chart.js components
ChartJS.register(
  CategoryScale,
  LinearScale,
  BarElement,
  LineElement,
  PointElement,
  ArcElement,
  Title,
  Tooltip,
  Legend
);

export default function DashboardPage() {
  const { isAuthenticated, isLoading: authLoading } = useAuth();
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const [mostBorrowed, setMostBorrowed] = useState<BorrowCount[]>([]);
  const [monthlyTrends, setMonthlyTrends] = useState<MonthlyTrend[]>([]);
  const [categoryDistribution, setCategoryDistribution] = useState<CategoryDistribution[]>([]);

  useEffect(() => {
    if (!authLoading && isAuthenticated) {
      setLoading(true);
      Promise.all([
        api.getMostBorrowed(),
        api.getMonthlyTrends(),
        api.getCategoryDistribution(),
      ]).then(([borrowedData, trendsData, categoryData]) => {
        setMostBorrowed(borrowedData);
        setMonthlyTrends(trendsData);
        setCategoryDistribution(categoryData);
        setError(null);
      }).catch(err => {
        console.error("Failed to load dashboard data:", err);
        setError(err.message || "Failed to load dashboard data");
      }).finally(() => {
        setLoading(false);
      });
    } else if (!authLoading && !isAuthenticated) {
      setLoading(false); // Stop loading if not authenticated
    }
  }, [authLoading, isAuthenticated]);

  if (authLoading || loading) {
    return <div className="p-8 text-center">Loading dashboard...</div>;
  }
  if (!isAuthenticated) {
    return null; 
  }
  if (error) {
    return <div className="p-8 text-center text-red-500">Error: {error}</div>;
  }

  // Chart Data Preparation
  const mostBorrowedChartData = {
    labels: mostBorrowed.map(item => item.book_title),
    datasets: [{
      label: 'Borrow Count',
      data: mostBorrowed.map(item => item.borrows),
      backgroundColor: 'rgba(54, 162, 235, 0.6)',
      borderColor: 'rgba(54, 162, 235, 1)',
      borderWidth: 1,
    }],
  };

  const monthlyTrendsChartData = {
    labels: monthlyTrends.map(item => item.month),
    datasets: [{
      label: 'Lending Count',
      data: monthlyTrends.map(item => item.count),
      fill: false,
      borderColor: 'rgb(75, 192, 192)',
      tension: 0.1,
    }],
  };

  const categoryDistributionChartData = {
    labels: categoryDistribution.map(item => item.category),
    datasets: [{
      label: 'Book Count by Category',
      data: categoryDistribution.map(item => item.count),
      backgroundColor: [
        'rgba(255, 99, 132, 0.6)',
        'rgba(54, 162, 235, 0.6)',
        'rgba(255, 206, 86, 0.6)',
        'rgba(75, 192, 192, 0.6)',
        'rgba(153, 102, 255, 0.6)',
        'rgba(255, 159, 64, 0.6)',
      ],
      borderColor: [
        'rgba(255, 99, 132, 1)',
        'rgba(54, 162, 235, 1)',
        'rgba(255, 206, 86, 1)',
        'rgba(75, 192, 192, 1)',
        'rgba(153, 102, 255, 1)',
        'rgba(255, 159, 64, 1)',
      ],
      borderWidth: 1,
    }],
  };

  const chartOptions = {
    responsive: true,
    maintainAspectRatio: false,
    plugins: {
      legend: { position: 'top' as const },
      title: { display: true, text: 'Chart Title' }, // Set specific titles below
    },
  };

  return (
    <div className="p-8">
      <h1 className="text-3xl font-bold mb-6">Analytics Dashboard</h1>
      <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
        
        {/* Most Borrowed Books Chart */} 
        <div className="bg-white p-4 rounded shadow" style={{ height: '400px' }}>
          <h2 className="text-xl font-semibold mb-4 text-center">Most Borrowed Books</h2>
          <div style={{ height: 'calc(100% - 40px)' }}> {/* Adjust height accounting for title */}
            <Bar 
              options={ { ...chartOptions, plugins: { ...chartOptions.plugins, title: { display: true, text: 'Top 10 Most Borrowed Books'} } } }
              data={mostBorrowedChartData} 
            />
          </div>
        </div>

        {/* Monthly Lending Trends Chart */} 
        <div className="bg-white p-4 rounded shadow" style={{ height: '400px' }}>
          <h2 className="text-xl font-semibold mb-4 text-center">Monthly Lending Trends</h2>
          <div style={{ height: 'calc(100% - 40px)' }}>
            <Line 
              options={ { ...chartOptions, plugins: { ...chartOptions.plugins, title: { display: true, text: 'Lends per Month'} } } }
              data={monthlyTrendsChartData} 
            />
          </div>
        </div>

        {/* Books by Category Chart */} 
        <div className="bg-white p-4 rounded shadow md:col-span-2" style={{ height: '400px' }}>
          <h2 className="text-xl font-semibold mb-4 text-center">Books by Category</h2>
           <div style={{ height: 'calc(100% - 40px)', width:'50%', margin: 'auto' }}> {/* Center pie chart */}
             <Pie 
               options={ { ...chartOptions, plugins: { ...chartOptions.plugins, title: { display: true, text: 'Book Distribution by Category'} } } }
               data={categoryDistributionChartData} 
             />
          </div>
        </div>

      </div>
    </div>
  );
}

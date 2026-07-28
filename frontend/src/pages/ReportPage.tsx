import React, { useState } from 'react';
import {
  BarChart3, Download, Calendar, TrendingUp, TrendingDown,
  DollarSign, ShoppingCart, Users, FileText, Filter, Package
} from 'lucide-react';
import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer, LineChart, Line } from 'recharts';

const weeklyData = [
  { day: 'Sen', sales: 2100000, transactions: 42 },
  { day: 'Sel', sales: 1850000, transactions: 38 },
  { day: 'Rab', sales: 2400000, transactions: 51 },
  { day: 'Kam', sales: 1950000, transactions: 40 },
  { day: 'Jum', sales: 2800000, transactions: 58 },
  { day: 'Sab', sales: 3200000, transactions: 65 },
  { day: 'Min', sales: 1500000, transactions: 30 },
];

const topProducts = [
  { name: 'Beef Burger', qty: 85, revenue: 2975000 },
  { name: 'Nasi Goreng Spesial', qty: 62, revenue: 2790000 },
  { name: 'Kopi Susu', qty: 78, revenue: 1950000 },
  { name: 'Chicken Wings', qty: 45, revenue: 1575000 },
  { name: 'French Fries', qty: 92, revenue: 1656000 },
];

export function ReportPage() {
  const [period, setPeriod] = useState('weekly');

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between flex-wrap gap-4">
        <div>
          <h1 className="text-2xl font-bold text-white">Laporan Penjualan</h1>
          <p className="text-gray-400 mt-1">Analisis performa bisnis</p>
        </div>
        <div className="flex items-center gap-3">
          <div className="flex bg-gray-800 rounded-lg p-1">
            {['daily', 'weekly', 'monthly'].map((p) => (
              <button
                key={p}
                onClick={() => setPeriod(p)}
                className={`px-3 py-1.5 rounded-md text-xs font-medium transition-colors ${
                  period === p ? 'bg-blue-600 text-white' : 'text-gray-400 hover:text-white'
                }`}
              >
                {p === 'daily' ? 'Harian' : p === 'weekly' ? 'Mingguan' : 'Bulanan'}
              </button>
            ))}
          </div>
          <button className="pos-btn-secondary text-sm">
            <Download className="w-4 h-4 mr-2" />
            Export
          </button>
        </div>
      </div>

      {/* Summary Cards */}
      <div className="grid grid-cols-2 lg:grid-cols-4 gap-4">
        {[
          { label: 'Total Penjualan', value: 'Rp 16.810.000', change: '+15.3%', icon: DollarSign, good: true },
          { label: 'Total Transaksi', value: '324', change: '+8.7%', icon: ShoppingCart, good: true },
          { label: 'Rata-rata/Transaksi', value: 'Rp 51.882', change: '+6.1%', icon: TrendingUp, good: true },
          { label: 'Produk Terjual', value: '1.247', change: '+12.4%', icon: Package, good: true },
        ].map((card) => (
          <div key={card.label} className="pos-card p-4">
            <div className="flex items-start justify-between mb-2">
              <p className="text-xs text-gray-400">{card.label}</p>
              <card.icon className="w-4 h-4 text-gray-500" />
            </div>
            <p className="text-lg font-bold text-white">{card.value}</p>
            <div className="flex items-center gap-1 mt-1">
              {card.good ? <TrendingUp className="w-3 h-3 text-emerald-500" /> : <TrendingDown className="w-3 h-3 text-red-500" />}
              <span className={`text-xs ${card.good ? 'text-emerald-500' : 'text-red-500'}`}>{card.change}</span>
            </div>
          </div>
        ))}
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Sales Chart */}
        <div className="lg:col-span-2 pos-card p-5">
          <h3 className="font-semibold text-white mb-4">Grafik Penjualan</h3>
          <div className="h-72">
            <ResponsiveContainer width="100%" height="100%">
              <BarChart data={weeklyData}>
                <CartesianGrid strokeDasharray="3 3" stroke="#1f2937" />
                <XAxis dataKey="day" stroke="#6b7280" fontSize={12} />
                <YAxis stroke="#6b7280" fontSize={12} />
                <Tooltip
                  contentStyle={{ backgroundColor: '#1f2937', border: '1px solid #374151', borderRadius: '8px' }}
                  labelStyle={{ color: '#f3f4f6' }}
                />
                <Bar dataKey="sales" fill="#3b82f6" radius={[4, 4, 0, 0]} />
              </BarChart>
            </ResponsiveContainer>
          </div>
        </div>

        {/* Top Products */}
        <div className="pos-card p-5">
          <h3 className="font-semibold text-white mb-4">Produk Terlaris</h3>
          <div className="space-y-3">
            {topProducts.map((product, idx) => (
              <div key={product.name} className="flex items-center gap-3">
                <div className="w-7 h-7 rounded-lg bg-gray-800 flex items-center justify-center text-xs font-bold text-gray-400">
                  {idx + 1}
                </div>
                <div className="flex-1 min-w-0">
                  <p className="text-sm font-medium text-white truncate">{product.name}</p>
                  <p className="text-xs text-gray-500">{product.qty} terjual</p>
                </div>
                <p className="text-sm font-bold text-white">Rp {product.revenue.toLocaleString()}</p>
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
  );
}

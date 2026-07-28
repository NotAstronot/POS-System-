import React from 'react';
import {
  TrendingUp, DollarSign, ShoppingBag, Users, Clock,
  ArrowUp, ArrowDown, Package, AlertTriangle, Activity
} from 'lucide-react';
import { useStore } from '../store/useStore';

const stats = [
  { label: 'Penjualan Hari Ini', value: 'Rp 2.450.000', change: '+12.5%', icon: DollarSign, color: 'emerald', increase: true },
  { label: 'Total Transaksi', value: '48', change: '+8.2%', icon: ShoppingBag, color: 'blue', increase: true },
  { label: 'Rata-rata Pesanan', value: 'Rp 51.042', change: '+3.1%', icon: TrendingUp, color: 'violet', increase: true },
  { label: 'Pelanggan', value: '36', change: '-2.4%', icon: Users, color: 'amber', increase: false },
];

const recentOrders = [
  { no: 'ORD-001', total: 85000, payment: 'Tunai', time: '10:32', status: 'completed' },
  { no: 'ORD-002', total: 45000, payment: 'QRIS', time: '10:15', status: 'completed' },
  { no: 'ORD-003', total: 120000, payment: 'Debit', time: '09:45', status: 'completed' },
  { no: 'ORD-004', total: 35000, payment: 'Tunai', time: '09:20', status: 'completed' },
  { no: 'ORD-005', total: 67000, payment: 'E-Wallet', time: '08:55', status: 'completed' },
];

const lowStockItems = [
  { product: 'Beef Burger', stock: 3, min: 5 },
  { product: 'Nasi Goreng Spesial', stock: 2, min: 5 },
  { product: 'French Fries', stock: 4, min: 10 },
];

export function DashboardPage() {
  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-white">Dashboard</h1>
        <p className="text-gray-400 mt-1">Ringkasan penjualan dan performa outlet</p>
      </div>

      {/* Stats grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        {stats.map((stat) => {
          const Icon = stat.icon;
          const colorClasses = {
            emerald: 'from-emerald-600/20 to-emerald-600/5 border-emerald-800/30 text-emerald-400',
            blue: 'from-blue-600/20 to-blue-600/5 border-blue-800/30 text-blue-400',
            violet: 'from-violet-600/20 to-violet-600/5 border-violet-800/30 text-violet-400',
            amber: 'from-amber-600/20 to-amber-600/5 border-amber-800/30 text-amber-400',
          }[stat.color];

          return (
            <div key={stat.label} className={`pos-card p-5 bg-gradient-to-br ${colorClasses}`}>
              <div className="flex items-start justify-between">
                <div>
                  <p className="text-sm text-gray-400">{stat.label}</p>
                  <p className="text-2xl font-bold text-white mt-1">{stat.value}</p>
                </div>
                <div className="p-3 rounded-xl bg-gray-800/50">
                  <Icon className="w-5 h-5" />
                </div>
              </div>
              <div className="flex items-center gap-1 mt-3">
                {stat.increase
                  ? <ArrowUp className="w-3.5 h-3.5 text-emerald-500" />
                  : <ArrowDown className="w-3.5 h-3.5 text-red-500" />
                }
                <span className={`text-xs ${stat.increase ? 'text-emerald-500' : 'text-red-500'}`}>
                  {stat.change} dari kemarin
                </span>
              </div>
            </div>
          );
        })}
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Recent Orders */}
        <div className="lg:col-span-2 pos-card p-5">
          <div className="flex items-center justify-between mb-4">
            <h3 className="font-semibold text-white">Pesanan Terakhir</h3>
            <Clock className="w-4 h-4 text-gray-500" />
          </div>
          <div className="space-y-2">
            {recentOrders.map((order) => (
              <div key={order.no} className="flex items-center justify-between py-2.5 border-b border-gray-800 last:border-0">
                <div>
                  <p className="text-sm font-medium text-white">{order.no}</p>
                  <p className="text-xs text-gray-500">{order.payment} &middot; {order.time}</p>
                </div>
                <div className="text-right">
                  <p className="text-sm font-bold text-white">Rp {order.total.toLocaleString()}</p>
                  <span className="pos-badge-success text-[10px]">Selesai</span>
                </div>
              </div>
            ))}
          </div>
        </div>

        {/* Low Stock Alert */}
        <div className="pos-card p-5">
          <div className="flex items-center justify-between mb-4">
            <h3 className="font-semibold text-white flex items-center gap-2">
              <AlertTriangle className="w-4 h-4 text-amber-400" />
              Stok Menipis
            </h3>
          </div>
          {lowStockItems.length === 0 ? (
            <p className="text-gray-500 text-sm">Semua stok aman</p>
          ) : (
            <div className="space-y-3">
              {lowStockItems.map((item) => (
                <div key={item.product} className="flex items-center justify-between p-3 bg-red-900/10 border border-red-800/30 rounded-xl">
                  <div>
                    <p className="text-sm font-medium text-white">{item.product}</p>
                    <p className="text-xs text-gray-400">Min: {item.min}</p>
                  </div>
                  <span className="text-sm font-bold text-red-400">{item.stock}</span>
                </div>
              ))}
            </div>
          )}

          <div className="mt-6">
            <h4 className="text-sm font-medium text-gray-300 mb-3">Penjualan per Metode</h4>
            <div className="space-y-2">
              {[
                { method: 'Tunai', amount: 'Rp 1.120.000', pct: 45.7 },
                { method: 'QRIS', amount: 'Rp 780.000', pct: 31.8 },
                { method: 'Debit/Kredit', amount: 'Rp 550.000', pct: 22.5 },
              ].map((pm) => (
                <div key={pm.method}>
                  <div className="flex justify-between text-xs mb-1">
                    <span className="text-gray-400">{pm.method}</span>
                    <span className="text-white">{pm.amount}</span>
                  </div>
                  <div className="w-full h-1.5 bg-gray-800 rounded-full overflow-hidden">
                    <div
                      className="h-full rounded-full bg-gradient-to-r from-blue-600 to-blue-400"
                      style={{ width: `${pm.pct}%` }}
                    />
                  </div>
                </div>
              ))}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

import React, { useState, useEffect } from 'react';
import {
  TrendingUp, DollarSign, ShoppingBag, Clock,
  ArrowUp, ArrowDown, Package, AlertTriangle
} from 'lucide-react';
import { productService, revenueService, type Product, type RevenueSummary, type RecentOrder } from '../services/api';

const methodLabel: Record<string, string> = {
  cash: 'Tunai',
  qris: 'QRIS',
  debit: 'Debit',
  credit: 'Kredit',
  ewallet: 'E-Wallet',
  paylater: 'PayLater',
};

export function DashboardPage() {
  const [lowStockItems, setLowStockItems] = useState<Product[]>([]);
  const [revenue, setRevenue] = useState<RevenueSummary | null>(null);
  const [recentOrders, setRecentOrders] = useState<RecentOrder[]>([]);

  const load = () => {
    productService.lowStock(5)
      .then((res) => setLowStockItems(res.data.data || []))
      .catch(() => {});
    revenueService.summary('daily')
      .then((res) => {
        setRevenue(res.data.data);
        setRecentOrders(res.data.data.recent_orders || []);
      })
      .catch(() => {});
  };

  useEffect(() => {
    load();
    const timer = setInterval(load, 30000); // real-time: refresh tiap 30 detik
    return () => clearInterval(timer);
  }, []);

  const stats = [
    {
      label: 'Penjualan Hari Ini',
      value: revenue ? `Rp ${Math.round(revenue.total_sales).toLocaleString()}` : '-',
      change: revenue ? `${revenue.change_pct >= 0 ? '+' : ''}${revenue.change_pct.toFixed(1)}%` : '-',
      icon: DollarSign,
      color: 'success',
      increase: (revenue?.change_pct ?? 0) >= 0,
    },
    {
      label: 'Total Transaksi',
      value: revenue ? String(revenue.total_transactions) : '-',
      change: `${revenue?.products_sold ?? '-'} produk terjual`,
      icon: ShoppingBag,
      color: 'brand',
      increase: true,
    },
    {
      label: 'Rata-rata Pesanan',
      value: revenue ? `Rp ${Math.round(revenue.avg_per_tx).toLocaleString()}` : '-',
      change: `dari ${revenue?.total_transactions ?? 0} transaksi`,
      icon: TrendingUp,
      color: 'brand',
      increase: true,
    },
    {
      label: 'Produk Terjual',
      value: revenue ? String(revenue.products_sold) : '-',
      change: 'unit hari ini',
      icon: Package,
      color: 'warning',
      increase: true,
    },
  ];

  const colorClasses = {
    success: 'bg-success/10 border-success/20 text-success',
    brand: 'bg-primary-50 border-primary-200 text-brand',
    warning: 'bg-warning/10 border-warning/20 text-warning-dark',
    danger: 'bg-danger/10 border-danger/20 text-danger',
  };

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-neutral-900">Dashboard</h1>
        <p className="text-neutral-500 mt-1">Ringkasan penjualan dan performa outlet</p>
      </div>

      {/* Stats grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        {stats.map((stat) => {
          const Icon = stat.icon;
          const colorClass = colorClasses[stat.color as keyof typeof colorClasses];

          return (
            <div key={stat.label} className={`pos-card p-5 ${colorClass}`}>
              <div className="flex items-start justify-between">
                <div>
                  <p className="text-sm text-neutral-500">{stat.label}</p>
                  <p className="text-2xl font-bold text-neutral-900 mt-1">{stat.value}</p>
                </div>
                <div className="p-3 rounded-xl bg-neutral-100">
                  <Icon className="w-5 h-5" />
                </div>
              </div>
              <div className="flex items-center gap-1 mt-3">
                {stat.increase
                  ? <ArrowUp className="w-3.5 h-3.5 text-success" />
                  : <ArrowDown className="w-3.5 h-3.5 text-danger" />
                }
                <span className={`text-xs ${stat.increase ? 'text-success' : 'text-danger'}`}>
                  {stat.change}
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
            <h3 className="font-semibold text-neutral-900">Pesanan Terakhir</h3>
            <Clock className="w-4 h-4 text-neutral-500" />
          </div>
          <div className="space-y-2">
            {recentOrders.length === 0 && (
              <p className="text-neutral-500 text-sm py-4 text-center">Belum ada pesanan hari ini</p>
            )}
            {recentOrders.map((order) => (
              <div key={order.id} className="flex items-center justify-between py-2.5 border-b border-neutral-200 last:border-0">
                <div>
                  <p className="text-sm font-medium text-neutral-900">{order.order_number || `#${order.id}`}</p>
                  <p className="text-xs text-neutral-500">
                    {methodLabel[order.payment_method] || order.payment_method} &middot; {order.cashier_name || 'Kasir'} &middot; {order.item_count} item
                  </p>
                </div>
                <div className="text-right">
                  <p className="text-sm font-bold text-neutral-900">Rp {Math.round(order.total).toLocaleString()}</p>
                  <span className="pos-badge-success text-[10px]">Selesai</span>
                </div>
              </div>
            ))}
          </div>
        </div>

        {/* Low Stock Alert */}
        <div className="pos-card p-5">
          <div className="flex items-center justify-between mb-4">
            <h3 className="font-semibold text-neutral-900 flex items-center gap-2">
              <AlertTriangle className="w-4 h-4 text-warning" />
              Stok Menipis
            </h3>
          </div>
          {lowStockItems.length === 0 ? (
            <p className="text-neutral-500 text-sm">Semua stok aman</p>
          ) : (
            <div className="space-y-3">
              {lowStockItems.map((item) => (
                <div key={item.id} className="flex items-center justify-between p-3 bg-danger/5 border border-danger/20 rounded-xl">
                  <div>
                    <p className="text-sm font-medium text-neutral-900">{item.name}</p>
                    <p className="text-xs text-neutral-500">Min: 5</p>
                  </div>
                  <span className="text-sm font-bold text-danger">{item.stock}</span>
                </div>
              ))}
            </div>
          )}

          <div className="mt-6">
            <h4 className="text-sm font-medium text-neutral-600 mb-3">Penjualan per Metode</h4>
            <div className="space-y-2">
              {(() => {
                const byMethod = new Map<string, number>();
                recentOrders.forEach((o) => {
                  const m = methodLabel[o.payment_method] || o.payment_method;
                  byMethod.set(m, (byMethod.get(m) || 0) + o.total);
                });
                const total = Array.from(byMethod.values()).reduce((s, v) => s + v, 0);
                const list = Array.from(byMethod.entries());
                if (list.length === 0) {
                  return <p className="text-neutral-500 text-xs">Belum ada data penjualan</p>;
                }
                return list.map(([method, amount]) => (
                  <div key={method}>
                    <div className="flex justify-between text-xs mb-1">
                      <span className="text-neutral-500">{method}</span>
                      <span className="text-neutral-900">Rp {Math.round(amount).toLocaleString()}</span>
                    </div>
                    <div className="w-full h-1.5 bg-neutral-200 rounded-full overflow-hidden">
                      <div
                        className="h-full rounded-full bg-brand"
                        style={{ width: `${total > 0 ? (amount / total) * 100 : 0}%` }}
                      />
                    </div>
                  </div>
                ));
              })()}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
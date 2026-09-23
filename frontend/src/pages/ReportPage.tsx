import React, { useState, useEffect, useMemo, useCallback } from 'react';
import {
  BarChart3, Download, Calendar, DollarSign, ShoppingCart,
  TrendingUp, Package, Receipt, Store, Landmark, Wallet,
  Loader2, ShoppingBag, PieChart, Layers
} from 'lucide-react';
import {
  BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer,
  PieChart as RePieChart, Pie, Cell
} from 'recharts';
import * as XLSX from 'xlsx';
import {
  reportService,
  type ReportSummaryData, type Outlet,
} from '../services/api';

const PERIODS: { value: string; label: string }[] = [
  { value: 'today', label: 'Today' },
  { value: 'yesterday', label: 'Yesterday' },
  { value: 'this_week', label: 'This Week' },
  { value: 'last_week', label: 'Last Week' },
  { value: 'this_month', label: 'This Month' },
  { value: 'last_month', label: 'Last Month' },
  { value: 'this_year', label: 'This Year' },
  { value: 'last_year', label: 'Last Year' },
];

const PERIOD_LABEL: Record<string, string> = Object.fromEntries(PERIODS.map((p) => [p.value, p.label]));

const BUCKET_LABEL: Record<string, string> = {
  hourly: 'per jam',
  daily: 'per hari',
  monthly: 'per bulan',
};

const CHART_COLORS = [
  '#3b82f6', '#10b981', '#f59e0b', '#ef4444', '#8b5cf6',
  '#06b6d4', '#ec4899', '#84cc16', '#f97316', '#64748b',
  '#14b8a6', '#a855f7',
];

const fmtRp = (n: number) => 'Rp ' + Math.round(n).toLocaleString('id-ID');
const fmtNum = (n: number) => Math.round(n).toLocaleString('id-ID');

const tooltipStyle = {
  backgroundColor: '#1f2937',
  border: '1px solid #374151',
  borderRadius: '8px',
  fontSize: 12,
};

export function ReportPage() {
  const [period, setPeriod] = useState('this_month');
  const [outletId, setOutletId] = useState<string>('');
  const [outlets, setOutlets] = useState<Outlet[]>([]);
  const [report, setReport] = useState<ReportSummaryData | null>(null);
  const [loading, setLoading] = useState(true);

  const loadOutlets = useCallback(async () => {
    try {
      const res = await reportService.outlets();
      setOutlets(res.data.data || []);
    } catch {}
  }, []);

  useEffect(() => {
    loadOutlets();
  }, [loadOutlets]);

  useEffect(() => {
    let cancelled = false;
    setLoading(true);
    setReport(null);
    reportService.summary({
      period,
      ...(outletId ? { outlet_id: Number(outletId) } : {}),
    })
      .then((res) => {
        if (!cancelled) setReport(res.data.data);
      })
      .catch(() => {
        if (!cancelled) setReport(null);
      })
      .finally(() => {
        if (!cancelled) setLoading(false);
      });
    return () => { cancelled = true; };
  }, [period, outletId]);

  const kpis = report?.kpis;
  const selectedOutlet = outlets.find((o) => String(o.id) === outletId);
  const rangeLabel = report
    ? `${report.range_start} s.d. ${report.range_end}`
    : '';

  const dailyChartData = useMemo(() => {
    if (!report) return [];
    const mode = report.bucket_mode;
    return report.daily_sales.map((d) => {
      let label = d.label;
      if (mode === 'daily') {
        const [y, m, day] = d.label.split('-');
        label = `${day}/${m}`;
      } else if (mode === 'monthly') {
        const [y, m] = d.label.split('-');
        const months = ['Jan', 'Feb', 'Mar', 'Apr', 'Mei', 'Jun', 'Jul', 'Agu', 'Sep', 'Okt', 'Nov', 'Des'];
        label = months[Number(m) - 1] || d.label;
      }
      return { label, gross: d.gross_sales, net: d.net_sales, tx: d.transactions };
    });
  }, [report]);

  const dowChartData = useMemo(
    () => (report ? report.day_of_week_sales.map((d) => ({ day: d.day, gross: d.gross_sales, tx: d.transactions })) : []),
    [report]
  );

  const hourlyChartData = useMemo(
    () => (report ? report.hourly_sales.map((h) => ({ hour: h.hour, gross: h.gross_sales, tx: h.transactions })) : []),
    [report]
  );

  const exportExcel = () => {
    if (!report) return;
    const wb = XLSX.utils.book_new();
    const title = `Laporan Penjualan - ${PERIOD_LABEL[report.period]}`;

    const kpiRows: any[][] = [
      [title],
      [`Outlet: ${selectedOutlet?.name || 'Semua Outlet'} | Periode: ${rangeLabel}`],
      [],
      ['Indikator', 'Nilai'],
      ['Gross Sales', report.kpis.gross_sales],
      ['Net Sales', report.kpis.net_sales],
      ['Gross Profit', report.kpis.gross_profit],
      ['Transaksi', report.kpis.transactions],
      ['Rata-rata / Transaksi', report.kpis.avg_sale_per_tx],
      ['Gross Margin (%)', report.kpis.gross_margin_pct],
      ['Item Terjual', report.kpis.items_sold],
    ];
    const wsKpi = XLSX.utils.aoa_to_sheet(kpiRows);
    XLSX.utils.book_append_sheet(wb, wsKpi, 'Ringkasan');

    const itemRows: any[][] = [
      [title],
      [],
      ['Item', 'Terjual', 'Gross Sales', 'Net Sales', 'Gross Profit'],
      ...report.items.map((it) => [it.name, it.items_sold, it.gross_sales, it.net_sales, it.gross_profit]),
    ];
    const wsItems = XLSX.utils.aoa_to_sheet(itemRows);
    XLSX.utils.book_append_sheet(wb, wsItems, 'Item Summary');

    const dailyRows: any[][] = [
      [title],
      [],
      ['Label', 'Gross Sales', 'Net Sales', 'Transaksi'],
      ...report.daily_sales.map((d) => [d.label, d.gross_sales, d.net_sales, d.transactions]),
    ];
    const wsDaily = XLSX.utils.aoa_to_sheet(dailyRows);
    XLSX.utils.book_append_sheet(wb, wsDaily, 'Grafik Harian');

    XLSX.writeFile(wb, `Laporan-${report.period}-${new Date().toISOString().slice(0, 10)}.xlsx`);
  };

  const renderDonut = (
    title: string,
    icon: React.ReactNode,
    data: { name: string; value: number; pct: number }[] | undefined,
    valueFmt: (n: number) => string
  ) => {
    const rows = data || [];
    return (
      <div className="pos-card p-5">
        <h3 className="font-semibold text-white mb-1 flex items-center gap-2">
          {icon}
          {title}
        </h3>
        {rows.length === 0 ? (
          <p className="text-sm text-gray-500 py-10 text-center">Tidak ada data</p>
        ) : (
          <div className="flex flex-col md:flex-row items-center gap-6">
            <div className="relative w-56 h-56 shrink-0">
              <ResponsiveContainer width="100%" height="100%">
                <RePieChart>
                  <Pie
                    data={rows}
                    dataKey="value"
                    nameKey="name"
                    cx="50%"
                    cy="50%"
                    innerRadius={58}
                    outerRadius={88}
                    paddingAngle={2}
                    label={(e: any) => (e.pct >= 3 ? `${e.pct.toFixed(0)}%` : '')}
                    labelLine={false}
                  >
                    {rows.map((_, i) => (
                      <Cell key={i} fill={CHART_COLORS[i % CHART_COLORS.length]} />
                    ))}
                  </Pie>
                  <Tooltip
                    formatter={(v: any) => [valueFmt(Number(v)), '']}
                    contentStyle={tooltipStyle}
                    labelStyle={{ color: '#f3f4f6' }}
                  />
                </RePieChart>
              </ResponsiveContainer>
              <div className="absolute inset-0 flex flex-col items-center justify-center pointer-events-none">
                <p className="text-xl font-bold text-white">
                  {rows.length > 0 ? fmtNum(rows.reduce((s, r) => s + r.value, 0)) : '0'}
                </p>
                <p className="text-[10px] text-gray-500 uppercase tracking-wide">Total</p>
              </div>
            </div>
            {/* Legend / keterangan warna */}
            <div className="flex-1 w-full space-y-2">
              {rows.map((r, i) => (
                <div key={r.name} className="flex items-center gap-2 text-sm">
                  <span
                    className="w-3 h-3 rounded-full shrink-0"
                    style={{ backgroundColor: CHART_COLORS[i % CHART_COLORS.length] }}
                  />
                  <span className="flex-1 text-gray-300 truncate">{r.name}</span>
                  <span className="text-gray-400 text-xs">{valueFmt(r.value)}</span>
                  <span className="w-12 text-right font-semibold text-white">{r.pct.toFixed(1)}%</span>
                </div>
              ))}
            </div>
          </div>
        )}
      </div>
    );
  };

  return (
    <div className="space-y-6">
      {/* ===== Header & Filter ===== */}
      <div className="pos-card p-4">
        <div className="flex flex-wrap items-end gap-4">
          <div>
            <h1 className="text-xl font-bold text-white flex items-center gap-2">
              <BarChart3 className="w-5 h-5 text-blue-400" />
              Laporan Penjualan
            </h1>
            <p className="text-gray-500 text-sm mt-0.5">
              {selectedOutlet ? `Outlet: ${selectedOutlet.name}` : 'Semua Outlet (All Outlets)'}
              {' | '}
              {PERIOD_LABEL[period]}
              {rangeLabel ? ` | ${rangeLabel}` : ''}
            </p>
          </div>
          <div className="flex-1" />
          {/* Pilihan outlet */}
          <div>
            <label className="block text-[11px] uppercase tracking-wide text-gray-500 mb-1 font-medium">
              <Store className="w-3 h-3 inline mr-1" />
              Outlet
            </label>
            <select
              value={outletId}
              onChange={(e) => setOutletId(e.target.value)}
              className="pos-input text-sm min-w-[180px]"
            >
              <option value="">All Outlets</option>
              {outlets.map((o) => (
                <option key={o.id} value={String(o.id)}>{o.name}</option>
              ))}
            </select>
          </div>
          {/* Pilihan periode */}
          <div>
            <label className="block text-[11px] uppercase tracking-wide text-gray-500 mb-1 font-medium">
              <Calendar className="w-3 h-3 inline mr-1" />
              Periode
            </label>
            <select
              value={period}
              onChange={(e) => setPeriod(e.target.value)}
              className="pos-input text-sm min-w-[160px]"
            >
              {PERIODS.map((p) => (
                <option key={p.value} value={p.value}>{p.label}</option>
              ))}
            </select>
          </div>
          <button onClick={exportExcel} disabled={!report} className="pos-btn-secondary text-sm disabled:opacity-40">
            <Download className="w-4 h-4 mr-2" />
            Export Excel
          </button>
        </div>
      </div>

      {loading ? (
        <div className="flex items-center justify-center py-24">
          <Loader2 className="w-8 h-8 text-blue-500 animate-spin" />
        </div>
      ) : !report ? (
        <div className="pos-card p-10 text-center text-gray-500">
          Tidak ada data untuk periode & outlet yang dipilih.
        </div>
      ) : (
        <>
          {/* ===== Summary (KPI) ===== */}
          <div className="grid grid-cols-2 md:grid-cols-3 xl:grid-cols-6 gap-3">
            {[
              { label: 'Gross Sales', value: fmtRp(kpis!.gross_sales), icon: DollarSign, color: 'text-blue-400' },
              { label: 'Net Sales', value: fmtRp(kpis!.net_sales), icon: Landmark, color: 'text-emerald-400' },
              { label: 'Gross Profit', value: fmtRp(kpis!.gross_profit), icon: Wallet, color: 'text-violet-400' },
              { label: 'Transactions', value: fmtNum(kpis!.transactions), icon: ShoppingCart, color: 'text-amber-400' },
              { label: 'Avg Sale / Tx', value: fmtRp(kpis!.avg_sale_per_tx), icon: TrendingUp, color: 'text-cyan-400' },
              { label: 'Gross Margin', value: `${kpis!.gross_margin_pct.toFixed(1)}%`, icon: PieChart, color: 'text-pink-400' },
            ].map((card) => (
              <div key={card.label} className="pos-card p-4">
                <div className="flex items-start justify-between mb-2">
                  <p className="text-[11px] text-gray-400 uppercase tracking-wide">{card.label}</p>
                  <card.icon className={`w-4 h-4 ${card.color}`} />
                </div>
                <p className="text-lg font-bold text-white truncate">{card.value}</p>
                <div className="mt-1 flex items-center gap-1.5">
                  {kpis!.change_pct >= 0 ? (
                    <TrendingUp className="w-3 h-3 text-emerald-500" />
                  ) : (
                    <TrendingUp className="w-3 h-3 text-red-500 rotate-180" />
                  )}
                  <span className={`text-xs ${kpis!.change_pct >= 0 ? 'text-emerald-500' : 'text-red-500'}`}>
                    {kpis!.change_pct >= 0 ? '+' : ''}{kpis!.change_pct.toFixed(1)}%
                  </span>
                  <span className="text-[10px] text-gray-600">vs periode sebelumnya</span>
                </div>
              </div>
            ))}
          </div>

          {/* ===== Daily Gross Sales Amount ===== */}
          <div className="pos-card p-5">
            <h3 className="font-semibold text-white mb-1">Daily Gross Sales Amount</h3>
            <p className="text-xs text-gray-500 mb-4">
              Penjualan kotor {BUCKET_LABEL[report.bucket_mode]} — grafik menyesuaikan periode yang dipilih ({PERIOD_LABEL[report.period]})
            </p>
            <div className="h-72">
              <ResponsiveContainer width="100%" height="100%">
                <BarChart data={dailyChartData}>
                  <CartesianGrid strokeDasharray="3 3" stroke="#1f2937" />
                  <XAxis dataKey="label" stroke="#6b7280" fontSize={11} interval={report.bucket_mode === 'daily' ? 3 : 0} />
                  <YAxis stroke="#6b7280" fontSize={11} tickFormatter={(v) => (v >= 1000000 ? `${(v / 1000000).toFixed(1)}jt` : v >= 1000 ? `${Math.round(v / 1000)}rb` : String(v))} />
                  <Tooltip formatter={(v: any) => [fmtRp(Number(v)), 'Gross Sales']} contentStyle={tooltipStyle} labelStyle={{ color: '#f3f4f6' }} />
                  <Bar dataKey="gross" fill="#3b82f6" radius={[4, 4, 0, 0]} name="Gross Sales" />
                </BarChart>
              </ResponsiveContainer>
            </div>
          </div>

          {/* ===== Day of Week + Hourly ===== */}
          <div className="grid grid-cols-1 xl:grid-cols-2 gap-6">
            <div className="pos-card p-5">
              <h3 className="font-semibold text-white mb-4">Day Of The Week Gross Sales Amount</h3>
              <div className="h-64">
                <ResponsiveContainer width="100%" height="100%">
                  <BarChart data={dowChartData}>
                    <CartesianGrid strokeDasharray="3 3" stroke="#1f2937" />
                    <XAxis dataKey="day" stroke="#6b7280" fontSize={11} />
                    <YAxis stroke="#6b7280" fontSize={11} tickFormatter={(v) => (v >= 1000 ? `${Math.round(v / 1000)}rb` : String(v))} />
                    <Tooltip formatter={(v: any) => [fmtRp(Number(v)), 'Gross Sales']} contentStyle={tooltipStyle} labelStyle={{ color: '#f3f4f6' }} />
                    <Bar dataKey="gross" fill="#10b981" radius={[4, 4, 0, 0]} />
                  </BarChart>
                </ResponsiveContainer>
              </div>
            </div>
            <div className="pos-card p-5">
              <h3 className="font-semibold text-white mb-4">Hourly Gross Sales Amount</h3>
              <div className="h-64">
                <ResponsiveContainer width="100%" height="100%">
                  <BarChart data={hourlyChartData}>
                    <CartesianGrid strokeDasharray="3 3" stroke="#1f2937" />
                    <XAxis dataKey="hour" stroke="#6b7280" fontSize={10} interval={2} />
                    <YAxis stroke="#6b7280" fontSize={11} tickFormatter={(v) => (v >= 1000 ? `${Math.round(v / 1000)}rb` : String(v))} />
                    <Tooltip formatter={(v: any) => [fmtRp(Number(v)), 'Gross Sales']} contentStyle={tooltipStyle} labelStyle={{ color: '#f3f4f6' }} />
                    <Bar dataKey="gross" fill="#f59e0b" radius={[4, 4, 0, 0]} />
                  </BarChart>
                </ResponsiveContainer>
              </div>
            </div>
          </div>

          {/* ===== Item Summary ===== */}
          <div className="pos-card p-5">
            <h3 className="font-semibold text-white mb-4 flex items-center gap-2">
              <Package className="w-4 h-4 text-blue-400" />
              Item Summary
            </h3>
            {report.items.length === 0 ? (
              <p className="text-sm text-gray-500 py-6 text-center">Tidak ada penjualan item</p>
            ) : (
              <div className="overflow-x-auto">
                <table className="w-full text-sm">
                  <thead>
                    <tr className="border-b border-gray-800 text-left text-[11px] uppercase tracking-wide text-gray-400">
                      <th className="px-3 py-2.5">Item</th>
                      <th className="px-3 py-2.5 text-right">Items Sold</th>
                      <th className="px-3 py-2.5 text-right">Gross Sales</th>
                      <th className="px-3 py-2.5 text-right">Net Sales</th>
                      <th className="px-3 py-2.5 text-right">Gross Profit</th>
                    </tr>
                  </thead>
                  <tbody className="divide-y divide-gray-800/70">
                    {report.items.map((it) => (
                      <tr key={it.name}>
                        <td className="px-3 py-2.5 text-white">{it.name}</td>
                        <td className="px-3 py-2.5 text-right text-gray-300">{fmtNum(it.items_sold)}</td>
                        <td className="px-3 py-2.5 text-right text-white">{fmtRp(it.gross_sales)}</td>
                        <td className="px-3 py-2.5 text-right text-emerald-400">{fmtRp(it.net_sales)}</td>
                        <td className="px-3 py-2.5 text-right font-semibold text-violet-400">{fmtRp(it.gross_profit)}</td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            )}
          </div>

          {/* ===== Category By Volume + Category By Sales ===== */}
          <div className="grid grid-cols-1 xl:grid-cols-2 gap-6">
            {renderDonut(
              'Category By Volume',
              <ShoppingBag className="w-4 h-4 text-blue-400" />,
              report.category_by_volume.map((c) => ({ name: c.category, value: c.value, pct: c.pct })),
              fmtNum
            )}
            {renderDonut(
              'Category By Sales',
              <Landmark className="w-4 h-4 text-emerald-400" />,
              report.category_by_sales.map((c) => ({ name: c.category, value: c.value, pct: c.pct })),
              fmtRp
            )}
          </div>

          {/* ===== Top Item By Category ===== */}
          <div className="pos-card p-5">
            <h3 className="font-semibold text-white mb-4 flex items-center gap-2">
              <Layers className="w-4 h-4 text-amber-400" />
              Top Item By Category
            </h3>
            {report.top_items_by_category.length === 0 ? (
              <p className="text-sm text-gray-500 py-6 text-center">Tidak ada data</p>
            ) : (
              <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-4">
                {report.top_items_by_category.map((cat) => (
                  <div key={cat.category} className="bg-gray-800/50 border border-gray-700/60 rounded-xl p-4">
                    <p className="font-semibold text-white mb-3 flex items-center gap-2">
                      <Receipt className="w-4 h-4 text-blue-400" />
                      {cat.category}
                    </p>
                    <div className="space-y-2.5">
                      {cat.items.map((it, i) => (
                        <div key={it.name} className="flex items-center gap-2 text-sm">
                          <span
                            className="w-5 h-5 rounded-md flex items-center justify-center text-[10px] font-bold shrink-0"
                            style={{
                              backgroundColor: CHART_COLORS[i % CHART_COLORS.length] + '33',
                              color: CHART_COLORS[i % CHART_COLORS.length],
                            }}
                          >
                            {i + 1}
                          </span>
                          <span className="flex-1 text-gray-200 truncate">{it.name}</span>
                          <span className="text-xs text-gray-400">{fmtNum(it.qty)} pcs</span>
                          <span className="text-xs font-semibold text-white w-24 text-right">{fmtRp(it.gross_sales)}</span>
                        </div>
                      ))}
                    </div>
                  </div>
                ))}
              </div>
            )}
          </div>
        </>
      )}
    </div>
  );
}
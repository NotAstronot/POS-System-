import React from 'react';
import { Outlet, NavLink, useNavigate, useLocation } from 'react-router-dom';
import {
  ShoppingCart, LayoutDashboard, Receipt, LogOut, Store,
  ChevronDown, Menu, X, Printer, BarChart3, Package, Shield, Building2, Users, FolderOpen, Truck,
  Banknote, RotateCcw, Landmark, BookOpen, List, ScrollText, CalendarCheck, Briefcase, ClipboardList, Percent,
  Globe, Factory, Landmark as BankIcon, FileSpreadsheet
} from 'lucide-react';
import { useState } from 'react';
import { useStore } from '../store/useStore';

export function Layout() {
  const navigate = useNavigate();
  const location = useLocation();
  const { user, logout, outletName } = useStore();
  const [sidebarOpen, setSidebarOpen] = useState(false);
  const [perusahaanOpen, setPerusahaanOpen] = useState(false);
  const [asetTetapOpen, setAsetTetapOpen] = useState(false);
  const [karyawanOpen, setKaryawanOpen] = useState(false);
  const [keuanganOpen, setKeuanganOpen] = useState(false);
  const [bukuBesarOpen, setBukuBesarOpen] = useState(false);
  const [pembelianOpen, setPembelianOpen] = useState(false);
  const [penjualanOpen, setPenjualanOpen] = useState(false);
  const [inventoryOpen, setInventoryOpen] = useState(false);
  const [smartlinkOpen, setSmartlinkOpen] = useState(false);

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  const isAdmin = user?.role === 'admin';
  const perms = user?.permissions || [];

  const hasPermission = (perm: string) => isAdmin || perms.includes(perm);

  const navItems = [
    { to: '/pos', label: 'POS Kasir', icon: ShoppingCart, perm: 'pos_access' },
    { to: '/dashboard', label: 'Dashboard', icon: LayoutDashboard, perm: 'pos_access' },
    { to: '/reports', label: 'Laporan', icon: BarChart3, perm: 'report_view' },
    { to: '/menu', label: 'Menu', icon: Package, perm: 'menu_manage' },
  ].filter((item) => hasPermission(item.perm));

  const perusahaanSubItems = [
    { to: '/branches', label: 'Cabang', icon: Building2 },
    { to: '/admin', label: 'Kelola User', icon: Shield },
  ];

  const asetTetapSubItems = [
    { to: '/fixed-assets', label: 'Pencatatan Aset', icon: ClipboardList },
    { to: '/depreciation', label: 'Penyusutan Otomatis', icon: Percent },
  ];

  const perusahaanAllPaths = [...perusahaanSubItems, ...asetTetapSubItems].map((s) => s.to);

  const karyawanSubItems = [
    { to: '/employees', label: 'Daftar Karyawan', icon: Users },
    { to: '/departments', label: 'Departemen', icon: FolderOpen },
    { to: '/commissions', label: 'Komisi', icon: Receipt },
    { to: '/salaries', label: 'Gaji', icon: Receipt },
  ];

  const keuanganSubItems = [
    { to: '/cash', label: 'Kas', icon: Receipt },
    { to: '/debts', label: 'Hutang', icon: Receipt },
    { to: '/receivables', label: 'Piutang', icon: Receipt },
    { to: '/withdrawals', label: 'Penarikan Dana', icon: Receipt },
    { to: '/bank-transfers', label: 'Transfer Bank', icon: Banknote },
    { to: '/bank-reconciliations', label: 'Rekonsiliasi Bank', icon: RotateCcw },
  ];

  const bukuBesarSubItems = [
    { to: '/chart-of-accounts', label: 'Akun Perkiraan', icon: List },
    { to: '/journal-vouchers', label: 'Jurnal Umum', icon: ScrollText },
    { to: '/period-end', label: 'Proses Akhir Bulan', icon: CalendarCheck },
  ];

  const pembelianSubItems = [
    { to: '/suppliers', label: 'Supplier', icon: Truck },
    { to: '/purchase-orders', label: 'Pesanan Pembelian', icon: ShoppingCart },
    { to: '/purchase-receive', label: 'Penerimaan Barang', icon: Package },
    { to: '/purchase-invoices', label: 'Faktur Pembelian', icon: Receipt },
    { to: '/purchase-payments', label: 'Pembayaran Pembelian', icon: Receipt },
    { to: '/purchase-returns', label: 'Retur Pembelian', icon: Receipt },
  ];

  const penjualanSubItems = [
    { to: '/customers', label: 'Pelanggan', icon: Users },
    { to: '/customer-categories', label: 'Kategori Pelanggan', icon: FolderOpen },
    { to: '/sales-categories', label: 'Kategori Penjualan', icon: FolderOpen },
    { to: '/sales-quotations', label: 'Penawaran Penjualan', icon: Receipt },
    { to: '/sales-orders', label: 'Pesanan Penjualan', icon: ShoppingCart },
    { to: '/delivery-orders', label: 'Pengiriman Pesanan', icon: Package },
    { to: '/sales-invoices', label: 'Faktur Penjualan', icon: Receipt },
    { to: '/sales-receipts', label: 'Penerimaan Penjualan', icon: Receipt },
    { to: '/sales-returns', label: 'Retur Penjualan', icon: Receipt },
  ];

  const inventorySubItems = [
    { to: '/items', label: 'Barang & Jasa', icon: Package },
    { to: '/warehouses', label: 'Gudang', icon: Building2 },
    { to: '/stock-transfers', label: 'Transfer Barang', icon: Truck },
    { to: '/stock-adjustments', label: 'Penyesuaian Persediaan', icon: Receipt },
    { to: '/price-changes', label: 'Penyesuaian Harga Jual', icon: Receipt },
  ];

  const smartlinkSubItems = [
    { to: '/smartlink/ecommerce', label: 'SmartLink e-Commerce', icon: Globe, perm: 'smartlink_manage' },
    { to: '/smartlink/ebanking', label: 'SmartLink e-Banking', icon: BankIcon, perm: 'smartlink_manage' },
    { to: '/smartlink/tax', label: 'SmartLink Tax / e-Faktur', icon: FileSpreadsheet, perm: 'smartlink_manage' },
    { to: '/smartlink/manufacturing', label: 'Modul Manufaktur', icon: Factory, perm: 'manufacturing_manage' },
  ];
  const smartlinkVisible = smartlinkSubItems.some((s) => hasPermission(s.perm));

  return (
    <div className="flex h-screen overflow-hidden">
      {/* Sidebar */}
      <aside className={`pos-sidebar ${sidebarOpen ? 'translate-x-0' : '-translate-x-full'}`}>
        <div className="flex flex-col h-full">
          {/* Logo */}
          <div className="flex items-center gap-3 px-6 py-5 border-b border-neutral-200">
            <div className="w-10 h-10 rounded-xl bg-brand flex items-center justify-center">
              <Store className="w-5 h-5 text-white" />
            </div>
            <div>
              <h1 className="font-bold text-neutral-900 text-lg leading-tight">POS System</h1>
              <p className="text-xs text-neutral-500">{outletName || 'Outlet'}</p>
            </div>
          </div>

          {/* Nav */}
          <nav className="flex-1 px-3 py-4 space-y-1 overflow-y-auto">
            {navItems.map((item) => (
              <NavLink
                key={item.to}
                to={item.to}
                onClick={() => setSidebarOpen(false)}
                className={({ isActive }) =>
                  isActive
                    ? 'pos-nav-item-active'
                    : 'pos-nav-item-inactive'
                }
              >
                <item.icon className="w-5 h-5" />
                {item.label}
              </NavLink>
            ))}

            {hasPermission('user_manage') && (
              <div>
                <button
                  onClick={() => setPerusahaanOpen(!perusahaanOpen)}
                  className={`w-full flex items-center justify-between pos-nav-item ${
                    perusahaanAllPaths.some((p) => location.pathname === p)
                      ? 'pos-nav-item-active'
                      : 'pos-nav-item-inactive'
                  }`}
                >
                  <span className="flex items-center gap-3"><Landmark className="w-5 h-5" /> Perusahaan</span>
                  <ChevronDown className={`w-4 h-4 transition-transform ${perusahaanOpen ? 'rotate-180' : ''}`} />
                </button>
                {perusahaanOpen && (
                  <div className="ml-4 mt-1 space-y-1">
                    {perusahaanSubItems.map((sub) => (
                      <NavLink
                        key={sub.to}
                        to={sub.to}
                        onClick={() => setSidebarOpen(false)}
                        className={({ isActive }) =>
                          isActive ? 'pos-nav-subitem-active' : 'pos-nav-subitem-inactive'
                        }
                      >
                        <sub.icon className="w-4 h-4" />
                        {sub.label}
                      </NavLink>
                    ))}
                    <div>
                      <button
                        onClick={() => setAsetTetapOpen(!asetTetapOpen)}
                        className={`w-full flex items-center justify-between pos-nav-subitem ${
                          asetTetapSubItems.some((sub) => location.pathname === sub.to)
                            ? 'pos-nav-subitem-active'
                            : 'pos-nav-subitem-inactive'
                        }`}
                      >
                        <span className="flex items-center gap-3"><Briefcase className="w-4 h-4" /> Aset Tetap</span>
                        <ChevronDown className={`w-4 h-4 transition-transform ${asetTetapOpen ? 'rotate-180' : ''}`} />
                      </button>
                      {asetTetapOpen && (
                        <div className="ml-4 mt-1 space-y-1">
                          {asetTetapSubItems.map((sub) => (
                            <NavLink
                              key={sub.to}
                              to={sub.to}
                              onClick={() => setSidebarOpen(false)}
                              className={({ isActive }) =>
                                isActive ? 'pos-nav-subitem-active' : 'pos-nav-subitem-inactive'
                              }
                            >
                              <sub.icon className="w-4 h-4" />
                              {sub.label}
                            </NavLink>
                          ))}
                        </div>
                      )}
                    </div>
                  </div>
                )}
                <button
                  onClick={() => setKaryawanOpen(!karyawanOpen)}
                  className={`w-full flex items-center justify-between pos-nav-item mt-1 ${
                    karyawanSubItems.some((sub) => location.pathname === sub.to)
                      ? 'pos-nav-item-active'
                      : 'pos-nav-item-inactive'
                  }`}
                >
                  <span className="flex items-center gap-3"><Users className="w-5 h-5" /> Karyawan</span>
                  <ChevronDown className={`w-4 h-4 transition-transform ${karyawanOpen ? 'rotate-180' : ''}`} />
                </button>
                {karyawanOpen && (
                  <div className="ml-4 mt-1 space-y-1">
                    {karyawanSubItems.map((sub) => (
                      <NavLink
                        key={sub.to}
                        to={sub.to}
                        onClick={() => setSidebarOpen(false)}
                        className={({ isActive }) =>
                          isActive ? 'pos-nav-subitem-active' : 'pos-nav-subitem-inactive'
                        }
                      >
                        <sub.icon className="w-4 h-4" />
                        {sub.label}
                      </NavLink>
                    ))}
                  </div>
                )}
                <button
                  onClick={() => setKeuanganOpen(!keuanganOpen)}
                  className={`w-full flex items-center justify-between pos-nav-item mt-1 ${
                    keuanganSubItems.some((sub) => location.pathname === sub.to)
                      ? 'pos-nav-item-active'
                      : 'pos-nav-item-inactive'
                  }`}
                >
                  <span className="flex items-center gap-3"><Receipt className="w-5 h-5" /> Keuangan</span>
                  <ChevronDown className={`w-4 h-4 transition-transform ${keuanganOpen ? 'rotate-180' : ''}`} />
                </button>
                {keuanganOpen && (
                  <div className="ml-4 mt-1 space-y-1">
                    {keuanganSubItems.map((sub) => (
                      <NavLink
                        key={sub.to}
                        to={sub.to}
                        onClick={() => setSidebarOpen(false)}
                        className={({ isActive }) =>
                          isActive ? 'pos-nav-subitem-active' : 'pos-nav-subitem-inactive'
                        }
                      >
                        <sub.icon className="w-4 h-4" />
                        {sub.label}
                      </NavLink>
                    ))}
                  </div>
                )}
                <button
                  onClick={() => setBukuBesarOpen(!bukuBesarOpen)}
                  className={`w-full flex items-center justify-between pos-nav-item mt-1 ${
                    bukuBesarSubItems.some((sub) => location.pathname === sub.to)
                      ? 'pos-nav-item-active'
                      : 'pos-nav-item-inactive'
                  }`}
                >
                  <span className="flex items-center gap-3"><BookOpen className="w-5 h-5" /> Buku Besar</span>
                  <ChevronDown className={`w-4 h-4 transition-transform ${bukuBesarOpen ? 'rotate-180' : ''}`} />
                </button>
                {bukuBesarOpen && (
                  <div className="ml-4 mt-1 space-y-1">
                    {bukuBesarSubItems.map((sub) => (
                      <NavLink
                        key={sub.to}
                        to={sub.to}
                        onClick={() => setSidebarOpen(false)}
                        className={({ isActive }) =>
                          isActive ? 'pos-nav-subitem-active' : 'pos-nav-subitem-inactive'
                        }
                      >
                        <sub.icon className="w-4 h-4" />
                        {sub.label}
                      </NavLink>
                    ))}
                  </div>
                )}
                <button
                  onClick={() => setPembelianOpen(!pembelianOpen)}
                  className={`w-full flex items-center justify-between pos-nav-item mt-1 ${
                    pembelianSubItems.some((sub) => location.pathname === sub.to)
                      ? 'pos-nav-item-active'
                      : 'pos-nav-item-inactive'
                  }`}
                >
                  <span className="flex items-center gap-3"><Truck className="w-5 h-5" /> Pembelian</span>
                  <ChevronDown className={`w-4 h-4 transition-transform ${pembelianOpen ? 'rotate-180' : ''}`} />
                </button>
                {pembelianOpen && (
                  <div className="ml-4 mt-1 space-y-1">
                    {pembelianSubItems.map((sub) => (
                      <NavLink
                        key={sub.to}
                        to={sub.to}
                        onClick={() => setSidebarOpen(false)}
                        className={({ isActive }) =>
                          isActive ? 'pos-nav-subitem-active' : 'pos-nav-subitem-inactive'
                        }
                      >
                        <sub.icon className="w-4 h-4" />
                        {sub.label}
                      </NavLink>
                    ))}
                  </div>
                )}
                <button
                  onClick={() => setPenjualanOpen(!penjualanOpen)}
                  className={`w-full flex items-center justify-between pos-nav-item mt-1 ${
                    penjualanSubItems.some((sub) => location.pathname === sub.to)
                      ? 'pos-nav-item-active'
                      : 'pos-nav-item-inactive'
                  }`}
                >
                  <span className="flex items-center gap-3"><ShoppingCart className="w-5 h-5" /> Penjualan</span>
                  <ChevronDown className={`w-4 h-4 transition-transform ${penjualanOpen ? 'rotate-180' : ''}`} />
                </button>
                {penjualanOpen && (
                  <div className="ml-4 mt-1 space-y-1">
                    {penjualanSubItems.map((sub) => (
                      <NavLink
                        key={sub.to}
                        to={sub.to}
                        onClick={() => setSidebarOpen(false)}
                        className={({ isActive }) =>
                          isActive ? 'pos-nav-subitem-active' : 'pos-nav-subitem-inactive'
                        }
                      >
                        <sub.icon className="w-4 h-4" />
                        {sub.label}
                      </NavLink>
                    ))}
                  </div>
                )}
              </div>
            )}

            {hasPermission('inventory_manage') && (
              <div>
                <button
                  onClick={() => setInventoryOpen(!inventoryOpen)}
                  className={`w-full flex items-center justify-between pos-nav-item mt-1 ${
                    inventorySubItems.some((sub) => location.pathname === sub.to)
                      ? 'pos-nav-item-active'
                      : 'pos-nav-item-inactive'
                  }`}
                >
                  <span className="flex items-center gap-3"><Package className="w-5 h-5" /> Persediaan & Gudang</span>
                  <ChevronDown className={`w-4 h-4 transition-transform ${inventoryOpen ? 'rotate-180' : ''}`} />
                </button>
                {inventoryOpen && (
                  <div className="ml-4 mt-1 space-y-1">
                    {inventorySubItems.map((sub) => (
                      <NavLink
                        key={sub.to}
                        to={sub.to}
                        onClick={() => setSidebarOpen(false)}
                        className={({ isActive }) =>
                          isActive ? 'pos-nav-subitem-active' : 'pos-nav-subitem-inactive'
                        }
                      >
                        <sub.icon className="w-4 h-4" />
                        {sub.label}
                      </NavLink>
                    ))}
                  </div>
                )}
              </div>
            )}
          {smartlinkVisible && (
              <div>
                <button
                  onClick={() => setSmartlinkOpen(!smartlinkOpen)}
                  className={`w-full flex items-center justify-between pos-nav-item mt-1 ${
                    smartlinkSubItems.some((s) => hasPermission(s.perm) && location.pathname === s.to)
                      ? 'pos-nav-item-active'
                      : 'pos-nav-item-inactive'
                  }`}
                >
                  <span className="flex items-center gap-3"><Store className="w-5 h-5" /> Store & SmartLink</span>
                  <ChevronDown className={`w-4 h-4 transition-transform ${smartlinkOpen ? 'rotate-180' : ''}`} />
                </button>
                {smartlinkOpen && (
                  <div className="ml-4 mt-1 space-y-1">
                    {smartlinkSubItems.filter((s) => hasPermission(s.perm)).map((sub) => (
                      <NavLink
                        key={sub.to}
                        to={sub.to}
                        onClick={() => setSidebarOpen(false)}
                        className={({ isActive }) =>
                          isActive ? 'pos-nav-subitem-active' : 'pos-nav-subitem-inactive'
                        }
                      >
                        <sub.icon className="w-4 h-4" />
                        {sub.label}
                      </NavLink>
                    ))}
                  </div>
                )}
              </div>
            )}
          </nav>

          {/* User */}
          <div className="px-4 py-4 border-t border-neutral-200">
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-3">
                <div className="pos-avatar">
                  {user?.name?.charAt(0)?.toUpperCase() || '?'}
                </div>
                <div>
                  <p className="text-sm font-medium text-neutral-900">{user?.name || 'Guest'}</p>
                  <p className="text-xs text-neutral-500 capitalize">{user?.role || ''}</p>
                </div>
              </div>
              <button onClick={handleLogout} className="p-2 rounded-lg text-neutral-500 hover:text-danger hover:bg-neutral-100 transition-colors">
                <LogOut className="w-4 h-4" />
              </button>
            </div>
          </div>
        </div>
      </aside>

      {/* Overlay */}
      {sidebarOpen && <div className="fixed inset-0 bg-black/50 z-40" onClick={() => setSidebarOpen(false)} />}

      {/* Main content */}
      <div className="flex-1 flex flex-col min-w-0 bg-neutral-50">
        {/* Top bar */}
        <header className="pos-header">
          <button onClick={() => setSidebarOpen(true)} className="p-2 rounded-lg text-neutral-500 hover:text-neutral-900 hover:bg-neutral-100">
            <Menu className="w-5 h-5" />
          </button>
          <div className="flex-1" />
          <div className="flex items-center gap-3 text-sm">
            <span className="text-neutral-500">Shift #{useStore.getState().activeShift?.shift_number || '-'}</span>
            <span className="pos-status-online" />
            <div className="flex items-center gap-2 ml-2 pl-3 border-l border-neutral-200">
              <div className="pos-avatar">
                {user?.name?.charAt(0)?.toUpperCase() || '?'}
              </div>
              <span className="text-neutral-900 font-medium">{user?.name || 'Guest'}</span>
            </div>
          </div>
        </header>

        {/* Page content */}
        <main className="pos-main">
          <Outlet />
        </main>
      </div>
    </div>
  );
}
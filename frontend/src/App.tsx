import React from 'react';
import { BrowserRouter, Routes, Route, Navigate, Outlet } from 'react-router-dom';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { Toaster } from 'react-hot-toast';
import { LoginPage } from './pages/LoginPage';
import { POSPage } from './pages/POSPage';
import { DashboardPage } from './pages/DashboardPage';
import { ReportPage } from './pages/ReportPage';
import { MenuPage } from './pages/MenuPage';
import { AdminPage } from './pages/AdminPage';
import { BranchPage } from './pages/BranchPage';
import { DepartmentPage } from './pages/DepartmentPage';
import { DaftarKaryawanPage } from './pages/DaftarKaryawanPage';
import { KomisiPage } from './pages/KomisiPage';
import { GajiPage } from './pages/GajiPage';
import KasPage from './pages/KasPage';
import HutangPage from './pages/HutangPage';
import PiutangPage from './pages/PiutangPage';
import PenarikanPage from './pages/PenarikanPage';
import SupplierPage from './pages/SupplierPage';
import PesananPembelianPage from './pages/PesananPembelianPage';
import PenerimaanBarangPage from './pages/PenerimaanBarangPage';
import FakturPembelianPage from './pages/FakturPembelianPage';
import PembayaranPembelianPage from './pages/PembayaranPembelianPage';
import ReturPembelianPage from './pages/ReturPembelianPage';
import CustomerCategoryPage from './pages/CustomerCategoryPage';
import SalesCategoryPage from './pages/SalesCategoryPage';
import CustomerPage from './pages/CustomerPage';
import PenawaranPenjualanPage from './pages/PenawaranPenjualanPage';
import PesananPenjualanPage from './pages/PesananPenjualanPage';
import PengirimanPesananPage from './pages/PengirimanPesananPage';
import FakturPenjualanPage from './pages/FakturPenjualanPage';
import PenerimaanPenjualanPage from './pages/PenerimaanPenjualanPage';
import ReturPenjualanPage from './pages/ReturPenjualanPage';
import WarehousePage from './pages/WarehousePage';
import ItemPage from './pages/ItemPage';
import TransferPage from './pages/TransferPage';
import StockAdjustmentPage from './pages/StockAdjustmentPage';
import PriceChangePage from './pages/PriceChangePage';
import TransferBankPage from './pages/TransferBankPage';
import RekonsiliasiBankPage from './pages/RekonsiliasiBankPage';
import ChartOfAccountsPage from './pages/ChartOfAccountsPage';
import JournalVoucherPage from './pages/JournalVoucherPage';
import PeriodEndPage from './pages/PeriodEndPage';
import PencatatanAsetPage from './pages/PencatatanAsetPage';
import PenyusutanOtomatisPage from './pages/PenyusutanOtomatisPage';
import SmartLinkECommercePage from './pages/SmartLinkECommercePage';
import SmartLinkEBankingPage from './pages/SmartLinkEBankingPage';
import SmartLinkTaxPage from './pages/SmartLinkTaxPage';
import SmartLinkManufacturingPage from './pages/SmartLinkManufacturingPage';
import { Layout } from './components/Layout';
import { useStore } from './store/useStore';

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      retry: 1,
      staleTime: 30000,
      refetchOnWindowFocus: false,
    },
  },
});

function PrivateRoute() {
  const user = useStore((s) => s.user);
  return user ? <Outlet /> : <Navigate to="/login" replace />;
}

function AdminRoute() {
  const user = useStore((s) => s.user);
  if (!user) return <Navigate to="/login" replace />;
  if (user.role !== 'admin') return <Navigate to="/pos" replace />;
  return <Outlet />;
}

function PermissionRoute({ permission }: { permission: string }) {
  const user = useStore((s) => s.user);
  if (!user) return <Navigate to="/login" replace />;
  if (user.role !== 'admin' && !(user.permissions || []).includes(permission)) {
    return <Navigate to="/pos" replace />;
  }
  return <Outlet />;
}

export default function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <BrowserRouter>
        <Toaster
          position="top-right"
          toastOptions={{
            style: {
              background: '#1f2937',
              color: '#f3f4f6',
              border: '1px solid #374151',
            },
          }}
        />
        <Routes>
          <Route path="/login" element={<LoginPage />} />
          <Route element={<Layout />}>
            <Route element={<PrivateRoute />}>
              <Route path="/" element={<Navigate to="/pos" replace />} />
              <Route path="pos" element={<POSPage />} />
              <Route element={<PermissionRoute permission="pos_access" />}>
                <Route path="dashboard" element={<DashboardPage />} />
              </Route>
              <Route element={<PermissionRoute permission="report_view" />}>
                <Route path="reports" element={<ReportPage />} />
              </Route>
              <Route element={<PermissionRoute permission="menu_manage" />}>
                <Route path="menu" element={<MenuPage />} />
              </Route>
            </Route>
            <Route element={<AdminRoute />}>
              <Route path="admin" element={<AdminPage />} />
              <Route path="branches" element={<BranchPage />} />
              <Route path="employees" element={<DaftarKaryawanPage />} />
              <Route path="departments" element={<DepartmentPage />} />
              <Route path="commissions" element={<KomisiPage />} />
              <Route path="salaries" element={<GajiPage />} />
              <Route path="cash" element={<KasPage />} />
              <Route path="debts" element={<HutangPage />} />
              <Route path="receivables" element={<PiutangPage />} />
              <Route path="withdrawals" element={<PenarikanPage />} />
              <Route path="bank-transfers" element={<TransferBankPage />} />
              <Route path="bank-reconciliations" element={<RekonsiliasiBankPage />} />
              <Route path="chart-of-accounts" element={<ChartOfAccountsPage />} />
              <Route path="journal-vouchers" element={<JournalVoucherPage />} />
              <Route path="period-end" element={<PeriodEndPage />} />
              <Route path="fixed-assets" element={<PencatatanAsetPage />} />
              <Route path="depreciation" element={<PenyusutanOtomatisPage />} />
              <Route path="suppliers" element={<SupplierPage />} />
              <Route path="purchase-orders" element={<PesananPembelianPage />} />
              <Route path="purchase-receive" element={<PenerimaanBarangPage />} />
              <Route path="purchase-invoices" element={<FakturPembelianPage />} />
              <Route path="purchase-payments" element={<PembayaranPembelianPage />} />
              <Route path="purchase-returns" element={<ReturPembelianPage />} />
              <Route path="customer-categories" element={<CustomerCategoryPage />} />
              <Route path="sales-categories" element={<SalesCategoryPage />} />
              <Route path="customers" element={<CustomerPage />} />
              <Route path="sales-quotations" element={<PenawaranPenjualanPage />} />
              <Route path="sales-orders" element={<PesananPenjualanPage />} />
              <Route path="delivery-orders" element={<PengirimanPesananPage />} />
              <Route path="sales-invoices" element={<FakturPenjualanPage />} />
              <Route path="sales-receipts" element={<PenerimaanPenjualanPage />} />
              <Route path="sales-returns" element={<ReturPenjualanPage />} />
            </Route>
            <Route element={<PermissionRoute permission="inventory_manage" />}>
              <Route path="warehouses" element={<WarehousePage />} />
              <Route path="items" element={<ItemPage />} />
              <Route path="stock-transfers" element={<TransferPage />} />
              <Route path="stock-adjustments" element={<StockAdjustmentPage />} />
              <Route path="price-changes" element={<PriceChangePage />} />
            </Route>
            <Route element={<PermissionRoute permission="smartlink_manage" />}>
              <Route path="smartlink/ecommerce" element={<SmartLinkECommercePage />} />
              <Route path="smartlink/ebanking" element={<SmartLinkEBankingPage />} />
              <Route path="smartlink/tax" element={<SmartLinkTaxPage />} />
            </Route>
            <Route element={<PermissionRoute permission="manufacturing_manage" />}>
              <Route path="smartlink/manufacturing" element={<SmartLinkManufacturingPage />} />
            </Route>
          </Route>
        </Routes>
      </BrowserRouter>
    </QueryClientProvider>
  );
}

import axios from 'axios';

const api = axios.create({
  baseURL: '/api/v1',
  timeout: 10000,
  headers: { 'Content-Type': 'application/json' },
});

api.interceptors.request.use((config) => {
  const stored = localStorage.getItem('pos-store');
  if (stored) {
    try {
      const state = JSON.parse(stored);
      if (state.state?.token) {
        config.headers.Authorization = `Bearer ${state.state.token}`;
      }
    } catch {}
  }
  return config;
});

api.interceptors.response.use(
  (res) => res,
  (err) => {
    if (err.response?.status === 401) {
      localStorage.removeItem('pos-store');
      window.location.href = '/login';
    }
    return Promise.reject(err);
  }
);

export interface Product {
  id: string;
  code: string;
  name: string;
  description: string;
  base_price: number;
  unit: string;
  image_url: string;
  category_id: string;
  has_variants: boolean;
  tax_rate: number;
  sort_order: number;
  stock: number;
  barcode: string;
  purchase_price: number;
  min_stock: number;
  is_service: boolean;
}

export interface Category {
  id: string;
  name: string;
  icon: string;
  color: string;
  sort_order: number;
}

export interface ProductVariant {
  id: string;
  name: string;
  group_name: string;
  price_adj: number;
}

export interface OrderResult {
  order: any;
  receipt: any;
}

// Product APIs
export const productService = {
  list: () =>
    api.get<{ data: Product[] }>('/products'),
  search: (outletId: string, q: string) =>
    api.get<{ data: Product[] }>('/products/search', { params: { outlet_id: outletId, q } }),
  barcode: (outletId: string, code: string) =>
    api.get<{ data: Product }>('/products/barcode', { params: { outlet_id: outletId, code } }),
  byCategory: (outletId: string, categoryId: string) =>
    api.get<{ data: Product[] }>(`/products/category/${categoryId}`, { params: { outlet_id: outletId } }),
  create: (data: any) =>
    api.post<{ data: Product }>('/products', data),
  update: (id: string, data: any) =>
    api.put<{ data: Product }>(`/products/${id}`, data),
  delete: (id: string) =>
    api.delete(`/products/${id}`),
  lowStock: (threshold?: number) =>
    api.get<{ data: Product[] }>('/products/low-stock', { params: threshold ? { threshold } : {} }),
};

export const availabilityService = {
  list: (date: string) =>
    api.get<{ data: any[] }>('/availability', { params: { date } }),
  get: (productId: string, date: string) =>
    api.get<{ data: any }>(`/availability/${productId}`, { params: { date } }),
  set: (productId: string, date: string, isAvailable: boolean) =>
    api.post<{ data: any }>('/availability', { product_id: productId, date, is_available: isAvailable }),
};

// Category APIs
export const categoryService = {
  list: (outletId: string) =>
    api.get<{ data: Category[] }>('/categories', { params: { outlet_id: outletId } }),
  create: (name: string) =>
    api.post<{ data: Category }>('/categories', { name }),
};

// Order APIs
export const orderService = {
  create: (data: any) =>
    api.post<{ success: boolean; data: OrderResult }>('/orders', data),
  sync: (data: any) =>
    api.post<{ success: boolean; data: any }>('/sync/orders', data),
  get: (id: string) => api.get(`/orders/${id}`),
  recent: (limit: number = 10) =>
    api.get<{ data: RecentOrder[] }>('/orders/recent', { params: { limit } }),
};

export interface RecentOrder {
  id: number;
  order_number: string;
  total: number;
  payment_method: string;
  cashier_name: string;
  item_count: number;
  status: string;
  created_at: string;
}

export interface RevenueBucket {
  label: string;
  sales: number;
  transactions: number;
}

export interface TopProduct {
  name: string;
  qty: number;
  revenue: number;
}

export interface RevenueSummary {
  period: string;
  total_sales: number;
  total_transactions: number;
  products_sold: number;
  avg_per_tx: number;
  change_pct: number;
  buckets: RevenueBucket[];
  top_products: TopProduct[];
  recent_orders: RecentOrder[];
}

export const revenueService = {
  summary: (period: 'daily' | 'weekly' | 'monthly' = 'daily') =>
    api.get<{ data: RevenueSummary }>('/revenue/summary', { params: { period } }),
};

// Report APIs
export interface Outlet {
  id: number;
  name: string;
}

export interface ReportKPIs {
  gross_sales: number;
  net_sales: number;
  gross_profit: number;
  transactions: number;
  avg_sale_per_tx: number;
  gross_margin_pct: number;
  items_sold: number;
  prev_gross_sales: number;
  change_pct: number;
}

export interface SalesPoint {
  label: string;
  gross_sales: number;
  net_sales: number;
  transactions: number;
}

export interface DayOfWeekSales {
  day: string;
  gross_sales: number;
  transactions: number;
}

export interface HourlySales {
  hour: string;
  gross_sales: number;
  transactions: number;
}

export interface ReportItem {
  name: string;
  items_sold: number;
  gross_sales: number;
  net_sales: number;
  gross_profit: number;
}

export interface CategoryStat {
  category: string;
  value: number;
  pct: number;
}

export interface TopItem {
  name: string;
  qty: number;
  gross_sales: number;
}

export interface CategoryTopItems {
  category: string;
  items: TopItem[];
}

export interface ReportSummaryData {
  period: string;
  range_start: string;
  range_end: string;
  bucket_mode: 'hourly' | 'daily' | 'monthly';
  kpis: ReportKPIs;
  daily_sales: SalesPoint[];
  day_of_week_sales: DayOfWeekSales[];
  hourly_sales: HourlySales[];
  items: ReportItem[];
  category_by_volume: CategoryStat[];
  category_by_sales: CategoryStat[];
  top_items_by_category: CategoryTopItems[];
}

export const reportService = {
  outlets: () => api.get<{ success: boolean; data: Outlet[] }>('/reports/outlets'),
  summary: (params: { period: string; outlet_id?: number }) =>
    api.get<{ success: boolean; data: ReportSummaryData }>('/reports/summary', { params }),
};

// Shift APIs
export const shiftService = {
  open: (data: { outlet_id: string; user_id: string; cash_start: number }) =>
    api.post('/shifts/open', data),
  close: (id: string, data: { cash_end: number; notes: string }) =>
    api.post(`/shifts/${id}/close`, data),
  active: (outletId: string, userId: string) =>
    api.get('/shifts/active', { params: { outlet_id: outletId, user_id: userId } }),
};

// Stock APIs
export const stockService = {
  get: (productId: string, outletId: string) =>
    api.get(`/stock/${productId}`, { params: { outlet_id: outletId } }),
  lowStock: (outletId: string) =>
    api.get('/stock/low', { params: { outlet_id: outletId } }),
};

// Payment APIs
export const paymentService = {
  generateQRIS: (amount: number) =>
    api.post('/payments/qris', { amount }),
};

// Auth APIs
export const authService = {
  login: (data: { username: string; password: string; tenant_id: string; outlet_id: string }) =>
    api.post('/auth/login', data),
};

// Admin APIs
export interface User {
  id: number;
  username: string;
  name: string;
  role: string;
  outlet_id: string;
  tenant_id: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface Permission {
  id: number;
  name: string;
  description: string;
}

export interface Branch {
  id: number;
  name: string;
  address: string;
  phone: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface Department {
  id: number;
  name: string;
  description: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface Employee {
  id: number;
  nik: string;
  name: string;
  email: string;
  phone: string;
  department_id: number | null;
  dept_name: string;
  position: string;
  hire_date: string;
  salary_base: number;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface Commission {
  id: number;
  employee_id: number;
  employee_name: string;
  commission_type: string;
  rate: number;
  description: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface Salary {
  id: number;
  employee_id: number;
  employee_name: string;
  pay_period: string;
  base_salary: number;
  commission_total: number;
  deduction: number;
  net_salary: number;
  status: string;
  paid_at: string | null;
  created_at: string;
  updated_at: string;
}

export interface CashTransaction {
  id: number;
  type: string;
  amount: number;
  description: string;
  category: string;
  reference: string;
  created_at: string;
}

export interface Debt {
  id: number;
  creditor_name: string;
  amount: number;
  paid_amount: number;
  remaining_amount: number;
  description: string;
  due_date: string | null;
  status: string;
  created_at: string;
  updated_at: string;
}

export interface Receivable {
  id: number;
  debtor_name: string;
  amount: number;
  paid_amount: number;
  remaining_amount: number;
  description: string;
  due_date: string | null;
  status: string;
  created_at: string;
  updated_at: string;
}

export interface Withdrawal {
  id: number;
  amount: number;
  description: string;
  method: string;
  status: string;
  requested_by: number | null;
  approved_by: number | null;
  approved_at: string | null;
  created_at: string;
  updated_at: string;
}

export interface BankTransfer {
  id: number;
  transfer_number: string;
  from_account_name: string;
  from_account_type: string;
  to_account_name: string;
  to_account_type: string;
  amount: number;
  transfer_date: string;
  notes: string;
  status: string;
  created_by: number | null;
  created_at: string;
  updated_at: string;
}

export interface BankReconciliationItem {
  id?: number;
  reconciliation_id?: number;
  description: string;
  transaction_date: string;
  amount: number;
  type: string;
  is_matched: boolean;
}

export interface BankReconciliation {
  id: number;
  reconciliation_number: string;
  account_name: string;
  account_type: string;
  period_start: string;
  period_end: string;
  book_balance: number;
  statement_balance: number;
  difference: number;
  notes: string;
  status: string;
  created_by: number | null;
  created_at: string;
  updated_at: string;
  items: BankReconciliationItem[];
}

export interface TransactionLogItem {
  product_name: string;
  quantity: number;
  price: number;
  subtotal: number;
}

export interface Supplier {
  id: number;
  name: string;
  contact_name: string;
  phone: string;
  email: string;
  npwp: string;
  address: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface PurchaseOrderItem {
  id: number;
  purchase_order_id: number;
  product_name: string;
  quantity: number;
  unit: string;
  unit_price: number;
  subtotal: number;
  received_qty: number;
}

export interface PurchaseOrder {
  id: number;
  po_number: string;
  supplier_id: number;
  supplier_name: string;
  order_date: string;
  expected_date: string | null;
  status: string;
  notes: string;
  total_amount: number;
  created_by: number | null;
  created_at: string;
  updated_at: string;
  items: PurchaseOrderItem[];
}

export interface PurchaseInvoice {
  id: number;
  invoice_number: string;
  purchase_order_id: number;
  po_number: string;
  supplier_id: number;
  supplier_name: string;
  invoice_date: string;
  due_date: string | null;
  total_amount: number;
  paid_amount: number;
  remaining_amount: number;
  status: string;
  notes: string;
  created_at: string;
  updated_at: string;
}

export interface PurchasePayment {
  id: number;
  payment_number: string;
  purchase_invoice_id: number;
  invoice_number: string;
  supplier_id: number;
  supplier_name: string;
  payment_date: string;
  amount: number;
  payment_method: string;
  notes: string;
  created_by: number | null;
  created_at: string;
}

export interface PurchaseReturnItem {
  id: number;
  purchase_return_id: number;
  product_name: string;
  quantity: number;
  unit_price: number;
  subtotal: number;
  reason: string;
}

export interface PurchaseReturn {
  id: number;
  return_number: string;
  purchase_order_id: number;
  po_number: string;
  supplier_id: number;
  supplier_name: string;
  return_date: string;
  reason: string;
  status: string;
  notes: string;
  created_by: number | null;
  created_at: string;
  updated_at: string;
  items: PurchaseReturnItem[];
}

export interface CustomerCategory {
  id: number;
  name: string;
  description: string;
  price_scheme: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface SalesCategory {
  id: number;
  name: string;
  description: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface Customer {
  id: number;
  name: string;
  phone: string;
  email: string;
  npwp: string;
  address: string;
  customer_category_id: number | null;
  customer_category_name: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface SalesQuotationItem {
  id: number;
  quotation_id: number;
  product_name: string;
  quantity: number;
  unit: string;
  unit_price: number;
  subtotal: number;
}

export interface SalesQuotation {
  id: number;
  quotation_number: string;
  customer_id: number;
  customer_name: string;
  sales_category_id: number | null;
  sales_category_name: string;
  quotation_date: string;
  valid_until: string | null;
  status: string;
  notes: string;
  total_amount: number;
  created_by: number | null;
  created_at: string;
  updated_at: string;
  items: SalesQuotationItem[];
}

export interface SalesOrderItem {
  id: number;
  sales_order_id: number;
  product_name: string;
  quantity: number;
  unit: string;
  unit_price: number;
  subtotal: number;
  delivered_qty: number;
}

export interface SalesOrder {
  id: number;
  order_number: string;
  customer_id: number;
  customer_name: string;
  quotation_id: number | null;
  sales_category_id: number | null;
  sales_category_name: string;
  order_date: string;
  expected_date: string | null;
  status: string;
  notes: string;
  total_amount: number;
  created_by: number | null;
  created_at: string;
  updated_at: string;
  items: SalesOrderItem[];
}

export interface DeliveryOrderItem {
  id: number;
  delivery_order_id: number;
  product_name: string;
  quantity: number;
  unit: string;
}

export interface DeliveryOrder {
  id: number;
  do_number: string;
  sales_order_id: number;
  order_number: string;
  customer_id: number;
  customer_name: string;
  delivery_date: string;
  status: string;
  notes: string;
  created_by: number | null;
  created_at: string;
  updated_at: string;
  items: DeliveryOrderItem[];
}

export interface SalesInvoice {
  id: number;
  invoice_number: string;
  sales_order_id: number;
  order_number: string;
  customer_id: number;
  customer_name: string;
  invoice_date: string;
  due_date: string | null;
  total_amount: number;
  paid_amount: number;
  remaining_amount: number;
  status: string;
  notes: string;
  created_at: string;
  updated_at: string;
}

export interface SalesReceipt {
  id: number;
  receipt_number: string;
  sales_invoice_id: number;
  invoice_number: string;
  customer_id: number;
  customer_name: string;
  receipt_date: string;
  amount: number;
  payment_method: string;
  notes: string;
  created_by: number | null;
  created_at: string;
}

export interface SalesReturnItem {
  id: number;
  sales_return_id: number;
  product_name: string;
  quantity: number;
  unit_price: number;
  subtotal: number;
  reason: string;
}

export interface SalesReturn {
  id: number;
  return_number: string;
  sales_order_id: number;
  order_number: string;
  customer_id: number;
  customer_name: string;
  return_date: string;
  reason: string;
  status: string;
  notes: string;
  total_amount: number;
  created_by: number | null;
  created_at: string;
  updated_at: string;
  items: SalesReturnItem[];
}

export interface TransactionLog {
  id: number;
  order_id: number;
  amount: number;
  payment_method: string;
  status: string;
  created_at: string;
  order_total: number;
  order_status: string;
  order_created: string;
  user_name: string;
  user_role: string;
  shift_number: number;
  shift_opened: string;
  items: TransactionLogItem[];
}

export interface ChartOfAccount {
  id: number;
  code: string;
  name: string;
  category: string;
  normal_balance: string;
  description: string;
  is_active: boolean;
  balance: number;
  created_at: string;
  updated_at: string;
}

export interface JournalEntryItem {
  id?: number;
  account_id: number;
  account_code: string;
  account_name: string;
  debit: number;
  credit: number;
}

export interface JournalEntry {
  id: number;
  entry_number: string;
  entry_date: string;
  description: string;
  reference: string;
  status: string;
  created_by: number | null;
  created_at: string;
  updated_at: string;
  items: JournalEntryItem[];
}

export interface FixedAsset {
  id: number;
  asset_code: string;
  name: string;
  asset_type: string;
  purchase_date: string;
  cost: number;
  salvage_value: number;
  useful_life_months: number;
  depreciation_method: string;
  accumulated_depreciation: number;
  depreciation_account_id: number;
  depreciation_account_code: string;
  depreciation_account_name: string;
  accumulated_account_id: number;
  accumulated_account_code: string;
  accumulated_account_name: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface AssetDepreciationSchedule extends FixedAsset {
  monthly_depreciation: number;
  first_depreciation: number;
  remaining_value: number;
  months_remaining: number;
}

export interface FxDifference {
  id: number;
  period: string;
  currency: string;
  description: string;
  gain_amount: number;
  loss_amount: number;
  created_at: string;
}

export interface PeriodClosing {
  id: number;
  period: string;
  status: string;
  depreciation_total: number;
  fx_gain: number;
  fx_loss: number;
  notes: string;
  closed_by: number | null;
  closed_at: string | null;
  created_at: string;
  updated_at: string;
}

export interface DepreciationResult {
  period: string;
  assets_count: number;
  depreciation_total: number;
}

export const adminService = {
  listUsers: () => api.get<{ data: User[] }>('/admin/users'),
  getUser: (id: number) => api.get<{ data: User }>(`/admin/users/${id}`),
  createUser: (data: { username: string; name: string; password: string; role: string }) =>
    api.post('/admin/users', data),
  updateUser: (id: number, data: { name?: string; role?: string; is_active?: boolean }) =>
    api.put(`/admin/users/${id}`, data),
  updatePassword: (id: number, password: string) =>
    api.put(`/admin/users/${id}/password`, { password }),
  deleteUser: (id: number) => api.delete(`/admin/users/${id}`),
  getUserPermissions: (id: number) => api.get<{ data: Permission[] }>(`/admin/users/${id}/permissions`),
  setUserPermissions: (id: number, permissionIds: number[]) =>
    api.put(`/admin/users/${id}/permissions`, { permission_ids: permissionIds }),
  getAllPermissions: () => api.get<{ data: Permission[] }>('/admin/permissions'),
  getTransactionLogsByUsername: (username: string) => api.get<{ data: TransactionLog[] }>('/admin/transactions/log', { params: { username } }),
  upload: (file: File) => {
    const fd = new FormData();
    fd.append('file', file);
    return api.post<{ data: { url: string; filename: string } }>('/admin/upload', fd, {
      headers: { 'Content-Type': 'multipart/form-data' },
    });
  },
  listBranches: () => api.get<{ data: Branch[] }>('/admin/branches'),
  getBranch: (id: number) => api.get<{ data: Branch }>(`/admin/branches/${id}`),
  createBranch: (data: { name: string; address: string; phone: string }) =>
    api.post<{ data: Branch }>('/admin/branches', data),
  updateBranch: (id: number, data: { name?: string; address?: string; phone?: string; is_active?: boolean }) =>
    api.put<{ data: Branch }>(`/admin/branches/${id}`, data),
  deleteBranch: (id: number) => api.delete(`/admin/branches/${id}`),
  listDepartments: () => api.get<{ data: Department[] }>('/admin/departments'),
  getDepartment: (id: number) => api.get<{ data: Department }>(`/admin/departments/${id}`),
  createDepartment: (data: { name: string; description: string }) =>
    api.post<{ data: Department }>('/admin/departments', data),
  updateDepartment: (id: number, data: { name?: string; description?: string; is_active?: boolean }) =>
    api.put<{ data: Department }>(`/admin/departments/${id}`, data),
  deleteDepartment: (id: number) => api.delete(`/admin/departments/${id}`),
  listEmployees: () => api.get<{ data: Employee[] }>('/admin/employees'),
  getEmployee: (id: number) => api.get<{ data: Employee }>(`/admin/employees/${id}`),
  createEmployee: (data: { nik: string; name: string; email: string; phone: string; department_id: number | null; position: string; hire_date: string; salary_base: number }) =>
    api.post<{ data: Employee }>('/admin/employees', data),
  updateEmployee: (id: number, data: any) =>
    api.put<{ data: Employee }>(`/admin/employees/${id}`, data),
  deleteEmployee: (id: number) => api.delete(`/admin/employees/${id}`),
  listCommissions: () => api.get<{ data: Commission[] }>('/admin/commissions'),
  getCommission: (id: number) => api.get<{ data: Commission }>(`/admin/commissions/${id}`),
  createCommission: (data: { employee_id: number; commission_type: string; rate: number; description: string }) =>
    api.post<{ data: Commission }>('/admin/commissions', data),
  updateCommission: (id: number, data: any) =>
    api.put<{ data: Commission }>(`/admin/commissions/${id}`, data),
  deleteCommission: (id: number) => api.delete(`/admin/commissions/${id}`),
  listSalaries: () => api.get<{ data: Salary[] }>('/admin/salaries'),
  getSalary: (id: number) => api.get<{ data: Salary }>(`/admin/salaries/${id}`),
  createSalary: (data: { employee_id: number; pay_period: string; base_salary: number; commission_total: number; deduction: number }) =>
    api.post<{ data: Salary }>('/admin/salaries', data),
  updateSalary: (id: number, data: any) =>
    api.put<{ data: Salary }>(`/admin/salaries/${id}`, data),
  deleteSalary: (id: number) => api.delete(`/admin/salaries/${id}`),
  markSalaryPaid: (id: number) => api.put(`/admin/salaries/${id}/pay`),
  listCash: () => api.get<{ data: CashTransaction[] }>('/admin/cash'),
  getCash: (id: number) => api.get<{ data: CashTransaction }>(`/admin/cash/${id}`),
  createCash: (data: { type: string; amount: number; description: string; category: string; reference: string }) =>
    api.post<{ data: CashTransaction }>('/admin/cash', data),
  deleteCash: (id: number) => api.delete(`/admin/cash/${id}`),
  listDebts: () => api.get<{ data: Debt[] }>('/admin/debts'),
  getDebt: (id: number) => api.get<{ data: Debt }>(`/admin/debts/${id}`),
  createDebt: (data: { creditor_name: string; amount: number; description: string; due_date: string | null }) =>
    api.post<{ data: Debt }>('/admin/debts', data),
  updateDebt: (id: number, data: any) => api.put<{ data: Debt }>(`/admin/debts/${id}`, data),
  deleteDebt: (id: number) => api.delete(`/admin/debts/${id}`),
  listReceivables: () => api.get<{ data: Receivable[] }>('/admin/receivables'),
  getReceivable: (id: number) => api.get<{ data: Receivable }>(`/admin/receivables/${id}`),
  createReceivable: (data: { debtor_name: string; amount: number; description: string; due_date: string | null }) =>
    api.post<{ data: Receivable }>('/admin/receivables', data),
  updateReceivable: (id: number, data: any) => api.put<{ data: Receivable }>(`/admin/receivables/${id}`, data),
  deleteReceivable: (id: number) => api.delete(`/admin/receivables/${id}`),
  listWithdrawals: () => api.get<{ data: Withdrawal[] }>('/admin/withdrawals'),
  getWithdrawal: (id: number) => api.get<{ data: Withdrawal }>(`/admin/withdrawals/${id}`),
  createWithdrawal: (data: { amount: number; description: string; method: string }) =>
    api.post<{ data: Withdrawal }>('/admin/withdrawals', data),
  updateWithdrawal: (id: number, data: any) => api.put<{ data: Withdrawal }>(`/admin/withdrawals/${id}`, data),
  deleteWithdrawal: (id: number) => api.delete(`/admin/withdrawals/${id}`),
  approveWithdrawal: (id: number) => api.put(`/admin/withdrawals/${id}/approve`),
  rejectWithdrawal: (id: number) => api.put(`/admin/withdrawals/${id}/reject`),
  listBankTransfers: () => api.get<{ data: BankTransfer[] }>('/admin/bank-transfers'),
  getBankTransfer: (id: number) => api.get<{ data: BankTransfer }>(`/admin/bank-transfers/${id}`),
  createBankTransfer: (data: { from_account_name: string; from_account_type: string; to_account_name: string; to_account_type: string; amount: number; transfer_date: string; notes: string }) =>
    api.post<{ data: BankTransfer }>('/admin/bank-transfers', data),
  updateBankTransfer: (id: number, data: any) => api.put<{ data: BankTransfer }>(`/admin/bank-transfers/${id}`, data),
  deleteBankTransfer: (id: number) => api.delete(`/admin/bank-transfers/${id}`),
  listBankReconciliations: () => api.get<{ data: BankReconciliation[] }>('/admin/bank-reconciliations'),
  getBankReconciliation: (id: number) => api.get<{ data: BankReconciliation }>(`/admin/bank-reconciliations/${id}`),
  createBankReconciliation: (data: { account_name: string; account_type: string; period_start: string; period_end: string; book_balance: number; statement_balance: number; notes: string; items: BankReconciliationItem[] }) =>
    api.post<{ data: BankReconciliation }>('/admin/bank-reconciliations', data),
  deleteBankReconciliation: (id: number) => api.delete(`/admin/bank-reconciliations/${id}`),
  listChartOfAccounts: () => api.get<{ data: ChartOfAccount[] }>('/admin/chart-of-accounts'),
  getChartOfAccount: (id: number) => api.get<{ data: ChartOfAccount }>(`/admin/chart-of-accounts/${id}`),
  createChartOfAccount: (data: { code: string; name: string; category: string; normal_balance: string; description: string; is_active: boolean }) =>
    api.post<{ data: ChartOfAccount }>('/admin/chart-of-accounts', data),
  updateChartOfAccount: (id: number, data: any) => api.put<{ data: ChartOfAccount }>(`/admin/chart-of-accounts/${id}`, data),
  deleteChartOfAccount: (id: number) => api.delete(`/admin/chart-of-accounts/${id}`),
  listJournalEntries: () => api.get<{ data: JournalEntry[] }>('/admin/journal-entries'),
  getJournalEntry: (id: number) => api.get<{ data: JournalEntry }>(`/admin/journal-entries/${id}`),
  createJournalEntry: (data: { entry_date: string; description: string; reference: string; items: { account_id: number; debit: number; credit: number }[] }) =>
    api.post<{ data: JournalEntry }>('/admin/journal-entries', data),
  deleteJournalEntry: (id: number) => api.delete(`/admin/journal-entries/${id}`),
  listFixedAssets: () => api.get<{ data: FixedAsset[] }>('/admin/fixed-assets'),
  getFixedAsset: (id: number) => api.get<{ data: FixedAsset }>(`/admin/fixed-assets/${id}`),
  listDepreciationSchedule: () => api.get<{ data: AssetDepreciationSchedule[] }>('/admin/fixed-assets/schedule'),
  createFixedAsset: (data: any) => api.post<{ data: FixedAsset }>('/admin/fixed-assets', data),
  updateFixedAsset: (id: number, data: any) => api.put<{ data: FixedAsset }>(`/admin/fixed-assets/${id}`, data),
  deleteFixedAsset: (id: number) => api.delete(`/admin/fixed-assets/${id}`),
  listPeriodEnd: () => api.get<{ data: PeriodClosing[] }>('/admin/period-end'),
  getPeriodEnd: (id: number) => api.get<{ data: PeriodClosing }>(`/admin/period-end/${id}`),
  openPeriod: (period: string) => api.post('/admin/period-end/open', { period }),
  listFxDifferences: (id: number) => api.get<{ data: FxDifference[] }>(`/admin/period-end/${id}/fx-differences`),
  createFxDifference: (data: { period: string; currency: string; description: string; gain_amount: number; loss_amount: number }) =>
    api.post<{ data: FxDifference }>('/admin/period-end/fx-differences', data),
  runDepreciation: (id: number) => api.post<{ data: DepreciationResult }>(`/admin/period-end/${id}/depreciation`),
  closePeriod: (id: number) => api.post(`/admin/period-end/${id}/close`),
  listSuppliers: () => api.get<{ data: Supplier[] }>('/admin/suppliers'),
  getSupplier: (id: number) => api.get<{ data: Supplier }>(`/admin/suppliers/${id}`),
  createSupplier: (data: { name: string; contact_name: string; phone: string; email: string; address: string }) =>
    api.post<{ data: Supplier }>('/admin/suppliers', data),
  updateSupplier: (id: number, data: any) => api.put<{ data: Supplier }>(`/admin/suppliers/${id}`, data),
  deleteSupplier: (id: number) => api.delete(`/admin/suppliers/${id}`),
  listPurchaseOrders: () => api.get<{ data: PurchaseOrder[] }>('/admin/purchase-orders'),
  getPurchaseOrder: (id: number) => api.get<{ data: PurchaseOrder }>(`/admin/purchase-orders/${id}`),
  createPurchaseOrder: (data: any) => api.post<{ data: PurchaseOrder }>('/admin/purchase-orders', data),
  updatePurchaseOrder: (id: number, data: any) => api.put<{ data: PurchaseOrder }>(`/admin/purchase-orders/${id}`, data),
  deletePurchaseOrder: (id: number) => api.delete(`/admin/purchase-orders/${id}`),
  receivePurchaseOrderItems: (id: number, items: { id: number; quantity: number }[], warehouseId?: number) =>
    api.put(`/admin/purchase-orders/${id}/receive`, { items, warehouse_id: warehouseId }),
  listPurchaseInvoices: () => api.get<{ data: PurchaseInvoice[] }>('/admin/purchase-invoices'),
  getPurchaseInvoice: (id: number) => api.get<{ data: PurchaseInvoice }>(`/admin/purchase-invoices/${id}`),
  createPurchaseInvoice: (data: any) => api.post<{ data: PurchaseInvoice }>('/admin/purchase-invoices', data),
  updatePurchaseInvoice: (id: number, data: any) => api.put<{ data: PurchaseInvoice }>(`/admin/purchase-invoices/${id}`, data),
  deletePurchaseInvoice: (id: number) => api.delete(`/admin/purchase-invoices/${id}`),
  listPurchasePayments: () => api.get<{ data: PurchasePayment[] }>('/admin/purchase-payments'),
  getPurchasePayment: (id: number) => api.get<{ data: PurchasePayment }>(`/admin/purchase-payments/${id}`),
  createPurchasePayment: (data: any) => api.post<{ data: PurchasePayment }>('/admin/purchase-payments', data),
  deletePurchasePayment: (id: number) => api.delete(`/admin/purchase-payments/${id}`),
  listPurchaseReturns: () => api.get<{ data: PurchaseReturn[] }>('/admin/purchase-returns'),
  getPurchaseReturn: (id: number) => api.get<{ data: PurchaseReturn }>(`/admin/purchase-returns/${id}`),
  createPurchaseReturn: (data: any) => api.post<{ data: PurchaseReturn }>('/admin/purchase-returns', data),
  updatePurchaseReturn: (id: number, data: any) => api.put<{ data: PurchaseReturn }>(`/admin/purchase-returns/${id}`, data),
  deletePurchaseReturn: (id: number) => api.delete(`/admin/purchase-returns/${id}`),
  listCustomerCategories: () => api.get<{ data: CustomerCategory[] }>('/admin/customer-categories'),
  getCustomerCategory: (id: number) => api.get<{ data: CustomerCategory }>(`/admin/customer-categories/${id}`),
  createCustomerCategory: (data: { name: string; description: string; price_scheme: string }) =>
    api.post<{ data: CustomerCategory }>('/admin/customer-categories', data),
  updateCustomerCategory: (id: number, data: any) => api.put<{ data: CustomerCategory }>(`/admin/customer-categories/${id}`, data),
  deleteCustomerCategory: (id: number) => api.delete(`/admin/customer-categories/${id}`),
  listSalesCategories: () => api.get<{ data: SalesCategory[] }>('/admin/sales-categories'),
  getSalesCategory: (id: number) => api.get<{ data: SalesCategory }>(`/admin/sales-categories/${id}`),
  createSalesCategory: (data: { name: string; description: string }) =>
    api.post<{ data: SalesCategory }>('/admin/sales-categories', data),
  updateSalesCategory: (id: number, data: any) => api.put<{ data: SalesCategory }>(`/admin/sales-categories/${id}`, data),
  deleteSalesCategory: (id: number) => api.delete(`/admin/sales-categories/${id}`),
  listCustomers: () => api.get<{ data: Customer[] }>('/admin/customers'),
  getCustomer: (id: number) => api.get<{ data: Customer }>(`/admin/customers/${id}`),
  createCustomer: (data: { name: string; phone: string; email: string; address: string; customer_category_id: number | null }) =>
    api.post<{ data: Customer }>('/admin/customers', data),
  updateCustomer: (id: number, data: any) => api.put<{ data: Customer }>(`/admin/customers/${id}`, data),
  deleteCustomer: (id: number) => api.delete(`/admin/customers/${id}`),
  listSalesQuotations: () => api.get<{ data: SalesQuotation[] }>('/admin/sales-quotations'),
  getSalesQuotation: (id: number) => api.get<{ data: SalesQuotation }>(`/admin/sales-quotations/${id}`),
  createSalesQuotation: (data: any) => api.post<{ data: SalesQuotation }>('/admin/sales-quotations', data),
  updateSalesQuotation: (id: number, data: any) => api.put<{ data: SalesQuotation }>(`/admin/sales-quotations/${id}`, data),
  deleteSalesQuotation: (id: number) => api.delete(`/admin/sales-quotations/${id}`),
  listSalesOrders: () => api.get<{ data: SalesOrder[] }>('/admin/sales-orders'),
  getSalesOrder: (id: number) => api.get<{ data: SalesOrder }>(`/admin/sales-orders/${id}`),
  createSalesOrder: (data: any) => api.post<{ data: SalesOrder }>('/admin/sales-orders', data),
  updateSalesOrder: (id: number, data: any) => api.put<{ data: SalesOrder }>(`/admin/sales-orders/${id}`, data),
  deleteSalesOrder: (id: number) => api.delete(`/admin/sales-orders/${id}`),
  listDeliveryOrders: () => api.get<{ data: DeliveryOrder[] }>('/admin/delivery-orders'),
  getDeliveryOrder: (id: number) => api.get<{ data: DeliveryOrder }>(`/admin/delivery-orders/${id}`),
  createDeliveryOrder: (data: any) => api.post<{ data: DeliveryOrder }>('/admin/delivery-orders', data),
  updateDeliveryOrder: (id: number, data: any) => api.put<{ data: DeliveryOrder }>(`/admin/delivery-orders/${id}`, data),
  deleteDeliveryOrder: (id: number) => api.delete(`/admin/delivery-orders/${id}`),
  listSalesInvoices: () => api.get<{ data: SalesInvoice[] }>('/admin/sales-invoices'),
  getSalesInvoice: (id: number) => api.get<{ data: SalesInvoice }>(`/admin/sales-invoices/${id}`),
  createSalesInvoice: (data: any) => api.post<{ data: SalesInvoice }>('/admin/sales-invoices', data),
  updateSalesInvoice: (id: number, data: any) => api.put<{ data: SalesInvoice }>(`/admin/sales-invoices/${id}`, data),
  deleteSalesInvoice: (id: number) => api.delete(`/admin/sales-invoices/${id}`),
  listSalesReceipts: () => api.get<{ data: SalesReceipt[] }>('/admin/sales-receipts'),
  getSalesReceipt: (id: number) => api.get<{ data: SalesReceipt }>(`/admin/sales-receipts/${id}`),
  createSalesReceipt: (data: any) => api.post<{ data: SalesReceipt }>('/admin/sales-receipts', data),
  deleteSalesReceipt: (id: number) => api.delete(`/admin/sales-receipts/${id}`),
  listSalesReturns: () => api.get<{ data: SalesReturn[] }>('/admin/sales-returns'),
  getSalesReturn: (id: number) => api.get<{ data: SalesReturn }>(`/admin/sales-returns/${id}`),
  createSalesReturn: (data: any) => api.post<{ data: SalesReturn }>('/admin/sales-returns', data),
  updateSalesReturn: (id: number, data: any) => api.put<{ data: SalesReturn }>(`/admin/sales-returns/${id}`, data),
  deleteSalesReturn: (id: number) => api.delete(`/admin/sales-returns/${id}`),
};

export interface Warehouse {
  id: number;
  code: string;
  name: string;
  branch_id: number | null;
  branch_name: string;
  address: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface ItemVariant {
  id?: number;
  name: string;
  barcode: string;
  additional_price: number;
  stock: number;
  is_active?: boolean;
}

export interface ItemStock {
  warehouse_id: number;
  warehouse_name: string;
  quantity: number;
}

export interface Item {
  id: number;
  code: string;
  name: string;
  barcode: string;
  description: string;
  category_id: number | null;
  category_name: string;
  purchase_price: number;
  base_price: number;
  unit: string;
  min_stock: number;
  is_service: boolean;
  has_variants: boolean;
  tax_rate: number;
  stock: number;
  initial_stock?: number;
  image_url: string;
  variants: ItemVariant[];
  stock_by_warehouse: ItemStock[];
}

export interface StockTransferItem {
  id?: number;
  product_id: number;
  product_name: string;
  quantity: number;
  unit: string;
}

export interface StockTransfer {
  id: number;
  transfer_number: string;
  from_warehouse_id: number;
  from_warehouse_name: string;
  to_warehouse_id: number;
  to_warehouse_name: string;
  transfer_date: string;
  status: string;
  notes: string;
  created_by: number | null;
  created_at: string;
  updated_at: string;
  items: StockTransferItem[];
}

export interface StockAdjustmentItem {
  id?: number;
  product_id: number;
  product_name: string;
  system_qty: number;
  actual_qty: number;
  difference: number;
  reason: string;
}

export interface StockAdjustment {
  id: number;
  adjustment_number: string;
  warehouse_id: number;
  warehouse_name: string;
  adjustment_date: string;
  status: string;
  notes: string;
  created_by: number | null;
  created_at: string;
  updated_at: string;
  items: StockAdjustmentItem[];
}

export interface PriceChangeItem {
  id?: number;
  product_id: number;
  product_name: string;
  old_price: number;
  new_price: number;
}

export interface PriceChange {
  id: number;
  price_change_number: string;
  change_date: string;
  status: string;
  notes: string;
  created_by: number | null;
  created_at: string;
  items: PriceChangeItem[];
}

export interface StockMovement {
  id: number;
  product_id: number;
  product_name: string;
  warehouse_id: number;
  warehouse_name: string;
  quantity: number;
  movement_type: string;
  reference_type: string;
  reference_id: number | null;
  note: string;
  created_by: number | null;
  created_at: string;
}

export const inventoryService = {
  listWarehouses: () => api.get<{ data: Warehouse[] }>('/admin/warehouses'),
  getWarehouse: (id: number) => api.get<{ data: Warehouse }>(`/admin/warehouses/${id}`),
  createWarehouse: (data: any) => api.post<{ data: Warehouse }>('/admin/warehouses', data),
  updateWarehouse: (id: number, data: any) => api.put<{ data: Warehouse }>(`/admin/warehouses/${id}`, data),
  deleteWarehouse: (id: number) => api.delete(`/admin/warehouses/${id}`),
  warehouseStock: (id: number) => api.get<{ data: ItemStock[] }>(`/admin/warehouses/${id}/stock`),
  listItems: () => api.get<{ data: Item[] }>('/admin/items'),
  getItem: (id: number) => api.get<{ data: Item }>(`/admin/items/${id}`),
  createItem: (data: any) => api.post<{ data: Product }>('/admin/items', data),
  updateItem: (id: number, data: any) => api.put<{ data: Product }>(`/admin/items/${id}`, data),
  deleteItem: (id: number) => api.delete(`/admin/items/${id}`),
  listTransfers: () => api.get<{ data: StockTransfer[] }>('/admin/stock-transfers'),
  getTransfer: (id: number) => api.get<{ data: StockTransfer }>(`/admin/stock-transfers/${id}`),
  createTransfer: (data: any) => api.post<{ data: StockTransfer }>('/admin/stock-transfers', data),
  updateTransfer: (id: number, data: any) => api.put<{ data: StockTransfer }>(`/admin/stock-transfers/${id}`, data),
  sendTransfer: (id: number) => api.put(`/admin/stock-transfers/${id}/send`),
  receiveTransfer: (id: number) => api.put(`/admin/stock-transfers/${id}/receive`),
  cancelTransfer: (id: number) => api.put(`/admin/stock-transfers/${id}/cancel`),
  deleteTransfer: (id: number) => api.delete(`/admin/stock-transfers/${id}`),
  listAdjustments: () => api.get<{ data: StockAdjustment[] }>('/admin/stock-adjustments'),
  getAdjustment: (id: number) => api.get<{ data: StockAdjustment }>(`/admin/stock-adjustments/${id}`),
  createAdjustment: (data: any) => api.post<{ data: StockAdjustment }>('/admin/stock-adjustments', data),
  deleteAdjustment: (id: number) => api.delete(`/admin/stock-adjustments/${id}`),
  listPriceChanges: () => api.get<{ data: PriceChange[] }>('/admin/price-changes'),
  getPriceChange: (id: number) => api.get<{ data: PriceChange }>(`/admin/price-changes/${id}`),
  createPriceChange: (data: any) => api.post<{ data: PriceChange }>('/admin/price-changes', data),
  deletePriceChange: (id: number) => api.delete(`/admin/price-changes/${id}`),
  listMovements: () => api.get<{ data: StockMovement[] }>('/admin/stock-movements'),
};

export interface MarketplaceConnection {
  id: number;
  platform: string;
  shop_name: string;
  api_token: string;
  account_name: string;
  customer_id: number | null;
  customer_name: string;
  sales_category_id: number | null;
  sales_category_name: string;
  shipping_fee_account: string;
  commission_account: string;
  is_active: boolean;
  last_sync_at: string | null;
  created_at: string;
  updated_at: string;
}

export interface MarketplaceOrder {
  id: number;
  connection_id: number;
  platform: string;
  shop_name: string;
  marketplace_order_id: string;
  order_date: string;
  customer_name: string;
  product_id: number | null;
  product_name: string;
  quantity: number;
  unit_price: number;
  subtotal: number;
  shipping_fee: number;
  platform_fee: number;
  grand_total: number;
  status: string;
  sales_order_id: number | null;
  sales_order_number: string;
  created_at: string;
}

export interface BankStatementLine {
  id?: number;
  import_id?: number;
  transaction_date: string;
  description: string;
  reference: string;
  amount: number;
  type: string;
  is_reconciled?: boolean;
  reconciliation_id?: number | null;
}

export interface BankStatementImport {
  id: number;
  import_number: string;
  account_name: string;
  account_type: string;
  statement_date: string;
  source: string;
  file_name: string;
  total_rows: number;
  status: string;
  created_at: string;
  lines: BankStatementLine[];
}

export interface EFakturLine {
  id?: number;
  export_id?: number;
  reference_number: string;
  party_name: string;
  npwp: string;
  line_date: string | null;
  taxable_base: number;
  tax_amount: number;
  status: string;
}

export interface EFakturExport {
  id: number;
  export_number: string;
  export_type: string;
  period: string;
  status: string;
  total_rows: number;
  total_dpp: number;
  total_tax: number;
  file_path: string;
  created_at: string;
  lines: EFakturLine[];
}

export interface BomItem {
  id?: number;
  bom_id?: number;
  product_id: number;
  product_name: string;
  quantity_required: number;
  unit: string;
  cost_per_unit: number;
  estimated_cost: number;
}

export interface BillOfMaterials {
  id: number;
  bom_number: string;
  product_id: number;
  product_name: string;
  quantity_output: number;
  unit: string;
  overhead_cost: number;
  cost_per_unit: number;
  status: string;
  notes: string;
  created_at: string;
  updated_at: string;
  items: BomItem[];
}

export interface WorkOrderItem {
  id?: number;
  work_order_id?: number;
  product_id: number;
  product_name: string;
  quantity_required: number;
  unit: string;
  cost_per_unit: number;
  subtotal: number;
}

export interface WorkOrder {
  id: number;
  work_order_number: string;
  bom_id: number;
  bom_number: string;
  product_id: number;
  product_name: string;
  quantity: number;
  warehouse_id: number | null;
  warehouse_name: string;
  scheduled_date: string;
  status: string;
  actual_cost: number;
  notes: string;
  created_at: string;
  updated_at: string;
  items: WorkOrderItem[];
}

export const smartlinkService = {
  // ---- e-Commerce / Marketplace ----
  listMarketplaceConnections: () => api.get<{ data: MarketplaceConnection[] }>('/admin/smartlink/marketplace/connections'),
  createMarketplaceConnection: (data: any) => api.post<{ data: MarketplaceConnection }>('/admin/smartlink/marketplace/connections', data),
  updateMarketplaceConnection: (id: number, data: any) => api.put<{ data: MarketplaceConnection }>(`/admin/smartlink/marketplace/connections/${id}`, data),
  deleteMarketplaceConnection: (id: number) => api.delete(`/admin/smartlink/marketplace/connections/${id}`),
  listMarketplaceOrders: () => api.get<{ data: MarketplaceOrder[] }>('/admin/smartlink/marketplace/orders'),
  importMarketplaceOrders: (connectionId: number, orders: any[]) =>
    api.post<{ data: { imported: number } }>('/admin/smartlink/marketplace/orders/import', { connection_id: connectionId, orders }),
  postMarketplaceOrder: (id: number) => api.post(`/admin/smartlink/marketplace/orders/${id}/post`),
  deleteMarketplaceOrder: (id: number) => api.delete(`/admin/smartlink/marketplace/orders/${id}`),

  // ---- e-Banking ----
  listBankImports: () => api.get<{ data: BankStatementImport[] }>('/admin/smartlink/bank/imports'),
  importBankLines: (data: { account_name: string; account_type: string; statement_date: string; source: string; file_name: string; lines: BankStatementLine[] }) =>
    api.post<{ data: BankStatementImport }>('/admin/smartlink/bank/imports', data),
  reconcileBankImport: (id: number, bookBalance: number) =>
    api.post(`/admin/smartlink/bank/imports/${id}/reconcile`, { book_balance: bookBalance }),
  deleteBankImport: (id: number) => api.delete(`/admin/smartlink/bank/imports/${id}`),

  // ---- e-Faktur / Pajak ----
  listTaxExports: () => api.get<{ data: EFakturExport[] }>('/admin/smartlink/tax/exports'),
  generateTaxExport: (exportType: string, period: string) =>
    api.post<{ data: EFakturExport }>('/admin/smartlink/tax/exports/generate', { export_type: exportType, period }),
  deleteTaxExport: (id: number) => api.delete(`/admin/smartlink/tax/exports/${id}`),

  // ---- Manufaktur: BOM ----
  listBoms: () => api.get<{ data: BillOfMaterials[] }>('/admin/manufacturing/boms'),
  getBom: (id: number) => api.get<{ data: BillOfMaterials }>(`/admin/manufacturing/boms/${id}`),
  createBom: (data: any) => api.post<{ data: BillOfMaterials }>('/admin/manufacturing/boms', data),
  updateBom: (id: number, data: any) => api.put<{ data: BillOfMaterials }>(`/admin/manufacturing/boms/${id}`, data),
  recalculateBom: (id: number) => api.post<{ data: BillOfMaterials }>(`/admin/manufacturing/boms/${id}/recalculate`),
  deleteBom: (id: number) => api.delete(`/admin/manufacturing/boms/${id}`),

  // ---- Manufaktur: Work Order ----
  listWorkOrders: () => api.get<{ data: WorkOrder[] }>('/admin/manufacturing/work-orders'),
  getWorkOrder: (id: number) => api.get<{ data: WorkOrder }>(`/admin/manufacturing/work-orders/${id}`),
  createWorkOrder: (data: any) => api.post<{ data: WorkOrder }>('/admin/manufacturing/work-orders', data),
  startWorkOrder: (id: number) => api.put(`/admin/manufacturing/work-orders/${id}/start`),
  completeWorkOrder: (id: number) => api.put(`/admin/manufacturing/work-orders/${id}/complete`),
  cancelWorkOrder: (id: number) => api.put(`/admin/manufacturing/work-orders/${id}/cancel`),
  deleteWorkOrder: (id: number) => api.delete(`/admin/manufacturing/work-orders/${id}`),
};

export default api;

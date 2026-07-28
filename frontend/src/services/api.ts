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
};

// Order APIs
export const orderService = {
  create: (data: any) =>
    api.post<{ success: boolean; data: OrderResult }>('/orders', data),
  get: (id: string) => api.get(`/orders/${id}`),
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

export default api;

import { create } from 'zustand';
import { persist } from 'zustand/middleware';

export interface CartItem {
  id: string;
  productId: string;
  name: string;
  code: string;
  price: number;
  quantity: number;
  variantLabel: string;
  note: string;
  customizations: { name: string; price: number }[];
  subtotal: number;
}

export interface OrderState {
  // Auth
  user: { id: string; name: string; role: string; outletId: string; tenantId: string } | null;
  outletName: string;
  token: string;

  // Cart
  items: CartItem[];
  orderType: 'dine_in' | 'takeaway';
  tableNumber: string;
  customerName: string;
  customerPhone: string;
  discountCode: string;
  discountAmount: number;

  // Shift
  activeShift: { id: string; shiftNumber: number; cashStart: number } | null;

  // Actions
  setUser: (user: any) => void;
  logout: () => void;
  addItem: (item: Omit<CartItem, 'id' | 'subtotal'>) => void;
  removeItem: (id: string) => void;
  updateQuantity: (id: string, qty: number) => void;
  clearCart: () => void;
  setOrderType: (type: 'dine_in' | 'takeaway') => void;
  setTableNumber: (n: string) => void;
  setCustomerName: (n: string) => void;
  setCustomerPhone: (p: string) => void;
  setDiscountCode: (c: string) => void;
  setDiscountAmount: (a: number) => void;
  setActiveShift: (s: any) => void;
  setOutletName: (n: string) => void;
}

export const useStore = create<OrderState>()(
  persist(
    (set, get) => ({
      user: null,
      outletName: '',
      token: '',
      items: [],
      orderType: 'dine_in',
      tableNumber: '',
      customerName: '',
      customerPhone: '',
      discountCode: '',
      discountAmount: 0,
      activeShift: null,

      setUser: (user) => set({ user }),
      setToken: (token) => set({ token }),
      logout: () => set({ user: null, token: '', items: [], activeShift: null }),

      addItem: (item) =>
        set((state) => {
          const existing = state.items.find(
            (i) => i.productId === item.productId && i.variantLabel === item.variantLabel
          );
          if (existing) {
            return {
              items: state.items.map((i) =>
                i.id === existing.id
                  ? { ...i, quantity: i.quantity + 1, subtotal: (i.quantity + 1) * i.price }
                  : i
              ),
            };
          }
          return {
            items: [
              ...state.items,
              {
                ...item,
                id: crypto.randomUUID(),
                subtotal: item.price * item.quantity,
              },
            ],
          };
        }),

      removeItem: (id) =>
        set((state) => ({ items: state.items.filter((i) => i.id !== id) })),

      updateQuantity: (id, qty) =>
        set((state) => ({
          items: state.items.map((i) =>
            i.id === id ? { ...i, quantity: Math.max(0, qty), subtotal: Math.max(0, qty) * i.price } : i
          ).filter((i) => i.quantity > 0),
        })),

      clearCart: () =>
        set({
          items: [],
          tableNumber: '',
          customerName: '',
          customerPhone: '',
          discountCode: '',
          discountAmount: 0,
        }),

      setOrderType: (orderType) => set({ orderType }),
      setTableNumber: (tableNumber) => set({ tableNumber }),
      setCustomerName: (customerName) => set({ customerName }),
      setCustomerPhone: (customerPhone) => set({ customerPhone }),
      setDiscountCode: (discountCode) => set({ discountCode }),
      setDiscountAmount: (discountAmount) => set({ discountAmount }),
      setActiveShift: (activeShift) => set({ activeShift }),
      setOutletName: (outletName) => set({ outletName }),
    }),
    { name: 'pos-store' }
  )
);

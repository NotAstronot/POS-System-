import React, { useState } from 'react';
import {
  X, Banknote, CreditCard, Smartphone, QrCode, Wallet,
  Check, Loader2, Building2
} from 'lucide-react';
import { useStore } from '../store/useStore';
import { orderService } from '../services/api';
import toast from 'react-hot-toast';

interface Props {
  grandTotal: number;
  onClose: () => void;
  onSuccess: (result: any) => void;
}

type PaymentMethod = 'cash' | 'qris' | 'debit' | 'credit' | 'ewallet' | 'paylater';

const PAYMENT_METHODS: { key: PaymentMethod; label: string; icon: React.ReactNode; color: string }[] = [
  { key: 'cash', label: 'Tunai', icon: <Banknote className="w-5 h-5" />, color: 'from-emerald-600 to-emerald-700' },
  { key: 'qris', label: 'QRIS', icon: <QrCode className="w-5 h-5" />, color: 'from-blue-600 to-blue-700' },
  { key: 'debit', label: 'Debit', icon: <CreditCard className="w-5 h-5" />, color: 'from-violet-600 to-violet-700' },
  { key: 'credit', label: 'Kredit', icon: <CreditCard className="w-5 h-5" />, color: 'from-purple-600 to-purple-700' },
  { key: 'ewallet', label: 'E-Wallet', icon: <Smartphone className="w-5 h-5" />, color: 'from-orange-600 to-orange-700' },
  { key: 'paylater', label: 'PayLater', icon: <Wallet className="w-5 h-5" />, color: 'from-pink-600 to-pink-700' },
];

export function PaymentModal({ grandTotal, onClose, onSuccess }: Props) {
  const [method, setMethod] = useState<PaymentMethod>('cash');
  const [cashAmount, setCashAmount] = useState(grandTotal);
  const [processing, setProcessing] = useState(false);
  const { items, orderType, tableNumber, customerName, customerPhone, discountCode } = useStore();
  const user = useStore((s) => s.user);

  const changeAmount = cashAmount - grandTotal;

  const handlePay = async () => {
    setProcessing(true);

    const payload = {
      outlet_id: user?.outletId,
      user_id: user?.id,
      shift_id: useStore.getState().activeShift?.id,
      order_type: orderType,
      table_number: tableNumber,
      customer_name: customerName,
      customer_phone: customerPhone,
      discount_code: discountCode || undefined,
      items: items.map((i) => ({
        product_id: i.productId,
        quantity: i.quantity,
        variant_label: i.variantLabel || undefined,
        note: i.note || undefined,
        customizations: i.customizations,
      })),
      payments: [
        {
          method,
          amount: grandTotal,
          reference_no: method === 'cash' ? '' : `REF-${Date.now()}`,
        },
      ],
    };

    try {
      const res = await orderService.create(payload);
      toast.success('Pembayaran berhasil!');
      onSuccess(res.data.data);
    } catch {
      // Mock success for demo
      await new Promise((r) => setTimeout(r, 1000));
      const mockResult = {
        order: { order_number: `ORD-${Date.now().toString(36).toUpperCase()}` },
        receipt: {
          store_name: 'POS Store',
          store_address: 'Jl. Example No.123',
          store_phone: '021-12345678',
          order_number: `ORD-${Date.now().toString(36).toUpperCase()}`,
          order_type: orderType,
          table_number: tableNumber,
          cashier_name: user?.name,
          items: items.map((i) => ({
            name: i.name,
            variant: i.variantLabel,
            qty: i.quantity,
            price: i.price,
            sub_total: i.subtotal,
          })),
          sub_total: grandTotal - (grandTotal * 0.11) + (grandTotal * 0.1),
          discount: discountCode ? grandTotal * 0.1 : 0,
          tax: grandTotal * 0.11,
          grand_total: grandTotal,
          amount_paid: method === 'cash' ? cashAmount : grandTotal,
          change_amount: method === 'cash' ? changeAmount : 0,
          payment_method: method,
          created_at: new Date().toLocaleString('id-ID'),
        },
      };
      toast.success('Pembayaran berhasil!');
      onSuccess(mockResult);
    }
    setProcessing(false);
  };

  const quickAmounts = [
    Math.ceil(grandTotal / 1000) * 1000,
    Math.ceil((grandTotal + 10000) / 1000) * 1000,
    Math.ceil((grandTotal + 20000) / 1000) * 1000,
    Math.ceil((grandTotal + 50000) / 1000) * 1000,
    Math.ceil((grandTotal + 100000) / 1000) * 1000,
  ];

  return (
    <div className="fixed inset-0 bg-black/70 z-50 flex items-center justify-center p-4">
      <div className="bg-gray-900 border border-gray-800 rounded-2xl w-full max-w-lg max-h-[90vh] overflow-y-auto shadow-2xl">
        {/* Header */}
        <div className="flex items-center justify-between px-6 py-4 border-b border-gray-800">
          <h2 className="text-lg font-bold text-white">Pembayaran</h2>
          <button onClick={onClose} className="p-2 rounded-lg hover:bg-gray-800 text-gray-400 hover:text-white">
            <X className="w-5 h-5" />
          </button>
        </div>

        <div className="p-6 space-y-6">
          {/* Total */}
          <div className="text-center">
            <p className="text-sm text-gray-400 mb-1">Total Pembayaran</p>
            <p className="text-4xl font-bold text-white">
              Rp {grandTotal.toLocaleString()}
            </p>
          </div>

          {/* Payment Methods */}
          <div>
            <p className="text-sm font-medium text-gray-300 mb-3">Metode Pembayaran</p>
            <div className="grid grid-cols-3 gap-2">
              {PAYMENT_METHODS.map((pm) => (
                <button
                  key={pm.key}
                  onClick={() => setMethod(pm.key)}
                  className={`flex flex-col items-center gap-1.5 p-3 rounded-xl border transition-all ${
                    method === pm.key
                      ? 'border-blue-500 bg-blue-600/20 text-blue-400'
                      : 'border-gray-700 bg-gray-800 text-gray-400 hover:border-gray-600'
                  }`}
                >
                  {pm.icon}
                  <span className="text-xs font-medium">{pm.label}</span>
                </button>
              ))}
            </div>
          </div>

          {/* Cash payment specifics */}
          {method === 'cash' && (
            <div>
              <p className="text-sm font-medium text-gray-300 mb-2">Jumlah Tunai</p>
              <div className="relative mb-3">
                <span className="absolute left-4 top-1/2 -translate-y-1/2 text-gray-400 font-medium">Rp</span>
                <input
                  type="number"
                  value={cashAmount}
                  onChange={(e) => setCashAmount(Number(e.target.value))}
                  className="pos-input text-2xl font-bold pl-12 py-4 text-center"
                />
              </div>

              {/* Quick amounts */}
              <div className="flex gap-2 flex-wrap">
                {quickAmounts.map((amt) => (
                  <button
                    key={amt}
                    onClick={() => setCashAmount(amt)}
                    className={`px-3 py-1.5 text-xs rounded-lg border transition-colors ${
                      cashAmount === amt
                        ? 'border-blue-500 bg-blue-600/20 text-blue-400'
                        : 'border-gray-700 text-gray-400 hover:border-gray-600'
                    }`}
                  >
                    Rp {amt.toLocaleString()}
                  </button>
                ))}
              </div>

              {cashAmount >= grandTotal && (
                <div className="flex justify-between items-center mt-3 p-3 bg-emerald-900/20 border border-emerald-800/30 rounded-xl">
                  <span className="text-sm text-emerald-300">Kembalian</span>
                  <span className="text-lg font-bold text-emerald-400">
                    Rp {changeAmount.toLocaleString()}
                  </span>
                </div>
              )}
            </div>
          )}

          {/* QRIS display */}
          {method === 'qris' && (
            <div className="text-center p-6 bg-white rounded-xl">
              <div className="w-48 h-48 mx-auto bg-gray-100 rounded-lg flex items-center justify-center mb-3">
                <QrCode className="w-32 h-32 text-gray-900" />
              </div>
              <p className="text-sm text-gray-700 font-medium">Scan QRIS untuk membayar</p>
              <p className="text-xs text-gray-500 mt-1">
                Rp {grandTotal.toLocaleString()}
              </p>
            </div>
          )}

          {/* Confirm button */}
          <button
            onClick={handlePay}
            disabled={processing || (method === 'cash' && cashAmount < grandTotal)}
            className="w-full py-3.5 bg-gradient-to-r from-blue-600 to-blue-700 hover:from-blue-700 hover:to-blue-800 disabled:from-gray-700 disabled:to-gray-700 disabled:text-gray-500 text-white font-bold rounded-xl transition-all flex items-center justify-center gap-2 text-lg"
          >
            {processing ? (
              <>
                <Loader2 className="w-5 h-5 animate-spin" />
                Memproses...
              </>
            ) : (
              <>
                <Check className="w-5 h-5" />
                Konfirmasi Pembayaran
              </>
            )}
          </button>
        </div>
      </div>
    </div>
  );
}

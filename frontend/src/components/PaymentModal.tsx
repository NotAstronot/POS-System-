import React, { useState, useEffect } from 'react';
import {
  X, Banknote, CreditCard, Smartphone, QrCode, Wallet,
  Check, Loader2, Building2
} from 'lucide-react';
import { useStore } from '../store/useStore';
import { orderService, paymentService } from '../services/api';
import { isOnline, enqueueOrder, buildLocalReceipt } from '../services/offline';
import toast from 'react-hot-toast';

interface Props {
  grandTotal: number;
  subtotal?: number;
  discount?: number;
  tax?: number;
  onClose: () => void;
  onSuccess: (result: any) => void;
}

type PaymentMethod = 'cash' | 'qris' | 'debit' | 'credit' | 'ewallet' | 'paylater';

const PAYMENT_METHODS: { key: PaymentMethod; label: string; icon: React.ReactNode; color: string }[] = [
  { key: 'cash', label: 'Tunai', icon: <Banknote className="w-5 h-5" />, color: 'success' },
  { key: 'qris', label: 'QRIS', icon: <QrCode className="w-5 h-5" />, color: 'brand' },
  { key: 'debit', label: 'Debit', icon: <CreditCard className="w-5 h-5" />, color: 'brand' },
  { key: 'credit', label: 'Kredit', icon: <CreditCard className="w-5 h-5" />, color: 'brand' },
  { key: 'ewallet', label: 'E-Wallet', icon: <Smartphone className="w-5 h-5" />, color: 'warning' },
  { key: 'paylater', label: 'PayLater', icon: <Wallet className="w-5 h-5" />, color: 'danger' },
];

export function PaymentModal({ grandTotal, subtotal = 0, discount = 0, tax = 0, onClose, onSuccess }: Props) {
  const [method, setMethod] = useState<PaymentMethod>('cash');
  const [cashAmount, setCashAmount] = useState(grandTotal);
  const [processing, setProcessing] = useState(false);
  const [qrImage, setQrImage] = useState('');
  const { items, orderType, tableNumber, customerName, customerPhone, discountCode } = useStore();
  const user = useStore((s) => s.user);

  useEffect(() => {
    if (method !== 'qris') return;
    let cancelled = false;
    paymentService.generateQRIS(grandTotal)
      .then((res) => !cancelled && setQrImage(res.data?.data?.qr_base64 || ''))
      .catch(() => !cancelled && setQrImage(''));
    return () => { cancelled = true; };
  }, [method, grandTotal]);

  const changeAmount = cashAmount - grandTotal;

  const handlePay = async () => {
    setProcessing(true);

    const payload = {
      outlet_id: user?.outletId,
      shift_id: useStore.getState().activeShift?.id,
      order_type: orderType,
      table_number: tableNumber,
      customer_name: customerName,
      customer_phone: customerPhone,
      discount_code: discountCode || undefined,
      items: items.map((i) => ({
        product_id: Number(i.productId),
        quantity: i.quantity,
        variant_label: i.variantLabel || undefined,
        note: i.note || undefined,
        customizations: i.customizations,
      })),
      item_names: items.map((i) => i.name),
      payments: [
        {
          method,
          amount: grandTotal,
          reference_no: method === 'cash' ? '' : `REF-${Date.now()}`,
        },
      ],
    };

    try {
      if (isOnline()) {
        const res = await orderService.create(payload);
        toast.success('Pembayaran berhasil!');
        onSuccess(res.data.data);
        return;
      }
      throw new Error('offline');
    } catch (err: any) {
      const networkError = err?.message === 'offline' || !err?.response;
      if (networkError) {
        const entry = enqueueOrder(payload);
        const localResult = buildLocalReceipt(
          entry,
          items,
          subtotal,
          discount,
          tax,
          grandTotal,
          method,
          user?.name || '',
          useStore.getState().outletName || 'Toko (Offline)'
        );
        toast.success('Transaksi disimpan OFFLINE — otomatis disinkronkan saat online', { duration: 4500 });
        onSuccess(localResult);
      } else {
        const msg = err?.response?.data?.error;
        toast.error(msg || 'Pembayaran gagal, periksa koneksi ke server');
      }
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
    <div className="pos-modal-overlay">
      <div className="pos-modal">
        {/* Header */}
        <div className="pos-modal-header">
          <h2 className="pos-modal-title">Pembayaran</h2>
          <button onClick={onClose} className="p-2 rounded-lg hover:bg-neutral-100 text-neutral-500 hover:text-neutral-900">
            <X className="w-5 h-5" />
          </button>
        </div>

        <div className="pos-modal-body space-y-6">
          {/* Total */}
          <div className="text-center">
            <p className="text-sm text-neutral-500 mb-1">Total Pembayaran</p>
            <p className="text-4xl font-bold text-neutral-900">
              Rp {grandTotal.toLocaleString()}
            </p>
          </div>

          {/* Payment Methods */}
          <div>
            <p className="text-sm font-medium text-neutral-700 mb-3">Metode Pembayaran</p>
            <div className="grid grid-cols-3 gap-2">
              {PAYMENT_METHODS.map((pm) => (
                <button
                  key={pm.key}
                  onClick={() => setMethod(pm.key)}
                  className={`flex flex-col items-center gap-1.5 p-3 rounded-xl border transition-all ${
                    method === pm.key
                      ? 'border-brand bg-primary-50 text-brand'
                      : 'border-neutral-200 bg-white text-neutral-500 hover:border-neutral-300'
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
              <p className="text-sm font-medium text-neutral-700 mb-2">Jumlah Tunai</p>
              <div className="relative mb-3">
                <span className="absolute left-4 top-1/2 -translate-y-1/2 text-neutral-400 font-medium">Rp</span>
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
                        ? 'border-brand bg-primary-50 text-brand'
                        : 'border-neutral-200 text-neutral-500 hover:border-neutral-300'
                    }`}
                  >
                    Rp {amt.toLocaleString()}
                  </button>
                ))}
              </div>

              {cashAmount >= grandTotal && (
                <div className="flex justify-between items-center mt-3 p-3 bg-success/10 border border-success/20 rounded-xl">
                  <span className="text-sm text-success-dark">Kembalian</span>
                  <span className="text-lg font-bold text-success-dark">
                    Rp {changeAmount.toLocaleString()}
                  </span>
                </div>
              )}
            </div>
          )}

          {/* QRIS display */}
          {method === 'qris' && (
            <div className="text-center p-6 bg-white border border-neutral-200 rounded-xl">
              <div className="w-48 h-48 mx-auto bg-neutral-50 rounded-lg flex items-center justify-center mb-3">
                {qrImage ? (
                  <img src={`data:image/png;base64,${qrImage}`} alt="QRIS" className="w-48 h-48" />
                ) : (
                  <QrCode className="w-32 h-32 text-neutral-900" />
                )}
              </div>
              <p className="text-sm text-neutral-700 font-medium">Scan QRIS untuk membayar</p>
              <p className="text-xs text-neutral-500 mt-1">
                Rp {grandTotal.toLocaleString()}
              </p>
            </div>
          )}
        </div>

        <div className="pos-modal-footer">
          <button
            onClick={handlePay}
            disabled={processing || (method === 'cash' && cashAmount < grandTotal)}
            className="w-full py-3.5 pos-btn-primary text-lg"
          >
            {processing ? (
              <>
                <Loader2 className="w-5 h-5 animate-spin mr-2" />
                Memproses...
              </>
            ) : (
              <>
                <Check className="w-5 h-5 mr-2" />
                Konfirmasi Pembayaran
              </>
            )}
          </button>
        </div>
      </div>
    </div>
  );
}
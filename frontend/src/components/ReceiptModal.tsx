import React, { useRef } from 'react';
import { X, Printer, Share2, Download, CheckCircle } from 'lucide-react';
import toast from 'react-hot-toast';

interface ReceiptItem {
  name: string;
  variant?: string;
  qty: number;
  price: number;
  sub_total: number;
}

interface ReceiptData {
  store_name: string;
  store_address: string;
  store_phone: string;
  order_number: string;
  order_type: string;
  table_number?: string;
  cashier_name?: string;
  items: ReceiptItem[];
  sub_total: number;
  discount: number;
  tax: number;
  grand_total: number;
  amount_paid: number;
  change_amount: number;
  payment_method: string;
  created_at: string;
}

interface Props {
  receipt: ReceiptData;
  onClose: () => void;
}

export function ReceiptModal({ receipt, onClose }: Props) {
  const printRef = useRef<HTMLDivElement>(null);

  const handlePrint = () => {
    window.print();
  };

  const handleShareWA = () => {
    const text = `*${receipt.store_name}*\n${receipt.store_address}\nTelp: ${receipt.store_phone}\n\n` +
      `No. Pesanan: ${receipt.order_number}\n` +
      `Tanggal: ${receipt.created_at}\n` +
      `Kasir: ${receipt.cashier_name || '-'}\n` +
      `${receipt.order_type === 'dine_in' ? `Meja: ${receipt.table_number || '-'}` : 'Takeaway'}\n\n` +
      `--- DAFTAR ITEM ---\n` +
      receipt.items.map((i) =>
        `  ${i.name}${i.variant ? ` (${i.variant})` : ''}  x${i.qty}  Rp ${i.sub_total.toLocaleString()}`
      ).join('\n') +
      `\n\nSub Total: Rp ${receipt.sub_total.toLocaleString()}\n` +
      `Diskon: Rp ${receipt.discount.toLocaleString()}\n` +
      `Pajak: Rp ${receipt.tax.toLocaleString()}\n` +
      `*Grand Total: Rp ${receipt.grand_total.toLocaleString()}*\n\n` +
      `Terima kasih telah berbelanja!`;

    const waUrl = `https://wa.me/?text=${encodeURIComponent(text)}`;
    window.open(waUrl, '_blank');
  };

  return (
    <div className="pos-modal-overlay">
      <div className="pos-modal max-w-md">
        {/* Action Bar */}
        <div className="pos-modal-header">
          <div className="flex items-center gap-2">
            <CheckCircle className="w-5 h-5 text-success" />
            <h2 className="pos-modal-title">Pembayaran Berhasil</h2>
          </div>
          <div className="flex items-center gap-2">
            <button onClick={handleShareWA} className="p-2 rounded-lg hover:bg-neutral-100 text-neutral-500 hover:text-success">
              <Share2 className="w-4 h-4" />
            </button>
            <button onClick={handlePrint} className="p-2 rounded-lg hover:bg-neutral-100 text-neutral-500 hover:text-brand">
              <Printer className="w-4 h-4" />
            </button>
            <button onClick={onClose} className="p-2 rounded-lg hover:bg-neutral-100 text-neutral-500 hover:text-neutral-900">
              <X className="w-5 h-5" />
            </button>
          </div>
        </div>

        {/* Receipt Content */}
        <div className="pos-modal-body" ref={printRef}>
          <div id="receipt-print-area" className="max-w-[80mm] mx-auto bg-white text-black p-4 rounded-lg">
            {/* Header */}
            <div className="text-center border-b border-dashed border-neutral-300 pb-3 mb-3">
              <h3 className="font-bold text-lg">{receipt.store_name}</h3>
              <p className="text-xs text-neutral-600">{receipt.store_address}</p>
              <p className="text-xs text-neutral-600">Telp: {receipt.store_phone}</p>
            </div>

            {/* Info */}
            <div className="text-xs space-y-0.5 mb-3">
              <div className="flex justify-between">
                <span className="text-neutral-600">No. Pesanan</span>
                <span className="font-medium">{receipt.order_number}</span>
              </div>
              <div className="flex justify-between">
                <span className="text-neutral-600">Tgl/Jam</span>
                <span>{receipt.created_at}</span>
              </div>
              <div className="flex justify-between">
                <span className="text-neutral-600">Kasir</span>
                <span>{receipt.cashier_name || '-'}</span>
              </div>
              <div className="flex justify-between">
                <span className="text-neutral-600">Tipe</span>
                <span>{receipt.order_type === 'dine_in' ? `Dine In (Meja ${receipt.table_number || '-'})` : 'Takeaway'}</span>
              </div>
            </div>

            {/* Items */}
            <div className="border-t border-dashed border-neutral-300 pt-2 mb-2">
              <div className="flex justify-between text-xs font-bold text-neutral-700 mb-1">
                <span className="flex-1">Item</span>
                <span className="w-10 text-right">Qty</span>
                <span className="w-20 text-right">Harga</span>
              </div>
              {receipt.items.map((item, idx) => (
                <div key={idx} className="text-xs py-1">
                  <div className="flex justify-between">
                    <span className="flex-1">{item.name}{item.variant ? ` (${item.variant})` : ''}</span>
                    <span className="w-10 text-right">{item.qty}</span>
                    <span className="w-20 text-right">Rp {item.sub_total.toLocaleString()}</span>
                  </div>
                </div>
              ))}
            </div>

            {/* Totals */}
            <div className="border-t border-dashed border-neutral-300 pt-2 space-y-0.5 text-xs">
              <div className="flex justify-between">
                <span>Sub Total</span>
                <span>Rp {receipt.sub_total.toLocaleString()}</span>
              </div>
              {receipt.discount > 0 && (
                <div className="flex justify-between text-danger">
                  <span>Diskon</span>
                  <span>-Rp {receipt.discount.toLocaleString()}</span>
                </div>
              )}
              <div className="flex justify-between">
                <span>Pajak (11%)</span>
                <span>Rp {receipt.tax.toLocaleString()}</span>
              </div>
              <div className="flex justify-between text-sm font-bold border-t border-neutral-300 pt-1 mt-1">
                <span>Grand Total</span>
                <span>Rp {receipt.grand_total.toLocaleString()}</span>
              </div>
              <div className="flex justify-between text-success">
                <span>Bayar ({receipt.payment_method})</span>
                <span>Rp {receipt.amount_paid.toLocaleString()}</span>
              </div>
              {receipt.change_amount > 0 && (
                <div className="flex justify-between text-success">
                  <span>Kembali</span>
                  <span>Rp {receipt.change_amount.toLocaleString()}</span>
                </div>
              )}
            </div>

            {/* Footer */}
            <div className="text-center text-xs text-neutral-600 mt-4 pt-3 border-t border-dashed border-neutral-300">
              <p className="font-bold">Terima Kasih!</p>
              <p>Barang yang sudah dibeli tidak dapat dikembalikan</p>
            </div>
          </div>
        </div>

        {/* Bottom actions */}
        <div className="pos-modal-footer flex-col sm:flex-row gap-3">
          <button onClick={onClose} className="w-full sm:flex-1 pos-btn-secondary text-sm py-3">
            Pesanan Baru
          </button>
          <button onClick={handlePrint} className="w-full sm:flex-1 pos-btn-primary text-sm py-3">
            <Printer className="w-4 h-4 mr-2" />
            Cetak Struk
          </button>
        </div>
      </div>
    </div>
  );
}
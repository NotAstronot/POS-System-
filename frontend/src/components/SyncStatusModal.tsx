import React, { useEffect, useState } from 'react';
import { X, RefreshCw, UploadCloud, AlertTriangle, CheckCircle2, Trash2, Wifi, WifiOff } from 'lucide-react';
import toast from 'react-hot-toast';
import {
  getQueue, getFailed, dismissFailed, isOnline,
  type OfflineOrder, type FailedOrder,
} from '../services/offline';

function orderSummary(o: OfflineOrder | FailedOrder) {
  const items = (o.payload?.items || []) as { product_id: number; quantity: number }[];
  const names = (o.payload?.item_names as string[] | undefined) || [];
  const label = items
    .map((it, i) => `${names[i] || `Produk #${it.product_id}`} x${it.quantity}`)
    .join(', ');
  return label || '-';
}

export function SyncStatusModal({ onClose, onSync, syncing, onCountsChange }: Props) {
  const [queue, setQueue] = useState<OfflineOrder[]>([]);
  const [failed, setFailed] = useState<FailedOrder[]>([]);
  const [online, setOnline] = useState(isOnline());

  const reload = () => {
    setQueue(getQueue());
    setFailed(getFailed());
  };

  useEffect(() => {
    reload();
    const on = () => setOnline(true);
    const off = () => setOnline(false);
    window.addEventListener('online', on);
    window.addEventListener('offline', off);
    return () => {
      window.removeEventListener('online', on);
      window.removeEventListener('offline', off);
    };
  }, []);

  const handleSync = async () => {
    await onSync();
    reload();
    onCountsChange();
  };

  const handleDismiss = (id: string) => {
    dismissFailed(id);
    reload();
    onCountsChange();
    toast('Order gagal dihapus dari daftar', { icon: '🗑️' });
  };

  const totalPending = queue.reduce((sum, o) => sum + (o.payload?.payments?.[0]?.amount || 0), 0);

  return (
    <div className="pos-modal-overlay">
      <div className="pos-modal max-w-lg max-h-[85vh] flex flex-col">
        {/* Header */}
        <div className="pos-modal-header">
          <div className="flex items-center gap-2">
            <UploadCloud className="w-5 h-5 text-brand" />
            <h2 className="pos-modal-title">Sinkronisasi Offline</h2>
          </div>
          <div className="flex items-center gap-2">
            <span className={`flex items-center gap-1.5 text-xs font-medium ${online ? 'text-success' : 'text-danger'}`}>
              {online ? <Wifi className="w-4 h-4" /> : <WifiOff className="w-4 h-4" />}
              {online ? 'Online' : 'Offline'}
            </span>
            <button onClick={onClose} className="p-2 rounded-lg hover:bg-neutral-100 text-neutral-500 hover:text-neutral-900">
              <X className="w-5 h-5" />
            </button>
          </div>
        </div>

        <div className="pos-modal-body space-y-5 overflow-y-auto">
          {/* Pending queue */}
          <div>
            <div className="flex items-center justify-between mb-2">
              <p className="text-sm font-semibold text-neutral-700">
                Antrian ({queue.length})
              </p>
              {queue.length > 0 && (
                <p className="text-xs text-neutral-500">
                  Total: <strong className="text-neutral-900">Rp {totalPending.toLocaleString()}</strong>
                </p>
              )}
            </div>
            {queue.length === 0 ? (
              <div className="flex items-center gap-2 p-3 rounded-xl bg-success/10 border border-success/20 text-success text-sm">
                <CheckCircle2 className="w-4 h-4" />
                Tidak ada order tertunda — semua sudah tersinkron.
              </div>
            ) : (
              <div className="space-y-2">
                {queue.map((o) => (
                  <div key={o.client_order_id} className="p-3 rounded-xl bg-neutral-50 border border-neutral-200 text-sm">
                    <div className="flex justify-between items-center gap-2">
                      <span className="font-mono text-xs text-brand">{o.client_order_id}</span>
                      <span className="text-xs text-neutral-500">{new Date(o.created_at).toLocaleString('id-ID')}</span>
                    </div>
                    <p className="text-neutral-700 mt-1">{orderSummary(o)}</p>
                  </div>
                ))}
              </div>
            )}
          </div>

          {/* Failed list */}
          {failed.length > 0 && (
            <div>
              <p className="text-sm font-semibold text-danger mb-2 flex items-center gap-1.5">
                <AlertTriangle className="w-4 h-4" />
                Ditolak Server ({failed.length})
              </p>
              <div className="space-y-2">
                {failed.map((f) => (
                  <div key={f.client_order_id} className="p-3 rounded-xl bg-danger/5 border border-danger/20 text-sm">
                    <div className="flex justify-between items-center gap-2">
                      <span className="font-mono text-xs text-danger">{f.client_order_id}</span>
                      <button
                        onClick={() => handleDismiss(f.client_order_id)}
                        className="p-1 text-neutral-500 hover:text-danger"
                        title="Hapus dari daftar"
                      >
                        <Trash2 className="w-3.5 h-3.5" />
                      </button>
                    </div>
                    <p className="text-neutral-700 mt-1">{orderSummary(f)}</p>
                    <p className="text-xs text-danger mt-1">{f.reason}</p>
                    <p className="text-[11px] text-neutral-500 mt-1">
                      Periksa stok/produk, lalu coba jual kembali secara manual.
                    </p>
                  </div>
                ))}
              </div>
            </div>
          )}

          <p className="text-xs text-neutral-500 leading-relaxed">
            Saat offline, transaksi disimpan di perangkat ini dan otomatis dikirim saat
            koneksi kembali. Server mengecek stok saat sync: order yang stoknya tidak
            mencukupi akan ditolak dan muncul di daftar ini.
          </p>
        </div>

        {/* Actions */}
        <div className="pos-modal-footer flex-col sm:flex-row gap-3">
          <button onClick={onClose} className="w-full sm:flex-1 pos-btn-secondary text-sm py-3">
            Tutup
          </button>
          <button
            onClick={handleSync}
            disabled={syncing || !online || queue.length === 0}
            className="w-full sm:flex-1 pos-btn-primary text-sm py-3 disabled:opacity-50"
          >
            {syncing ? (
              <>
                <RefreshCw className="w-4 h-4 animate-spin mr-2" />
                Menyinkronkan...
              </>
            ) : (
              <>
                <UploadCloud className="w-4 h-4 mr-2" />
                Sync Sekarang
              </>
            )}
          </button>
        </div>
      </div>
    </div>
  );
}

interface Props {
  onClose: () => void;
  onSync: () => Promise<void>;
  syncing: boolean;
  onCountsChange: () => void;
}
import React, { useState, useEffect } from 'react';
import { X, Clock, DollarSign, Calendar, CheckCircle, AlertTriangle } from 'lucide-react';
import { useStore } from '../store/useStore';
import { shiftService } from '../services/api';
import toast from 'react-hot-toast';

interface Props {
  onClose: () => void;
}

export function ShiftModal({ onClose }: Props) {
  const user = useStore((s) => s.user);
  const setActiveShift = useStore((s) => s.setActiveShift);
  const activeShift = useStore((s) => s.activeShift);
  const [cashStart, setCashStart] = useState(500000);
  const [cashEnd, setCashEnd] = useState(0);
  const [loading, setLoading] = useState(false);

  // Get active shift on mount
  useEffect(() => {
    if (!activeShift && user) {
      // Try to fetch existing active shift
      shiftService.active(user.outletId, user.id)
        .then((res) => {
          if (res.data.data) {
            setActiveShift(res.data.data);
          }
        })
        .catch(() => {});
    }
  }, [user]);

  const handleOpenShift = async () => {
    if (!user) return;
    setLoading(true);
    try {
      const res = await shiftService.open({
        outlet_id: user.outletId,
        user_id: user.id,
        cash_start: cashStart,
      });
      setActiveShift(res.data.data);
      toast.success(`Shift #${res.data.data.shift_number} berhasil dibuka!`);
      onClose();
    } catch (err: any) {
      toast.error(err?.response?.data?.error || 'Gagal membuka shift, periksa koneksi ke server');
    }
    setLoading(false);
  };

  const handleCloseShift = async () => {
    if (!activeShift) return;
    setLoading(true);
    try {
      await shiftService.close(activeShift.id, {
        cash_end: cashEnd,
        notes: '',
      });
      setActiveShift(null);
      toast.success('Shift berhasil ditutup!');
      onClose();
    } catch (err: any) {
      toast.error(err?.response?.data?.error || 'Gagal menutup shift');
    }
    setLoading(false);
  };

  if (activeShift) {
    const openedAt = new Date(activeShift.opened_at || Date.now());
    const duration = Math.floor((Date.now() - openedAt.getTime()) / (1000 * 60 * 60));
    const durationMins = Math.floor((Date.now() - openedAt.getTime()) / (1000 * 60)) % 60;

    return (
      <div className="pos-modal-overlay">
        <div className="pos-modal max-w-md">
          <div className="pos-modal-header">
            <h2 className="pos-modal-title">Shift Aktif</h2>
            <button onClick={onClose} className="p-2 rounded-lg hover:bg-neutral-100 text-neutral-500 hover:text-neutral-900">
              <X className="w-5 h-5" />
            </button>
          </div>

          <div className="pos-modal-body space-y-4">
            <div className="pos-card p-4 space-y-3">
              <div className="flex items-center justify-between">
                <span className="text-sm text-neutral-500">Shift #</span>
                <span className="font-bold text-neutral-900">{activeShift.shift_number || '-'}</span>
              </div>
              <div className="flex items-center justify-between">
                <span className="text-sm text-neutral-500">Status</span>
                <span className="pos-badge-success">Aktif</span>
              </div>
              <div className="flex items-center justify-between">
                <span className="text-sm text-neutral-500">Durasi</span>
                <span className="text-neutral-900">{duration} jam {durationMins} menit</span>
              </div>
              <div className="flex items-center justify-between">
                <span className="text-sm text-neutral-500">Kas Awal</span>
                <span className="text-neutral-900 font-medium">Rp {(activeShift.cash_start || 0).toLocaleString()}</span>
              </div>
            </div>

            <div>
              <label className="block text-sm text-neutral-700 mb-2">Kas Akhir (untuk penutupan)</label>
              <div className="relative">
                <span className="absolute left-4 top-1/2 -translate-y-1/2 text-neutral-400">Rp</span>
                <input
                  type="number"
                  value={cashEnd}
                  onChange={(e) => setCashEnd(Number(e.target.value))}
                  className="pos-input pl-12 py-3 text-lg font-bold"
                  placeholder="0"
                />
              </div>
            </div>

            <button
              onClick={handleCloseShift}
              disabled={loading}
              className="w-full py-3 pos-btn-danger disabled:opacity-50 flex items-center justify-center gap-2"
            >
              {loading ? 'Memproses...' : 'Tutup Shift'}
            </button>
          </div>
        </div>
      </div>
    );
  }

  // Open Shift Form
  return (
    <div className="pos-modal-overlay">
      <div className="pos-modal max-w-md">
        <div className="pos-modal-header">
          <h2 className="pos-modal-title">Buka Shift Baru</h2>
          <button onClick={onClose} className="p-2 rounded-lg hover:bg-neutral-100 text-neutral-500 hover:text-neutral-900">
            <X className="w-5 h-5" />
          </button>
        </div>

        <div className="pos-modal-body space-y-5">
          <div className="text-center">
            <Calendar className="w-12 h-12 text-brand mx-auto mb-2" />
            <p className="text-neutral-900 font-medium">{user?.name || 'Kasir'}</p>
            <p className="text-neutral-500 text-sm mt-1">Masukkan jumlah kas awal untuk memulai shift</p>
          </div>

          <div>
            <label className="block text-sm font-medium text-neutral-700 mb-2">Kas Awal (Cash Drawer)</label>
            <div className="relative">
              <span className="absolute left-4 top-1/2 -translate-y-1/2 text-neutral-400 font-medium">Rp</span>
              <input
                type="number"
                value={cashStart}
                onChange={(e) => setCashStart(Number(e.target.value))}
                className="pos-input text-2xl font-bold pl-12 py-4 text-center"
              />
            </div>
          </div>

          <div className="flex gap-2">
            {[200000, 300000, 500000, 1000000, 2000000].map((amt) => (
              <button
                key={amt}
                onClick={() => setCashStart(amt)}
                className={`flex-1 py-2 text-xs rounded-lg border transition-colors ${
                  cashStart === amt
                    ? 'border-brand bg-primary-50 text-brand'
                    : 'border-neutral-200 text-neutral-500 hover:border-neutral-300'
                }`}
              >
                Rp {amt.toLocaleString()}
              </button>
            ))}
          </div>

          <button
            onClick={handleOpenShift}
            disabled={loading || cashStart <= 0}
            className="w-full py-3.5 pos-btn-primary text-lg"
          >
            {loading ? 'Memproses...' : 'Buka Shift'}
          </button>

          <p className="text-center text-xs text-neutral-500 flex items-center justify-center gap-1">
            <AlertTriangle className="w-3 h-3" />
            Pastikan jumlah kas awal sudah sesuai dengan uang di laci
          </p>
        </div>
      </div>
    </div>
  );
}
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
    } catch {
      // Mock for demo
      const mockShift = {
        id: crypto.randomUUID(),
        outlet_id: user.outletId,
        user_id: user.id,
        shift_number: Math.floor(Math.random() * 100) + 1,
        cash_start: cashStart,
        status: 'open',
        opened_at: new Date().toISOString(),
      };
      setActiveShift(mockShift);
      toast.success(`Shift #${mockShift.shift_number} berhasil dibuka!`);
      onClose();
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
    } catch {
      setActiveShift(null);
      toast.success('Shift berhasil ditutup!');
      onClose();
    }
    setLoading(false);
  };

  if (activeShift) {
    const openedAt = new Date(activeShift.opened_at || activeShift.opened_at);
    const duration = Math.floor((Date.now() - openedAt.getTime()) / (1000 * 60 * 60));
    const durationMins = Math.floor((Date.now() - openedAt.getTime()) / (1000 * 60)) % 60;

    return (
      <div className="fixed inset-0 bg-black/70 z-50 flex items-center justify-center p-4">
        <div className="bg-gray-900 border border-gray-800 rounded-2xl w-full max-w-md shadow-2xl">
          <div className="flex items-center justify-between px-6 py-4 border-b border-gray-800">
            <h2 className="font-bold text-white">Shift Aktif</h2>
            <button onClick={onClose} className="p-2 rounded-lg hover:bg-gray-800 text-gray-400 hover:text-white">
              <X className="w-5 h-5" />
            </button>
          </div>

          <div className="p-6 space-y-4">
            <div className="pos-card p-4 space-y-3">
              <div className="flex items-center justify-between">
                <span className="text-sm text-gray-400">Shift #</span>
                <span className="font-bold text-white">{activeShift.shift_number}</span>
              </div>
              <div className="flex items-center justify-between">
                <span className="text-sm text-gray-400">Status</span>
                <span className="pos-badge-success">Aktif</span>
              </div>
              <div className="flex items-center justify-between">
                <span className="text-sm text-gray-400">Durasi</span>
                <span className="text-white">{duration} jam {durationMins} menit</span>
              </div>
              <div className="flex items-center justify-between">
                <span className="text-sm text-gray-400">Kas Awal</span>
                <span className="text-white font-medium">Rp {activeShift.cash_start.toLocaleString()}</span>
              </div>
            </div>

            <div>
              <label className="block text-sm text-gray-300 mb-2">Kas Akhir (untuk penutupan)</label>
              <div className="relative">
                <span className="absolute left-4 top-1/2 -translate-y-1/2 text-gray-400">Rp</span>
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
              className="w-full py-3 bg-red-600 hover:bg-red-700 disabled:bg-gray-700 text-white font-medium rounded-xl transition-colors flex items-center justify-center gap-2"
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
    <div className="fixed inset-0 bg-black/70 z-50 flex items-center justify-center p-4">
      <div className="bg-gray-900 border border-gray-800 rounded-2xl w-full max-w-md shadow-2xl">
        <div className="flex items-center justify-between px-6 py-4 border-b border-gray-800">
          <h2 className="font-bold text-white">Buka Shift Baru</h2>
          <button onClick={onClose} className="p-2 rounded-lg hover:bg-gray-800 text-gray-400 hover:text-white">
            <X className="w-5 h-5" />
          </button>
        </div>

        <div className="p-6 space-y-5">
          <div className="text-center">
            <Calendar className="w-12 h-12 text-blue-500 mx-auto mb-2" />
            <p className="text-gray-400 text-sm">Masukkan jumlah kas awal untuk memulai shift</p>
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-300 mb-2">Kas Awal (Cash Drawer)</label>
            <div className="relative">
              <span className="absolute left-4 top-1/2 -translate-y-1/2 text-gray-400 font-medium">Rp</span>
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
                    ? 'border-blue-500 bg-blue-600/20 text-blue-400'
                    : 'border-gray-700 text-gray-400 hover:border-gray-600'
                }`}
              >
                Rp {amt.toLocaleString()}
              </button>
            ))}
          </div>

          <button
            onClick={handleOpenShift}
            disabled={loading || cashStart <= 0}
            className="w-full py-3.5 bg-gradient-to-r from-blue-600 to-blue-700 hover:from-blue-700 hover:to-blue-800 disabled:from-gray-700 disabled:to-gray-700 disabled:text-gray-500 text-white font-bold rounded-xl transition-all flex items-center justify-center gap-2 text-lg"
          >
            {loading ? 'Memproses...' : 'Buka Shift'}
          </button>

          <p className="text-center text-xs text-gray-500">
            <AlertTriangle className="w-3 h-3 inline mr-1" />
            Pastikan jumlah kas awal sudah sesuai dengan uang di laci
          </p>
        </div>
      </div>
    </div>
  );
}

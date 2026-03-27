'use client';

import { useEffect, useState } from 'react';
import { apiClient } from '@/lib/api-client';
import { formatCurrency } from '@/lib/formatters';

interface PriceAlert {
  id: string;
  ticker: string;
  condition: string;
  target_price: string;
  triggered: boolean;
  triggered_at: string | null;
  created_at: string;
}

export default function AlertsPage() {
  const [alerts, setAlerts] = useState<PriceAlert[]>([]);
  const [loading, setLoading] = useState(true);

  // Create form
  const [showForm, setShowForm] = useState(false);
  const [ticker, setTicker] = useState('');
  const [condition, setCondition] = useState('above');
  const [targetPrice, setTargetPrice] = useState('');
  const [creating, setCreating] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    loadAlerts();
  }, []);

  async function loadAlerts() {
    setLoading(true);
    try {
      const data = await apiClient.getAlerts();
      setAlerts(data || []);
    } catch (err) {
      console.error('Failed to load alerts:', err);
    } finally {
      setLoading(false);
    }
  }

  async function handleCreate() {
    if (!ticker || !targetPrice) return;
    setCreating(true);
    setError(null);

    try {
      await apiClient.createAlert(ticker.toUpperCase(), condition, targetPrice);
      setShowForm(false);
      setTicker('');
      setTargetPrice('');
      await loadAlerts();
    } catch (err: any) {
      setError(err.message || 'Failed to create alert');
    } finally {
      setCreating(false);
    }
  }

  async function handleDelete(id: string) {
    try {
      await apiClient.deleteAlert(id);
      setAlerts((prev) => prev.filter((a) => a.id !== id));
    } catch (err) {
      console.error('Failed to delete alert:', err);
    }
  }

  const activeAlerts = alerts.filter((a) => !a.triggered);
  const triggeredAlerts = alerts.filter((a) => a.triggered);

  return (
    <div className="p-6 max-w-4xl mx-auto space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold">Price Alerts</h1>
        <button
          onClick={() => setShowForm(!showForm)}
          className="px-4 py-2 rounded-lg bg-[#50E3C2] text-[#0D1117] text-sm font-semibold hover:bg-[#3BC4A7] transition-colors"
        >
          {showForm ? 'Cancel' : 'New Alert'}
        </button>
      </div>

      {/* Create Alert Form */}
      {showForm && (
        <div className="rounded-xl bg-[#161B22] border border-[#30363D] p-6 space-y-4">
          <h2 className="font-semibold">Create Alert</h2>

          <div className="grid grid-cols-3 gap-4">
            <div>
              <label className="block text-xs text-[#6E7681] mb-1.5">Ticker</label>
              <input
                type="text"
                value={ticker}
                onChange={(e) => setTicker(e.target.value.toUpperCase())}
                placeholder="PIPE"
                className="w-full rounded-lg bg-[#0D1117] border border-[#30363D] px-4 py-2.5 text-white placeholder-[#6E7681] focus:outline-none focus:border-[#50E3C2] text-sm"
              />
            </div>
            <div>
              <label className="block text-xs text-[#6E7681] mb-1.5">Condition</label>
              <select
                value={condition}
                onChange={(e) => setCondition(e.target.value)}
                className="w-full rounded-lg bg-[#0D1117] border border-[#30363D] px-4 py-2.5 text-white text-sm focus:outline-none focus:border-[#50E3C2]"
              >
                <option value="above">Price goes above</option>
                <option value="below">Price goes below</option>
              </select>
            </div>
            <div>
              <label className="block text-xs text-[#6E7681] mb-1.5">Target Price</label>
              <input
                type="text"
                value={targetPrice}
                onChange={(e) => setTargetPrice(e.target.value)}
                placeholder="100.00"
                className="w-full rounded-lg bg-[#0D1117] border border-[#30363D] px-4 py-2.5 text-white placeholder-[#6E7681] focus:outline-none focus:border-[#50E3C2] text-sm"
              />
            </div>
          </div>

          {error && (
            <div className="rounded-lg bg-red-500/10 text-red-400 px-4 py-2.5 text-sm">{error}</div>
          )}

          <button
            onClick={handleCreate}
            disabled={creating || !ticker || !targetPrice}
            className="px-6 py-2.5 rounded-lg bg-[#50E3C2] text-[#0D1117] text-sm font-semibold hover:bg-[#3BC4A7] transition-colors disabled:opacity-40 disabled:cursor-not-allowed"
          >
            {creating ? 'Creating...' : 'Create Alert'}
          </button>
        </div>
      )}

      {loading ? (
        <div className="p-8 text-center text-[#8B949E]">Loading alerts...</div>
      ) : alerts.length === 0 ? (
        <div className="rounded-xl bg-[#161B22] border border-[#30363D] p-8 text-center text-[#8B949E]">
          No price alerts yet. Create one to get notified when a stock hits your target price.
        </div>
      ) : (
        <>
          {/* Active Alerts */}
          {activeAlerts.length > 0 && (
            <div className="space-y-3">
              <h2 className="text-sm font-semibold text-[#8B949E] uppercase tracking-wider">Active</h2>
              <div className="space-y-2">
                {activeAlerts.map((alert) => (
                  <div key={alert.id} className="rounded-xl bg-[#161B22] border border-[#30363D] px-5 py-4 flex items-center justify-between">
                    <div className="flex items-center gap-4">
                      <span className="rounded bg-[#50E3C2]/10 px-2 py-0.5 text-xs font-mono font-bold text-[#50E3C2]">
                        {alert.ticker}
                      </span>
                      <span className="text-sm text-[#8B949E]">
                        {alert.condition === 'above' ? '↑ Above' : '↓ Below'}{' '}
                        <span className="text-white font-medium">{formatCurrency(alert.target_price)}</span>
                      </span>
                    </div>
                    <button
                      onClick={() => handleDelete(alert.id)}
                      className="text-xs text-red-400 hover:text-red-300 transition-colors"
                    >
                      Delete
                    </button>
                  </div>
                ))}
              </div>
            </div>
          )}

          {/* Triggered Alerts */}
          {triggeredAlerts.length > 0 && (
            <div className="space-y-3">
              <h2 className="text-sm font-semibold text-[#8B949E] uppercase tracking-wider">Triggered</h2>
              <div className="space-y-2">
                {triggeredAlerts.map((alert) => (
                  <div key={alert.id} className="rounded-xl bg-[#161B22] border border-emerald-500/20 px-5 py-4 flex items-center justify-between opacity-70">
                    <div className="flex items-center gap-4">
                      <span className="rounded bg-emerald-500/10 px-2 py-0.5 text-xs font-mono font-bold text-emerald-400">
                        {alert.ticker}
                      </span>
                      <span className="text-sm text-[#8B949E]">
                        {alert.condition === 'above' ? '↑ Above' : '↓ Below'}{' '}
                        <span className="text-white font-medium">{formatCurrency(alert.target_price)}</span>
                      </span>
                      <span className="text-xs text-emerald-400">Triggered</span>
                    </div>
                    <button
                      onClick={() => handleDelete(alert.id)}
                      className="text-xs text-[#6E7681] hover:text-red-400 transition-colors"
                    >
                      Remove
                    </button>
                  </div>
                ))}
              </div>
            </div>
          )}
        </>
      )}
    </div>
  );
}

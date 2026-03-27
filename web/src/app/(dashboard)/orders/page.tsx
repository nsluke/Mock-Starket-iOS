'use client';

import { useEffect, useState } from 'react';
import { apiClient } from '@/lib/api-client';
import { formatCurrency } from '@/lib/formatters';
import type { Order } from '@/types/portfolio';

export default function OrdersPage() {
  const [orders, setOrders] = useState<Order[]>([]);
  const [loading, setLoading] = useState(true);

  // Create order form
  const [showForm, setShowForm] = useState(false);
  const [ticker, setTicker] = useState('');
  const [side, setSide] = useState('buy');
  const [orderType, setOrderType] = useState('limit');
  const [shares, setShares] = useState('');
  const [limitPrice, setLimitPrice] = useState('');
  const [stopPrice, setStopPrice] = useState('');
  const [creating, setCreating] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    loadOrders();
  }, []);

  async function loadOrders() {
    setLoading(true);
    try {
      const data = await apiClient.getOrders();
      setOrders(data || []);
    } catch (err) {
      console.error('Failed to load orders:', err);
    } finally {
      setLoading(false);
    }
  }

  async function handleCreate() {
    if (!ticker || !shares || parseInt(shares) <= 0) return;
    setCreating(true);
    setError(null);

    try {
      await apiClient.createOrder(
        ticker.toUpperCase(),
        side,
        orderType,
        parseInt(shares),
        orderType !== 'stop' ? limitPrice || undefined : undefined,
        orderType !== 'limit' ? stopPrice || undefined : undefined,
      );
      setShowForm(false);
      setTicker('');
      setShares('');
      setLimitPrice('');
      setStopPrice('');
      await loadOrders();
    } catch (err: any) {
      setError(err.message || 'Failed to create order');
    } finally {
      setCreating(false);
    }
  }

  async function handleCancel(id: string) {
    try {
      await apiClient.cancelOrder(id);
      setOrders((prev) => prev.filter((o) => o.id !== id));
    } catch (err) {
      console.error('Failed to cancel order:', err);
    }
  }

  const orderTypeLabel: Record<string, string> = {
    limit: 'Limit',
    stop: 'Stop',
    stop_limit: 'Stop Limit',
  };

  return (
    <div className="p-6 max-w-4xl mx-auto space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold">Orders</h1>
        <button
          onClick={() => setShowForm(!showForm)}
          className="px-4 py-2 rounded-lg bg-[#50E3C2] text-[#0D1117] text-sm font-semibold hover:bg-[#3BC4A7] transition-colors"
        >
          {showForm ? 'Cancel' : 'New Order'}
        </button>
      </div>

      {/* Create Order Form */}
      {showForm && (
        <div className="rounded-xl bg-[#161B22] border border-[#30363D] p-6 space-y-4">
          <h2 className="font-semibold">Create Order</h2>

          <div className="grid grid-cols-2 gap-4">
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
              <label className="block text-xs text-[#6E7681] mb-1.5">Shares</label>
              <input
                type="number"
                min="1"
                value={shares}
                onChange={(e) => setShares(e.target.value)}
                placeholder="10"
                className="w-full rounded-lg bg-[#0D1117] border border-[#30363D] px-4 py-2.5 text-white placeholder-[#6E7681] focus:outline-none focus:border-[#50E3C2] text-sm"
              />
            </div>
          </div>

          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="block text-xs text-[#6E7681] mb-1.5">Side</label>
              <div className="flex rounded-lg overflow-hidden border border-[#30363D]">
                <button
                  onClick={() => setSide('buy')}
                  className={`flex-1 py-2 text-sm font-semibold ${side === 'buy' ? 'bg-emerald-500/20 text-emerald-400' : 'text-[#8B949E]'}`}
                >
                  Buy
                </button>
                <button
                  onClick={() => setSide('sell')}
                  className={`flex-1 py-2 text-sm font-semibold ${side === 'sell' ? 'bg-red-500/20 text-red-400' : 'text-[#8B949E]'}`}
                >
                  Sell
                </button>
              </div>
            </div>
            <div>
              <label className="block text-xs text-[#6E7681] mb-1.5">Order Type</label>
              <select
                value={orderType}
                onChange={(e) => setOrderType(e.target.value)}
                className="w-full rounded-lg bg-[#0D1117] border border-[#30363D] px-4 py-2.5 text-white text-sm focus:outline-none focus:border-[#50E3C2]"
              >
                <option value="limit">Limit</option>
                <option value="stop">Stop</option>
                <option value="stop_limit">Stop Limit</option>
              </select>
            </div>
          </div>

          <div className="grid grid-cols-2 gap-4">
            {orderType !== 'stop' && (
              <div>
                <label className="block text-xs text-[#6E7681] mb-1.5">Limit Price</label>
                <input
                  type="text"
                  value={limitPrice}
                  onChange={(e) => setLimitPrice(e.target.value)}
                  placeholder="100.00"
                  className="w-full rounded-lg bg-[#0D1117] border border-[#30363D] px-4 py-2.5 text-white placeholder-[#6E7681] focus:outline-none focus:border-[#50E3C2] text-sm"
                />
              </div>
            )}
            {orderType !== 'limit' && (
              <div>
                <label className="block text-xs text-[#6E7681] mb-1.5">Stop Price</label>
                <input
                  type="text"
                  value={stopPrice}
                  onChange={(e) => setStopPrice(e.target.value)}
                  placeholder="95.00"
                  className="w-full rounded-lg bg-[#0D1117] border border-[#30363D] px-4 py-2.5 text-white placeholder-[#6E7681] focus:outline-none focus:border-[#50E3C2] text-sm"
                />
              </div>
            )}
          </div>

          {error && (
            <div className="rounded-lg bg-red-500/10 text-red-400 px-4 py-2.5 text-sm">{error}</div>
          )}

          <button
            onClick={handleCreate}
            disabled={creating || !ticker || !shares}
            className="w-full py-3 rounded-lg bg-[#50E3C2] text-[#0D1117] text-sm font-semibold hover:bg-[#3BC4A7] transition-colors disabled:opacity-40 disabled:cursor-not-allowed"
          >
            {creating ? 'Creating...' : 'Place Order'}
          </button>
        </div>
      )}

      {/* Orders List */}
      <div className="rounded-xl bg-[#161B22] border border-[#30363D] overflow-hidden">
        {loading ? (
          <div className="p-8 text-center text-[#8B949E]">Loading orders...</div>
        ) : orders.length === 0 ? (
          <div className="p-8 text-center text-[#8B949E]">
            No open orders. Create a limit or stop order to get started.
          </div>
        ) : (
          <table className="w-full">
            <thead>
              <tr className="border-b border-[#30363D] text-xs text-[#8B949E] uppercase">
                <th className="text-left px-4 py-3">Stock</th>
                <th className="text-left px-4 py-3">Type</th>
                <th className="text-left px-4 py-3">Side</th>
                <th className="text-right px-4 py-3">Shares</th>
                <th className="text-right px-4 py-3 hidden sm:table-cell">Limit</th>
                <th className="text-right px-4 py-3 hidden sm:table-cell">Stop</th>
                <th className="text-right px-4 py-3">Action</th>
              </tr>
            </thead>
            <tbody>
              {orders.map((order) => (
                <tr key={order.id} className="border-b border-[#21262D] hover:bg-[#21262D] transition-colors">
                  <td className="px-4 py-3">
                    <span className="rounded bg-[#50E3C2]/10 px-2 py-0.5 text-xs font-mono font-bold text-[#50E3C2]">
                      {order.ticker}
                    </span>
                  </td>
                  <td className="px-4 py-3 text-sm text-[#8B949E]">
                    {orderTypeLabel[order.order_type] || order.order_type}
                  </td>
                  <td className="px-4 py-3">
                    <span className={`text-xs font-semibold uppercase ${order.side === 'buy' ? 'text-emerald-400' : 'text-red-400'}`}>
                      {order.side}
                    </span>
                  </td>
                  <td className="px-4 py-3 text-right text-sm">{order.shares}</td>
                  <td className="px-4 py-3 text-right text-sm text-[#8B949E] hidden sm:table-cell">
                    {order.limit_price ? formatCurrency(order.limit_price) : '-'}
                  </td>
                  <td className="px-4 py-3 text-right text-sm text-[#8B949E] hidden sm:table-cell">
                    {order.stop_price ? formatCurrency(order.stop_price) : '-'}
                  </td>
                  <td className="px-4 py-3 text-right">
                    <button
                      onClick={() => handleCancel(order.id)}
                      className="text-xs text-red-400 hover:text-red-300 transition-colors"
                    >
                      Cancel
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </div>
    </div>
  );
}

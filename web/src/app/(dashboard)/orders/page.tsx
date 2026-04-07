'use client';

import { useState } from 'react';
import { PageTransition } from '@/components/ui/PageTransition';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { useOrders, useCreateOrder, useCancelOrder } from '@/hooks/use-orders';
import { formatCurrency } from '@/lib/formatters';
import { orderSchema, type OrderFormValues } from '@/lib/schemas';
import { FormInput } from '@/components/ui/FormInput';

export default function OrdersPage() {
  const { data: orders = [], isLoading } = useOrders();
  const createOrder = useCreateOrder();
  const cancelOrder = useCancelOrder();

  const [showForm, setShowForm] = useState(false);
  const { register, handleSubmit, watch, reset, formState: { errors } } = useForm<OrderFormValues>({
    resolver: zodResolver(orderSchema),
    defaultValues: { side: 'buy', order_type: 'limit' },
  });
  const orderType = watch('order_type');
  const side = watch('side');

  function onSubmitOrder(data: OrderFormValues) {
    createOrder.mutate(
      {
        ticker: data.ticker,
        side: data.side,
        orderType: data.order_type,
        shares: data.shares,
        limitPrice: data.limit_price || undefined,
        stopPrice: data.stop_price || undefined,
      },
      {
        onSuccess: () => {
          setShowForm(false);
          reset();
        },
      }
    );
  }

  const orderTypeLabel: Record<string, string> = {
    limit: 'Limit',
    stop: 'Stop',
    stop_limit: 'Stop Limit',
  };

  return (
    <PageTransition>
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
        <form onSubmit={handleSubmit(onSubmitOrder)} className="rounded-xl bg-[#161B22] border border-[#30363D] p-6 space-y-4">
          <h2 className="font-semibold">Create Order</h2>

          <div className="grid grid-cols-2 gap-4">
            <FormInput label="Ticker" placeholder="PIPE" error={errors.ticker?.message} {...register('ticker')} />
            <FormInput label="Shares" type="number" min="1" placeholder="10" error={errors.shares?.message} {...register('shares')} />
          </div>

          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="block text-xs text-[#6E7681] mb-1.5">Side</label>
              <div className="flex rounded-lg overflow-hidden border border-[#30363D]">
                <label className={`flex-1 py-2 text-sm font-semibold text-center cursor-pointer ${side === 'buy' ? 'bg-emerald-500/20 text-emerald-400' : 'text-[#8B949E]'}`}>
                  <input type="radio" value="buy" className="sr-only" {...register('side')} />
                  Buy
                </label>
                <label className={`flex-1 py-2 text-sm font-semibold text-center cursor-pointer ${side === 'sell' ? 'bg-red-500/20 text-red-400' : 'text-[#8B949E]'}`}>
                  <input type="radio" value="sell" className="sr-only" {...register('side')} />
                  Sell
                </label>
              </div>
            </div>
            <div>
              <label className="block text-xs text-[#6E7681] mb-1.5">Order Type</label>
              <select
                {...register('order_type')}
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
              <FormInput label="Limit Price" placeholder="100.00" error={errors.limit_price?.message} {...register('limit_price')} />
            )}
            {orderType !== 'limit' && (
              <FormInput label="Stop Price" placeholder="95.00" error={errors.stop_price?.message} {...register('stop_price')} />
            )}
          </div>

          {createOrder.isError && (
            <div className="rounded-lg bg-red-500/10 text-red-400 px-4 py-2.5 text-sm">
              {createOrder.error?.message || 'Failed to create order'}
            </div>
          )}

          <button
            type="submit"
            disabled={createOrder.isPending}
            className="w-full py-3 rounded-lg bg-[#50E3C2] text-[#0D1117] text-sm font-semibold hover:bg-[#3BC4A7] transition-colors disabled:opacity-40 disabled:cursor-not-allowed"
          >
            {createOrder.isPending ? 'Creating...' : 'Place Order'}
          </button>
        </form>
      )}

      {/* Orders List */}
      <div className="rounded-xl bg-[#161B22] border border-[#30363D] overflow-hidden">
        {isLoading ? (
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
                      onClick={() => cancelOrder.mutate(order.id)}
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
    </PageTransition>
  );
}

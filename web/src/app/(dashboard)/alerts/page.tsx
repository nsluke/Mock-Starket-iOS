'use client';

import { useState } from 'react';
import { PageTransition } from '@/components/ui/PageTransition';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { useAlerts, useCreateAlert, useDeleteAlert } from '@/hooks/use-alerts';
import { formatCurrency } from '@/lib/formatters';
import { alertSchema, type AlertFormValues } from '@/lib/schemas';
import { FormInput } from '@/components/ui/FormInput';

export default function AlertsPage() {
  const { data: alerts = [], isLoading } = useAlerts();
  const createAlert = useCreateAlert();
  const deleteAlert = useDeleteAlert();

  const [showForm, setShowForm] = useState(false);
  const { register, handleSubmit, reset, formState: { errors } } = useForm<AlertFormValues>({
    resolver: zodResolver(alertSchema),
    defaultValues: { condition: 'above' },
  });

  function onSubmitAlert(data: AlertFormValues) {
    createAlert.mutate(
      { ticker: data.ticker, condition: data.condition, targetPrice: data.target_price },
      {
        onSuccess: () => {
          setShowForm(false);
          reset();
        },
      }
    );
  }

  const activeAlerts = alerts.filter((a) => !a.triggered);
  const triggeredAlerts = alerts.filter((a) => a.triggered);

  return (
    <PageTransition>
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
        <form onSubmit={handleSubmit(onSubmitAlert)} className="rounded-xl bg-[#161B22] border border-[#30363D] p-6 space-y-4">
          <h2 className="font-semibold">Create Alert</h2>

          <div className="grid grid-cols-3 gap-4">
            <FormInput label="Ticker" placeholder="PIPE" error={errors.ticker?.message} {...register('ticker')} />
            <div>
              <label className="block text-xs text-[#6E7681] mb-1.5">Condition</label>
              <select
                {...register('condition')}
                className="w-full rounded-lg bg-[#0D1117] border border-[#30363D] px-4 py-2.5 text-white text-sm focus:outline-none focus:border-[#50E3C2]"
              >
                <option value="above">Price goes above</option>
                <option value="below">Price goes below</option>
              </select>
            </div>
            <FormInput label="Target Price" placeholder="100.00" error={errors.target_price?.message} {...register('target_price')} />
          </div>

          {createAlert.isError && (
            <div className="rounded-lg bg-red-500/10 text-red-400 px-4 py-2.5 text-sm">
              {createAlert.error?.message || 'Failed to create alert'}
            </div>
          )}

          <button
            type="submit"
            disabled={createAlert.isPending}
            className="px-6 py-2.5 rounded-lg bg-[#50E3C2] text-[#0D1117] text-sm font-semibold hover:bg-[#3BC4A7] transition-colors disabled:opacity-40 disabled:cursor-not-allowed"
          >
            {createAlert.isPending ? 'Creating...' : 'Create Alert'}
          </button>
        </form>
      )}

      {isLoading ? (
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
                        {alert.condition === 'above' ? '\u2191 Above' : '\u2193 Below'}{' '}
                        <span className="text-white font-medium">{formatCurrency(alert.target_price)}</span>
                      </span>
                    </div>
                    <button
                      onClick={() => deleteAlert.mutate(alert.id)}
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
                        {alert.condition === 'above' ? '\u2191 Above' : '\u2193 Below'}{' '}
                        <span className="text-white font-medium">{formatCurrency(alert.target_price)}</span>
                      </span>
                      <span className="text-xs text-emerald-400">Triggered</span>
                    </div>
                    <button
                      onClick={() => deleteAlert.mutate(alert.id)}
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
    </PageTransition>
  );
}

'use client';

import { useEffect, useState, lazy, Suspense } from 'react';
import { useParams, useRouter } from 'next/navigation';
import Link from 'next/link';
import { apiClient } from '@/lib/api-client';
import { formatCurrency, formatPercent, formatVolume, priceChangeColor, priceChangeBg } from '@/lib/formatters';
import { useMarketStore } from '@/stores/market-store';
import type { Stock, PricePoint } from '@/types/stock';

const PriceChart = lazy(() => import('@/components/charts/PriceChart'));

export default function StockDetailPage() {
  const { ticker } = useParams<{ ticker: string }>();
  const router = useRouter();
  const { stocks } = useMarketStore();

  const [stock, setStock] = useState<Stock | null>(null);
  const [history, setHistory] = useState<PricePoint[]>([]);
  const [interval, setInterval] = useState('1m');
  const [loading, setLoading] = useState(true);

  // Trade form state
  const [side, setSide] = useState<'buy' | 'sell'>('buy');
  const [shares, setShares] = useState('');
  const [trading, setTrading] = useState(false);
  const [tradeResult, setTradeResult] = useState<{ success: boolean; message: string } | null>(null);

  useEffect(() => {
    async function load() {
      setLoading(true);
      try {
        const [stockData, historyData] = await Promise.all([
          apiClient.getStock(ticker),
          apiClient.getStockHistory(ticker, interval),
        ]);
        setStock(stockData);
        setHistory(historyData || []);
      } catch {
        console.error('Failed to load stock');
      } finally {
        setLoading(false);
      }
    }
    load();
  }, [ticker, interval]);

  // Keep live price from market store
  const liveStock = stocks.find((s) => s.ticker === ticker);
  const currentPrice = liveStock ? parseFloat(liveStock.current_price) : stock ? parseFloat(stock.current_price) : 0;
  const dayOpen = stock ? parseFloat(stock.day_open) : 0;
  const change = currentPrice - dayOpen;
  const changePct = dayOpen !== 0 ? (change / dayOpen) * 100 : 0;
  const totalCost = parseFloat(shares || '0') * currentPrice;

  async function handleTrade() {
    const qty = parseInt(shares);
    if (!qty || qty <= 0) return;

    setTrading(true);
    setTradeResult(null);
    try {
      await apiClient.executeTrade(ticker, side, qty);
      setTradeResult({
        success: true,
        message: `${side === 'buy' ? 'Bought' : 'Sold'} ${qty} shares of ${ticker} at ${formatCurrency(currentPrice)}`,
      });
      setShares('');
    } catch (err: any) {
      setTradeResult({ success: false, message: err.message || 'Trade failed' });
    } finally {
      setTrading(false);
    }
  }

  if (loading) {
    return <div className="p-6 text-center text-[#8B949E]">Loading...</div>;
  }

  if (!stock) {
    return <div className="p-6 text-center text-[#8B949E]">Stock not found</div>;
  }

  return (
    <div className="p-6 max-w-6xl mx-auto space-y-6">
      {/* Back link */}
      <Link href="/market" className="text-sm text-[#8B949E] hover:text-white transition-colors">
        &larr; Back to Market
      </Link>

      {/* Header */}
      <div className="flex items-start justify-between">
        <div>
          <div className="flex items-center gap-3 mb-1">
            <span className="rounded bg-[#50E3C2]/10 px-2.5 py-1 text-sm font-mono font-bold text-[#50E3C2]">
              {stock.ticker}
            </span>
            <span className="text-sm text-[#8B949E]">{stock.sector}</span>
          </div>
          <h1 className="text-2xl font-bold">{stock.name}</h1>
          {stock.description && (
            <p className="text-sm text-[#8B949E] mt-1 max-w-xl">{stock.description}</p>
          )}
        </div>
      </div>

      {/* Price & Stats */}
      <div className="rounded-xl bg-[#161B22] border border-[#30363D] p-6">
        <div className="flex items-baseline gap-4 mb-4">
          <span className="text-4xl font-bold">{formatCurrency(currentPrice)}</span>
          <span className={`text-lg font-semibold ${priceChangeColor(changePct)}`}>
            {change >= 0 ? '+' : ''}{formatCurrency(Math.abs(change))} ({formatPercent(changePct)})
          </span>
        </div>

        <div className="grid grid-cols-2 sm:grid-cols-4 gap-4 text-sm">
          <div>
            <p className="text-[#6E7681]">Open</p>
            <p className="font-medium">{formatCurrency(stock.day_open)}</p>
          </div>
          <div>
            <p className="text-[#6E7681]">High</p>
            <p className="font-medium text-emerald-400">{formatCurrency(stock.day_high)}</p>
          </div>
          <div>
            <p className="text-[#6E7681]">Low</p>
            <p className="font-medium text-red-400">{formatCurrency(stock.day_low)}</p>
          </div>
          <div>
            <p className="text-[#6E7681]">Volume</p>
            <p className="font-medium">{formatVolume(stock.volume)}</p>
          </div>
        </div>
      </div>

      {/* Price Chart */}
      <div className="rounded-xl bg-[#161B22] border border-[#30363D] p-6">
        <div className="flex items-center justify-between mb-4">
          <h2 className="font-semibold">Price Chart</h2>
          <div className="flex gap-1">
            {['1m', '5m', '1h'].map((iv) => (
              <button
                key={iv}
                onClick={() => setInterval(iv)}
                className={`px-3 py-1 rounded text-xs font-medium transition-colors ${
                  interval === iv
                    ? 'bg-[#50E3C2]/10 text-[#50E3C2]'
                    : 'text-[#8B949E] hover:text-white'
                }`}
              >
                {iv}
              </button>
            ))}
          </div>
        </div>

        {history.length > 0 ? (
          <Suspense fallback={<div className="h-[300px] flex items-center justify-center text-[#8B949E]">Loading chart...</div>}>
            <PriceChart
              data={history.map((p) => ({
                time: p.recorded_at,
                open: parseFloat(p.open as any),
                high: parseFloat(p.high as any),
                low: parseFloat(p.low as any),
                close: parseFloat(p.close as any),
              }))}
            />
          </Suspense>
        ) : (
          <div className="h-[300px] flex items-center justify-center text-[#6E7681] text-sm">
            No price history yet. Data will appear as the market simulation runs.
          </div>
        )}
      </div>

      {/* Trade Panel */}
      <div className="rounded-xl bg-[#161B22] border border-[#30363D] p-6">
        <h2 className="font-semibold mb-4">Trade {stock.ticker}</h2>

        {/* Buy / Sell toggle */}
        <div className="flex rounded-lg overflow-hidden border border-[#30363D] mb-4">
          <button
            onClick={() => setSide('buy')}
            className={`flex-1 py-2.5 text-sm font-semibold transition-colors ${
              side === 'buy'
                ? 'bg-emerald-500/20 text-emerald-400 border-r border-[#30363D]'
                : 'text-[#8B949E] hover:text-white border-r border-[#30363D]'
            }`}
          >
            Buy
          </button>
          <button
            onClick={() => setSide('sell')}
            className={`flex-1 py-2.5 text-sm font-semibold transition-colors ${
              side === 'sell'
                ? 'bg-red-500/20 text-red-400'
                : 'text-[#8B949E] hover:text-white'
            }`}
          >
            Sell
          </button>
        </div>

        {/* Shares input */}
        <div className="mb-4">
          <label className="block text-xs text-[#6E7681] mb-1.5">Shares</label>
          <input
            type="number"
            min="1"
            value={shares}
            onChange={(e) => setShares(e.target.value)}
            placeholder="0"
            className="w-full rounded-lg bg-[#0D1117] border border-[#30363D] px-4 py-3 text-white placeholder-[#6E7681] focus:outline-none focus:border-[#50E3C2]"
          />
        </div>

        {/* Estimated total */}
        {parseFloat(shares || '0') > 0 && (
          <div className="flex justify-between text-sm mb-4 px-1">
            <span className="text-[#8B949E]">Estimated Total</span>
            <span className="font-semibold">{formatCurrency(totalCost)}</span>
          </div>
        )}

        {/* Submit */}
        <button
          onClick={handleTrade}
          disabled={trading || !shares || parseInt(shares) <= 0}
          className={`w-full py-3 rounded-lg text-sm font-semibold transition-colors disabled:opacity-40 disabled:cursor-not-allowed ${
            side === 'buy'
              ? 'bg-emerald-500 hover:bg-emerald-600 text-white'
              : 'bg-red-500 hover:bg-red-600 text-white'
          }`}
        >
          {trading ? 'Processing...' : `${side === 'buy' ? 'Buy' : 'Sell'} ${stock.ticker}`}
        </button>

        {/* Trade result */}
        {tradeResult && (
          <div
            className={`mt-3 rounded-lg px-4 py-3 text-sm ${
              tradeResult.success
                ? 'bg-emerald-500/10 text-emerald-400'
                : 'bg-red-500/10 text-red-400'
            }`}
          >
            {tradeResult.message}
          </div>
        )}
      </div>
    </div>
  );
}

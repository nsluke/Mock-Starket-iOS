'use client';

import { useState } from 'react';
import { useParams } from 'next/navigation';
import Link from 'next/link';
import dynamic from 'next/dynamic';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { PageTransition } from '@/components/ui/PageTransition';
import { useStock, useStockHistory, useETFHoldings } from '@/hooks/use-stocks';
import { useExecuteTrade, usePortfolio } from '@/hooks/use-portfolio';
import { useCreateOrder } from '@/hooks/use-orders';
import { useWatchlist, useAddToWatchlist, useRemoveFromWatchlist } from '@/hooks/use-watchlist';
import { formatCurrency, formatPercent, formatVolume, priceChangeColor } from '@/lib/formatters';
import { tradeSchema, type TradeFormValues } from '@/lib/schemas';
import { useMarketStore } from '@/stores/market-store';
import { PieChart, COLORS } from '@/components/charts/PieChart';

import { OptionChainTable } from '@/components/options/OptionChainTable';

const PriceChart = dynamic(() => import('@/components/charts/PriceChart'), { ssr: false });

export default function StockDetailPage() {
  const { ticker } = useParams<{ ticker: string }>();
  const { stocks } = useMarketStore();

  const { data: stock, isLoading } = useStock(ticker);
  const [interval, setInterval] = useState('1m');
  const { data: history = [] } = useStockHistory(ticker, interval);
  const { data: etfHoldings = [] } = useETFHoldings(ticker, stock?.asset_type);
  const { data: watchlist = [] } = useWatchlist();
  const addToWatchlist = useAddToWatchlist();
  const removeFromWatchlist = useRemoveFromWatchlist();
  const isWatched = watchlist.includes(ticker);

  // Portfolio data for sell-all
  const { data: portfolio } = usePortfolio();
  const position = portfolio?.positions?.find((p: any) => p.ticker === ticker);
  const ownedShares = position ? Number(position.shares) : 0;

  // Trade form
  const [side, setSide] = useState<'buy' | 'sell'>('buy');
  const [orderMode, setOrderMode] = useState<'market' | 'limit' | 'stop'>('market');
  const executeTrade = useExecuteTrade();
  const createOrder = useCreateOrder();
  const { register, handleSubmit, watch, reset, formState: { errors } } = useForm<TradeFormValues>({
    resolver: zodResolver(tradeSchema),
  });
  const sharesValue = watch('shares');

  // Keep live price from market store
  const liveStock = stocks.find((s) => s.ticker === ticker);
  const currentPrice = liveStock ? parseFloat(liveStock.current_price) : stock ? parseFloat(stock.current_price) : 0;
  const dayOpen = stock ? parseFloat(stock.day_open) : 0;
  const change = currentPrice - dayOpen;
  const changePct = dayOpen !== 0 ? (change / dayOpen) * 100 : 0;
  const totalCost = (sharesValue || 0) * currentPrice;

  const [limitPrice, setLimitPrice] = useState('');
  const [stopPrice, setStopPrice] = useState('');

  function onTrade(data: TradeFormValues) {
    if (orderMode === 'market') {
      executeTrade.mutate(
        { ticker, side, shares: data.shares },
        { onSuccess: () => reset() }
      );
    } else {
      createOrder.mutate(
        {
          ticker,
          side,
          orderType: orderMode === 'limit' ? 'limit' : 'stop',
          shares: data.shares,
          limitPrice: orderMode === 'limit' ? limitPrice : undefined,
          stopPrice: orderMode === 'stop' ? stopPrice : undefined,
        },
        { onSuccess: () => { reset(); setLimitPrice(''); setStopPrice(''); } }
      );
    }
  }

  function handleQuickBuy(shares: number) {
    executeTrade.mutate({ ticker, side: 'buy', shares });
  }

  function handleSellAll() {
    if (ownedShares > 0) {
      executeTrade.mutate({ ticker, side: 'sell', shares: ownedShares });
    }
  }

  if (isLoading) {
    return <div className="p-6 text-center text-[#8B949E]">Loading...</div>;
  }

  if (!stock) {
    return <div className="p-6 text-center text-[#8B949E]">Stock not found</div>;
  }

  return (
    <PageTransition>
    <div className="p-6 max-w-6xl mx-auto space-y-6">
      {/* Back link */}
      <Link href="/market" className="text-sm text-[#8B949E] hover:text-white transition-colors">
        &larr; Back to Market
      </Link>

      {/* Header */}
      <div className="flex items-start justify-between">
        <div>
          <button
            onClick={() =>
              isWatched
                ? removeFromWatchlist.mutate(ticker)
                : addToWatchlist.mutate(ticker)
            }
            className={`float-right ml-3 mt-1 px-3 py-1.5 rounded-lg text-xs font-medium border transition-colors ${
              isWatched
                ? 'border-[#50E3C2]/30 text-[#50E3C2] hover:bg-[#50E3C2]/10'
                : 'border-[#30363D] text-[#8B949E] hover:text-white hover:border-[#50E3C2]'
            }`}
          >
            {isWatched ? '★ Watching' : '☆ Watch'}
          </button>
          <div className="flex items-center gap-3 mb-1">
            <span className="rounded bg-[#50E3C2]/10 px-2.5 py-1 text-sm font-mono font-bold text-[#50E3C2]">
              {stock.ticker}
            </span>
            <span className="text-sm text-[#8B949E]">{stock.sector}</span>
            <span className={`text-xs font-medium px-2 py-0.5 rounded ${
              stock.asset_type === 'crypto' ? 'bg-orange-400/10 text-orange-400' :
              stock.asset_type === 'commodity' ? 'bg-yellow-400/10 text-yellow-400' :
              stock.asset_type === 'etf' ? 'bg-purple-400/10 text-purple-400' :
              'bg-[#8B949E]/10 text-[#8B949E]'
            }`}>
              {stock.asset_type?.toUpperCase() || 'STOCK'}
            </span>
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
          <PriceChart
            data={history.map((p) => ({
              time: p.recorded_at,
              open: parseFloat(p.open as any),
              high: parseFloat(p.high as any),
              low: parseFloat(p.low as any),
              close: parseFloat(p.close as any),
            }))}
          />
        ) : (
          <div className="h-[300px] flex items-center justify-center text-[#6E7681] text-sm">
            No price history yet. Data will appear as the market simulation runs.
          </div>
        )}
      </div>

      {/* ETF Holdings */}
      {stock.asset_type === 'etf' && etfHoldings.length > 0 && (
        <div className="rounded-xl bg-[#161B22] border border-[#30363D] p-6">
          <h2 className="font-semibold mb-4">Holdings</h2>
          <div className="flex flex-col md:flex-row gap-6">
            <div>
              <PieChart
                data={etfHoldings.map((h, i) => ({
                  label: h.ticker,
                  value: parseFloat(h.weight) * 100,
                  color: COLORS[i % COLORS.length],
                }))}
                size={140}
              />
            </div>
            <div className="flex-1 space-y-3">
              {etfHoldings.map((h) => (
                <div key={h.ticker} className="flex items-center justify-between">
                  <div className="flex items-center gap-3">
                    <Link href={`/stock/${h.ticker}`}>
                      <span className="rounded bg-[#50E3C2]/10 px-2 py-0.5 text-xs font-mono font-bold text-[#50E3C2] hover:bg-[#50E3C2]/20 transition-colors">
                        {h.ticker}
                      </span>
                    </Link>
                    <span className="text-sm text-[#8B949E]">{h.name}</span>
                  </div>
                  <div className="flex items-center gap-4">
                    <span className="text-sm text-[#8B949E]">{h.price ? formatCurrency(h.price) : '-'}</span>
                    <span className="text-sm font-semibold">{(parseFloat(h.weight) * 100).toFixed(0)}%</span>
                  </div>
                </div>
              ))}
            </div>
          </div>
        </div>
      )}

      {/* Trade Panel */}
      <div className="rounded-xl bg-[#161B22] border border-[#30363D] p-6">
        <h2 className="font-semibold mb-4">Trade {stock.ticker}</h2>

        {/* Position info */}
        {ownedShares > 0 && (
          <div className="flex items-center justify-between rounded-lg bg-[#0D1117] border border-[#30363D] px-4 py-3 mb-4 text-sm">
            <span className="text-[#8B949E]">You own <span className="text-white font-semibold">{ownedShares}</span> shares</span>
            <button
              onClick={handleSellAll}
              disabled={executeTrade.isPending}
              className="px-3 py-1 rounded-lg bg-red-500/10 text-red-400 text-xs font-semibold hover:bg-red-500/20 transition-colors disabled:opacity-40"
            >
              Sell All
            </button>
          </div>
        )}

        {/* Quick Buy Buttons */}
        <div className="mb-4">
          <label className="block text-xs text-[#6E7681] mb-2">Quick Buy</label>
          <div className="grid grid-cols-4 gap-2">
            {[1, 10, 100, 1000].map((qty) => (
              <button
                key={qty}
                onClick={() => handleQuickBuy(qty)}
                disabled={executeTrade.isPending}
                className="py-2 rounded-lg bg-emerald-500/10 text-emerald-400 text-sm font-semibold hover:bg-emerald-500/20 transition-colors disabled:opacity-40"
              >
                {qty}
              </button>
            ))}
          </div>
          <p className="text-xs text-[#6E7681] mt-1">
            1 share = {formatCurrency(currentPrice)}
          </p>
        </div>

        <div className="border-t border-[#30363D] pt-4 mt-4">
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

          {/* Order type toggle */}
          <div className="flex gap-1 mb-4">
            {(['market', 'limit', 'stop'] as const).map((mode) => (
              <button
                key={mode}
                onClick={() => setOrderMode(mode)}
                className={`px-3 py-1.5 rounded-lg text-xs font-medium transition-colors ${
                  orderMode === mode
                    ? 'bg-[#21262D] text-white'
                    : 'text-[#8B949E] hover:text-white'
                }`}
              >
                {mode.charAt(0).toUpperCase() + mode.slice(1)}
              </button>
            ))}
          </div>

          {/* Trade form */}
          <form onSubmit={handleSubmit(onTrade)}>
            <div className="mb-4">
              <label className="block text-xs text-[#6E7681] mb-1.5">Shares</label>
              <input
                type="number"
                min="1"
                placeholder="0"
                {...register('shares')}
                className={`w-full rounded-lg bg-[#0D1117] border px-4 py-3 text-white placeholder-[#6E7681] focus:outline-none ${
                  errors.shares ? 'border-red-400' : 'border-[#30363D] focus:border-[#50E3C2]'
                }`}
              />
              {errors.shares && (
                <p className="text-xs text-red-400 mt-1">{errors.shares.message}</p>
              )}
            </div>

            {/* Limit price input */}
            {orderMode === 'limit' && (
              <div className="mb-4">
                <label className="block text-xs text-[#6E7681] mb-1.5">
                  Limit Price {side === 'buy' ? '(buy at or below)' : '(sell at or above)'}
                </label>
                <input
                  type="number"
                  step="0.01"
                  min="0.01"
                  placeholder={currentPrice.toFixed(2)}
                  value={limitPrice}
                  onChange={(e) => setLimitPrice(e.target.value)}
                  className="w-full rounded-lg bg-[#0D1117] border border-[#30363D] px-4 py-3 text-white placeholder-[#6E7681] focus:outline-none focus:border-[#50E3C2]"
                />
              </div>
            )}

            {/* Stop price input */}
            {orderMode === 'stop' && (
              <div className="mb-4">
                <label className="block text-xs text-[#6E7681] mb-1.5">
                  Stop Price {side === 'buy' ? '(trigger above)' : '(trigger below)'}
                </label>
                <input
                  type="number"
                  step="0.01"
                  min="0.01"
                  placeholder={currentPrice.toFixed(2)}
                  value={stopPrice}
                  onChange={(e) => setStopPrice(e.target.value)}
                  className="w-full rounded-lg bg-[#0D1117] border border-[#30363D] px-4 py-3 text-white placeholder-[#6E7681] focus:outline-none focus:border-[#50E3C2]"
                />
              </div>
            )}

            {totalCost > 0 && orderMode === 'market' && (
              <div className="flex justify-between text-sm mb-4 px-1">
                <span className="text-[#8B949E]">Estimated Total</span>
                <span className="font-semibold">{formatCurrency(totalCost)}</span>
              </div>
            )}

            <button
              type="submit"
              disabled={executeTrade.isPending || createOrder.isPending}
              className={`w-full py-3 rounded-lg text-sm font-semibold transition-colors disabled:opacity-40 disabled:cursor-not-allowed ${
                side === 'buy'
                  ? 'bg-emerald-500 hover:bg-emerald-600 text-white'
                  : 'bg-red-500 hover:bg-red-600 text-white'
              }`}
            >
              {(executeTrade.isPending || createOrder.isPending)
                ? 'Processing...'
                : orderMode === 'market'
                  ? `${side === 'buy' ? 'Buy' : 'Sell'} ${stock.ticker}`
                  : `Place ${orderMode} ${side} order`
              }
            </button>
          </form>
        </div>
      </div>

      {/* Options Chain */}
      {stock.asset_type === 'stock' && (
        <div className="rounded-xl bg-[#161B22] border border-[#30363D] p-6">
          <h2 className="font-semibold mb-4">Options Chain</h2>
          <OptionChainTable ticker={ticker} underlyingPrice={currentPrice} />
        </div>
      )}
    </div>
    </PageTransition>
  );
}

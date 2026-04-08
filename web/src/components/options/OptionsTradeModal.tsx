'use client';

import { useState } from 'react';
import { useExecuteOptionsTrade } from '@/hooks/use-options';
import { GreeksDisplay } from './GreeksDisplay';
import { PayoffDiagram } from './PayoffDiagram';
import { formatCurrency } from '@/lib/formatters';
import type { OptionContract } from '@/types/options';

interface OptionsTradeModalProps {
  contract: OptionContract;
  underlyingPrice: number;
  onClose: () => void;
}

type Side = 'buy_to_open' | 'sell_to_open';

const sideInfo: Record<Side, { label: string; description: string; color: string }> = {
  buy_to_open: {
    label: 'Buy to Open',
    description: 'Buy this contract. Pay the ask price as premium. Max loss = premium paid.',
    color: 'text-emerald-400',
  },
  sell_to_open: {
    label: 'Sell to Open (Write)',
    description: 'Write/sell this contract. Collect premium but take on obligation. Requires collateral.',
    color: 'text-orange-400',
  },
};

export function OptionsTradeModal({ contract, underlyingPrice, onClose }: OptionsTradeModalProps) {
  const [side, setSide] = useState<Side>('buy_to_open');
  const [quantity, setQuantity] = useState(1);
  const [step, setStep] = useState<'configure' | 'review'>('configure');
  const executeTrade = useExecuteOptionsTrade();

  const isCall = contract.option_type === 'call';
  const strike = parseFloat(contract.strike_price);
  const bid = parseFloat(contract.bid_price);
  const ask = parseFloat(contract.ask_price);
  const mark = parseFloat(contract.mark_price);
  const price = side === 'buy_to_open' ? ask : bid;
  const totalCost = price * quantity * 100;
  const isLong = side === 'buy_to_open';

  // Break-even
  const breakEven = isCall ? strike + price : strike - price;

  // Moneyness
  const itm = isCall ? underlyingPrice > strike : underlyingPrice < strike;
  const atm = Math.abs(underlyingPrice - strike) / underlyingPrice < 0.01;
  const moneyness = atm ? 'ATM' : itm ? 'ITM' : 'OTM';
  const moneynessColor = atm ? 'text-yellow-400' : itm ? 'text-emerald-400' : 'text-[#6E7681]';

  const handleSubmit = () => {
    executeTrade.mutate(
      { contractId: contract.id, side, quantity },
      {
        onSuccess: () => {
          onClose();
        },
      }
    );
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm" onClick={onClose}>
      <div
        className="w-full max-w-lg max-h-[90vh] overflow-y-auto rounded-xl bg-[#161B22] border border-[#30363D] shadow-2xl"
        onClick={(e) => e.stopPropagation()}
      >
        {/* Header */}
        <div className="flex items-center justify-between p-4 border-b border-[#30363D]">
          <div>
            <h3 className="text-lg font-bold text-white">{contract.ticker} {contract.option_type.toUpperCase()}</h3>
            <p className="text-sm text-[#8B949E]">
              ${strike.toFixed(2)} strike · Exp {new Date(contract.expiration).toLocaleDateString()}
              <span className={`ml-2 font-medium ${moneynessColor}`}>{moneyness}</span>
            </p>
          </div>
          <button onClick={onClose} className="text-[#6E7681] hover:text-white text-xl">✕</button>
        </div>

        <div className="p-4 space-y-4">
          {step === 'configure' ? (
            <>
              {/* Price info */}
              <div className="grid grid-cols-3 gap-3">
                <div className="rounded-lg bg-[#0D1117] border border-[#30363D] p-3 text-center">
                  <p className="text-[10px] text-[#6E7681] uppercase">Bid</p>
                  <p className="text-sm font-semibold text-white">{formatCurrency(bid)}</p>
                </div>
                <div className="rounded-lg bg-[#0D1117] border border-[#50E3C2]/30 p-3 text-center">
                  <p className="text-[10px] text-[#6E7681] uppercase">Mark</p>
                  <p className="text-sm font-semibold text-[#50E3C2]">{formatCurrency(mark)}</p>
                </div>
                <div className="rounded-lg bg-[#0D1117] border border-[#30363D] p-3 text-center">
                  <p className="text-[10px] text-[#6E7681] uppercase">Ask</p>
                  <p className="text-sm font-semibold text-white">{formatCurrency(ask)}</p>
                </div>
              </div>

              {/* Side selector */}
              <div>
                <p className="text-xs text-[#6E7681] mb-2 uppercase font-medium">Order Side</p>
                <div className="grid grid-cols-2 gap-2">
                  {(Object.keys(sideInfo) as Side[]).map((s) => (
                    <button
                      key={s}
                      onClick={() => setSide(s)}
                      className={`rounded-lg border p-3 text-left transition-colors ${
                        side === s
                          ? 'border-[#50E3C2] bg-[#50E3C2]/5'
                          : 'border-[#30363D] hover:border-[#50E3C2]/30'
                      }`}
                    >
                      <p className={`text-sm font-medium ${sideInfo[s].color}`}>{sideInfo[s].label}</p>
                      <p className="text-[10px] text-[#6E7681] mt-1">{sideInfo[s].description}</p>
                    </button>
                  ))}
                </div>
              </div>

              {/* Quantity */}
              <div>
                <p className="text-xs text-[#6E7681] mb-2 uppercase font-medium">Contracts</p>
                <div className="flex items-center gap-3">
                  <button
                    onClick={() => setQuantity(Math.max(1, quantity - 1))}
                    className="w-10 h-10 rounded-lg bg-[#0D1117] border border-[#30363D] text-white font-bold hover:border-[#50E3C2]/50"
                  >
                    -
                  </button>
                  <input
                    type="number"
                    min={1}
                    max={100}
                    value={quantity}
                    onChange={(e) => setQuantity(Math.max(1, Math.min(100, parseInt(e.target.value) || 1)))}
                    className="w-20 text-center rounded-lg bg-[#0D1117] border border-[#30363D] py-2 text-white font-semibold focus:outline-none focus:border-[#50E3C2]"
                  />
                  <button
                    onClick={() => setQuantity(Math.min(100, quantity + 1))}
                    className="w-10 h-10 rounded-lg bg-[#0D1117] border border-[#30363D] text-white font-bold hover:border-[#50E3C2]/50"
                  >
                    +
                  </button>
                  <span className="text-xs text-[#6E7681]">= {quantity * 100} shares</span>
                </div>
              </div>

              {/* Cost summary */}
              <div className="rounded-lg bg-[#0D1117] border border-[#30363D] p-3 space-y-2">
                <div className="flex justify-between text-sm">
                  <span className="text-[#8B949E]">{isLong ? 'Total Cost' : 'Premium Received'}</span>
                  <span className={`font-semibold ${isLong ? 'text-red-400' : 'text-emerald-400'}`}>
                    {isLong ? '-' : '+'}{formatCurrency(totalCost)}
                  </span>
                </div>
                <div className="flex justify-between text-sm">
                  <span className="text-[#8B949E]">Break-even at expiry</span>
                  <span className="text-white font-semibold">{formatCurrency(breakEven)}</span>
                </div>
                {!isLong && (
                  <div className="flex justify-between text-sm">
                    <span className="text-[#8B949E]">Collateral required</span>
                    <span className="text-orange-400 font-semibold">
                      {contract.option_type === 'put'
                        ? formatCurrency(strike * quantity * 100)
                        : `${quantity * 100} shares`}
                    </span>
                  </div>
                )}
              </div>

              {/* Risk warning for writing */}
              {!isLong && (
                <div className="rounded-lg bg-orange-400/5 border border-orange-400/20 p-3">
                  <p className="text-xs text-orange-400 font-medium">⚠ Risk Warning</p>
                  <p className="text-[11px] text-[#8B949E] mt-1">
                    Writing options can result in losses greater than the premium received.
                    {contract.option_type === 'call'
                      ? ' Covered calls require you to own the underlying shares as collateral.'
                      : ' Cash-secured puts require cash equal to strike × 100 × contracts held as collateral.'}
                  </p>
                </div>
              )}

              {/* Greeks */}
              <GreeksDisplay
                delta={contract.delta}
                gamma={contract.gamma}
                theta={contract.theta}
                vega={contract.vega}
                rho={contract.rho}
              />

              {/* Payoff diagram */}
              <PayoffDiagram
                optionType={contract.option_type}
                strike={strike}
                premium={price}
                isLong={isLong}
                underlyingPrice={underlyingPrice}
                quantity={quantity}
              />

              <button
                onClick={() => setStep('review')}
                className="w-full py-3 rounded-lg bg-[#50E3C2] text-black font-semibold hover:bg-[#50E3C2]/90 transition-colors"
              >
                Review Order
              </button>
            </>
          ) : (
            <>
              {/* Review step */}
              <div className="rounded-lg bg-[#0D1117] border border-[#30363D] p-4 space-y-3">
                <h4 className="text-sm font-semibold text-white mb-3">Order Summary</h4>
                <div className="space-y-2 text-sm">
                  <div className="flex justify-between">
                    <span className="text-[#8B949E]">Action</span>
                    <span className={`font-medium ${sideInfo[side].color}`}>{sideInfo[side].label}</span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-[#8B949E]">Contract</span>
                    <span className="text-white">{contract.ticker} ${strike} {contract.option_type.toUpperCase()}</span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-[#8B949E]">Expiration</span>
                    <span className="text-white">{new Date(contract.expiration).toLocaleDateString()}</span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-[#8B949E]">Quantity</span>
                    <span className="text-white">{quantity} contract{quantity > 1 ? 's' : ''} ({quantity * 100} shares)</span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-[#8B949E]">Price per contract</span>
                    <span className="text-white">{formatCurrency(price)}</span>
                  </div>
                  <hr className="border-[#30363D]" />
                  <div className="flex justify-between text-base">
                    <span className="text-[#8B949E] font-medium">{isLong ? 'Total Debit' : 'Total Credit'}</span>
                    <span className={`font-bold ${isLong ? 'text-red-400' : 'text-emerald-400'}`}>
                      {formatCurrency(totalCost)}
                    </span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-[#8B949E]">Break-even</span>
                    <span className="text-yellow-400 font-semibold">{formatCurrency(breakEven)}</span>
                  </div>
                </div>
              </div>

              <p className="text-[11px] text-[#6E7681] text-center">
                This is a simulated trade for educational purposes. No real money is involved.
              </p>

              <div className="flex gap-3">
                <button
                  onClick={() => setStep('configure')}
                  className="flex-1 py-3 rounded-lg border border-[#30363D] text-[#8B949E] font-medium hover:text-white transition-colors"
                >
                  Back
                </button>
                <button
                  onClick={handleSubmit}
                  disabled={executeTrade.isPending}
                  className="flex-1 py-3 rounded-lg bg-[#50E3C2] text-black font-semibold hover:bg-[#50E3C2]/90 transition-colors disabled:opacity-50"
                >
                  {executeTrade.isPending ? 'Executing...' : 'Confirm Trade'}
                </button>
              </div>
            </>
          )}
        </div>
      </div>
    </div>
  );
}

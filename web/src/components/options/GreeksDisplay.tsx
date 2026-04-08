'use client';

import { useState } from 'react';

interface GreeksDisplayProps {
  delta: string;
  gamma: string;
  theta: string;
  vega: string;
  rho: string;
  compact?: boolean;
}

const greekInfo: Record<string, { label: string; description: string; example: string }> = {
  delta: {
    label: 'Delta (Δ)',
    description: 'Measures how much the option price changes for every $1 move in the underlying stock.',
    example: 'A delta of 0.50 means the option gains ~$0.50 when the stock rises $1. Calls have positive delta (0 to 1), puts have negative delta (-1 to 0).',
  },
  gamma: {
    label: 'Gamma (Γ)',
    description: 'The rate of change of delta — how fast delta itself moves as the stock price changes.',
    example: 'High gamma means delta changes rapidly. Options near the strike price (ATM) have the highest gamma.',
  },
  theta: {
    label: 'Theta (Θ)',
    description: 'Time decay — how much value the option loses each day, all else equal.',
    example: 'A theta of -0.05 means the option loses $0.05/day. Time decay accelerates as expiration approaches. This hurts buyers and helps sellers.',
  },
  vega: {
    label: 'Vega (ν)',
    description: 'How much the option price changes for a 1% increase in implied volatility.',
    example: 'A vega of 0.10 means the option gains $0.10 if implied volatility rises by 1%. Higher volatility = more expensive options.',
  },
  rho: {
    label: 'Rho (ρ)',
    description: 'Sensitivity to interest rate changes. Usually the smallest greek.',
    example: 'A rho of 0.05 means the option gains $0.05 if interest rates rise by 1%. Generally less important for short-term options.',
  },
};

export function GreeksDisplay({ delta, gamma, theta, vega, rho, compact }: GreeksDisplayProps) {
  const [activeTooltip, setActiveTooltip] = useState<string | null>(null);

  const greeks = [
    { key: 'delta', value: parseFloat(delta) },
    { key: 'gamma', value: parseFloat(gamma) },
    { key: 'theta', value: parseFloat(theta) },
    { key: 'vega', value: parseFloat(vega) },
    { key: 'rho', value: parseFloat(rho) },
  ];

  if (compact) {
    return (
      <div className="flex gap-3 text-xs">
        {greeks.slice(0, 4).map((g) => (
          <div key={g.key} className="flex items-center gap-1">
            <span className="text-[#6E7681] uppercase">{g.key[0]}:</span>
            <span className="text-[#8B949E] font-mono">{g.value.toFixed(4)}</span>
          </div>
        ))}
      </div>
    );
  }

  return (
    <div className="space-y-2">
      <p className="text-xs text-[#6E7681] font-medium uppercase tracking-wider">The Greeks</p>
      <div className="grid grid-cols-5 gap-2">
        {greeks.map((g) => {
          const info = greekInfo[g.key];
          return (
            <div key={g.key} className="relative">
              <button
                onClick={() => setActiveTooltip(activeTooltip === g.key ? null : g.key)}
                className="w-full text-left rounded-lg bg-[#0D1117] border border-[#30363D] p-2 hover:border-[#50E3C2]/50 transition-colors"
              >
                <div className="flex items-center gap-1 mb-1">
                  <span className="text-[10px] text-[#6E7681] uppercase">{g.key}</span>
                  <span className="text-[10px] text-[#50E3C2] cursor-help">?</span>
                </div>
                <span className={`text-sm font-mono font-semibold ${
                  g.key === 'theta' && g.value < 0 ? 'text-red-400' :
                  g.key === 'delta' ? (g.value > 0 ? 'text-emerald-400' : 'text-red-400') :
                  'text-white'
                }`}>
                  {g.value.toFixed(4)}
                </span>
              </button>

              {activeTooltip === g.key && (
                <div className="absolute z-50 top-full left-0 mt-1 w-72 rounded-lg bg-[#1C2128] border border-[#30363D] p-3 shadow-xl">
                  <p className="text-sm font-semibold text-white mb-1">{info.label}</p>
                  <p className="text-xs text-[#8B949E] mb-2">{info.description}</p>
                  <p className="text-xs text-[#6E7681] italic">{info.example}</p>
                </div>
              )}
            </div>
          );
        })}
      </div>
    </div>
  );
}

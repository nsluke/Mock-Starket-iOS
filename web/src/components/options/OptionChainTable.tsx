'use client';

import { useState } from 'react';
import { useOptionExpirations, useOptionChain } from '@/hooks/use-options';
import { EducationalBanner } from './EducationalBanner';
import { OptionsTradeModal } from './OptionsTradeModal';
import { GreeksDisplay } from './GreeksDisplay';
import { formatCurrency } from '@/lib/formatters';
import type { OptionContract } from '@/types/options';

interface OptionChainTableProps {
  ticker: string;
  underlyingPrice: number;
}

const columnTooltips: Record<string, string> = {
  Bid: 'The highest price a buyer is willing to pay for this contract right now.',
  Ask: 'The lowest price a seller is willing to accept for this contract right now.',
  Mark: 'The midpoint between bid and ask — the estimated fair value.',
  Vol: 'Number of contracts traded today.',
  OI: 'Open Interest — total number of outstanding contracts.',
  IV: 'Implied Volatility — the market\'s expectation of future price swings.',
};

export function OptionChainTable({ ticker, underlyingPrice }: OptionChainTableProps) {
  const { data: expirations = [] } = useOptionExpirations(ticker);
  const [selectedExp, setSelectedExp] = useState<string | undefined>();
  const [selectedContract, setSelectedContract] = useState<OptionContract | null>(null);
  const [hoveredCol, setHoveredCol] = useState<string | null>(null);

  // Auto-select first expiration
  const activeExp = selectedExp || expirations[0];
  const { data: chain, isLoading } = useOptionChain(ticker, activeExp);

  // Build strike-indexed map for side-by-side display
  const strikes = new Set<string>();
  const callMap = new Map<string, OptionContract>();
  const putMap = new Map<string, OptionContract>();

  if (chain) {
    chain.calls.forEach((c) => {
      strikes.add(c.strike_price);
      callMap.set(c.strike_price, c);
    });
    chain.puts.forEach((p) => {
      strikes.add(p.strike_price);
      putMap.set(p.strike_price, p);
    });
  }

  const sortedStrikes = [...strikes].sort((a, b) => parseFloat(a) - parseFloat(b));

  function moneynessLabel(strike: number, optionType: 'call' | 'put'): { label: string; color: string; bg: string } {
    const diff = Math.abs(underlyingPrice - strike) / underlyingPrice;
    if (diff < 0.01) return { label: 'ATM', color: 'text-yellow-400', bg: 'bg-yellow-400/5' };
    const itm = optionType === 'call' ? underlyingPrice > strike : underlyingPrice < strike;
    if (itm) return { label: 'ITM', color: 'text-emerald-400', bg: 'bg-emerald-400/5' };
    return { label: 'OTM', color: 'text-[#6E7681]', bg: '' };
  }

  function ContractCell({ contract, side }: { contract: OptionContract | undefined; side: 'call' | 'put' }) {
    if (!contract) return <td colSpan={4} className="px-2 py-2" />;
    const strike = parseFloat(contract.strike_price);
    const { bg } = moneynessLabel(strike, side);

    return (
      <>
        <td className={`px-2 py-2 text-right text-xs font-mono ${bg}`}>
          {formatCurrency(contract.bid_price)}
        </td>
        <td className={`px-2 py-2 text-right text-xs font-mono ${bg}`}>
          {formatCurrency(contract.ask_price)}
        </td>
        <td className={`px-2 py-2 text-right text-xs font-mono text-[#6E7681] ${bg}`}>
          {contract.volume}
        </td>
        <td className={`px-2 py-2 text-right text-xs font-mono text-[#6E7681] ${bg}`}>
          {(parseFloat(contract.implied_vol) * 100).toFixed(1)}%
        </td>
      </>
    );
  }

  function ColHeader({ label }: { label: string }) {
    const tip = columnTooltips[label];
    return (
      <th
        className="px-2 py-2 text-right text-[10px] font-medium uppercase text-[#6E7681] relative cursor-help"
        onMouseEnter={() => tip && setHoveredCol(label)}
        onMouseLeave={() => setHoveredCol(null)}
      >
        {label}
        {hoveredCol === label && tip && (
          <div className="absolute z-50 top-full right-0 mt-1 w-52 rounded-lg bg-[#1C2128] border border-[#30363D] p-2 text-left text-[11px] text-[#8B949E] font-normal normal-case shadow-xl">
            {tip}
          </div>
        )}
      </th>
    );
  }

  return (
    <div className="space-y-4">
      <EducationalBanner />

      {/* Expiration tabs */}
      {expirations.length > 0 && (
        <div className="flex gap-1.5 overflow-x-auto pb-1">
          {expirations.map((exp) => {
            const d = new Date(exp);
            const label = d.toLocaleDateString('en-US', { month: 'short', day: 'numeric' });
            const isActive = exp === activeExp;
            return (
              <button
                key={exp}
                onClick={() => setSelectedExp(exp)}
                className={`px-3 py-1.5 rounded-lg text-xs font-medium whitespace-nowrap transition-colors ${
                  isActive
                    ? 'bg-[#50E3C2]/10 text-[#50E3C2] border border-[#50E3C2]/30'
                    : 'bg-[#21262D] text-[#8B949E] hover:text-white'
                }`}
              >
                {label}
              </button>
            );
          })}
        </div>
      )}

      {isLoading && (
        <div className="text-center py-8 text-[#6E7681] text-sm">Loading option chain...</div>
      )}

      {!isLoading && sortedStrikes.length === 0 && (
        <div className="text-center py-8 text-[#6E7681] text-sm">
          No options available yet. Contracts will appear as the simulation runs.
        </div>
      )}

      {/* Chain table */}
      {sortedStrikes.length > 0 && (
        <div className="overflow-x-auto rounded-lg border border-[#30363D]">
          <table className="w-full text-sm">
            <thead>
              <tr className="border-b border-[#30363D] bg-[#0D1117]">
                <th colSpan={4} className="px-3 py-2 text-left text-xs font-semibold text-emerald-400 border-r border-[#30363D]">
                  CALLS
                </th>
                <th className="px-3 py-2 text-center text-xs font-semibold text-[#8B949E]">
                  Strike
                </th>
                <th colSpan={4} className="px-3 py-2 text-right text-xs font-semibold text-red-400 border-l border-[#30363D]">
                  PUTS
                </th>
              </tr>
              <tr className="border-b border-[#21262D] bg-[#0D1117]">
                <ColHeader label="Bid" />
                <ColHeader label="Ask" />
                <ColHeader label="Vol" />
                <ColHeader label="IV" />
                <th className="px-3 py-2 text-center text-[10px] font-medium uppercase text-[#6E7681]" />
                <ColHeader label="Bid" />
                <ColHeader label="Ask" />
                <ColHeader label="Vol" />
                <ColHeader label="IV" />
              </tr>
            </thead>
            <tbody>
              {sortedStrikes.map((strikeStr) => {
                const strike = parseFloat(strikeStr);
                const call = callMap.get(strikeStr);
                const put = putMap.get(strikeStr);
                const callM = moneynessLabel(strike, 'call');
                const putM = moneynessLabel(strike, 'put');
                const isATM = Math.abs(underlyingPrice - strike) / underlyingPrice < 0.01;

                return (
                  <tr
                    key={strikeStr}
                    className={`border-b border-[#21262D] hover:bg-[#21262D]/50 transition-colors ${
                      isATM ? 'bg-yellow-400/5' : ''
                    }`}
                  >
                    {/* Call side */}
                    <td
                      colSpan={4}
                      className={`cursor-pointer ${callM.bg}`}
                      onClick={() => call && setSelectedContract(call)}
                    >
                      {call && (
                        <table className="w-full">
                          <tbody>
                            <tr>
                              <ContractCell contract={call} side="call" />
                            </tr>
                          </tbody>
                        </table>
                      )}
                    </td>

                    {/* Strike */}
                    <td className="px-3 py-2 text-center border-x border-[#30363D]">
                      <div className="flex items-center justify-center gap-1.5">
                        <span className="text-sm font-semibold text-white">${strike.toFixed(2)}</span>
                        {(callM.label !== 'OTM' || putM.label !== 'OTM') && (
                          <span className={`text-[9px] font-bold ${isATM ? 'text-yellow-400' : callM.label === 'ITM' ? 'text-emerald-400' : 'text-red-400'}`}>
                            {isATM ? 'ATM' : callM.label === 'ITM' ? '◄ ITM' : 'ITM ►'}
                          </span>
                        )}
                      </div>
                    </td>

                    {/* Put side */}
                    <td
                      colSpan={4}
                      className={`cursor-pointer ${putM.bg}`}
                      onClick={() => put && setSelectedContract(put)}
                    >
                      {put && (
                        <table className="w-full">
                          <tbody>
                            <tr>
                              <ContractCell contract={put} side="put" />
                            </tr>
                          </tbody>
                        </table>
                      )}
                    </td>
                  </tr>
                );
              })}
            </tbody>
          </table>
        </div>
      )}

      {/* Greeks for selected row (quick preview) */}
      {selectedContract && (
        <OptionsTradeModal
          contract={selectedContract}
          underlyingPrice={underlyingPrice}
          onClose={() => setSelectedContract(null)}
        />
      )}
    </div>
  );
}

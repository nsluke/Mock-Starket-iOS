'use client';

import { useOptionsPositions, useExecuteOptionsTrade } from '@/hooks/use-options';
import { formatCurrency, formatPercent, priceChangeColor } from '@/lib/formatters';
import { GreeksDisplay } from './GreeksDisplay';

export function OptionsPositionsTable() {
  const { data: positions = [], isLoading } = useOptionsPositions();
  const closeTrade = useExecuteOptionsTrade();

  if (isLoading) {
    return <div className="p-8 text-center text-[#6E7681] text-sm">Loading options positions...</div>;
  }

  if (positions.length === 0) {
    return (
      <div className="p-8 text-center text-[#8B949E]">
        <p>No options positions yet.</p>
        <p className="text-sm text-[#6E7681] mt-1">Visit a stock detail page to browse the option chain.</p>
      </div>
    );
  }

  return (
    <div className="overflow-x-auto">
      <table className="w-full">
        <thead>
          <tr className="border-b border-[#30363D] text-xs text-[#8B949E] uppercase">
            <th className="text-left px-4 py-3">Contract</th>
            <th className="text-left px-4 py-3 hidden sm:table-cell">Type</th>
            <th className="text-right px-4 py-3">Qty</th>
            <th className="text-right px-4 py-3 hidden sm:table-cell">Avg Cost</th>
            <th className="text-right px-4 py-3">Value</th>
            <th className="text-right px-4 py-3">P&L</th>
            <th className="text-right px-4 py-3">Action</th>
          </tr>
        </thead>
        <tbody>
          {positions.map((pos) => {
            const contract = pos.contract;
            const pnl = parseFloat(pos.pnl);
            const pnlPct = parseFloat(pos.pnl_pct);
            const closeSide = pos.is_long ? 'sell_to_close' : 'buy_to_close';
            const absQty = Math.abs(pos.quantity);

            return (
              <tr key={pos.id} className="border-b border-[#21262D] hover:bg-[#21262D] transition-colors">
                <td className="px-4 py-3">
                  <div className="flex items-center gap-2">
                    <span className="rounded bg-[#50E3C2]/10 px-2 py-0.5 text-xs font-mono font-bold text-[#50E3C2]">
                      {contract.ticker}
                    </span>
                    <span className="text-sm text-white">
                      ${parseFloat(contract.strike_price).toFixed(2)} {contract.option_type.toUpperCase()}
                    </span>
                  </div>
                  <p className="text-[10px] text-[#6E7681] mt-0.5">
                    Exp {new Date(contract.expiration).toLocaleDateString()}
                  </p>
                </td>
                <td className="px-4 py-3 hidden sm:table-cell">
                  <span className={`text-xs font-medium px-2 py-0.5 rounded ${
                    pos.is_long ? 'bg-emerald-400/10 text-emerald-400' : 'bg-orange-400/10 text-orange-400'
                  }`}>
                    {pos.is_long ? 'LONG' : 'SHORT'}
                  </span>
                </td>
                <td className="px-4 py-3 text-right text-sm">
                  {absQty}
                </td>
                <td className="px-4 py-3 text-right text-sm text-[#8B949E] hidden sm:table-cell">
                  {formatCurrency(pos.avg_cost)}
                </td>
                <td className="px-4 py-3 text-right text-sm font-semibold text-white">
                  {formatCurrency(pos.market_value)}
                </td>
                <td className="px-4 py-3 text-right">
                  <div className={`text-sm font-semibold ${priceChangeColor(pnl)}`}>
                    {pnl >= 0 ? '+' : ''}{formatCurrency(pnl)}
                  </div>
                  <div className={`text-xs ${priceChangeColor(pnlPct)}`}>
                    {formatPercent(pnlPct)}
                  </div>
                </td>
                <td className="px-4 py-3 text-right">
                  <button
                    onClick={() =>
                      closeTrade.mutate({
                        contractId: pos.contract_id,
                        side: closeSide,
                        quantity: absQty,
                      })
                    }
                    disabled={closeTrade.isPending}
                    className="text-xs font-medium px-3 py-1.5 rounded-lg border border-[#30363D] text-[#8B949E] hover:text-white hover:border-[#50E3C2]/50 transition-colors disabled:opacity-50"
                  >
                    Close
                  </button>
                </td>
              </tr>
            );
          })}
        </tbody>
      </table>
    </div>
  );
}

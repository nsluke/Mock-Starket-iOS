'use client';

import { useState } from 'react';

export function EducationalBanner() {
  const [expanded, setExpanded] = useState(false);
  const [dismissed, setDismissed] = useState(false);

  if (dismissed) return null;

  return (
    <div className="rounded-lg bg-[#50E3C2]/5 border border-[#50E3C2]/20 p-4">
      <div className="flex items-start justify-between gap-3">
        <div className="flex-1">
          <div className="flex items-center gap-2 mb-1">
            <span className="text-sm">📚</span>
            <p className="text-sm font-medium text-[#50E3C2]">New to Options?</p>
          </div>
          <p className="text-xs text-[#8B949E]">
            Options give you the right to buy (call) or sell (put) a stock at a set price (strike)
            by a set date (expiration). Each contract represents 100 shares.
          </p>

          {expanded && (
            <div className="mt-3 space-y-3 text-xs text-[#8B949E]">
              <div>
                <p className="text-white font-medium mb-1">Calls vs Puts</p>
                <p><span className="text-emerald-400 font-medium">Calls</span> profit when the stock goes up. You pay a premium for the right to buy at the strike price.</p>
                <p className="mt-1"><span className="text-red-400 font-medium">Puts</span> profit when the stock goes down. You pay a premium for the right to sell at the strike price.</p>
              </div>

              <div>
                <p className="text-white font-medium mb-1">ITM / ATM / OTM</p>
                <p><span className="text-emerald-400">In-the-Money (ITM)</span>: The option has intrinsic value. For calls: stock price &gt; strike. For puts: stock price &lt; strike.</p>
                <p className="mt-1"><span className="text-yellow-400">At-the-Money (ATM)</span>: Strike price ≈ current stock price.</p>
                <p className="mt-1"><span className="text-[#6E7681]">Out-of-the-Money (OTM)</span>: No intrinsic value yet. Cheaper but riskier.</p>
              </div>

              <div>
                <p className="text-white font-medium mb-1">Buying vs Writing (Selling)</p>
                <p><span className="font-medium">Buying</span>: Limited risk (you can only lose the premium paid). Unlimited profit potential for calls.</p>
                <p className="mt-1"><span className="font-medium">Writing/Selling</span>: You collect premium upfront but take on obligation. Requires collateral (margin). Covered calls require owning the shares; cash-secured puts require holding cash equal to strike × 100.</p>
              </div>

              <div>
                <p className="text-white font-medium mb-1">Break-Even</p>
                <p>For calls: Strike Price + Premium Paid. For puts: Strike Price - Premium Paid.</p>
              </div>
            </div>
          )}
        </div>
        <button
          onClick={() => setDismissed(true)}
          className="text-[#6E7681] hover:text-white text-xs shrink-0"
        >
          ✕
        </button>
      </div>
      <button
        onClick={() => setExpanded(!expanded)}
        className="mt-2 text-xs text-[#50E3C2] hover:underline"
      >
        {expanded ? 'Show less' : 'Learn more about options →'}
      </button>
    </div>
  );
}

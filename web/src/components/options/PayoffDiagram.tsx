'use client';

import { useMemo } from 'react';

interface PayoffDiagramProps {
  optionType: 'call' | 'put';
  strike: number;
  premium: number;
  isLong: boolean;
  underlyingPrice: number;
  quantity: number;
}

export function PayoffDiagram({ optionType, strike, premium, isLong, underlyingPrice, quantity }: PayoffDiagramProps) {
  const multiplier = 100 * quantity;

  const { points, breakEven, maxLoss, maxProfit } = useMemo(() => {
    const range = strike * 0.4;
    const low = Math.max(0, strike - range);
    const high = strike + range;
    const step = (high - low) / 100;
    const pts: { x: number; y: number }[] = [];

    for (let price = low; price <= high; price += step) {
      let pnl: number;
      if (optionType === 'call') {
        const intrinsic = Math.max(price - strike, 0);
        pnl = isLong ? (intrinsic - premium) * multiplier : (premium - intrinsic) * multiplier;
      } else {
        const intrinsic = Math.max(strike - price, 0);
        pnl = isLong ? (intrinsic - premium) * multiplier : (premium - intrinsic) * multiplier;
      }
      pts.push({ x: price, y: pnl });
    }

    let be: number;
    if (optionType === 'call') {
      be = strike + premium;
    } else {
      be = strike - premium;
    }

    const ml = isLong ? -premium * multiplier : undefined; // undefined = potentially unlimited for short calls
    const mp = isLong
      ? (optionType === 'call' ? Infinity : (strike - premium) * multiplier)
      : premium * multiplier;

    return { points: pts, breakEven: be, maxLoss: ml, maxProfit: mp };
  }, [optionType, strike, premium, isLong, multiplier]);

  // SVG dimensions
  const width = 320;
  const height = 140;
  const padding = { top: 10, right: 15, bottom: 25, left: 50 };
  const plotW = width - padding.left - padding.right;
  const plotH = height - padding.top - padding.bottom;

  if (points.length === 0) return null;

  const minX = points[0].x;
  const maxX = points[points.length - 1].x;
  const minY = Math.min(...points.map((p) => p.y));
  const maxY = Math.max(...points.map((p) => p.y));
  const yRange = Math.max(maxY - minY, 1);

  const toSVG = (x: number, y: number) => ({
    sx: padding.left + ((x - minX) / (maxX - minX)) * plotW,
    sy: padding.top + (1 - (y - minY) / yRange) * plotH,
  });

  const pathD = points
    .map((p, i) => {
      const { sx, sy } = toSVG(p.x, p.y);
      return `${i === 0 ? 'M' : 'L'} ${sx} ${sy}`;
    })
    .join(' ');

  // Zero line
  const zeroY = toSVG(0, 0).sy;

  // Break-even mark
  const beSVG = toSVG(breakEven, 0);

  // Current price mark
  const cpSVG = toSVG(underlyingPrice, 0);

  return (
    <div>
      <p className="text-xs text-[#6E7681] mb-2 font-medium uppercase tracking-wider">Payoff at Expiration</p>
      <svg width={width} height={height} className="w-full" viewBox={`0 0 ${width} ${height}`}>
        {/* Zero line */}
        <line
          x1={padding.left} y1={zeroY} x2={width - padding.right} y2={zeroY}
          stroke="#30363D" strokeWidth={1} strokeDasharray="4,4"
        />

        {/* Payoff line */}
        <path d={pathD} fill="none" stroke="#50E3C2" strokeWidth={2} />

        {/* Fill profit area */}
        {points.map((p, i) => {
          if (i === 0) return null;
          const prev = points[i - 1];
          if ((p.y > 0 && prev.y > 0) || (p.y < 0 && prev.y < 0)) {
            const { sx: x1, sy: y1 } = toSVG(prev.x, prev.y);
            const { sx: x2, sy: y2 } = toSVG(p.x, p.y);
            const color = p.y > 0 ? 'rgba(74,222,128,0.1)' : 'rgba(248,113,113,0.1)';
            return (
              <polygon
                key={i}
                points={`${x1},${y1} ${x2},${y2} ${x2},${zeroY} ${x1},${zeroY}`}
                fill={color}
              />
            );
          }
          return null;
        })}

        {/* Break-even marker */}
        <line x1={beSVG.sx} y1={padding.top} x2={beSVG.sx} y2={height - padding.bottom} stroke="#FBBF24" strokeWidth={1} strokeDasharray="3,3" />
        <text x={beSVG.sx} y={height - 5} textAnchor="middle" className="fill-[#FBBF24] text-[9px]">
          BE ${breakEven.toFixed(0)}
        </text>

        {/* Current price marker */}
        <circle cx={cpSVG.sx} cy={cpSVG.sy} r={3} fill="#50E3C2" />
        <text x={cpSVG.sx} y={height - 5} textAnchor="middle" className="fill-[#8B949E] text-[9px]">
          Now
        </text>

        {/* Y-axis labels */}
        <text x={padding.left - 5} y={toSVG(0, maxY).sy + 3} textAnchor="end" className="fill-[#6E7681] text-[9px]">
          +${(maxY).toFixed(0)}
        </text>
        <text x={padding.left - 5} y={zeroY + 3} textAnchor="end" className="fill-[#6E7681] text-[9px]">
          $0
        </text>
        <text x={padding.left - 5} y={toSVG(0, minY).sy + 3} textAnchor="end" className="fill-[#6E7681] text-[9px]">
          -${Math.abs(minY).toFixed(0)}
        </text>
      </svg>

      <div className="flex justify-between text-[10px] text-[#6E7681] mt-1 px-1">
        <span>Max loss: {maxLoss !== undefined ? `$${Math.abs(maxLoss).toFixed(0)}` : 'Unlimited'}</span>
        <span>Break-even: ${breakEven.toFixed(2)}</span>
        <span>Max profit: {maxProfit === Infinity ? 'Unlimited' : `$${maxProfit.toFixed(0)}`}</span>
      </div>
    </div>
  );
}

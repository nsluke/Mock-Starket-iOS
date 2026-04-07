'use client';

import { useMemo } from 'react';

interface PieSlice {
  label: string;
  value: number;
  color: string;
}

interface PieChartProps {
  data: PieSlice[];
  size?: number;
}

const COLORS = [
  '#50E3C2', '#4ADE80', '#F87171', '#60A5FA', '#FBBF24',
  '#A78BFA', '#FB923C', '#34D399', '#F472B6', '#38BDF8',
  '#818CF8', '#FCA5A1', '#86EFAC', '#FCD34D', '#C4B5FD',
  '#FDBA74', '#6EE7B7', '#F9A8D4', '#7DD3FC', '#A5B4FC',
];

export function PieChart({ data, size = 140 }: PieChartProps) {
  const total = useMemo(() => data.reduce((sum, d) => sum + d.value, 0), [data]);

  const slices = useMemo(() => {
    let cumulative = 0;
    return data
      .filter((d) => d.value > 0)
      .sort((a, b) => b.value - a.value)
      .map((d, i) => {
        const startAngle = (cumulative / total) * 360;
        const angle = (d.value / total) * 360;
        cumulative += d.value;
        return {
          ...d,
          color: d.color || COLORS[i % COLORS.length],
          startAngle,
          angle,
          pct: (d.value / total) * 100,
        };
      });
  }, [data, total]);

  const r = size / 2;
  const cx = r;
  const cy = r;
  const outerR = r - 2;
  const innerR = r * 0.55; // donut

  function arcPath(startAngle: number, endAngle: number, outer: number, inner: number) {
    const startRad = ((startAngle - 90) * Math.PI) / 180;
    const endRad = ((endAngle - 90) * Math.PI) / 180;

    const x1 = cx + outer * Math.cos(startRad);
    const y1 = cy + outer * Math.sin(startRad);
    const x2 = cx + outer * Math.cos(endRad);
    const y2 = cy + outer * Math.sin(endRad);
    const x3 = cx + inner * Math.cos(endRad);
    const y3 = cy + inner * Math.sin(endRad);
    const x4 = cx + inner * Math.cos(startRad);
    const y4 = cy + inner * Math.sin(startRad);

    const largeArc = endAngle - startAngle > 180 ? 1 : 0;

    return `M ${x1} ${y1} A ${outer} ${outer} 0 ${largeArc} 1 ${x2} ${y2} L ${x3} ${y3} A ${inner} ${inner} 0 ${largeArc} 0 ${x4} ${y4} Z`;
  }

  if (slices.length === 0) return null;

  return (
    <div className="flex items-center gap-4">
      <svg width={size} height={size} viewBox={`0 0 ${size} ${size}`}>
        {slices.map((s, i) => {
          // Handle single-slice case (full circle)
          if (s.angle >= 359.99) {
            return (
              <circle
                key={i}
                cx={cx}
                cy={cy}
                r={(outerR + innerR) / 2}
                fill="none"
                stroke={s.color}
                strokeWidth={outerR - innerR}
              />
            );
          }
          return (
            <path
              key={i}
              d={arcPath(s.startAngle, s.startAngle + s.angle, outerR, innerR)}
              fill={s.color}
              className="transition-opacity hover:opacity-80"
            />
          );
        })}
      </svg>
      <div className="flex flex-col gap-1 min-w-0 overflow-y-auto max-h-40">
        {slices.map((s, i) => (
          <div key={i} className="flex items-center gap-2 text-xs">
            <div
              className="w-2 h-2 rounded-full shrink-0"
              style={{ backgroundColor: s.color }}
            />
            <span className="text-[#8B949E] truncate">{s.label}</span>
            <span className="text-white font-medium ml-auto whitespace-nowrap">{s.pct.toFixed(1)}%</span>
          </div>
        ))}
      </div>
    </div>
  );
}

export { COLORS };

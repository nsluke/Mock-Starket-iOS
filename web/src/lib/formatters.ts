export function formatCurrency(value: string | number): string {
  const num = typeof value === 'string' ? parseFloat(value) : value;
  return new Intl.NumberFormat('en-US', {
    style: 'currency',
    currency: 'USD',
    minimumFractionDigits: 2,
    maximumFractionDigits: 2,
  }).format(num);
}

export function formatPercent(value: string | number): string {
  const num = typeof value === 'string' ? parseFloat(value) : value;
  const sign = num >= 0 ? '+' : '';
  return `${sign}${num.toFixed(2)}%`;
}

export function formatCompact(value: string | number): string {
  const num = typeof value === 'string' ? parseFloat(value) : value;
  const abs = Math.abs(num);
  if (abs >= 1_000_000_000) return `$${(num / 1_000_000_000).toFixed(1)}B`;
  if (abs >= 1_000_000) return `$${(num / 1_000_000).toFixed(1)}M`;
  if (abs >= 1_000) return `$${(num / 1_000).toFixed(1)}K`;
  return formatCurrency(num);
}

export function formatVolume(value: number): string {
  if (value >= 1_000_000) return `${(value / 1_000_000).toFixed(1)}M`;
  if (value >= 1_000) return `${(value / 1_000).toFixed(1)}K`;
  return value.toString();
}

export function priceChangeColor(change: string | number): string {
  const num = typeof change === 'string' ? parseFloat(change) : change;
  if (num > 0) return 'text-emerald-400';
  if (num < 0) return 'text-red-400';
  return 'text-gray-400';
}

export function priceChangeBg(change: string | number): string {
  const num = typeof change === 'string' ? parseFloat(change) : change;
  if (num > 0) return 'bg-emerald-400/10 text-emerald-400';
  if (num < 0) return 'bg-red-400/10 text-red-400';
  return 'bg-gray-400/10 text-gray-400';
}

import { describe, it, expect } from 'vitest';
import { formatCurrency, formatPercent, formatCompact, formatVolume, priceChangeColor, priceChangeBg } from '../formatters';

describe('formatCurrency', () => {
  it('formats numbers as USD', () => {
    expect(formatCurrency(1234.56)).toBe('$1,234.56');
    expect(formatCurrency(0)).toBe('$0.00');
  });

  it('handles string input', () => {
    expect(formatCurrency('155.50')).toBe('$155.50');
  });

  it('formats negative values', () => {
    expect(formatCurrency(-500)).toBe('-$500.00');
  });
});

describe('formatPercent', () => {
  it('adds + sign for positive values', () => {
    expect(formatPercent(5.25)).toBe('+5.25%');
  });

  it('shows - sign for negative values', () => {
    expect(formatPercent(-3.1)).toBe('-3.10%');
  });

  it('handles zero', () => {
    expect(formatPercent(0)).toBe('+0.00%');
  });

  it('handles string input', () => {
    expect(formatPercent('12.5')).toBe('+12.50%');
  });
});

describe('formatCompact', () => {
  it('formats billions', () => {
    expect(formatCompact(2500000000)).toBe('$2.5B');
  });

  it('formats millions', () => {
    expect(formatCompact(1500000)).toBe('$1.5M');
  });

  it('formats thousands', () => {
    expect(formatCompact(45000)).toBe('$45.0K');
  });

  it('falls back to formatCurrency for small values', () => {
    expect(formatCompact(500)).toBe('$500.00');
  });
});

describe('formatVolume', () => {
  it('formats millions', () => {
    expect(formatVolume(1250000)).toBe('1.3M');
  });

  it('formats thousands', () => {
    expect(formatVolume(5000)).toBe('5.0K');
  });

  it('returns raw number for small values', () => {
    expect(formatVolume(500)).toBe('500');
  });
});

describe('priceChangeColor', () => {
  it('returns green for positive', () => {
    expect(priceChangeColor(5)).toBe('text-emerald-400');
  });

  it('returns red for negative', () => {
    expect(priceChangeColor(-3)).toBe('text-red-400');
  });

  it('returns gray for zero', () => {
    expect(priceChangeColor(0)).toBe('text-gray-400');
  });

  it('handles string input', () => {
    expect(priceChangeColor('1.5')).toBe('text-emerald-400');
    expect(priceChangeColor('-0.5')).toBe('text-red-400');
  });
});

describe('priceChangeBg', () => {
  it('returns green bg for positive', () => {
    expect(priceChangeBg(5)).toContain('emerald');
  });

  it('returns red bg for negative', () => {
    expect(priceChangeBg(-3)).toContain('red');
  });
});

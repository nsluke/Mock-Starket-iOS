import { describe, it, expect } from 'vitest';
import { tradeSchema, orderSchema, alertSchema, profileSchema } from '../schemas';

describe('tradeSchema', () => {
  it('accepts valid shares', () => {
    const result = tradeSchema.safeParse({ shares: 10 });
    expect(result.success).toBe(true);
  });

  it('coerces string to number', () => {
    const result = tradeSchema.safeParse({ shares: '25' });
    expect(result.success).toBe(true);
    if (result.success) expect(result.data.shares).toBe(25);
  });

  it('rejects zero shares', () => {
    const result = tradeSchema.safeParse({ shares: 0 });
    expect(result.success).toBe(false);
  });

  it('rejects negative shares', () => {
    const result = tradeSchema.safeParse({ shares: -5 });
    expect(result.success).toBe(false);
  });

  it('rejects non-integer shares', () => {
    const result = tradeSchema.safeParse({ shares: 1.5 });
    expect(result.success).toBe(false);
  });
});

describe('orderSchema', () => {
  it('accepts valid limit order', () => {
    const result = orderSchema.safeParse({
      ticker: 'PIPE',
      side: 'buy',
      order_type: 'limit',
      shares: 10,
      limit_price: '150.00',
    });
    expect(result.success).toBe(true);
  });

  it('accepts valid stop order', () => {
    const result = orderSchema.safeParse({
      ticker: 'PIPE',
      side: 'sell',
      order_type: 'stop',
      shares: 5,
      stop_price: '140.00',
    });
    expect(result.success).toBe(true);
  });

  it('accepts valid stop_limit order', () => {
    const result = orderSchema.safeParse({
      ticker: 'PIPE',
      side: 'buy',
      order_type: 'stop_limit',
      shares: 10,
      limit_price: '150.00',
      stop_price: '145.00',
    });
    expect(result.success).toBe(true);
  });

  it('rejects limit order without limit_price', () => {
    const result = orderSchema.safeParse({
      ticker: 'PIPE',
      side: 'buy',
      order_type: 'limit',
      shares: 10,
    });
    expect(result.success).toBe(false);
  });

  it('rejects stop order without stop_price', () => {
    const result = orderSchema.safeParse({
      ticker: 'PIPE',
      side: 'sell',
      order_type: 'stop',
      shares: 5,
    });
    expect(result.success).toBe(false);
  });

  it('rejects empty ticker', () => {
    const result = orderSchema.safeParse({
      ticker: '',
      side: 'buy',
      order_type: 'limit',
      shares: 10,
      limit_price: '100.00',
    });
    expect(result.success).toBe(false);
  });
});

describe('alertSchema', () => {
  it('accepts valid alert', () => {
    const result = alertSchema.safeParse({
      ticker: 'PIPE',
      condition: 'above',
      target_price: '160.00',
    });
    expect(result.success).toBe(true);
  });

  it('rejects invalid price format', () => {
    const result = alertSchema.safeParse({
      ticker: 'PIPE',
      condition: 'below',
      target_price: 'abc',
    });
    expect(result.success).toBe(false);
  });

  it('rejects empty target price', () => {
    const result = alertSchema.safeParse({
      ticker: 'PIPE',
      condition: 'above',
      target_price: '',
    });
    expect(result.success).toBe(false);
  });
});

describe('profileSchema', () => {
  it('accepts valid name', () => {
    const result = profileSchema.safeParse({ display_name: 'TraderJoe' });
    expect(result.success).toBe(true);
  });

  it('rejects empty name', () => {
    const result = profileSchema.safeParse({ display_name: '' });
    expect(result.success).toBe(false);
  });

  it('rejects name over 30 chars', () => {
    const result = profileSchema.safeParse({ display_name: 'a'.repeat(31) });
    expect(result.success).toBe(false);
  });
});

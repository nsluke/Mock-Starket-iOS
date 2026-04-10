import { describe, it, expect, vi, beforeEach } from 'vitest';

// We test that the API client properly encodes ticker symbols in URL paths.
// Tickers like "X:BTCUSD" contain colons that break URL routing if not encoded.

describe('APIClient URL encoding', () => {
  let fetchSpy: ReturnType<typeof vi.fn>;

  beforeEach(() => {
    fetchSpy = vi.fn().mockResolvedValue({
      ok: true,
      json: () => Promise.resolve({}),
    });
    vi.stubGlobal('fetch', fetchSpy);
  });

  // Dynamic import so env var and fetch stub are in place
  async function getClient() {
    const mod = await import('@/lib/api-client');
    return mod.apiClient;
  }

  it('encodes crypto tickers with colons in getStock', async () => {
    const client = await getClient();
    await client.getStock('X:BTCUSD');

    const url: string = fetchSpy.mock.calls[0][0];
    // The colon must be encoded so the server sees the full ticker as one path segment
    expect(url).toContain(encodeURIComponent('X:BTCUSD'));
    expect(url).not.toContain('/stocks/X:BTCUSD');
  });

  it('encodes crypto tickers in getStockHistory', async () => {
    const client = await getClient();
    await client.getStockHistory('X:ETHUSD', '1m');

    const url: string = fetchSpy.mock.calls[0][0];
    expect(url).toContain(encodeURIComponent('X:ETHUSD'));
  });

  it('encodes crypto tickers in getETFHoldings', async () => {
    const client = await getClient();
    await client.getETFHoldings('X:BTCUSD');

    const url: string = fetchSpy.mock.calls[0][0];
    expect(url).toContain(encodeURIComponent('X:BTCUSD'));
  });

  it('encodes crypto tickers in getOptionExpirations', async () => {
    const client = await getClient();
    await client.getOptionExpirations('X:BTCUSD');

    const url: string = fetchSpy.mock.calls[0][0];
    expect(url).toContain(encodeURIComponent('X:BTCUSD'));
  });

  it('encodes crypto tickers in getOptionChain', async () => {
    const client = await getClient();
    await client.getOptionChain('X:BTCUSD', '2026-05-01');

    const url: string = fetchSpy.mock.calls[0][0];
    expect(url).toContain(encodeURIComponent('X:BTCUSD'));
  });

  it('leaves normal tickers unchanged', async () => {
    const client = await getClient();
    await client.getStock('AAPL');

    const url: string = fetchSpy.mock.calls[0][0];
    expect(url).toContain('/stocks/AAPL');
  });

  it('handles tickers with dots (BRK.B)', async () => {
    const client = await getClient();
    await client.getStock('BRK.B');

    const url: string = fetchSpy.mock.calls[0][0];
    // Dots are safe in URL paths, but encoding them is also fine
    expect(url).toMatch(/\/stocks\/(BRK\.B|BRK%2EB)/);
  });

  it('encodes tickers in watchlist removal', async () => {
    const client = await getClient();
    await client.removeFromWatchlist('X:BTCUSD');

    const url: string = fetchSpy.mock.calls[0][0];
    expect(url).toContain(encodeURIComponent('X:BTCUSD'));
  });
});

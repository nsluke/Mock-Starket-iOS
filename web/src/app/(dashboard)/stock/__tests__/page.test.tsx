import React from 'react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, waitFor } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { http, HttpResponse } from 'msw';
import { server } from '@/test/mocks/server';
import { useMarketStore } from '@/stores/market-store';

// Mock next/navigation
const mockTicker = { ticker: 'AAPL' };
vi.mock('next/navigation', () => ({
  useParams: () => mockTicker,
  useRouter: () => ({ push: vi.fn(), back: vi.fn() }),
  usePathname: () => '/stock/AAPL',
}));

// Mock next/dynamic to render children directly
vi.mock('next/dynamic', () => ({
  default: () => {
    return function DynamicMock() {
      return <div data-testid="price-chart-mock" />;
    };
  },
}));

// Mock next/link
vi.mock('next/link', () => ({
  default: ({ children, href, ...props }: any) => (
    <a href={href} {...props}>{children}</a>
  ),
}));

// Import the component using the @ alias which resolves through vitest.config
import StockDetailPage from '@/app/(dashboard)/stock/[ticker]/page';

const API_URL = 'http://localhost';

const mockStock = {
  ticker: 'AAPL',
  name: 'Apple Inc.',
  sector: 'Technology',
  asset_type: 'stock',
  base_price: '185.00',
  current_price: '192.50',
  day_open: '190.00',
  day_high: '194.00',
  day_low: '189.00',
  prev_close: '189.50',
  volume: 52000000,
  volatility: '0.02',
  description: 'Consumer electronics and software company',
  logo_url: 'https://example.com/aapl.png',
};

const mockCryptoStock = {
  ticker: 'X:BTCUSD',
  name: 'Bitcoin',
  sector: 'Crypto',
  asset_type: 'crypto',
  base_price: '45000.00',
  current_price: '67500.00',
  day_open: '66000.00',
  day_high: '68000.00',
  day_low: '65500.00',
  prev_close: '65000.00',
  volume: 1500000,
  volatility: '0.05',
  description: 'Decentralized digital currency',
};

const mockStockMinimal = {
  ticker: 'TEST',
  name: 'Test Corp',
  sector: 'Technology',
  asset_type: 'stock',
  base_price: '100.00',
  current_price: '105.00',
  day_open: '100.00',
  day_high: '106.00',
  day_low: '99.00',
  prev_close: '99.50',
  volume: 100000,
  volatility: '0.02',
  // No description, no logo_url
};

function createQueryClient() {
  return new QueryClient({
    defaultOptions: {
      queries: { retry: false, gcTime: 0 },
      mutations: { retry: false },
    },
  });
}

function renderStockPage() {
  const queryClient = createQueryClient();
  return render(
    <QueryClientProvider client={queryClient}>
      <StockDetailPage />
    </QueryClientProvider>
  );
}

function setupHandlers(stock: any = mockStock) {
  server.use(
    http.get(`${API_URL}/api/v1/stocks/:ticker`, () =>
      HttpResponse.json(stock)
    ),
    http.get(`${API_URL}/api/v1/stocks/:ticker/history`, () =>
      HttpResponse.json([])
    ),
    http.get(`${API_URL}/api/v1/stocks/:ticker/options/expirations`, () =>
      HttpResponse.json([])
    ),
    http.get(`${API_URL}/api/v1/stocks/:ticker/options`, () =>
      HttpResponse.json({ calls: [], puts: [] })
    ),
    http.get(`${API_URL}/api/v1/portfolio`, () =>
      HttpResponse.json({
        portfolio: { id: 'p1', user_id: 'u1', cash: '100000.00' },
        positions: [],
        net_worth: '100000.00',
        invested: '0',
      })
    ),
    http.get(`${API_URL}/api/v1/watchlist`, () =>
      HttpResponse.json([])
    ),
    http.get(`${API_URL}/api/v1/market/status`, () =>
      HttpResponse.json({ is_open: true, session: 'regular', next_open: '', next_close: '' })
    ),
  );
}

describe('StockDetailPage', () => {
  beforeEach(() => {
    mockTicker.ticker = 'AAPL';
    useMarketStore.setState({ stocks: [] });
  });

  it('renders stock name and price after loading', async () => {
    setupHandlers(mockStock);
    renderStockPage();

    await waitFor(() => {
      expect(screen.getByText('Apple Inc.')).toBeInTheDocument();
    });

    expect(screen.getByText('Technology')).toBeInTheDocument();
  });

  it('shows loading state initially', () => {
    setupHandlers(mockStock);
    renderStockPage();

    expect(screen.getByText('Loading...')).toBeInTheDocument();
  });

  it('shows "Stock not found" on 404', async () => {
    server.use(
      http.get(`${API_URL}/api/v1/stocks/:ticker`, () =>
        new HttpResponse('not found', { status: 404 })
      ),
      http.get(`${API_URL}/api/v1/stocks/:ticker/history`, () =>
        HttpResponse.json([])
      ),
      http.get(`${API_URL}/api/v1/stocks/:ticker/options/expirations`, () =>
        HttpResponse.json([])
      ),
      http.get(`${API_URL}/api/v1/stocks/:ticker/options`, () =>
        HttpResponse.json({ calls: [], puts: [] })
      ),
      http.get(`${API_URL}/api/v1/portfolio`, () =>
        HttpResponse.json({
          portfolio: { id: 'p1', user_id: 'u1', cash: '100000.00' },
          positions: [],
          net_worth: '100000.00',
          invested: '0',
        })
      ),
      http.get(`${API_URL}/api/v1/watchlist`, () =>
        HttpResponse.json([])
      ),
    );
    renderStockPage();

    await waitFor(() => {
      expect(screen.getByText('Stock not found')).toBeInTheDocument();
    });
  });

  it('renders stock with missing optional fields (no description, no logo)', async () => {
    setupHandlers(mockStockMinimal);
    mockTicker.ticker = 'TEST';
    renderStockPage();

    await waitFor(() => {
      expect(screen.getByText('Test Corp')).toBeInTheDocument();
    });

    // Should not crash when description is undefined
    expect(screen.queryByText('undefined')).not.toBeInTheDocument();
  });

  it('renders crypto stock page without crashing', async () => {
    mockTicker.ticker = 'X:BTCUSD';
    setupHandlers(mockCryptoStock);
    renderStockPage();

    await waitFor(() => {
      expect(screen.getByText('Bitcoin')).toBeInTheDocument();
    });

    // Options chain should NOT appear for crypto
    expect(screen.queryByText('Options Chain')).not.toBeInTheDocument();
  });

  it('handles empty market store gracefully', async () => {
    useMarketStore.setState({ stocks: [] });
    setupHandlers(mockStock);
    renderStockPage();

    await waitFor(() => {
      expect(screen.getByText('Apple Inc.')).toBeInTheDocument();
    });

    // Should still show price from API data even if market store is empty
    expect(screen.queryByText('$NaN')).not.toBeInTheDocument();
  });

  it('shows Quick Buy buttons', async () => {
    setupHandlers(mockStock);
    renderStockPage();

    await waitFor(() => {
      expect(screen.getByText('Quick Buy')).toBeInTheDocument();
    });

    expect(screen.getByText('10')).toBeInTheDocument();
    expect(screen.getByText('100')).toBeInTheDocument();
    expect(screen.getByText('1000')).toBeInTheDocument();
  });

  it('shows buy/sell toggle', async () => {
    setupHandlers(mockStock);
    renderStockPage();

    await waitFor(() => {
      expect(screen.getByRole('heading', { name: /Trade AAPL/i })).toBeInTheDocument();
    });

    expect(screen.getByText('Buy')).toBeInTheDocument();
    expect(screen.getByText('Sell')).toBeInTheDocument();
  });

  it('shows Options Chain section for stock asset type', async () => {
    setupHandlers(mockStock);
    renderStockPage();

    await waitFor(() => {
      expect(screen.getByText('Options Chain')).toBeInTheDocument();
    });
  });

  it('handles stock with zero day_open without NaN', async () => {
    const stockZeroOpen = {
      ...mockStock,
      day_open: '0',
      current_price: '100.00',
    };
    setupHandlers(stockZeroOpen);
    renderStockPage();

    await waitFor(() => {
      expect(screen.getByText('Apple Inc.')).toBeInTheDocument();
    });

    // change % should be 0 (not NaN or Infinity) when day_open is 0
    expect(screen.queryByText('NaN')).not.toBeInTheDocument();
    expect(screen.queryByText('Infinity')).not.toBeInTheDocument();
  });

  it('handles portfolio with null positions array', async () => {
    server.use(
      http.get(`${API_URL}/api/v1/stocks/:ticker`, () =>
        HttpResponse.json(mockStock)
      ),
      http.get(`${API_URL}/api/v1/stocks/:ticker/history`, () =>
        HttpResponse.json([])
      ),
      http.get(`${API_URL}/api/v1/stocks/:ticker/options/expirations`, () =>
        HttpResponse.json([])
      ),
      http.get(`${API_URL}/api/v1/stocks/:ticker/options`, () =>
        HttpResponse.json({ calls: [], puts: [] })
      ),
      http.get(`${API_URL}/api/v1/portfolio`, () =>
        HttpResponse.json({
          portfolio: { id: 'p1', user_id: 'u1', cash: '100000.00' },
          positions: null,
          net_worth: '100000.00',
          invested: '0',
        })
      ),
      http.get(`${API_URL}/api/v1/watchlist`, () =>
        HttpResponse.json([])
      ),
    );
    renderStockPage();

    await waitFor(() => {
      expect(screen.getByText('Apple Inc.')).toBeInTheDocument();
    });

    // Should not crash with null positions
    expect(screen.queryByText('You own')).not.toBeInTheDocument();
  });
});

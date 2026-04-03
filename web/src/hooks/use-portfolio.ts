import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { apiClient } from '@/lib/api-client';
import { toast } from '@/lib/toast';
import type { PortfolioResponse, Trade } from '@/types/portfolio';

export function usePortfolio() {
  return useQuery<PortfolioResponse>({
    queryKey: ['portfolio'],
    queryFn: () => apiClient.getPortfolio(),
  });
}

export function usePortfolioHistory(limit = 100) {
  return useQuery<any[]>({
    queryKey: ['portfolio-history', limit],
    queryFn: () => apiClient.getPortfolioHistory(limit),
  });
}

export function useTradeHistory(limit = 50, offset = 0) {
  return useQuery<Trade[]>({
    queryKey: ['trades', limit, offset],
    queryFn: async () => {
      const data = await apiClient.getTradeHistory(limit, offset);
      return data || [];
    },
  });
}

export function useExecuteTrade() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ ticker, side, shares }: { ticker: string; side: string; shares: number }) =>
      apiClient.executeTrade(ticker, side, shares),
    onSuccess: (_data, variables) => {
      const action = variables.side === 'buy' ? 'Bought' : 'Sold';
      toast.success(`${action} ${variables.shares} shares of ${variables.ticker}`);
      queryClient.invalidateQueries({ queryKey: ['portfolio'] });
      queryClient.invalidateQueries({ queryKey: ['trades'] });
    },
    onError: (error: Error) => {
      toast.error(error.message || 'Trade failed');
    },
  });
}

import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { apiClient } from '@/lib/api-client';
import { toast } from '@/lib/toast';
import type { OptionChainResponse, OptionPosition, OptionTrade, OptionOrder } from '@/types/options';

export function useOptionExpirations(ticker: string) {
  return useQuery<string[]>({
    queryKey: ['option-expirations', ticker],
    queryFn: () => apiClient.getOptionExpirations(ticker),
    enabled: !!ticker,
  });
}

export function useOptionChain(ticker: string, expiration?: string) {
  return useQuery<OptionChainResponse>({
    queryKey: ['option-chain', ticker, expiration],
    queryFn: () => apiClient.getOptionChain(ticker, expiration),
    enabled: !!ticker,
    refetchInterval: 10000, // refresh every 10s
  });
}

export function useOptionsPositions() {
  return useQuery<OptionPosition[]>({
    queryKey: ['options-positions'],
    queryFn: () => apiClient.getOptionsPositions(),
  });
}

export function useOptionsTradeHistory(limit = 50, offset = 0) {
  return useQuery<OptionTrade[]>({
    queryKey: ['options-trades', limit, offset],
    queryFn: () => apiClient.getOptionsTradeHistory(limit, offset),
  });
}

export function useOptionsOrders() {
  return useQuery<OptionOrder[]>({
    queryKey: ['options-orders'],
    queryFn: () => apiClient.getOptionsOrders(),
  });
}

export function useExecuteOptionsTrade() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ contractId, side, quantity }: { contractId: string; side: string; quantity: number }) =>
      apiClient.executeOptionsTrade(contractId, side, quantity),
    onSuccess: (_data, variables) => {
      const sideLabel = variables.side.replace(/_/g, ' ');
      toast.success(`Options trade executed: ${sideLabel} ${variables.quantity} contract(s)`);
      queryClient.invalidateQueries({ queryKey: ['options-positions'] });
      queryClient.invalidateQueries({ queryKey: ['options-trades'] });
      queryClient.invalidateQueries({ queryKey: ['portfolio'] });
    },
    onError: (error: Error) => {
      toast.error(error.message || 'Options trade failed');
    },
  });
}

export function useCancelOptionsOrder() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => apiClient.cancelOptionsOrder(id),
    onSuccess: () => {
      toast.success('Options order cancelled');
      queryClient.invalidateQueries({ queryKey: ['options-orders'] });
    },
    onError: (error: Error) => {
      toast.error(error.message || 'Failed to cancel order');
    },
  });
}

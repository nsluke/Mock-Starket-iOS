import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { apiClient } from '@/lib/api-client';
import { toast } from '@/lib/toast';

export function useWatchlist() {
  return useQuery<string[]>({
    queryKey: ['watchlist'],
    queryFn: async () => {
      const data = await apiClient.getWatchlist();
      return data || [];
    },
  });
}

export function useAddToWatchlist() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (ticker: string) => apiClient.addToWatchlist(ticker),
    onMutate: async (ticker) => {
      await queryClient.cancelQueries({ queryKey: ['watchlist'] });
      const previous = queryClient.getQueryData<string[]>(['watchlist']);
      queryClient.setQueryData<string[]>(['watchlist'], (old) =>
        old ? [...old, ticker] : [ticker]
      );
      return { previous };
    },
    onError: (_err, _ticker, context) => {
      queryClient.setQueryData(['watchlist'], context?.previous);
      toast.error('Failed to add to watchlist');
    },
    onSuccess: (_data, ticker) => {
      toast.success(`${ticker} added to watchlist`);
    },
    onSettled: () => {
      queryClient.invalidateQueries({ queryKey: ['watchlist'] });
    },
  });
}

export function useRemoveFromWatchlist() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (ticker: string) => apiClient.removeFromWatchlist(ticker),
    onMutate: async (ticker) => {
      await queryClient.cancelQueries({ queryKey: ['watchlist'] });
      const previous = queryClient.getQueryData<string[]>(['watchlist']);
      queryClient.setQueryData<string[]>(['watchlist'], (old) =>
        old?.filter((t) => t !== ticker)
      );
      return { previous };
    },
    onError: (_err, _ticker, context) => {
      queryClient.setQueryData(['watchlist'], context?.previous);
      toast.error('Failed to remove from watchlist');
    },
    onSuccess: (_data, ticker) => {
      toast.success(`${ticker} removed from watchlist`);
    },
    onSettled: () => {
      queryClient.invalidateQueries({ queryKey: ['watchlist'] });
    },
  });
}

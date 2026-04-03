import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { apiClient } from '@/lib/api-client';
import { toast } from '@/lib/toast';

export interface PriceAlert {
  id: string;
  ticker: string;
  condition: string;
  target_price: string;
  triggered: boolean;
  triggered_at: string | null;
  created_at: string;
}

export function useAlerts() {
  return useQuery<PriceAlert[]>({
    queryKey: ['alerts'],
    queryFn: async () => {
      const data = await apiClient.getAlerts();
      return data || [];
    },
  });
}

export function useCreateAlert() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (params: { ticker: string; condition: string; targetPrice: string }) =>
      apiClient.createAlert(params.ticker, params.condition, params.targetPrice),
    onSuccess: () => {
      toast.success('Price alert created');
      queryClient.invalidateQueries({ queryKey: ['alerts'] });
    },
    onError: (error: Error) => {
      toast.error(error.message || 'Failed to create alert');
    },
  });
}

export function useDeleteAlert() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => apiClient.deleteAlert(id),
    onMutate: async (id) => {
      await queryClient.cancelQueries({ queryKey: ['alerts'] });
      const previous = queryClient.getQueryData<PriceAlert[]>(['alerts']);
      queryClient.setQueryData<PriceAlert[]>(['alerts'], (old) =>
        old?.filter((a) => a.id !== id)
      );
      return { previous };
    },
    onError: (_err, _id, context) => {
      queryClient.setQueryData(['alerts'], context?.previous);
      toast.error('Failed to delete alert');
    },
    onSettled: () => {
      queryClient.invalidateQueries({ queryKey: ['alerts'] });
    },
  });
}

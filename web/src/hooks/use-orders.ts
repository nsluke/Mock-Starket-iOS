import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { apiClient } from '@/lib/api-client';
import { toast } from '@/lib/toast';
import type { Order } from '@/types/portfolio';

export function useOrders() {
  return useQuery<Order[]>({
    queryKey: ['orders'],
    queryFn: async () => {
      const data = await apiClient.getOrders();
      return data || [];
    },
  });
}

export function useCreateOrder() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (params: {
      ticker: string;
      side: string;
      orderType: string;
      shares: number;
      limitPrice?: string;
      stopPrice?: string;
    }) =>
      apiClient.createOrder(
        params.ticker,
        params.side,
        params.orderType,
        params.shares,
        params.limitPrice,
        params.stopPrice
      ),
    onSuccess: () => {
      toast.success('Order placed successfully');
      queryClient.invalidateQueries({ queryKey: ['orders'] });
    },
    onError: (error: Error) => {
      toast.error(error.message || 'Failed to create order');
    },
  });
}

export function useCancelOrder() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => apiClient.cancelOrder(id),
    onMutate: async (id) => {
      await queryClient.cancelQueries({ queryKey: ['orders'] });
      const previous = queryClient.getQueryData<Order[]>(['orders']);
      queryClient.setQueryData<Order[]>(['orders'], (old) =>
        old?.filter((o) => o.id !== id)
      );
      return { previous };
    },
    onError: (_err, _id, context) => {
      queryClient.setQueryData(['orders'], context?.previous);
      toast.error('Failed to cancel order');
    },
    onSuccess: () => {
      toast.success('Order cancelled');
    },
    onSettled: () => {
      queryClient.invalidateQueries({ queryKey: ['orders'] });
    },
  });
}

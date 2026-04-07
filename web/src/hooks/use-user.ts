import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { apiClient } from '@/lib/api-client';
import { toast } from '@/lib/toast';
import type { User } from '@/types/user';

export function useCurrentUser() {
  return useQuery<User>({
    queryKey: ['user'],
    queryFn: () => apiClient.getMe(),
  });
}

export function useUpdateProfile() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (displayName: string) =>
      apiClient.request('/api/v1/auth/me', {
        method: 'PUT',
        body: JSON.stringify({ display_name: displayName }),
      }),
    onSuccess: () => {
      toast.success('Profile updated');
      queryClient.invalidateQueries({ queryKey: ['user'] });
    },
    onError: () => {
      toast.error('Failed to update profile');
    },
  });
}

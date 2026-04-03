import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { apiClient } from '@/lib/api-client';
import { toast } from '@/lib/toast';

interface DailyChallenge {
  id: string;
  date: string;
  challenge_type: string;
  description: string;
  reward_cash: string;
}

interface UserProgress {
  completed: boolean;
  completed_at: string | null;
  claimed: boolean;
}

export interface ChallengeData {
  challenge: DailyChallenge | null;
  progress: UserProgress | null;
}

export function useTodaysChallenge() {
  return useQuery<ChallengeData>({
    queryKey: ['challenge'],
    queryFn: async () => {
      try {
        const data = await apiClient.getTodaysChallenge();
        return { challenge: data.challenge, progress: data.progress };
      } catch {
        return { challenge: null, progress: null };
      }
    },
  });
}

export function useCheckChallenge() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: () => apiClient.checkChallenge(),
    onSuccess: (result) => {
      if (result.completed) {
        toast.success('Challenge completed! Claim your reward.');
      } else {
        toast.info('Not yet complete. Keep trading!');
      }
      queryClient.invalidateQueries({ queryKey: ['challenge'] });
    },
    onError: () => {
      toast.error('Failed to check progress');
    },
  });
}

export function useClaimChallenge() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => apiClient.claimChallenge(id),
    onSuccess: () => {
      toast.success('Reward claimed! Cash added to your portfolio.');
      queryClient.invalidateQueries({ queryKey: ['challenge'] });
      queryClient.invalidateQueries({ queryKey: ['portfolio'] });
    },
    onError: (error: Error) => {
      toast.error(error.message || 'Failed to claim reward');
    },
  });
}

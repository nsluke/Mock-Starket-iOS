import { useQuery } from '@tanstack/react-query';
import { apiClient } from '@/lib/api-client';
import type { LeaderboardEntry } from '@/types/user';

export function useLeaderboard(period: string) {
  return useQuery<LeaderboardEntry[]>({
    queryKey: ['leaderboard', period],
    queryFn: async () => {
      const data = await apiClient.getLeaderboard(period);
      return data || [];
    },
  });
}

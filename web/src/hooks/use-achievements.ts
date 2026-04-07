import { useQuery } from '@tanstack/react-query';
import { apiClient } from '@/lib/api-client';
import type { Achievement } from '@/types/user';

export interface UserAchievement {
  id: string;
  achievement_id: string;
  earned_at: string;
}

export function useAchievements() {
  return useQuery<{ achievements: Achievement[]; earned: UserAchievement[] }>({
    queryKey: ['achievements'],
    queryFn: async () => {
      const [achievements, earned] = await Promise.all([
        apiClient.getAchievements(),
        apiClient.getMyAchievements(),
      ]);
      return {
        achievements: achievements || [],
        earned: earned || [],
      };
    },
  });
}

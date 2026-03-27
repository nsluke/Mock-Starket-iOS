'use client';

import { useEffect, useState } from 'react';
import { apiClient } from '@/lib/api-client';
import type { Achievement } from '@/types/user';

interface UserAchievement {
  id: string;
  achievement_id: string;
  earned_at: string;
}

const categoryLabels: Record<string, string> = {
  trading: 'Trading',
  portfolio: 'Portfolio',
  social: 'Social',
  streak: 'Streaks',
  special: 'Special',
  skill: 'Skill',
};

export default function AchievementsPage() {
  const [achievements, setAchievements] = useState<Achievement[]>([]);
  const [earned, setEarned] = useState<UserAchievement[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    async function load() {
      setLoading(true);
      try {
        const [allAchievements, myAchievements] = await Promise.all([
          apiClient.getAchievements(),
          apiClient.getMyAchievements(),
        ]);
        setAchievements(allAchievements || []);
        setEarned(myAchievements || []);
      } catch (err) {
        console.error('Failed to load achievements:', err);
      } finally {
        setLoading(false);
      }
    }
    load();
  }, []);

  const earnedIds = new Set(earned.map((e) => e.achievement_id));
  const earnedCount = earnedIds.size;
  const totalCount = achievements.length;

  // Group by category
  const grouped = achievements.reduce<Record<string, Achievement[]>>((acc, a) => {
    if (!acc[a.category]) acc[a.category] = [];
    acc[a.category].push(a);
    return acc;
  }, {});

  if (loading) {
    return <div className="p-6 text-center text-[#8B949E]">Loading achievements...</div>;
  }

  return (
    <div className="p-6 max-w-4xl mx-auto space-y-6">
      <div className="flex items-baseline justify-between">
        <h1 className="text-2xl font-bold">Achievements</h1>
        <span className="text-sm text-[#8B949E]">
          {earnedCount} / {totalCount} earned
        </span>
      </div>

      {/* Progress bar */}
      <div className="rounded-full bg-[#21262D] h-2 overflow-hidden">
        <div
          className="h-full bg-[#50E3C2] rounded-full transition-all duration-500"
          style={{ width: totalCount > 0 ? `${(earnedCount / totalCount) * 100}%` : '0%' }}
        />
      </div>

      {/* Achievement Categories */}
      {Object.entries(grouped).map(([category, items]) => (
        <div key={category} className="space-y-3">
          <h2 className="text-sm font-semibold text-[#8B949E] uppercase tracking-wider">
            {categoryLabels[category] || category}
          </h2>

          <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
            {items.map((achievement) => {
              const isEarned = earnedIds.has(achievement.id);
              const earnedData = earned.find((e) => e.achievement_id === achievement.id);

              return (
                <div
                  key={achievement.id}
                  className={`rounded-xl border p-4 transition-colors ${
                    isEarned
                      ? 'bg-[#161B22] border-[#50E3C2]/30'
                      : 'bg-[#161B22] border-[#30363D] opacity-60'
                  }`}
                >
                  <div className="flex items-start gap-3">
                    <div
                      className={`text-2xl flex-shrink-0 ${isEarned ? '' : 'grayscale'}`}
                    >
                      {isEarned ? '&#x2714;' : '&#x1F512;'}
                    </div>
                    <div className="min-w-0">
                      <h3 className={`font-semibold text-sm ${isEarned ? 'text-white' : 'text-[#8B949E]'}`}>
                        {achievement.name}
                      </h3>
                      <p className="text-xs text-[#6E7681] mt-0.5">{achievement.description}</p>
                      {isEarned && earnedData && (
                        <p className="text-xs text-[#50E3C2] mt-1">
                          Earned {new Date(earnedData.earned_at).toLocaleDateString()}
                        </p>
                      )}
                    </div>
                  </div>
                </div>
              );
            })}
          </div>
        </div>
      ))}
    </div>
  );
}

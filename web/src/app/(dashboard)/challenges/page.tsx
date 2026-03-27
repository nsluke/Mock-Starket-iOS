'use client';

import { useEffect, useState } from 'react';
import { apiClient } from '@/lib/api-client';
import { formatCurrency } from '@/lib/formatters';

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

const challengeIcons: Record<string, string> = {
  trade_count: '🔄',
  buy_stock: '🛒',
  sell_stock: '💰',
  profit_target: '📈',
  diversify: '🎯',
  volume_trader: '📊',
};

export default function ChallengesPage() {
  const [challenge, setChallenge] = useState<DailyChallenge | null>(null);
  const [progress, setProgress] = useState<UserProgress | null>(null);
  const [loading, setLoading] = useState(true);
  const [checking, setChecking] = useState(false);
  const [claiming, setClaiming] = useState(false);
  const [message, setMessage] = useState<string | null>(null);

  useEffect(() => {
    loadChallenge();
  }, []);

  async function loadChallenge() {
    setLoading(true);
    try {
      const data = await apiClient.getTodaysChallenge();
      setChallenge(data.challenge);
      setProgress(data.progress);
    } catch {
      setChallenge(null);
      setProgress(null);
    } finally {
      setLoading(false);
    }
  }

  async function handleCheck() {
    setChecking(true);
    setMessage(null);
    try {
      const result = await apiClient.checkChallenge();
      if (result.completed) {
        setMessage('Challenge completed! Claim your reward below.');
        await loadChallenge();
      } else {
        setMessage('Not yet complete. Keep trading!');
      }
    } catch {
      setMessage('Failed to check progress.');
    } finally {
      setChecking(false);
    }
  }

  async function handleClaim() {
    if (!challenge) return;
    setClaiming(true);
    setMessage(null);
    try {
      await apiClient.claimChallenge(challenge.id);
      setMessage('Reward claimed! Cash has been added to your portfolio.');
      await loadChallenge();
    } catch (err: any) {
      setMessage(err.message || 'Failed to claim reward.');
    } finally {
      setClaiming(false);
    }
  }

  if (loading) {
    return <div className="p-6 text-center text-[#8B949E]">Loading challenge...</div>;
  }

  return (
    <div className="p-6 max-w-2xl mx-auto space-y-6">
      <h1 className="text-2xl font-bold">Daily Challenge</h1>

      {!challenge ? (
        <div className="rounded-xl bg-[#161B22] border border-[#30363D] p-8 text-center">
          <p className="text-4xl mb-4">🎯</p>
          <p className="text-[#8B949E]">No challenge available today. Check back later!</p>
        </div>
      ) : (
        <div className="rounded-xl bg-[#161B22] border border-[#30363D] p-8 space-y-6">
          {/* Challenge header */}
          <div className="text-center space-y-3">
            <p className="text-5xl">{challengeIcons[challenge.challenge_type] || '🎯'}</p>
            <h2 className="text-xl font-bold">{challenge.description}</h2>
            <div className="inline-flex items-center gap-2 rounded-full bg-emerald-500/10 px-4 py-2 text-sm font-medium text-emerald-400">
              <span>💵</span>
              Reward: {formatCurrency(challenge.reward_cash)}
            </div>
          </div>

          {/* Status */}
          <div className="border-t border-[#30363D] pt-6">
            {progress?.claimed ? (
              <div className="text-center space-y-3">
                <div className="inline-flex items-center gap-2 rounded-full bg-emerald-500/10 px-6 py-3 text-emerald-400 font-semibold">
                  <span>✅</span> Reward Claimed!
                </div>
                <p className="text-sm text-[#8B949E]">
                  Come back tomorrow for a new challenge.
                </p>
              </div>
            ) : progress?.completed ? (
              <div className="text-center space-y-4">
                <p className="text-sm text-emerald-400 font-medium">Challenge completed!</p>
                <button
                  onClick={handleClaim}
                  disabled={claiming}
                  className="w-full py-3 rounded-lg bg-emerald-500 hover:bg-emerald-600 text-white text-sm font-semibold transition-colors disabled:opacity-50"
                >
                  {claiming ? 'Claiming...' : 'Claim Reward'}
                </button>
              </div>
            ) : (
              <div className="text-center space-y-4">
                <div className="flex items-center justify-center gap-2 text-sm text-[#8B949E]">
                  <span className="w-2 h-2 rounded-full bg-yellow-400 animate-pulse"></span>
                  In Progress
                </div>
                <button
                  onClick={handleCheck}
                  disabled={checking}
                  className="w-full py-3 rounded-lg bg-[#21262D] hover:bg-[#30363D] text-white text-sm font-semibold transition-colors disabled:opacity-50"
                >
                  {checking ? 'Checking...' : 'Check Progress'}
                </button>
              </div>
            )}

            {message && (
              <p className="text-center text-sm mt-4 text-[#8B949E]">{message}</p>
            )}
          </div>
        </div>
      )}
    </div>
  );
}

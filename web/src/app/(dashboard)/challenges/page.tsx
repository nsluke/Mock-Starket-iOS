'use client';

import { PageTransition } from '@/components/ui/PageTransition';
import { useTodaysChallenge, useCheckChallenge, useClaimChallenge } from '@/hooks/use-challenges';
import { formatCurrency } from '@/lib/formatters';

const challengeIcons: Record<string, string> = {
  trade_count: '\uD83D\uDD04',
  buy_stock: '\uD83D\uDED2',
  sell_stock: '\uD83D\uDCB0',
  profit_target: '\uD83D\uDCC8',
  diversify: '\uD83C\uDFAF',
  volume_trader: '\uD83D\uDCCA',
};

export default function ChallengesPage() {
  const { data, isLoading } = useTodaysChallenge();
  const checkChallenge = useCheckChallenge();
  const claimChallenge = useClaimChallenge();

  const challenge = data?.challenge;
  const progress = data?.progress;

  if (isLoading) {
    return <div className="p-6 text-center text-[#8B949E]">Loading challenge...</div>;
  }

  return (
    <PageTransition>
    <div className="p-6 max-w-2xl mx-auto space-y-6">
      <h1 className="text-2xl font-bold">Daily Challenge</h1>

      {!challenge ? (
        <div className="rounded-xl bg-[#161B22] border border-[#30363D] p-8 text-center">
          <p className="text-4xl mb-4">{'\uD83C\uDFAF'}</p>
          <p className="text-[#8B949E]">No challenge available today. Check back later!</p>
        </div>
      ) : (
        <div className="rounded-xl bg-[#161B22] border border-[#30363D] p-8 space-y-6">
          {/* Challenge header */}
          <div className="text-center space-y-3">
            <p className="text-5xl">{challengeIcons[challenge.challenge_type] || '\uD83C\uDFAF'}</p>
            <h2 className="text-xl font-bold">{challenge.description}</h2>
            <div className="inline-flex items-center gap-2 rounded-full bg-emerald-500/10 px-4 py-2 text-sm font-medium text-emerald-400">
              <span>{'\uD83D\uDCB5'}</span>
              Reward: {formatCurrency(challenge.reward_cash)}
            </div>
          </div>

          {/* Status */}
          <div className="border-t border-[#30363D] pt-6">
            {progress?.claimed ? (
              <div className="text-center space-y-3">
                <div className="inline-flex items-center gap-2 rounded-full bg-emerald-500/10 px-6 py-3 text-emerald-400 font-semibold">
                  <span>{'\u2705'}</span> Reward Claimed!
                </div>
                <p className="text-sm text-[#8B949E]">
                  Come back tomorrow for a new challenge.
                </p>
              </div>
            ) : progress?.completed ? (
              <div className="text-center space-y-4">
                <p className="text-sm text-emerald-400 font-medium">Challenge completed!</p>
                <button
                  onClick={() => claimChallenge.mutate(challenge.id)}
                  disabled={claimChallenge.isPending}
                  className="w-full py-3 rounded-lg bg-emerald-500 hover:bg-emerald-600 text-white text-sm font-semibold transition-colors disabled:opacity-50"
                >
                  {claimChallenge.isPending ? 'Claiming...' : 'Claim Reward'}
                </button>
              </div>
            ) : (
              <div className="text-center space-y-4">
                <div className="flex items-center justify-center gap-2 text-sm text-[#8B949E]">
                  <span className="w-2 h-2 rounded-full bg-yellow-400 animate-pulse"></span>
                  In Progress
                </div>
                <button
                  onClick={() => checkChallenge.mutate()}
                  disabled={checkChallenge.isPending}
                  className="w-full py-3 rounded-lg bg-[#21262D] hover:bg-[#30363D] text-white text-sm font-semibold transition-colors disabled:opacity-50"
                >
                  {checkChallenge.isPending ? 'Checking...' : 'Check Progress'}
                </button>
              </div>
            )}
          </div>
        </div>
      )}
    </div>
    </PageTransition>
  );
}

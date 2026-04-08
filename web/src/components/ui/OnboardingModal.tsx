'use client';

import { useState, useCallback } from 'react';
import { TrendingUp, DollarSign, BarChart3, Zap, Trophy, Target, X, ChevronRight, ChevronLeft } from 'lucide-react';

interface OnboardingModalProps {
  onClose: () => void;
}

const steps = [
  {
    icon: TrendingUp,
    iconColor: 'text-[#50E3C2]',
    iconBg: 'bg-[#50E3C2]/10',
    title: 'Real Stocks. Real Prices.',
    description: 'Mock Starket uses live market data from real exchanges. The prices you see for Apple, Tesla, Bitcoin, and 40+ other assets are the same prices Wall Street is trading right now.',
    details: [
      'Live data from NYSE, NASDAQ & crypto exchanges',
      'Real companies: AAPL, NVDA, TSLA, MSFT & more',
      'The perfect way to learn how markets actually work',
    ],
  },
  {
    icon: DollarSign,
    iconColor: 'text-emerald-400',
    iconBg: 'bg-emerald-400/10',
    title: 'Paper Trade with $100K',
    description: 'You get $100,000 in virtual cash to invest however you want. Build a portfolio, test strategies, and learn from your wins and losses — all without risking a single real dollar.',
    details: [
      'Stocks, ETFs, crypto — all with real prices',
      'Practice before you invest real money',
      'No sign-ups, no credit card, no catch',
    ],
  },
  {
    icon: BarChart3,
    iconColor: 'text-blue-400',
    iconBg: 'bg-blue-400/10',
    title: 'Learn Real Trading Skills',
    description: 'Use the same order types that professional traders use. Understand how market orders, limit orders, and stop losses work before you trade with real money.',
    details: [
      'Market orders execute at the current price',
      'Limit orders let you name your price',
      'Stop orders help you manage risk',
    ],
  },
  {
    icon: Trophy,
    iconColor: 'text-amber-500',
    iconBg: 'bg-amber-500/20',
    title: 'Compete & Build Confidence',
    description: 'Track your performance on the leaderboard, earn achievements as you learn, and take on daily challenges to sharpen your skills.',
    details: [
      'Daily challenges with cash rewards',
      '20 achievements from beginner to pro',
      'See how you stack up against other learners',
    ],
  },
  {
    icon: Target,
    iconColor: 'text-yellow-400',
    iconBg: 'bg-yellow-400/10',
    title: 'Your Market Toolkit',
    description: 'Set price alerts on real stocks, build a watchlist of companies you want to follow, and watch your portfolio grow with real market movements.',
    details: [
      'Price alerts for real market moves',
      'Watchlist to follow your favorite stocks',
      'Portfolio charts tracking real performance',
    ],
  },
];

export function OnboardingModal({ onClose }: OnboardingModalProps) {
  const [step, setStep] = useState(0);
  const current = steps[step];
  const Icon = current.icon;
  const isLast = step === steps.length - 1;

  const handleNext = useCallback(() => {
    if (isLast) {
      onClose();
    } else {
      setStep((s) => s + 1);
    }
  }, [isLast, onClose]);

  const handleBack = useCallback(() => {
    setStep((s) => s - 1);
  }, []);

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/60 px-4">
      <div className="relative w-full max-w-md rounded-2xl bg-[#161B22] border border-[#30363D] overflow-hidden">
        {/* Close button */}
        <button
          onClick={onClose}
          className="absolute top-4 right-4 text-[#6E7681] hover:text-white transition-colors"
        >
          <X className="w-5 h-5" />
        </button>

        <div className="p-8 pt-10">
          {/* Icon */}
          <div className={`w-14 h-14 rounded-xl ${current.iconBg} flex items-center justify-center mb-6`}>
            <Icon className={`w-7 h-7 ${current.iconColor}`} />
          </div>

          {/* Content */}
          <h2 className="text-2xl font-bold text-white mb-3">{current.title}</h2>
          <p className="text-[#8B949E] leading-relaxed mb-6">{current.description}</p>

          {/* Details */}
          <ul className="space-y-3 mb-8">
            {current.details.map((detail, i) => (
              <li key={i} className="flex items-start gap-3 text-sm">
                <div className="mt-1.5 w-1.5 h-1.5 rounded-full bg-[#50E3C2] shrink-0" />
                <span className="text-[#E6EDF3]">{detail}</span>
              </li>
            ))}
          </ul>

          {/* Progress dots */}
          <div className="flex items-center justify-center gap-2 mb-6">
            {steps.map((_, i) => (
              <div
                key={i}
                className={`h-1.5 rounded-full transition-all ${
                  i === step ? 'w-6 bg-[#50E3C2]' : 'w-1.5 bg-[#30363D]'
                }`}
              />
            ))}
          </div>

          {/* Navigation */}
          <div className="flex gap-3">
            {step > 0 && (
              <button
                onClick={handleBack}
                className="flex items-center justify-center gap-2 rounded-xl bg-[#21262D] px-5 py-3 text-sm font-semibold text-white transition-colors hover:bg-[#30363D]"
              >
                <ChevronLeft className="w-4 h-4" />
                Back
              </button>
            )}
            <button
              onClick={handleNext}
              className="flex-1 flex items-center justify-center gap-2 rounded-xl bg-[#50E3C2] px-5 py-3 text-sm font-semibold text-black transition-transform hover:scale-[1.02] active:scale-[0.98]"
            >
              {isLast ? 'Start Trading' : 'Next'}
              {!isLast && <ChevronRight className="w-4 h-4" />}
              {isLast && <TrendingUp className="w-4 h-4" />}
            </button>
          </div>

          {/* Skip link */}
          {!isLast && (
            <button
              onClick={onClose}
              className="w-full mt-3 text-xs text-[#6E7681] hover:text-[#8B949E] transition-colors"
            >
              Skip intro
            </button>
          )}
        </div>
      </div>
    </div>
  );
}

'use client';

import { useState, useEffect, useRef } from 'react';
import { Clock } from 'lucide-react';
import { wsClient } from '@/lib/websocket-client';

const TICK_INTERVAL = 30; // seconds

export function TickCountdown() {
  const [secondsLeft, setSecondsLeft] = useState(TICK_INTERVAL);
  const timerRef = useRef<ReturnType<typeof setInterval> | null>(null);

  useEffect(() => {
    function resetTimer() {
      setSecondsLeft(TICK_INTERVAL);
    }

    // Reset countdown whenever we receive a price update
    wsClient.on('price_batch', resetTimer);

    // Tick down every second
    timerRef.current = setInterval(() => {
      setSecondsLeft((s) => (s > 0 ? s - 1 : 0));
    }, 1000);

    return () => {
      wsClient.off('price_batch', resetTimer);
      if (timerRef.current) clearInterval(timerRef.current);
    };
  }, []);

  const pct = (secondsLeft / TICK_INTERVAL) * 100;

  return (
    <div className="flex items-center gap-2 text-xs text-[#8B949E]">
      <Clock className="w-3 h-3" />
      <div className="flex items-center gap-2 min-w-0">
        <span className="whitespace-nowrap">
          Next tick {secondsLeft}s
        </span>
        <div className="w-16 h-1.5 rounded-full bg-[#21262D] overflow-hidden">
          <div
            className="h-full rounded-full bg-[#50E3C2] transition-all duration-1000 ease-linear"
            style={{ width: `${pct}%` }}
          />
        </div>
      </div>
    </div>
  );
}

'use client';

import { useEffect, useState } from 'react';
import { wsClient, type ConnectionStatus as Status } from '@/lib/websocket-client';
import { cn } from '@/lib/cn';

const statusConfig: Record<Status, { color: string; label: string }> = {
  connected: { color: 'bg-emerald-400', label: 'Connected' },
  reconnecting: { color: 'bg-yellow-400 animate-pulse', label: 'Reconnecting...' },
  disconnected: { color: 'bg-red-400', label: 'Disconnected' },
};

export function ConnectionStatus() {
  const [status, setStatus] = useState<Status>(wsClient.status);

  useEffect(() => {
    return wsClient.onStatusChange(setStatus);
  }, []);

  const config = statusConfig[status];

  return (
    <div className="flex items-center gap-2 text-xs text-[#6E7681]">
      <span className={cn('w-2 h-2 rounded-full', config.color)} />
      {config.label}
    </div>
  );
}

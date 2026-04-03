import { Skeleton } from './Skeleton';

export function CardSkeleton() {
  return (
    <div className="rounded-xl bg-[#161B22] border border-[#30363D] p-4 space-y-2">
      <Skeleton className="h-3 w-16" />
      <Skeleton className="h-6 w-24" />
      <Skeleton className="h-3 w-20" />
    </div>
  );
}

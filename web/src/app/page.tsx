'use client';

import { useEffect } from 'react';
import { useRouter } from 'next/navigation';

export default function Home() {
  const router = useRouter();

  useEffect(() => {
    // Check if user is authenticated
    const token = localStorage.getItem('mockstarket_token');
    if (token) {
      router.replace('/market');
    } else {
      router.replace('/login');
    }
  }, [router]);

  return (
    <div className="flex h-screen items-center justify-center bg-[#0D1117]">
      <div className="text-center">
        <div className="text-5xl mb-4">📈</div>
        <div className="text-lg text-[#8B949E]">Loading...</div>
      </div>
    </div>
  );
}

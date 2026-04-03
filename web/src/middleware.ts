import { NextResponse } from 'next/server';
import type { NextRequest } from 'next/server';

const publicPaths = ['/login', '/'];

export function middleware(request: NextRequest) {
  const { pathname } = request.nextUrl;
  const token = request.cookies.get('mockstarket_token')?.value;

  // Allow public paths
  if (publicPaths.includes(pathname)) {
    // If logged in and on login page, redirect to market
    if (token && pathname === '/login') {
      return NextResponse.redirect(new URL('/market', request.url));
    }
    return NextResponse.next();
  }

  // Protect dashboard routes
  if (!token) {
    return NextResponse.redirect(new URL('/login', request.url));
  }

  return NextResponse.next();
}

export const config = {
  matcher: [
    '/((?!_next/static|_next/image|favicon.ico|api).*)',
  ],
};

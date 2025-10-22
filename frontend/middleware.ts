import { NextResponse } from 'next/server'
import type { NextRequest } from 'next/server'

const protectedRoutes = [
  '/dashboard',
  '/badges',
  '/chat-bot',
  '/electricity-tracker',
  '/leaderboard',
  '/missions',
  '/profile',
  '/shop',
  '/vehicle-tracker',
]

const authRoutes = ['/login', '/register']

export function middleware(req: NextRequest) {
  const token = req.cookies.get('token')?.value
  const { pathname } = req.nextUrl

  // Jika user belum login tapi akses halaman protected
  if (!token && protectedRoutes.some(route => pathname.startsWith(route))) {
    return NextResponse.redirect(new URL('/login', req.url))
  }

  // Jika user sudah login tapi akses login/register
  if (token && authRoutes.some(route => pathname.startsWith(route))) {
    return NextResponse.redirect(new URL('/dashboard', req.url))
  }

  return NextResponse.next()
}

export const config = {
  matcher: [
    '/dashboard/:path*',
    '/badges/:path*',
    '/chat-bot/:path*',
    '/electricity-tracker/:path*',
    '/leaderboard/:path*',
    '/missions/:path*',
    '/profile/:path*',
    '/shop/:path*',
    '/vehicle-tracker/:path*',
    '/login',
    '/register',
  ],
}

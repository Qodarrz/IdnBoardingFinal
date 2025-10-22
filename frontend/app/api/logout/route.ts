import { NextResponse } from 'next/server'

export async function POST() {
  const res = NextResponse.json({ success: true })

  // hapus cookie token
  res.cookies.set('token', '', {
    httpOnly: true,
    secure: process.env.NODE_ENV === 'production',
    path: '/',
    expires: new Date(0),
  })

  return res
}

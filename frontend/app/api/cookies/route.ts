import { NextResponse } from 'next/server'

export async function POST(req: Request) {
  const { token } = await req.json()

  if (token) {
    const res = NextResponse.json({ success: true })

    res.cookies.set('token', token, {
      httpOnly: true, 
      secure: process.env.NODE_ENV === 'production',
      path: '/',
      maxAge: 60 * 60 * 24,
    })

    return res
  }

  return NextResponse.json({ success: false, message: 'Invalid credentials' }, { status: 401 })
}

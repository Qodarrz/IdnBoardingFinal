"use client"

import { useRouter } from "next/navigation"

export default function LogoutAction() {
  const router = useRouter()

  const handleLogout = async () => {
    localStorage.removeItem('islogin')
    localStorage.removeItem('authtoken')
    await fetch("/api/logout", { method: "POST" })
    router.push("/login")
  }

  return handleLogout
}

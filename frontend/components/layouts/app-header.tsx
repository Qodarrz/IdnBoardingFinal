// di file AppHeader.tsx
"use client"

import Image from "next/image"
import { ChevronDown, LogOut, User } from "lucide-react"
import { IconBellFilled, IconShoppingBag } from "@tabler/icons-react"
import { Button } from "@/components/ui/button"
import { Separator } from "@/components/ui/separator"
import { SidebarTrigger } from "@/components/ui/sidebar"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import { useAuthMe } from "@/helpers/AuthMe"
import LogoutAction from "@/helpers/LogoutAction"
import Link from "next/link"

// ⬇️ import helper baru
import { useNotifications, timeAgo } from "@/helpers/fetchNotifications"

export function AppHeader() {
  const { data: dataUser, loading, error } = useAuthMe()
  const logout = LogoutAction()

  // ⬇️ panggil hook notifikasi (contoh polling tiap 30 detik)
  const {
    data: notif,
    loading: notifLoading,
    error: notifError,
    unreadCount,
    refetch,
  } = useNotifications({
    refreshInterval: 30000,
    getToken: () => localStorage.getItem("authtoken"),
  })

  if (loading) {
    console.log("loadingg..")
    // return <header className="p-4">Loading...</header>
  }

  if (error) {
    console.log("error get me..", error)
    // return <header className="p-4 text-red-500">Gagal memuat data</header>
  }

  return (
    <header className="flex py-4 shrink-0 items-center gap-2 border-b transition-[width,height] ease-linear group-has-data-[collapsible=icon]/sidebar-wrapper:h-(--header-height)">
      <div className="flex w-full items-center gap-1 px-4 lg:gap-2 lg:px-6">
        <SidebarTrigger className="-ml-1" />
        <Separator
          orientation="vertical"
          className="mx-2 data-[orientation=vertical]:h-4"
        />
        <h1 className="text-base font-medium">Dashboard</h1>

        <div className="ml-auto flex items-center gap-4">
          {/* 🔔 Notifikasi (dinamis) */}
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="ghost" size="icon" className="relative" aria-label="Notifikasi" onClick={refetch}>
                <IconBellFilled />
                {unreadCount > 0 && (
                  <>
                    <span className="absolute -top-0.5 -right-0.5 h-2 w-2 rounded-full bg-red-500" />
                    <span className="absolute -bottom-1 -right-1 min-w-5 h-5 px-1 rounded-full bg-red-500 text-white text-[10px] leading-5 text-center">
                      {unreadCount > 99 ? "99+" : unreadCount}
                    </span>
                  </>
                )}
              </Button>
            </DropdownMenuTrigger>

            <DropdownMenuContent align="end" className="w-80 p-0">
              <div className="px-3 py-2">
                <DropdownMenuLabel>Notifikasi</DropdownMenuLabel>
                <p className="text-xs text-muted-foreground">
                  {notifLoading ? "Memuat..." : notifError ? "Gagal memuat notifikasi" : `${(notif?.length || 0)} total`}
                </p>
              </div>
              <DropdownMenuSeparator />

              {/* List notifikasi */}
              <div className="max-h-80 overflow-auto">
                {!notifLoading && !notifError && (notif?.length ?? 0) === 0 && (
                  <div className="px-3 py-6 text-sm text-muted-foreground text-center">
                    Belum ada notifikasi.
                  </div>
                )}

                {notif?.map((n) => (
                  <DropdownMenuItem key={n.id} className="py-3 px-3 flex items-start gap-2">
                    <div className={`mt-1 h-2 w-2 rounded-full ${n.is_read ? "bg-muted" : "bg-blue-500"}`} />
                    <div className="flex-1">
                      <div className="text-sm font-medium">{n.title || "Notifikasi"}</div>
                      <div className="text-xs text-muted-foreground">{n.message}</div>
                      <div className="mt-1 text-[10px] text-muted-foreground">{timeAgo(n.created_at)}</div>
                    </div>
                  </DropdownMenuItem>
                ))}
              </div>

              <DropdownMenuSeparator />
              <div className="px-3 py-2 flex items-center justify-between">
                <Button variant="ghost" size="sm" onClick={refetch}>
                  Muat ulang
                </Button>
                <Button variant="ghost" size="sm">
                  Lihat semua
                </Button>
              </div>
            </DropdownMenuContent>
          </DropdownMenu>

          {/* 👤 Profile */}
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="ghost" className="flex items-center gap-2">
                <Image
                  src="/images/profile.png"
                  alt="Avatar"
                  width={32}
                  height={32}
                  className="rounded-full"
                />
                <div className="hidden sm:flex flex-col items-start">
                  <span className="text-sm font-medium">{dataUser?.data?.username ?? 'Guest'}</span>
                  <span className="text-xs text-muted-foreground">{dataUser?.data?.role ?? 'User'}</span>
                </div>
                <ChevronDown className="h-4 w-4" />
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end" className="w-56">
              <DropdownMenuLabel>Akun Saya</DropdownMenuLabel>
              <DropdownMenuSeparator />
              <DropdownMenuItem>
                <Link href="/profile" className="flex gap-1">
                  <User className="mr-2 h-4 w-4" />
                  <span>Profil</span>
                </Link>
              </DropdownMenuItem>
              <DropdownMenuItem>
                <Link href="/order" className="flex gap-1">
                  <IconShoppingBag className="mr-2 h-4 w-4" />
                  <span>Order</span>
                </Link>
              </DropdownMenuItem>
              <DropdownMenuSeparator />
              <DropdownMenuItem onClick={() => logout()}>
                <LogOut className="mr-2 h-4 w-4" />
                <span>Keluar</span>
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        </div>
      </div>
    </header>
  )
}

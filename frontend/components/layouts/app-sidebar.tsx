"use client"

import * as React from "react"
import {
  IconShoppingCartFilled,
  IconLayoutDashboardFilled,
  IconCarFilled,
  IconDeviceMobileFilled,
  IconSparkles,
  IconTargetArrow,
  IconTrophyFilled,
  IconAwardFilled,
  IconLogout
} from "@tabler/icons-react"

import { NavMain } from "@/components/shared/nav-main"
import { NavSecondary } from "@/components/shared/nav-secondary"
import {
  Sidebar,
  SidebarContent,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from "@/components/ui/sidebar"
import Link from "next/link"
import Image from "next/image"
import { NavTracker } from "../shared/nav-tracker"
import { NavGamifikasi } from "../shared/nav-gamifikasi"

const data = {
  user: {
    name: "shadcn",
    email: "m@example.com",
    avatar: "/avatars/shadcn.jpg",
  },
  navMain: [
    {
      title: "Dashboard",
      url: "/dashboard",
      icon: IconLayoutDashboardFilled,
    },
    {
      title: "ChatBot AI",
      url: "/chat-bot",
      icon: IconSparkles,
    },
    {
      title: "Toko Penukaran",
      url: "/shop",
      icon: IconShoppingCartFilled,
    },
  ],
  navSecondary: [
    {
      title: "Logout",
      url: "#",
      icon: IconLogout,
    },
  ],
  navTracker: [
    {
      title: "Karbon Kendaraan",
      url: "/vehicle-tracker",
      icon: IconCarFilled,
    },
    {
      title: "Karbon Alat Elektronik",
      url: "/electricity-tracker",
      icon: IconDeviceMobileFilled,
    },
  ],
  navGamifikasi: [
    {
      title: "Misi",
      url: "/missions",
      icon: IconTargetArrow,
    },
    {
      title: "Lencana Saya",
      url: "/badges",
      icon: IconAwardFilled,
    },
    {
      title: "Papan Peringkat",
      url: "/leaderboard",
      icon: IconTrophyFilled,
    },
  ],
}

export function AppSidebar({ ...props }: React.ComponentProps<typeof Sidebar>) {
  return (
    <Sidebar collapsible="offcanvas" {...props}>
      <SidebarHeader>
        <SidebarMenu>
          <SidebarMenuItem>
            <SidebarMenuButton
              asChild
              className="data-[slot=sidebar-menu-button]:!p-1.5"
            >
              <Link href="/">
                <Image src="/icons/main-logo.svg" alt="Logo" width={150} height={150} />
              </Link>
            </SidebarMenuButton>
          </SidebarMenuItem>
        </SidebarMenu>
      </SidebarHeader>
      <SidebarContent>
        <NavMain items={data.navMain} />
        <NavTracker items={data.navTracker} />
        <NavGamifikasi items={data.navGamifikasi} />
        <NavSecondary items={data.navSecondary} className="mt-auto" />
      </SidebarContent>
    </Sidebar>
  )
}

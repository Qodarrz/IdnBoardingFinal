"use client"

import { ReactNode } from "react"
import { AppHeader } from "@/components/layouts/app-header"
import { SidebarInset, SidebarProvider } from "@/components/ui/sidebar"
import { AppSidebar } from "@/components/layouts/app-sidebar"

export default function DashboardLayout({ children }: { children: ReactNode }) {
    return (
        <SidebarProvider
            style={
                {
                    "--sidebar-width": "calc(var(--spacing) * 72)",
                    "--header-height": "calc(var(--spacing) * 12)",
                } as React.CSSProperties
            }
        >   
            <AppSidebar variant="inset" />
            <SidebarInset>
                <AppHeader />
                <div className="flex flex-1 flex-col">{children}</div>
            </SidebarInset>
        </SidebarProvider>
    )
}

"use client"

import { useEffect, useMemo, useRef, useState } from "react"

export type NotificationItem = {
    id: number
    user_id: number
    title: string
    message: string
    type: string
    is_read: boolean
    created_at: string
}

type ApiResponse = {
    status: boolean
    message: string
    data: NotificationItem[]
}

type Options = {
    refreshInterval?: number
    endpoint?: string
    getToken?: () => string | null | undefined | Promise<string | null | undefined>
    extraHeaders?: Record<string, string>
}

export function useNotifications(options?: Options) {
    const {
        refreshInterval = 0,
        endpoint = `${process.env.NEXT_PUBLIC_API_URL}/api/custom/my-data`,
        getToken,
        extraHeaders = {},
    } = options || {}

    const [data, setData] = useState<NotificationItem[] | null>(null)
    const [loading, setLoading] = useState<boolean>(true)
    const [error, setError] = useState<string | null>(null)

    const abortRef = useRef<AbortController | null>(null)
    const timerRef = useRef<number | null>(null)

    async function fetchNotifications(signal?: AbortSignal) {
        try {
            setError(null)
            if (!data) setLoading(true)

            // ambil token (jika disediakan)
            const tokenValue = getToken ? await getToken() : undefined
            if (getToken && !tokenValue) {
                throw new Error("Token tidak tersedia. Pastikan sudah login.")
            }

            const res = await fetch(endpoint, {
                method: "GET",
                headers: {
                    "Accept": "application/json",
                    ...(tokenValue ? { Authorization: `Bearer ${tokenValue}` } : {}),
                    ...extraHeaders,
                },
                // NOTE: biasanya tak perlu credentials jika pakai Bearer token
                signal,
            })

            if (!res.ok) {
                // opsional: khusus 401
                if (res.status === 401) {
                    throw new Error("Unauthorized (401). Token invalid/kedaluwarsa.")
                }
                const text = await res.text().catch(() => "")
                throw new Error(text || `HTTP ${res.status}`)
            }

            const json: ApiResponse = await res.json()
            if (!json.status) throw new Error(json.message || "Gagal mengambil notifikasi")

            setData(Array.isArray(json.data) ? json.data : [])
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        } catch (e: any) {
            if (e?.name !== "AbortError") {
                setError(e?.message || "Terjadi kesalahan")
            }
        } finally {
            setLoading(false)
        }
    }

    useEffect(() => {
        abortRef.current?.abort()
        const controller = new AbortController()
        abortRef.current = controller

        fetchNotifications(controller.signal)

        if (refreshInterval > 0) {
            timerRef.current = window.setInterval(() => {
                fetchNotifications(controller.signal)
            }, refreshInterval) as unknown as number
        }

        return () => {
            controller.abort()
            if (timerRef.current) clearInterval(timerRef.current)
        }
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [endpoint, refreshInterval]) // token diambil dinamis via getToken()

    const unreadCount = useMemo(
        () => (data || []).reduce((acc, n) => acc + (n.is_read ? 0 : 1), 0),
        [data]
    )

    return {
        data,
        loading,
        error,
        unreadCount,
        refetch: () => fetchNotifications(abortRef.current?.signal),
    }
}

// util kecil untuk tampilan waktu relatif (ID)
export function timeAgo(iso: string) {
    const now = new Date().getTime()
    const ts = new Date(iso).getTime()
    const diff = Math.max(0, now - ts)

    const s = Math.floor(diff / 1000)
    if (s < 60) return `${s} detik lalu`
    const m = Math.floor(s / 60)
    if (m < 60) return `${m} menit lalu`
    const h = Math.floor(m / 60)
    if (h < 24) return `${h} jam lalu`
    const d = Math.floor(h / 24)
    if (d < 30) return `${d} hari lalu`
    const mo = Math.floor(d / 30)
    if (mo < 12) return `${mo} bln lalu`
    const y = Math.floor(mo / 12)
    return `${y} thn lalu`
}

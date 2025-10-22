"use client"

import { useEffect, useState } from "react"
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card"
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table"
import Image from "next/image"
import { IconCheck } from "@tabler/icons-react"
import { GetLeaderboard } from "@/helpers/GetLeaderboard" // Import the helper function
import { fetchProfile } from "@/helpers/fetchProfile"

type LeaderboardUser = {
    rank: number
    user: {
        id: number
        username: string
        full_name: string
        avatar_url: string
    }
    total_points: number
    completed_missions: number
    score: number
    carbon_reduction_g?: number
}

export default function LeaderboardPage() {
    const [leaderboard, setLeaderboard] = useState<LeaderboardUser[]>([])
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    const [profile, setProfile] = useState<any>(null)
    const [loading, setLoading] = useState(true)

    useEffect(() => {
        const fetchData = async () => {
            const token = localStorage.getItem("authtoken")
            if (!token) {
                console.error("Token tidak ditemukan. Silakan login ulang.")
                return
            }

            try {
                // Fetch profile data
                const profileData = await fetchProfile(token)
                setProfile(profileData)

                // Fetch leaderboard data
                const leaderboardData = await GetLeaderboard(token)
                if (leaderboardData && Array.isArray(leaderboardData)) {
                    setLeaderboard(leaderboardData) // Update leaderboard data
                }
            } catch (err) {
                console.error("Error fetching data:", err)
            } finally {
                setLoading(false)
            }
        }

        fetchData()
    }, [])

    if (loading) return <p className="p-6">Loading...</p>

    const topThree = leaderboard.slice(0, 3)

    return (
        <div className="p-6 space-y-6">
            {/* Statistik di atas */}
            <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
                <Card className="flex justify-center">
                    <CardContent>
                        <div className="flex w-full mb-3 justify-between">
                            <p className="text-4xl font-semibold">{profile?.missions.length} Misi</p>
                            <div className="bg-green-700 p-2 rounded-full aspect-square text-white">
                                <IconCheck size={32} />
                            </div>
                        </div>
                        <CardTitle className="font-normal">Misi Yang Terselesaikan</CardTitle>
                    </CardContent>
                </Card>
                <Card className="flex justify-center">
                    <CardContent>
                        <div className="flex w-full mb-3 justify-between">
                            <p className="text-4xl font-semibold">{profile?.user.total_points}</p>
                            <div className="rounded-full aspect-square">
                                <Image src="/icons/green-point.svg" alt="point" width={52} height={52} />
                            </div>
                        </div>
                        <CardTitle className="font-normal">Poin Yang Didapatkan</CardTitle>
                    </CardContent>
                </Card>
                <Card className="flex col-span-2 justify-center">
                    <CardContent className="flex justify-between items-center">
                        <div className="space-y-6">
                            <CardTitle className="text-4xl">Papan Peringkat</CardTitle>
                            <p>Peringkat Selama Agustus 2025</p>
                        </div>
                        <Image src="/icons/trophy.svg" alt="point" width={80} height={80} />
                    </CardContent>
                </Card>
            </div>

            {/* Top 3 Leaderboard */}
            <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                {topThree.map((user) => (
                    <Card
                        key={user.user.id}
                        className={user.rank === 1 ? "border-yellow-400 border-2" : ""}
                    >
                        <CardHeader className="flex flex-row items-center gap-4">
                            <Image
                                src={user.user.avatar_url || "/images/profile.png"}
                                alt={user.user.username}
                                width={48}
                                height={48}
                                className="rounded-full"
                            />
                            <CardTitle className="text-2xl">
                                {user.user.full_name || user.user.username}
                            </CardTitle>
                        </CardHeader>
                        <CardContent className="flex justify-between">
                            <div className="flex flex-col">
                                <p>Poin Didapatkan</p>
                                <p className="text-2xl font-bold">{user.score}</p>
                            </div>
                            <div className="flex flex-col">
                                <p>Misi Terselesaikan</p>
                                <p className="text-2xl font-bold">{user.completed_missions}</p>
                            </div>
                        </CardContent>
                    </Card>
                ))}
            </div>

            {/* Tabel Ranking Lainnya */}
            <Card>
                <CardHeader>
                    <CardTitle>Daftar Peringkat Lainnya</CardTitle>
                </CardHeader>
                <CardContent>
                    <Table>
                        <TableHeader>
                            <TableRow>
                                <TableHead>Rank</TableHead>
                                <TableHead>Nama</TableHead>
                                <TableHead>Poin Yang Didapatkan</TableHead>
                                <TableHead>Misi Yang Terselesaikan</TableHead>
                                <TableHead>Total Poin Peringkat</TableHead>
                            </TableRow>
                        </TableHeader>
                        <TableBody>
                            {leaderboard.map((user) => (
                                <TableRow key={user.user.id}>
                                    <TableCell>{user.rank}</TableCell>
                                    <TableCell className="flex items-center gap-2">
                                        <Image
                                            src={user.user.avatar_url || "/images/profile.png"}
                                            alt={user.user.username}
                                            width={32}
                                            height={32}
                                            className="rounded-full"
                                        />
                                        {user.user.full_name || user.user.username}
                                    </TableCell>
                                    <TableCell>{user.score}</TableCell>
                                    <TableCell>{user.completed_missions}</TableCell>
                                    <TableCell className="text-primary font-semibold">
                                        {user.total_points}
                                    </TableCell>
                                </TableRow>
                            ))}
                        </TableBody>
                    </Table>
                </CardContent>
            </Card>
        </div>
    )
}

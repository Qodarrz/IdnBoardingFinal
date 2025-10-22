"use client";

import { useEffect, useState } from "react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Progress } from "@/components/ui/progress";
import { Badge } from "@/components/ui/badge";
import Image from "next/image";
import { GetMissions } from "@/helpers/GetMissions";

interface Mission {
    mission_id: number;
    title: string;
    description: string;
    target_value: number;
    progress: number;
    is_completed: boolean;
    points_reward: number;
}

export default function MissionPage() {
    const [missions, setMissions] = useState<Mission[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        const fetchMissionsData = async () => {
            try {
                const token = localStorage.getItem("authtoken");
                if (!token) {
                    setError("Token tidak ditemukan. Silakan login ulang.");
                    setLoading(false);
                    return;
                }

                const json = await GetMissions(token); // Use the helper function here
                if (json.status) {
                    setMissions(json.data);
                } else {
                    setError(json.message || "Gagal mengambil data");
                }
            } catch (error: unknown) {
                if (error instanceof Error) {
                    console.error("Error fetching missions:", error.message);
                    setError(error.message);
                } else {
                    console.error("Unknown error:", error);
                    setError("Terjadi kesalahan tak dikenal");
                }
            } finally {
                setLoading(false);
            }
        };

        fetchMissionsData();
    }, []);

    return (
        <div className="p-6 space-y-6">
            <div className="grid lg:grid-cols-6 gap-6">
                {/* Header Section */}
                <div className="flex flex-col lg:col-span-4 md:flex-row items-center justify-between bg-card p-6 rounded-2xl shadow-sm">
                    <div>
                        <h1 className="text-2xl font-semibold">Selesaikan Misi Rutin Anda</h1>
                        <p className="text-muted-foreground max-w-sm">
                            Anda telah berkontribusi mengurangi emisi karbon sebesar{" "}
                            <span className="font-medium">34,54 kg CO₂e</span> bulan ini!
                        </p>
                    </div>
                    <Image src="/icons/alarm.svg" alt="Alarm Icon" width={100} height={100} />
                </div>

                {/* Summary Section */}
                <Card className="lg:col-span-2">
                    <CardHeader>
                        <CardTitle>Misi Harian Selesai</CardTitle>
                    </CardHeader>
                    <CardContent>
                        <p className="text-3xl font-bold">
                            {missions.filter((m) => m.is_completed).length}{" "}
                            <span className="text-base font-normal">Selesai</span>
                        </p>
                        <p className="text-sm text-muted-foreground">Terselesaikan</p>
                    </CardContent>
                </Card>
            </div>

            {/* State Handler */}
            {loading && <p className="text-center text-muted-foreground">Loading misi...</p>}
            {error && <p className="text-center text-red-500">{error}</p>}

            {/* Missions Grid */}
            {!loading && !error && (
                <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
                    {missions.map((mission) => (
                        <Card key={mission.mission_id} className="flex flex-col">
                            <CardHeader className="flex flex-row items-center justify-between">
                                <Badge variant="secondary">{mission.is_completed ? "Selesai" : "Belum"}</Badge>
                                <Badge variant="secondary">
                                    <Image
                                        src="/icons/green-point.svg"
                                        alt="Poin Icon"
                                        width={12}
                                        height={12}
                                    />
                                    <span className="text-xs font-medium">{mission.points_reward} Poin</span>
                                </Badge>
                            </CardHeader>
                            <CardContent className="space-y-2">
                                <h3 className="font-semibold">{mission.title}</h3>
                                <p className="text-sm text-muted-foreground">{mission.description}</p>
                                <Progress value={(mission.progress / mission.target_value) * 100} />
                                <p className="text-xs text-muted-foreground">
                                    {mission.progress} / {mission.target_value}
                                </p>
                            </CardContent>
                        </Card>
                    ))}
                </div>
            )}
        </div>
    );
}

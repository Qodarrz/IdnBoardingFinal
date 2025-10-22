"use client"

import { useEffect, useState } from "react"
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card"
import { Dialog, DialogContent, DialogFooter, DialogHeader, DialogTitle } from "@/components/ui/dialog"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { IconPencil, IconTrendingUp } from "@tabler/icons-react"
import Image from "next/image"
import { fetchProfile } from "@/helpers/fetchProfile"
import { patchProfile } from "@/helpers/editProfile"

// Type definitions for the data structure
interface Vehicle {
    id: number
    vehicle_type: string
    fuel_type: string
    name: string
    total_carbon_emission_g: number
}

interface Electronic {
    id: number
    device_name: string
    device_type: string
    power_watts: number
    total_carbon_emission_g: number
}

interface BadgeType {
    id: number
    name: string
    image_url: string
    description: string
    redeemed_at: string
}

interface Mission {
    id: number
    title: string
    description: string
    mission_type: string
    points_reward: number
    target_value: number
    progress: number
    completed_at: string
    created_at: string
}

interface ProfileData {
    user: {
        id: number
        username: string
        email: string
        full_name: string
        avatar_url: string
        total_points: number
    }
    vehicles: Vehicle[]
    electronics: Electronic[]
    missions: Mission[]
    badges: BadgeType[]
}

export default function Profile() {
    const [open, setOpen] = useState(false)
    const [profile, setProfile] = useState<ProfileData | null>(null)
    const [loading, setLoading] = useState(true)
    const [name, setName] = useState("");
    const [email, setEmail] = useState("");
    const [username, setUsername] = useState("");
    const [gender, setGender] = useState<"male" | "female" | "other" | "">("");
    const [birthdate, setBirthdate] = useState("");      // "YYYY-MM-DD"
    const [avatarFile, setAvatarFile] = useState<File | null>(null);
    const [saving, setSaving] = useState(false);

    useEffect(() => {
        const fetchData = async () => {
            const token = localStorage.getItem("authtoken")
            if (!token) {
                console.error("Token tidak ditemukan")
                return
            }

            const profileData = await fetchProfile(token)
            setProfile(profileData)
            setLoading(false)
        }

        fetchData()
    }, [])

    useEffect(() => {
        const init = async () => {
            const token = localStorage.getItem("authtoken");
            if (!token) return;
            const profileData = await fetchProfile(token);
            setProfile(profileData);
            // seed form
            setName(profileData?.user?.full_name ?? "");
            setEmail(profileData?.user?.email ?? "");
            setUsername(profileData?.user?.username ?? "");
            setGender((profileData as any)?.user?.gender ?? ""); // kalau ada
            const bd = (profileData as any)?.user?.birthdate || (profileData as any)?.profile?.birthdate;
            if (bd) setBirthdate(String(bd).slice(0, 10)); // ke YYYY-MM-DD
            setLoading(false);
        };
        init();
    }, []);

    if (loading) return <p className="p-6">Loading...</p>

    const vehicles = profile?.vehicles || []
    const electronics = profile?.electronics || []

    // Calculate total carbon emissions in kilograms
    const totalCarbonEmission = (
        (vehicles.reduce((acc: number, vehicle: Vehicle) => acc + (vehicle.total_carbon_emission_g || 0), 0) +
            electronics.reduce((acc: number, device: Electronic) => acc + (device.total_carbon_emission_g || 0), 0)) / 1000
    ).toFixed(2)

    async function onSave() {
        try {
            setSaving(true);
            const token = localStorage.getItem("authtoken");
            if (!token) throw new Error("Token tidak ditemukan");

            const fd = new FormData();
            fd.append("full_name", name);
            fd.append("username", username);
            if (gender) fd.append("gender", gender);
            if (birthdate) fd.append("birthdate", birthdate);
            if (avatarFile) fd.append("avatar", avatarFile);

            const resp = await patchProfile(token, fd);
            // ...update state dari resp
            setOpen(false);
        } catch (e: any) {
            console.error(e);
            alert(e.message || "Failed to fetch");
        } finally {
            setSaving(false);
        }
    }


    return (
        <div className="p-6 space-y-6">
            <Card className="flex flex-col md:flex-row items-center justify-between gap-6 p-6 rounded-2xl shadow-sm">
                <div className="flex items-center gap-4">
                    <Avatar className="h-16 w-16">
                        <AvatarImage src={profile?.user?.avatar_url || "/images/profile.png"} alt="Avatar" />
                        <AvatarFallback>{profile?.user?.full_name.slice(0, 2).toUpperCase() || "NA"}</AvatarFallback>
                    </Avatar>
                    <div>
                        <h2 className="text-lg font-semibold">{profile?.user?.full_name || profile?.user?.username}</h2>
                        <p className="text-sm text-muted-foreground">{profile?.user?.email}</p>
                        <div className="flex gap-2 mt-2">
                            <span className="px-3 py-1 text-xs rounded-full border">Peringkat 22</span>
                            <span className="px-3 py-1 text-xs rounded-full border">Streak 12 Hari</span>
                        </div>
                    </div>
                </div>

                <Button variant="outline" size="sm" onClick={() => setOpen(true)}>
                    Edit data <IconPencil className="ml-2 h-4 w-4" />
                </Button>
            </Card>

            {/* Modal Edit */}
            <Dialog open={open} onOpenChange={setOpen}>
                <DialogContent>
                    <DialogHeader>
                        <DialogTitle>Edit Data Profil</DialogTitle>
                    </DialogHeader>

                    <div className="space-y-4">
                        <div>
                            <Label htmlFor="name">Nama</Label>
                            <Input id="name" value={name} onChange={(e) => setName(e.target.value)} />
                        </div>
                        <div>
                            <Label htmlFor="username">Username</Label>
                            <Input id="username" value={username} onChange={(e) => setUsername(e.target.value)} />
                        </div>
                        <div>
                            <Label htmlFor="email">Email</Label>
                            <Input id="email" type="email" value={email} onChange={(e) => setEmail(e.target.value)} />
                        </div>
                        <div className="grid grid-cols-2 gap-4">
                            <div>
                                <Label htmlFor="gender">Gender</Label>
                                <select
                                    id="gender"
                                    className="mt-2 w-full border rounded-md h-10 px-3"
                                    value={gender}
                                    onChange={(e) => setGender(e.target.value as any)}
                                >
                                    <option value="">- pilih -</option>
                                    <option value="male">male</option>
                                    <option value="female">female</option>
                                    <option value="other">other</option>
                                </select>
                            </div>
                            <div>
                                <Label htmlFor="birthdate">Birthdate</Label>
                                <Input id="birthdate" type="date" value={birthdate} onChange={(e) => setBirthdate(e.target.value)} />
                            </div>
                        </div>

                        <div>
                            <Label htmlFor="avatar">Avatar (opsional)</Label>
                            <Input
                                id="avatar"
                                type="file"
                                accept="image/*"
                                onChange={(e) => setAvatarFile(e.target.files?.[0] ?? null)}
                            />
                        </div>
                    </div>

                    <DialogFooter className="mt-4">
                        <Button variant="outline" onClick={() => setOpen(false)}>Batal</Button>
                        <Button onClick={onSave} disabled={saving}>{saving ? "Menyimpan..." : "Simpan"}</Button>
                    </DialogFooter>
                </DialogContent>
            </Dialog>


            {/* Statistik + badge section */}
            <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-6">
                <Card>
                    <CardHeader className="flex gap-4 mb-2">
                        <CardTitle className="font-normal text-sm">Total Karbon Yang Dihasilkan</CardTitle>
                        <Badge variant="outline" className="mt-2">
                            <IconTrendingUp className="w-4 h-4 mr-1" /> +11.0%
                        </Badge>
                    </CardHeader>
                    <CardContent>
                        <CardDescription className="text-[28px] font-semibold text-black mb-2">
                            {totalCarbonEmission} Kg CO₂e
                        </CardDescription>
                    </CardContent>
                    <CardFooter>
                        <p className="text-sm text-muted-foreground">Kenaikan Bulan Ini</p>
                    </CardFooter>
                </Card>
                <Card>
                    <CardHeader className="flex gap-4 mb-2">
                        <CardTitle className="font-normal text-sm">Total Misi Yang Diselesaikan</CardTitle>
                        <Badge variant="outline" className="mt-2">
                            <IconTrendingUp className="w-4 h-4 mr-1" /> +11.0%
                        </Badge>
                    </CardHeader>
                    <CardContent>
                        <CardDescription className="text-[28px] font-semibold text-black mb-2">
                            {profile?.missions.length} Misi
                        </CardDescription>
                    </CardContent>
                    <CardFooter>
                        <p className="text-sm text-muted-foreground">Kenaikan Bulan Ini</p>
                    </CardFooter>
                </Card>
                <Card>
                    <CardHeader className="flex gap-4 mb-2">
                        <CardTitle className="font-normal text-sm">Total Misi Yang Diselesaikan</CardTitle>
                        <Badge variant="outline" className="mt-2">
                            <IconTrendingUp className="w-4 h-4 mr-1" /> +11.0%
                        </Badge>
                    </CardHeader>
                    <CardContent>
                        <CardDescription className="text-[28px] font-semibold text-black mb-2">
                            {profile?.user.total_points} Point
                        </CardDescription>
                    </CardContent>
                    <CardFooter>
                        <p className="text-sm text-muted-foreground">Kenaikan Bulan Ini</p>
                    </CardFooter>
                </Card>
            </div>

            {/* Display badges dynamically */}
            <div className="rounded-xl shadow border p-6">
                <h1 className="text-2xl font-bold mb-6">Lencana yang Telah Dimiliki</h1>
                <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-4 gap-6">
                    {profile?.badges && profile?.badges.length > 0 ? (
                        profile?.badges.map((badge: BadgeType) => (
                            <Card key={badge.id}>
                                <CardHeader className="flex flex-col items-center text-center">
                                    <div className="w-16 h-16 relative mb-2">
                                        <Image src={badge.image_url} alt={badge.name} fill className="object-contain" />
                                    </div>
                                    <CardTitle className="text-lg">{badge.name}</CardTitle>
                                </CardHeader>
                                <CardContent className="text-sm text-center text-muted-foreground">
                                    {badge.description}
                                </CardContent>
                            </Card>
                        ))
                    ) : (
                        <p className="col-span-4">Selesaikan Misi untuk mendapatkan Badge!</p>
                    )}
                </div>
            </div>
        </div>
    )
}

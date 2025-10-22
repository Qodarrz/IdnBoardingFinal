"use client"

import { IconTrendingUp, IconArrowUpRight, IconArrowDownRight } from "@tabler/icons-react"
import { Badge } from "@/components/ui/badge"
import {
    Card,
    CardContent,
    CardDescription,
    CardFooter,
    CardHeader,
    CardTitle,
} from "@/components/ui/card"
import { GetCarbonElectronicLogs } from "@/helpers/GetCarbonElectronicLog"
import { useEffect, useState } from "react"
import Image from "next/image"

export function SectionSumCards({ point }: number) {
    const { data: dataElectricty } = GetCarbonElectronicLogs()

    const [totalEmission, setTotalEmission] = useState(0)
    const [averageEmission, setAverageEmission] = useState(0)
    const [highestEmission, setHighestEmission] = useState(0)

    useEffect(() => {
        if (dataElectricty?.data && dataElectricty?.data?.length > 0) {
            // Hitung total emisi karbon
            const total = dataElectricty.data.reduce((acc, d) => acc + d.CarbonEmission, 0)
            setTotalEmission(total)

            // Hitung rata-rata karbon per log
            const average = total / dataElectricty.data.length
            setAverageEmission(average)

            // Hitung karbon tertinggi
            const highest = Math.max(...dataElectricty.data.map(d => d.CarbonEmission))
            setHighestEmission(highest)
        }
    }, [dataElectricty?.data])

    return (
        <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
            {/* Card Total Emission */}
            <Card className="gap-0">
                <CardHeader className="flex gap-4 mb-2">
                    <CardTitle className="font-normal text-sm">Total Karbon Yang Dihasilkan</CardTitle>
                    <Badge variant="outline" className="mt-2">
                        <IconTrendingUp className="w-4 h-4 mr-1" /> +11.0%
                    </Badge>
                </CardHeader>
                <CardContent>
                    <CardDescription className="text-[28px] whitespace-nowrap font-semibold text-black mb-4">
                        {totalEmission.toFixed(2)} kg CO₂e
                    </CardDescription>
                </CardContent>
                <CardFooter>
                    <p className="text-sm text-muted-foreground">Kenaikan Bulan Ini</p>
                </CardFooter>
            </Card>

            {/* Card Average Emission */}
            <Card className="gap-0">
                <CardHeader className="flex gap-4 mb-2">
                    <CardTitle className="font-normal text-sm">Rata-rata Karbon per Perangkat</CardTitle>
                    <Badge variant="outline" className="mt-2">
                        <IconArrowUpRight className="w-4 h-4 mr-1" /> +5.0%
                    </Badge>
                </CardHeader>
                <CardContent>
                    <CardDescription className="text-[28px] whitespace-nowrap font-semibold text-black mb-4">
                        {averageEmission.toFixed(2)} kg CO₂e
                    </CardDescription>
                </CardContent>
                <CardFooter>
                    <p className="text-sm text-muted-foreground">Rata-rata bulan ini</p>
                </CardFooter>
            </Card>

            {/* Card Highest Emission */}
            <Card className="gap-0">
                <CardHeader className="flex gap-4 mb-2">
                    <CardTitle className="font-normal text-sm">Karbon Tertinggi</CardTitle>
                    <Badge variant="outline" className="mt-2">
                        <IconArrowDownRight className="w-4 h-4 mr-1" /> -2.0%
                    </Badge>
                </CardHeader>
                <CardContent>
                    <CardDescription className="text-[28px] whitespace-nowrap font-semibold text-black mb-4">
                        {highestEmission.toFixed(2)} kg CO₂e
                    </CardDescription>
                </CardContent>
                <CardFooter>
                    <p className="text-sm text-muted-foreground">Tertinggi bulan ini</p>
                </CardFooter>
            </Card>

            <Card className="gap-0">
                <CardHeader className="flex gap-4 mb-2 justify-between">
                    <CardTitle className="font-normal text-sm">Total Point</CardTitle>
                    <Image src="/icons/green-point.svg" alt="point" width={32} height={32} />
                </CardHeader>
                <CardContent>
                    <CardDescription className="text-[28px] whitespace-nowrap font-semibold text-black mb-4">
                        {point} Poin
                    </CardDescription>
                </CardContent>
                <CardFooter>
                    <p className="text-sm text-muted-foreground">Total Poin Yang Dikumpulkan</p>
                </CardFooter>
            </Card>
        </div>
    )
}

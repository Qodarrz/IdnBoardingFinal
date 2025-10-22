/* eslint-disable @typescript-eslint/no-explicit-any */
"use client"

import { useState, useEffect } from "react"
import Image from "next/image"
import {
  Card,
  CardHeader,
  CardTitle,
  CardFooter,
} from "@/components/ui/card"
import { LineChart, Line, XAxis, YAxis, Tooltip, ResponsiveContainer } from "recharts"
import { SectionSumCards } from "@/components/shared/section-summary-cards"
import { fetchProfile } from "@/helpers/fetchProfile"

// Helper untuk membuat data grafik berdasarkan emisi kendaraan dan elektronik
const formatCarbonData = (profileData: any) => {
  const vehicleCarbon = profileData?.monthly_vehicle_carbon || []
  const electronicCarbon = profileData?.monthly_electronic_carbon || []

  // Gabungkan data kendaraan dan elektronik
  const combinedData = [...vehicleCarbon, ...electronicCarbon]

  // Gabungkan data yang memiliki bulan yang sama, jumlahkan emisi karbonnya
  const carbonMap = new Map()

  combinedData.forEach(item => {
    const month = new Date(item.month).toLocaleString('default', { month: 'short' })
    if (carbonMap.has(month)) {
      carbonMap.set(month, carbonMap.get(month) + item.total_carbon_emission_g)
    } else {
      carbonMap.set(month, item.total_carbon_emission_g)
    }
  })


  // Tentukan bulan pertama (bulan paling awal) dan bulan terakhir
  const firstMonth = new Date(Math.min(...combinedData.map(item => new Date(item.month).getTime())))
  const lastMonth = new Date()

  // Generate data bulan yang hilang (misalnya Juli jika tidak ada data)
  const completeData = []
  const currentMonth = firstMonth

  while (currentMonth <= lastMonth) {
    const monthName = currentMonth.toLocaleString('default', { month: 'short' })
    const value = carbonMap.has(monthName) ? carbonMap.get(monthName) : 0
    completeData.push({ name: monthName, value })

    // Tambahkan satu bulan
    currentMonth.setMonth(currentMonth.getMonth() + 1)
  }

  // Ambil hanya 6 bulan terakhir
  return completeData.slice(-6)
}

export default function DashboardPage() {
  const [carbonData, setCarbonData] = useState<any[]>([])
  const [loading, setLoading] = useState(true)
  const [point, setPoint] = useState<any>([])
  const [error, setError] = useState<string | null>(null) 

  useEffect(() => {
    const loadData = async () => {
      const token = localStorage.getItem("authtoken")
      if (!token) {
        setError("Token tidak ditemukan.")
        setLoading(false)
        return
      }

      const profile = await fetchProfile(token)

      if (profile) {
        const user = profile.user
        const data = formatCarbonData(profile)
        setCarbonData(data)
        setPoint(user.total_points)
      } else {
        setError("Gagal memuat data.")
      }
      setLoading(false)
    }

    loadData()
  }, [])

  if (loading) console.log("loading")
  if (error) console.log(error)

  return (
    <div className="p-6 space-y-6">
      {/* Header Greeting */}
      <div className="flex flex-col md:flex-row items-center justify-between gap-6 bg-card p-6 rounded-2xl shadow-sm">
        <Image src="/images/dashboard.png" alt="" width={200} height={200} className="object-contain md:hidden" />
        <div>
          <h1 className="text-2xl font-semibold">
            Hallo, <span className="text-primary">JiePrass</span>
          </h1>
          <p className="text-muted-foreground max-w-sm">
            Pantau jejak emisi karbon Anda setiap hari dan mulailah langkah kecil menuju gaya hidup yang lebih ramah lingkungan.
          </p>
        </div>
        <div className="gap-6 justify-center hidden md:flex">
          <Image src="/images/dashboard.png" alt="" width={200} height={200} className="object-contain" />
          <Image src="/images/dashboard-2.png" alt="" width={200} height={200} className="object-contain" />
        </div>
      </div>

      {/* Statistik Cards */}
      <SectionSumCards point={point} />

      {/* Grafik Tren */}
      <Card>
        <CardHeader>
          <CardTitle>Tren Pengurangan Karbon Bulanan</CardTitle>
        </CardHeader>
        <CardFooter>
          <div className="w-full overflow-x-auto">
            <div className="h-64 min-w-[600px]">
              <ResponsiveContainer width="100%" height="100%">
                <LineChart data={carbonData}>
                  <XAxis dataKey="name" />
                  <YAxis />
                  <Tooltip />
                  <Line type="monotone" dataKey="value" stroke="#16a34a" />
                </LineChart>
              </ResponsiveContainer>
            </div>
          </div>
        </CardFooter>
      </Card>
    </div>
  )
}

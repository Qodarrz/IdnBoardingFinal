"use client"

import { useEffect, useState, } from "react"
import Image from "next/image"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger } from "@/components/ui/dialog"
import { ScrollArea, ScrollBar } from "@/components/ui/scroll-area"
import { PostVehicleTracker } from "@/helpers/PostVehicleTracker"
import { useAuthMe } from "@/helpers/AuthMe"
import { GetVehicleTracker } from "@/helpers/GetVehicleTracker"
import { PostVehicleTrackerLog } from "@/helpers/PostVehicleTrackerLog"

type Vehicle = {
  id?: number
  ID?: number
  title: string
  type: "car" | "motorcycle" | "public_transport" | "walk"
  image: string
  total: string
  percentage: string
  active: boolean
  startLocation?: { lat: number; lng: number }
  endLocation?: { lat: number; lng: number }
  distance?: number // dalam meter
  watchId?: number | null
}

export default function VehicleSlider() {
  const { data: dataMe } = useAuthMe()
  const { data: vehicle } = GetVehicleTracker()

  const [vehicles, setVehicles] = useState<Vehicle[]>([])

  useEffect(() => {
    if (vehicle?.data && vehicle?.data?.length > 0) {
      console.log(vehicle.data)
      const initialDevices = vehicle?.data
        ?.filter((item) => item.VehicleType !== "public_transport")
        .map((d) => ({
          id: d.ID,
          title: d.Name,
          image: "/icons/motor.svg",
          type: d.type,
          percentage: "0%",
          total: (d.LatestLog !== null ? (parseFloat(d.LatestLog.CarbonEmission).toFixed(2)) : 0) + " kg CO₂e",
          active: false,
          watchId: null,
        }))
      setVehicles(initialDevices ?? [])
    }
  }, [vehicle?.data])

  const [open, setOpen] = useState(false)
  const [newVehicle, setNewVehicle] = useState({
    name: "",
    vehicle_type: "car",
    fuel_type: "petrol",
  })

  // fungsi hitung jarak haversine (meter) ini generate oleh Ai
  const haversine = (lat1: number, lon1: number, lat2: number, lon2: number) => {
    const R = 6371e3
    const toRad = (x: number) => (x * Math.PI) / 180
    const φ1 = toRad(lat1)
    const φ2 = toRad(lat2)
    const Δφ = toRad(lat2 - lat1)
    const Δλ = toRad(lon2 - lon1)
    const a =
      Math.sin(Δφ / 2) ** 2 +
      Math.cos(φ1) * Math.cos(φ2) * Math.sin(Δλ / 2) ** 2
    const c = 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1 - a))
    return R * c
  }

  // aktifkan GPS untuk 1 kendaraan
  const handleToggleVehicle = (v: Vehicle) => {
    if (!navigator.geolocation) {
      alert("Geolocation tidak didukung browser ini.");
      return;
    }

    if (!v.active) {
      // Matikan kendaraan lain dulu biar tidak semuanya aktif
      setVehicles((prev) =>
        prev.map((item) =>
          item.id === v.ID
            ? item // biarkan kendaraan yang dipilih tetap
            : { ...item, active: false, watchId: null }
        )
      );

      // Aktifkan kendaraan ini
      const watchId = navigator.geolocation.watchPosition(
        (pos) => {
          const { latitude, longitude } = pos.coords;
          setVehicles((prev) =>
            prev.map((item) =>
              item.id === v.id
                ? {
                    ...item,
                    active: true,
                    startLocation: item.startLocation ?? {
                      lat: latitude,
                      lng: longitude,
                    }, // snap awal sekali
                  }
                : item
            )
          );
        },
        (err) => {
          alert("Gagal ambil lokasi: " + err.message);
        },
        { enableHighAccuracy: true, maximumAge: 0, timeout: 10000 }
      );

      // simpan watchId ke kendaraan ini
      setVehicles((prev) =>
        prev.map((item) =>
          item.id === v.id ? { ...item, active: true, watchId } : item
        )
      );
    } else {
      // Matikan kendaraan ini
      if (v.watchId) {
        navigator.geolocation.clearWatch(v.watchId);
      }

      navigator.geolocation.getCurrentPosition(
        async (pos) => {
          const { latitude, longitude } = pos.coords;
          const distance =
            v.startLocation &&
            haversine(
              v.startLocation.lat,
              v.startLocation.lng,
              latitude,
              longitude
            );

          setVehicles((prev) =>
            prev.map((item) =>
              item.id === v.id
                ? {
                    ...item,
                    active: false,
                    endLocation: { lat: latitude, lng: longitude },
                    distance: distance ?? 0,
                    watchId: null,
                  }
                : item
            )
          );

          const res = await PostVehicleTrackerLog({
            vehicle_id: v.id,
            distance_km: parseFloat(((distance ?? 0) / 1000).toFixed(2)),
            start_lat: v.startLocation?.lat ?? 0,
            start_lon: v.startLocation?.lng ?? 0,
            end_lat: latitude,
            end_lon: longitude,
          });

          if (res) console.log(res);

          if (distance) {
            alert(
              `Jarak tempuh kendaraan "${v.title}" adalah ${(distance / 1000).toFixed(
                2
              )} KM`
            );
          }
        },
        (err) => {
          alert("Gagal ambil lokasi akhir: " + err.message);
        }
      );
    }
  };


  // tambah kendaraan baru
  const handleAddVehicle = async () => {
    if (!newVehicle.name) return
    try {
      const res = await PostVehicleTracker({
        user_id: dataMe?.data?.id ? parseInt(dataMe.data.id) : null,
        ...newVehicle,
      })

      if (res) {
        const newData: Vehicle = {
          id: vehicles.length + 1,
          title: newVehicle.name,
          type: newVehicle.vehicle_type as
            | "car"
            | "motorcycle"
            | "public_transport"
            | "walk",
          image: "/icons/motor.svg",
          total: `0 kg CO₂e`,
          percentage: "0%",
          active: false,
          watchId: null,
        }
        setVehicles([...vehicles, newData])
        setNewVehicle({ name: "", vehicle_type: "car", fuel_type: "petrol" })
        setOpen(false)
      }
    } catch (err) {
      console.error("Error add vehicle:", err)
    }
  }

  return (
    <section className="space-y-6">
      <h2 className="text-xl font-semibold">Kendaraan Pribadi</h2>
      <div className="w-full overflow-hidden">
        <ScrollArea className="max-w-[166vh] whitespace-nowrap">
          <div className="flex gap-4 pb-4">
            {vehicles.map((v) => (
              <Card key={v.id} className="w-64 flex-shrink-0">
                <CardHeader className="flex flex-col items-center justify-center">
                  <Image
                    src={v.image}
                    alt={v.title}
                    width={80}
                    height={80}
                    className="aspect-square rounded-full p-1 object-contain bg-[#ECECEC]"
                  />
                  <CardTitle className="mt-2 text-center">{v.title}</CardTitle>
                </CardHeader>
                <CardContent className="text-center space-y-2">
                  <p className="text-sm text-muted-foreground">
                    Total Karbon Terakhir <br /> {v.title}
                  </p>
                  <p className="text-lg font-semibold">{v.total}</p>
                  {/* <p className="text-xs text-gray-500">{v.percentage}</p> */}
                  <Button
                    variant={v.active ? "destructive" : "outline"}
                    className="w-full"
                    onClick={() => handleToggleVehicle(v)}
                  >
                    {v.active ? "Matikan Kendaraan" : "Aktifkan Kendaraan"}
                  </Button>

                  {v.startLocation && (
                    <p className="mt-2 text-xs">
                      Start: {v.startLocation.lat.toFixed(5)},{" "}
                      {v.startLocation.lng.toFixed(5)}
                    </p>
                  )}
                  {v.endLocation && (
                    <p className="mt-2 text-xs">
                      End: {v.endLocation.lat.toFixed(5)},{" "}
                      {v.endLocation.lng.toFixed(5)}
                    </p>
                  )}
                  {v.distance !== undefined && !v.active && (
                    <p className="mt-2 text-xs font-semibold">
                      Jarak: {(v.distance / 1000).toFixed(2)} KM
                    </p>
                  )}
                </CardContent>
              </Card>
            ))}

            {/* Tambah Kendaraan */}
            <Dialog open={open} onOpenChange={setOpen}>
              <DialogTrigger asChild>
                <Card className="w-64 flex-shrink-0 flex items-center justify-center cursor-pointer hover:bg-muted">
                  <CardContent className="text-center p-6">
                    <p className="text-4xl">+</p>
                    <p className="mt-2 text-sm">Tambahkan Kendaraan Anda</p>
                  </CardContent>
                </Card>
              </DialogTrigger>
              <DialogContent>
                <DialogHeader>
                  <DialogTitle>Tambah Kendaraan Baru</DialogTitle>
                </DialogHeader>
                <div className="space-y-4">
                  <input
                    type="text"
                    placeholder="Nama Kendaraan"
                    value={newVehicle.name}
                    onChange={(e) =>
                      setNewVehicle({ ...newVehicle, name: e.target.value })
                    }
                    className="w-full border rounded-md p-2"
                  />
                  <select
                    className="w-full border rounded-md p-2"
                    value={newVehicle.vehicle_type}
                    onChange={(e) =>
                      setNewVehicle({
                        ...newVehicle,
                        vehicle_type: e.target.value,
                      })
                    }
                  >
                    <option value="car">Car</option>
                    <option value="motorcycle">Motorcycle</option>
                    <option value="bicycle">Bicycle</option>
                    <option value="public_transport">Public Transport</option>
                    <option value="walk">Walk</option>
                  </select>
                  <select
                    className="w-full border rounded-md p-2"
                    value={newVehicle.fuel_type}
                    onChange={(e) =>
                      setNewVehicle({
                        ...newVehicle,
                        fuel_type: e.target.value,
                      })
                    }
                  >
                    <option value="petrol">Petrol</option>
                    <option value="diesel">Diesel</option>
                    <option value="electric">Electric</option>
                    <option value="none">None</option>
                  </select>
                  <Button className="w-full" onClick={handleAddVehicle}>
                    Simpan
                  </Button>
                </div>
              </DialogContent>
            </Dialog>
          </div>
          <ScrollBar orientation="horizontal" />
        </ScrollArea>
      </div>
    </section>
  )
}

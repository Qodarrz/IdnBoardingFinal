"use client"

import { useState, useEffect } from "react"
import { SectionElectricCards } from "@/components/shared/section-electric-cards"
import DeviceTable from "@/components/shared/device-table"
import { IconArrowsSort, IconFilterFilled, IconPlus, IconSearch } from "@tabler/icons-react"
import { Input } from "@/components/ui/input"
import { Button } from "@/components/ui/button"
import {
    Dialog,
    DialogContent,
    DialogHeader,
    DialogTitle,
    DialogTrigger,
    DialogFooter,
} from "@/components/ui/dialog"
import { Label } from "@/components/ui/label"
import { PostElectricityDevice } from "@/helpers/PostElectricityDevice"
import { PostCarbonElectronicLog } from "@/helpers/PostCarbonElectronicLog"
import { GetElectricityDevice } from "@/helpers/GetElectricityDevice"
import { GetCarbonElectronicLogs } from "@/helpers/GetCarbonElectronicLog"
import { useAuthMe } from "@/helpers/AuthMe"
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table"

// Device Types
type DeviceTypes = {
    id: number;
    name: string;
    type: string;
    power_watts: number;
    status: string;
}

// Carbon Log Types
type CarbonElectronicLogTypes = {
    id: number;
    device_id: number;
    duration_hours: number;
    created_at: string;
    carbon_emission: number; // Tambahkan properti emisi karbon
}

const DEVICE_TYPES = [
    "Kulkas", "Lampu", "Mesin Cuci", "TV", "Komputer", "Laptop", "Smartphone", "Microwave", "Fan", "AC"
];

export default function ElectrictyTracker() {
    const { data: dataMe } = useAuthMe()
    const { data: dataCarbonElectronicLogs, loading: loadingCarbonLogs, error: errorCarbonLogs } = GetCarbonElectronicLogs()
    const { data: dataElectricty, loading: loadingDevices, error: errorDevices } = GetElectricityDevice()
    const [devices, setDevices] = useState<DeviceTypes[]>([])
    const [carbonElectronicLogs, setCarbonElectronicLogs] = useState<CarbonElectronicLogTypes[]>([])
    const [search, setSearch] = useState("");
    const [open, setOpen] = useState(false)
    const [openCarbon, setOpenCarbon] = useState(false)
    const [newDevice, setNewDevice] = useState({ name: "", type: "", power_watts: "" })
    const [newCarbonElectronicLog, setNewCarbonElectronicLog] = useState({ device_id: 0, duration_hours: "" })

    // Ambil data perangkat yang sudah ada
    useEffect(() => {
        if (dataElectricty && dataElectricty?.data && dataElectricty?.data?.length > 0) {
            const initialDevices = dataElectricty?.data?.map((d) => ({
                id: d.id,
                name: d.device_name,
                type: d.device_type,
                power_watts: d.power_watts,
                status: d.power_watts <= 20 ? "Rendah" : d.power_watts >= 20 ? 'Normal' : 'Tinggi',
            }));
            setDevices(initialDevices ?? [])
        }
    }, [dataElectricty])

    // Ambil data carbon logs dan perbaiki pemetaan data
    useEffect(() => {
        if (dataCarbonElectronicLogs && dataCarbonElectronicLogs?.data && dataCarbonElectronicLogs?.data?.length > 0) {
            console.log("Raw Carbon Logs Data:", dataCarbonElectronicLogs.data);  // Cek data yang diterima
            const initialCarbonElectronicLogs = dataCarbonElectronicLogs?.data?.map((log) => ({
                id: log.ID,  // Pemetaan ID
                device_id: log.DeviceID,  // Pemetaan DeviceID
                duration_hours: log.DurationHours,  // Pemetaan DurationHours
                created_at: log.LoggedAt,  // Pemetaan LoggedAt
                carbon_emission: log.CarbonEmission,  // Pemetaan CarbonEmission
            }));
            setCarbonElectronicLogs(initialCarbonElectronicLogs ?? []);
        }
    }, [dataCarbonElectronicLogs]);

    // Fungsi untuk menambahkan perangkat baru
    const handleAddDevice = async () => {
        if (!newDevice.name || !newDevice.power_watts || !newDevice.type) return

        const device = {
            id: devices.length + 1,
            name: newDevice.name,
            type: newDevice.type,
            power_watts: parseFloat(newDevice.power_watts),
            status: parseFloat(newDevice.power_watts) <= 20 ? "Rendah" : parseFloat(newDevice.power_watts) >= 20 ? 'Normal' : 'Tinggi',
        }

        try {
            const result = await PostElectricityDevice({
                device_name: newDevice.name,
                power_watts: parseFloat(newDevice.power_watts),
                device_type: newDevice.type,
                user_id: dataMe?.data?.id ? parseInt(dataMe.data.id) : null
            })

            console.log(result)
            setDevices([...devices, device])
            setNewDevice({ name: "", type: '', power_watts: "" })
            setOpen(false)
        } catch (error) {
            console.error('Error while adding device:', error)
            setOpen(false)
        }
    }

    // Fungsi untuk menambahkan carbon log baru
    const handleAddCarbonElectronicLog = async () => {
        if (!newCarbonElectronicLog.device_id || !newCarbonElectronicLog.duration_hours) return

        const carbonElectronicLog = {
            id: carbonElectronicLogs.length + 1,
            device_id: newCarbonElectronicLog.device_id,
            duration_hours: parseFloat(newCarbonElectronicLog.duration_hours),
            created_at: new Date().toISOString(),
            carbon_emission: parseFloat(newCarbonElectronicLog.duration_hours) * 0.0275 // Asumsi: perhitungan emisi karbon berdasarkan durasi
        }

        try {
            const result = await PostCarbonElectronicLog({
                device_id: newCarbonElectronicLog.device_id,
                duration_hours: parseFloat(newCarbonElectronicLog.duration_hours),
            })

            console.log(result)
            setCarbonElectronicLogs([...carbonElectronicLogs, carbonElectronicLog])
            setNewCarbonElectronicLog({ device_id: 0, duration_hours: "" })
            setOpen(false)
        } catch (error) {
            console.error('Error while adding carbon log:', error)
            setOpen(false)
        }
    }

    // Apply search, filter, and sort for devices
    const filteredDevices = devices!
        .filter((d) => d.name?.toLowerCase().includes(search.toLowerCase()))

    // Apply search, filter, and sort for carbon logs
    const filteredCarbonElectronicLogs = carbonElectronicLogs!
        .filter((log) => {
            // Pastikan log.device_id ada
            if (!log.device_id || !log.duration_hours || !log.created_at) return false;  // Cek semua properti yang diperlukan
            return log.device_id.toString().includes(search.toLowerCase());
        });

    // Format Tanggal
    const formatDate = (dateString: string) => {
        const date = new Date(dateString);
        const day = String(date.getDate()).padStart(2, '0');
        const month = String(date.getMonth() + 1).padStart(2, '0');
        const year = date.getFullYear();
        const formattedDate = `${day}/${month}/${year}`;

        const hours = String(date.getHours()).padStart(2, '0');
        const minutes = String(date.getMinutes()).padStart(2, '0');
        const formattedTime = `${hours}:${minutes}`;

        return { formattedDate, formattedTime };
    }

    if (loadingDevices || loadingCarbonLogs) {
        return <div>Loading...</div>
    }

    if (errorDevices || errorCarbonLogs) {
        return <div>Error: {errorDevices || errorCarbonLogs}</div>
    }

    return (
        <div className="p-6 space-y-6">
            {/* Statistik Cards */}
            <SectionElectricCards />

            {/* Action Bar for Devices */}
            <div className="flex justify-between items-center bg-gray-100 px-6 py-3 rounded-md w-full overflow-x-auto mb-6">
                <h2 className="font-semibold text-lg whitespace-nowrap">Daftar Alat Elektronik</h2>
                <div className="flex gap-4 items-center">
                    <Dialog open={open} onOpenChange={setOpen}>
                        <DialogTrigger asChild>
                            <Button size="icon" variant="ghost">
                                <IconPlus />
                            </Button>
                        </DialogTrigger>
                        <DialogContent>
                            <DialogHeader>
                                <DialogTitle>Tambah Device Baru</DialogTitle>
                            </DialogHeader>
                            <div className="space-y-4">
                                <div>
                                    <Label>Nama Device</Label>
                                    <Input
                                        value={newDevice.name}
                                        onChange={(e) =>
                                            setNewDevice({ ...newDevice, name: e.target.value })
                                        }
                                        placeholder="Contoh: AC Ruang Tamu"
                                    />
                                </div>
                                <div>
                                    <Label>Tipe Device</Label>
                                    <select
                                        value={newDevice.type}
                                        onChange={(e) =>
                                            setNewDevice({ ...newDevice, type: e.target.value })
                                        }
                                        className="w-full border px-3 py-2 rounded-md"
                                    >
                                        <option value="">Pilih Tipe Perangkat</option>
                                        {DEVICE_TYPES.map((type, index) => (
                                            <option key={index} value={type.toLocaleLowerCase()}>{type}</option>
                                        ))}
                                    </select>
                                </div>
                                <div>
                                    <Label>Konsumsi</Label>
                                    <Input
                                        value={newDevice.power_watts}
                                        onChange={(e) =>
                                            setNewDevice({ ...newDevice, power_watts: e.target.value })
                                        }
                                        placeholder="Contoh: 20 kWh"
                                    />
                                </div>
                            </div>
                            <DialogFooter>
                                <Button variant="outline" onClick={() => setOpen(false)}>Batal</Button>
                                <Button onClick={handleAddDevice}>Tambah</Button>
                            </DialogFooter>
                        </DialogContent>
                    </Dialog>

                    <Button size="icon" variant="ghost">
                        <IconFilterFilled />
                    </Button>
                    <Button size="icon" variant="ghost">
                        <IconArrowsSort />
                    </Button>
                </div>

                {/* SearchBar */}
                <div className="w-64 relative">
                    <div className="absolute left-2 -translate-y-1/2 top-1/2 opacity-60">
                        <IconSearch />
                    </div>
                    <Input
                        className="bg-gray-50 pl-10"
                        placeholder="Cari perangkat..."
                        value={search}
                        onChange={(e) => setSearch(e.target.value)}
                    />
                </div>
            </div>

            {/* Device Table */}
            <DeviceTable
                devices={filteredDevices}
                onUpdate={() => { }}
                onDelete={() => { }}
            />

            {/* Action Bar for Carbon Logs */}
            <div className="flex justify-between items-center bg-gray-100 px-6 py-3 rounded-md w-full overflow-x-auto mb-6">
                <h2 className="font-semibold text-lg whitespace-nowrap">Daftar Carbon Logs</h2>
                <div className="flex gap-4 items-center">
                    <Dialog open={openCarbon} onOpenChange={setOpenCarbon}>
                        <DialogTrigger asChild>
                            <Button size="icon" variant="ghost">
                                <IconPlus />
                            </Button>
                        </DialogTrigger>
                        <DialogContent>
                            <DialogHeader>
                                <DialogTitle>Tambah Carbon Log</DialogTitle>
                            </DialogHeader>
                            <div className="space-y-4">
                                <div>
                                    <Label>Device</Label>
                                    <select
                                        value={newCarbonElectronicLog.device_id}
                                        onChange={(e) =>
                                            setNewCarbonElectronicLog({ ...newCarbonElectronicLog, device_id: parseInt(e.target.value) })
                                        }
                                        className="w-full border px-3 py-2 rounded-md"
                                    >
                                        <option value="">Pilih Device</option>
                                        {devices.map((device) => (
                                            <option key={device.id} value={device.id}>
                                                {device.name} ({device.type}) {/* Display name and type */}
                                            </option>
                                        ))}
                                    </select>
                                </div>
                                <div>
                                    <Label>Durasi (Jam)</Label>
                                    <Input
                                        type="number"
                                        value={newCarbonElectronicLog.duration_hours}
                                        onChange={(e) =>
                                            setNewCarbonElectronicLog({ ...newCarbonElectronicLog, duration_hours: e.target.value })
                                        }
                                        placeholder="Masukkan durasi dalam jam"
                                    />
                                </div>
                            </div>
                            <DialogFooter>
                                <Button variant="outline" onClick={() => setOpenCarbon(false)}>Batal</Button>
                                <Button onClick={handleAddCarbonElectronicLog}>Tambah</Button>
                            </DialogFooter>
                        </DialogContent>
                    </Dialog>

                    <Button size="icon" variant="ghost">
                        <IconFilterFilled />
                    </Button>
                    <Button size="icon" variant="ghost">
                        <IconArrowsSort />
                    </Button>
                </div>

                {/* SearchBar */}
                <div className="w-64 relative">
                    <div className="absolute left-2 -translate-y-1/2 top-1/2 opacity-60">
                        <IconSearch />
                    </div>
                    <Input
                        className="bg-gray-50 pl-10"
                        placeholder="Cari carbon logs..."
                        value={search}
                        onChange={(e) => setSearch(e.target.value)}
                    />
                </div>
            </div>

            {/* Carbon Log Table */}
            <div className="rounded-md border">
                <Table>
                    <TableHeader>
                        <TableRow>
                            <TableHead>Device ID</TableHead>
                            <TableHead>Perangkat</TableHead>
                            <TableHead>Durasi (Jam)</TableHead>
                            <TableHead>Emisi</TableHead>
                            <TableHead>Jam</TableHead>
                            <TableHead>Tanggal</TableHead>
                        </TableRow>
                    </TableHeader>
                    <TableBody>
                        {filteredCarbonElectronicLogs.length > 0 ? filteredCarbonElectronicLogs.map((log) => {
                            const device = devices.find((d) => d.id === log.device_id);  // Temukan nama perangkat berdasarkan device_id
                            const { formattedDate, formattedTime } = formatDate(log.created_at); // Pecah tanggal dan waktu
                            return (
                                <TableRow key={log.id}>
                                    <TableCell>{log.device_id}</TableCell>
                                    <TableCell>{device ? device.name : "Tidak ditemukan"}</TableCell>
                                    <TableCell>{log.duration_hours}</TableCell>
                                    <TableCell>{log.carbon_emission}</TableCell>
                                    <TableCell>{formattedTime}</TableCell>
                                    <TableCell>{formattedDate}</TableCell>
                                </TableRow>
                            );
                        }) : (
                            <TableRow>
                                <TableCell colSpan={6} className="text-center">Tidak ada data</TableCell>
                            </TableRow>
                        )}
                    </TableBody>
                </Table>
            </div>
        </div>
    )
}

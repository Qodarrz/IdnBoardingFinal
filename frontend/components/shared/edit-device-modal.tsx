import {
    Dialog, DialogContent, DialogHeader, DialogTitle, DialogFooter
} from "@/components/ui/dialog"
import { Input } from "@/components/ui/input"
import { Button } from "@/components/ui/button"

type Device = {
    id: number
    name: string
    type: string
    power_watts: number
    status: string // Pastikan status ada di tipe Device
}

interface EditDeviceProps {
    editOpen: boolean
    setEditOpen: (open: boolean) => void
    selectedDevice: Device | null
    setSelectedDevice: (device: Device) => void // setSelectedDevice harus menerima tipe Device
    handleEdit: () => void
}

const deviceTypes = [
    "kulkas", "lampu", "mesin cuci", "tv", "komputer", "laptop", "smartphone", "microwave", "fan", "ac"
]

export default function EditDeviceModal({
    editOpen, setEditOpen, selectedDevice, setSelectedDevice, handleEdit
}: EditDeviceProps) {

    // Pastikan selectedDevice tidak null sebelum merender
    if (!selectedDevice) return null;

    return (
        <Dialog open={editOpen} onOpenChange={setEditOpen}>
            <DialogContent>
                <DialogHeader>
                    <DialogTitle>Edit Device</DialogTitle>
                </DialogHeader>
                <div className="space-y-4">
                    <Input
                        value={selectedDevice.name}
                        onChange={(e) =>
                            setSelectedDevice({ ...selectedDevice, name: e.target.value })
                        }
                        placeholder="Nama device"
                    />
                    {/* Input untuk konsumsi (power_watts) */}
                    <Input
                        value={selectedDevice.power_watts}
                        onChange={(e) =>
                            setSelectedDevice({
                                ...selectedDevice,
                                power_watts: parseInt(e.target.value)
                            })
                        }
                        placeholder="Konsumsi (kWh)"
                        type="number" // Menambahkan type number untuk input konsumsi
                    />
                    {/* Dropdown untuk memilih tipe device */}
                    <select
                        value={selectedDevice.type}
                        onChange={(e) =>
                            setSelectedDevice({
                                ...selectedDevice,
                                type: e.target.value
                            })
                        }
                        className="w-full border p-2 rounded-md"
                    >
                        {deviceTypes.map((deviceType) => (
                            <option key={deviceType} value={deviceType}>
                                {deviceType}
                            </option>
                        ))}
                    </select>
                </div>
                <DialogFooter>
                    <Button variant="outline" onClick={() => setEditOpen(false)}>
                        Batal
                    </Button>
                    <Button onClick={handleEdit}>Simpan</Button>
                </DialogFooter>
            </DialogContent>
        </Dialog>
    )
}

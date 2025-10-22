"use client"

import { useState } from "react"
import {
  Table, TableBody, TableCell, TableHead, TableHeader, TableRow
} from "@/components/ui/table"
import { Badge } from "@/components/ui/badge"
import {
  DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import { IconDots, IconPencil, IconTrash } from "@tabler/icons-react"
import {
  Dialog, DialogContent, DialogHeader, DialogTitle, DialogFooter
} from "@/components/ui/dialog"
import { Button } from "@/components/ui/button"
import { UpdateElectricityDevice } from "@/helpers/UpdateElectricityDevice"
import { DeleteElectricityDevice } from "@/helpers/DeleteElectricityDevice"
import EditDeviceModal from "./edit-device-modal"

type DeviceTypes = {
  id: number
  name: string
  type: string
  power_watts: number
  status: string
}

interface DeviceTableProps {
  devices: DeviceTypes[]; 
  onUpdate: (updatedDevice: DeviceTypes) => void;
  onDelete: (id: number) => void;
}

export default function DeviceTable({ devices, onUpdate, onDelete }: DeviceTableProps) {
  const [editOpen, setEditOpen] = useState(false)
  const [deleteOpen, setDeleteOpen] = useState(false)
  const [selectedDevice, setSelectedDevice] = useState<DeviceTypes | null>(null)

  console.log(devices)

  // Submit Edit
  const handleEdit = async () => {
    if (!selectedDevice) return; // Pastikan selectedDevice ada

    // Membuat payload berdasarkan selectedDevice
    const payload = {
      device_name: selectedDevice.name,
      power_watts: selectedDevice.power_watts,
      device_type: selectedDevice.type, // Menyertakan type perangkat
    };

    try {
      // Memanggil helper untuk update perangkat
      const result = await UpdateElectricityDevice(selectedDevice.id, payload);

      // Mengupdate perangkat di state
      onUpdate(result.data);
      setEditOpen(false);
    } catch (error) {
      console.error("Error updating device:", error);
    }
  };

  // Submit Delete
  const handleDelete = async () => {
    if (!selectedDevice) return
    try {
      await DeleteElectricityDevice(selectedDevice.id)
      onDelete(selectedDevice.id)
      setDeleteOpen(false)
    } catch (e) {
      console.error(e)
    }
  }

  return (
    <div className="rounded-md border">
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>Perangkat</TableHead>
            <TableHead>Tipe</TableHead>
            <TableHead>Konsumsi (kWh)</TableHead>
            <TableHead>Status</TableHead>
            <TableHead className="w-[40px]" />
          </TableRow>
        </TableHeader>
        <TableBody>
          {devices.length > 0 ? devices.map((device) => (
            <TableRow key={device.id}>
              <TableCell>{device.name}</TableCell>
              <TableCell>{device.type}</TableCell>
              <TableCell>{device.power_watts}</TableCell>
              <TableCell>
                <Badge className={device.status === "Aktif" || device.status === "Rendah"
                  ? "bg-emerald-100 text-emerald-700 hover:bg-emerald-200"
                  : "bg-red-100 text-red-700 hover:bg-red-200"}>
                  {device.status}
                </Badge>
              </TableCell>
              <TableCell>
                <DropdownMenu>
                  <DropdownMenuTrigger asChild>
                    <button className="p-1 rounded hover:bg-muted">
                      <IconDots size={18} />
                    </button>
                  </DropdownMenuTrigger>
                  <DropdownMenuContent align="end">
                    <DropdownMenuItem
                      onClick={() => {
                        setSelectedDevice(device)
                        setEditOpen(true)
                      }}
                    >
                      <IconPencil className="mr-2 h-4 w-4" /> Edit
                    </DropdownMenuItem>
                    <DropdownMenuItem
                      className="text-red-600 focus:text-red-600"
                      onClick={() => {
                        setSelectedDevice(device)
                        setDeleteOpen(true)
                      }}
                    >
                      <IconTrash className="mr-2 h-4 w-4" /> Hapus
                    </DropdownMenuItem>
                  </DropdownMenuContent>
                </DropdownMenu>
              </TableCell>
            </TableRow>
          )) : (
            <TableRow>
              <TableCell colSpan={7} className="text-center">
                Loading..
              </TableCell>
            </TableRow>
          )}
        </TableBody>
      </Table>

      {/* Modal Edit */}
      <EditDeviceModal
        editOpen={editOpen}
        setEditOpen={setEditOpen}
        selectedDevice={selectedDevice}
        setSelectedDevice={setSelectedDevice}
        handleEdit={handleEdit}
      />


      {/* Modal Delete */}
      <Dialog open={deleteOpen} onOpenChange={setDeleteOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Hapus Device</DialogTitle>
          </DialogHeader>
          <p>Yakin ingin menghapus {selectedDevice?.name}?</p>
          <DialogFooter>
            <Button variant="outline" onClick={() => setDeleteOpen(false)}>Batal</Button>
            <Button variant="destructive" onClick={handleDelete}>Hapus</Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  )
}

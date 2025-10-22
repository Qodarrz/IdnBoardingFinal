import { Button } from "../ui/button";
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from "../ui/dialog";
import { Input } from "../ui/input";
import { Label } from "../ui/label";
import { Select } from "../ui/select";

interface AddDeviceModalProps {
    open: boolean;
    setOpen: React.Dispatch<React.SetStateAction<boolean>>;
    newDevice: { name: string; type: string; power_watts: string };
    setNewDevice: React.Dispatch<React.SetStateAction<{ name: string; type: string; power_watts: string }>>;
    handleAdd: () => void;
}

const AddDeviceModal: React.FC<AddDeviceModalProps> = ({ open, setOpen, newDevice, setNewDevice, handleAdd }) => {
    return (
        <Dialog open={open} onOpenChange={setOpen}>
            <DialogContent>
                <DialogHeader>
                    <DialogTitle>Tambah Device Baru</DialogTitle>
                    <DialogDescription>
                        Isi detail perangkat untuk menambahkannya ke daftar.
                    </DialogDescription>
                </DialogHeader>
                <div className="space-y-4">
                    <div>
                        <Label>Nama Device</Label>
                        <Input
                            value={newDevice.name}
                            onChange={(e) => setNewDevice({ ...newDevice, name: e.target.value })}
                            placeholder="Contoh: AC Ruang Tamu"
                        />
                    </div>
                    <div>
                        <Label>Tipe Device</Label>
                        <Select
                            value={newDevice.type}
                            onValueChange={(value) => setNewDevice({ ...newDevice, type: value })} // Correct event handler
                        >
                            <option value="kulkas">Kulkas</option>
                            <option value="lampu">Lampu</option>
                            <option value="mesin cuci">Mesin Cuci</option>
                            <option value="tv">TV</option>
                            <option value="komputer">Komputer</option>
                            <option value="laptop">Laptop</option>
                            <option value="smartphone">Smartphone</option>
                            <option value="microwave">Microwave</option>
                            <option value="fan">Fan</option>
                            <option value="ac">AC</option>
                        </Select>
                    </div>
                    <div>
                        <Label>Konsumsi</Label>
                        <Input
                            value={newDevice.power_watts}
                            onChange={(e) => setNewDevice({ ...newDevice, power_watts: e.target.value })}
                            placeholder="Contoh: 20 kWh"
                        />
                    </div>
                </div>
                <DialogFooter>
                    <Button variant="outline" onClick={() => setOpen(false)}>
                        Batal
                    </Button>
                    <Button onClick={handleAdd}>Tambah</Button>
                </DialogFooter>
            </DialogContent>
        </Dialog>
    );
};

export default AddDeviceModal;

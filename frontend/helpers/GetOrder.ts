// helpers/GetOrder.ts

export interface Product {
    id: number;
    name: string;
    description: string;
    price_points: number;
    stock: number;
    status: string;
    image_url: string;
    created_at: string;
}

export async function getOrders(): Promise<Product[]> {
    try {
        const res = await fetch("https://backend-phi-murex-10.vercel.app/api/store/items");
        if (!res.ok) {
            throw new Error("Terjadi kesalahan saat mengambil data.");
        }
        const data = await res.json();
        return data.data;  // Mengembalikan data produk
    } catch (error) {
        throw new Error(error.message || "Tidak dapat mengambil data.");
    }
}

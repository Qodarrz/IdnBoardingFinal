// pages/orders.tsx
"use client";

import { useEffect, useState } from "react";
import { getOrders, Product } from "@/helpers/GetOrder";  // Import tipe data Product
import Image from "next/image";
import { Card, CardContent } from "@/components/ui/card";

export default function Orders() {
    const [products, setProducts] = useState<Product[]>([]);  // Tentukan tipe data produk
    const [loading, setLoading] = useState<boolean>(true);
    const [error, setError] = useState<string | null>(null);  // Tipe data error adalah string atau null

    useEffect(() => {
        // Mengambil data dari helper
        const fetchData = async () => {
            try {
                const data = await getOrders();  // Menggunakan helper getOrders
                setProducts(data);  // Menyimpan data produk
            } catch (err: any) {  // Menambahkan penanganan error dengan tipe any
                setError(err.message);  // Menyimpan error jika terjadi
            } finally {
                setLoading(false);
            }
        };

        fetchData();
    }, []);

    if (loading) {
        return (
            <div className="flex justify-center items-center h-screen">
                <div className="spinner">Loading...</div>
            </div>
        );
    }

    if (error) {
        return (
            <div className="text-center text-red-500">
                <p>{error}</p>
            </div>
        );
    }

    return (
        <div className="max-w-7xl mx-auto px-6 py-12">
            <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-6">
                {products.map((item) => (
                    <Card key={item.id} className="overflow-hidden shadow-lg rounded-lg p-4">
                        <div className="relative flex items-center justify-center mb-4">
                            <Image
                                src={item.image_url}
                                alt={item.name}
                                width={150}
                                height={150}
                                className="object-contain w-full"
                            />
                        </div>
                        <CardContent>
                            <h3 className="text-xl font-medium">{item.name}</h3>
                            <div className="flex justify-between items-center mb-2">
                                <p className="text-sm text-muted-foreground">{item.description}</p>
                                <div className="flex gap-2 items-center">
                                    <p className="font-semibold">{item.price_points} Points</p>
                                </div>
                            </div>
                        </CardContent>
                    </Card>
                ))}
            </div>
        </div>
    );
}

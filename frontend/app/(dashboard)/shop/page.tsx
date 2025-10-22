"use client"

import { useEffect, useState } from "react"
import { Card, CardContent } from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { Button } from "@/components/ui/button"
import Image from "next/image"
import { Search } from "lucide-react"
import Swal from "sweetalert2"

type Product = {
    id: number
    name: string
    description: string
    price_points: number
    stock: number
    status: string
    image_url: string
    created_at: string
}

type ApiResponse = {
    status: boolean
    message: string
    data: Product[]
}

export default function RewardStorePage() {
    const [search, setSearch] = useState("")
    const [products, setProducts] = useState<Product[]>([])
    const [loading, setLoading] = useState(true)
    const [page, setPage] = useState(1)
    const perPage = 8

    // Fetch products
    useEffect(() => {
        const fetchProducts = async () => {
            setLoading(true)
            try {
                const res = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/store/items`)
                const data: ApiResponse = await res.json()
                if (data.status) {
                    setProducts(data.data)
                }
            } catch (error) {
                console.error("Gagal fetch data:", error)
            } finally {
                setLoading(false)
            }
        }
        fetchProducts()
    }, [])

    // Function tukar poin
    async function handleExchange(productId: number) {
        try {
            const token = localStorage.getItem("authtoken")
            if (!token) {
                Swal.fire("Error", "Token tidak ditemukan, silakan login ulang.", "error")
                return
            }

            const res = await fetch(
                `${process.env.NEXT_PUBLIC_API_URL}/api/store/orders/${productId}`,
                {
                    method: "POST",
                    headers: {
                        "Content-Type": "application/json",
                        Authorization: `Bearer ${token}`,
                    },
                }
            )

            const result = await res.json()

            if (res.ok && result.status) {
                Swal.fire("Berhasil", result.message || "Produk berhasil ditukar!", "success")
            } else {
                Swal.fire("Gagal", result.message || "Gagal menukar produk.", "error")
            }
        } catch (error: unknown) {
            if (error instanceof Error) {
                Swal.fire("Error", error.message, "error")
            } else {
                Swal.fire("Error", "Terjadi kesalahan tak dikenal", "error")
            }
        }
    }

    // Search & pagination
    const filteredProducts = products.filter((p) =>
        p.name.toLowerCase().includes(search.toLowerCase())
    )
    const totalPages = Math.ceil(filteredProducts.length / perPage)
    const startIndex = (page - 1) * perPage
    const paginatedProducts = filteredProducts.slice(startIndex, startIndex + perPage)

    return (
        <div className="p-6 space-y-6">
            {/* Header */}
            <div className="flex flex-col md:flex-row md:items-center md:justify-between gap-4">
                <h1 className="text-xl font-semibold">Toko Penukaran Poin</h1>
                <div className="relative max-w-sm">
                    <Input
                        placeholder="Cari Barang"
                        value={search}
                        onChange={(e) => {
                            setSearch(e.target.value)
                            setPage(1)
                        }}
                        className="pl-9"
                    />
                    <Search className="absolute left-2 top-2.5 h-4 w-4 text-muted-foreground" />
                </div>
            </div>

            {/* Loading State */}
            {loading && <p className="text-center text-muted-foreground">Memuat data...</p>}

            {/* Products Grid */}
            {!loading && paginatedProducts.length === 0 && (
                <p className="text-center text-muted-foreground">Tidak ada produk ditemukan</p>
            )}

            <div className="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-4 gap-6">
                {paginatedProducts.map((product) => (
                    <Card key={product.id} className="overflow-hidden border-none shadow-none gap-2">
                        <div className="relative flex items-center justify-center transition aspect-square">
                            <Image
                                src={product.image_url}
                                alt={product.name}
                                width={150}
                                height={150}
                                className="object-contain w-full"
                            />
                        </div>
                        <CardContent className="p-0">
                            <h3 className="text-xl font-medium">{product.name}</h3>
                            <div className="flex justify-between items-center mb-2 gap-8">
                                <p className="text-sm text-muted-foreground">{product.description}</p>
                                <div className="flex gap-2 items-center">
                                    <p className="font-semibold">{product.price_points}</p>
                                    <Image
                                        src="/icons/green-point.svg"
                                        alt="Green Point"
                                        width={16}
                                        height={16}
                                    />
                                </div>
                            </div>
                            <Button
                                className="w-full rounded-full"
                                onClick={() => handleExchange(product.id)}
                            >
                                Tukar Point
                            </Button>
                        </CardContent>
                    </Card>
                ))}
            </div>

            {/* Pagination */}
            {!loading && totalPages > 1 && (
                <div className="flex justify-center items-center gap-4">
                    <Button
                        variant="outline"
                        size="sm"
                        disabled={page === 1}
                        onClick={() => setPage((p) => Math.max(1, p - 1))}
                    >
                        Sebelumnya
                    </Button>

                    <div className="flex gap-2">
                        {Array.from({ length: totalPages }).map((_, i) => (
                            <Button
                                key={i}
                                variant={page === i + 1 ? "default" : "outline"}
                                size="sm"
                                onClick={() => setPage(i + 1)}
                            >
                                {i + 1}
                            </Button>
                        ))}
                    </div>

                    <Button
                        variant="outline"
                        size="sm"
                        disabled={page === totalPages}
                        onClick={() => setPage((p) => Math.min(totalPages, p + 1))}
                    >
                        Selanjutnya
                    </Button>
                </div>
            )}
        </div>
    )
}

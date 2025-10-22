"use client"

import Image from "next/image"
import { useEffect, useRef, useState } from "react"
import { Badge } from "../ui/badge"
import { Card, CardContent } from "../ui/card"
import { Button } from "../ui/button"
import { Swiper, SwiperSlide } from "swiper/react"
import { Navigation } from "swiper/modules"
import "swiper/css" // ❗ cukup ini; tak perlu css navigation biar panah default tak muncul
import { ChevronLeft, ChevronRight } from "lucide-react"

const reviews = [
    {
        id: 1,
        text: `"Dengan menggunakan GreenFlow, saya bisa mengurangi emisi karbon saya dan berbuat baik untuk bumi. Aplikasi yang luar biasa!"`,
        name: "Abdul Balmond",
        role: "Aktivis Lingkungan",
        avatar: "/images/profile3.png",
    },
    {
        id: 2,
        text: `"GreenFlow membantu saya untuk lebih sadar akan jejak karbon saya, dan sekarang saya merasa lebih berkontribusi pada lingkungan!"`,
        name: "Yusuf Terizla",
        role: "Aktivis Lingkungan",
        avatar: "/images/profile2.png",
    },
    {
        id: 3,
        text: `"Misi harian di GreenFlow menyenangkan dan bermanfaat. Saya bisa mendapatkan diskon produk daur ulang untuk membantu UMKM!"`,
        name: "Siti Miya",
        role: "Pengguna Setia",
        avatar: "/images/profile.png",
    },
]

export default function ReviewSection() {
    const prevRef = useRef<HTMLButtonElement>(null)
    const nextRef = useRef<HTMLButtonElement>(null)
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    const [swiperInst, setSwiperInst] = useState<any>(null)

    // 🔗 Bind tombol shadcn ke swiper ketika instancenya sudah ada
    useEffect(() => {
        if (!swiperInst) return
        if (!prevRef.current || !nextRef.current) return

        swiperInst.params.navigation.prevEl = prevRef.current
        swiperInst.params.navigation.nextEl = nextRef.current
        swiperInst.navigation.init()
        swiperInst.navigation.update()
    }, [swiperInst])

    return (
        <section className="relative flex flex-col items-center justify-center py-20 lg:px-0 px-4 overflow-hidden" id="testimoni">
            {/* Header */}
            <div className="flex flex-col gap-4 items-center mb-12 text-center">    
                <Badge variant="outline" className="text-gray-700 border-gray-300 px-4 py-2 font-light text-base">
                    <Image src="/icons/stars-black.svg" alt="stars" width={20} height={20} className="inline-block mr-2" />
                    Testimoni Aplikasi GreenFlow
                </Badge>
                <h2 className="text-4xl md:text-5xl font-medium">
                    Apa Kata Mereka Tentang <span className="text-primary">GreenFlow?</span>
                </h2>
                <p className="text-gray-500">Lihat bagaimana GreenFlow telah membantu pengguna untuk lebih sadar akan jejak karbon mereka</p>
            </div>

            {/* Carousel */}
            <div className="relative w-full">
                <Swiper
                    modules={[Navigation]}
                    onSwiper={setSwiperInst}
                    navigation={false}
                    loop={true}
                    loopAdditionalSlides={reviews.length}
                    centeredSlides={true}
                    slidesPerView={1.5}
                    spaceBetween={30}
                    breakpoints={{
                        0: { slidesPerView: 1, spaceBetween: 16, centeredSlides: true },
                        1024: { slidesPerView: 1.5, spaceBetween: 30, centeredSlides: true },
                    }}
                    className="pb-12"
                >
                    {reviews.map((r) => (
                        <SwiperSlide
                            key={r.id}
                            className="!h-auto transition-all duration-300 scale-95 opacity-50 
                            [&.swiper-slide-active]:scale-100 [&.swiper-slide-active]:opacity-100"
                        >
                            <Card className="p-8 shadow-md rounded-2xl">
                                <CardContent className="flex flex-col gap-6">
                                    <Image src="/icons/quote.svg" alt="" width={64} height={64} />
                                    <p className="text-lg text-gray-700">{r.text}</p>
                                    <div className="flex items-center gap-3">
                                        <Image src={r.avatar} alt={r.name} width={48} height={48} className="rounded-full object-cover aspect-square" />
                                        <div className="text-left">
                                            <p className="font-semibold">{r.name}</p>
                                            <p className="text-sm text-gray-500">{r.role}</p>
                                        </div>
                                    </div>
                                </CardContent>
                            </Card>
                        </SwiperSlide>
                    ))}
                </Swiper>

                {/* Fade kiri & kanan */}
                <div className="hidden lg:block pointer-events-none absolute left-0 top-0 h-full w-32 bg-gradient-to-r from-white to-transparent z-10" />
                <div className="hidden lg:block pointer-events-none absolute right-0 top-0 h-full w-32 bg-gradient-to-l from-white to-transparent z-10" />

                {/* Tombol shadcn di bawah, tanpa absolute */}
                <div className="mt-6 flex justify-center gap-4">
                    <Button ref={prevRef} variant="outline" size="icon" className="rounded-full" aria-label="Prev review">
                        <ChevronLeft className="h-5 w-5" />
                    </Button>
                    <Button ref={nextRef} variant="outline" size="icon" className="rounded-full" aria-label="Next review">
                        <ChevronRight className="h-5 w-5" />
                    </Button>
                </div>
            </div>
        </section>
    )
}

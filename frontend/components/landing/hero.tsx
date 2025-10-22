import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"
import Image from "next/image"
import Link from "next/link"
import { ChevronRight } from "lucide-react"

export default function HeroSection() {
    return (
        <section className="pt-18">
            <div className="relative container flex items-center justify-center w-full overflow-hidden mx-auto" id="home">
                {/* Background */}
                <div
                    className="absolute inset-0 bg-cover bg-center mx-4 rounded-2xl"
                    style={{
                        backgroundImage: "url('/images/hero-bg.png')", // ganti dengan gambar kamu
                    }}
                />

                {/* Konten */}
                <div className="relative z-10 max-w-2xl text-center text-white px-8 py-24 pb-64">
                    <Badge
                        variant="outline"
                        className="mb-4 bg-white/10 text-white backdrop-blur-sm border-white/30 px-4 py-2 font-light text-base"
                    >
                        <Image src="/icons/stars.svg" alt="stars" width={24} height={24} className="inline-block mr-2" />
                        Selamat Datang Di GreenFlow
                    </Badge>

                    <h1 className="text-3xl md:text-5xl font-medium mb-4">
                        Mulai Hidup Sejahtera dan Sehat Bersama GreenFlow
                    </h1>

                    <p className="md:text-xl text-gray-200 mb-6">
                        Bergabunglah bersama kami untuk mengurangi jejak karbon dan buat Indonesia lebih hijau.
                    </p>

                    <div className="flex items-center justify-center gap-4">
                        <Button className="bg-white text-primary hover:bg-gray-200 rounded-full" asChild>
                            <Link href="/dashboard">
                                Masuk Dashboard <ChevronRight />
                            </Link>
                        </Button>
                        <Button variant="link" className="text-white underline-offset-4" asChild>
                            <Link href="/register">
                            Daftar
                            </Link>
                        </Button>
                    </div>
                </div>

                {/* Preview Dashboard */}
                <div className="absolute -bottom-28 md:-bottom-22 lg:bottom-0 left-1/2 -translate-x-1/2 z-10 w-full flex justify-end md:justify-center">
                    <picture>
                        {/* Desktop */}
                        <source
                            srcSet="/images/dashboard-preview.png"
                            media="(min-width: 1024px)"
                        />
                        {/* Tablet */}
                        <source
                            srcSet="/images/dashboard-preview-tab.png"
                            media="(min-width: 768px)"
                        />
                        {/* Mobile (default) */}
                        <Image
                            src="/images/dashboard-preview-hp.png"
                            alt="dashboard preview"
                            width={600}
                            height={400}
                            className="drop-shadow-2xl w-[333px] md:w-[555px] lg:w-[900px]"
                        />
                    </picture>
                </div>
            </div >
        </section>
    )
}

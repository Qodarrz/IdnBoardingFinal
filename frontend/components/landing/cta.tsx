import Link from "next/link";
import { Button } from "../ui/button";
import { ChevronRight } from "lucide-react";

export default function CTASection() {
    return (
        <section className="container mx-auto flex justify-center items-center relative rounded-2xl py-12 mb-16">
            <div
                className="absolute inset-0 bg-cover bg-center rounded-2xl mx-4"
                style={{
                    backgroundImage: "url('/images/cta-bg.png')", // ganti dengan gambar kamu
                }}
            />
            <div className="flex flex-col relative z-10 max-w-2xl justify-center items-center py-28 px-12">
                <div className="flex flex-col text-center">
                    <h1 className="text-4xl md:text-5xl font-medium mb-4 text-white">
                        Mulai Perjalanan Hijau Anda Sekarang!
                    </h1>

                    <p className="text-lg md:text-xl text-gray-200 mb-6">
                        Segera daftar di GreenFlow dan mulailah mengurangi jejak karbon Anda. Semua fitur aplikasi ini gratis!
                    </p>
                </div>

                <div className="flex flex-col md:flex-row items-center justify-center gap-4">
                    <Button className="bg-white text-black hover:bg-gray-200" asChild>
                        <Link href="/dashboard">
                            Masuk Dashboard <ChevronRight />
                        </Link>
                    </Button>
                    <Button variant="link" className="text-white underline-offset-4">
                        Lihat Misi Anda
                    </Button>
                </div>
            </div>
        </section>
    )
}
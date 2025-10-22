import Image from "next/image";
import Link from "next/link";
import { Badge } from "../ui/badge";
import { Button } from "../ui/button";
import { ChevronRight } from "lucide-react";
import { Card, CardContent, CardHeader, CardTitle } from "../ui/card";

export default function AboutSection() {
    return (
        <section className="container mx-auto py-16 px-4" id="about">
            <div className="flex flex-col justify-between gap-8 md:gap-12 lg:flex-row">
                <div className="flex flex-col justify-between gap-8 md:gap-16">
                    <div className="flex flex-col gap-2">
                        <Badge
                            variant="outline"
                            className="mb-4 text-gray-700 border-gray-700 px-4 py-2 font-light text-base"
                        >
                            <Image src="/icons/stars-black.svg" alt="stars" width={24} height={24} className="inline-block mr-2" />
                            Perkenalan Aplikasi GreenFlow
                        </Badge>
                        <h2 className="text-3xl md:text-6xl font-medium">Apa Itu <span className="text-primary">GreenFlow?</span></h2>
                    </div>
                    <Button className="rounded-full px-[32px] text-base w-fit" asChild>
                        <Link href="/dashboard">
                            Mulai dan Masuk Dashboard<ChevronRight />
                        </Link>
                    </Button>
                </div>
                <p className="text-3xl md:text-4xl lg:text-5xl text-gray-700 lg:w-1/2 md:max-w-2xl lg:items-center font-medium">
                    GreenFlow adalah aplikasi untuk memantau jejak karbon dan <span className="opacity-70">berpartisipasi dalam program pengurangan karbon.</span>
                </p>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-6 lg:grid-cols-12 gap-8 py-12">
                <Card className="md:col-span-6 flex w-full md:flex-row items-center rounded-2xl p-0 bg-gradient-to-tr from-[#C8EAD8] to-[#DFF1EF] lg:max-w-[600px]">
                    <CardContent className="px-8 pt-8 md:py-8 flex justify-between flex-col h-full">
                        <CardTitle className="text-3xl md:text-4xl font-medium mb-4">
                            Menjadikan Bumi Yang Lebih Asri dan Baik
                        </CardTitle>
                        <p className="text-sm md:text-lg">
                            Dengan berpartisipasi dalam GreenFlow, Anda berkontribusi pada penciptaan lingkungan yang lebih sehat dan asri
                        </p>
                    </CardContent>
                    <div className="w-full flex justify-center md:w-1/3 h-52 md:h-auto">
                        <Image
                            src="/images/flower.png" // ganti dengan gambar kamu
                            alt="feature"
                            className="object-contain h-full max-w-[320px] md:w-48"
                            width={500}
                            height={500}
                        />
                    </div>
                </Card>

                <Card className="rounded-2xl md:col-span-3 aspect-square bg-gradient-to-tr pr-16 from-[#27423B] to-primary text-white justify-between">
                    <CardHeader>
                        <CardTitle className="text-3xl font-medium">
                            Mengurangi Jejak Karbon Anda
                        </CardTitle>
                    </CardHeader>
                    <CardContent>
                        <p className="opacity-70">
                            Lacak dan kurangi emisi karbon dalam kegiatan sehari-hari dengan lebih mudah.
                        </p>
                    </CardContent>
                </Card>

                <Card className="rounded-2xl md:col-span-3 aspect-square bg-gradient-to-tr pr-16 from-[#27423B] to-primary text-white justify-between">
                    <CardHeader>
                        <CardTitle className="text-3xl font-medium">
                            Mendukung Konservasi Alam 
                        </CardTitle>
                    </CardHeader>
                    <CardContent>
                        <p className="opacity-70">
                            Anda mendukung konservasi alam dan mempromosikan keberlanjutan.
                        </p>
                    </CardContent>
                </Card>
            </div>
        </section>
    );
}
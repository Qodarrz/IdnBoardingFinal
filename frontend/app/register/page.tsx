"use client"

import { useEffect, useState } from "react"
import Image from "next/image"
import { motion, AnimatePresence } from "framer-motion"
import { RegisterForm } from "@/components/shared/register-form"

const slides = [
    {
        bg: "/images/auth-bg.png",
        logo: "/icons/bg-white-logo-only.svg",
        title: "Mulai Hidup Sehat dan Sejahtera Bersama Kami",
        desc: "Mulai Hidup Sehat dan Sejahtera Bersama Kami",
    },
    {
        bg: "/images/auth-bg2.png",
        logo: "/icons/bg-white-logo-only.svg",
        title: "Aksi Nyata untuk Bumi Hijau",
        desc: "Matikan lampu saat tidak digunakan dan kurangi jejak karbon mulai dari rumah.",
    },
    {
        bg: "/images/auth-bg3.png",
        logo: "/icons/bg-white-logo-only.svg",
        title: "Bersama Kurangi Emisi, Raih Masa Depan Cerah",
        desc: "Setiap aksi kecilmu adalah investasi bagi generasi mendatang.",
    },
]

export default function RegisterPage() {
    const [current, setCurrent] = useState(0)

    useEffect(() => {
        const interval = setInterval(() => {
            setCurrent((prev) => (prev + 1) % slides.length)
        }, 5000)
        return () => clearInterval(interval)
    }, [])

    return (
        <div className="grid min-h-svh lg:grid-cols-12">
            {/* Left Side - Carousel */}
            <div className="relative hidden lg:block lg:col-span-4 overflow-hidden">
                <AnimatePresence mode="sync">
                    <motion.div
                        key={current}
                        initial={{ opacity: 0 }}
                        animate={{ opacity: 1 }}
                        exit={{ opacity: 0 }}
                        transition={{ duration: 1 }}
                        className="absolute inset-0"
                    >
                        {/* Background Image */}
                        <Image
                            src={slides[current].bg}
                            alt="Background"
                            fill
                            className="object-cover dark:brightness-[0.4]"
                            priority
                        />

                        {/* Overlay Content */}
                        <div className="relative z-10 flex flex-col items-center justify-end text-center px-8 h-full text-white pb-24">
                            <motion.div
                                key={current + "-content"}
                                initial={{ y: 20, opacity: 0 }}
                                animate={{ y: 0, opacity: 1 }}
                                exit={{ opacity: 0 }}
                                transition={{ duration: 0.8 }}
                                className="flex flex-col items-center gap-4 max-w-xs"
                            >
                                <Image
                                    src={slides[current].logo}
                                    alt="Logo"
                                    width={48}
                                    height={48}
                                />
                                <h2 className="text-2xl font-semibold leading-snug max-w-xs">
                                    {slides[current].title}
                                </h2>
                                <p className="text-sm text-gray-200 max-w-sm">
                                    {slides[current].desc}
                                </p>
                            </motion.div>

                            {/* Dot Indicators */}
                            <div className="flex gap-2 mt-6">
                                {slides.map((_, i) => (
                                    <button
                                        key={i}
                                        onClick={() => setCurrent(i)}
                                        className={`h-2.5 w-2.5 rounded-full transition-colors duration-300 ${current === i
                                                ? "bg-primary"
                                                : "bg-white"
                                            }`}
                                    />
                                ))}
                            </div>
                        </div>
                    </motion.div>
                </AnimatePresence>
            </div>

            {/* Right Side - Login */}
            <div className="flex flex-col gap-4 p-6 md:p-10 lg:col-span-8">
                <div className="flex justify-center gap-2 md:justify-start">
                    <Image
                        src="/icons/main-logo.svg"
                        alt="Main Logo"
                        width={128}
                        height={128}
                    />
                </div>
                <div className="flex flex-1 items-center justify-center">
                    <div className="w-full max-w-sm">
                        <RegisterForm />
                    </div>
                </div>
            </div>
        </div>
    )
}

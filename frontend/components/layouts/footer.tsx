"use client"

import Image from "next/image";
import { ArrowUp } from "lucide-react";
import Link from "next/link";

export default function Footer() {
    const currentYear = new Date().getFullYear();

    const scrollToTop = () => {
        window.scrollTo({ top: 0, behavior: "smooth" });
    };

    return (
        <footer className="bg-[#3E3E3E] text-white py-8 overflow-hidden px-4">
            {/* Logo */}
            <div className="container mx-auto md:px-0 flex flex-col md:flex-row md:justify-between md:items-center gap-6">
                <div className="flex items-center">
                    <Image
                        src="/icons/white-logo.svg"
                        alt="Logo"
                        width={200}
                        height={200}
                        className="object-contain"
                    />
                </div>
            </div>

            {/* Info Section */}
            <div className="container mx-auto md:px-0 py-8 grid grid-cols-1 lg:grid-cols-2 gap-12 lg:gap-64">
                <div className="flex flex-col justify-between">
                    <h1 className="text-5xl md:text-5xl font-semibold">
                        Langkah Kecil untuk <span className="opacity-70">Bumi yang Lebih Besar!</span>
                    </h1>

                    <p className="text-xs md:text-sm hidden md:flex font-light opacity-70 mt-6">
                        &copy; {currentYear} GreenFlow all right reserved
                    </p>
                </div>

                <div className="grid grid-cols-1 md:grid-cols-2 gap-8 text-sm">
                    <div className="flex flex-col gap-2 text-lg">
                        <h2 className="opacity-60">Location</h2>
                        <p>Jl. Raya Puncak, Gadog No.619, Pandansari, Ciawi District, Bogor Regency</p>
                    </div>
                    <div className="flex flex-col gap-2 text-lg">
                        <h2 className="opacity-60">Email</h2>
                        <p>renjieprass@gmail.com</p>
                    </div>
                    <div className="flex flex-col gap-2 text-lg">
                        <h2 className="opacity-60">Contact</h2>
                        <p>+62 857 7025 3105</p>
                    </div>
                    <div className="flex flex-col gap-2 text-lg">
                        <h2 className="opacity-60">Social Media</h2>
                        <div className="flex gap-2">
                            <Link href="https://facebook.com/jie.env" target="_blank" rel="noopener noreferrer">
                                <Image src="/icons/facebook.svg" alt="Facebook" width={32} height={32} />
                            </Link>
                            <Link href="https://instagram.com/jie.env" target="_blank" rel="noopener noreferrer">
                                <Image src="/icons/instagram.svg" alt="Instagram" width={32} height={32} />
                            </Link>
                            <Link href="https://twitter.com/jie.env" target="_blank" rel="noopener noreferrer">
                                <Image src="/icons/twitter.svg" alt="Twitter" width={32} height={32} />
                            </Link>
                            <Link href="https://youtube.com/jie.env" target="_blank" rel="noopener noreferrer">
                                <Image src="/icons/youtube.svg" alt="YouTube" width={32} height={32} />
                            </Link>
                        </div>
                    </div>
                </div>
            </div>

            {/* Big Logo Text */}
            <div className="container mx-auto md:px-0">
                <div className="flex justify-center items-center">
                    <h1 className="text-[77px] md:text-[160px] lg:text-[266px] cursor-default text-text mb-6 leading-none font-semibold text-center">
                        GreenFlow
                    </h1>
                </div>

                <button
                    onClick={scrollToTop}
                    className="w-full rounded-xl border border-gray-700 px-4 md:px-8 py-4 flex items-center justify-between gap-4 bg-[#4C4C4C] text-sm hover:bg-neutral-700 transition cursor-pointer"
                >
                    <span className="text-left text-white">
                        Kembali Ke Halaman<br /> Paling Atas
                    </span>
                    <div className="bg-white text-black rounded-full p-2">
                        <ArrowUp className="w-4 h-4" />
                    </div>
                </button>
            </div>

            <p className="text-xs text-center text-subtle md:hidden mt-4">
                &copy; {currentYear} GreenFlow semua hak dilindungi.
            </p>
        </footer>
    );
}

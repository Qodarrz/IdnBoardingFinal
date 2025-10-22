"use client"

import Image from "next/image"
import { Badge } from "@/components/ui/badge"
import {
    Accordion,
    AccordionContent,
    AccordionItem,
    AccordionTrigger,
} from "@/components/ui/accordion"

export default function FAQSection() {
    const faqs = [
        {
            question: "Apa itu GreenFlow?",
            answer: "GreenFlow adalah aplikasi yang membantu Anda mengurangi jejak karbon dan berpartisipasi dalam program ramah lingkungan.",
        },
        {
            question: "Bagaimana cara menggunakan GreenFlow?",
            answer: "Cukup daftarkan akun, ikuti misi, dan lacak jejak karbon Anda.",
        },
        {
            question: "Apakah GreenFlow gratis?",
            answer: "Ya, GreenFlow sepenuhnya gratis dan dapat digunakan oleh siapa saja.",
        },
        {
            question: "Apa manfaat mengikuti misi di GreenFlow?",
            answer: "Dengan mengikuti misi harian, Anda mengurangi jejak karbon dan mendapatkan koin yang bisa ditukar dengan voucher diskon produk ramah lingkungan.",
        },
        {
            question: "Bagaimana saya bisa mengonversi koin?",
            answer: "Setelah mengumpulkan cukup koin, Anda bisa menukarkannya dengan voucher diskon produk dari UMKM Indonesia yang ramah lingkungan.",
        },
    ]

    return (
        <section className="container mx-auto py-20 px-4" id="faq">
            <div className="flex flex-col lg:flex-row justify-between gap-12">
                {/* Left Content */}
                <div className="flex flex-col gap-4 lg:w-1/2">
                    <Badge
                        variant="outline"
                        className="w-fit text-gray-700 border-gray-700 px-4 py-2 font-light text-base"
                    >
                        <Image
                            src="/icons/stars-black.svg"
                            alt="stars"
                            width={20}
                            height={20}
                            className="inline-block mr-2"
                        />
                        Fitur AI Terbaru dan Terbaik
                    </Badge>
                    <h2 className="text-4xl md:text-6xl font-medium leading-tight">
                        Cari Jawaban Dari <span className="text-primary">Pertanyaan</span> Anda 🤔
                    </h2>
                    <p className="text-gray-500 text-lg">
                        Temukan jawaban atas berbagai pertanyaan seputar penggunaan GreenFlow
                    </p>
                </div>

                {/* Right FAQ Accordion */}
                <div className="lg:flex-1">
                    <Accordion type="single" collapsible className="w-full space-y-8">
                        {faqs.map((faq, index) => (
                            <AccordionItem
                                key={index}
                                value={`item-${index}`}
                                className="border border-primary rounded-xl px-6"
                            >
                                <AccordionTrigger className="text-xl text-primary font-medium hover:no-underline">
                                    {faq.question}
                                </AccordionTrigger>
                                <AccordionContent className="text-gray-600 text-base pb-4">
                                    {faq.answer}
                                </AccordionContent>
                            </AccordionItem>
                        ))}
                    </Accordion>
                </div>
            </div>
        </section>
    )
}

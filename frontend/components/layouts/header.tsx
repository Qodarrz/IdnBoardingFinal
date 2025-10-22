"use client";

import { useEffect, useRef, useState } from "react";
import Link from "next/link";
import Image from "next/image";
import { Button } from "../ui/button";
import { motion, useAnimation, AnimatePresence } from "framer-motion";
import clsx from "clsx";
import { IconBubble } from "@tabler/icons-react"
import { Menu, X, Home, Info, Wrench, HelpCircle } from "lucide-react";

/* ========= NAVIGATION ITEMS ========= */
const navItems = [
    { id: "home", label: "Beranda", icon: Home },
    { id: "about", label: "Tentang Kami", icon: Info },
    { id: "feature", label: "Layanan Kami", icon: Wrench },
    { id: "testimoni", label: "Testimoni", icon: IconBubble },
    { id: "faq", label: "FAQ", icon: HelpCircle },
] as const;

type NavItem = (typeof navItems)[number];

export default function Header() {
    const [activeId, setActiveId] = useState<string>("home");
    const [menuOpen, setMenuOpen] = useState(false);
    const controls = useAnimation();
    const lastScrollY = useRef<number>(0);

    /* ---------------- SCROLL-SPY ---------------- */
    useEffect(() => {
        const observer = new IntersectionObserver(
            (entries) => {
                entries.forEach((entry) => {
                    if (entry.isIntersecting) setActiveId(entry.target.id);
                });
            },
            { rootMargin: "-50% 0px -50% 0px", threshold: 0 }
        );

        navItems.forEach(({ id }) => {
            const el = document.getElementById(id);
            if (el) observer.observe(el);
        });

        return () => observer.disconnect();
    }, []);

    /* ------------- HIDE/SHOW ON SCROLL ------------- */
    useEffect(() => {
        const handleScroll = () => {
            const currentY = window.scrollY;

            if (currentY > lastScrollY.current && currentY > 80) {
                controls.start({ y: "-100%" });
            } else {
                controls.start({ y: 0 });
            }

            lastScrollY.current = currentY;
        };

        window.addEventListener("scroll", handleScroll, { passive: true });
        return () => window.removeEventListener("scroll", handleScroll);
    }, [controls]);

    /* ------------- LOCK SCROLL WHEN SIDEBAR OPEN ------------- */
    useEffect(() => {
        if (menuOpen) document.body.style.overflow = "hidden";
        else document.body.style.overflow = "";
        return () => {
            document.body.style.overflow = "";
        };
    }, [menuOpen]);

    /* ------------- Helpers ------------- */
    const linkClass = (id: string) =>
        clsx(
            "transition-colors",
            activeId === id
                ? "text-primary underline font-medium"
                : "text-muted-foreground hover:text-primary"
        );

    const closeMenu = () => setMenuOpen(false);

    return (
        <>
            {/* ======= TOP NAVBAR ======= */}
            <motion.header
                animate={controls}
                transition={{ type: "spring", stiffness: 300, damping: 30 }}
                className="fixed top-0 left-0 right-0 z-50"
            >
                <div className="flex justify-between items-center bg-[#f9f9f9] py-4 px-6 lg:px-4 lg:mx-auto lg:container">
                    {/* Logo */}
                    <Link
                        href="#home"
                        className="text-2xl lg:text-3xl font-bold text-primary"
                    >
                        <Image
                            src="/icons/main-logo.svg"
                            alt="Logo"
                            width={200}
                            height={200}
                        />
                    </Link>

                    {/* Desktop Nav */}
                    <nav className="hidden lg:flex gap-8">
                        {navItems.map(({ id, label }) => (
                            <Link key={id} href={`#${id}`} className={linkClass(id)}>
                                {label}
                            </Link>
                        ))}
                    </nav>

                    {/* Dashboard Button & Hamburger */}
                    <div className="flex gap-4 items-center">
                        <Button className="rounded-full px-6 text-base hidden lg:flex" asChild>
                            <Link href="/dashboard">Dashboard</Link>
                        </Button>

                        {/* Hamburger (Mobile) */}
                        <button
                            className="lg:hidden text-text"
                            onClick={() => setMenuOpen(true)}
                            aria-label="Open menu"
                        >
                            <Menu size={24} />
                        </button>
                    </div>
                </div>
            </motion.header>

            {/* ======= MOBILE SIDEBAR ======= */}
            <AnimatePresence>
                {menuOpen && (
                    <motion.div
                        key="overlay"
                        className="fixed inset-0 z-60 bg-black/50 backdrop-blur-sm lg:hidden"
                        initial={{ opacity: 0 }}
                        animate={{ opacity: 1 }}
                        exit={{ opacity: 0 }}
                        onClick={closeMenu}
                    >
                        <motion.aside
                            className="absolute right-0 top-0 h-full w-3/4 max-w-xs bg-[#f9f9f9] dark:bg-dark-card-muted p-6 flex flex-col shadow-2xl"
                            initial={{ x: "100%" }}
                            animate={{ x: 0 }}
                            exit={{ x: "100%" }}
                            transition={{ type: "spring", stiffness: 300, damping: 30 }}
                            onClick={(e) => e.stopPropagation()}
                        >
                            <div className="flex justify-between items-center border-b border-text-sub mb-6 pb-2">
                                <Link
                                    href="#home"
                                    className="text-2xl font-bold text-primary"
                                    onClick={closeMenu}
                                >
                                    GreenFlow
                                </Link>
                                <button
                                    className="text-text"
                                    aria-label="Close menu"
                                    onClick={closeMenu}
                                >
                                    <X className="w-6 h-6" />
                                </button>
                            </div>

                            <div className="flex flex-col gap-6">
                                <div className="flex flex-col gap-3">
                                    <span className="text-text text-sm font-semibold">Menu</span>
                                    {/* Nav links with icons */}
                                    <nav className="flex flex-col gap-4">
                                        {navItems.map(({ id, label, icon: Icon }: NavItem) => (
                                            <Link
                                                key={id}
                                                href={`#${id}`}
                                                className={clsx(linkClass(id), "flex items-center gap-3")}
                                                onClick={closeMenu}
                                            >
                                                <Icon size={18} className="shrink-0" />
                                                {label}
                                            </Link>
                                        ))}
                                    </nav>
                                </div>

                                <div className="flex flex-col gap-3">
                                    <span className="text-text text-sm font-semibold">Akses Cepat</span>
                                    <Button className="rounded-full px-6 text-base" asChild>
                                        <Link href="/dashboard">Dashboard</Link>
                                    </Button>
                                </div>
                            </div>
                        </motion.aside>
                    </motion.div>
                )}
            </AnimatePresence>
        </>
    );
}

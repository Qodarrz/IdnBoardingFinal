"use client"

import { cn } from "@/lib/utils"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import Image from "next/image"
import Link from "next/link"
import { useState } from "react"
import Swal from "sweetalert2"

export function LoginForm({
    className,
    ...props
}: React.ComponentProps<"form">) {
    const [loading, setLoading] = useState(false)

    const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
        e.preventDefault()
        setLoading(true)

        const formData = new FormData(e.currentTarget)
        const email = formData.get("email")
        const password = formData.get("password")

        try {
            const res = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/auth/login`, {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({ email, password }),
            })

            const data = await res.json()
            console.log("Login Response:", data)

            if (res.ok && data) {
                // set cookie
                
                const res = await fetch('/api/cookies', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ token: data.data.token }),
                })

                if (!res.ok) {
                    alert('Login gagal, cookies gagal diset')
                    console.log('error cookie login')
                    console.log(res)
                }

                // storing to localStorage
                localStorage.setItem('islogin', 'done')
                localStorage.setItem('authtoken', data.data.token)

                Swal.fire({
                    toast: true,
                    position: "top-end",
                    icon: "success",
                    title: "Login berhasil!",
                    showConfirmButton: false,
                    timer: 2000,
                    timerProgressBar: true,
                }).then(() => window.location.href = '/dashboard')



            } else {
                Swal.fire({
                    icon: "error",
                    title: "Login gagal",
                    text: data.message || "Periksa kembali email dan password Anda.",
                })
            }
        } catch (error) {
            console.error(error)
            Swal.fire({
                icon: "error",
                title: "Terjadi Kesalahan",
                text: "Tidak dapat terhubung ke server.",
            })
        } finally {
            setLoading(false)
        }
    }

    return (
        <form
            onSubmit={handleSubmit}
            className={cn("flex flex-col gap-6", className)}
            {...props}
        >
            <div className="flex flex-col items-center gap-2 text-center">
                <h1 className="text-2xl font-bold">
                    Selamat Datang Kembali di <span className="text-primary">GreenFlow</span>👋
                </h1>
                <p className="text-muted-foreground text-sm text-balance">
                    Mulai Hidup Sehat dan Sejahtera Bersama Kami
                </p>
            </div>
            <div className="grid gap-6">
                <div className="grid gap-3">
                    <Label htmlFor="email">Email</Label>
                    <Input id="email" name="email" type="email" placeholder="example@email.com" required />
                </div>
                <div className="grid gap-3">
                    <div className="flex items-center">
                        <Label htmlFor="password">Password</Label>
                        <Link
                            href="#"
                            className="ml-auto text-sm underline-offset-4 hover:underline"
                        >
                            Forgot your password?
                        </Link>
                    </div>
                    <Input
                        id="password"
                        name="password"
                        type="password"
                        placeholder="At least 8 characters"
                        required
                    />
                </div>
                <Button type="submit" className="w-full" disabled={loading}>
                    {loading ? "Loading..." : "Sign In"}
                </Button>
                <div className="after:border-border relative text-center text-sm after:absolute after:inset-0 after:top-1/2 after:z-0 after:flex after:items-center after:border-t">
                    <span className="bg-background text-muted-foreground relative z-10 px-2">
                        Or continue with
                    </span>
                </div>
                <Button variant="outline" className="w-full">
                    <Image src="/icons/google.svg" alt="Google Icon" width={16} height={16} />
                    Sign In With Google
                </Button>
            </div>
            <div className="text-center text-sm">
                Don&apos;t have an account?{" "}
                <a href="/register" className="underline underline-offset-4">
                    Sign up
                </a>
            </div>
        </form>
    )
}

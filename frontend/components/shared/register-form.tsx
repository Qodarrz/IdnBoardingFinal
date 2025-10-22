"use client"

import { useState } from "react"
import { cn } from "@/lib/utils"
import Swal from "sweetalert2"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import Image from "next/image"
import Link from "next/link"

export function RegisterForm({
    className,
    ...props
}: React.ComponentProps<"form">) {
    const [loading, setLoading] = useState(false)
    const [error, setError] = useState<string | null>(null)
    const [success, setSuccess] = useState<string | null>(null)

    async function handleSubmit(e: React.FormEvent<HTMLFormElement>) {
        e.preventDefault()
        setError(null)
        setSuccess(null)
        setLoading(true)

        const formData = new FormData(e.currentTarget)
        const data = {
            username: formData.get("username"),
            email: formData.get("email"),
            password: formData.get("password"),
        }

        try {
            const res = await fetch(process.env.NEXT_PUBLIC_API_URL + "/api/auth/register", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify(data),
            })

            if (!res.ok) {
                throw new Error("Failed to register")
            }

            const result = await res.json()

            // ✅ Alert sukses
            Swal.fire({
                toast: true,
                position: "top-end",
                icon: "success",
                title: "Register berhasil! Silakan login.",
                showConfirmButton: false,
                timer: 3000,
                timerProgressBar: true,
            }).then(() => {
                window.location.href = '/login'
            })

            setSuccess("Register berhasil! Silakan login.")
            console.log("Register success:", result)

        } catch (err: unknown) {
            let message = "Terjadi kesalahan"
            if (err instanceof Error) {
                message = err.message
            }

            // ❌ Alert error
            Swal.fire({
                toast: true,
                position: "top-end",
                icon: "error",
                title: message,
                showConfirmButton: false,
                timer: 3000,
                timerProgressBar: true,
            })

            setError(message)
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
                    <Label htmlFor="username">Username</Label>
                    <Input name="username" id="username" type="text" placeholder="your name" required />
                </div>

                <div className="grid gap-3">
                    <Label htmlFor="email">Email</Label>
                    <Input name="email" id="email" type="email" placeholder="example@email.com" required />
                </div>

                <div className="grid gap-3">
                    <Label htmlFor="password">Password</Label>
                    <Input name="password" id="password" type="password" placeholder="At least 8 characters" required />
                </div>

                <div className="grid gap-3">
                    <Label htmlFor="confirm_password">Confirm Password</Label>
                    <Input name="confirm_password" id="confirm_password" type="password" placeholder="Confirm your password" required />
                </div>

                {error && <p className="text-red-500 text-sm text-center">{error}</p>}
                {success && <p className="text-green-500 text-sm text-center">{success}</p>}

                <Button type="submit" className="w-full" disabled={loading}>
                    {loading ? "Registering..." : "Register"}
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
                Have an account?{" "}
                <Link href="/login" className="underline underline-offset-4">
                    Sign In
                </Link>
            </div>
        </form>
    )
}

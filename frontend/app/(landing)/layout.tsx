import Header from "@/components/layouts/header"
import Footer from "@/components/layouts/footer"

export default function LandingLayout({ children }: { children: React.ReactNode }) {
    return (
        <>
            <Header />
            <main className="bg-[#F9F9F9]">{children}</main>
            <Footer />
        </>
    )
}

import AboutSection from "@/components/landing/about";
import CTASection from "@/components/landing/cta";
import FAQSection from "@/components/landing/faq";
import FeatureSection from "@/components/landing/feature";
import HeroSection from "@/components/landing/hero";
import ReviewSection from "@/components/landing/review";

export default function Home() {
  return (
    <div className="">
      <HeroSection />
      <AboutSection />
      <FeatureSection />
      <ReviewSection />
      <FAQSection />
      <CTASection />
    </div>
  );
}

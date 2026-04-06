import { Navbar } from "@/components/layouts/Navbar";
import { PageContainer } from "@/components/layouts/PageContainer";
import { BannerSection } from "./BannerSection";
import { FeaturesSection } from "./FeaturesSection";

export function HomePage() {
  return (
    <>
      <Navbar />
      <PageContainer withSidebar={false}>
        {/* Banner Section */}
        <BannerSection />
        {/* Features Section */}
        <FeaturesSection />
      </PageContainer>
    </>
  );
}

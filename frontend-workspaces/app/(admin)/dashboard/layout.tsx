import { PageContainer } from "@/components/layouts";
import { ReactNode } from "react";

interface PageProps {
  children: ReactNode;
}

export default function Page({ children }: PageProps) {
  return <PageContainer>{children}</PageContainer>;
}

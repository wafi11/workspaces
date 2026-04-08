"use client"
import { PageContainer } from "@/components/layouts";
import { navItemsUsers } from "@/data";
import { ReactNode } from "react";

interface PageProps {
  children: ReactNode;
}

export default function Page({ children }: PageProps) {
  return <PageContainer data={navItemsUsers}>{children}</PageContainer>;
}

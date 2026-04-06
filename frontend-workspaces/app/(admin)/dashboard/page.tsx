import { getMetadata } from "@/data/DataMetadata";
import { AdminDashboardPage } from "@/features/pages/dashboard/Dashboard";
import { Metadata } from "next";

export const metadata: Metadata = getMetadata("dashboard");

export default function Page() {
  return <AdminDashboardPage />;
}

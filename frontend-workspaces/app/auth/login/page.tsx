import { getMetadata } from "@/data/DataMetadata";
import { LoginPage } from "@/features/pages/auth";
import { Metadata } from "next";

export const metadata: Metadata = getMetadata("login");

export default function Page() {
  return <LoginPage />;
}

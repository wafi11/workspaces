import { getMetadata } from "@/data/DataMetadata";
import { RegisterPage } from "@/features/pages/auth";
import { Metadata } from "next";

export const metadata: Metadata = getMetadata("register");

export default function Page() {
  return <RegisterPage />;
}

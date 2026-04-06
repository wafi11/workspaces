import { getMetadata } from "@/data/DataMetadata";
import { HomePage } from "@/features/pages/home/Home";
import { Metadata } from "next";

export const metadata: Metadata = getMetadata("home");

export default function Home() {
  return <HomePage />;
}

import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  images: {
    remotePatterns: [
      {
        protocol : "http",
        port: "9000",
        hostname: "192.168.1.10",
      },
    ],
  },
};

export default nextConfig;

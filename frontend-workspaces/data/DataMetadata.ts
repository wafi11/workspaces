export type MetadataType = {
  title: string;
  description: string;
  url: string;
};

// Menggunakan Record agar mapping antara nama page dan datanya aman (Type-safe)
export const METADATA_CONFIG: Record<string, MetadataType> = {
  home: {
    title: "KubeSpace — Code, Stage, and Launch",
    description:
      "Platform workspace lengkap untuk membangun aplikasi dari development hingga produksi tanpa pusing infrastruktur.",
    url: "/",
  },
  login: {
    title: "Sign In | KubeSpace",
    description:
      "Masuk ke akun KubeSpace Anda untuk mengelola workspace dan deployment.",
    url: "/login",
  },
  register: {
    title: "Create an Account | KubeSpace",
    description:
      "Mulai bangun aplikasi Anda di cloud secara gratis dengan KubeSpace.",
    url: "/register",
  },
  features: {
    title: "Capabilities & Workflow | KubeSpace",
    description:
      "Jelajahi fitur unggulan: Isolated Workspaces, Instant Subdomains, dan Real-time Monitoring.",
    url: "/features",
  },
  dashboard: {
    title: "Console | KubeSpace",
    description:
      "Pantau project, log, dan resource usage aplikasi Anda secara real-time.",
    url: "/dashboard",
  },
  docs: {
    title: "Documentation | KubeSpace",
    description:
      "Panduan lengkap cara deploy aplikasi, integrasi database, dan manajemen environment.",
    url: "/docs",
  },
};

/**
 * Fungsi helper untuk generate metadata Next.js
 */
export function getMetadata(page: keyof typeof METADATA_CONFIG) {
  const data = METADATA_CONFIG[page];

  return {
    title: data.title,
    description: data.description,
    alternates: {
      canonical: data.url,
    },
    openGraph: {
      title: data.title,
      description: data.description,
      url: data.url,
      type: "website",
    },
    // Tambahkan konfigurasi default lainnya di sini
    metadataBase: new URL("https://kubespace.id"),
  };
}

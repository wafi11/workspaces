import { DOCS_URL, REGISTER_URL } from "@/constants";
import Link from "next/link";

export function BannerSection() {
  return (
    <section className="relative z-10 max-w-6xl mx-auto px-8 pt-40 pb-20 flex flex-col items-center text-center">
      {/* Status Badge */}
      <div className="inline-flex items-center gap-2 font-mono text-xs text-primary border border-border px-3 py-1.5 mb-10 bg-secondary tracking-widest">
        <span className="w-1.5 h-1.5 rounded-full bg-green-700 animate-pulse" />
        Active deployments: 1,024
      </div>

      {/* Hero Title */}
      <h1 className="flex flex-col text-[clamp(2.5rem,7vw,5.5rem)] font-semibold leading-[1.05] tracking-tight mb-6 text-foreground">
        <span>Code, Stage, and</span>
        <span className="font-mono bg-primary text-primary-foreground px-4 py-1 my-1 self-center">
          Launch
        </span>
        <span>without the infra fuss.</span>
      </h1>

      <p className="text-lg text-muted-foreground leading-relaxed max-w-lg mb-10 font-light">
        Platform workspace lengkap untuk membangun aplikasi dari baris kode
        pertama hingga skala produksi — tanpa perlu menyentuh konfigurasi
        infrastruktur yang rumit.
      </p>

      <div className="flex gap-4 mb-16 flex-wrap justify-center">
        <Link
          href={REGISTER_URL}
          className="font-mono font-bold text-sm bg-primary text-primary-foreground hover:opacity-90 px-7 py-3 transition-all hover:-translate-y-px flex items-center gap-2 group"
        >
          Start for free
          <span className="transition-transform group-hover:translate-x-1">
            →
          </span>
        </Link>
        <Link
          href={DOCS_URL}
          className="font-mono text-sm text-foreground border-2 border-foreground hover:bg-secondary px-7 py-3 transition-all"
        >
          Explore the docs
        </Link>
      </div>

      {/* Terminal UI */}
      <div className="w-full max-w-2xl border border-border bg-card text-left shadow-2xl">
        <div className="flex items-center gap-1.5 px-4 py-2.5 border-b border-border bg-secondary">
          <span className="w-2.5 h-2.5 rounded-full bg-[#ff5f57]" />
          <span className="w-2.5 h-2.5 rounded-full bg-[#ffbd2e]" />
          <span className="w-2.5 h-2.5 rounded-full bg-[#28ca41]" />
          <span className="font-mono text-xs text-muted-foreground ml-2">
            workspace --status --all
          </span>
        </div>
        <pre className="font-mono text-[11px] md:text-xs leading-loose px-6 py-5 text-primary overflow-x-auto">
          {`PROJECT           ENVIRONMENT    STATUS      UPTIME    ENDPOINT
main-app-v3       Production     Healthy     14d       app.vazz.id
main-app-staging  Staging        Running     5h        stage.vazz.id
dev-api-wafi      Development    Active      12m       dev-api.vazz.id
db-postgresql     Database       Connected   14d       internal:5432`}
        </pre>
      </div>
    </section>
  );
}

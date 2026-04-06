import { LOGIN_URL, REGISTER_URL } from "@/constants";
import Link from "next/link";

export function Navbar() {
  return (
    <nav className="fixed top-0 left-0 right-0 z-50 border-b border-border bg-background/85 backdrop-blur-md">
      <div className="mx-auto px-8 h-[60px] flex items-center gap-8">
        <div className="flex items-center gap-2">
          <span className="text-primary text-xl">⬡</span>
          <span className="font-mono font-bold text-sm tracking-tight text-foreground">
            KubeSpace
          </span>
        </div>
        <div className="flex gap-6 ml-auto">
          {["Features", "Pricing", "Docs"].map((l) => (
            <Link
              key={l}
              href={`#${l.toLowerCase()}`}
              className="text-sm text-muted-foreground hover:text-foreground transition-colors tracking-wide"
            >
              {l}
            </Link>
          ))}
        </div>
        <div className="flex gap-3 items-center">
          <Link
            href={LOGIN_URL}
            className="text-sm text-muted-foreground hover:text-foreground transition-colors px-3 py-1.5"
          >
            Sign in
          </Link>
          <Link
            href={REGISTER_URL}
            className="font-mono font-bold text-xs bg-primary text-primary-foreground hover:opacity-90 px-5 py-2 transition-all hover:-translate-y-px tracking-wide"
          >
            Get started
          </Link>
        </div>
      </div>
    </nav>
  );
}

import { motion } from "framer-motion";
import React from "react";

interface AuthContainerProps {
  children: React.ReactNode;
  title?: string;
  subtitle?: string;
}

export function AuthContainer({
  children,
  title = "Access your workspace",
  subtitle = "Bangun, uji, dan luncurkan aplikasi Anda dalam hitungan detik.",
}: AuthContainerProps) {
  const words = title.split(" ");

  return (
    <div className="min-h-screen w-full flex items-center justify-center bg-[var(--color-app-bg)] bg-[radial-gradient(var(--color-app-border)_1px,transparent_1px)] [background-size:24px_24px] p-4 font-sans">
      <div className="w-full max-w-[860px] grid grid-cols-1 md:grid-cols-2 bg-[var(--color-app-surface)] border border-[var(--color-app-border)] shadow-2xl rounded-2xl overflow-hidden">

        {/* Kiri: Form */}
        <div className="p-10 sm:p-12 flex flex-col justify-center order-2 md:order-1 bg-[var(--color-sidebar-bg)]">
          <div className="w-full max-w-sm mx-auto">

            {/* Logo */}
            <div className="flex items-center gap-2 mb-9">
              <span className="text-[#326ce5] text-xl leading-none">☸</span>
              <span className="font-mono font-bold text-[15px] tracking-tighter text-[var(--color-primary)]">
                Kube<span className="text-[#326ce5]">Space</span>
              </span>
            </div>

            {children}
          </div>
        </div>

        {/* Kanan: Branding */}
        <div className="hidden md:flex flex-col items-center justify-center bg-[var(--color-sidebar-surface)] p-12 text-center border-l border-[var(--color-sidebar-border)] relative overflow-hidden order-1 md:order-2">

          {/* Subtle glow — pakai opacity rendah biar ga lebay */}
          <div className="absolute -top-16 -right-16 w-64 h-64 bg-[#326ce5]/[0.06] blur-[80px] rounded-full pointer-events-none" />
          <div className="absolute -bottom-12 -left-12 w-48 h-48 bg-[#326ce5]/[0.04] blur-[60px] rounded-full pointer-events-none" />

          <div className="relative z-10 text-center">

            {/* Badge kubernetes */}
            <div className="inline-flex items-center gap-1.5 bg-[var(--color-sidebar-bg)] border border-[var(--color-sidebar-border)] rounded-full px-3 py-1 mb-8">
              <span className="text-[#326ce5] text-sm leading-none">☸</span>
              <span className="text-[9px] font-medium text-[var(--color-sidebar-text)] uppercase tracking-[0.12em]">
                kubernetes-native
              </span>
            </div>

            {/* Title — animate per-word, cleaner dari per-char */}
            <h2 className="text-[28px] font-bold tracking-[-0.04em] text-[var(--color-primary)] leading-[1.2] mb-4">
              {words.map((word, i) => (
                <motion.span
                  key={i}
                  className="inline-block mr-[0.25em]"
                  initial={{ opacity: 0, y: 6 }}
                  animate={{ opacity: 1, y: 0 }}
                  transition={{
                    duration: 0.3,
                    delay: i * 0.08,
                    repeat: Infinity,
                    repeatDelay: 5,
                  }}
                >
                  {/* Highlight kata terakhir */}
                  {i === words.length - 1 ? (
                    <span className="text-[#326ce5]">{word}</span>
                  ) : (
                    word
                  )}
                </motion.span>
              ))}
            </h2>

            <p className="text-[12px] text-[var(--color-sidebar-text)] leading-[1.75] max-w-[200px] mx-auto">
              {subtitle}
            </p>

            
          </div>

          <div className="absolute bottom-5 text-[9px] text-[var(--color-sidebar-border)] uppercase tracking-[0.25em] font-medium">
            &copy; 2026 KubeSpace Engine
          </div>
        </div>
      </div>

      {/* Footer */}
      <div className="fixed bottom-5 text-[10px] text-[var(--color-sidebar-text-muted)] text-center w-full">
        By continuing, you agree to our{" "}
        <a href="#" className="text-[var(--color-sidebar-text)] underline hover:text-[#326ce5] transition-colors">
          Terms
        </a>{" "}
        &{" "}
        <a href="#" className="text-[var(--color-sidebar-text)] underline hover:text-[#326ce5] transition-colors">
          Privacy Policy
        </a>
      </div>
    </div>
  );
}
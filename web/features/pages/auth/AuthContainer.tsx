import React from "react";
import { motion } from "framer-motion";

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
  const characters = Array.from(title);

  return (
    <div className="min-h-screen w-full flex items-center justify-center bg-[#0a0a0a] bg-[radial-gradient(#ffffff1a_1px,transparent_1px)] [background-size:20px_20px] p-4">
      {/* Main Card Container - Responsive Max Width */}
      <div className="w-full max-w-[1000px] min-h-[min(600px,90vh)] grid grid-cols-1 md:grid-cols-2 bg-[#1c1c1c] border border-border shadow-2xl rounded-xl overflow-hidden transition-all duration-300">
        {/* Sisi Kiri: Form Area */}
        <div className="p-6 sm:p-8 md:p-12 flex flex-col justify-center order-2 md:order-1">
          <div className="w-full max-w-sm mx-auto">
            <div className="flex items-center gap-2 mb-8 justify-center md:justify-start">
              <span className="text-primary text-2xl">⬡</span>
              <span className="font-mono font-bold text-lg tracking-tighter text-foreground">
                KubeSpace
              </span>
            </div>
            {children}
          </div>
        </div>

        {/* Sisi Kanan: Branding Area - Responsive Text */}
        <div className="hidden md:flex flex-col items-center justify-center bg-[#252525] p-8 lg:p-12 text-center border-l border-border relative overflow-hidden order-1 md:order-2">
          <div className="absolute top-0 right-0 w-32 h-32 bg-primary/10 blur-[80px]" />

          <div className="relative z-10 w-full">
            <span className="text-4xl mb-6 block animate-bounce">👍</span>

            {/* Typing Animation Title */}
            <h2 className="text-[clamp(1.5rem,4vw,2.5rem)] font-bold tracking-tight mb-4 text-foreground leading-tight">
              {characters.map((char, i) => (
                <motion.span
                  key={i}
                  initial={{ opacity: 0 }}
                  animate={{ opacity: 1 }}
                  transition={{
                    duration: 0.1,
                    delay: i * 0.05,
                    repeat: Infinity,
                    repeatDelay: 3, // Jeda sebelum ngetik ulang
                  }}
                >
                  {char}
                </motion.span>
              ))}
              <motion.span
                animate={{ opacity: [0, 1, 0] }}
                transition={{ duration: 0.8, repeat: Infinity }}
                className="text-primary ml-1"
              >
                |
              </motion.span>
            </h2>

            <p className="text-sm lg:text-base text-muted-foreground font-light leading-relaxed max-w-[280px] mx-auto">
              {subtitle}
            </p>
          </div>

          <div className="absolute bottom-8 text-[10px] text-muted-foreground/50 uppercase tracking-[0.2em]">
            &copy; 2026 KubeSpace Platform
          </div>
        </div>
      </div>

      {/* Legal Links - Hidden on very small screens */}
      <div className="fixed bottom-6 text-[10px] sm:text-[11px] text-muted-foreground text-center px-4">
        By clicking continue, you agree to our{" "}
        <a href="#" className="underline hover:text-foreground">
          Terms of Service
        </a>{" "}
        dan{" "}
        <a href="#" className="underline hover:text-foreground">
          Privacy Policy
        </a>
        .
      </div>
    </div>
  );
}

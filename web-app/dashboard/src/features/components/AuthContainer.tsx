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
  const characters = Array.from(title);

  return (
    <div className="min-h-screen w-full flex items-center justify-center bg-[#09090b] bg-[radial-gradient(#326ce515_1px,transparent_1px)] [background-size:24px_24px] p-4 font-sans">
      {/* Main Card */}
      <div className="w-full max-w-[1000px] min-h-[600px] grid grid-cols-1 md:grid-cols-2 bg-[#121214] border border-zinc-800 shadow-2xl rounded-2xl overflow-hidden">
        {/* Sisi Kiri: Form Area */}
        <div className="p-8 sm:p-12 flex flex-col justify-center order-2 md:order-1 bg-[#121214]">
          <div className="w-full max-w-sm mx-auto">
            {/* Logo KubeSpace */}
            <div className="flex items-center gap-2 mb-10 justify-center md:justify-start">
              <span className="text-[#326ce5] text-3xl">☸</span>
              <span className="font-mono font-bold text-xl tracking-tighter text-white">
                Kube<span className="text-[#326ce5]">Space</span>
              </span>
            </div>
            {children}
          </div>
        </div>

        {/* Sisi Kanan: Branding Area */}
        <div className="hidden md:flex flex-col items-center justify-center bg-[#18181b] p-12 text-center border-l border-zinc-800 relative overflow-hidden order-1 md:order-2">
          {/* K8s Ambient Glow */}
          <div className="absolute -top-20 -right-20 w-80 h-80 bg-[#326ce5]/10 blur-[100px] rounded-full" />
          <div className="absolute -bottom-20 -left-20 w-60 h-60 bg-[#326ce5]/5 blur-[80px] rounded-full" />

          <div className="relative z-10 w-full">
            {/* <div className="mb-8 inline-flex p-4 rounded-3xl bg-[#326ce5]/10 border border-[#326ce5]/20 shadow-inner">
              <span className="text-4xl animate-pulse">🚀</span>
            </div> */}

            {/* Typing Animation Title */}
            <h2 className="text-4xl font-bold tracking-tight mb-6 text-white leading-tight min-h-[80px]">
              {characters.map((char, i) => (
                <motion.span
                  key={i}
                  initial={{ opacity: 0 }}
                  animate={{ opacity: 1 }}
                  transition={{
                    duration: 0.05,
                    delay: i * 0.05,
                    repeat: Infinity,
                    repeatDelay: 4,
                  }}
                >
                  {char}
                </motion.span>
              ))}
              <motion.span
                animate={{ opacity: [0, 1, 0] }}
                transition={{ duration: 0.8, repeat: Infinity }}
                className="text-[#326ce5] ml-1"
              >
                _
              </motion.span>
            </h2>

            <p className="text-zinc-400 font-light leading-relaxed max-w-[300px] mx-auto text-sm lg:text-base">
              {subtitle}
            </p>
          </div>

          <div className="absolute bottom-8 text-[10px] text-zinc-600 uppercase tracking-[0.3em] font-medium">
            &copy; 2026 KubeSpace Engine
          </div>
        </div>
      </div>

      {/* Footer Links */}
      <div className="fixed bottom-6 text-[11px] text-zinc-500 text-center w-full">
        By continuing, you agree to our{" "}
        <a href="#" className="text-zinc-300 underline hover:text-[#326ce5]">
          Terms
        </a>{" "}
        &{" "}
        <a href="#" className="text-zinc-300 underline hover:text-[#326ce5]">
          Privacy Policy
        </a>
      </div>
    </div>
  );
}

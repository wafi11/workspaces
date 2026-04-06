import { FEATURES } from "@/data/DataHome";

export function FeaturesSection() {
  return (
    <section
      className="relative z-10 max-w-6xl mx-auto px-8 py-20"
      id="features"
    >
      <div className="font-mono text-xs text-muted-foreground tracking-widest mb-3">
        // capabilities
      </div>
      <h2 className="text-[clamp(1.8rem,3vw,2.5rem)] font-semibold tracking-tight mb-12 text-foreground">
        Everything you need to ship
      </h2>
      <div
        className="grid gap-px bg-border border border-border"
        style={{
          gridTemplateColumns: "repeat(auto-fill, minmax(320px, 1fr))",
        }}
      >
        {FEATURES.map((f, i) => (
          <div
            key={i}
            className="bg-card hover:bg-secondary p-8 transition-colors"
          >
            <span className="text-2xl text-primary block mb-4">{f.icon}</span>
            <h3 className="text-base font-semibold mb-2 tracking-tight text-card-foreground">
              {f.title}
            </h3>
            <p className="text-sm text-muted-foreground leading-relaxed">
              {f.desc}
            </p>
          </div>
        ))}
      </div>
    </section>
  );
}

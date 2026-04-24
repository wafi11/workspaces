import { useProfileProviders } from "@/features/api/profile";
import type { ProviderUser } from "@/types";
import { useState } from "react";

const PROVIDER_ICONS: Record<string, React.ReactNode> = {
  github: (
    <svg width="16" height="16" viewBox="0 0 16 16" fill="none">
      <path d="M8 0C3.58 0 0 3.58 0 8c0 3.54 2.29 6.53 5.47 7.59.4.07.55-.17.55-.38 0-.19-.01-.82-.01-1.49-2.01.37-2.53-.49-2.69-.94-.09-.23-.48-.94-.82-1.13-.28-.15-.68-.52-.01-.53.63-.01 1.08.58 1.23.82.72 1.21 1.87.87 2.33.66.07-.52.28-.87.51-1.07-1.78-.2-3.64-.89-3.64-3.95 0-.87.31-1.59.82-2.15-.08-.2-.36-1.02.08-2.12 0 0 .67-.21 2.2.82.64-.18 1.32-.27 2-.27.68 0 1.36.09 2 .27 1.53-1.04 2.2-.82 2.2-.82.44 1.1.16 1.92.08 2.12.51.56.82 1.27.82 2.15 0 3.07-1.87 3.75-3.65 3.95.29.25.54.73.54 1.48 0 1.07-.01 1.93-.01 2.2 0 .21.15.46.55.38A8.013 8.013 0 0016 8c0-4.42-3.58-8-8-8z" fill="#52525b"/>
    </svg>
  ),
  google: (
    <svg width="16" height="16" viewBox="0 0 20 20" fill="none">
      <path d="M19.6 10.23c0-.68-.06-1.36-.18-2H10v3.77h5.4a4.6 4.6 0 01-2 3.02v2.5h3.24c1.9-1.75 3-4.33 3-7.29z" fill="#4285F4"/>
      <path d="M10 20c2.7 0 4.96-.9 6.62-2.43l-3.24-2.5c-.9.6-2.04.96-3.38.96-2.6 0-4.8-1.76-5.59-4.12H1.07v2.6A10 10 0 0010 20z" fill="#34A853"/>
      <path d="M4.41 11.91A6.02 6.02 0 014.1 10c0-.66.11-1.3.31-1.91V5.49H1.07A10 10 0 000 10c0 1.61.39 3.13 1.07 4.51l3.34-2.6z" fill="#FBBC05"/>
      <path d="M10 3.96c1.47 0 2.79.5 3.83 1.5l2.86-2.86C14.95.99 12.7 0 10 0A10 10 0 001.07 5.49l3.34 2.6C5.2 5.72 7.4 3.96 10 3.96z" fill="#EA4335"/>
    </svg>
  ),
};

function getProviderName(providerId: string): string {
  if (providerId.startsWith("github")) return "GitHub";
  if (providerId.startsWith("google")) return "Google";
  if (providerId.startsWith("gitlab")) return "GitLab";
  return providerId.split("|")[0] ?? providerId;
}

function getProviderKey(providerId: string): string {
  return providerId.split("|")[0]?.toLowerCase() ?? "";
}

function ProviderIcon({ providerId }: { providerId: string }) {
  const key = getProviderKey(providerId);
  const icon = PROVIDER_ICONS[key];
  return (
    <div className="w-8 h-8 rounded-lg border border-[#1c1c1c] bg-[#111111] flex items-center justify-center flex-shrink-0">
      {icon ?? (
        <span className="text-[11px] font-mono text-[#3f3f46] uppercase">
          {key.slice(0, 2)}
        </span>
      )}
    </div>
  );
}

export function ProfileProviders() {
  const { data: providers } = useProfileProviders();

  return (
    <div className="p-6 w-full">
      <p className="text-[10px] uppercase tracking-widest text-[#2e2e2e] font-semibold mb-2.5">
        Connected providers
      </p>

      <div className="bg-[#0a0a0a] border border-[#1c1c1c] rounded-xl overflow-hidden">
        {/* Header */}
        <div className="flex items-center justify-between px-5 py-3 border-b border-[#161616]">
          <span className="text-xs font-medium text-[#52525b] tracking-wide">
            OAuth providers
          </span>
          <span className="text-[10px] text-[#3f3f46] bg-[#111111] border border-[#1c1c1c] rounded px-1.5 py-0.5 font-mono">
            {providers?.length ?? 0} connected
          </span>
        </div>

        {/* List */}
        {!providers || providers.length === 0 ? (
          <div className="py-8 text-center">
            <p className="text-xs text-[#2e2e2e]">no providers connected</p>
          </div>
        ) : (
          <div className="divide-y divide-[#111111]">
            {providers.map((p) => (
              <ProviderRow key={p.provider_id} provider={p} />
            ))}
          </div>
        )}

        {/* Footer */}
        <div className="flex justify-end px-5 py-2.5 border-t border-[#111111]">
          <button className="flex items-center gap-1.5 px-2.5 py-1 text-[11px] text-[#52525b] border border-[#1c1c1c] rounded-md hover:text-[#a1a1aa] hover:border-[#3f3f46] hover:bg-[#111111] transition-all">
            <span className="text-base leading-none">+</span>
            Add provider
          </button>
        </div>
      </div>
    </div>
  );
}

function ProviderRow({ provider }: { provider: ProviderUser }) {
  const [hovered, setHovered] = useState(false);

  return (
    <div
      className="flex items-center gap-3 px-5 py-3 hover:bg-[#0d0d0d] transition-colors"
      onMouseEnter={() => setHovered(true)}
      onMouseLeave={() => setHovered(false)}
    >
      <ProviderIcon providerId={provider.provider_id} />

      <div className="flex-1 min-w-0">
        <p className="text-[13px] font-medium text-[#a1a1aa]">
          {provider.name || getProviderName(provider.provider_id)}
        </p>
        <p className="text-[11px] font-mono text-[#3f3f46] truncate mt-0.5">
          {provider.provider_id}
        </p>
      </div>

        <div className="flex items-center gap-1.5 text-[10px] text-[#3f3f46]">
          <span className="w-1.5 h-1.5 rounded-full bg-green-500" />
          connected
        </div>
      
    </div>
  );
}
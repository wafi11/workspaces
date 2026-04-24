import type { User } from "@/types";
import { GetInitials } from "../../GetInitials";
import { ButtonUploadPhoto } from "../../ButtonUploadPhoto";
import { formatDate, TimeAgo } from "@/utils/formatDate";
import { CopyableValue } from "../../CopyValue";
import { ExternalIcon } from "../../ExternalIcons";

interface ProfileSectionProps {
    user: User;
    isAdmin: boolean;
}
export function ProfileSection({ user, isAdmin }: ProfileSectionProps){
    return (
        <div className="p-6 w-full">
      <p className="text-[10px] uppercase tracking-widest text-[#2e2e2e] font-semibold mb-2.5 pl-0.5">
        Account
      </p>

      <div className="bg-[#0a0a0a] border border-[#1c1c1c] rounded-xl overflow-hidden">
        {/* Header */}
        <div className="relative flex items-center gap-4 px-6 py-5 border-b border-[#161616]">
          <div
            className="absolute inset-x-0 top-0 h-px"
            style={{ background: "linear-gradient(90deg, transparent, #2e2e2e 30%, #3f3f46 60%, transparent)" }}
          />

          {/* Avatar */}
          <div className="relative flex-shrink-0">
            <div className="size-20 rounded-[10px] relative border border-[#232323] bg-[#111111] flex items-center justify-center text-xl font-medium text-[#a1a1aa] overflow-hidden">
              {user.avatar_url ? (
                <img src={user.avatar_url} alt={user.username} className="w-full h-full object-cover" />
              ) : (
                <GetInitials username={user.username} />
              )}
            </div>
            <span className="absolute -bottom-0.5 -right-0.5 w-2.5 h-2.5 rounded-full bg-green-500 border-2 border-[#0a0a0a]" />
            <ButtonUploadPhoto onClick={() => {}} />
          </div>

          {/* Name / role */}
          <div className="flex-1 min-w-0">
            <p className="text-[15px] font-medium text-[#e4e4e7] truncate tracking-[-0.02em]">
              {user.username}
            </p>
            <p className="text-xs text-[#52525b] mt-0.5 font-mono">
              @{user.username}
            </p>
            <span
              className={`inline-flex items-center gap-1 mt-1.5 px-1.5 py-0.5 rounded text-[10px] font-medium uppercase tracking-wider border
                ${isAdmin
                  ? "bg-[#0f0a1f] text-[#a78bfa] border-[#2e1f54]"
                  : "bg-[#161616] text-[#71717a] border-[#232323]"
                }`}
            >
              <span className={`w-1 h-1 rounded-full ${isAdmin ? "bg-[#a78bfa]" : "bg-[#3f3f46]"}`} />
              {user.role.toLowerCase()}
            </span>
          </div>
        </div>

        {/* Details */}
        <div className="px-6 divide-y divide-[#111111]">
          {[
            { label: "Email", value: user.email, highlight: true },
            { label: "Member since", value: formatDate(user.created_at) },
            { label: "Last updated", value: TimeAgo(user.updated_at) },
          ].map(({ label, value, highlight }) => (
            <div key={label} className="flex items-center justify-between py-2.5">
              <span className="text-[11px] uppercase tracking-[0.08em] text-[#3f3f46] font-medium">
                {label}
              </span>
              <CopyableValue value={value} highlight={highlight} />
            </div>
          ))}
        </div>

      </div>
    </div>
    )
}
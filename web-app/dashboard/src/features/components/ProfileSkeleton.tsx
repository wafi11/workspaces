
export function ProfileSkeleton() {
  return (
    <div className="py-6 max-w-lg">
      <div className="bg-[#0a0a0a] border border-[#1c1c1c] rounded-xl p-6 space-y-4 animate-pulse">
        <div className="flex gap-4">
          <div className="w-14 h-14 rounded-[10px] bg-[#161616]" />
          <div className="flex-1 space-y-2">
            <div className="h-4 w-32 bg-[#161616] rounded" />
            <div className="h-3 w-24 bg-[#111111] rounded" />
          </div>
        </div>
        {[...Array(4)].map((_, i) => (
          <div key={i} className="flex justify-between py-2 border-t border-[#111111]">
            <div className="h-2.5 w-20 bg-[#111111] rounded" />
            <div className="h-2.5 w-32 bg-[#161616] rounded" />
          </div>
        ))}
      </div>
    </div>
  );
}
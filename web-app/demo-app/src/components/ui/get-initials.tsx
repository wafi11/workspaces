export function GetInitials({ username }: { username: string }) {
  return (
    <div className="w-10 h-10 rounded-full bg-[#1c1c1c] flex items-center justify-center text-sm font-medium text-white">
      {username.slice(0, 2).toUpperCase()}
    </div>
  )
}
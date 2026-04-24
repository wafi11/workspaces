import { Camera } from "lucide-react";

interface ButtonUploadPhotoProps {
    onClick: () => void;
}

export function ButtonUploadPhoto({ onClick }: ButtonUploadPhotoProps) {
  return (
    <button
      onClick={onClick}
      className="absolute inset-0 rounded-[10px] flex items-center justify-center bg-black/0 hover:bg-black/50 transition-all group"
    >
      <Camera className="size-4 text-white opacity-0 group-hover:opacity-100 transition-opacity" />
    </button>
  )
}
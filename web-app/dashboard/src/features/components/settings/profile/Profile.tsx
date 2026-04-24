import { useProfile } from "@/features/api";
import { ProfileSkeleton } from "../../ProfileSkeleton";
import { ProfileSection } from "./ProfileSection";
import { ProfileProviders } from "./ProfileProviders";

export function ProfilePage() {
  const { data: user } = useProfile();

  if (!user) return <ProfileSkeleton />;

  const isAdmin = user.role === "admin";

  return (
   <>
    <ProfileSection user={user} isAdmin={isAdmin} />
    <ProfileProviders />
   </>
  );
}


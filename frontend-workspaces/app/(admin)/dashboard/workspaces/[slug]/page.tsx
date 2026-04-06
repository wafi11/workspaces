import { WorkspaceDetails } from "@/features/pages/dashboard/workspaces/WorkspaceDetails";

interface PageProps {
  slug: string;
}
export default async function Page({ params }: { params: PageProps }) {
  const { slug } = await params;
  return <WorkspaceDetails slug={slug} />;
}

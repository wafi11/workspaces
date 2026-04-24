import { ADMIN_ROLE } from "@/constants";
import { useProfile } from "@/features/api";
import { ButtonCreate } from "@/features/components/ButtonCreate";
import { ListTemplates } from "@/features/components/templates/ListTemplates";
import { LoadingScreen } from "@/features/layout/loadingScreen";
import { MainContainer } from "@/features/layout/MainContainer";
import { TopbarAdmin } from "@/features/layout/TopbarDashboard";
import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/templates/")({
  component: RouteComponent,
});

function RouteComponent() {
  const {data} = useProfile()
  if (!data) {
    return <LoadingScreen />
  }
  return (
   <MainContainer >
         <TopbarAdmin title="List Templates">
          {
            data.role  === ADMIN_ROLE && (
              <ButtonCreate label="Add" to="templates" className={"justify-end items-end"}/>
            )
          }
         </TopbarAdmin>
         <ListTemplates  profile={data}/>
       </MainContainer>
  );
}

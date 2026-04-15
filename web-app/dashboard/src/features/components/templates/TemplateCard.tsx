import { ADMIN_ROLE } from "@/constants"
import type { Templates, User } from "@/types"
import { getCategoryStyle } from "@/utils/categoryTemplateColor"
import { Button } from "@radix-ui/themes"
import { useRouter } from "@tanstack/react-router"

interface TemplateCardProps {
    template : Templates
    profile : User
}
export function TemplateCard({template,profile} : TemplateCardProps){
    const {navigate}  = useRouter()
    return (
        <div
            key={template.id}
            className="group relative flex flex-col rounded-xl border border-sidebar-border bg-sidebar-surface p-5 transition-all duration-200 hover:border-sidebar-text/30 hover:bg-white/2 cursor-pointer"
          >
            {/* Icon & Category */}
            <div className="flex items-start justify-between mb-4">
              <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-black border border-sidebar-border group-hover:border-sidebar-text/50 transition-colors">
                {template.icon ? (
                  <img src={template.icon} alt="" className="h-6 w-6 object-contain" />
                ) : (
                  <span className="text-xs font-bold text-sidebar-text-active">
                    {template.name.charAt(0)}
                  </span>
                )}
              </div>
              <span className={`inline-flex items-center rounded-md px-2 py-0.5 text-[10px] font-medium ring-1 ring-inset transition-all ${getCategoryStyle(template.category)}`}>
                {template.category.toUpperCase()}
              </span>
            </div>

            {/* Content */}
            <div className="flex-1">
              <h3 className="text-sm font-medium text-sidebar-text-active mb-1 group-hover:text-white transition-colors">
                {template.name}
              </h3>
              <p className="text-xs leading-relaxed text-sidebar-text line-clamp-2">
                {template.description}
              </p>
            </div>

            {/* Footer / Meta */}
            <div className="mt-5 flex items-center justify-between">
              {
                profile.role === ADMIN_ROLE && (
                 <Button onClick={()  => navigate({to : `/templates/${template.id}`})} className="flex hover:bg-sidebar-text py-1 px-1.5 rounded-md items-center gap-1 text-[10px] font-medium text-sidebar-primary  transition-colors">
                  <span>Template Details</span>
                  <svg width="12" height="12" viewBox="0 0 15 15" fill="none" xmlns="http://www.w3.org/2000/svg">
                    <path d="M8.14645 3.14645C8.34171 2.95118 8.65829 2.95118 8.85355 3.14645L12.8536 7.14645C13.0488 7.34171 13.0488 7.65829 12.8536 7.85355L8.85355 11.8536C8.65829 12.0488 8.34171 12.0488 8.14645 11.8536C7.95118 11.6583 7.95118 11.3417 8.14645 11.1464L11.2929 8H2.5C2.22386 8 2 7.77614 2 7.5C2 7.22386 2.22386 7 2.5 7H11.2929L8.14645 3.85355C7.95118 3.65829 7.95118 3.34171 8.14645 3.14645Z" fill="currentColor" fillRule="evenodd" clipRule="evenodd"></path>
                  </svg>
                </Button>
                )
              }
              <Button onClick={()  => navigate({to : `/workspaces/create?templateId=${template.id}`})} className="flex hover:bg-sidebar-text py-1 px-1.5 rounded-md items-center gap-1 text-[10px] font-medium text-sidebar-primary  transition-colors">
                <span>Use Template</span>
                <svg width="12" height="12" viewBox="0 0 15 15" fill="none" xmlns="http://www.w3.org/2000/svg">
                  <path d="M8.14645 3.14645C8.34171 2.95118 8.65829 2.95118 8.85355 3.14645L12.8536 7.14645C13.0488 7.34171 13.0488 7.65829 12.8536 7.85355L8.85355 11.8536C8.65829 12.0488 8.34171 12.0488 8.14645 11.8536C7.95118 11.6583 7.95118 11.3417 8.14645 11.1464L11.2929 8H2.5C2.22386 8 2 7.77614 2 7.5C2 7.22386 2.22386 7 2.5 7H11.2929L8.14645 3.85355C7.95118 3.65829 7.95118 3.34171 8.14645 3.14645Z" fill="currentColor" fillRule="evenodd" clipRule="evenodd"></path>
                </svg>
              </Button>
            </div>
          </div>
    )
}
import { cn } from "@/lib/utils";
import React, { forwardRef } from "react";
import { SidebarProfile } from "./SidebarProfile";


type PageContainerProps = {
  withSidebar?: boolean;
  withFooter?: boolean;
};

export const PageContainer = forwardRef<
  HTMLElement,
  React.HTMLAttributes<HTMLElement> & PageContainerProps
>(
  (
    {
      className,
      children,
      withSidebar = true,
      withFooter = true,
      ...props
    },
    ref,
  ) => {

    return (
      <>
       <div className="flex flex-1 mt-0 overflow-hidden">
        {withSidebar && <SidebarProfile />}
        <main
          ref={ref}
          className={cn(
            "flex-1 overflow-y-auto min-w-0",
            withSidebar && "ml-[260px]",
            className
          )}
          {...props}
        >
          <div className="p-6">{children}</div>
          {withFooter && (
            <footer className="flex min-h-16 border-t-2 p-4">
              <p className="w-full text-center text-muted-foreground">
                © {new Date().getFullYear()} Wafiuddin. All rights reserved
              </p>
            </footer>
          )}
        </main>
      </div>
        </>
    );
  },
);

PageContainer.displayName = "PageContainer";
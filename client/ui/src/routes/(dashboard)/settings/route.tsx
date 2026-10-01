import { ChatNavActions } from "#/components/chat-nav-actions";
import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator,
} from "#/components/ui/breadcrumb";
import { Separator } from "#/components/ui/separator";
import { SidebarTrigger } from "#/components/ui/sidebar";
import {
  createFileRoute,
  Outlet,
  useRouterState,
} from "@tanstack/react-router";

export const Route = createFileRoute("/(dashboard)/settings")({
  component: RouteComponent,
});

function RouteComponent() {
  const matches = useRouterState({ select: (s) => s.matches });

  return (
    <div className="flex flex-col h-screen overflow-hidden">
      <header className="flex h-14 shrink-0 items-center gap-2 border-b bg-background sticky top-0 z-10">
        <div className="flex flex-1 items-center gap-2 px-3">
          <SidebarTrigger />
          <Separator orientation="vertical" />
          <Breadcrumb>
            <BreadcrumbList>
              {matches.map((match, index) => {
                const staticData = match.staticData as
                  { breadcrumb?: string } | undefined;
                const breadcrumbName = staticData?.breadcrumb;

                if (!breadcrumbName) return null;

                const isLast = index === matches.length - 1;

                return (
                  <BreadcrumbItem key={match.id}>
                    {isLast ? (
                      <BreadcrumbPage className="line-clamp-1">
                        {breadcrumbName}
                      </BreadcrumbPage>
                    ) : (
                      <>
                        <BreadcrumbLink href={match.pathname}>
                          {breadcrumbName}
                        </BreadcrumbLink>
                        <BreadcrumbSeparator />
                      </>
                    )}
                  </BreadcrumbItem>
                );
              })}
            </BreadcrumbList>
          </Breadcrumb>
        </div>
      </header>

      <div className="flex-1  flex flex-col overflow-y-auto">
        <Outlet />
      </div>
    </div>
  );
}

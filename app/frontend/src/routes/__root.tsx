import { Outlet, createRootRoute } from '@tanstack/react-router'
import '../styles.css'
import { AppSidebar } from "@/components/app-sidebar"
import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
} from "@/components/ui/breadcrumb"
import {
  SidebarInset,
  SidebarProvider,
} from "@/components/ui/sidebar"

export const Route = createRootRoute({
  component: RootComponent,
})

function RootComponent() {
  return (
    <>
      <SidebarProvider className="h-screen w-screen overflow-hidden">
        <AppSidebar />
        <SidebarInset className="relative h-screen w-full overflow-hidden flex flex-col">
          <header className="flex h-16 shrink-0 items-center gap-2 border-b">
            <div className="flex items-center gap-2 px-4">
              <Breadcrumb>
                <BreadcrumbList>
                  <BreadcrumbItem className="hidden md:block">
                    <BreadcrumbLink href="#">
                      Chat
                    </BreadcrumbLink>
                  </BreadcrumbItem>
                </BreadcrumbList>
              </Breadcrumb>
            </div>
          </header>

          <div className="flex-1 overflow-hidden relative flex flex-col">
            <Outlet />
          </div>
        </SidebarInset>
      </SidebarProvider>
    </>
  )
}
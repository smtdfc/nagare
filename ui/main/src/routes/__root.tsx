import { Outlet, createRootRoute } from '@tanstack/react-router'
import { AppSidebar } from '@/components/app-sidebar'
import { SidebarInset, SidebarProvider } from '@/components/ui/sidebar'
import { DynamicBreadcrumbs } from '@/components/dynamic-breadcrumbs.tsx'

export const Route = createRootRoute({
  component: RootComponent,
})

function RootComponent() {
  return (
    <>
      <SidebarProvider className="h-screen w-screen overflow-hidden">
        <AppSidebar />
        <SidebarInset className="relative h-screen w-full overflow-hidden flex flex-col">
          <header className="flex h-16 shrink-0 items-center gap-2">
            <div className="flex items-center gap-2 px-4">
              <DynamicBreadcrumbs />
            </div>
          </header>
          <div className="flex-1 relative flex flex-col overflow-y-auto">
            <Outlet />
          </div>
        </SidebarInset>
      </SidebarProvider>
    </>
  )
}

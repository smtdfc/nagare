import { Outlet, createRootRoute } from '@tanstack/react-router'
import '../styles.css'
import { AppSidebar } from "@/components/app-sidebar"
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
          </header>
          <div className="flex-1 overflow-hidden relative flex flex-col">
            <Outlet />
          </div>
        </SidebarInset>
      </SidebarProvider>
    </>
  )
}
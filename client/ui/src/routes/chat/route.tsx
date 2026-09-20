import { createFileRoute, Outlet} from '@tanstack/react-router'
import { NavActions } from "#/components/nav-actions.tsx"
import {
    Breadcrumb,
    BreadcrumbItem,
    BreadcrumbList,
    BreadcrumbPage,
} from "#/components/ui/breadcrumb.tsx"
import { Separator } from "#/components/ui/separator.tsx"
import {
    SidebarTrigger,
} from "#/components/ui/sidebar.tsx"
import {Skeleton} from "#/components/ui/skeleton.tsx";

export const Route = createFileRoute('/chat')({
    component: RouteComponent,
})

function RouteComponent() {
    return (
        <div className="flex flex-col h-screen overflow-hidden">
            <header className="flex h-14 shrink-0 items-center gap-2 border-b bg-background sticky top-0 z-10">
                <div className="flex flex-1 items-center gap-2 px-3">
                    <SidebarTrigger />
                    <Separator
                        orientation="vertical"
                    />
                    <Breadcrumb>
                        <BreadcrumbList>
                            <BreadcrumbItem>
                                <BreadcrumbPage className="line-clamp-1">
                                    <Skeleton className="h-5 w-75 rounded-full" />
                                </BreadcrumbPage>
                            </BreadcrumbItem>
                        </BreadcrumbList>
                    </Breadcrumb>
                </div>
                <div className="ml-auto px-3">
                    <NavActions />
                </div>
            </header>

            <div className="flex-1 overflow-y-auto flex flex-col gap-4 px-4 py-10">
                <Outlet/>
            </div>
        </div>
    )
}
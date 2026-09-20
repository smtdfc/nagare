import { createFileRoute } from '@tanstack/react-router'
import LoadingLayout from "#/components/loading-layout.tsx";


export const Route = createFileRoute('/chat/$id')({
    component: RouteComponent,
})

function RouteComponent() {
    return(
        <>
         <LoadingLayout/>
        </>
    )
}

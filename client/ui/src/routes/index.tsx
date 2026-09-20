import { createFileRoute } from '@tanstack/react-router'


export const Route = createFileRoute('/')({
  component: RouteComponent,
})

function RouteComponent() {
  return <h1>ddd</h1>
}

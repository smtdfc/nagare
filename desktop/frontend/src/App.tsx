import { RouterProvider } from '@tanstack/react-router'
import { TooltipProvider } from './components/ui/tooltip'
import { Toaster } from './components/ui/toast'

export default function App({ router }: any) {
  return (
    <TooltipProvider>
      <Toaster />
      <RouterProvider router={router} />
    </TooltipProvider>
  )
}

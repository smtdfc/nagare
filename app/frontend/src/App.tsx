import { RouterProvider } from '@tanstack/react-router';
import { TooltipProvider } from './components/ui/tooltip';

export default function App({ router }: any) {
    return (
        <TooltipProvider>
            <RouterProvider router={router} />
        </TooltipProvider>
    );
}
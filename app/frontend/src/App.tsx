import { RouterProvider } from '@tanstack/react-router';
import { TooltipProvider } from './components/ui/tooltip';
import { initWebsocketConnection } from './services/chat';
import { useEffect } from 'react';



export default function App({ router }: any) {
    useEffect(() => {
        setTimeout(() => { initWebsocketConnection() }, 2000);
    }, []);

    return (
        <TooltipProvider>
            <RouterProvider router={router} />
        </TooltipProvider>
    );
}
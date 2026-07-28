import { RouterProvider } from '@tanstack/react-router';
import { TooltipProvider } from './components/ui/tooltip';
import { initWebsocketConnection } from './services/chat';
import { useEffect, useState } from 'react';



export default function App({ router }: any) {
    const [isReady, setIsReady] = useState(false);
    useEffect(() => {
        setTimeout(() => { initWebsocketConnection().then(() => { setIsReady(true) }) }, 2000);
    }, []);

    return (
        <TooltipProvider>
            {isReady ? <RouterProvider router={router} /> : "Loading ..."}
        </TooltipProvider>
    );
}
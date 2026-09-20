import { Spinner } from "#/components/ui/spinner.tsx";

interface LoadingLayoutProps {
    text?: string;
}

export default function LoadingLayout({ text = "Loading ..." }: LoadingLayoutProps) {
    return (
        <div className="h-screen w-full flex flex-col items-center justify-center gap-4">
            <Spinner className="size-8" />
            {text && (
                <p className="text-sm text-muted-foreground animate-pulse">
                    {text}
                </p>
            )}
        </div>
    );
}
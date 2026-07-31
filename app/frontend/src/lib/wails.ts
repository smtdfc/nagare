export function waitForWails(maxRetries = 20, intervalMs = 50): Promise<void> {
    return new Promise((resolve, reject) => {
        let retries = 0;
        const interval = setInterval(() => {
            if (typeof window !== 'undefined' && (window as any).go) {
                clearInterval(interval);
                resolve();
            } else {
                retries++;
                if (retries >= maxRetries) {
                    clearInterval(interval);
                    reject(new Error("Wails bindings failed to load!"));
                }
            }
        }, intervalMs);
    });
}
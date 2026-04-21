export async function api<T>(
    url: string,
    options: RequestInit = {}
): Promise<T> {
    const token = localStorage.getItem("token");
    const headers = new Headers(options.headers || {});
    headers.set("Content-Type", "application/json");
    if (token) {
        headers.set("Authorization", `Bearer ${token}`);
    }

    const res = await fetch(url, {...options, headers});
    if (!res.ok) {
        const err = new Error(await res.text());
        (err as any).status = res.status;
        throw err;
    }

    return res.json().catch(() => ({} as T));
}

export async function apiBlob(url: string, options: RequestInit = {}): Promise<Blob> {
    const token = localStorage.getItem("token");
    const headers = new Headers(options.headers || {});
    if (token) {
        headers.set("Authorization", `Bearer ${token}`);
    }

    const res = await fetch(url, {...options, headers});
    if (!res.ok) {
        throw new Error(await res.text());
    }

    return res.blob();
}

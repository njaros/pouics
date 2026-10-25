import.meta.env.VITE_API_URL = "http://localhost:3333"
const TOKEN_KEY = "pouics.token"

export function getToken(): string | null {
  return localStorage.getItem(TOKEN_KEY);
}

export function setToken(token: string): void {
  localStorage.setItem(TOKEN_KEY, token);
}

export function clearToken(): void {
  localStorage.removeItem(TOKEN_KEY);
}

export class ApiError extends Error {

    public readonly status: number;
    public readonly name: string;

    constructor(
    message: string,
    status: number,
    ) {
        super(message);
        this.name = "ApiError";
        this.status = status;
    }
}

/**
 * 
 * @param path Url path. It does not include the server url, neither the firts '/'.
 * @param init Headers. "Content-Type", "application/json" and the bearer are automatically set.
 * @returns response or error as a Promise.
 */
export async function request<T>(path: string, init?: RequestInit): Promise<T> {
    const headers = new Headers(init?.headers);
    headers.set("Content-Type", "application/json");
    const token = getToken();
    if (token != null) {
        headers.set("Authorization", `Bearer ${token}`);
    }

    const response = await fetch(`${import.meta.env.VITE_API_URL}/${path}`, { ...init, headers });
    if (!response.ok) {
        const message = await response
            .json()
            .then((body: { error?: string }) => body.error ?? response.statusText)
            .catch(() => response.statusText);
        throw new ApiError(message, response.status);
    }
    return response.json() as Promise<T>;
}
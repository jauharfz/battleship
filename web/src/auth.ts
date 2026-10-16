import { GUEST_TOKEN_KEY } from "./constants";

export async function getGuestJWT(): Promise<string> {
  const cachedToken = sessionStorage.getItem(GUEST_TOKEN_KEY);
  if (cachedToken && !isTokenExpired(cachedToken)) {
    return cachedToken;
  }
  try {
    const res = await fetch("/api/v1/auth/guest", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
    });
    if (!res.ok) {
      throw new Error(
        `failed to fetch response from register guest, status: ${res.status}`,
      );
    }
    const payload = await res.json();
    const token = payload.data.token;
    sessionStorage.setItem(GUEST_TOKEN_KEY, token);
    return token;
  } catch (error) {
    console.error("failed to get guest JWT ", error);
    throw error;
  }
}

export function isTokenExpired(token: string): boolean {
  try {
    const payloadB64 = token.split(".")[1];
    const payload = JSON.parse(atob(payloadB64));

    if (!payload.exp) return true;

    return Date.now() >= payload.exp * 1000 - 8000;
  } catch {
    return true;
  }
}

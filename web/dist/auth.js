const GUEST_TOKEN_KEY = "guest_jwt";
export async function getGuestJWT() {
    const cachedToken = sessionStorage.getItem(GUEST_TOKEN_KEY);
    if (cachedToken) {
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
            throw new Error(`failed to fetch response from register guest, status: ${res.status}`);
        }
        const payload = await res.json();
        const token = payload.data.token;
        sessionStorage.setItem(GUEST_TOKEN_KEY, token);
        return token;
    }
    catch (error) {
        console.error("failed to get guest JWT ", error);
        throw error;
    }
}

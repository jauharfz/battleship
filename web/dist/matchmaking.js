export async function enterMatchmaking(token, roomID) {
    try {
        console.log(roomID);
        const res = await fetch("/api/v1/matchmaking", {
            method: "POST",
            headers: {
                Authorization: `Bearer ${token}`,
                "Content-Type": "application/json",
            },
            body: JSON.stringify(roomID ?? ""),
        });
        if (!res.ok) {
            console.error();
            throw new Error(`failed to fetch response from enter matchmaking, status: ${res.status}`);
        }
        const payload = await res.json();
        if (payload.status == "error") {
            throw new Error(payload.error);
        }
        return;
    }
    catch (error) {
        console.error("failed to enter matchmaking ", error);
        throw error;
    }
}
export async function poolMatchmaking(token) {
    try {
        while (true) {
            const res = await fetch("/api/v1/matchmaking/status", {
                method: "GET",
                headers: {
                    "Content-Type": "application/json",
                    Authorization: `Bearer ${token}`,
                },
            });
            if (!res.ok) {
                throw new Error(`failed to fetch response from check status: ${res.status}`);
            }
            const payload = await res.json();
            if (payload.status == "error") {
                throw new Error(payload.error);
            }
            if (payload.status == "timedout") {
                throw new Error(payload.status);
            }
            if (payload.status == "matched") {
                return payload.data.ws_token;
            }
            await new Promise((r) => setTimeout(r, 1000));
        }
    }
    catch (error) {
        console.error("failed to pool matchmaking ", error);
        throw error;
    }
}

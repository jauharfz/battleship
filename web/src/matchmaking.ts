export async function enterMatchmaking(token: string): Promise<void> {
  try {
    const res = await fetch("/api/v1/matchmaking", {
      method: "POST",
      headers: {
        Authorization: `Bearer ${token}`,
        "Content-Type": "application/json",
      },
    });
    if (!res.ok) {
      throw new Error(
        `failed to fetch response from enter matchmaking, status: ${res.status}`,
      );
    }
    return;
  } catch (error) {
    console.error("failed to enter matchmaking ", error);
    throw error;
  }
}

export async function poolMatchmaking(token: string): Promise<string> {
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
        throw new Error(
          `failed to fetch response from check status: ${res.status}`,
        );
      }
      const payload = await res.json();
      if (payload.status == "matched") {
        return payload.data.ws_token;
      }
      await new Promise((r) => setTimeout(r, 1000));
    }
  } catch (error) {
    console.error("failed to pool matchmaking ", error);
    throw error;
  }
}

var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    function adopt(value) { return value instanceof P ? value : new P(function (resolve) { resolve(value); }); }
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : adopt(result.value).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
export function enterMatchmaking(token) {
    return __awaiter(this, void 0, void 0, function* () {
        try {
            const res = yield fetch("/api/v1/matchmaking", {
                method: "POST",
                headers: {
                    Authorization: `Bearer ${token}`,
                    "Content-Type": "application/json",
                },
            });
            if (!res.ok) {
                throw new Error(`failed to fetch response from enter matchmaking, status: ${res.status}`);
            }
            return;
        }
        catch (error) {
            console.error("failed to enter matchmaking ", error);
            throw error;
        }
    });
}
export function poolMatchmaking(token) {
    return __awaiter(this, void 0, void 0, function* () {
        try {
            while (true) {
                const res = yield fetch("/api/v1/matchmaking/status", {
                    method: "GET",
                    headers: {
                        "Content-Type": "application/json",
                        Authorization: `Bearer ${token}`,
                    },
                });
                if (!res.ok) {
                    throw new Error(`failed to fetch response from check status: ${res.status}`);
                }
                const payload = yield res.json();
                if (payload.status == "matched") {
                    return payload.data.ws_token;
                }
                yield new Promise((r) => setTimeout(r, 1000));
            }
        }
        catch (error) {
            console.error("failed to pool matchmaking ", error);
            throw error;
        }
    });
}

var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    function adopt(value) { return value instanceof P ? value : new P(function (resolve) { resolve(value); }); }
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : adopt(result.value).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
const GUEST_TOKEN_KEY = "guest_jwt";
export function getGuestJWT() {
    return __awaiter(this, void 0, void 0, function* () {
        const cachedToken = sessionStorage.getItem(GUEST_TOKEN_KEY);
        if (cachedToken) {
            return cachedToken;
        }
        try {
            const res = yield fetch("/api/v1/auth/guest", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
            });
            if (!res.ok) {
                throw new Error(`failed to fetch response from register guest, status: ${res.status}`);
            }
            const payload = yield res.json();
            const token = payload.data.token;
            sessionStorage.setItem(GUEST_TOKEN_KEY, token);
            return token;
        }
        catch (error) {
            console.error("failed to get guest JWT ", error);
            throw error;
        }
    });
}

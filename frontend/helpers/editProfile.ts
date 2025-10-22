// helpers/patchProfile.ts
export async function patchProfile(token: string, form: FormData) {

    const res = await fetch(process.env.NEXT_PUBLIC_API_URL + `/api/user/profile`, {
        method: "PATCH",
        headers: { Authorization: `Bearer ${token}` },
        body: form,
    });

    const json = await res.json().catch(() => ({}));
    if (!res.ok || json?.status === false) {
        throw new Error(json?.message || `Gagal memperbarui profil (HTTP ${res.status})`);
    }
    return json;
}

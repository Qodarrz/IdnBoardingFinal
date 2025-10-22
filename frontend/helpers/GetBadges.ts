// utils/api.ts

export const GetBadges = async (token: string) => {
    const response = await fetch(process.env.NEXT_PUBLIC_API_URL + "/api/badges", {
        method: "GET", 
        headers: {
            "Authorization": `Bearer ${token}`, 
            "Content-Type": "application/json",
        },
    });

    if (!response.ok) {
        throw new Error("Failed to fetch badges");
    }

    const data = await response.json();
    return data.data.data; 
};

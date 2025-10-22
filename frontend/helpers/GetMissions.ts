// helpers/api.ts

export const GetMissions = async (token: string) => {
    const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/custom/mission-progress`, {
        headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`,
        },
    });

    if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
    }

    const data = await response.json();
    return data;
};

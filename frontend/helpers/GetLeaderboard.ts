export const GetLeaderboard = async (token: string) => {
    const res = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/custom/leaderboard`, {
        headers: {
            Authorization: `Bearer ${token}`,
        },
    });

    const data = await res.json();
    if (!data.status) {
        throw new Error(data.message || "Failed to fetch leaderboard");
    }

    return data.data.leaderboard;
};

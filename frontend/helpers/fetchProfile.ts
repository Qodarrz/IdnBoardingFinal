// helpers/fetchProfile.ts

interface Vehicle {
    id: number;
    vehicle_type: string;
    fuel_type: string;
    name: string;
    total_carbon_emission_g: number;
    created_at: string;
    total_logs: number;
}

interface Electronic {
    id: number;
    device_name: string;
    device_type: string;
    power_watts: number;
    total_carbon_emission_g: number;
    created_at: string;
    total_logs: number;
}

interface BadgeType {
    id: number;
    name: string;
    image_url: string;
    description: string;
    redeemed_at: string;
}

interface Mission {
    id: number;
    title: string;
    description: string;
    mission_type: string;
    points_reward: number;
    target_value: number;
    progress: number;
    completed_at: string;
    created_at: string;
}

interface PointHistory {
    id: number;
    amount: number;
    direction: string;
    source: string;
    reference_type: string;
    reference_id: number;
    created_at: string;
}

interface ActivityLog {
    id: number;
    activity: string;
    created_at: string;
}

interface OrderItem {
    id: number;
    item_name: string;
    qty: number;
    price_each_points: number;
}

interface Order {
    id: number;
    total_points: number;
    status: string;
    created_at: string;
    items: OrderItem[];
}

interface MonthlyCarbonEmission {
    month: string;
    total_carbon_emission_g: number;
}

interface ProfileData {
    user: {
        id: number;
        username: string;
        email: string;
        role: string;
        full_name: string;
        avatar_url: string;
        birthdate: string | null;
        gender: string;
        total_points: number;
        created_at: string;
    };
    vehicles: Vehicle[];
    electronics: Electronic[];
    missions: Mission[];
    badges: BadgeType[];
    point_history: PointHistory[];
    activity_logs: ActivityLog[];
    orders: Order[];
    monthly_vehicle_carbon: MonthlyCarbonEmission[];
    monthly_electronic_carbon: MonthlyCarbonEmission[];
}

export const fetchProfile = async (token: string): Promise<ProfileData | null> => {
    try {
        const res = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/custom/my-data`, {
            headers: {
                Authorization: `Bearer ${token}`,
            },
        })

        const data = await res.json()
        if (data.status) {
            return data.data
        } else {
            console.error("Gagal ambil profil:", data.message)
            return null
        }
    } catch (err) {
        console.error("Error fetch profil:", err)
        return null
    }
}

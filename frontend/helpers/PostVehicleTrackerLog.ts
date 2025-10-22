/* eslint-disable */
interface ElectronicPayload {
  start_lat: string | number;
  start_lon: string | number;
  end_lat: string | number;
  end_lon: string | number;
  distance_km: number;
  // duration_minutes: number;
  vehicle_id?: number | null;
}

interface ResponseData {
  status: string;
  data?: any;
  message: string;
}

export async function PostVehicleTrackerLog(payload: ElectronicPayload): Promise<ResponseData> {
  try {
    const res = await fetch(
      process.env.NEXT_PUBLIC_API_URL + "/api/carbon/vehicle-log",
      {
        method: "POST",
        headers: {
          "Authorization": "Bearer " + localStorage.getItem("authtoken"),
          "Content-Type": "application/json",
        },
        body: JSON.stringify(payload),
      }
    );

    // if (!res.ok) {
    //   throw new Error("Gagal menambahkan data elektronik");
    // }

    return await res.json();
  } catch (err: unknown) {
    let message = "Terjadi kesalahan";
    if (err instanceof Error) {
      message = err.message;
    }
    return { status: "error", message };
  }
}

/* eslint-disable */
interface ElectronicPayload {
  name: string;
  vehicle_type: string;
  fuel_type: string;
  user_id?: number | null;
}

interface ResponseData {
  status: string;
  data?: any;
  message: string;
}

export async function PostVehicleTracker(payload: ElectronicPayload): Promise<ResponseData> {
  try {
    const res = await fetch(
      process.env.NEXT_PUBLIC_API_URL + "/api/carbon/vehicle",
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

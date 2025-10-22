/* eslint-disable */
interface CarbonLogPayload {
  device_id: number;
  duration_hours: number;
}

interface ResponseData {
  status: string;
  data?: any;
  message: string;
}

export async function PostCarbonElectronicLog(payload: CarbonLogPayload): Promise<ResponseData> {
  try {
    const res = await fetch(
      process.env.NEXT_PUBLIC_API_URL + "/api/carbon/electronics-log", // Ganti endpoint sesuai dengan yang dibutuhkan
      {
        method: "POST",
        headers: {
          "Authorization": "Bearer " + localStorage.getItem("authtoken"),
          "Content-Type": "application/json",
        },
        body: JSON.stringify(payload),
      }
    );

    // Menangani status respon jika perlu
    if (!res.ok) {
      throw new Error("Gagal menambahkan log karbon");
    }

    return await res.json();
  } catch (err: unknown) {
    let message = "Terjadi kesalahan";
    if (err instanceof Error) {
      message = err.message;
    }
    return { status: "error", message };
  }
}

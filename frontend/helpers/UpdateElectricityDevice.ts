/* eslint-disable */

export async function UpdateElectricityDevice(id: number, payload: any) {
  try {
    const res = await fetch(
      process.env.NEXT_PUBLIC_API_URL + `/api/carbon/electronics/${id}`,
      {
        method: "PATCH",
        headers: {
          "Authorization": "Bearer " + localStorage.getItem("authtoken"),
          "Content-Type": "application/json",
        },
        body: JSON.stringify(payload), // Mengirimkan payload yang sudah diubah
      }
    );

    // Memeriksa jika response OK
    if (!res.ok) {
      const errorMessage = await res.text();
      throw new Error(`Error: ${errorMessage}`);
    }

    const data = await res.json();
    if (data.status) {
      return data;
    } else {
      throw new Error(data.message || "Gagal memperbarui perangkat");
    }
  } catch (err) {
    console.error("Update error:", err);
    throw err;
  }
}

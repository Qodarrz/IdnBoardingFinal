export async function DeleteElectricityDevice(id: number) {
  try {
    const res = await fetch(
      process.env.NEXT_PUBLIC_API_URL + `/api/carbon/electronics/${id}`,
      {
        method: "DELETE",
        headers: {
          "Authorization": "Bearer " + localStorage.getItem("authtoken"),
          "Content-Type": "application/json",
        },
      }
    );

    if (!res.ok) throw new Error("Gagal hapus device");
    return await res.json();
  } catch (err) {
    console.error("Delete error:", err);
    throw err;
  }
}

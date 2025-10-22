import { useState, useEffect } from "react";

interface DeviceTypes {
  id: number;
  user_id: number;
  device_name: string;
  device_type: string;
  power_watts: number;
  created_at: string;
}

interface ResponseData {
  status: boolean;
  message: string;
  data: DeviceTypes[];
}

export function GetElectricityDevice() {
  const [data, setData] = useState<ResponseData | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    const fetchElectricityDevices = async () => {
      try {
        const res = await fetch(
          process.env.NEXT_PUBLIC_API_URL + "/api/carbon/electronics",
          {
            headers: {
              "Authorization": "Bearer " + localStorage.getItem("authtoken"),
              "Content-Type": "application/json"
            },
          }
        )

        const result = await res.json()
        if (result.status) {
          setData(result)
        } else {
          setError(result.message)
        }
      } catch (err: unknown) {
        let message = "Terjadi kesalahan"
        if (err instanceof Error) {
          message = err.message
        }
        setError(message)
      } finally {
        setLoading(false)
      }
    }

    fetchElectricityDevices()
  }, [])

  return { data, loading, error }
}
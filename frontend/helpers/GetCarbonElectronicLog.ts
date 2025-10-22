import * as React from "react"
/* eslint-disable */

interface ResponseData {
  status: string;
  data?: any[]; // Bisa disesuaikan dengan data yang diperlukan untuk carbon log
  message: string;
}

export function GetCarbonElectronicLogs() {
  const [data, setData] = React.useState<ResponseData | null>(null)
  const [loading, setLoading] = React.useState(true)
  const [error, setError] = React.useState<string | null>(null)

  React.useEffect(() => {
    const fetchCarbonLogs = async () => {
      try {
        const res = await fetch(
          process.env.NEXT_PUBLIC_API_URL + "/api/carbon/electronics/logs",
          {
            headers: {
              "Authorization": "Bearer " + localStorage.getItem("authtoken"),
              "Content-Type": "application/json"
            },
          }
        )

        const result = await res.json()
        setData(result)
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

    fetchCarbonLogs()
  }, [])

  return { data, loading, error }
}

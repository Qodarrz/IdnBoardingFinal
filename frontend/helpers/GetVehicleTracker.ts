import * as React from "react"
/* eslint-disable */

interface ResponseData {
  status: string;
  data?: any[];
  message: string;
}

export function GetVehicleTracker() {
  const [data, setData] = React.useState<ResponseData | null>(null)
  const [loading, setLoading] = React.useState(true)
  const [error, setError] = React.useState<string | null>(null)

  React.useEffect(() => {
    const fetchAuth = async () => {
      try {
        const res = await fetch(
          process.env.NEXT_PUBLIC_API_URL + "/api/carbon/vehicles",
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

    fetchAuth()
  }, [])

  return { data, loading, error }
}

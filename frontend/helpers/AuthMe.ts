import * as React from "react"
interface Response1 {
    id?: string;
    username?: string;
    role?: string;
}
interface Response {
    status: string;
    message: string;
    data: Response1; 
}

export function useAuthMe() {
  const [data, setData] = React.useState<Response | null>(null)
  const [loading, setLoading] = React.useState(true)
  const [error, setError] = React.useState<string | null>(null)

  React.useEffect(() => {
    const fetchAuth = async () => {
      try {
        const res = await fetch(
          process.env.NEXT_PUBLIC_API_URL + "/api/auth/me",
          {
            headers: {
              "Authorization": "Bearer " + localStorage.getItem("authtoken"),
              "Content-Type": "application/json"
            },
          }
        )

        // if (!res.ok) {
        //   throw new Error("Failed to fetch api/auth/me")
        // }

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

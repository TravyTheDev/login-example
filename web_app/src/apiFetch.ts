const API_BASE = 'http://localhost:8000'

export async function apiFetch(endpoint: string, options: RequestInit = {}): Promise<Response> {
  const defaultOptions: RequestInit = {
    ...options,
    headers: {
      'Content-Type': 'application/json',
      ...options.headers,
    },
    credentials: 'include',
  }

  let response = await fetch(`${API_BASE}${endpoint}`, defaultOptions)

  if (response.status === 401 && !endpoint.includes('/renew_token')) {
    try {
      const refreshResponse = await fetch(`${API_BASE}/users/renew_token`, {
        method: 'POST',
        credentials: 'include',
      })

      if (refreshResponse.ok) {
        response = await fetch(`${API_BASE}${endpoint}`, defaultOptions)
      } else {
        throw new Error('Session expired')
      }
    } catch (refreshError) {
      throw new Error('Authentication failed')
    }
  }

  return response
}
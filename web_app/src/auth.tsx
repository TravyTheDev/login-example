import React, { createContext, useContext, useState, useEffect } from 'react'
import { apiFetch } from './apiFetch'

interface User {
  id: string
  username: string
  email: string
}

interface AuthState {
  isAuthenticated: boolean
  user: User | null
  login: (email: string, password: string) => Promise<void>
  logout: () => void
}

const AuthContext = createContext<AuthState | undefined>(undefined)

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<User | null>(null)
  const [isAuthenticated, setIsAuthenticated] = useState(false)
  const [isLoading, setIsLoading] = useState(true)

  const getMe = async () => {
    try {
        const res = await apiFetch('/users/me')
        
        if (!res.ok) {
            throw new Error('Failed to fetch user profile')
        }
        const user = await res.json()
        setUser(user)
        setIsAuthenticated(true)
    } catch (error) {
        setIsAuthenticated(false)
        throw new Error('Authentication failed')
    } finally {
        setIsLoading(false)
    }
  }

  useEffect(() => {
    getMe()
  }, [])

  if (isLoading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        Loading...
      </div>
    )
  }

  const login = async (email: string, password: string) => {
    const response = await apiFetch('/users/login', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ email, password }),
    })
    
    if (response.ok) {
      const userData = await response.json()
      setUser(userData)
      setIsAuthenticated(true)
    } else {
      throw new Error('Authentication failed')
    }
  }

  const logout = async () => {
    await apiFetch('/users/logout', {method: 'POST'})
    setUser(null)
    setIsAuthenticated(false)
  }

  return (
    <AuthContext.Provider value={{ isAuthenticated, user, login, logout }}>
      {children}
    </AuthContext.Provider>
  )
}

export function useAuth() {
  const context = useContext(AuthContext)
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider')
  }
  return context
}
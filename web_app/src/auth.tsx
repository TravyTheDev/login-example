import React, { createContext, useContext, useState, useEffect } from 'react'
import { apiFetch } from './apiFetch'

export type ImpersonatedByUser = {
  id: string
  username: string
  email: string
}

interface User {
  id: string
  username: string
  email: string
  impersonatedBy?: ImpersonatedByUser
}

export interface AuthState {
  isAuthenticated: boolean
  user: User | null
  login: (email: string, password: string) => Promise<void>
  logout: () => void
  loginAs: (email: string) => Promise<void> 
  logoutFromLoginAs: () => Promise<void>
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
        console.log(user)
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
    const response = await apiFetch('/auth/login', {
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
    await apiFetch('/auth/logout', {method: 'POST'})
    setUser(null)
    setIsAuthenticated(false)
  }

  const loginAs = async (email: string) => {
    const response = await apiFetch('/auth/login_as', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ targetEmail: email }),
    })
    if (response.ok) {
      const userData = await response.json()
      console.log(userData, 'hey')
      setUser(userData)
      setIsAuthenticated(true)
    } else {
      throw new Error('loginAs failed')
    }
  }

  const logoutFromLoginAs = async () => {
    const response = await apiFetch('/auth/logout_from_loginAs', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
    })
    
    if (response.ok) {
      const userData = await response.json()
      setUser(userData)
      setIsAuthenticated(true)
    } else {
      throw new Error('Authentication failed')
    }
  }

  return (
    <AuthContext.Provider value={{ isAuthenticated, user, login, logout, loginAs, logoutFromLoginAs }}>
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
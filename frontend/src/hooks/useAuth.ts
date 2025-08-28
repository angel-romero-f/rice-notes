import { useState, useEffect, useCallback } from 'react'
import { useApi } from './useApi'

interface User {
  email: string
  name: string
  picture: string
}

interface AuthState {
  user: User | null
  isAuthenticated: boolean
  isLoading: boolean
  error: string | null
}

export function useAuth() {
  const [authState, setAuthState] = useState<AuthState>({
    user: null,
    isAuthenticated: false,
    isLoading: true,
    error: null
  })

  const { get: apiGet } = useApi<User>()

  const checkAuthStatus = useCallback(async () => {
    try {
      setAuthState(prev => ({ ...prev, isLoading: true, error: null }))
      
      const userData = await apiGet('/api/auth/me')
      
      setAuthState({
        user: userData,
        isAuthenticated: true,
        isLoading: false,
        error: null
      })
    } catch (error) {
      // User is not authenticated or token is invalid
      setAuthState({
        user: null,
        isAuthenticated: false,
        isLoading: false,
        error: error instanceof Error ? error.message : 'Authentication failed'
      })
    }
  }, [apiGet])

  const login = useCallback(() => {
    // Redirect to Google OAuth login endpoint
    window.location.href = `${process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8081'}/api/auth/google`
  }, [])

  const logout = useCallback(async () => {
    try {
      setAuthState(prev => ({ ...prev, isLoading: true }))
      
      // Clear the JWT cookie by calling a logout endpoint (if implemented)
      // For now, we'll just clear the client state and redirect
      setAuthState({
        user: null,
        isAuthenticated: false,
        isLoading: false,
        error: null
      })

      // Use window.location for navigation to home page (static routing)
      window.location.href = '/'
    } catch (error) {
      setAuthState(prev => ({
        ...prev,
        isLoading: false,
        error: error instanceof Error ? error.message : 'Logout failed'
      }))
    }
  }, [])

  // Check authentication status on mount
  useEffect(() => {
    checkAuthStatus()
  }, [checkAuthStatus])

  return {
    ...authState,
    login,
    logout,
    checkAuthStatus
  }
}
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

// JWT token management
const TOKEN_KEY = 'rice_notes_jwt'

const setToken = (token: string) => {
  localStorage.setItem(TOKEN_KEY, token)
}

const getToken = (): string | null => {
  return localStorage.getItem(TOKEN_KEY)
}

const removeToken = () => {
  localStorage.removeItem(TOKEN_KEY)
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
      
      // Check if we have a token from URL (OAuth callback)
      const urlParams = new URLSearchParams(window.location.search)
      const tokenFromUrl = urlParams.get('token')
      
      if (tokenFromUrl) {
        // Store token and clean URL
        setToken(tokenFromUrl)
        window.history.replaceState({}, document.title, window.location.pathname)
      }

      // Check if we have a stored token
      const token = getToken()
      if (!token) {
        setAuthState({
          user: null,
          isAuthenticated: false,
          isLoading: false,
          error: null
        })
        return
      }
      
      const userData = await apiGet('/api/auth/me')
      
      setAuthState({
        user: userData,
        isAuthenticated: true,
        isLoading: false,
        error: null
      })
    } catch (error) {
      // Token is invalid, remove it
      removeToken()
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
    window.location.href = `${process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'}/api/auth/google`
  }, [])

  const logout = useCallback(async () => {
    try {
      setAuthState(prev => ({ ...prev, isLoading: true }))
      
      // Clear the stored token
      removeToken()
      
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
    checkAuthStatus,
    getToken // Export for use by API calls
  }
}
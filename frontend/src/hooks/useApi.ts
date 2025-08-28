import { useState, useCallback } from 'react'

interface ApiOptions {
  method?: 'GET' | 'POST' | 'PUT' | 'DELETE' | 'PATCH'
  headers?: Record<string, string>
  body?: any
}

interface ApiResponse<T> {
  data: T | null
  error: string | null
  loading: boolean
}

interface ApiError {
  error: string
  message?: string
}

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'

export function useApi<T = any>() {
  const [response, setResponse] = useState<ApiResponse<T>>({
    data: null,
    error: null,
    loading: false
  })

  const request = useCallback(async (endpoint: string, options: ApiOptions = {}) => {
    setResponse(prev => ({ ...prev, loading: true, error: null }))

    try {
      const { method = 'GET', headers = {}, body } = options

      const config: RequestInit = {
        method,
        headers: {
          'Content-Type': 'application/json',
          ...headers
        },
        credentials: 'include' // Include cookies for authentication
      }

      if (body && method !== 'GET') {
        config.body = typeof body === 'string' ? body : JSON.stringify(body)
      }

      const url = `${API_BASE_URL}${endpoint}`
      const fetchResponse = await fetch(url, config)

      if (!fetchResponse.ok) {
        let errorMessage = `HTTP ${fetchResponse.status}`
        
        try {
          const errorData: ApiError = await fetchResponse.json()
          errorMessage = errorData.message || errorData.error || errorMessage
        } catch {
          // Fallback to status text if JSON parsing fails
          errorMessage = fetchResponse.statusText || errorMessage
        }

        throw new Error(errorMessage)
      }

      // Handle empty responses (like 204 No Content)
      let data: T | null = null
      const contentType = fetchResponse.headers.get('content-type')
      
      if (contentType && contentType.includes('application/json')) {
        data = await fetchResponse.json()
      }

      setResponse({ data, error: null, loading: false })
      return data

    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : 'An unexpected error occurred'
      setResponse({ data: null, error: errorMessage, loading: false })
      throw error
    }
  }, [])

  const get = useCallback((endpoint: string, options: Omit<ApiOptions, 'method'> = {}) => {
    return request(endpoint, { ...options, method: 'GET' })
  }, [request])

  const post = useCallback((endpoint: string, body?: any, options: Omit<ApiOptions, 'method' | 'body'> = {}) => {
    return request(endpoint, { ...options, method: 'POST', body })
  }, [request])

  const put = useCallback((endpoint: string, body?: any, options: Omit<ApiOptions, 'method' | 'body'> = {}) => {
    return request(endpoint, { ...options, method: 'PUT', body })
  }, [request])

  const del = useCallback((endpoint: string, options: Omit<ApiOptions, 'method'> = {}) => {
    return request(endpoint, { ...options, method: 'DELETE' })
  }, [request])

  const patch = useCallback((endpoint: string, body?: any, options: Omit<ApiOptions, 'method' | 'body'> = {}) => {
    return request(endpoint, { ...options, method: 'PATCH', body })
  }, [request])

  return {
    ...response,
    request,
    get,
    post,
    put,
    delete: del,
    patch
  }
}
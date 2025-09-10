const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'
const TOKEN_KEY = 'rice_notes_jwt'

const getToken = (): string | null => {
  return localStorage.getItem(TOKEN_KEY)
}

export interface Note {
  id: string
  title: string
  course_id: string
  file_name: string
  uploaded_at: string
}

export interface ApiError {
  error: string
  message: string
}

export interface UploadProgress {
  loaded: number
  total: number
  percentage: number
}

/**
 * Upload a note with file and metadata
 */
export async function uploadNote(
  file: File,
  title: string,
  courseId: string,
  onProgress?: (progress: UploadProgress) => void
): Promise<Note> {
  return new Promise((resolve, reject) => {
    const formData = new FormData()
    formData.append('file', file)
    formData.append('title', title)
    formData.append('course_id', courseId)

    const xhr = new XMLHttpRequest()

    // Track upload progress
    if (onProgress) {
      xhr.upload.addEventListener('progress', (event) => {
        if (event.lengthComputable) {
          const progress: UploadProgress = {
            loaded: event.loaded,
            total: event.total,
            percentage: Math.round((event.loaded / event.total) * 100)
          }
          onProgress(progress)
        }
      })
    }

    xhr.onload = () => {
      if (xhr.status >= 200 && xhr.status < 300) {
        try {
          const note: Note = JSON.parse(xhr.responseText)
          resolve(note)
        } catch (error) {
          reject(new Error('Invalid response format'))
        }
      } else {
        try {
          const errorData: ApiError = JSON.parse(xhr.responseText)
          reject(new Error(errorData.message || errorData.error || `HTTP ${xhr.status}`))
        } catch {
          reject(new Error(`Upload failed with status ${xhr.status}`))
        }
      }
    }

    xhr.onerror = () => {
      reject(new Error('Network error occurred during upload'))
    }

    xhr.open('POST', `${API_BASE_URL}/api/notes`)
    xhr.withCredentials = true // Include JWT cookie (fallback)
    
    // Add Authorization header if we have a token
    const token = getToken()
    if (token) {
      xhr.setRequestHeader('Authorization', `Bearer ${token}`)
    }
    
    xhr.send(formData)
  })
}

/**
 * Fetch user's notes
 */
export async function fetchNotes(courseId?: string): Promise<Note[]> {
  const url = new URL(`${API_BASE_URL}/api/notes`)
  if (courseId) {
    url.searchParams.append('course_id', courseId)
  }

  const headers: Record<string, string> = {
    'Content-Type': 'application/json'
  }

  // Add Authorization header if we have a token
  const token = getToken()
  if (token) {
    headers.Authorization = `Bearer ${token}`
  }

  const response = await fetch(url.toString(), {
    method: 'GET',
    credentials: 'include', // fallback for cookies
    headers
  })

  if (!response.ok) {
    const errorData: ApiError = await response.json().catch(() => ({ 
      error: 'unknown_error', 
      message: `HTTP ${response.status}` 
    }))
    throw new Error(errorData.message || errorData.error)
  }

  return response.json()
}

/**
 * Delete a note by ID
 */
export async function deleteNote(noteId: string): Promise<void> {
  const headers: Record<string, string> = {
    'Content-Type': 'application/json'
  }

  // Add Authorization header if we have a token
  const token = getToken()
  if (token) {
    headers.Authorization = `Bearer ${token}`
  }

  const response = await fetch(`${API_BASE_URL}/api/notes/${noteId}`, {
    method: 'DELETE',
    credentials: 'include', // fallback for cookies
    headers
  })

  if (!response.ok) {
    const errorData: ApiError = await response.json().catch(() => ({ 
      error: 'unknown_error', 
      message: `HTTP ${response.status}` 
    }))
    throw new Error(errorData.message || errorData.error)
  }
}

/**
 * Get note download URL
 */
export async function getNoteDownloadUrl(noteId: string): Promise<string> {
  const headers: Record<string, string> = {}

  // Add Authorization header if we have a token
  const token = getToken()
  if (token) {
    headers.Authorization = `Bearer ${token}`
  }

  const response = await fetch(`${API_BASE_URL}/api/notes/${noteId}`, {
    method: 'GET',
    credentials: 'include', // fallback for cookies
    redirect: 'manual', // Don't follow redirects automatically
    headers
  })

  if (response.status === 302) {
    const location = response.headers.get('Location')
    if (location) {
      return location
    }
  }

  if (!response.ok) {
    const errorData: ApiError = await response.json().catch(() => ({ 
      error: 'unknown_error', 
      message: `HTTP ${response.status}` 
    }))
    throw new Error(errorData.message || errorData.error)
  }

  throw new Error('No redirect URL found')
}

/**
 * Validate file before upload
 */
export function validateFile(file: File): { valid: boolean; error?: string } {
  // Check file type
  if (file.type !== 'application/pdf') {
    return { valid: false, error: 'Only PDF files are allowed' }
  }

  // Check file size (10MB limit)
  const maxSize = 10 * 1024 * 1024 // 10MB in bytes
  if (file.size > maxSize) {
    return { valid: false, error: 'File size must be less than 10MB' }
  }

  return { valid: true }
}

/**
 * Format file size for display
 */
export function formatFileSize(bytes: number): string {
  if (bytes === 0) return '0 Bytes'

  const k = 1024
  const sizes = ['Bytes', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))

  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}
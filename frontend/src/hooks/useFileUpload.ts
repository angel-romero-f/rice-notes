import { useState, useCallback } from 'react'
import { uploadNote, validateFile, type Note, type UploadProgress } from '@/lib/api'

export interface UploadState {
  isUploading: boolean
  progress: UploadProgress | null
  error: string | null
  uploadedNote: Note | null
}

export interface UseFileUploadReturn {
  uploadState: UploadState
  upload: (file: File, title: string, courseId: string) => Promise<Note>
  reset: () => void
  validateFile: (file: File) => { valid: boolean; error?: string }
}

export function useFileUpload(): UseFileUploadReturn {
  const [uploadState, setUploadState] = useState<UploadState>({
    isUploading: false,
    progress: null,
    error: null,
    uploadedNote: null
  })

  const reset = useCallback(() => {
    setUploadState({
      isUploading: false,
      progress: null,
      error: null,
      uploadedNote: null
    })
  }, [])

  const upload = useCallback(async (file: File, title: string, courseId: string): Promise<Note> => {
    // Validate inputs
    if (!file || !title.trim() || !courseId.trim()) {
      throw new Error('File, title, and course ID are required')
    }

    // Validate file
    const fileValidation = validateFile(file)
    if (!fileValidation.valid) {
      throw new Error(fileValidation.error)
    }

    // Reset state and start upload
    setUploadState({
      isUploading: true,
      progress: { loaded: 0, total: file.size, percentage: 0 },
      error: null,
      uploadedNote: null
    })

    try {
      const note = await uploadNote(
        file,
        title.trim(),
        courseId.trim(),
        (progress) => {
          setUploadState(prev => ({
            ...prev,
            progress
          }))
        }
      )

      setUploadState({
        isUploading: false,
        progress: { loaded: file.size, total: file.size, percentage: 100 },
        error: null,
        uploadedNote: note
      })

      return note
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : 'Upload failed'
      
      setUploadState({
        isUploading: false,
        progress: null,
        error: errorMessage,
        uploadedNote: null
      })

      throw error
    }
  }, [])

  return {
    uploadState,
    upload,
    reset,
    validateFile
  }
}
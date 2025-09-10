'use client'

import { useState, useRef, useCallback, useEffect } from 'react'
import { useFileUpload } from '@/hooks/useFileUpload'
import { formatFileSize } from '@/lib/api'
import CourseSelect from './CourseSelect'

interface UploadModalProps {
  isOpen: boolean
  onClose: () => void
  onSuccess?: (noteId: string) => void
}

const styles = {
  overlay: 'fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4',
  modal: 'bg-white rounded-lg shadow-xl w-full max-w-md mx-auto max-h-[90vh] overflow-y-auto',
  header: 'flex items-center justify-between p-6 border-b border-gray-200',
  title: 'text-xl font-semibold text-gray-900',
  closeButton: 'text-gray-400 hover:text-gray-600 transition-colors',
  content: 'p-6 space-y-6',
  
  // File upload area
  uploadArea: 'border-2 border-dashed border-gray-300 rounded-lg p-8 text-center transition-colors hover:border-gray-400',
  uploadAreaDragOver: 'border-blue-500 bg-blue-50',
  uploadAreaError: 'border-red-300 bg-red-50',
  uploadIcon: 'w-12 h-12 mx-auto mb-4 text-gray-400',
  uploadText: 'text-lg font-medium text-gray-900 mb-2',
  uploadSubtext: 'text-sm text-gray-500 mb-4',
  browseButton: 'text-blue-600 hover:text-blue-700 font-medium cursor-pointer',
  
  // Selected file
  selectedFile: 'bg-gray-50 border border-gray-200 rounded-lg p-4 flex items-center justify-between',
  fileInfo: 'flex items-center space-x-3',
  fileIcon: 'w-8 h-8 text-red-600',
  fileName: 'font-medium text-gray-900',
  fileSize: 'text-sm text-gray-500',
  removeButton: 'text-red-600 hover:text-red-700 font-medium text-sm',
  
  // Form fields
  fieldGroup: 'space-y-2',
  label: 'block text-sm font-medium text-gray-700',
  requiredMark: 'text-red-500 ml-1',
  input: 'w-full px-4 py-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent',
  inputError: 'border-red-500 focus:ring-red-500',
  errorMessage: 'mt-1 text-sm text-red-600',
  
  // Progress
  progressContainer: 'space-y-3',
  progressBar: 'w-full bg-gray-200 rounded-full h-2',
  progressFill: 'h-2 bg-blue-600 rounded-full transition-all duration-300',
  progressText: 'text-sm text-gray-600 text-center',
  
  // Buttons
  buttonGroup: 'flex space-x-3 pt-4 border-t border-gray-200',
  cancelButton: 'flex-1 px-4 py-3 border border-gray-300 text-gray-700 rounded-lg hover:bg-gray-50 font-medium transition-colors',
  uploadButton: 'flex-1 px-4 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:bg-gray-300 disabled:cursor-not-allowed font-medium transition-colors',
  
  // Success state
  successContainer: 'text-center py-8',
  successIcon: 'w-16 h-16 mx-auto mb-4 text-green-500',
  successTitle: 'text-xl font-semibold text-gray-900 mb-2',
  successMessage: 'text-gray-600 mb-6'
}

export default function UploadModal({ isOpen, onClose, onSuccess }: UploadModalProps) {
  const [selectedFile, setSelectedFile] = useState<File | null>(null)
  const [title, setTitle] = useState('')
  const [courseId, setCourseId] = useState('')
  const [dragOver, setDragOver] = useState(false)
  const [errors, setErrors] = useState<{ title?: string; courseId?: string; file?: string }>({})
  
  const fileInputRef = useRef<HTMLInputElement>(null)
  const { uploadState, upload, reset, validateFile } = useFileUpload()

  // Reset form when modal closes
  useEffect(() => {
    if (!isOpen) {
      setSelectedFile(null)
      setTitle('')
      setCourseId('')
      setErrors({})
      reset()
    }
  }, [isOpen, reset])

  // Handle escape key
  useEffect(() => {
    const handleEscape = (e: KeyboardEvent) => {
      if (e.key === 'Escape' && isOpen && !uploadState.isUploading) {
        onClose()
      }
    }

    if (isOpen) {
      document.addEventListener('keydown', handleEscape)
      return () => document.removeEventListener('keydown', handleEscape)
    }
  }, [isOpen, onClose, uploadState.isUploading])

  const handleFileSelect = useCallback((files: FileList | null) => {
    if (!files || files.length === 0) return

    const file = files[0]
    const validation = validateFile(file)
    
    if (!validation.valid) {
      setErrors(prev => ({ ...prev, file: validation.error }))
      return
    }

    setSelectedFile(file)
    setErrors(prev => ({ ...prev, file: undefined }))

    // Auto-generate title from filename if title is empty
    if (!title) {
      const fileNameWithoutExt = file.name.replace(/\.[^/.]+$/, '')
      setTitle(fileNameWithoutExt)
    }
  }, [title, validateFile])

  const handleDragOver = useCallback((e: React.DragEvent) => {
    e.preventDefault()
    setDragOver(true)
  }, [])

  const handleDragLeave = useCallback((e: React.DragEvent) => {
    e.preventDefault()
    setDragOver(false)
  }, [])

  const handleDrop = useCallback((e: React.DragEvent) => {
    e.preventDefault()
    setDragOver(false)
    handleFileSelect(e.dataTransfer.files)
  }, [handleFileSelect])

  const handleFileInputChange = useCallback((e: React.ChangeEvent<HTMLInputElement>) => {
    handleFileSelect(e.target.files)
  }, [handleFileSelect])

  const validateForm = (): boolean => {
    const newErrors: { title?: string; courseId?: string; file?: string } = {}

    if (!title.trim()) {
      newErrors.title = 'Title is required'
    }

    if (!courseId.trim()) {
      newErrors.courseId = 'Course ID is required'
    }

    if (!selectedFile) {
      newErrors.file = 'Please select a PDF file'
    }

    setErrors(newErrors)
    return Object.keys(newErrors).length === 0
  }

  const handleUpload = async () => {
    if (!validateForm() || !selectedFile) return

    try {
      const note = await upload(selectedFile, title, courseId)
      onSuccess?.(note.id)
      onClose()
    } catch (error) {
      // Error is handled by the upload hook
      console.error('Upload failed:', error)
    }
  }

  if (!isOpen) return null

  // Success state
  if (uploadState.uploadedNote && !uploadState.isUploading) {
    return (
      <div className={styles.overlay} onClick={onClose}>
        <div className={styles.modal} onClick={e => e.stopPropagation()}>
          <div className={styles.successContainer}>
            <svg className={styles.successIcon} fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            <h3 className={styles.successTitle}>Upload Successful!</h3>
            <p className={styles.successMessage}>
              Your note "{uploadState.uploadedNote.title}" has been uploaded successfully.
            </p>
            <button
              onClick={onClose}
              className={styles.uploadButton}
            >
              Done
            </button>
          </div>
        </div>
      </div>
    )
  }

  return (
    <div className={styles.overlay} onClick={onClose}>
      <div className={styles.modal} onClick={e => e.stopPropagation()}>
        {/* Header */}
        <div className={styles.header}>
          <h2 className={styles.title}>Upload Notes</h2>
          <button
            onClick={onClose}
            className={styles.closeButton}
            disabled={uploadState.isUploading}
            aria-label="Close modal"
          >
            <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>

        {/* Content */}
        <div className={styles.content}>
          {/* File Upload Area */}
          <div className={styles.fieldGroup}>
            <label className={styles.label}>
              PDF File<span className={styles.requiredMark}>*</span>
            </label>
            
            {!selectedFile ? (
              <div
                className={`${styles.uploadArea} ${dragOver ? styles.uploadAreaDragOver : ''} ${errors.file ? styles.uploadAreaError : ''}`}
                onDragOver={handleDragOver}
                onDragLeave={handleDragLeave}
                onDrop={handleDrop}
                onClick={() => fileInputRef.current?.click()}
              >
                <svg className={styles.uploadIcon} fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12" />
                </svg>
                <div className={styles.uploadText}>Drop your PDF here</div>
                <div className={styles.uploadSubtext}>
                  or{' '}
                  <span className={styles.browseButton}>browse files</span>
                </div>
                <div className={styles.uploadSubtext}>Maximum file size: 10MB</div>
              </div>
            ) : (
              <div className={styles.selectedFile}>
                <div className={styles.fileInfo}>
                  <svg className={styles.fileIcon} fill="currentColor" viewBox="0 0 24 24">
                    <path d="M14,2H6A2,2 0 0,0 4,4V20A2,2 0 0,0 6,22H18A2,2 0 0,0 20,20V8L14,2M18,20H6V4H13V9H18V20Z" />
                  </svg>
                  <div>
                    <div className={styles.fileName}>{selectedFile.name}</div>
                    <div className={styles.fileSize}>{formatFileSize(selectedFile.size)}</div>
                  </div>
                </div>
                <button
                  onClick={() => setSelectedFile(null)}
                  className={styles.removeButton}
                  disabled={uploadState.isUploading}
                >
                  Remove
                </button>
              </div>
            )}
            
            <input
              ref={fileInputRef}
              type="file"
              accept=".pdf,application/pdf"
              onChange={handleFileInputChange}
              className="hidden"
            />
            
            {errors.file && <div className={styles.errorMessage}>{errors.file}</div>}
          </div>

          {/* Title Field */}
          <div className={styles.fieldGroup}>
            <label htmlFor="title" className={styles.label}>
              Title<span className={styles.requiredMark}>*</span>
            </label>
            <input
              id="title"
              type="text"
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              placeholder="Enter a descriptive title for your notes"
              className={`${styles.input} ${errors.title ? styles.inputError : ''}`}
              disabled={uploadState.isUploading}
            />
            {errors.title && <div className={styles.errorMessage}>{errors.title}</div>}
          </div>

          {/* Course Select */}
          <div className={styles.fieldGroup}>
            <label htmlFor="courseId" className={styles.label}>
              Course<span className={styles.requiredMark}>*</span>
            </label>
            <CourseSelect
              value={courseId}
              onChange={setCourseId}
              error={errors.courseId}
              placeholder="Select or enter course code (e.g., COMP101)"
            />
          </div>

          {/* Upload Progress */}
          {uploadState.isUploading && uploadState.progress && (
            <div className={styles.progressContainer}>
              <div className={styles.progressBar}>
                <div 
                  className={styles.progressFill}
                  style={{ width: `${uploadState.progress.percentage}%` }}
                />
              </div>
              <div className={styles.progressText}>
                Uploading... {uploadState.progress.percentage}%
              </div>
            </div>
          )}

          {/* Error Message */}
          {uploadState.error && (
            <div className={styles.errorMessage}>{uploadState.error}</div>
          )}

          {/* Action Buttons */}
          <div className={styles.buttonGroup}>
            <button
              onClick={onClose}
              className={styles.cancelButton}
              disabled={uploadState.isUploading}
            >
              Cancel
            </button>
            <button
              onClick={handleUpload}
              className={styles.uploadButton}
              disabled={uploadState.isUploading || !selectedFile || !title.trim() || !courseId.trim()}
            >
              {uploadState.isUploading ? 'Uploading...' : 'Upload Note'}
            </button>
          </div>
        </div>
      </div>
    </div>
  )
}
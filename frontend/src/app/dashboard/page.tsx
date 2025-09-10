'use client'

import { useState, useEffect } from 'react'
import { useAuth } from '@/hooks/useAuth'
import { mockNotes } from '@/components/MockNotes'
import NoteCard from '@/components/ui/NoteCard'
import UploadModal from '@/components/ui/UploadModal'
import { fetchNotes, type Note } from '@/lib/api'
import { redirect } from 'next/navigation'

const styles = {
  container: 'min-h-screen bg-gray-50',
  header: 'bg-white border-b border-gray-200 shadow-sm',
  headerContent: 'max-w-7xl mx-auto px-4 sm:px-6 lg:px-8',
  headerTop: 'flex justify-between items-center py-4',
  logo: 'text-2xl font-bold text-blue-600 flex items-center gap-2',
  userSection: 'flex items-center gap-4',
  userInfo: 'flex items-center gap-3',
  userAvatar: 'w-8 h-8 rounded-full',
  userName: 'text-gray-700 font-medium',
  logoutButton: 'text-gray-500 hover:text-gray-700 px-3 py-1 rounded border border-gray-300 hover:bg-gray-50 transition-colors',
  searchSection: 'py-4 border-t border-gray-100',
  searchContainer: 'flex items-center gap-4',
  searchInput: 'flex-1 px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500',
  uploadButton: 'bg-blue-600 text-white px-6 py-2 rounded-lg hover:bg-blue-700 transition-colors font-medium',
  main: 'max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8',
  welcomeSection: 'mb-8',
  welcomeTitle: 'text-3xl font-bold text-gray-900 mb-2',
  welcomeSubtitle: 'text-gray-600',
  notesSection: '',
  sectionTitle: 'text-xl font-semibold text-gray-900 mb-6',
  notesGrid: 'grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6',
  loadingContainer: 'flex justify-center items-center min-h-screen',
  loadingSpinner: 'animate-spin h-8 w-8 border-2 border-blue-600 border-t-transparent rounded-full',
  errorContainer: 'flex justify-center items-center min-h-screen',
  errorMessage: 'text-red-600 text-lg',
  
  // Notification styles
  notification: 'fixed top-4 right-4 z-50 max-w-md bg-white border border-gray-200 rounded-lg shadow-lg p-4',
  notificationSuccess: 'border-green-500 bg-green-50',
  notificationIcon: 'w-5 h-5 mr-3 text-green-500',
  notificationContent: 'flex items-center',
  notificationText: 'text-sm font-medium text-green-900',
  notificationClose: 'ml-auto text-green-500 hover:text-green-700 cursor-pointer'
}

export default function DashboardPage() {
  const { user, isAuthenticated, isLoading, logout } = useAuth()
  const [isUploadModalOpen, setIsUploadModalOpen] = useState(false)
  const [userNotes, setUserNotes] = useState<Note[]>([])
  const [notesLoading, setNotesLoading] = useState(false)
  const [notesError, setNotesError] = useState<string | null>(null)
  const [showNotification, setShowNotification] = useState(false)
  const [notificationMessage, setNotificationMessage] = useState('')

  // Load user's notes on authentication
  useEffect(() => {
    if (isAuthenticated && user) {
      loadUserNotes()
    }
  }, [isAuthenticated, user])

  const loadUserNotes = async () => {
    try {
      setNotesLoading(true)
      setNotesError(null)
      const notes = await fetchNotes()
      setUserNotes(notes)
    } catch (error) {
      setNotesError(error instanceof Error ? error.message : 'Failed to load notes')
      console.error('Failed to load notes:', error)
    } finally {
      setNotesLoading(false)
    }
  }

  const showSuccessNotification = (message: string) => {
    setNotificationMessage(message)
    setShowNotification(true)
    setTimeout(() => setShowNotification(false), 5000)
  }

  const handleUploadSuccess = (noteId: string) => {
    // Refresh the notes list
    loadUserNotes()
    showSuccessNotification('Note uploaded successfully!')
  }

  const handleUploadClick = () => {
    setIsUploadModalOpen(true)
  }

  if (isLoading) {
    return (
      <div className={styles.loadingContainer}>
        <div className={styles.loadingSpinner}></div>
      </div>
    )
  }

  // Redirect unauthenticated users using Next.js redirect
  if (!isAuthenticated || !user) {
    redirect('/')
  }

  return (
    <div className={styles.container}>
      {/* Header */}
      <header className={styles.header}>
        <div className={styles.headerContent}>
          <div className={styles.headerTop}>
            <div className={styles.logo}>
              <span>📝</span>
              <span>ricenotes</span>
            </div>
            
            <div className={styles.userSection}>
              <div className={styles.userInfo}>
                {user.picture && (
                  <img 
                    src={user.picture} 
                    alt={user.name || 'User'}
                    className={styles.userAvatar}
                  />
                )}
                <span className={styles.userName}>{user.name || 'User'}</span>
              </div>
              <button 
                onClick={logout}
                className={styles.logoutButton}
              >
                Sign Out
              </button>
            </div>
          </div>

          <div className={styles.searchSection}>
            <div className={styles.searchContainer}>
              <input
                type="text"
                placeholder="Search for courses, topics, or notes..."
                className={styles.searchInput}
              />
              <button 
                onClick={handleUploadClick}
                className={styles.uploadButton}
              >
                Upload Notes
              </button>
            </div>
          </div>
        </div>
      </header>

      {/* Main Content */}
      <main className={styles.main}>
        <div className={styles.welcomeSection}>
          <h1 className={styles.welcomeTitle}>
            Welcome back, {user.name ? user.name.split(' ')[0] : 'User'}! 👋
          </h1>
          <p className={styles.welcomeSubtitle}>
            Discover study materials shared by your fellow Rice students
          </p>
        </div>

        {/* User's Notes Section */}
        {userNotes && userNotes.length > 0 && (
          <div className={styles.notesSection}>
            <h2 className={styles.sectionTitle}>Your Notes</h2>
            {notesLoading ? (
              <div className="flex justify-center py-8">
                <div className={styles.loadingSpinner}></div>
              </div>
            ) : notesError ? (
              <div className={styles.errorMessage}>{notesError}</div>
            ) : (
              <div className={styles.notesGrid}>
                {userNotes?.map((note) => (
                  <NoteCard 
                    key={note.id} 
                    note={{
                      id: note.id,
                      title: note.title,
                      course: note.course_id,
                      author: user?.name || 'You',
                      rating: 5, // Default rating for user's own notes
                      thumbnailUrl: `data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='200' height='240' viewBox='0 0 200 240' fill='%23f3f4f6'%3E%3Crect width='200' height='240' fill='%23f3f4f6'/%3E%3Cpath d='M50 80h100v8H50zM50 100h100v8H50zM50 120h100v8H50zM50 140h80v8H50z' fill='%23d1d5db'/%3E%3Ctext x='100' y='180' text-anchor='middle' font-family='sans-serif' font-size='14' fill='%239ca3af'%3EPDF%3C/text%3E%3C/svg%3E` // SVG placeholder
                    }} 
                    showHeart={false} 
                  />
                ))}
              </div>
            )}
          </div>
        )}

        {/* Popular Notes Section */}
        <div className={styles.notesSection}>
          <h2 className={styles.sectionTitle}>
            {userNotes && userNotes.length > 0 ? 'Popular Notes' : 'Featured Notes'}
          </h2>
          <div className={styles.notesGrid}>
            {mockNotes.map((note) => (
              <NoteCard key={note.id} note={note} showHeart={true} />
            ))}
          </div>
        </div>
      </main>

      {/* Upload Modal */}
      <UploadModal
        isOpen={isUploadModalOpen}
        onClose={() => setIsUploadModalOpen(false)}
        onSuccess={handleUploadSuccess}
      />

      {/* Success Notification */}
      {showNotification && (
        <div className={`${styles.notification} ${styles.notificationSuccess}`}>
          <div className={styles.notificationContent}>
            <svg className={styles.notificationIcon} fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            <span className={styles.notificationText}>{notificationMessage}</span>
            <button
              onClick={() => setShowNotification(false)}
              className={styles.notificationClose}
            >
              <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>
        </div>
      )}
    </div>
  )
}
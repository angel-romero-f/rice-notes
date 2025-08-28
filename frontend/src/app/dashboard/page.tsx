'use client'

import { useAuth } from '@/hooks/useAuth'
import { mockNotes } from '@/components/MockNotes'
import NoteCard from '@/components/ui/NoteCard'
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
  errorMessage: 'text-red-600 text-lg'
}

export default function DashboardPage() {
  const { user, isAuthenticated, isLoading, logout } = useAuth()

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
                    alt={user.name}
                    className={styles.userAvatar}
                  />
                )}
                <span className={styles.userName}>{user.name}</span>
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
              <button className={styles.uploadButton}>
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
            Welcome back, {user.name.split(' ')[0]}! 👋
          </h1>
          <p className={styles.welcomeSubtitle}>
            Discover study materials shared by your fellow Rice students
          </p>
        </div>

        <div className={styles.notesSection}>
          <h2 className={styles.sectionTitle}>Popular Notes</h2>
          <div className={styles.notesGrid}>
            {mockNotes.map((note) => (
              <NoteCard key={note.id} note={note} showHeart={true} />
            ))}
          </div>
        </div>
      </main>
    </div>
  )
}
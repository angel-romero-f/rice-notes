'use client'

import React from 'react'
import { redirect } from 'next/navigation'
import NoteCard from '@/components/ui/NoteCard'
import GoogleSignIn from '@/components/auth/GoogleSignIn'
import { mockNotes } from '@/components/MockNotes'
import { useAuth } from '@/hooks/useAuth'

const styles = {
  container: 'h-screen flex overflow-hidden',
  leftSection: 'flex-[1.2] relative bg-gradient-to-br from-blue-500 via-blue-600 to-indigo-700 text-white px-8 lg:px-12 py-6 flex flex-col items-center justify-center',
  logo: 'absolute top-6 left-8 lg:left-12 text-white font-bold text-base flex items-center gap-2',
  heroSection: 'flex flex-col items-center justify-center h-full gap-8 text-center',
  heroContent: 'max-w-4xl',
  heroTitle: 'text-2xl lg:text-3xl font-bold mb-3 leading-none whitespace-nowrap',
  heroSubtitle: 'text-base lg:text-lg text-blue-100 whitespace-nowrap',
  mockupContainer: 'flex justify-center',
  mockupPreview: 'bg-white/15 backdrop-blur-md rounded-xl p-3 border border-white/20 shadow-xl max-w-lg',
  mockupHeader: 'flex items-center justify-between mb-3',
  mockupLogo: 'text-white font-bold text-sm flex items-center gap-1',
  mockupSearch: 'flex-1 mx-3 bg-white/20 rounded px-2 py-1.5 text-white placeholder-blue-200 text-sm border border-white/20',
  mockupButton: 'text-white text-sm bg-blue-500 px-2 py-1.5 rounded hover:bg-blue-400 transition-colors',
  mockupGrid: 'grid grid-cols-2 gap-2',
  rightSection: 'flex-[0.8] bg-white flex flex-col justify-center px-6 lg:px-8 py-6',
  welcomeSection: 'max-w-md mx-auto text-center',
  welcomeTitle: 'text-2xl lg:text-3xl font-bold text-gray-900 mb-4 flex items-center justify-center gap-2',
  welcomeSubtitle: 'text-gray-600 mb-8 leading-relaxed',
  riceHighlight: 'text-blue-600 font-semibold'
}


export default function HomePage() {
  const { isAuthenticated, isLoading } = useAuth()

  // Show loading state while checking authentication
  if (isLoading) {
    return (
      <div className="h-screen flex items-center justify-center">
        <div className="animate-spin h-8 w-8 border-2 border-blue-600 border-t-transparent rounded-full"></div>
      </div>
    )
  }

  // Redirect authenticated users using Next.js redirect
  if (isAuthenticated) {
    redirect('/dashboard')
  }

  // Show landing page for unauthenticated users
  return (
    <div className={styles.container}>
      {/* Left Section - Hero */}
      <div className={styles.leftSection}>
        {/* Logo in top left */}
        <div className={styles.logo}>
          <span>📝</span>
          <span>ricenotes</span>
        </div>

        <div className={styles.heroSection}>
          {/* Hero content in upper left */}
          <div className={styles.heroContent}>
            <h1 className={styles.heroTitle}>
              Find your perfect study notes at <span className="text-blue-100 font-extrabold">Rice</span>.
            </h1>
            <p className={styles.heroSubtitle}>
              Connect with Rice students & discover ideal study materials!
            </p>
          </div>

          {/* Mockup Preview in lower section */}
          <div className={styles.mockupContainer}>
            <div className={styles.mockupPreview}>
              <div className={styles.mockupHeader}>
                <div className={styles.mockupLogo}>
                  <span>📝</span>
                  <span>ricenotes</span>
                </div>
                <input 
                  className={styles.mockupSearch}
                  placeholder="Search for courses, topics..."
                  readOnly
                />
                <button className={styles.mockupButton}>
                  Post a Listing
                </button>
              </div>
              <div className={styles.mockupGrid}>
                {mockNotes.map((note) => (
                  <NoteCard key={note.id} note={note} showHeart={true} />
                ))}
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* Right Section - Auth */}
      <div className={styles.rightSection}>
        <div className={styles.welcomeSection}>
          <h2 className={styles.welcomeTitle}>
            <span>📚</span>
            Welcome to <span className={styles.riceHighlight}>ricenotes</span>
          </h2>
          <p className={styles.welcomeSubtitle}>
            Please sign in through your Rice Google account to access your account!
          </p>
          <GoogleSignIn />
        </div>
      </div>
    </div>
  )
}
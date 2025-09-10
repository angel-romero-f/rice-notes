'use client'

import React from 'react'
import Link from 'next/link'

const styles = {
  container: 'h-screen flex overflow-hidden',
  leftSection: 'flex-[1.2] relative bg-gradient-to-br from-red-500 via-red-600 to-pink-700 text-white px-8 lg:px-12 py-6 flex flex-col items-center justify-center',
  logo: 'absolute top-6 left-8 lg:left-12 text-white font-bold text-base flex items-center gap-2',
  heroSection: 'flex flex-col items-center justify-center h-full gap-8 text-center',
  heroContent: 'max-w-4xl',
  heroTitle: 'text-2xl lg:text-3xl font-bold mb-3 leading-tight',
  heroSubtitle: 'text-base lg:text-lg text-red-100',
  iconContainer: 'mb-6',
  icon: 'text-6xl',
  rightSection: 'flex-[0.8] bg-white flex flex-col justify-center px-6 lg:px-8 py-6',
  errorSection: 'max-w-md mx-auto text-center',
  errorTitle: 'text-2xl lg:text-3xl font-bold text-gray-900 mb-4 flex items-center justify-center gap-2',
  errorSubtitle: 'text-gray-600 mb-6 leading-relaxed',
  errorDetails: 'text-gray-500 text-sm mb-8 leading-relaxed',
  riceHighlight: 'text-red-600 font-semibold',
  emailHighlight: 'font-mono text-blue-600 bg-blue-50 px-2 py-1 rounded',
  backButton: 'inline-flex items-center gap-2 px-6 py-3 bg-blue-600 text-white font-semibold rounded-lg hover:bg-blue-700 transition-colors',
  backIcon: 'text-lg'
}

export default function UnauthorizedPage() {
  return (
    <div className={styles.container}>
      {/* Left Section - Error Illustration */}
      <div className={styles.leftSection}>
        {/* Logo in top left */}
        <div className={styles.logo}>
          <span>📝</span>
          <span>ricenotes</span>
        </div>

        <div className={styles.heroSection}>
          <div className={styles.heroContent}>
            <div className={styles.iconContainer}>
              <div className={styles.icon}>🚫</div>
            </div>
            <h1 className={styles.heroTitle}>
              Access Restricted
            </h1>
            <p className={styles.heroSubtitle}>
              Only Rice University students can access this platform
            </p>
          </div>
        </div>
      </div>

      {/* Right Section - Error Details */}
      <div className={styles.rightSection}>
        <div className={styles.errorSection}>
          <h2 className={styles.errorTitle}>
            <span>🏫</span>
            <span className={styles.riceHighlight}>Rice Students Only</span>
          </h2>
          <p className={styles.errorSubtitle}>
            We're sorry, but <strong>ricenotes</strong> is exclusively available to Rice University students.
          </p>
          <p className={styles.errorDetails}>
            To access this platform, you need to sign in with a valid Rice University Google account ending with <span className={styles.emailHighlight}>@rice.edu</span>
          </p>
          <p className={styles.errorDetails}>
            If you're a Rice student and believe this is an error, please contact your IT administrator or try signing in with your official Rice email address.
          </p>
          <Link href="/" className={styles.backButton}>
            <span className={styles.backIcon}>←</span>
            Back to Sign In
          </Link>
        </div>
      </div>
    </div>
  )
}
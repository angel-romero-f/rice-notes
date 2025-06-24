export interface Course {
  id: string
  code: string
  name: string
  department: string
}

export interface User {
  id: string
  email: string
  name: string
  profileImage?: string
}

export interface Note {
  id: string
  title: string
  description: string
  course: Course
  uploadedBy: User
  uploadedAt: Date
  fileType: 'pdf' | 'image' | 'markdown'
  thumbnailUrl?: string
  tags: string[]
  rating: number
  downloadCount: number
}

export interface MockNotePreview {
  id: string
  title: string
  course: string
  author: string
  rating: number
  price?: string
  thumbnailUrl: string
}
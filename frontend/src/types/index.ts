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
  description?: string
  course?: Course
  course_id?: string // Backend format
  uploadedBy?: User
  uploadedAt?: Date
  uploaded_at?: string // Backend format
  fileType?: 'pdf' | 'image' | 'markdown'
  file_name?: string // Backend format
  thumbnailUrl?: string
  tags?: string[]
  rating?: number
  downloadCount?: number
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
import { MockNotePreview } from '@/types'

interface NoteCardProps {
  note: MockNotePreview
  showHeart?: boolean
}

const styles = {
  card: 'bg-white rounded shadow-sm border border-gray-200 overflow-hidden hover:shadow-md transition-shadow',
  imageContainer: 'relative aspect-[4/3] bg-gradient-to-br from-blue-50 to-blue-100',
  image: 'w-full h-full object-cover',
  content: 'p-1.5',
  title: 'font-semibold text-gray-900 text-xs mb-0.5',
  course: 'text-xs text-gray-500 mb-0.5',
  meta: 'flex items-center justify-between text-xs text-gray-400',
  rating: 'flex items-center gap-1',
  heart: 'absolute top-1 right-1 w-3 h-3 text-gray-400 hover:text-red-500 transition-colors cursor-pointer',
  placeholder: 'w-full h-full flex items-center justify-center text-blue-400 font-medium text-xs'
}

export default function NoteCard({ note, showHeart = true }: NoteCardProps) {
  return (
    <div className={styles.card}>
      <div className={styles.imageContainer}>
        {note.thumbnailUrl ? (
          <img 
            src={note.thumbnailUrl} 
            alt={note.title}
            className={styles.image}
          />
        ) : (
          <div className={styles.placeholder}>
            📝 {note.course}
          </div>
        )}
        {showHeart && (
          <svg className={styles.heart} fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z" />
          </svg>
        )}
      </div>
      <div className={styles.content}>
        <h3 className={styles.title}>{note.title}</h3>
        <p className={styles.course}>{note.course}</p>
        <div className={styles.meta}>
          <span>{note.author}</span>
          <div className={styles.rating}>
            <span>⭐</span>
            <span>{note.rating}</span>
          </div>
        </div>
        {note.price && (
          <div className="mt-2 text-sm font-semibold text-blue-600">
            {note.price}
          </div>
        )}
      </div>
    </div>
  )
}
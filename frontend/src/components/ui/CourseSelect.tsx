'use client'

import { useState, useRef, useEffect } from 'react'

export interface Course {
  code: string
  name: string
  department: string
}

interface CourseSelectProps {
  value: string
  onChange: (courseId: string) => void
  error?: string
  placeholder?: string
  className?: string
}

// Common Rice University courses
const RICE_COURSES: Course[] = [
  // Computer Science
  { code: 'COMP101', name: 'Owls, Machines and Algorithms', department: 'COMP' },
  { code: 'COMP140', name: 'Computational Thinking', department: 'COMP' },
  { code: 'COMP160', name: 'Introduction to Programming', department: 'COMP' },
  { code: 'COMP182', name: 'Algorithmic Thinking', department: 'COMP' },
  { code: 'COMP215', name: 'Introduction to Program Design', department: 'COMP' },
  { code: 'COMP280', name: 'Discrete Mathematics', department: 'COMP' },
  { code: 'COMP310', name: 'Data Structures and Algorithms', department: 'COMP' },
  { code: 'COMP321', name: 'Introduction to Computer Systems', department: 'COMP' },
  { code: 'COMP382', name: 'Reasoning about Algorithms', department: 'COMP' },
  { code: 'COMP421', name: 'Operating Systems and Concurrent Programming', department: 'COMP' },
  { code: 'COMP450', name: 'Algorithm Design and Analysis', department: 'COMP' },
  
  // Mathematics
  { code: 'MATH101', name: 'Single Variable Calculus I', department: 'MATH' },
  { code: 'MATH102', name: 'Single Variable Calculus II', department: 'MATH' },
  { code: 'MATH211', name: 'Ordinary Differential Equations', department: 'MATH' },
  { code: 'MATH212', name: 'Multivariable Calculus', department: 'MATH' },
  { code: 'MATH220', name: 'Linear Algebra', department: 'MATH' },
  { code: 'MATH355', name: 'Linear Algebra', department: 'MATH' },
  { code: 'MATH381', name: 'Probability Theory', department: 'MATH' },
  { code: 'MATH382', name: 'Statistics', department: 'MATH' },
  
  // Physics
  { code: 'PHYS101', name: 'Mechanics', department: 'PHYS' },
  { code: 'PHYS102', name: 'Electricity and Magnetism', department: 'PHYS' },
  { code: 'PHYS201', name: 'Vibrations, Waves, and Optics', department: 'PHYS' },
  { code: 'PHYS202', name: 'Thermal Physics and Quantum Mechanics', department: 'PHYS' },
  
  // Chemistry
  { code: 'CHEM121', name: 'General Chemistry I', department: 'CHEM' },
  { code: 'CHEM122', name: 'General Chemistry II', department: 'CHEM' },
  { code: 'CHEM211', name: 'Organic Chemistry I', department: 'CHEM' },
  { code: 'CHEM212', name: 'Organic Chemistry II', department: 'CHEM' },
  
  // Economics
  { code: 'ECON100', name: 'Introduction to Economics', department: 'ECON' },
  { code: 'ECON211', name: 'Principles of Microeconomics', department: 'ECON' },
  { code: 'ECON212', name: 'Principles of Macroeconomics', department: 'ECON' },
  
  // English
  { code: 'ENGL103', name: 'First-Year Writing Intensive Seminar', department: 'ENGL' },
  { code: 'ENGL291', name: 'Intermediate Fiction Writing', department: 'ENGL' },
  { code: 'ENGL292', name: 'Intermediate Poetry Writing', department: 'ENGL' },
  
  // Other common courses
  { code: 'BIOC201', name: 'Introduction to Biochemistry and Cell Biology', department: 'BIOC' },
  { code: 'ELEC220', name: 'Fundamentals of Electrical Engineering', department: 'ELEC' },
  { code: 'MECH370', name: 'Introduction to Mechanical Engineering Design', department: 'MECH' },
  { code: 'PSYC101', name: 'Introduction to Psychology', department: 'PSYC' },
  { code: 'HIST103', name: 'Introduction to History', department: 'HIST' }
]

const styles = {
  container: 'relative',
  inputContainer: 'relative',
  input: 'w-full px-4 py-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent',
  inputError: 'border-red-500 focus:ring-red-500',
  dropdown: 'absolute z-10 w-full mt-1 bg-white border border-gray-200 rounded-lg shadow-lg max-h-60 overflow-auto',
  dropdownItem: 'px-4 py-3 hover:bg-gray-50 cursor-pointer border-b border-gray-100 last:border-b-0',
  dropdownItemActive: 'bg-blue-50 text-blue-700',
  courseCode: 'font-semibold text-gray-900',
  courseName: 'text-sm text-gray-600 mt-1',
  customOption: 'px-4 py-3 text-blue-600 font-medium border-t border-gray-200 hover:bg-blue-50 cursor-pointer',
  errorMessage: 'mt-1 text-sm text-red-600',
  noResults: 'px-4 py-3 text-gray-500 text-sm'
}

export default function CourseSelect({ 
  value, 
  onChange, 
  error, 
  placeholder = 'Select or enter course code (e.g., COMP101)',
  className = ''
}: CourseSelectProps) {
  const [isOpen, setIsOpen] = useState(false)
  const [searchTerm, setSearchTerm] = useState('')
  const [highlightedIndex, setHighlightedIndex] = useState(-1)
  const inputRef = useRef<HTMLInputElement>(null)
  const dropdownRef = useRef<HTMLDivElement>(null)

  // Filter courses based on search term
  const filteredCourses = RICE_COURSES.filter(course =>
    course.code.toLowerCase().includes(searchTerm.toLowerCase()) ||
    course.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
    course.department.toLowerCase().includes(searchTerm.toLowerCase())
  )

  // Validate course code format (e.g., COMP101, MATH220)
  const validateCourseCode = (code: string): boolean => {
    const courseRegex = /^[A-Z]{3,4}\d{3}$/
    return courseRegex.test(code.toUpperCase())
  }

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const inputValue = e.target.value
    setSearchTerm(inputValue)
    onChange(inputValue)
    setIsOpen(true)
    setHighlightedIndex(-1)
  }

  const handleCourseSelect = (courseCode: string) => {
    onChange(courseCode)
    setSearchTerm(courseCode)
    setIsOpen(false)
    setHighlightedIndex(-1)
  }

  const handleCustomCourse = () => {
    const customCode = searchTerm.toUpperCase()
    if (validateCourseCode(customCode)) {
      handleCourseSelect(customCode)
    } else {
      // Show format hint
      if (inputRef.current) {
        inputRef.current.focus()
      }
    }
  }

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (!isOpen) {
      if (e.key === 'Enter' || e.key === 'ArrowDown') {
        setIsOpen(true)
        setHighlightedIndex(0)
        e.preventDefault()
      }
      return
    }

    switch (e.key) {
      case 'ArrowDown':
        setHighlightedIndex(prev => 
          prev < filteredCourses.length - 1 ? prev + 1 : prev
        )
        e.preventDefault()
        break
      case 'ArrowUp':
        setHighlightedIndex(prev => prev > 0 ? prev - 1 : -1)
        e.preventDefault()
        break
      case 'Enter':
        if (highlightedIndex >= 0 && highlightedIndex < filteredCourses.length) {
          handleCourseSelect(filteredCourses[highlightedIndex].code)
        } else if (filteredCourses.length === 0 && searchTerm) {
          handleCustomCourse()
        }
        e.preventDefault()
        break
      case 'Escape':
        setIsOpen(false)
        setHighlightedIndex(-1)
        if (inputRef.current) {
          inputRef.current.blur()
        }
        e.preventDefault()
        break
    }
  }

  // Close dropdown when clicking outside
  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (
        dropdownRef.current && 
        !dropdownRef.current.contains(event.target as Node) &&
        inputRef.current &&
        !inputRef.current.contains(event.target as Node)
      ) {
        setIsOpen(false)
        setHighlightedIndex(-1)
      }
    }

    document.addEventListener('mousedown', handleClickOutside)
    return () => document.removeEventListener('mousedown', handleClickOutside)
  }, [])

  // Update search term when value changes externally
  useEffect(() => {
    setSearchTerm(value)
  }, [value])

  const inputClasses = `${styles.input} ${error ? styles.inputError : ''} ${className}`

  return (
    <div className={styles.container}>
      <div className={styles.inputContainer}>
        <input
          ref={inputRef}
          type="text"
          value={searchTerm}
          onChange={handleInputChange}
          onKeyDown={handleKeyDown}
          onFocus={() => setIsOpen(true)}
          placeholder={placeholder}
          className={inputClasses}
          aria-expanded={isOpen}
          aria-haspopup="listbox"
          role="combobox"
        />
        
        {isOpen && (
          <div ref={dropdownRef} className={styles.dropdown} role="listbox">
            {filteredCourses.length > 0 ? (
              filteredCourses.map((course, index) => (
                <div
                  key={course.code}
                  className={`${styles.dropdownItem} ${
                    index === highlightedIndex ? styles.dropdownItemActive : ''
                  }`}
                  onClick={() => handleCourseSelect(course.code)}
                  role="option"
                  aria-selected={index === highlightedIndex}
                >
                  <div className={styles.courseCode}>{course.code}</div>
                  <div className={styles.courseName}>{course.name}</div>
                </div>
              ))
            ) : (
              <div className={styles.noResults}>
                No courses found matching "{searchTerm}"
              </div>
            )}
            
            {/* Custom course option */}
            {searchTerm && 
             filteredCourses.length === 0 && 
             validateCourseCode(searchTerm) && (
              <div
                className={styles.customOption}
                onClick={handleCustomCourse}
                role="option"
              >
                Use "{searchTerm.toUpperCase()}" as custom course
              </div>
            )}
          </div>
        )}
      </div>
      
      {error && (
        <div className={styles.errorMessage}>{error}</div>
      )}
    </div>
  )
}
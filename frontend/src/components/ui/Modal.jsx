import { forwardRef, useEffect } from 'react'
import { createPortal } from 'react-dom'
import { cn } from '../../utils/cn'

const Modal = forwardRef(({
  isOpen = false,
  onClose,
  title,
  children,
  size = 'md',
  closeOnOverlayClick = true,
  closeOnEscape = true,
  showCloseButton = true,
  className,
  overlayClassName,
  contentClassName,
  ...props
}, ref) => {
  const sizeClasses = {
    xs: 'max-w-sm',
    sm: 'max-w-md',
    md: 'max-w-lg',
    lg: 'max-w-2xl',
    xl: 'max-w-4xl',
    full: 'max-w-full mx-4'
  }

  // Handle escape key
  useEffect(() => {
    if (!isOpen || !closeOnEscape) return

    const handleEscape = (event) => {
      if (event.key === 'Escape') {
        onClose?.()
      }
    }

    document.addEventListener('keydown', handleEscape)
    return () => document.removeEventListener('keydown', handleEscape)
  }, [isOpen, closeOnEscape, onClose])

  // Prevent body scroll when modal is open
  useEffect(() => {
    if (isOpen) {
      document.body.style.overflow = 'hidden'
    } else {
      document.body.style.overflow = 'unset'
    }

    return () => {
      document.body.style.overflow = 'unset'
    }
  }, [isOpen])

  if (!isOpen) return null

  const overlayClasses = cn(
    'fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50',
    overlayClassName
  )

  const contentClasses = cn(
    'bg-white rounded-lg shadow-xl w-full relative',
    sizeClasses[size],
    contentClassName
  )

  const handleOverlayClick = (e) => {
    if (closeOnOverlayClick && e.target === e.currentTarget) {
      onClose?.()
    }
  }

  const modalContent = (
    <div className={overlayClasses} onClick={handleOverlayClick}>
      <div 
        ref={ref}
        className={contentClasses}
        {...props}
      >
        {/* Header */}
        {(title || showCloseButton) && (
          <div className="flex items-center justify-between p-6 border-b border-gray-200">
            <h2 className="text-lg font-semibold text-gray-900">
              {title}
            </h2>
            {showCloseButton && (
              <button
                type="button"
                onClick={onClose}
                className="text-gray-400 hover:text-gray-600 focus:outline-none focus:text-gray-600 transition-colors"
              >
                <svg className="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
                  <path fillRule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clipRule="evenodd" />
                </svg>
              </button>
            )}
          </div>
        )}

        {/* Content */}
        <div className={cn(
          'p-6',
          !(title || showCloseButton) && 'pt-6',
          className
        )}>
          {children}
        </div>
      </div>
    </div>
  )

  // Render modal in portal
  return createPortal(modalContent, document.body)
})

Modal.displayName = 'Modal'

export default Modal
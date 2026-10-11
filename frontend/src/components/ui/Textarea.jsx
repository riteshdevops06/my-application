import { forwardRef } from 'react'
import { cn } from '../../utils/cn'

const Textarea = forwardRef(({
  label,
  error,
  helperText,
  rows = 3,
  resize = true,
  disabled = false,
  required = false,
  maxLength,
  showCounter = false,
  className,
  ...props
}, ref) => {
  const textareaClasses = cn(
    // Base styles
    'block w-full rounded-md border shadow-sm transition-colors',
    'placeholder:text-gray-400 focus:outline-none focus:ring-1',
    'disabled:bg-gray-50 disabled:text-gray-500 disabled:cursor-not-allowed',
    // Size styles
    'px-3 py-2 text-sm',
    // Variant styles
    error 
      ? 'border-red-500 focus:border-red-500 focus:ring-red-500'
      : 'border-gray-300 focus:border-primary-500 focus:ring-primary-500',
    // Resize behavior
    !resize && 'resize-none',
    className
  )

  const currentLength = props.value ? props.value.length : 0
  const isOverLimit = maxLength && currentLength > maxLength

  return (
    <div className="space-y-2">
      {label && (
        <label className="block text-sm font-medium text-gray-700">
          {label}
          {required && <span className="text-red-500 ml-1">*</span>}
        </label>
      )}
      
      <div className="relative">
        <textarea
          ref={ref}
          rows={rows}
          disabled={disabled}
          required={required}
          maxLength={maxLength}
          className={textareaClasses}
          {...props}
        />
        
        {showCounter && maxLength && (
          <div className="absolute bottom-2 right-2 text-xs text-gray-400 bg-white px-1">
            <span className={isOverLimit ? 'text-red-500' : ''}>
              {currentLength}
            </span>
            /{maxLength}
          </div>
        )}
      </div>
      
      {error && (
        <p className="text-sm text-red-600 flex items-center">
          <svg className="w-4 h-4 mr-1" fill="currentColor" viewBox="0 0 20 20">
            <path fillRule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7 4a1 1 0 11-2 0 1 1 0 012 0zm-1-9a1 1 0 00-1 1v4a1 1 0 102 0V6a1 1 0 00-1-1z" clipRule="evenodd" />
          </svg>
          {error}
        </p>
      )}
      
      {helperText && !error && (
        <p className="text-sm text-gray-500">{helperText}</p>
      )}
      
      {showCounter && maxLength && !error && (
        <p className="text-sm text-gray-500 text-right">
          {currentLength} / {maxLength} characters
        </p>
      )}
    </div>
  )
})

Textarea.displayName = 'Textarea'

export default Textarea
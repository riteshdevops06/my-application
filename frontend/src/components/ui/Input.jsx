import { forwardRef } from 'react'
import { cn } from '../../utils/cn'

const Input = forwardRef(({
  type = 'text',
  label,
  error,
  helperText,
  leftIcon,
  rightIcon,
  size = 'md',
  variant = 'default',
  disabled = false,
  required = false,
  className,
  ...props
}, ref) => {
  const sizeClasses = {
    sm: 'px-3 py-2 text-sm',
    md: 'px-3 py-2 text-sm',
    lg: 'px-4 py-3 text-base'
  }

  const variantClasses = {
    default: 'border-gray-300 focus:border-primary-500 focus:ring-primary-500',
    error: 'border-red-500 focus:border-red-500 focus:ring-red-500',
    success: 'border-green-500 focus:border-green-500 focus:ring-green-500'
  }

  const inputClasses = cn(
    // Base styles
    'block w-full rounded-md border shadow-sm transition-colors',
    'placeholder:text-gray-400 focus:outline-none focus:ring-1',
    'disabled:bg-gray-50 disabled:text-gray-500 disabled:cursor-not-allowed',
    // Size styles
    sizeClasses[size],
    // Variant styles
    error ? variantClasses.error : variantClasses[variant],
    // Icon padding adjustments
    leftIcon && 'pl-10',
    rightIcon && 'pr-10',
    className
  )

  const containerClasses = cn(
    'relative',
    disabled && 'opacity-60'
  )

  const iconClasses = cn(
    'absolute top-1/2 transform -translate-y-1/2 text-gray-400',
    size === 'sm' ? 'w-4 h-4' : size === 'lg' ? 'w-6 h-6' : 'w-5 h-5'
  )

  return (
    <div className="space-y-2">
      {label && (
        <label className="block text-sm font-medium text-gray-700">
          {label}
          {required && <span className="text-red-500 ml-1">*</span>}
        </label>
      )}
      
      <div className={containerClasses}>
        {leftIcon && (
          <div className={cn(iconClasses, 'left-3')}>
            {leftIcon}
          </div>
        )}
        
        <input
          ref={ref}
          type={type}
          disabled={disabled}
          required={required}
          className={inputClasses}
          {...props}
        />
        
        {rightIcon && (
          <div className={cn(iconClasses, 'right-3')}>
            {rightIcon}
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
    </div>
  )
})

Input.displayName = 'Input'

export default Input
import { forwardRef } from 'react'
import { cn } from '../../utils/cn'

const Button = forwardRef(({
  children,
  variant = 'primary',
  size = 'md',
  disabled = false,
  loading = false,
  fullWidth = false,
  leftIcon,
  rightIcon,
  className,
  ...props
}, ref) => {
  const baseClasses = 'inline-flex items-center justify-center rounded-md font-medium transition-all focus:outline-none focus:ring-2 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed'
  
  const sizeClasses = {
    xs: 'px-2.5 py-1.5 text-xs',
    sm: 'px-3 py-2 text-sm',
    md: 'px-4 py-2 text-sm',
    lg: 'px-4 py-3 text-base',
    xl: 'px-6 py-3 text-base'
  }
  
  const variantClasses = {
    primary: 'bg-primary-600 text-white hover:bg-primary-700 focus:ring-primary-500 shadow-sm',
    secondary: 'bg-white text-gray-700 border border-gray-300 hover:bg-gray-50 focus:ring-primary-500 shadow-sm',
    success: 'bg-green-600 text-white hover:bg-green-700 focus:ring-green-500 shadow-sm',
    danger: 'bg-red-600 text-white hover:bg-red-700 focus:ring-red-500 shadow-sm',
    warning: 'bg-yellow-600 text-white hover:bg-yellow-700 focus:ring-yellow-500 shadow-sm',
    ghost: 'text-gray-700 hover:bg-gray-100 focus:ring-gray-500',
    link: 'text-primary-600 hover:text-primary-700 focus:ring-primary-500 underline-offset-4 hover:underline'
  }
  
  const buttonClasses = cn(
    baseClasses,
    sizeClasses[size],
    variantClasses[variant],
    fullWidth && 'w-full',
    loading && 'cursor-wait',
    className
  )
  
  const iconSize = size === 'xs' ? 'w-3 h-3' : size === 'sm' ? 'w-4 h-4' : 'w-5 h-5'
  
  return (
    <button
      ref={ref}
      disabled={disabled || loading}
      className={buttonClasses}
      {...props}
    >
      {loading ? (
        <svg className={cn(iconSize, 'animate-spin mr-2')} fill="none" viewBox="0 0 24 24">
          <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" />
          <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
        </svg>
      ) : leftIcon ? (
        <span className={cn(iconSize, 'mr-2')}>
          {leftIcon}
        </span>
      ) : null}
      
      {children}
      
      {rightIcon && !loading && (
        <span className={cn(iconSize, 'ml-2')}>
          {rightIcon}
        </span>
      )}
    </button>
  )
})

Button.displayName = 'Button'

export default Button
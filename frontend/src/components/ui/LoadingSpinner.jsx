import { forwardRef } from 'react'
import { cn } from '../../utils/cn'

const LoadingSpinner = forwardRef(({
  size = 'md',
  color = 'primary',
  text,
  className,
  ...props
}, ref) => {
  const sizeClasses = {
    xs: 'w-3 h-3',
    sm: 'w-4 h-4',
    md: 'w-6 h-6',
    lg: 'w-8 h-8',
    xl: 'w-12 h-12'
  }
  
  const colorClasses = {
    primary: 'text-primary-600',
    white: 'text-white',
    gray: 'text-gray-600',
    success: 'text-green-600',
    warning: 'text-yellow-600',
    danger: 'text-red-600'
  }

  const spinnerClasses = cn(
    'animate-spin',
    sizeClasses[size],
    colorClasses[color],
    className
  )

  const Spinner = () => (
    <svg 
      className={spinnerClasses} 
      fill="none" 
      viewBox="0 0 24 24" 
      {...props}
      ref={ref}
    >
      <circle 
        className="opacity-25" 
        cx="12" 
        cy="12" 
        r="10" 
        stroke="currentColor" 
        strokeWidth="4"
      />
      <path 
        className="opacity-75" 
        fill="currentColor" 
        d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
      />
    </svg>
  )

  if (text) {
    return (
      <div className="flex items-center justify-center space-x-2">
        <Spinner />
        <span className="text-sm text-gray-600">{text}</span>
      </div>
    )
  }

  return <Spinner />
})

LoadingSpinner.displayName = 'LoadingSpinner'

export default LoadingSpinner
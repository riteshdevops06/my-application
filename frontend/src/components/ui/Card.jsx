import { forwardRef } from 'react'
import { cn } from '../../utils/cn'

const Card = forwardRef(({
  children,
  variant = 'default',
  padding = 'default',
  shadow = 'default',
  hover = false,
  className,
  ...props
}, ref) => {
  const baseClasses = 'bg-white rounded-lg border'
  
  const variantClasses = {
    default: 'border-gray-200',
    elevated: 'border-gray-200 shadow-sm',
    outlined: 'border-2 border-gray-300',
    ghost: 'border-transparent'
  }
  
  const paddingClasses = {
    none: '',
    sm: 'p-3',
    default: 'p-4',
    md: 'p-6',
    lg: 'p-8'
  }
  
  const shadowClasses = {
    none: '',
    sm: 'shadow-sm',
    default: 'shadow',
    md: 'shadow-md',
    lg: 'shadow-lg',
    xl: 'shadow-xl'
  }

  const cardClasses = cn(
    baseClasses,
    variantClasses[variant],
    paddingClasses[padding],
    shadowClasses[shadow],
    hover && 'hover:shadow-md transition-shadow duration-200',
    className
  )

  return (
    <div ref={ref} className={cardClasses} {...props}>
      {children}
    </div>
  )
})

Card.displayName = 'Card'

// Compound components for better organization
const CardHeader = forwardRef(({ children, className, ...props }, ref) => (
  <div ref={ref} className={cn('pb-4 border-b border-gray-200', className)} {...props}>
    {children}
  </div>
))

CardHeader.displayName = 'CardHeader'

const CardTitle = forwardRef(({ children, className, ...props }, ref) => (
  <h3 ref={ref} className={cn('text-lg font-semibold text-gray-900', className)} {...props}>
    {children}
  </h3>
))

CardTitle.displayName = 'CardTitle'

const CardDescription = forwardRef(({ children, className, ...props }, ref) => (
  <p ref={ref} className={cn('text-sm text-gray-600 mt-1', className)} {...props}>
    {children}
  </p>
))

CardDescription.displayName = 'CardDescription'

const CardContent = forwardRef(({ children, className, ...props }, ref) => (
  <div ref={ref} className={cn('pt-4', className)} {...props}>
    {children}
  </div>
))

CardContent.displayName = 'CardContent'

const CardFooter = forwardRef(({ children, className, ...props }, ref) => (
  <div ref={ref} className={cn('pt-4 border-t border-gray-200', className)} {...props}>
    {children}
  </div>
))

CardFooter.displayName = 'CardFooter'

// Attach compound components to main Card component
Card.Header = CardHeader
Card.Title = CardTitle
Card.Description = CardDescription
Card.Content = CardContent
Card.Footer = CardFooter

export default Card
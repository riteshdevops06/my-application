import { forwardRef } from 'react'
import { cn } from '../../utils/cn'

const Container = forwardRef(({
  children,
  size = 'default',
  padding = 'default',
  center = true,
  className,
  ...props
}, ref) => {
  const sizeClasses = {
    sm: 'max-w-3xl',
    default: 'max-w-7xl',
    lg: 'max-w-screen-xl',
    xl: 'max-w-screen-2xl',
    full: 'max-w-full'
  }
  
  const paddingClasses = {
    none: '',
    sm: 'px-4 sm:px-6',
    default: 'px-4 sm:px-6 lg:px-8',
    lg: 'px-6 sm:px-8 lg:px-12'
  }

  const containerClasses = cn(
    'w-full',
    sizeClasses[size],
    paddingClasses[padding],
    center && 'mx-auto',
    className
  )

  return (
    <div ref={ref} className={containerClasses} {...props}>
      {children}
    </div>
  )
})

Container.displayName = 'Container'

export default Container
import { ReactNode } from 'react';
import classNames from 'classnames';

interface TypographyProps {
  variant?: 'h1' | 'h2' | 'h3' | 'h4' | 'body' | 'caption';
  align?: 'left' | 'center' | 'right';
  color?: 'primary' | 'secondary' | 'success' | 'error' | 'default';
  className?: string;
  children: ReactNode;
}

export function Typography({
  variant = 'body',
  align = 'left',
  color = 'default',
  className,
  children,
}: TypographyProps) {
  const variantClasses = {
    h1: 'text-4xl font-bold',
    h2: 'text-3xl font-extrabold',
    h3: 'text-2xl font-semibold',
    h4: 'text-xl font-semibold',
    body: 'text-base',
    caption: 'text-sm',
  };

  const alignClasses = {
    left: 'text-left',
    center: 'text-center',
    right: 'text-right',
  };

  const colorClasses = {
    primary: 'text-primary-600',
    secondary: 'text-gray-600',
    success: 'text-green-600',
    error: 'text-red-600',
    default: 'text-gray-900',
  };

  const Component = variant.startsWith('h') ? variant : 'p';

  return (
    <Component
      className={classNames(
        variantClasses[variant],
        alignClasses[align],
        colorClasses[color],
        className
      )}
    >
      {children}
    </Component>
  );
}

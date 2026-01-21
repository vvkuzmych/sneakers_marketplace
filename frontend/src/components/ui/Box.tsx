import { ReactNode } from 'react';
import classNames from 'classnames';

interface BoxProps {
  children: ReactNode;
  className?: string;
  mt?: number; // margin-top (1-12)
  mb?: number; // margin-bottom
  p?: number;  // padding
  flex?: boolean;
  flexDirection?: 'row' | 'column';
  alignItems?: 'start' | 'center' | 'end';
  justifyContent?: 'start' | 'center' | 'end' | 'between';
  gap?: number;
}

export function Box({
  children,
  className,
  mt,
  mb,
  p,
  flex,
  flexDirection = 'row',
  alignItems,
  justifyContent,
  gap,
}: BoxProps) {
  const marginTopClasses = {
    1: 'mt-1',
    2: 'mt-2',
    3: 'mt-3',
    4: 'mt-4',
    6: 'mt-6',
    8: 'mt-8',
    12: 'mt-12',
  };

  const marginBottomClasses = {
    1: 'mb-1',
    2: 'mb-2',
    3: 'mb-3',
    4: 'mb-4',
    6: 'mb-6',
    8: 'mb-8',
    12: 'mb-12',
  };

  const paddingClasses = {
    1: 'p-1',
    2: 'p-2',
    3: 'p-3',
    4: 'p-4',
    6: 'p-6',
    8: 'p-8',
    12: 'p-12',
  };

  const gapClasses = {
    1: 'gap-1',
    2: 'gap-2',
    3: 'gap-3',
    4: 'gap-4',
    6: 'gap-6',
    8: 'gap-8',
  };

  return (
    <div
      className={classNames(
        mt && marginTopClasses[mt as keyof typeof marginTopClasses],
        mb && marginBottomClasses[mb as keyof typeof marginBottomClasses],
        p && paddingClasses[p as keyof typeof paddingClasses],
        flex && 'flex',
        flexDirection === 'column' && 'flex-col',
        alignItems === 'center' && 'items-center',
        alignItems === 'start' && 'items-start',
        alignItems === 'end' && 'items-end',
        justifyContent === 'center' && 'justify-center',
        justifyContent === 'between' && 'justify-between',
        gap && gapClasses[gap as keyof typeof gapClasses],
        className
      )}
    >
      {children}
    </div>
  );
}

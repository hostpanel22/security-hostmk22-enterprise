import type { HTMLAttributes, PropsWithChildren } from 'react';

export function Card({ children, className = '', ...props }: PropsWithChildren<HTMLAttributes<HTMLDivElement>>) {
  return (
    <div
      {...props}
      className={['rounded-lg border border-slate-800 bg-slate-900/80 p-4 shadow-sm', className].join(' ')}
    >
      {children}
    </div>
  );
}

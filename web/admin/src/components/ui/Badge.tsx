import type { HTMLAttributes, PropsWithChildren } from 'react';

export function Badge({ children, className = '', ...props }: PropsWithChildren<HTMLAttributes<HTMLSpanElement>>) {
  return (
    <span
      {...props}
      className={['inline-flex items-center rounded-full border border-slate-700 bg-slate-800 px-2.5 py-1 text-xs font-medium text-slate-200', className].join(' ')}
    >
      {children}
    </span>
  );
}

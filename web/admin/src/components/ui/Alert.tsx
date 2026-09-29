import type { HTMLAttributes, PropsWithChildren } from 'react';

export function Alert({ children, className = '', ...props }: PropsWithChildren<HTMLAttributes<HTMLDivElement>>) {
  return (
    <div
      {...props}
      className={['rounded-md border border-sky-800 bg-sky-950/60 p-3 text-sm text-sky-100', className].join(' ')}
    >
      {children}
    </div>
  );
}

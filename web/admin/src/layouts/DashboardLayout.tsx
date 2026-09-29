import { useState, type FormEvent } from 'react';
import { Outlet, NavLink } from 'react-router-dom';
import { useQueryClient } from '@tanstack/react-query';
import { useAuth } from '../features/auth/AuthContext';
import { apiPost, type ApiResponse } from '../services/api';
import { Button } from '../components/ui/Button';
import { Input } from '../components/ui/Input';

const navItems = [
  { label: 'Overview', to: '/' },
  { label: 'Tenants', to: '/tenants' },
  { label: 'Users', to: '/users' },
  { label: 'Sessions', to: '/sessions' },
  { label: 'RBAC', to: '/rbac' },
  { label: 'Permissions', to: '/rbac' },
  { label: 'API Keys', to: '/api-keys' },
  { label: 'Secrets', to: '/secrets' },
  { label: 'Fraud Detection', to: '/fraud' },
  { label: 'Audit Logs', to: '/audit' },
  { label: 'Settings', to: '/settings' },
  { label: 'Monitoring', to: '/monitoring' },
];

export function DashboardLayout() {
  const { accessToken, setSession, clearSession } = useAuth();
  const queryClient = useQueryClient();
  const [tenantSlug, setTenantSlug] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [mfaCode, setMfaCode] = useState('');
  const [loginError, setLoginError] = useState('');
  const [isSigningIn, setIsSigningIn] = useState(false);

  const signIn = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    setLoginError('');
    setIsSigningIn(true);
    try {
      const response = await apiPost<ApiResponse<{ access_token: string; refresh_token: string }>>('/api/v1/auth/login', {
        tenant_slug: tenantSlug,
        email,
        password,
        mfa_code: mfaCode,
      });
      const tokens = response.data;
      if (!tokens?.access_token || !tokens.refresh_token) {
        throw new Error(response.error ?? 'The authentication response did not include both tokens.');
      }
      setSession({ accessToken: tokens.access_token, refreshToken: tokens.refresh_token });
      await queryClient.invalidateQueries();
      setPassword('');
      setMfaCode('');
    } catch (error) {
      setLoginError(error instanceof Error ? error.message : 'Unable to sign in.');
    } finally {
      setIsSigningIn(false);
    }
  };

  const signOut = () => {
    clearSession();
    queryClient.clear();
  };

  if (!accessToken) {
    return (
      <main className="flex min-h-screen items-center justify-center bg-slate-950 px-4 py-12 text-slate-100">
        <section className="w-full max-w-md border border-slate-800 bg-slate-900 p-6">
          <p className="text-xs font-semibold uppercase tracking-[0.2em] text-sky-400">HostMK Security Framework</p>
          <h1 className="mt-3 text-2xl font-semibold">Sign in to the console</h1>
          <form className="mt-6 space-y-4" onSubmit={signIn}>
            <label className="block space-y-1 text-sm text-slate-300">
              <span>Tenant slug</span>
              <Input autoComplete="organization" value={tenantSlug} onChange={(event) => setTenantSlug(event.target.value)} required />
            </label>
            <label className="block space-y-1 text-sm text-slate-300">
              <span>Email</span>
              <Input type="email" autoComplete="username" value={email} onChange={(event) => setEmail(event.target.value)} required />
            </label>
            <label className="block space-y-1 text-sm text-slate-300">
              <span>Password</span>
              <Input type="password" autoComplete="current-password" value={password} onChange={(event) => setPassword(event.target.value)} required />
            </label>
            <label className="block space-y-1 text-sm text-slate-300">
              <span>MFA code, if enabled</span>
              <Input inputMode="numeric" autoComplete="one-time-code" value={mfaCode} onChange={(event) => setMfaCode(event.target.value)} />
            </label>
            {loginError ? <p role="alert" className="border border-rose-900 bg-rose-950/40 p-3 text-sm text-rose-200">{loginError}</p> : null}
            <Button type="submit" disabled={isSigningIn} className="w-full">{isSigningIn ? 'Signing in…' : 'Sign in'}</Button>
          </form>
        </section>
      </main>
    );
  }

  return (
    <div className="min-h-screen bg-slate-950 text-slate-100">
      <div className="flex min-h-screen flex-col lg:flex-row">
        <aside className="w-full border-b border-slate-800 bg-slate-900 lg:w-72 lg:border-b-0 lg:border-r">
          <div className="flex items-center justify-between border-b border-slate-800 px-4 py-4 lg:block">
            <div>
              <div className="text-xs font-semibold uppercase tracking-[0.2em] text-sky-400">HostMK</div>
              <div className="mt-1 text-lg font-semibold text-slate-100">Security Framework</div>
            </div>
          </div>

          <nav className="flex flex-col gap-1 p-3">
            {navItems.map((item) => (
              <NavLink
                key={item.label}
                to={item.to}
                end={item.to === '/'}
                className={({ isActive }) =>
                  [
                    'rounded-md px-3 py-2 text-sm transition',
                    isActive ? 'bg-sky-500/10 text-sky-300 ring-1 ring-sky-500/30' : 'text-slate-300 hover:bg-slate-800 hover:text-slate-100',
                  ].join(' ')
                }
              >
                {item.label}
              </NavLink>
            ))}
          </nav>
        </aside>

        <div className="flex min-w-0 flex-1 flex-col">
          <header className="border-b border-slate-800 bg-slate-900/95 px-4 py-3 backdrop-blur">
            <div className="flex items-center justify-between gap-4">
              <div>
                <div className="text-xs uppercase tracking-[0.2em] text-slate-500">Dashboard</div>
                <div className="text-lg font-semibold text-slate-100">Enterprise Console</div>
              </div>

              <div className="flex items-center gap-3">
                <span className="text-xs text-slate-400">Authenticated</span>
                <Button type="button" onClick={signOut} className="px-3 py-1.5 text-xs">Sign out</Button>
              </div>
            </div>
          </header>

          <div className="border-b border-slate-800 bg-slate-950/60 px-4 py-3 text-sm text-slate-400">
            <div className="flex flex-wrap items-center gap-2">
              <span>Home</span>
              <span>/</span>
              <span className="text-slate-200">Overview</span>
            </div>
          </div>

          <main className="flex-1 p-4 md:p-6">
            <div className="mb-4 flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
              <div>
                <h1 className="text-2xl font-semibold text-slate-100">Overview</h1>
                <p className="text-sm text-slate-400">Security operations and tenant health</p>
              </div>
              <div className="flex items-center gap-2">
                <span className="rounded-full border border-emerald-500/30 bg-emerald-500/10 px-2.5 py-1 text-xs text-emerald-300">System online</span>
              </div>
            </div>

            <div className="mb-6 rounded-xl border border-slate-800 bg-slate-900/80 p-3 text-sm text-slate-300">
              Notification center is available for operational alerts.
            </div>

            <Outlet />
          </main>
        </div>
      </div>
    </div>
  );
}

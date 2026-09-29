import { createContext, useContext, useMemo, useState } from 'react';
import type { ReactNode } from 'react';
import { authStorage, createEmptyAuthState } from '@/lib/auth';

interface AuthContextValue {
  accessToken: string | null;
  refreshToken: string | null;
  currentUser: Record<string, unknown> | null;
  permissions: string[];
  setSession: (session: {
    accessToken?: string | null;
    refreshToken?: string | null;
    user?: Record<string, unknown> | null;
    permissions?: string[];
  }) => void;
  clearSession: () => void;
}

interface AuthSessionState {
  accessToken: string | null;
  refreshToken: string | null;
  user: Record<string, unknown> | null;
  permissions: string[];
}

const AuthContext = createContext<AuthContextValue | undefined>(undefined);

export function AuthProvider({ children }: { children: ReactNode }) {
  const [state, setState] = useState<AuthSessionState>(() => ({
    accessToken: authStorage.getAccessToken(),
    refreshToken: authStorage.getRefreshToken(),
    user: null,
    permissions: [],
  }));

  const setSession = (session: {
    accessToken?: string | null;
    refreshToken?: string | null;
    user?: Record<string, unknown> | null;
    permissions?: string[];
  }) => {
    setState((current) => {
      const next: AuthSessionState = {
        ...current,
        ...session,
      };

      if (session.accessToken !== undefined) {
        authStorage.setAccessToken(session.accessToken);
      }

      if (session.refreshToken !== undefined) {
        authStorage.setRefreshToken(session.refreshToken);
      }

      return next;
    });
  };

  const clearSession = () => {
    authStorage.clear();
    setState(createEmptyAuthState());
  };

  const value = useMemo<AuthContextValue>(
    () => ({
      accessToken: state.accessToken,
      refreshToken: state.refreshToken,
      currentUser: state.user,
      permissions: state.permissions,
      setSession,
      clearSession,
    }),
    [state.accessToken, state.refreshToken, state.user, state.permissions],
  );

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth() {
  const context = useContext(AuthContext);

  if (!context) {
    throw new Error('useAuth must be used within AuthProvider');
  }

  return context;
}

export interface AuthState {
  accessToken: string | null;
  refreshToken: string | null;
  user: Record<string, unknown> | null;
  permissions: string[];
}

const STORAGE_KEYS = {
  access: 'hostmk.accessToken',
  refresh: 'hostmk.refreshToken',
};

export const authStorage = {
  getAccessToken: () => localStorage.getItem(STORAGE_KEYS.access),
  setAccessToken: (token: string | null) => {
    if (token) {
      localStorage.setItem(STORAGE_KEYS.access, token);
      return;
    }
    localStorage.removeItem(STORAGE_KEYS.access);
  },
  getRefreshToken: () => localStorage.getItem(STORAGE_KEYS.refresh),
  setRefreshToken: (token: string | null) => {
    if (token) {
      localStorage.setItem(STORAGE_KEYS.refresh, token);
      return;
    }
    localStorage.removeItem(STORAGE_KEYS.refresh);
  },
  clear: () => {
    localStorage.removeItem(STORAGE_KEYS.access);
    localStorage.removeItem(STORAGE_KEYS.refresh);
  },
};

export const createEmptyAuthState = (): AuthState => ({
  accessToken: null,
  refreshToken: null,
  user: null,
  permissions: [],
});

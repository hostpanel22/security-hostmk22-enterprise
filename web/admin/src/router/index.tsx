import { createBrowserRouter } from 'react-router-dom';
import { DashboardLayout } from '../layouts/DashboardLayout';
import { OverviewPage } from '../pages/overview/OverviewPage';
import { TenantsPage } from '../pages/tenants/TenantsPage';
import { UsersPage } from '../pages/users/UsersPage';
import { SessionsPage } from '../pages/sessions/SessionsPage';
import { RbacPage } from '../pages/rbac/RbacPage';
import { ApiKeysPage } from '../pages/api-keys/ApiKeysPage';
import { SecretsPage } from '../pages/secrets/SecretsPage';
import { FraudPage } from '../pages/fraud/FraudPage';
import { AuditPage } from '../pages/audit/AuditPage';
import { SettingsPage } from '../pages/settings/SettingsPage';
import { MonitoringPage } from '../pages/monitoring/MonitoringPage';
import { LandingPage } from '../pages/landing/LandingPage';

export const appRouter = createBrowserRouter([
  { path: '/landing', element: <LandingPage /> },
  {
    path: '/',
    element: <DashboardLayout />,
    children: [
      { index: true, element: <OverviewPage /> },
      { path: 'tenants', element: <TenantsPage /> },
      { path: 'users', element: <UsersPage /> },
      { path: 'sessions', element: <SessionsPage /> },
      { path: 'rbac', element: <RbacPage /> },
      { path: 'api-keys', element: <ApiKeysPage /> },
      { path: 'secrets', element: <SecretsPage /> },
      { path: 'fraud', element: <FraudPage /> },
      { path: 'audit', element: <AuditPage /> },
      { path: 'settings', element: <SettingsPage /> },
      { path: 'monitoring', element: <MonitoringPage /> },
    ],
  },
]);

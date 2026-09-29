import { useMemo, useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { fetchTenants, type TenantRecord } from '../../services/tenants';

export function TenantsPage() {
  const [query, setQuery] = useState('');
  const [page, setPage] = useState(1);

  const { data: tenants = [], isLoading, isError, error } = useQuery({
    queryKey: ['tenants'],
    queryFn: fetchTenants,
  });

  const visibleTenants = useMemo(() => {
    const normalized = query.trim().toLowerCase();
    const filtered = normalized
      ? tenants.filter((tenant: TenantRecord) => [tenant.name, tenant.slug, tenant.status, tenant.plan].filter(Boolean).join(' ').toLowerCase().includes(normalized))
      : tenants;
    return filtered.slice((page - 1) * 10, page * 10);
  }, [query, page, tenants]);

  if (isLoading) {
    return <div className="rounded-xl border border-slate-800 bg-slate-900/80 p-6 text-slate-300">Loading tenants…</div>;
  }

  if (isError) {
    return <div className="rounded-xl border border-rose-900/60 bg-rose-950/30 p-6 text-rose-200">Unable to load tenants: {error instanceof Error ? error.message : 'Unknown error'}</div>;
  }

  return (
    <div className="space-y-4">
      <div className="rounded-xl border border-slate-800 bg-slate-900/80 p-4">
        <div className="flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
          <input
            value={query}
            onChange={(event) => {
              setQuery(event.target.value);
              setPage(1);
            }}
            placeholder="Search tenants"
            className="w-full rounded-md border border-slate-700 bg-slate-950 px-3 py-2 text-sm text-slate-100 outline-none ring-0 placeholder:text-slate-500 md:max-w-sm"
          />
          <div className="text-sm text-slate-400">{tenants.length} records</div>
        </div>
      </div>

      <div className="overflow-x-auto rounded-xl border border-slate-800 bg-slate-900/80">
        <table className="min-w-full text-left text-sm text-slate-200">
          <thead className="bg-slate-800/80 text-slate-300">
            <tr>
              <th className="px-4 py-3">Name</th>
              <th className="px-4 py-3">Status</th>
              <th className="px-4 py-3">Plan</th>
              <th className="px-4 py-3">Branding</th>
              <th className="px-4 py-3">Quotas</th>
            </tr>
          </thead>
          <tbody>
            {visibleTenants.length ? (
              visibleTenants.map((tenant: TenantRecord) => (
                <tr key={tenant.id ?? tenant.slug ?? 'tenant'} className="border-t border-slate-800">
                  <td className="px-4 py-3">
                    <div className="font-medium text-slate-100">{tenant.name ?? 'Unnamed tenant'}</div>
                    <div className="text-xs text-slate-500">{tenant.slug ?? '—'}</div>
                  </td>
                  <td className="px-4 py-3">
                    <span className="rounded-full border border-slate-700 bg-slate-800 px-2.5 py-1 text-xs uppercase text-slate-200">
                      {tenant.status ?? 'unknown'}
                    </span>
                  </td>
                  <td className="px-4 py-3 text-slate-300">{tenant.plan ?? '—'}</td>
                  <td className="px-4 py-3 text-slate-300">{tenant.branding ? JSON.stringify(tenant.branding) : '—'}</td>
                  <td className="px-4 py-3 text-slate-300">{tenant.quotas ? JSON.stringify(tenant.quotas) : '—'}</td>
                </tr>
              ))
            ) : (
              <tr>
                <td colSpan={5} className="px-4 py-10 text-center text-slate-400">No tenants found.</td>
              </tr>
            )}
          </tbody>
        </table>
      </div>

      <div className="flex items-center justify-between rounded-xl border border-slate-800 bg-slate-900/80 p-3 text-sm text-slate-300">
        <button
          type="button"
          onClick={() => setPage((prev) => Math.max(1, prev - 1))}
          className="rounded-md border border-slate-700 bg-slate-800 px-3 py-1.5 disabled:cursor-not-allowed disabled:opacity-40"
          disabled={page === 1}
        >
          Previous
        </button>
        <span>Page {page}</span>
        <button
          type="button"
          onClick={() => setPage((prev) => prev + 1)}
          className="rounded-md border border-slate-700 bg-slate-800 px-3 py-1.5 disabled:cursor-not-allowed disabled:opacity-40"
          disabled={visibleTenants.length < 10}
        >
          Next
        </button>
      </div>
    </div>
  );
}

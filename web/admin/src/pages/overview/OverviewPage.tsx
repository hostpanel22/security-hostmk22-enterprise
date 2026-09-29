import { useQuery } from '@tanstack/react-query';
import { Button } from '../../components/ui/Button';
import type { StandardApiError } from '../../services/api';
import { fetchOverviewCounts } from '../../services/overview';

export function OverviewPage() {
  const { data, isLoading, isError, error, refetch, isFetching } = useQuery({
    queryKey: ['overview'],
    queryFn: fetchOverviewCounts,
  });
  const requestError = error as StandardApiError | null;

  const sections = [
    { title: 'Security Overview', items: data?.security ?? [] },
    { title: 'Tenant Overview', items: data?.tenants ?? [] },
    { title: 'Fraud Overview', items: data?.fraud ?? [] },
    { title: 'Audit Overview', items: data?.audit ?? [] },
  ];

  if (isLoading) {
    return <div className="rounded-xl border border-slate-800 bg-slate-900/60 p-6 text-slate-300">Loading dashboard metrics…</div>;
  }

  if (isError) {
    return (
      <div className="rounded-xl border border-rose-900/60 bg-rose-950/30 p-6 text-rose-200">
        <p>Unable to load overview metrics: {error instanceof Error ? error.message : 'Unknown error'}</p>
        {requestError?.url ? (
          <div className="mt-3 space-y-1 font-mono text-xs text-rose-100">
            <p>{requestError.method} {requestError.url}</p>
            <p>HTTP {requestError.status} {requestError.statusText}</p>
          </div>
        ) : null}
        {requestError?.details !== undefined ? (
          <details className="mt-3 text-sm">
            <summary className="cursor-pointer">Response body</summary>
            <pre className="mt-2 overflow-x-auto whitespace-pre-wrap break-words text-xs">{typeof requestError.details === 'string' ? requestError.details : JSON.stringify(requestError.details, null, 2)}</pre>
          </details>
        ) : null}
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between gap-4">
        <p className="text-sm text-slate-400">Live totals calculated from the current tenant API lists.</p>
        <Button type="button" onClick={() => void refetch()} disabled={isFetching} aria-label="Refresh overview">
          {isFetching ? 'Refreshing…' : 'Refresh'}
        </Button>
      </div>
      <div className="grid gap-4 md:grid-cols-2 xl:grid-cols-4">
        {sections.flatMap((section) =>
          section.items.map((item) => (
            <div key={`${section.title}-${item.label}`} className="rounded-xl border border-slate-800 bg-slate-900/80 p-4">
              <div className="text-sm text-slate-400">{item.label}</div>
              <div className="mt-3 text-3xl font-semibold text-slate-100">{item.value}</div>
              {item.hint ? <div className="mt-2 text-xs text-slate-500">{item.hint}</div> : null}
            </div>
          )),
        )}
        {sections.every((section) => section.items.length === 0) ? (
          <div className="border border-slate-800 bg-slate-900/70 p-6 text-sm text-slate-300">No overview metrics are available from the current API data.</div>
        ) : null}
      </div>
    </div>
  );
}

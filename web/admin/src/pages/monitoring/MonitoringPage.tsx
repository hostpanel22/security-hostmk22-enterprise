import { useQuery } from '@tanstack/react-query';
import { fetchMonitoringStatus, type MonitoringItem } from '../../services/monitoring';

export function MonitoringPage() {
  const { data: health = [], isLoading, isError, error } = useQuery({
    queryKey: ['monitoring'],
    queryFn: fetchMonitoringStatus,
  });

  if (isLoading) {
    return <div className="rounded-xl border border-slate-800 bg-slate-900/80 p-6 text-slate-300">Loading monitoring status…</div>;
  }

  if (isError) {
    return <div className="rounded-xl border border-rose-900/60 bg-rose-950/30 p-6 text-rose-200">Unable to load monitoring data: {error instanceof Error ? error.message : 'Unknown error'}</div>;
  }

  return (
    <div className="space-y-4">
      <div className="grid gap-4 md:grid-cols-3">
        {health.length ? health.map((item: MonitoringItem) => (
          <div key={item.name ?? 'service'} className="rounded-xl border border-slate-800 bg-slate-900/80 p-4">
            <div className="text-sm text-slate-400">{item.name ?? 'Service'}</div>
            <div className="mt-3 text-2xl font-semibold text-slate-100">{item.status ?? 'unknown'}</div>
            <div className="mt-2 text-xs text-slate-500">{item.detail ?? 'No details'}</div>
          </div>
        )) : <div className="rounded-xl border border-slate-800 bg-slate-900/80 p-6 text-slate-400 md:col-span-3">No monitoring data available.</div>}
      </div>
    </div>
  );
}

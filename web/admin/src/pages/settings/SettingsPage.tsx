import { useQuery } from '@tanstack/react-query';
import { fetchSettings, type SettingRecord } from '../../services/settings';

export function SettingsPage() {
  const { data: settings = [], isLoading, isError, error } = useQuery({
    queryKey: ['settings'],
    queryFn: fetchSettings,
  });

  if (isLoading) {
    return <div className="rounded-xl border border-slate-800 bg-slate-900/80 p-6 text-slate-300">Loading settings…</div>;
  }

  if (isError) {
    return <div className="rounded-xl border border-rose-900/60 bg-rose-950/30 p-6 text-rose-200">Unable to load settings: {error instanceof Error ? error.message : 'Unknown error'}</div>;
  }

  return (
    <div className="space-y-4">
      <div className="overflow-x-auto rounded-xl border border-slate-800 bg-slate-900/80">
        <table className="min-w-full text-left text-sm text-slate-200">
          <thead className="bg-slate-800/80 text-slate-300">
            <tr>
              <th className="px-4 py-3">Key</th>
              <th className="px-4 py-3">Value</th>
            </tr>
          </thead>
          <tbody>
            {settings.length ? settings.map((setting: SettingRecord) => (
              <tr key={setting.id ?? setting.key ?? 'setting'} className="border-t border-slate-800">
                <td className="px-4 py-3 text-slate-100">{setting.key ?? '—'}</td>
                <td className="px-4 py-3 text-slate-300">{typeof setting.value === 'object' ? JSON.stringify(setting.value) : String(setting.value ?? '—')}</td>
              </tr>
            )) : <tr><td colSpan={2} className="px-4 py-10 text-center text-slate-400">No settings available.</td></tr>}
          </tbody>
        </table>
      </div>
    </div>
  );
}

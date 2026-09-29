import { useQuery } from '@tanstack/react-query';
import { fetchFraudRules, fetchFraudViolations, type FraudRuleRecord, type FraudViolationRecord } from '../../services/fraud';

export function FraudPage() {
  const rulesQuery = useQuery({ queryKey: ['fraud', 'rules'], queryFn: fetchFraudRules });
  const violationsQuery = useQuery({ queryKey: ['fraud', 'violations'], queryFn: fetchFraudViolations });

  if (rulesQuery.isLoading || violationsQuery.isLoading) return <div className="rounded-xl border border-slate-800 bg-slate-900/80 p-6 text-slate-300">Loading fraud data...</div>;
  if (rulesQuery.isError || violationsQuery.isError) {
    const error = rulesQuery.error ?? violationsQuery.error;
    return <div className="rounded-xl border border-rose-900/60 bg-rose-950/30 p-6 text-rose-200">Unable to load fraud data: {error instanceof Error ? error.message : 'Unknown error'}</div>;
  }

  const rules = rulesQuery.data ?? [];
  const events = violationsQuery.data ?? [];
  return <div className="grid gap-4 xl:grid-cols-2"><section className="rounded-xl border border-slate-800 bg-slate-900/80 p-4"><div className="mb-3 text-lg font-semibold text-slate-100">Rules</div><div className="space-y-2">{rules.length ? rules.map((rule: FraudRuleRecord) => <div key={rule.id ?? rule.name ?? 'fraud-rule'} className="rounded-lg border border-slate-800 bg-slate-950/60 p-3 text-sm"><div className="font-medium text-slate-100">{rule.name ?? 'Unknown rule'}</div><div className="mt-1 text-slate-400">{[rule.field, rule.operator, rule.threshold].filter(Boolean).join(' ') || 'No condition returned by the API'}</div><div className="mt-2 text-xs text-slate-500">Enabled: {rule.enabled ? 'Yes' : 'No'} · Action: {rule.action ?? '—'} · Score: {rule.score ?? '—'}</div></div>) : <div className="text-slate-400">No fraud rules available.</div>}</div></section><section className="rounded-xl border border-slate-800 bg-slate-900/80 p-4"><div className="mb-3 text-lg font-semibold text-slate-100">Events</div><div className="space-y-2">{events.length ? events.map((event: FraudViolationRecord) => <div key={event.id ?? `${event.score ?? 'score'}-${event.created_at ?? 'date'}`} className="rounded-lg border border-slate-800 bg-slate-950/60 p-3 text-sm"><div className="flex items-center justify-between gap-2"><span className="font-medium text-slate-100">Score: {event.score ?? '—'}</span><span className="rounded-full border border-amber-700 bg-amber-900/40 px-2 py-1 text-[10px] uppercase text-amber-200">{event.level ?? 'unknown'}</span></div><div className="mt-1 text-slate-400">{event.created_at ?? '—'}</div></div>) : <div className="text-slate-400">No fraud events available.</div>}</div></section></div>;
}

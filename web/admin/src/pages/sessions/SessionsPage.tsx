import { useMemo, useState } from 'react';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { Button } from '../../components/ui/Button';
import { Input } from '../../components/ui/Input';
import { fetchSessions, revokeSession } from '../../services/sessions';

const pageSize = 10;

export function SessionsPage() {
  const queryClient = useQueryClient();
  const [search, setSearch] = useState('');
  const [stateFilter, setStateFilter] = useState('');
  const [page, setPage] = useState(1);
  const sessionsQuery = useQuery({ queryKey: ['sessions'], queryFn: fetchSessions });
  const revokeMutation = useMutation({
    mutationFn: revokeSession,
    onSuccess: () => void queryClient.invalidateQueries({ queryKey: ['sessions'] }),
  });
  const sessions = sessionsQuery.data ?? [];
  const filtered = useMemo(() => {
    const term = search.trim().toLowerCase();
    return sessions.filter((session) => (!stateFilter || (stateFilter === 'revoked' ? Boolean(session.revoked_at) : !session.revoked_at)) && (!term || [session.id, session.user_id, session.tenant_id, session.ip_address, session.user_agent].filter(Boolean).join(' ').toLowerCase().includes(term)));
  }, [search, sessions, stateFilter]);
  const pageCount = Math.max(1, Math.ceil(filtered.length / pageSize));
  const visible = filtered.slice((page - 1) * pageSize, page * pageSize);

  if (sessionsQuery.isLoading) return <div role="status" className="border border-slate-800 bg-slate-900 p-6 text-slate-300">Loading sessions…</div>;
  if (sessionsQuery.isError) return <div role="alert" className="border border-rose-900 bg-rose-950/30 p-6 text-rose-200">Unable to load sessions: {sessionsQuery.error.message}</div>;

  return <div className="space-y-4">
    <div className="flex flex-col gap-3 border-b border-slate-800 pb-4 sm:flex-row"><Input aria-label="Search sessions" placeholder="Search session, user, tenant, IP or agent" value={search} onChange={(event) => { setSearch(event.target.value); setPage(1); }} /><select aria-label="Filter session state" value={stateFilter} onChange={(event) => { setStateFilter(event.target.value); setPage(1); }} className="rounded-md border border-slate-700 bg-slate-950 px-3 py-2 text-sm"><option value="">All sessions</option><option value="current">Current</option><option value="revoked">Revoked</option></select></div>
    {revokeMutation.isError ? <div role="alert" className="border border-rose-900 bg-rose-950/30 p-3 text-sm text-rose-200">Unable to revoke session: {revokeMutation.error.message}</div> : null}
    <div className="overflow-x-auto border border-slate-800"><table className="min-w-full text-left text-sm text-slate-200"><thead className="bg-slate-900 text-xs uppercase text-slate-400"><tr>{['Session ID', 'User ID', 'Tenant ID', 'IP address', 'User agent', 'Created', 'Expires', 'Revoked', 'Action'].map((heading) => <th key={heading} className="px-3 py-3">{heading}</th>)}</tr></thead><tbody>{visible.map((session) => <tr key={session.id} className="border-t border-slate-800"><td className="px-3 py-3">{session.id ?? '—'}</td><td className="px-3 py-3">{session.user_id ?? '—'}</td><td className="px-3 py-3">{session.tenant_id ?? '—'}</td><td className="px-3 py-3">{session.ip_address ?? '—'}</td><td className="max-w-xs break-words px-3 py-3">{session.user_agent ?? '—'}</td><td className="px-3 py-3">{session.created_at ?? '—'}</td><td className="px-3 py-3">{session.expires_at ?? '—'}</td><td className="px-3 py-3">{session.revoked_at ?? 'Active'}</td><td className="px-3 py-3"><Button type="button" className="border-rose-800 px-2 py-1 text-xs text-rose-200" disabled={Boolean(session.revoked_at) || revokeMutation.isPending || !session.id} onClick={() => { if (session.id && window.confirm(`Revoke session ${session.id}?`)) revokeMutation.mutate(session.id); }}>{revokeMutation.isPending ? 'Revoking…' : 'Revoke'}</Button></td></tr>)}{visible.length === 0 ? <tr><td colSpan={9} className="px-4 py-12 text-center text-slate-400">{sessions.length ? 'No sessions match the current filters.' : 'No sessions are recorded for this tenant.'}</td></tr> : null}</tbody></table></div>
    <div className="flex items-center justify-between text-sm text-slate-400"><span>{filtered.length} results · page {page} of {pageCount}</span><div className="flex gap-2"><Button type="button" disabled={page <= 1} onClick={() => setPage((value) => value - 1)}>Previous</Button><Button type="button" disabled={page >= pageCount} onClick={() => setPage((value) => value + 1)}>Next</Button></div></div>
  </div>;
}

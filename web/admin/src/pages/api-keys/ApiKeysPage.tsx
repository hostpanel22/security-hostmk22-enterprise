import { useMemo, useState, type FormEvent } from 'react';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { Button } from '../../components/ui/Button';
import { Input } from '../../components/ui/Input';
import { createApiKey, fetchApiKeys, revokeApiKey, rotateApiKey, type ApiKeyRecord } from '../../services/apiKeys';

const pageSize = 10;

export function ApiKeysPage() {
  const client = useQueryClient();
  const [search, setSearch] = useState('');
  const [stateFilter, setStateFilter] = useState('');
  const [page, setPage] = useState(1);
  const [formOpen, setFormOpen] = useState(false);
  const [newKey, setNewKey] = useState<{ name?: string; key?: string } | null>(null);
  const keysQuery = useQuery({ queryKey: ['api-keys'], queryFn: fetchApiKeys });
  const refresh = () => client.invalidateQueries({ queryKey: ['api-keys'] });
  const createMutation = useMutation({ mutationFn: createApiKey, onSuccess: (response) => { setFormOpen(false); setNewKey({ name: response.data?.name, key: response.data?.key }); void refresh(); } });
  const rotateMutation = useMutation({ mutationFn: rotateApiKey, onSuccess: (response) => { setNewKey({ name: 'Rotated API key', key: response.data?.key }); void refresh(); } });
  const revokeMutation = useMutation({ mutationFn: revokeApiKey, onSuccess: () => void refresh() });
  const keys = keysQuery.data ?? [];
  const filtered = useMemo(() => {
    const term = search.trim().toLowerCase();
    return keys.filter((key) => (!stateFilter || (stateFilter === 'revoked' ? Boolean(key.revoked_at) : !key.revoked_at && (!key.expires_at || new Date(key.expires_at).getTime() > Date.now()))) && (!term || [key.id, key.tenant_id, key.name, key.prefix, ...(key.scopes ?? [])].filter(Boolean).join(' ').toLowerCase().includes(term)));
  }, [keys, search, stateFilter]);
  const pageCount = Math.max(1, Math.ceil(filtered.length / pageSize));
  const shown = filtered.slice((page - 1) * pageSize, page * pageSize);
  const mutationError = createMutation.error ?? rotateMutation.error ?? revokeMutation.error;

  const submit = (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    const form = new FormData(event.currentTarget);
    const rawLimit = String(form.get('rate_limit') ?? '');
    const expiration = String(form.get('expires_at') ?? '');
    createMutation.mutate({ name: String(form.get('name') ?? ''), scopes: String(form.get('scopes') ?? '').split(',').map((scope) => scope.trim()).filter(Boolean), ...(rawLimit ? { rate_limit: Number(rawLimit) } : {}), ...(expiration ? { expires_at: new Date(expiration).toISOString() } : {}) });
  };

  if (keysQuery.isLoading) return <div role="status" className="border border-slate-800 bg-slate-900 p-6 text-slate-300">Loading API keys…</div>;
  if (keysQuery.isError) return <div role="alert" className="border border-rose-900 bg-rose-950/30 p-6 text-rose-200">Unable to load API keys: {keysQuery.error.message}</div>;

  return <div className="space-y-4">
    <div className="flex flex-col gap-3 border-b border-slate-800 pb-4 md:flex-row"><Input aria-label="Search API keys" placeholder="Search name, ID, tenant, prefix or scope" value={search} onChange={(event) => { setSearch(event.target.value); setPage(1); }} /><select aria-label="Filter API key status" value={stateFilter} onChange={(event) => { setStateFilter(event.target.value); setPage(1); }} className="rounded-md border border-slate-700 bg-slate-950 px-3 py-2 text-sm"><option value="">All keys</option><option value="active">Active</option><option value="revoked">Revoked</option></select><Button type="button" onClick={() => setFormOpen((open) => !open)}>{formOpen ? 'Close' : 'Create API key'}</Button></div>
    {formOpen ? <form onSubmit={submit} className="grid gap-3 border border-slate-800 bg-slate-900/70 p-4 md:grid-cols-2"><Input name="name" aria-label="Key name" placeholder="Key name" required /><Input name="scopes" aria-label="Scopes" placeholder="Scopes separated by commas" /><Input name="rate_limit" aria-label="Rate limit" type="number" min="1" placeholder="Rate limit (optional)" /><Input name="expires_at" aria-label="Expiration" type="datetime-local" /><Button type="submit" disabled={createMutation.isPending}>{createMutation.isPending ? 'Creating…' : 'Create API key'}</Button></form> : null}
    {newKey ? <section className="border border-amber-700 bg-amber-950/30 p-4 text-amber-100"><div className="flex items-start justify-between gap-3"><div><h2 className="font-semibold">Copy this key now</h2><p className="mt-1 text-xs">The plaintext is only returned by create and rotate operations.</p></div><Button type="button" onClick={() => setNewKey(null)}>Dismiss</Button></div><p className="mt-3 break-all font-mono text-sm">{newKey.key || 'The server did not return a plaintext key.'}</p></section> : null}
    {mutationError ? <div role="alert" className="border border-rose-900 bg-rose-950/30 p-3 text-sm text-rose-200">Operation failed: {mutationError.message}</div> : null}
    <div className="overflow-x-auto border border-slate-800"><table className="min-w-full text-left text-sm"><thead className="bg-slate-900 text-xs uppercase text-slate-400"><tr>{['ID', 'Tenant', 'Name', 'Prefix', 'Scopes', 'Rate limit', 'Created', 'Updated', 'Last used', 'Expires', 'Revoked', 'Actions'].map((heading) => <th key={heading} className="px-3 py-3">{heading}</th>)}</tr></thead><tbody>{shown.map((key: ApiKeyRecord) => <tr key={key.id} className="border-t border-slate-800"><td className="px-3 py-3">{key.id}</td><td className="px-3 py-3">{key.tenant_id ?? '—'}</td><td className="px-3 py-3">{key.name}</td><td className="px-3 py-3">{key.prefix}</td><td className="px-3 py-3">{key.scopes?.join(', ') || '—'}</td><td className="px-3 py-3">{key.rate_limit ?? '—'}</td><td className="px-3 py-3">{key.created_at ?? '—'}</td><td className="px-3 py-3">{key.updated_at ?? '—'}</td><td className="px-3 py-3">{key.last_used_at ?? '—'}</td><td className="px-3 py-3">{key.expires_at ?? 'Never'}</td><td className="px-3 py-3">{key.revoked_at ?? 'Active'}</td><td className="space-x-2 px-3 py-3"><Button type="button" disabled={!key.id || Boolean(key.revoked_at) || rotateMutation.isPending} className="px-2 py-1 text-xs" onClick={() => key.id && window.confirm(`Rotate key ${key.name}? The old key will stop working.`) && rotateMutation.mutate(key.id)}>Rotate</Button><Button type="button" disabled={!key.id || Boolean(key.revoked_at) || revokeMutation.isPending} className="border-rose-800 px-2 py-1 text-xs text-rose-200" onClick={() => key.id && window.confirm(`Revoke key ${key.name}?`) && revokeMutation.mutate(key.id)}>Revoke</Button></td></tr>)}{shown.length === 0 ? <tr><td colSpan={12} className="px-4 py-12 text-center text-slate-400">{keys.length ? 'No keys match the selected filters.' : 'No API keys have been created for this tenant.'}</td></tr> : null}</tbody></table></div>
    <div className="flex items-center justify-between text-sm text-slate-400"><span>{filtered.length} keys · page {page} of {pageCount}</span><div className="flex gap-2"><Button type="button" disabled={page <= 1} onClick={() => setPage((value) => value - 1)}>Previous</Button><Button type="button" disabled={page >= pageCount} onClick={() => setPage((value) => value + 1)}>Next</Button></div></div>
  </div>;
}

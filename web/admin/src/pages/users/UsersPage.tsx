import { useMemo, useState, type FormEvent } from 'react';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { Button } from '../../components/ui/Button';
import { Input } from '../../components/ui/Input';
import { fetchRoles } from '../../services/rbac';
import { createUser, fetchUsers, updateUser, userAction, type UserRecord } from '../../services/users';

const pageSize = 10;
const statuses = ['active', 'pending', 'inactive', 'suspended', 'blocked'];

export function UsersPage() {
  const queryClient = useQueryClient();
  const [search, setSearch] = useState('');
  const [statusFilter, setStatusFilter] = useState('');
  const [page, setPage] = useState(1);
  const [createOpen, setCreateOpen] = useState(false);
  const [editing, setEditing] = useState<UserRecord | null>(null);
  const [mfaCredentials, setMfaCredentials] = useState<{ secret?: string; otpauth_url?: string } | null>(null);
  const usersQuery = useQuery({ queryKey: ['users'], queryFn: fetchUsers });
  const rolesQuery = useQuery({ queryKey: ['rbac', 'roles'], queryFn: fetchRoles });
  const refreshUsers = () => queryClient.invalidateQueries({ queryKey: ['users'] });
  const createMutation = useMutation({ mutationFn: createUser, onSuccess: () => { setCreateOpen(false); void refreshUsers(); } });
  const updateMutation = useMutation({ mutationFn: ({ id, user }: { id: string; user: Parameters<typeof updateUser>[1] }) => updateUser(id, user), onSuccess: () => { setEditing(null); void refreshUsers(); } });
  const actionMutation = useMutation({ mutationFn: ({ id, action }: { id: string; action: Parameters<typeof userAction>[1] }) => userAction(id, action), onSuccess: (result, variables) => { if (variables.action === 'enable-mfa') setMfaCredentials(result); void refreshUsers(); } });
  const users = usersQuery.data ?? [];
  const filteredUsers = useMemo(() => {
    const term = search.trim().toLowerCase();
    return users.filter((user) => (!statusFilter || user.status === statusFilter) && (!term || [user.id, user.name, user.email, user.tenant_id, ...(user.role_ids ?? [])].filter(Boolean).join(' ').toLowerCase().includes(term)));
  }, [search, statusFilter, users]);
  const pageCount = Math.max(1, Math.ceil(filteredUsers.length / pageSize));
  const visibleUsers = filteredUsers.slice((page - 1) * pageSize, page * pageSize);
  const mutationError = createMutation.error ?? updateMutation.error ?? actionMutation.error;

  const submitCreate = (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    const form = new FormData(event.currentTarget);
    createMutation.mutate({
      name: String(form.get('name') ?? ''),
      email: String(form.get('email') ?? ''),
      password: String(form.get('password') ?? ''),
      status: String(form.get('status') ?? 'active'),
      role_ids: form.getAll('role_ids').map(String),
    });
  };

  const submitUpdate = (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    if (!editing?.id) return;
    const form = new FormData(event.currentTarget);
    updateMutation.mutate({
      id: editing.id,
      user: { name: String(form.get('name') ?? ''), status: String(form.get('status') ?? ''), role_ids: form.getAll('role_ids').map(String) },
    });
  };

  const confirmAction = (user: UserRecord, action: Parameters<typeof userAction>[1], prompt: string) => {
    if (user.id && window.confirm(prompt)) actionMutation.mutate({ id: user.id, action });
  };

  if (usersQuery.isLoading) return <div role="status" className="border border-slate-800 bg-slate-900 p-6 text-slate-300">Loading users…</div>;
  if (usersQuery.isError) return <div role="alert" className="border border-rose-900 bg-rose-950/30 p-6 text-rose-200">Unable to load users: {usersQuery.error.message}</div>;

  return (
    <div className="space-y-4">
      <div className="flex flex-col gap-3 border-b border-slate-800 pb-4 sm:flex-row sm:items-center sm:justify-between">
        <div className="flex flex-1 flex-col gap-3 sm:flex-row">
          <Input aria-label="Search users" placeholder="Search name, email, tenant, role or ID" value={search} onChange={(event) => { setSearch(event.target.value); setPage(1); }} />
          <select aria-label="Filter by status" className="rounded-md border border-slate-700 bg-slate-950 px-3 py-2 text-sm" value={statusFilter} onChange={(event) => { setStatusFilter(event.target.value); setPage(1); }}>
            <option value="">All statuses</option>{statuses.map((status) => <option key={status} value={status}>{status}</option>)}
          </select>
        </div>
        <Button type="button" onClick={() => setCreateOpen((open) => !open)}>{createOpen ? 'Close form' : 'Create user'}</Button>
      </div>

      {createOpen ? <form onSubmit={submitCreate} className="grid gap-3 border border-slate-800 bg-slate-900/70 p-4 md:grid-cols-2">
        <Input name="name" aria-label="Name" placeholder="Name" required />
        <Input name="email" aria-label="Email" type="email" placeholder="Email" required />
        <Input name="password" aria-label="Temporary password" type="password" placeholder="Temporary password" autoComplete="new-password" required minLength={8} />
        <select name="status" aria-label="Initial status" defaultValue="active" className="rounded-md border border-slate-700 bg-slate-950 px-3 py-2 text-sm"><option value="active">Active</option><option value="pending">Pending</option></select>
        <label className="space-y-1 text-sm text-slate-300 md:col-span-2"><span>Roles</span><select name="role_ids" multiple aria-label="Roles" className="min-h-24 w-full rounded-md border border-slate-700 bg-slate-950 p-2">{(rolesQuery.data ?? []).map((role) => <option key={role.id} value={role.id}>{role.name}</option>)}</select></label>
        {rolesQuery.isError ? <p role="alert" className="text-sm text-rose-200 md:col-span-2">Roles could not be loaded: {rolesQuery.error.message}</p> : null}
        <Button type="submit" disabled={createMutation.isPending}>{createMutation.isPending ? 'Creating…' : 'Create user'}</Button>
      </form> : null}

      {editing ? <form onSubmit={submitUpdate} className="grid gap-3 border border-sky-800 bg-slate-900/70 p-4 md:grid-cols-2">
        <h2 className="font-semibold md:col-span-2">Edit {editing.email}</h2>
        <Input name="name" aria-label="Edit name" defaultValue={editing.name ?? ''} required />
        <select name="status" aria-label="Edit status" defaultValue={editing.status ?? 'active'} className="rounded-md border border-slate-700 bg-slate-950 px-3 py-2 text-sm">{statuses.map((status) => <option key={status} value={status}>{status}</option>)}</select>
        <label className="space-y-1 text-sm text-slate-300 md:col-span-2"><span>Roles</span><select name="role_ids" multiple aria-label="Edit roles" defaultValue={editing.role_ids ?? []} className="min-h-24 w-full rounded-md border border-slate-700 bg-slate-950 p-2">{(rolesQuery.data ?? []).map((role) => <option key={role.id} value={role.id}>{role.name}</option>)}</select></label>
        <Button type="submit" disabled={updateMutation.isPending}>{updateMutation.isPending ? 'Saving…' : 'Save changes'}</Button>
        <Button type="button" onClick={() => setEditing(null)}>Cancel</Button>
      </form> : null}

      {mutationError ? <div role="alert" className="border border-rose-900 bg-rose-950/30 p-3 text-sm text-rose-200">Action failed: {mutationError.message}</div> : null}
      {mfaCredentials ? <section className="border border-amber-700 bg-amber-950/30 p-4 text-sm text-amber-100"><div className="flex justify-between"><h2 className="font-semibold">MFA enrollment secret</h2><Button type="button" onClick={() => setMfaCredentials(null)}>Dismiss</Button></div><p className="mt-2 break-all">{mfaCredentials.secret}</p><p className="mt-1 break-all text-xs">{mfaCredentials.otpauth_url}</p></section> : null}

      <div className="overflow-x-auto border border-slate-800">
        <table className="min-w-full text-left text-sm text-slate-200"><thead className="bg-slate-900 text-xs uppercase text-slate-400"><tr>{['User', 'Status', 'Tenant', 'Roles', 'MFA', 'Email verified', 'Last login', 'Created', 'Updated', 'Actions'].map((heading) => <th key={heading} className="px-3 py-3">{heading}</th>)}</tr></thead>
          <tbody>{visibleUsers.map((user) => <tr key={user.id} className="border-t border-slate-800 align-top"><td className="px-3 py-3"><div className="font-medium">{user.name || '—'}</div><div className="text-xs text-slate-400">{user.email}</div><div className="text-xs text-slate-500">{user.id}</div></td><td className="px-3 py-3">{user.status ?? 'unknown'}</td><td className="px-3 py-3">{user.tenant_id ?? '—'}</td><td className="px-3 py-3">{user.role_ids?.join(', ') || '—'}</td><td className="px-3 py-3">{user.mfa_enabled ? 'Enabled' : 'Disabled'}</td><td className="px-3 py-3">{user.email_verified_at ?? 'Not verified'}</td><td className="px-3 py-3">{user.last_login ?? '—'}</td><td className="px-3 py-3">{user.created_at ?? '—'}</td><td className="px-3 py-3">{user.updated_at ?? '—'}</td>
            <td className="min-w-48 space-y-2 px-3 py-3"><div className="flex flex-wrap gap-2"><Button type="button" className="px-2 py-1 text-xs" onClick={() => setEditing(user)}>Edit</Button><Button type="button" className="px-2 py-1 text-xs" onClick={() => confirmAction(user, user.status === 'active' ? 'suspend' : 'activate', `${user.status === 'active' ? 'Suspend' : 'Activate'} ${user.email}?`)}>{user.status === 'active' ? 'Suspend' : 'Activate'}</Button><Button type="button" className="px-2 py-1 text-xs" disabled={Boolean(user.email_verified_at)} onClick={() => confirmAction(user, 'verify-email', `Mark ${user.email} as verified?`)}>Verify email</Button><Button type="button" className="px-2 py-1 text-xs" disabled={Boolean(user.mfa_enabled)} onClick={() => confirmAction(user, 'enable-mfa', `Enable MFA for ${user.email}?`)}>Enable MFA</Button><Button type="button" className="border-rose-800 px-2 py-1 text-xs text-rose-200" onClick={() => confirmAction(user, 'delete', `Permanently delete ${user.email}?`)}>Delete</Button></div></td>
          </tr>)}
          {visibleUsers.length === 0 ? <tr><td colSpan={10} className="px-4 py-12 text-center text-slate-400">{users.length === 0 ? 'No users exist in this tenant yet.' : 'No users match the current search and filters.'}</td></tr> : null}</tbody>
        </table>
      </div>
      <div className="flex items-center justify-between text-sm text-slate-400"><span>{filteredUsers.length} results · page {page} of {pageCount}</span><div className="flex gap-2"><Button type="button" disabled={page <= 1} onClick={() => setPage((value) => value - 1)}>Previous</Button><Button type="button" disabled={page >= pageCount} onClick={() => setPage((value) => value + 1)}>Next</Button></div></div>
    </div>
  );
}
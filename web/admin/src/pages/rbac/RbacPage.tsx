import { useMemo, useState, type FormEvent } from 'react';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { Button } from '../../components/ui/Button';
import { Input } from '../../components/ui/Input';
import { checkPermission, createPermission, deletePermission, fetchPermissions } from '../../services/permissions';
import { createRole, deleteRole, fetchRoles, updateRole, type RoleRecord } from '../../services/rbac';

const pageSize = 10;

export function RbacPage() {
  const client = useQueryClient();
  const [roleSearch, setRoleSearch] = useState('');
  const [permissionSearch, setPermissionSearch] = useState('');
  const [rolePage, setRolePage] = useState(1);
  const [permissionPage, setPermissionPage] = useState(1);
  const [roleFormOpen, setRoleFormOpen] = useState(false);
  const [editingRole, setEditingRole] = useState<RoleRecord | null>(null);
  const [permissionFormOpen, setPermissionFormOpen] = useState(false);
  const [permissionResult, setPermissionResult] = useState<{ user_id: string; permission: string; allowed: boolean } | null>(null);
  const rolesQuery = useQuery({ queryKey: ['rbac', 'roles'], queryFn: fetchRoles });
  const permissionsQuery = useQuery({ queryKey: ['permissions'], queryFn: fetchPermissions });
  const refreshRoles = () => client.invalidateQueries({ queryKey: ['rbac', 'roles'] });
  const refreshPermissions = () => { void client.invalidateQueries({ queryKey: ['permissions'] }); void refreshRoles(); };
  const createRoleMutation = useMutation({ mutationFn: createRole, onSuccess: () => { setRoleFormOpen(false); void refreshRoles(); } });
  const updateRoleMutation = useMutation({ mutationFn: updateRole, onSuccess: () => { setEditingRole(null); void refreshRoles(); } });
  const deleteRoleMutation = useMutation({ mutationFn: deleteRole, onSuccess: () => void refreshRoles() });
  const createPermissionMutation = useMutation({ mutationFn: createPermission, onSuccess: () => { setPermissionFormOpen(false); refreshPermissions(); } });
  const deletePermissionMutation = useMutation({ mutationFn: deletePermission, onSuccess: refreshPermissions });
  const checkMutation = useMutation({ mutationFn: ({ userId, permission }: { userId: string; permission: string }) => checkPermission(userId, permission), onSuccess: setPermissionResult });
  const roles = rolesQuery.data ?? [];
  const permissions = permissionsQuery.data ?? [];
  const matchingRoles = useMemo(() => roles.filter((role) => [role.id, role.name, role.description, role.tenant_id, ...(role.permissions ?? [])].filter(Boolean).join(' ').toLowerCase().includes(roleSearch.trim().toLowerCase())), [roleSearch, roles]);
  const matchingPermissions = useMemo(() => permissions.filter((permission) => [permission.id, permission.name, permission.resource, permission.action, permission.description, permission.tenant_id].filter(Boolean).join(' ').toLowerCase().includes(permissionSearch.trim().toLowerCase())), [permissionSearch, permissions]);
  const rolePages = Math.max(1, Math.ceil(matchingRoles.length / pageSize));
  const permissionPages = Math.max(1, Math.ceil(matchingPermissions.length / pageSize));
  const shownRoles = matchingRoles.slice((rolePage - 1) * pageSize, rolePage * pageSize);
  const shownPermissions = matchingPermissions.slice((permissionPage - 1) * pageSize, permissionPage * pageSize);
  const mutationError = createRoleMutation.error ?? updateRoleMutation.error ?? deleteRoleMutation.error ?? createPermissionMutation.error ?? deletePermissionMutation.error ?? checkMutation.error;

  const submitRole = (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    const form = new FormData(event.currentTarget);
    const values = { name: String(form.get('name') ?? ''), description: String(form.get('description') ?? ''), permissions: form.getAll('permissions').map(String) };
    if (editingRole?.id) updateRoleMutation.mutate({ ...editingRole, ...values, id: editingRole.id });
    else createRoleMutation.mutate(values);
  };

  const submitPermission = (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    const form = new FormData(event.currentTarget);
    createPermissionMutation.mutate({ name: String(form.get('name') ?? ''), resource: String(form.get('resource') ?? ''), action: String(form.get('action') ?? ''), description: String(form.get('description') ?? '') });
  };

  const runCheck = (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    const form = new FormData(event.currentTarget);
    checkMutation.mutate({ userId: String(form.get('user_id') ?? ''), permission: String(form.get('permission') ?? '') });
  };

  if (rolesQuery.isLoading || permissionsQuery.isLoading) return <div role="status" className="border border-slate-800 bg-slate-900 p-6 text-slate-300">Loading roles and permissions…</div>;
  if (rolesQuery.isError || permissionsQuery.isError) return <div role="alert" className="border border-rose-900 bg-rose-950/30 p-6 text-rose-200">Unable to load RBAC data: {(rolesQuery.error ?? permissionsQuery.error)?.message}</div>;

  return <div className="space-y-6">
    {mutationError ? <div role="alert" className="border border-rose-900 bg-rose-950/30 p-3 text-sm text-rose-200">Operation failed: {mutationError.message}</div> : null}
    <form onSubmit={runCheck} className="flex flex-col gap-3 border-b border-slate-800 pb-4 md:flex-row"><Input name="user_id" aria-label="User ID to check" placeholder="User ID" required /><Input name="permission" aria-label="Permission to check" placeholder="Permission name" required /><Button type="submit" disabled={checkMutation.isPending}>{checkMutation.isPending ? 'Checking…' : 'Check access'}</Button></form>
    {permissionResult ? <p role="status" className="text-sm text-slate-300">{permissionResult.user_id} · {permissionResult.permission}: <strong className={permissionResult.allowed ? 'text-emerald-300' : 'text-rose-300'}>{permissionResult.allowed ? 'Allowed' : 'Denied'}</strong></p> : null}

    <section className="space-y-3">
      <div className="flex flex-col gap-3 sm:flex-row sm:items-center"><h2 className="text-lg font-semibold">Roles</h2><Input aria-label="Search roles" placeholder="Search roles" value={roleSearch} onChange={(event) => { setRoleSearch(event.target.value); setRolePage(1); }} /><Button type="button" onClick={() => { setEditingRole(null); setRoleFormOpen((open) => !open); }}>{roleFormOpen ? 'Close' : 'Create role'}</Button></div>
      {roleFormOpen || editingRole ? <form onSubmit={submitRole} className="grid gap-3 border border-slate-800 bg-slate-900/70 p-4 md:grid-cols-2"><h3 className="font-medium md:col-span-2">{editingRole ? `Edit ${editingRole.name}` : 'New role'}</h3><Input name="name" aria-label="Role name" placeholder="Role name" defaultValue={editingRole?.name ?? ''} required /><Input name="description" aria-label="Role description" placeholder="Description" defaultValue={editingRole?.description ?? ''} /><label className="space-y-1 text-sm text-slate-300 md:col-span-2"><span>Assigned permissions</span><select name="permissions" multiple aria-label="Assigned permissions" defaultValue={editingRole?.permissions ?? []} className="min-h-28 w-full rounded-md border border-slate-700 bg-slate-950 p-2">{permissions.map((permission) => <option key={permission.id} value={permission.name}>{permission.name} · {permission.resource}:{permission.action}</option>)}</select></label><Button type="submit" disabled={createRoleMutation.isPending || updateRoleMutation.isPending}>{editingRole ? 'Save role' : 'Create role'}</Button><Button type="button" onClick={() => { setEditingRole(null); setRoleFormOpen(false); }}>Cancel</Button></form> : null}
      <div className="overflow-x-auto border border-slate-800"><table className="min-w-full text-left text-sm"><thead className="bg-slate-900 text-xs uppercase text-slate-400"><tr>{['ID', 'Name', 'Description', 'Permissions', 'Tenant', 'Created', 'Updated', 'Actions'].map((heading) => <th key={heading} className="px-3 py-3">{heading}</th>)}</tr></thead><tbody>{shownRoles.map((role) => <tr key={role.id} className="border-t border-slate-800"><td className="px-3 py-3">{role.id}</td><td className="px-3 py-3">{role.name}</td><td className="px-3 py-3">{role.description ?? '—'}</td><td className="px-3 py-3">{role.permissions?.join(', ') || '—'}</td><td className="px-3 py-3">{role.tenant_id ?? '—'}</td><td className="px-3 py-3">{role.created_at ?? '—'}</td><td className="px-3 py-3">{role.updated_at ?? '—'}</td><td className="space-x-2 px-3 py-3"><Button type="button" className="px-2 py-1 text-xs" onClick={() => { setRoleFormOpen(false); setEditingRole(role); }}>Edit</Button><Button type="button" className="border-rose-800 px-2 py-1 text-xs text-rose-200" onClick={() => role.id && window.confirm(`Delete role ${role.name}?`) && deleteRoleMutation.mutate(role.id)}>Delete</Button></td></tr>)}{shownRoles.length === 0 ? <tr><td colSpan={8} className="px-4 py-10 text-center text-slate-400">{roles.length ? 'No roles match the search.' : 'No roles returned by the API.'}</td></tr> : null}</tbody></table></div>
      <div className="flex items-center justify-between text-sm text-slate-400"><span>{matchingRoles.length} roles · page {rolePage} of {rolePages}</span><div className="flex gap-2"><Button type="button" disabled={rolePage <= 1} onClick={() => setRolePage((value) => value - 1)}>Previous</Button><Button type="button" disabled={rolePage >= rolePages} onClick={() => setRolePage((value) => value + 1)}>Next</Button></div></div>
    </section>

    <section className="space-y-3 border-t border-slate-800 pt-5">
      <div className="flex flex-col gap-3 sm:flex-row sm:items-center"><h2 className="text-lg font-semibold">Permissions</h2><Input aria-label="Search permissions" placeholder="Search permissions" value={permissionSearch} onChange={(event) => { setPermissionSearch(event.target.value); setPermissionPage(1); }} /><Button type="button" onClick={() => setPermissionFormOpen((open) => !open)}>{permissionFormOpen ? 'Close' : 'Create permission'}</Button></div>
      {permissionFormOpen ? <form onSubmit={submitPermission} className="grid gap-3 border border-slate-800 bg-slate-900/70 p-4 md:grid-cols-2"><Input name="name" aria-label="Permission name" placeholder="Permission name" required /><Input name="resource" aria-label="Permission resource" placeholder="Resource" required /><Input name="action" aria-label="Permission action" placeholder="Action" required /><Input name="description" aria-label="Permission description" placeholder="Description" /><Button type="submit" disabled={createPermissionMutation.isPending}>{createPermissionMutation.isPending ? 'Creating…' : 'Create permission'}</Button></form> : null}
      <div className="overflow-x-auto border border-slate-800"><table className="min-w-full text-left text-sm"><thead className="bg-slate-900 text-xs uppercase text-slate-400"><tr>{['ID', 'Name', 'Resource', 'Action', 'Description', 'Tenant', 'Created', 'Updated', 'Actions'].map((heading) => <th key={heading} className="px-3 py-3">{heading}</th>)}</tr></thead><tbody>{shownPermissions.map((permission) => <tr key={permission.id} className="border-t border-slate-800"><td className="px-3 py-3">{permission.id}</td><td className="px-3 py-3">{permission.name}</td><td className="px-3 py-3">{permission.resource}</td><td className="px-3 py-3">{permission.action}</td><td className="px-3 py-3">{permission.description ?? '—'}</td><td className="px-3 py-3">{permission.tenant_id ?? '—'}</td><td className="px-3 py-3">{permission.created_at ?? '—'}</td><td className="px-3 py-3">{permission.updated_at ?? '—'}</td><td className="px-3 py-3"><Button type="button" className="border-rose-800 px-2 py-1 text-xs text-rose-200" onClick={() => permission.id && window.confirm(`Delete permission ${permission.name}?`) && deletePermissionMutation.mutate(permission.id)}>Delete</Button></td></tr>)}{shownPermissions.length === 0 ? <tr><td colSpan={9} className="px-4 py-10 text-center text-slate-400">{permissions.length ? 'No permissions match the search.' : 'No permissions returned by the API.'}</td></tr> : null}</tbody></table></div>
      <div className="flex items-center justify-between text-sm text-slate-400"><span>{matchingPermissions.length} permissions · page {permissionPage} of {permissionPages}</span><div className="flex gap-2"><Button type="button" disabled={permissionPage <= 1} onClick={() => setPermissionPage((value) => value - 1)}>Previous</Button><Button type="button" disabled={permissionPage >= permissionPages} onClick={() => setPermissionPage((value) => value + 1)}>Next</Button></div></div>
    </section>
  </div>;
}

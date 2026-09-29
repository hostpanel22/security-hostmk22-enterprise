export type PermissionName = string;

export const hasPermission = (permissions: PermissionName[], permission: string) =>
  permissions.includes(permission);

export const hasAnyPermission = (permissions: PermissionName[], values: string[]) =>
  values.some((permission) => permissions.includes(permission));

package permissions

type Subject struct {
	UserID string

	TenantID string

	Roles []*Role
}

type Checker struct {
}

func NewChecker() *Checker {
	return &Checker{}
}

func (c *Checker) Can(
	subject *Subject,
	resource string,
	action string,
) bool {

	if subject == nil {
		return false
	}

	for _, role := range subject.Roles {

		if role == nil {
			continue
		}

		if role.HasPermission(
			resource,
			action,
		) {
			return true
		}
	}

	return false
}

func (c *Checker) CanAny(
	subject *Subject,
	resource string,
	actions ...string,
) bool {

	if subject == nil {
		return false
	}

	for _, action := range actions {

		if c.Can(
			subject,
			resource,
			action,
		) {
			return true
		}
	}

	return false
}

func (c *Checker) CanAll(
	subject *Subject,
	resource string,
	actions ...string,
) bool {

	if subject == nil {
		return false
	}

	for _, action := range actions {

		if !c.Can(
			subject,
			resource,
			action,
		) {
			return false
		}
	}

	return true
}

func (c *Checker) HasRole(
	subject *Subject,
	roleName string,
) bool {

	if subject == nil {
		return false
	}

	roleName = normalize(
		roleName,
	)

	for _, role := range subject.Roles {

		if role == nil {
			continue
		}

		if normalize(
			role.Name,
		) == roleName {

			return true
		}
	}

	return false
}

func (c *Checker) HasPermission(
	subject *Subject,
	permissionKey string,
) bool {

	if subject == nil {
		return false
	}

	for _, role := range subject.Roles {

		if role == nil {
			continue
		}

		for _, permission := range role.Permissions {

			if permission == nil {
				continue
			}

			if permission.Key() == normalize(permissionKey) {
				return true
			}
		}
	}

	return false
}

func (c *Checker) RoleCount(
	subject *Subject,
) int {

	if subject == nil {
		return 0
	}

	return len(
		subject.Roles,
	)
}

func (c *Checker) PermissionCount(
	subject *Subject,
) int {

	if subject == nil {
		return 0
	}

	unique := make(
		map[string]struct{},
	)

	for _, role := range subject.Roles {

		if role == nil {
			continue
		}

		for _, permission := range role.Permissions {

			if permission == nil {
				continue
			}

			unique[permission.Key()] = struct{}{}
		}
	}

	return len(
		unique,
	)
}

func (c *Checker) ValidateTenant(
	subject *Subject,
	tenantID string,
) bool {

	if subject == nil {
		return false
	}

	return subject.TenantID ==
		tenantID
}

func (c *Checker) IsSuperAdmin(
	subject *Subject,
) bool {

	if subject == nil {
		return false
	}

	return c.HasRole(
		subject,
		"super-admin",
	)
}

package audit

import "time"

type Filter struct {
	ID string

	TenantID string

	ActorID string

	Category Category

	Action string

	Status Status

	Severity Severity

	Resource string

	Tag string

	From *time.Time

	To *time.Time
}

func NewFilter() Filter {
	return Filter{}
}

func (f Filter) Match(event *Event) bool {

	if event == nil {
		return false
	}

	if f.ID != "" && event.ID != f.ID {
		return false
	}

	if f.TenantID != "" &&
		event.TenantID != f.TenantID {
		return false
	}

	if f.ActorID != "" &&
		event.ActorID != f.ActorID {
		return false
	}

	if f.Category != "" &&
		event.Category != f.Category {
		return false
	}

	if f.Action != "" &&
		event.Action != f.Action {
		return false
	}

	if f.Status != "" &&
		event.Status != f.Status {
		return false
	}

	if f.Severity != "" &&
		event.Severity != f.Severity {
		return false
	}

	if f.Resource != "" &&
		event.Resource != f.Resource {
		return false
	}

	if f.Tag != "" &&
		!containsTag(event.Tags, f.Tag) {
		return false
	}

	if f.From != nil &&
		event.Timestamp.Before(*f.From) {
		return false
	}

	if f.To != nil &&
		event.Timestamp.After(*f.To) {
		return false
	}

	return true
}

func containsTag(tags []string, target string) bool {

	for _, tag := range tags {

		if tag == target {
			return true
		}
	}

	return false
}

func (f *Filter) ByTenant(
	tenantID string,
) *Filter {

	f.TenantID = tenantID

	return f
}

func (f *Filter) ByActor(
	actorID string,
) *Filter {

	f.ActorID = actorID

	return f
}

func (f *Filter) ByCategory(
	category Category,
) *Filter {

	f.Category = category

	return f
}

func (f *Filter) BySeverity(
	severity Severity,
) *Filter {

	f.Severity = severity

	return f
}

func (f *Filter) ByStatus(
	status Status,
) *Filter {

	f.Status = status

	return f
}

func (f *Filter) ByResource(
	resource string,
) *Filter {

	f.Resource = resource

	return f
}

func (f *Filter) ByTag(
	tag string,
) *Filter {

	f.Tag = tag

	return f
}

func (f *Filter) Between(
	from time.Time,
	to time.Time,
) *Filter {

	f.From = &from
	f.To = &to

	return f
}

func (f Filter) IsEmpty() bool {

	return f.ID == "" &&
		f.TenantID == "" &&
		f.ActorID == "" &&
		f.Category == "" &&
		f.Action == "" &&
		f.Status == "" &&
		f.Severity == "" &&
		f.Resource == "" &&
		f.Tag == "" &&
		f.From == nil &&
		f.To == nil
}

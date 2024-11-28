package role

type Role string

const (
	Unspecified Role = "unspecified"
	Owner       Role = "owner"
	Manager     Role = "manager"
	Employee    Role = "employee"
)

func (r Role) String() string {
	return string(r)
}

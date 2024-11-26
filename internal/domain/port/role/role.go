package role

type Role string

const (
	Owner    Role = "owner"
	Manager  Role = "manager"
	Employee Role = "employee"
)

func (r Role) String() string {
	return string(r)
}

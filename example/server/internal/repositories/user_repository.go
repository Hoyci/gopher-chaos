package repositories

type UserRepository interface {
	Save(user *User) error
	GetByID(id string) (*User, error)
	UpdateByID(id, name string) (*User, error)
	DeleteByID(id string) error
}

type User struct {
	ID    string
	Name  string
	Email string
	Age   int32
}

package domain

type User struct {
	ID       int    `db:"id"`
	Username string `db:"username"`
	Password string `db:"password"` // захешированный пароль
}

type UserRepository interface {
	Create(u *User) error
	GetByUsername(username string) (*User, error)
	GetByID(id int) (*User, error)
	Update(u *User) error
}

type UserUsecase interface {
	Register(username, password string) (*User, error)
	Authenticate(username, password string) (*User, error)
	GetProfile(id int) (*User, error)
	UpdateProfile(id int, username string) (*User, error)
}

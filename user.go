package chat

type User struct {
	Uid   uint64
	Name  string
	authk string
	authv string
}

type UserManager struct {
	authIndex map[string]*User
	guid      uint64
}

func NewUserManager() *UserManager {
	manager := &UserManager{
		authIndex: make(map[string]*User),
		guid:      10000,
	}
	return manager
}

func (man *UserManager) Auth(authk string, authv string) *User {
	user, ok := man.authIndex[authk]
	if !ok {
		return nil
	}
	return user
}

func (man *UserManager) AuthOrNew(authk string, authv string, name string) *User {
	user := man.Auth(authk, authv)
	if user == nil {
		user = man.NewUser(authk, authv, name)
	}
	return user
}

func (man *UserManager) NewUser(authk string, authv string, name string) *User {
	user := &User{
		authk: authk,
		authv: authv,
		Uid:   man.newUid(),
		Name:  name,
	}
	man.authIndex[authk] = user
	return user
}

func (man *UserManager) newUid() uint64 {
	man.guid++
	return man.guid
}

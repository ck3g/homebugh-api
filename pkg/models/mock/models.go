package mock

type Models struct {
	Users UserModel
}

func NewModels(db *interface{}) interface{} {
	return Models{
		Users: UserModel{},
	}
}

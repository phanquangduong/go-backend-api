package repo

// type UserRepo struct{}

// func NewUserRepo() *UserRepo {
// 	return &UserRepo{}
// }

// func (ur *UserRepo) GetInfoUser() string {
// 	return "qduongsayhi"
// }

// Interface version

type IUserRepository interface {
	GetUserByEmail(email string) bool
}

type userRepository struct {
}

func NewUserRepository() IUserRepository {
	return &userRepository{}
}

// GetUserByEmail implements IUserRepository.
func (u *userRepository) GetUserByEmail(email string) bool {
	// SELECT * FROM user WHERE email = '??' ORDER BY email
	// row := global.Mdb.Table(TableNameGoCrmUser).Where("usr_email = ?", email).First(&model.GoCrmUser{}).RowsAffected
	// user, err := up.sqlc.GetUserByEmailSQLC(ctx, email)
	// if err != nil {
	// 	fmt.Printf("GetUserByEmail error: %v\n", err)
	// 	return false
	// }

	return false
}

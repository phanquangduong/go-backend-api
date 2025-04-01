package repo

import (
	"go/go-backend-api/global"
	"go/go-backend-api/internal/database"
)

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
	sqlc *database.Queries
}

func NewUserRepository() IUserRepository {
	return &userRepository{
		sqlc: database.New(global.Mdbc),
	}
}

// GetUserByEmail implements IUserRepository.
func (up *userRepository) GetUserByEmail(email string) bool {
	// SELECT * FROM user WHERE email = '??' ORDER BY email
	// row := global.Mdb.Table(TableNameGoCrmUser).Where("usr_email = ?", email).First(&model.GoCrmUser{}).RowsAffected
	// user, err := up.sqlc.GetUserByEmailSQLC(ctx, email)
	// if err != nil {
	// 	fmt.Printf("GetUserByEmail error: %v\n", err)
	// 	return false
	// }

	user, err := up.sqlc.GetUserByEmailSQLC(ctx, email)

	if err != nil {
		return false
	}

	return user.UsrID != 0
}

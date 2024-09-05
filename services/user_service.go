package services

import (
	"github.com/ShinnosukeSuzuki/techtrain-mission-ca-tech-dojo-golang/models"
	"github.com/ShinnosukeSuzuki/techtrain-mission-ca-tech-dojo-golang/pkg/uuid"
	"github.com/ShinnosukeSuzuki/techtrain-mission-ca-tech-dojo-golang/repositories"
)

// サービス構造体を定義
type UserService struct {
	// userRepositoryを埋め込む
	repository repositories.UserRepository
}

// サービスのコンストラクタ
func NewUserService(r repositories.UserRepository) *UserService {
	return &UserService{repository: r}
}

// ハンドラー UserCreateHandler 用のサービスメソッド
func (s *UserService) UserCreateService(name string) (models.User, error) {

	id := uuid.GenerateUUID()
	token := uuid.GenerateUUID()

	newUser, err := s.repository.CreateUser(id, name, token)
	if err != nil {
		return models.User{}, err
	}

	return newUser, nil
}

// ハンドラー UserGetHandler 用のサービスメソッド
func (s *UserService) UserGetService(userId string) (models.User, error) {

	user, err := s.repository.GetUserById(userId)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

// ハンドラー UserUpdateHandler 用のサービスメソッド
func (s *UserService) UserUpdateService(userId, name string) error {

	err := s.repository.UpdateUserName(userId, name)
	if err != nil {
		return err
	}

	return nil
}

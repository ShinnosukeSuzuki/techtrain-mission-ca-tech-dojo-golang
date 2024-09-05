package services

import (
	"github.com/ShinnosukeSuzuki/techtrain-mission-ca-tech-dojo-golang/models"
	"github.com/ShinnosukeSuzuki/techtrain-mission-ca-tech-dojo-golang/pkg/uuid"
	"github.com/ShinnosukeSuzuki/techtrain-mission-ca-tech-dojo-golang/repositories"
)

// サービス構造体を定義
type UserService struct {
	// userRepositoryを埋め込む
	uRep repositories.UserRepository
}

// サービスのコンストラクタ
func NewUserService(r repositories.UserRepository) *UserService {
	return &UserService{uRep: r}
}

// ハンドラー CreateHandler 用のサービスメソッド
func (s *UserService) Create(name string) (models.User, error) {

	id := uuid.GenerateUUID()
	token := uuid.GenerateUUID()

	newUser, err := s.uRep.Create(id, name, token)
	if err != nil {
		return models.User{}, err
	}

	return newUser, nil
}

// ハンドラー GetHandler 用のサービスメソッド
func (s *UserService) Get(userId string) (models.User, error) {

	user, err := s.uRep.GetById(userId)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

// ハンドラー UpdateNameHandler 用のサービスメソッド
func (s *UserService) UpdateName(userId, name string) error {

	err := s.uRep.UpdateName(userId, name)
	if err != nil {
		return err
	}

	return nil
}

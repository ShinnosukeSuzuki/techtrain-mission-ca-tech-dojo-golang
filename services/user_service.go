package services

import (
	"github.com/ShinnosukeSuzuki/techtrain-mission-ca-tech-dojo-golang/dto"
	"github.com/ShinnosukeSuzuki/techtrain-mission-ca-tech-dojo-golang/models"
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

	eu, err := s.uRep.Create(name)
	if err != nil {
		return models.User{}, err
	}

	newUser := dtoToModel(eu)

	return newUser, nil
}

// ハンドラー GetHandler 用のサービスメソッド
func (s *UserService) Get(userID string) (models.User, error) {

	eu, err := s.uRep.GetById(userID)
	if err != nil {
		return models.User{}, err
	}

	user := dtoToModel(eu)

	return user, nil
}

// ハンドラー UpdateNameHandler 用のサービスメソッド
func (s *UserService) UpdateName(userID, name string) error {

	err := s.uRep.UpdateName(userID, name)
	if err != nil {
		return err
	}

	return nil
}

// dto.User から models.User に変換する
func dtoToModel(user dto.User) models.User {
	return models.User{
		ID:    user.ID,
		Name:  user.Name,
		Token: user.Token,
	}
}

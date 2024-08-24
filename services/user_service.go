package services

import (
	"github.com/ShinnosukeSuzuki/techtrain-mission-ca-tech-dojo-golang/common"
	"github.com/ShinnosukeSuzuki/techtrain-mission-ca-tech-dojo-golang/models"
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

	// tokenを生成
	token := common.GenerateToken()

	// ユーザーを作成
	newUser, err := s.repository.CreateUser(name, token)
	if err != nil {
		return models.User{}, err
	}

	return newUser, nil
}

// ハンドラー UserGetHandler 用のサービスメソッド
func (s *UserService) UserGetService(token string) (models.User, error) {

	// tokenを持つユーザーを取得
	user, err := s.repository.GetUserByToken(token)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

// ハンドラー UserUpdateHandler 用のサービスメソッド
func (s *UserService) UserUpdateService(token string, name string) error {

	// tokenを持つユーザーのnameを更新
	err := s.repository.UpdateUserNameByToken(token, name)
	if err != nil {
		return err
	}

	return nil
}

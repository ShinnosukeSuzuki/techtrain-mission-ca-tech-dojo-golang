package services

import (
	"github.com/ShinnosukeSuzuki/techtrain-mission-ca-tech-dojo-golang/common"
	"github.com/ShinnosukeSuzuki/techtrain-mission-ca-tech-dojo-golang/models"
	"github.com/ShinnosukeSuzuki/techtrain-mission-ca-tech-dojo-golang/repositories"
)

// ハンドラー UserCreateHandler 用のサービスメソッド
func (s *MyAppService) UserCreateService(name string) (models.User, error) {

	// tokenを生成
	token := common.GenerateToken()

	// ユーザーを作成
	newUser, err := repositories.CreateUser(s.db, name, token)
	if err != nil {
		return models.User{}, err
	}

	return newUser, nil
}

// ハンドラー UserGetHandler 用のサービスメソッド
func (s *MyAppService) UserGetService(token string) (models.User, error) {

	// tokenを持つユーザーを取得
	user, err := repositories.GetUserByToken(s.db, token)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

// ハンドラー UserUpdateHandler 用のサービスメソッド
func (s *MyAppService) UserUpdateService(token string, name string) error {

	// tokenを持つユーザーのnameを更新
	err := repositories.UpdateUserNameByToken(s.db, token, name)
	if err != nil {
		return err
	}

	return nil
}

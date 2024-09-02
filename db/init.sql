# ユーザ情報を格納するためのテーブル
CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(255) NOT NULL UNIQUE PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    token VARCHAR(255) NOT NULL UNIQUE
);

# キャラクター情報を格納するためのテーブル
CREATE TABLE IF NOT EXISTS characters (
    id VARCHAR(255) NOT NULL UNIQUE PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    probability FLOAT NOT NULL
);

# ユーザがキャラクターを引いた記録を格納するためのテーブル
CREATE TABLE IF NOT EXISTS users_characters (
    id VARCHAR(255) NOT NULL UNIQUE PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL,
    character_id VARCHAR(255) NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (character_id) REFERENCES characters(id)
);

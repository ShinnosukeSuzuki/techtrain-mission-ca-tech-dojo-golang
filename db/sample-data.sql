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

# ユーザ情報のsampleデータ
INSERT INTO users (id, name, token) VALUES 
('a7abf496-3f37-488d-a74e-49b2b555c987', '進之介', '076a5b2c-97a9-49c7-b426-41f861386f8b'),
('ae2547e8-c00e-43b7-9dfa-2799f24e6455', '太郎', '0adede71-78ff-4296-a3ec-8074cd56ad30');

# キャラクター情報のsampleデータ
INSERT INTO characters (id, name, probability) VALUES 
('3c8e0b3c-9dc9-426c-879f-72bf69221a23', 'ルフィ', 0.1),
('e4b12136-de5e-4760-abae-efeaabd0338a', 'ゾロ', 0.2),
('5af3b613-4892-440d-8635-9ae4a13a7b4a', 'サンジ', 0.2),
('73e1cc66-b667-4fa1-9e6d-60ae3ced2ee7', 'ナミ', 0.4),
('6406c5aa-5467-4e9a-86f6-11dd06a2110e', 'ウソップ', 0.5),
('754744d2-4bfe-4ff0-85b4-8ec852419131', 'チョッパー', 0.6),
('7abd4440-b137-44d7-94b5-a48490a899f1', 'フランキー', 0.7),
('d19b30c4-8825-410f-8a0a-f3651a8186ff', 'ブルック', 0.8),
('93be9bce-4225-4398-b0a9-9aa2ea348bbe', 'ジンベエ', 0.2),
('72a57228-1ae5-4a73-8ca2-738d27574a35', 'ロビン', 0.4);

# ユーザがキャラクターを引いた記録のsampleデータ
INSERT INTO users_characters (id, user_id, character_id) VALUES 
('c4334312-d6cb-4be6-851a-1ed8423a6f15', 'a7abf496-3f37-488d-a74e-49b2b555c987', '3c8e0b3c-9dc9-426c-879f-72bf69221a23'),
('80548478-4b03-4ef0-8b4f-7f2873e83d9e', 'a7abf496-3f37-488d-a74e-49b2b555c987', 'e4b12136-de5e-4760-abae-efeaabd0338a'),
('dab77755-91cc-4ad2-8f2b-7dbd06d8bdaa', 'a7abf496-3f37-488d-a74e-49b2b555c987', '5af3b613-4892-440d-8635-9ae4a13a7b4a'),
('db6ab0c6-9d58-48ed-96bf-127771e1d1ed', 'a7abf496-3f37-488d-a74e-49b2b555c987', '73e1cc66-b667-4fa1-9e6d-60ae3ced2ee7'),
('e2162fad-a268-4f65-8218-aa756daf1739', 'ae2547e8-c00e-43b7-9dfa-2799f24e6455', '6406c5aa-5467-4e9a-86f6-11dd06a2110e');

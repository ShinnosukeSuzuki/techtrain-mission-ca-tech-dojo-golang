import os
import pandas as pd
import boto3
from sqlalchemy.sql import text
from sqlalchemy.engine import create_engine


class MySQLClient:
    def __init__(self, server: str, database: str, username: str, password: str):
        self.server = server
        self.database = database
        self.username = username
        self.password = password
        connection_string = f"mysql+mysqlconnector://{username}:{password}@{server}/{database}"
        self.engine = create_engine(connection_string)
        self.conn = self.engine.connect()

    def close(self):
        self.conn.close()

    # csvファイルを読み込み、charactersテーブルにアップサートする
    def update_characters(self, df: pd.DataFrame):
        transaction = self.conn.begin()
        
        # 1行ずつ処理
        for _, row in df.iterrows():
            id = row['id']
            name = row['name']
            probability = row['probability']
            
            # アップサートクエリを生成 (idをUUID_TO_BIN()で変換)
            upsert_query = """
            INSERT INTO characters (id, name, probability)
            VALUES (UUID_TO_BIN(:id), :name, :probability)
            ON DUPLICATE KEY UPDATE
                name = VALUES(name),
                probability = VALUES(probability);
            """
            
            try:
                # パラメータをバインド
                params = {
                    'id': id,
                    'name': name,
                    'probability': probability
                }
                # クエリを実行
                self.conn.execute(text(upsert_query), params)
            except Exception as e:
                # エラーが発生した行を表示
                print(f"Error at row {id}: {row}")
                print(f"Error message: {e}")
        
        # トランザクションをコミット
        transaction.commit()

# DB接続情報は環境変数から取得
server = os.environ['DB_SERVER']
database = os.environ['DB_DATABASE']
username = os.environ['DB_USERNAME']
password = os.environ['DB_PASSWORD']


def lambda_handler(event, context):
    print("event: ", event)
    # S3バケットとオブジェクトキーを取得
    # S3バケットとオブジェクトキーを取得
    bucket = event['detail']['bucket']['name']
    key = event['detail']['object']['key']

    # S3からCSVファイルを読み込む
    s3 = boto3.client('s3')
    response = s3.get_object(Bucket = bucket, Key = key)

    status = response.get("ResponseMetadata", {}).get("HTTPStatusCode")

    if status != 200:
        print(f"Error: failed to get object, status code: {status}")
        return {
            'statusCode': 500,
            'body': 'Failed to get object'
        }
    
    print(f"Successful S3 get_object response. Status - {status}")

    # PandasでCSVデータを読み込む
    df = pd.read_csv(response.get("Body"))

    # MySQLに接続
    client = MySQLClient(server, database, username, password)

    # charactersテーブルをアップサート
    client.update_characters(df)
    print("charactersテーブルをアップサートしました")

    # MySQL接続をクローズ
    client.close()

    return {
        'statusCode': 200,
        'body': 'Success'
    }

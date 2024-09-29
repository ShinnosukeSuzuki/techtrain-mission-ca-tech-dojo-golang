from locust import HttpUser, task, constant_throughput

class GachaDrawApiTest(HttpUser):
    wait_time = constant_throughput(1)

    @task
    # /gacha/draw ガチャ実行APIにリクエストを送信する
    def gacha_draw(self):
        headers = {
            "x-token": "076a5b2c-97a9-49c7-b426-41f861386f8b"
        }
        self.client.post("/gacha/draw", headers = headers, json={"times": 10})

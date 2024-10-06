from locust import HttpUser, task, constant_throughput

class GachaDrawApiTest(HttpUser):
    wait_time = constant_throughput(1)

    @task
    # /gacha/draw ガチャ実行APIにリクエストを送信する
    def gacha_draw(self):
        headers = {
            "x-token": "c05b11df-5592-440d-9615-1911006ed112"
        }
        self.client.post("/gacha/draw", headers = headers, json={"times": 10})

from locust import HttpUser, task, constant_throughput

class GachaDrawApiTest(HttpUser):
    wait_time = constant_throughput(1)

    @task
    # /gacha/draw ガチャ実行APIにリクエストを送信する
    def gacha_draw(self):
        headers = {
            "x-token": "0192ebd2-5afd-7c38-b51d-51b705ed52c2"
        }
        self.client.post("/gacha/draw", headers = headers, json={"times": 10})

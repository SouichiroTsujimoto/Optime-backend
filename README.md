### 環境変数の設定

`.env.example`をコピーして`.env`を作成し、適切な値を設定

```bash
cp .env.example .env
```

### ビルド & 実行

```bash
sudo docker build -t optime-backend .

sudo docker run -d \
   --name optime-backend-server \
   -p 8080:8080 \
   --env-file .env \
   --restart always \
   optime-backend
```
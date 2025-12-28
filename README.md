### Optime
https://github.com/SouichiroTsujimoto/Optime

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

### AWS構成図
現在ver1の構成になっています。

ver2移行タスク(優先度順)
- NAT gatewayをfck-natのec2インスタンスに置き換え
- bastionからSSMを用いたアクセスへの移行

<img width="1337" height="521" alt="Optime Backend AWS (2)" src="https://github.com/user-attachments/assets/faacd2ef-262d-42bf-acea-ea30d00cd166" />

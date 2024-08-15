# pelith assignment
## 本地執行
`docker-compose.yml` 中的 database 參數需自行修改
```
docker compose up --build
```
## sql init
使用 ./migrate/init.sql，第一次執行 docker compose 時會一起 init
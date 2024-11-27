### start colima
```shell
colima start --kubernetes
```

### Run redis on docker
```shell
docker run -d --name redis-server -p 6379:6379 redis
# verify its running
docker ps
# connect to redis
# docker exec -it redis-server redis-cli
```
#### Redis DB default details:
- Addr:     "localhost:6379"
- Password: ""
- DB:       0
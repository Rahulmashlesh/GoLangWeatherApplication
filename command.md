# Redis can be accessed via port 6379 on the following DNS name from within your cluster:
    my-redis-master.default.svc.cluster.local
# To get your password run:
    export REDIS_PASSWORD=$(kubectl get secret  weather-redis -o jsonpath="{.data.redis-password}" | base64 -d)

# To connect to your Redis&reg; server:

# Run a Redis pod that you can use as a client and attach to the client pod.
   kubectl run redis-client --restart='Never'  --env REDIS_PASSWORD=$REDIS_PASSWORD  --image docker.io/bitnami/redis:7.2.5-debian-12-r0 --command -- sleep infinity
   kubectl exec --tty -i redis-client -- bash

# Connect using the Redis CLI:
   REDISCLI_AUTH="$REDIS_PASSWORD" redis-cli -h weather-redis-master

# To connect to your database from outside the cluster execute the following commands:
    kubectl port-forward --namespace default svc/weather-redis-master 6379:6379 &
    REDISCLI_AUTH="$REDIS_PASSWORD" redis-cli -h 127.0.0.1 -p 6379



kubectl run redis-client --restart='Never'  --env REDIS_PASSWORD="redis-p0"  --image docker.io/bitnami/redis:7.2.5-debian-12-r0 --command -- sleep infinity
kubectl exec --tty -i redis-client -- bash
I have no name!@redis-client:/$ redis-cli -h weather-redis-master.goweather
weather-redis-master.goweather:6379> ping
(error) NOAUTH Authentication required.
weather-redis-master.goweather:6379> auth
(error) ERR wrong number of arguments for 'auth' command
weather-redis-master.goweather:6379> auth redis-p0
OK
weather-redis-master.goweather:6379> ping
PONG
weather-redis-master.goweather:6379> 

# run redis on docker
$ docker run -d --name redis-stack -p 6379:6379 -p 8001:8001 redis/redis-stack:latest
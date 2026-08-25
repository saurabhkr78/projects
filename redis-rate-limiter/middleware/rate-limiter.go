package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type Limiter struct {
	RedisClient *redis.Client
}

func NewLimiter(redisClient *redis.Client) *Limiter {
	return &Limiter{
		RedisClient: redisClient,
	}
}
func (r *Limiter) RateLimit(incomingHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		//get the user context from the request context because redis method expects a context
		ctx := req.Context()
		//get the ip address of the user from the request
		ip := req.RemoteAddr
		//construct the key for redis as "rate-limit:<ip-address>"
		key := "rate-limit:" + ip

		//#####RACE_CONDITION#####
		//this creating a race condition because if two requests come at the same time then both will get the same value of the key and both will increment the value of the key by 1 and both will call the incoming handler to process the request
		//if the value of the key is less than 5 then increment the value of the key by 1
		//instead of get then incer
		//we do incr then get because incr is atomic operation and it will return the new value of the key after incrementing it by 1

		//now if two requests come at the same time then both will increment the value of the key by 1 and both will get the new value of the key after incrementing it by 1 and both will call the incoming handler to process the request
		//due to race condition created by incr method of redis after get method
		// we cannot call the client first to get the value of the key s
		//so first incr the value of the key by 1 and then get the value of the key from redis refer:notes.md

		err := r.RedisClient.Incr(ctx, key).Err()
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		//call the redis client to get the value of the key from redis
		//if the key does not exist in redis then it will return redis.Nil error
		count, err := r.RedisClient.Get(ctx, key).Result()
		if err == redis.Nil {
			//if the key does not exist in redis then set the value of the key to 1 and set the expiration time to 1 minute
			err = r.RedisClient.Set(ctx, key, 1, 1*time.Minute).Err()
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			//call the incoming handler to process the request
			incomingHandler.ServeHTTP(w, req)
			return
		} else if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		//if the key exists in redis then check if the value of the key is greater than 5
		if cnt, _ := strconv.Atoi(count); cnt >= 5 {
			http.Error(w, "Too Many Requests  slow down", http.StatusTooManyRequests)
			return
		}

		//call the incoming handler to process the request
		incomingHandler.ServeHTTP(w, req)

	})
}

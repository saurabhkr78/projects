package handler

import (
	"encoding/json"
	"net/http"
	"redis-otp-service/client"
	"redis-otp-service/dto"
	. "redis-otp-service/util"
	"time"

	"github.com/redis/go-redis/v9"
)

// inject the redis client into handler struct
type Handler struct {
	//redis client instance
	Client *client.Client
}

func NewHandler(Client *client.Client) *Handler {
	return &Handler{
		Client: Client,
	}
}

func (h *Handler) Requestotp(w http.ResponseWriter, r *http.Request) {
	//get the context from the request
	ctx := r.Context()
	//read request body and unmarshal it into RequestOTP struct
	var reqOTP dto.RequestOTP
	//get the new encoder to decode the request body into the struct
	err := json.NewDecoder(r.Body).Decode(&reqOTP)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	//now generate the otp and store it in redis with the user id as key and otp as value
	otp := GenerateOTP()
	//create the key for redis as opt for user by concatenating the user id of request body
	key := "otp:user:" + reqOTP.UserID

	//store the otp in redis with the key and value as otp and set the expiration time to 5 minutes
	//since redis returns a redis command object so we
	_, err = h.Client.RedisClient.Set(ctx, key, otp, 5*time.Minute).Result()
	if err != nil {
		http.Error(w, "Failed to store OTP in redis", http.StatusInternalServerError)
		return
	}

	//send the response back to the client with the otp and user id in json format
	//but before that crate a struct to hold the response data
	response := map[string]string{
		"message": "OTP generated successfully",
		"user_id": reqOTP.UserID,
		"otp":     otp,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}
func (h *Handler) Verifyotp(w http.ResponseWriter, r *http.Request) {
	//first get the context from the request
	ctx := r.Context()

	//create a varibale that can hold the request body data
	var reqVerifyOTP dto.VerifyOTPRequest

	//Note:what we want to receive in request body we should define a variable of that type and then decode the request body into that variable
	//then decode the request body into the variable
	err := json.NewDecoder(r.Body).Decode(&reqVerifyOTP)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	//build the key to get the otp from redis using the user id and otp sent by the client
	//exact same key what was used to store the otp in redis when the otp was generated
	key := "otp:user:" + reqVerifyOTP.UserID

	//then call the redis client to get the otp from redis and compare it with the otp sent by the client
	storedOTP, err := h.Client.RedisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		http.Error(w, "OTP expired or not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Failed to get OTP from redis", http.StatusInternalServerError)
		return
	}

	//compare the otp sent by the client with the otp stored in redis
	if storedOTP == reqVerifyOTP.OTP {
		response := map[string]string{
			"message": "OTP verified successfully",
			"user_id": reqVerifyOTP.UserID,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

}

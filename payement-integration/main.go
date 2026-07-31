package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/checkout/session"
)

// in order to make a payment request, we need to define a struct that will hold the necessary information for the payment. This struct will be used to serialize the payment request data into JSON format when sending it to the payment gateway API
type PaymentRequest struct {
	UserUUID string  `json:"user_uuid,omitempty"`
	Amount   float64 `json:"amount,omitempty"`
	PlanName string  `json:"plan_name,omitempty"`
}

// now we need to define a struct that will hold the response data from the payment gateway API. This struct will be used to deserialize the JSON response data into a Go struct when receiving the response from the payment gateway API
type APIResponse struct {
	Status     string      `json:"status,omitempty"`
	StatusCode int         `json:"status_code,omitempty"`
	Message    string      `json:"message,omitempty"`
	SessionURL StripeModel `json:"session_url,omitempty"`
	//why the session url is a struct? because the session url is a nested object in the response data from the payment gateway API. The session url contains the url that the user will be redirected to in order to complete the payment process. The session url is a nested object because it contains additional information about the payment session, such as the session id and the payment method used. By defining the session url as a struct, we can easily access this information when deserializing the JSON response data into a Go struct.
}

// now we need to define a struct that will hold the session url data from the payment gateway API. This struct will be used to deserialize the JSON response data into a Go struct when receiving the response from the payment gateway API
type StripeModel struct {
	UserUUID  string `json:"user_uuid,omitempty"`
	URL       string `json:"url,omitempty"`
	PlanName  string `json:"plan_name,omitempty"`  //jo hum kharid rahe hain uska plan name
	Price     string `json:"price,omitempty"`      //jo hum kharid rahe hain uska price
	CreatedAt string `json:"created_at,omitempty"` //jo hum kharid rahe hain uska created at
}

//the stripe doesnot take the amount in USD so we need to be in cents
//so we need a convert function

func ConvertToCents(price string, currency string) (int64, error) {
	var amount float64
	_, err := fmt.Sscanf(price, "%f", &amount)
	if err != nil {
		return 0, fmt.Errorf("Invalid Price format")
	}

	return int64(amount * 100), nil
}

// integration logic
func StripeSession(ctx context.Context, userUUID, amount, planName string) (*APIResponse, error) {
	ConvertedAmount, err := ConvertToCents(amount, "USD")
	if err != nil {
		return nil, fmt.Errorf("Failed to convert amount to cents: %v", err)
	}
	fmt.Println("converted amount--->", ConvertedAmount)

	//now we need stripe key
	stripe.Key = "sk_test_51SLLfQGkl90T1eK8no1VfMFi0VjmjIAFVfZJCuiIm7HN1cgA3gDixGImhPvsVD5UlEtNROvNlOmfjjWJYIKIgewR00U8ApsGbj"

	// ab hum stripe session create karenge
	//Session mein jo required information hai like price,id,etc
	// hum iske liye stripe ke checkout session params ka use karenge
	//stripe ke checkout session params mein hum user ka uuid, amount, plan name, currency, payment method, success url, cancel url, etc. define karenge
	params := &stripe.CheckoutSessionParams{
		Mode: stripe.String(string(stripe.CheckoutSessionModePayment)),
		LineItems: []*stripe.CheckoutSessionLineItemParams{{
			PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
				Currency: stripe.String("usd"),
				ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
					Name: stripe.String(planName),
				},
				UnitAmount: stripe.Int64(ConvertedAmount),
			},
			Quantity: stripe.Int64(1),
		}},
		AllowPromotionCodes: stripe.Bool(true),

		//success or cancel url
		SuccessURL: stripe.String("http://localhost:3000/payment_success.html"),
		CancelURL:  stripe.String("http://localhost:3000/Payment_failed.html"),

		//agar payment itne der ke baad nahi hoti hai to session expire ho jaye by default 30 minutes
		ExpiresAt: stripe.Int64(time.Now().Add(30 * 60).Unix()), // 30 minutes

		PaymentIntentData: &stripe.CheckoutSessionPaymentIntentDataParams{
			SetupFutureUsage: stripe.String("off_session"),
		},
	}
	//add metadata to the session
	params.AddMetadata("api_version", "2026-06-01")

	//now we need to create a stripe session
	sess, err := session.New(params)
	if err != nil {
		return nil, fmt.Errorf("Failed to create stripe session: %v", err)
	}
	//now we store the session data in a model
	StripeRecord := StripeModel{
		UserUUID:  userUUID,
		URL:       sess.URL,
		PlanName:  planName,
		Price:     amount,
		CreatedAt: time.Now().UTC().Format(time.RFC3339), //Everyone must send dates exactly like this YYYY-MM-DDTHH:MM:SSZa nd unix timestamp  Stores seconds since January 1, 1970 UTC. and not human readable format but fast and easy to compare and sort
	}
	fmt.Println("stripe model--->", StripeRecord)

	//now we need to return the response to the client
	response := &APIResponse{
		Status:     "success",
		StatusCode: 200,
		Message:    "checkout session created successfully",
		SessionURL: StripeModel{
			SessionURL: sess.URL,
			UserUUID:   userUUID,
			PlanName:   planName,
			Price:      amount,
			CreatedAt:  StripeModel.CreatedAt,
		},
	}
	fmt.Printf("checkout session URL---> %s\n", sess.URL)
	return response, nil
}

//now our code is ready to create a stripe session and return the session url to the client. The client can then redirect the user to the session url to complete the payment process.
//now we convert it into api

func main() {
	http.HandleFunc("/create-checkout-session", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only Post Method is allowed", http.StatusMethodNotAllowed)
			return
		}
		var req PaymentRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, fmt.Sprintf("Invalid request payload: %v", err), http.StatusBadRequest)
			return
		}
		ctx := r.Context()
		response, err := StripeSession(ctx, req.UserUUID, req.Amount, req.PlanName)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to create checkout session: %v", err), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})
	fmt.Println("server is running on ")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

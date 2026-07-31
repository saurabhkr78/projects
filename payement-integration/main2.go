package main


import ()

//customer walk into store and says i want a buy premium plan of netflix and he says i want to pay 10$ for it
//he fills the form that is struct
type PaymentRequest struct {
	UserUUID string  `json:"user_uuid,omitempty"`
	Amount   float64 `json:"amount,omitempty"`
	PlanName string  `json:"plan_name,omitempty"`
}
//customer says i want to pay 10$ for it but stripe doesnot take the amount in USD so we need to be in cents
//so we need a convert function
func ConvertToCents(price string, currency string) (int64, error) {
	var amount float64
	_, err := fmt.Sscanf(price, "%f", &amount)
	if err != nil {
		return 0, fmt.Errorf("Invalid Price format")
	}

	return int64(amount * 100), nil
}

//now customer go to the cashier now the cashier starts payment process he starts a session for this customer

//so this payment session requires some information to start with like userUUID,amount,planName so we need to pass these information to the cashier
func StripeSession(ctx context.Context, userUUID, amount, planName string) (*APIResponse, error) {
	convertedAmount, err := ConvertToCents(amount, "USD")//converting the amount to cents
	if err != nil {
		return nil, fmt.Errorf("Failed to convert amount to cents: %v", err)
	}
	fmt.Println("converted amount--->", convertedAmount)

// The shopping cart is ready.
//
// Now the cashier has to fill a long billing form before
// sending the customer to Stripe.
//
// Think of CheckoutSessionParams as a government billing form.
//
// Stripe asks questions like:
//
// Customer?
// What are they buying?
// How much?
// Which currency?
// Can coupons be applied?
// Where should the customer go after payment?
// Where should they go if payment fails?
// Should we remember this card?
// When should this payment session expire?
//
// Once every field is filled, the cashier sends the form
// to Stripe's billing counter to create a Checkout Session.
params := &stripe.CheckoutSessionParams{

	// ----------------------------
	// Shopping Cart
	// ----------------------------
	//
	// Every checkout begins with a shopping cart.
	// One customer can buy one or many products.
	LineItems: []*stripe.CheckoutSessionLineItemParams{
		{
			// Every product needs a price tag.
			PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{

				// Cashier asks:
				// "Which currency are we charging?"
				Currency: stripe.String(string(stripe.CurrencyUSD)),

				// Describe the product.
				ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{

					// Customer:
					// "I want Premium Plan."
					Name: stripe.String(planName),
				},

				// Price of ONE product.
				// Stripe accepts the smallest currency unit.
				//
				// $20 → 2000 cents
				UnitAmount: stripe.Int64(convertedAmount),
			},

			// Customer is buying one copy.
			Quantity: stripe.Int64(1),
		},
	},

	// ------------------------------------------------
	// Promotion Codes (Coupon Counter)
	// ------------------------------------------------
	//
	// Before taking payment, the cashier asks:
	//
	// "Do you have any discount coupon?"
	//
	// If yes,
	// Stripe checks whether the coupon is valid
	// and automatically reduces the final amount.
	AllowPromotionCodes: stripe.Bool(true),

	// ------------------------------------------------
	// Success URL
	// ------------------------------------------------
	//
	// Imagine the customer pays successfully.
	//
	// Stripe asks:
	//
	// "Where should I send the customer now?"
	//
	// Cashier replies:
	//
	// "Send them back to our Thank You page."
	SuccessURL: stripe.String(
		"http://localhost:3000/payment_success.html",
	),

	// ------------------------------------------------
	// Cancel URL
	// ------------------------------------------------
	//
	// Suppose the customer closes the payment page,
	// payment fails,
	// or presses Cancel.
	//
	// Stripe again asks:
	//
	// "Where should I send them?"
	//
	// Cashier:
	//
	// "Take them back to our Payment Failed page."
	CancelURL: stripe.String(
		"http://localhost:3000/payment_failed.html",
	),

	// ------------------------------------------------
	// Session Expiry
	// ------------------------------------------------
	//
	// Imagine the customer opens the payment page
	// and disappears for six hours.
	//
	// We don't want old payment links
	// staying alive forever.
	//
	// Cashier tells Stripe:
	//
	// "If payment isn't completed within
	// 30 minutes,
	// throw away this billing form."
	ExpiresAt: stripe.Int64(
		time.Now().Add(30 * time.Minute).Unix(),
	),

	// ------------------------------------------------
	// Payment Intent
	// ------------------------------------------------
	//
	// Think of PaymentIntent as Stripe opening
	// a payment file for this customer.
	//
	// Inside that payment file we can tell Stripe
	// what to do after this payment.
	//
	// In this case the cashier asks:
	//
	// "If the customer agrees,
	// please remember this card
	// so they don't have to enter it again
	// next time."
	//
	// This is useful for subscriptions,
	// one-click checkout,
	// or future automatic payments.
	PaymentIntentData: &stripe.CheckoutSessionPaymentIntentDataParams{

		// Save the customer's payment method
		// for future off-session payments.
		//
		// "Off-session" means the customer
		// doesn't need to be sitting in front
		// of the computer next time.
		//
		// Example:
		// Netflix monthly renewal.
		SetupFutureUsage: stripe.String("off_session"),
	},
}
// When Stripe later sends a webhook, these notes come back, making it easy to identify which order was paid.
params.AddMetadata("order_id", orderID)
params.AddMetadata("user_id", userUUID)
params.AddMetadata("plan", planName)


//"Everything is ready. Let me send your cart to the billing counter (Stripe)."
	sess,er:=session.New(params)
	if er!=nil{
		return nil, fmt.Errorf("Failed to create stripe session: %v", er)
	}

	// Now the cashier has a receipt from Stripe.
	// The receipt contains a URL that the cashier can give to the customer.
	// The customer can then go to that URL and complete the payment.
	// The cashier also stores/write the receipt in their own records. so tommorrow if the customer comes back and says "I want to see my receipt or i paid yesterday", the cashier can show it to them.
	StripeRecord := StripeModel{
		UserUUID:  userUUID,
		URL:       sess.URL,
		PlanName:  planName,
		Price:     amount,
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}

	//since stripe returns plan,url,price,url,cretedAt in nested object so we need to create a struct for it
type StripeModel struct {
	UserUUID  string `json:"user_uuid,omitempty"`
	URL       string `json:"url,omitempty"`
	PlanName  string `json:"plan_name,omitempty"`  //jo hum kharid rahe hain uska plan name
	Price     string `json:"price,omitempty"`      //jo hum kharid rahe hain uska price
	CreatedAt string `json:"created_at,omitempty"` //jo hum kharid rahe hain uska created at
}
//cahsier prepare the payment receipt format
type APIResponse struct {
	Status     string      `json:"status,omitempty"`
	StatusCode int         `json:"status_code,omitempty"`
	Message    string      `json:"message,omitempty"`
	SessionURL StripeModel `json:"session_url,omitempty"`
}

	//return receipt to the customer
	response := &APIResponse{
		Status:     "success",
		StatusCode: 200,
		Message:    "checkout session created successfully",
		SessionURL: StripeModel{
			SessionURL: sess.URL,
			UserUUID:   userUUID,
			PlanName:   planName,
			Price:      amount,
			CreatedAt:  StripeRecord.CreatedAt,
		},
	}
	return response, nil
}














int main(){

}
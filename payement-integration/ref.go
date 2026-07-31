
/*
   CheckoutSessionParams
   │
   ├── Shopping Cart
   │     ├── Product
   │     ├── Price
   │     ├── Currency
   │     └── Quantity
   │
   ├── Customer
   │     ├── Customer ID
   │     ├── Email
   │     └── Billing Address
   │
   ├── Discounts
   │     ├── Promotion Codes
   │     └── Coupons
   │
   ├── Payment
   │     ├── Payment Methods
   │     ├── Save Card
   │     └── Payment Intent
   │
   ├── Navigation
   │     ├── Success URL
   │     ├── Cancel URL
   │     └── Locale
   │
   ├── Security
   │     ├── Expiry
   │     ├── Metadata
   │     └── Automatic Tax
   │
   └── Business

   	├── Invoice
   	├── Tax ID
   	├── Shipping
   	└── Custom Fields
*/
convertedAmount, err := ConvertToCents(amount, "USD") //converting the amount to cents
if err != nil {
	return nil, fmt.Errorf("Failed to convert amount to cents: %v", err)
}
fmt.Println("converted amount--->", convertedAmount)

//now we need stripe key before creating the checkout session
stripe.Key = os.Getenv("STRIPE_SECRET_KEY")

params := &stripe.CheckoutSessionParams{

	// ==========================================================
	// SHOPPING CART
	// ==========================================================

	Mode: stripe.String(string(stripe.CheckoutSessionModePayment)),

	LineItems: []*stripe.CheckoutSessionLineItemParams{
		{
			PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{

				Currency: stripe.String(string(stripe.CurrencyUSD)),

				ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
					Name: stripe.String(planName),
				},

				// Price in cents
				UnitAmount: stripe.Int64(convertedAmount),
			},

			Quantity: stripe.Int64(1),
		},
	},

	// ==========================================================
	// CUSTOMER
	// ==========================================================

	// Existing Stripe Customer
	Customer: stripe.String(customerID),

	// Used if Customer is not supplied
	CustomerEmail: stripe.String(email),

	// Ask for billing address
	BillingAddressCollection: stripe.String(
		string(stripe.CheckoutSessionBillingAddressCollectionRequired),
	),

	// Ask for phone number
	PhoneNumberCollection: &stripe.CheckoutSessionPhoneNumberCollectionParams{
		Enabled: stripe.Bool(true),
	},

	// ==========================================================
	// DISCOUNTS
	// ==========================================================

	// Customer can enter promo code
	AllowPromotionCodes: stripe.Bool(true),

	// Apply a coupon automatically
	Discounts: []*stripe.CheckoutSessionDiscountParams{
		{
			Coupon: stripe.String(couponID),
		},
	},

	// ==========================================================
	// PAYMENT
	// ==========================================================

	PaymentMethodTypes: []*string{
		stripe.String("card"),
		// stripe.String("link"),
		// stripe.String("us_bank_account"),
	},

	PaymentIntentData: &stripe.CheckoutSessionPaymentIntentDataParams{

		// Save card for future payments
		SetupFutureUsage: stripe.String("off_session"),

		Metadata: map[string]string{
			"order_id": orderID,
		},
	},

	// ==========================================================
	// NAVIGATION
	// ==========================================================

	SuccessURL: stripe.String(
		"https://example.com/payment-success",
	),

	CancelURL: stripe.String(
		"https://example.com/payment-failed",
	),

	Locale: stripe.String("auto"),

	// ==========================================================
	// SECURITY
	// ==========================================================

	ExpiresAt: stripe.Int64(
		time.Now().Add(30 * time.Minute).Unix(),
	),

	AutomaticTax: &stripe.CheckoutSessionAutomaticTaxParams{
		Enabled: stripe.Bool(true),
	},

	// ==========================================================
	// BUSINESS
	// ==========================================================

	InvoiceCreation: &stripe.CheckoutSessionInvoiceCreationParams{
		Enabled: stripe.Bool(true),
	},

	TaxIDCollection: &stripe.CheckoutSessionTaxIDCollectionParams{
		Enabled: stripe.Bool(true),
	},

	ShippingAddressCollection: &stripe.CheckoutSessionShippingAddressCollectionParams{
		AllowedCountries: []*string{
			stripe.String("US"),
			stripe.String("CA"),
			stripe.String("IN"),
		},
	},

	CustomFields: []*stripe.CheckoutSessionCustomFieldParams{

		{
			Key: stripe.String("company_name"),

			Label: &stripe.CheckoutSessionCustomFieldLabelParams{
				Type:   stripe.String("custom"),
				Custom: stripe.String("Company Name"),
			},

			Type: stripe.String("text"),
		},

		{
			Key: stripe.String("employee_id"),

			Label: &stripe.CheckoutSessionCustomFieldLabelParams{
				Type:   stripe.String("custom"),
				Custom: stripe.String("Employee ID"),
			},

			Type: stripe.String("text"),
		},
	},
}

// ==========================================================
// METADATA
// ==========================================================

params.AddMetadata("user_id", userUUID)
params.AddMetadata("order_id", orderID)
params.AddMetadata("plan_name", planName)
params.AddMetadata("api_version", "v1")
params.AddMetadata("environment", "production")
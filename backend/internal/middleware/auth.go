package middleware

// Struttura completa per la risposta di Clerk
type ClerkUser struct {
	CreatedAt           int64                  `json:"created_at"`
	UpdatedAt           int64                  `json:"updated_at"`
	LastSignInAt        *int64                 `json:"last_sign_in_at"`
	PrimaryEmailAddress *ClerkEmailAddress     `json:"primary_email_address"`
	PrimaryPhoneNumber  *ClerkPhoneNumber      `json:"primary_phone_number"`
	Username            *string                `json:"username"`
	FirstName           *string                `json:"first_name"`
	LastName            *string                `json:"last_name"`
	EmailAddresses      []ClerkEmailAddress    `json:"email_addresses"`
	PhoneNumbers        []ClerkPhoneNumber     `json:"phone_numbers"`
	PublicMetadata      map[string]interface{} `json:"public_metadata"`
	PrivateMetadata     map[string]interface{} `json:"private_metadata"`
	UnsafeMetadata      map[string]interface{} `json:"unsafe_metadata"`
	ID                  string                 `json:"id"`
	ImageURL            string                 `json:"image_url"`
}

type ClerkEmailAddress struct {
	Verification *ClerkVerification `json:"verification"`
	ID           string             `json:"id"`
	EmailAddress string             `json:"email_address"`
}

type ClerkPhoneNumber struct {
	Verification *ClerkVerification `json:"verification"`
	ID           string             `json:"id"`
	PhoneNumber  string             `json:"phone_number"`
}

type ClerkVerification struct {
	Status   string `json:"status"`
	Strategy string `json:"strategy"`
}

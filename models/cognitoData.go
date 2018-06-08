package models

type CognitoData struct {
	Sub                 string `json:"sub" bson:"sub"`
	EmailVerified       bool   `json:"email_verified" bson:"email_verified"`
	Iss                 string `json:"iss" bson:"iss"`
	PhoneNumberVerified bool   `json:"phone_number_verified" bson:"phone_number_verified"`
	CognitoUsername     string `json:"cognito:username bson:"cognito:username"`
	Aud                 string `json:"aud" bson:"aud"`
	EventID             string `json:"event_id" bson:"event_id"`
	TokenUse            string `json:"token_use" bson:"token_use"`
	AuthTime            int    `json:"auth_time" bson:"auth_time"`
	PhoneNumber         string `json:"phone_number" bson:"phone_number"`
	Exp                 int    `json:"exp" bson:"exp"`
	Iat                 int    `json:"iat" bson:"iat"`
	Email               string `json:"email" bson:"email"`
	JwtToken            string `json:"jwtToken" bson:"jwtToken"`
}


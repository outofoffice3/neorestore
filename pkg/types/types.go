package types

type RegisterRequest struct {
	S3BucketName string `json:"s3BucketName"`
	Prefix       string `json:"prefix"`
	Region       string `json:"region"`
}

type RegisterResponse struct {
	Body string `json:"body"`
}

type ValidateRequest struct {
	S3BucketName string `json:"s3BucketName"`
	Prefix       string `json:"prefix"`
	Region       string `json:"region"`
}

type ValidateResponse struct {
	Body string `json:"body"`
}

type OnboardRequest struct {
	S3BucketName string `json:"s3BucketName"`
	Prefix       string `json:"prefix"`
	Region       string `json:"region"`
}

type OnboardResponse struct {
	Body string `json:"body"`
}

type FinalizeRequest struct {
	S3BucketName string `json:"s3BucketName"`
	Prefix       string `json:"prefix"`
	Region       string `json:"region"`
}

type FinalizeResponse struct {
	Body string `json:"body"`
}

type PrefixItem struct {
	Keys           PrefixItemKey `json:"keys"`
	DateRegistered string        `json:"dateRegistered"`
	DateRestored   string        `json:"dateRestored"`
	DeleteMarkerId string        `json:"deleteMarkerId"`
}

type PrefixItemKey struct {
	PK string `json:"pk"`
	SK string `json:"sk"`
}

type PrefixListItem struct {
	Keys     PrefixItemKey `json:"keys"`
	Prefixes []string      `json:"prefixes"`
}

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

type RestoreRequest struct {
	S3BucketName string `json:"s3BucketName"`
	Prefix       string `json:"prefix"`
	Region       string `json:"region"`
}

type ListenerEvent struct {
	S3ObjectMetadata S3ObjectMetadata `json:"s3ObjectMetadata"`
	EventMetadata    EventMetadata    `json:"eventMetadata"`
}

type ListenerResponse struct {
	Body string `json:"body"`
}

type S3ObjectMetadata struct {
	BucketName string `json:"bucketName"`
	Key        string `json:"key"`
	VersionId  string `json:"versionId"`
	Region     string `json:"region"`
}

type EventMetadata struct {
	Name    string `json:"name"`
	Time    string `json:"time"`
	Version string `json:"version"`
	Region  string `json:"region"`
}

type DeleteMarker struct {
	VersionId string `json:"versionId"`
}

type Prefix struct {
	Prefix string `json:"prefix"`
}

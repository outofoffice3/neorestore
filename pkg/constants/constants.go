package constants

const (
	NeoDDBSecretName    = "neoRestorer/ddb"
	ReservedPK          = "chj-aws-poc-neo"
	ReservedSK          = "neoRestorer"
	ManifestTableName   = "neoRestorer"
	ManifestPK          = "PK"
	ManifestSK          = "SK"
	ListenerFunctionArn = "arn:aws:lambda:us-east-1:017608207428:function:delete-listener-neo"
	DateRegisteredAtt   = "dateRegistered"
	DateRestoredAtt     = "dateRestored"
	DeleteMarkerIdAtt   = "deleteMarkerId"
	PrefixListAtt       = "prefixes"
)

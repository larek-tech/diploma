package s3

type Credentials struct {
	endpoint        string
	accessKeyID     string
	secretAccessKey string
	useSSL          bool
}

func NewCredentials(endpoint, accessKey, secretAccessKey string, useSSL ...bool) Credentials {
	ssl := false
	if len(useSSL) > 0 {
		ssl = true
	}
	return Credentials{
		endpoint:        endpoint,
		accessKeyID:     accessKey,
		secretAccessKey: secretAccessKey,
		useSSL:          ssl,
	}
}

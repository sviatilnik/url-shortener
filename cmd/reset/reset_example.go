package main

// generate:reset
type User struct {
	ID       int64
	Name     string
	Email    string
	Active   bool
	Tags     []string
	Settings map[string]string
	Profile  *Profile
}

// generate:reset
type Profile struct {
	Bio         string
	Avatar      string
	SocialLinks map[string]string
	Metadata    []byte
}

// generate:reset
type Cache struct {
	Data    map[string]interface{}
	Keys    []string
	Size    int
	MaxSize int
	Enabled bool
	Backend *Backend
}

// generate:reset
type Backend struct {
	URL     string
	Timeout int
	Retries int
	Headers map[string]string
	Options []string
}

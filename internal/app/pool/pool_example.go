package pool

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

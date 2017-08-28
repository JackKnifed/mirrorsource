package mirrorsource

type Target interface {
	Check(chan<- error)
	Latest() string
	List() []string
}

type VerifyAction interface {
	Verify(string) bool
}

type PostAction interface {
	Process() error
}

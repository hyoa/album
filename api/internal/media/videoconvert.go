package media

type VideoConverter interface {
	Convert(key string) error
}

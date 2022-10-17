package mock

type VideoConverter struct{}

func (*VideoConverter) Convert(key string) error {
	return nil
}

package mock

type Storage struct {
	Keys []string
}

func (s *Storage) MediaExist(key, bucket string) (bool, error) {
	for _, k := range s.Keys {
		if k == key {
			return true, nil
		}
	}

	return false, nil
}

func (s *Storage) SignUploadUri(key, bucket string) (string, error) {
	return "signeduri", nil
}

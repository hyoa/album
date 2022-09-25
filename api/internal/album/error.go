package album

type AlbumNotFoundError struct {
	message string
}

func (e *AlbumNotFoundError) Error() string {
	return e.message
}

type AlbumSaveError struct {
	message string
}

func (e *AlbumSaveError) Error() string {
	return e.message
}

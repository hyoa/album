package media

import (
	"errors"
	"fmt"
)

type mediaKind string

const (
	KindPhoto mediaKind = "photo"
	KindVideo mediaKind = "video"
)

type Media struct {
	Key    string
	Author string
	Kind   mediaKind
	Folder string
}

type MediaManager struct {
	mediaRepo MediaRepository
	storage   Storage
}

func CreateMediaManager(mediaRepo MediaRepository, storage Storage) MediaManager {
	return MediaManager{
		mediaRepo: mediaRepo,
		storage:   storage,
	}
}

func (mm *MediaManager) GetMediasByFolder(folder string) ([]Media, error) {
	return mm.mediaRepo.FindByFolder(folder)
}

func (mm *MediaManager) GetFoldersNames(name string) ([]string, error) {
	return mm.mediaRepo.FindFoldersName(name)
}

func (mm *MediaManager) DeleteFolder(name string) error {
	medias, errFind := mm.mediaRepo.FindByFolder(name)

	if errFind != nil {
		return fmt.Errorf("Unable to get medias by folder: %w", errFind)
	}

	for _, m := range medias {
		m.Folder = "none"
		mm.mediaRepo.Save(m)
	}

	return nil
}

func (mm *MediaManager) GetAll() ([]Media, error) {
	return mm.mediaRepo.FindAll()
}

func (mm *MediaManager) ChangeMediasFolder(keys []string, newFolder string) ([]Media, error) {
	medias, errFind := mm.mediaRepo.FindManyByKeys(keys)

	if errFind != nil {
		return make([]Media, 0), fmt.Errorf("Unable to get medias by keys: %w", errFind)
	}

	for k := range medias {
		medias[k].Folder = newFolder
		mm.mediaRepo.Save(medias[k])
	}

	return medias, nil
}

func (mm *MediaManager) ChangeFolderName(folderToRename, newFolder string) ([]Media, error) {
	medias, errFind := mm.mediaRepo.FindByFolder(folderToRename)

	if errFind != nil {
		return make([]Media, 0), fmt.Errorf("Unable to get medias by folder: %w", errFind)
	}

	for k := range medias {
		medias[k].Folder = newFolder
		mm.mediaRepo.Save(medias[k])
	}

	return medias, nil
}

func (mm *MediaManager) Ingest(key, author, folder string, kind mediaKind) (Media, error) {
	mediaWithKey, errFind := mm.mediaRepo.FindByKey(key)

	if errFind != nil {
		return Media{}, fmt.Errorf("Unable to fetch media with key %s: %w", key, errFind)
	}

	if mediaWithKey != (Media{}) {
		return Media{}, nil
	}

	mediaInStorage, errCheckStorage := mm.storage.MediaExist(key)

	if errCheckStorage != nil {
		return Media{}, fmt.Errorf("Unable to check if media %s exist in storage: %w", key, errFind)
	}

	if !mediaInStorage {
		return Media{}, errors.New(fmt.Sprintf("Media %s does not exist in storage", key))
	}

	media := Media{Key: key, Author: author, Kind: kind, Folder: folder}

	errSave := mm.mediaRepo.Save(media)

	if errSave != nil {
		return Media{}, fmt.Errorf("Unable to save media %s: %w", key, errFind)
	}

	return media, nil
}

func (mm *MediaManager) GetUploadSignedUri(key string) (string, error) {
	return mm.storage.SignUploadUri(key)
}

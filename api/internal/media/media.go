package media

import (
	"fmt"
	"os"
	"time"
)

type MediaKind string

const (
	KindPhoto MediaKind = "photo"
	KindVideo MediaKind = "video"
)

type Media struct {
	Key        string
	Author     string
	Kind       MediaKind
	Folder     string
	UploadDate int
	Visible    bool
}

type MediaManager struct {
	mediaRepo                                         MediaRepository
	storage                                           Storage
	videoConverter                                    VideoConverter
	bucketVideoRaw, bucketVideoFormatted, bucketImage string
}

type mediaError string

func (e mediaError) Error() string {
	return string(e)
}

const ErrGetByFolder = mediaError("cannot get medias by folder")
const ErrGetByKey = mediaError("cannot get media by key")
const ErrCheckStorage = mediaError("cannot check if media is in storage")
const ErrSaveMedia = mediaError("cannot save media")
const ErrConvert = mediaError("cannot convert media")

func CreateMediaManager(mediaRepo MediaRepository, storage Storage, videoConverter VideoConverter) MediaManager {
	return MediaManager{
		mediaRepo:            mediaRepo,
		storage:              storage,
		bucketVideoRaw:       os.Getenv("BUCKET_VIDEO_RAW"),
		bucketVideoFormatted: os.Getenv("BUCKET_VIDEO_FORMATTED"),
		bucketImage:          os.Getenv("BUCKET_IMAGE"),
		videoConverter:       videoConverter,
	}
}

func (mm *MediaManager) GetMediasByFolder(folder string) ([]Media, error) {
	return mm.mediaRepo.FindByFolder(folder)
}

func (mm *MediaManager) GetFolders(name string) ([]string, error) {
	return mm.mediaRepo.FindFoldersName(name)
}

func (mm *MediaManager) DeleteFolder(name string) error {
	medias, errFind := mm.mediaRepo.FindByFolder(name)

	if errFind != nil {
		return fmt.Errorf("%w", ErrGetByFolder)
	}

	for _, m := range medias {
		m.Folder = "none"
		mm.mediaRepo.Save(m)
	}

	return nil
}

func (mm *MediaManager) GetAll(name string) ([]Media, error) {
	return mm.mediaRepo.FindByFolder(name)
}

func (mm *MediaManager) ChangeMediasFolder(keys []string, newFolder string) ([]Media, error) {
	medias, errFind := mm.mediaRepo.FindManyByKeys(keys)

	if errFind != nil {
		return make([]Media, 0), fmt.Errorf("%w", ErrGetByKey)
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
		return make([]Media, 0), fmt.Errorf("%w", ErrGetByFolder)
	}

	for k := range medias {
		medias[k].Folder = newFolder
		mm.mediaRepo.Save(medias[k])
	}

	return medias, nil
}

func (mm *MediaManager) Ingest(key, author, folder string, kind MediaKind) (Media, error) {
	mediaWithKey, errFind := mm.mediaRepo.FindByKey(key)

	if errFind != nil {
		return Media{}, fmt.Errorf("%w", ErrGetByKey)
	}

	if mediaWithKey != (Media{}) {
		return Media{}, nil
	}

	mediaInStorage, errCheckStorage := mm.storage.MediaExist(key, getIngestBucketFromMediaKind(kind))

	if errCheckStorage != nil {
		return Media{}, fmt.Errorf("%w", ErrCheckStorage)
	}

	if !mediaInStorage {
		return Media{}, nil
	}

	visible := false
	if kind == KindPhoto {
		visible = true
	}

	media := Media{Key: key, Author: author, Kind: kind, Folder: folder, UploadDate: int(time.Now().Unix()), Visible: visible}

	errSave := mm.mediaRepo.Save(media)

	if errSave != nil {
		return Media{}, fmt.Errorf("%w", ErrSaveMedia)
	}

	if kind == KindVideo {
		errConvert := mm.videoConverter.Convert(key)

		if errConvert != nil {
			return Media{}, fmt.Errorf("%w", ErrConvert)
		}
	}

	return media, nil
}

func (mm *MediaManager) GetUploadSignedUri(key string, kind MediaKind) (string, error) {
	return mm.storage.SignUploadUri(key, getIngestBucketFromMediaKind(kind))
}

func (mm *MediaManager) AcknowledgeVideoConversion(key string) error {
	m, errFind := mm.mediaRepo.FindByKey(key)

	if m == (Media{}) {
		return nil
	}

	if errFind != nil {
		return fmt.Errorf("%w", ErrGetByKey)
	}

	m.Visible = true

	return mm.mediaRepo.Save(m)
}

func getIngestBucketFromMediaKind(kind MediaKind) string {
	if kind == KindPhoto {
		return os.Getenv("BUCKET_IMAGE")
	}

	return os.Getenv("BUCKET_VIDEO_RAW")
}

package media

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
)

type VideoConverter interface {
	Convert(key string) error
}

type CloudConvert struct{}

type ImportTask struct {
	Operation       string `json:"operation"`
	Bucket          string `json:"bucket"`
	Region          string `json:"region"`
	Key             string `json:"key"`
	Endpoint        string `json:"endpoint"`
	AccessKeyId     string `json:"access_key_id"`
	SecretAccessKey string `json:"secret_access_key"`
}

type ExportTask struct {
	Operation       string   `json:"operation"`
	Input           []string `json:"input"`
	Bucket          string   `json:"bucket"`
	Region          string   `json:"region"`
	Key             string   `json:"key"`
	Endpoint        string   `json:"endpoint"`
	AccessKeyId     string   `json:"access_key_id"`
	SecretAccessKey string   `json:"secret_access_key"`
}

type ConvertTask struct {
	Operation     string   `json:"operation"`
	InputFormat   string   `json:"input_format"`
	OutputFormat  string   `json:"output_format"`
	Engine        string   `json:"engine"`
	Input         []string `json:"input"`
	VideoCodec    string   `json:"video_codec"`
	Crf           int      `json:"crf"`
	Preset        string   `json:"preset"`
	Width         int      `json:"width"`
	Height        int      `json:"height"`
	Fit           string   `json:"fit"`
	SubtitlesMode string   `json:"subtitles_mode"`
	AudioCodec    string   `json:"audio_codec"`
	AudioBitrate  int      `json:"audio_bitrate"`
}

type Tasks struct {
	Import  ImportTask  `json:"import"`
	Convert ConvertTask `json:"convert"`
	Export  ExportTask  `json:"export"`
}

type Job struct {
	Tasks Tasks  `json:"tasks"`
	Tag   string `json:"tag"`
}

func NewCloudConvert() VideoConverter {
	return &CloudConvert{}
}

func (*CloudConvert) Convert(key string) error {
	importTask := ImportTask{
		Operation:       "import/s3",
		Bucket:          os.Getenv("BUCKET_VIDEO_RAW"),
		Endpoint:        os.Getenv("S3_ENDPOINT"),
		Region:          os.Getenv("S3_REGION"),
		Key:             key,
		AccessKeyId:     os.Getenv("AKID"),
		SecretAccessKey: os.Getenv("ASK"),
	}

	exportTask := ExportTask{
		Operation:       "export/s3",
		Bucket:          os.Getenv("BUCKET_VIDEO_FORMATTED"),
		Endpoint:        os.Getenv("S3_ENDPOINT"),
		Region:          os.Getenv("S3_REGION"),
		Key:             key,
		AccessKeyId:     os.Getenv("AKID"),
		SecretAccessKey: os.Getenv("ASK"),
		Input:           []string{"convert"},
	}

	converTask := ConvertTask{
		Operation:     "convert",
		InputFormat:   "mp4",
		OutputFormat:  "mp4",
		Engine:        "ffmpeg",
		Input:         []string{"import"},
		VideoCodec:    "x264",
		Crf:           23,
		Preset:        "medium",
		Width:         1280,
		Height:        720,
		Fit:           "scale",
		SubtitlesMode: "none",
		AudioCodec:    "aac",
		AudioBitrate:  128,
	}

	tasks := Tasks{
		Import:  importTask,
		Convert: converTask,
		Export:  exportTask,
	}

	job := Job{
		Tasks: tasks,
		Tag:   "jobbuilder",
	}

	jsonData, errMarshal := json.Marshal(job)

	if errMarshal != nil {
		return fmt.Errorf("unable to marshal cloudconvert job:%s ", errMarshal)
	}

	req, _ := http.NewRequest(http.MethodPost, os.Getenv("CLOUD_CONVERT_URI"), bytes.NewBuffer(jsonData))
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("CLOUD_CONVERT_TOKEN")))

	res, errPost := http.DefaultClient.Do(req)

	if res.StatusCode != 200 || errPost != nil {
		return errors.New("unable to create jobs on cloudconvert")
	}

	return nil
}

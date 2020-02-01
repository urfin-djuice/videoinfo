package main

import (
	"fmt"
	"github.com/zelenin/go-mediainfo"
	"strconv"
)

type MediaInfo struct {
	Name     string `json:"name"`
	Width    *int   `json:"width,omitempty"`
	Height   *int   `json:"height,omitempty"`
	BitRate  int    `json:"bit_rate"`
	Duration string `json:"duration"`
}

func newInformer(mediaInfo *mediainfo.File) *Informer {
	return &Informer{mediaInfo}
}

type Informer struct {
	mediaInfo *mediainfo.File
}

func (i Informer) getIntParameter(streamType mediainfo.StreamKind, streamNumber int, parameter string) int {
	strVal := i.mediaInfo.Parameter(streamType, streamNumber, parameter)
	if strVal == "" {
		return 0
	}
	intVal, err := strconv.Atoi(strVal)
	if err != nil {
		return 0
	}
	return intVal
}

func (i Informer) getStreamsCount(streamType mediainfo.StreamKind) int {
	return i.getIntParameter(streamType, 0, "StreamCount")
}

func (i Informer) getBitRate(streamType mediainfo.StreamKind, streamNumber int) int {
	return i.getIntParameter(streamType, streamNumber, "BitRate")
}

func (i Informer) getWidth(streamType mediainfo.StreamKind, streamNumber int) int {
	return i.getIntParameter(streamType, streamNumber, "Width")
}

func (i Informer) getHeight(streamType mediainfo.StreamKind, streamNumber int) int {
	return i.getIntParameter(streamType, streamNumber, "Height")
}

func (i Informer) getDuration(streamType mediainfo.StreamKind, streamNumber int) string {
	durationMs := i.getIntParameter(streamType, streamNumber, "Duration")
	return fmt.Sprintf("%.3fs", float64(durationMs) / 1000.0)
}

func (i Informer) getName(streamType mediainfo.StreamKind, streamNumber int) string {
	name := i.mediaInfo.Parameter(streamType, streamNumber, "CodecID/Hint")
	if name == "" {
		name = i.mediaInfo.Parameter(streamType, streamNumber, "CodecID")
	}
	return name
}

func (i Informer) getStreamInfo(streamType mediainfo.StreamKind, streamNumber int) MediaInfo {
	result := MediaInfo{
		Name:     i.getName(streamType, streamNumber),
		Duration: i.getDuration(streamType, streamNumber),
		BitRate:  i.getBitRate(streamType, streamNumber),
	}
	if streamType == mediainfo.StreamVideo {
		rHeight := i.getHeight(streamType, streamNumber)
		rWidth := i.getWidth(streamType, streamNumber)
		result.Height = &rHeight
		result.Width = &rWidth
	}
	return result
}

func (i Informer) getStreamsInTypeInfo(streamType mediainfo.StreamKind) []MediaInfo {
	streamsCount := i.getStreamsCount(streamType)
	result := []MediaInfo{}
	if streamsCount > 0 {
		for k := 0; k < streamsCount; k++ {
			result = append(result, i.getStreamInfo(streamType, k))
		}
	}
	return result
}

func (i Informer) GetInfo() map[string][]MediaInfo {
	return map[string][]MediaInfo{
		"video": i.getStreamsInTypeInfo(mediainfo.StreamVideo),
		"audio": i.getStreamsInTypeInfo(mediainfo.StreamAudio),
	}
}

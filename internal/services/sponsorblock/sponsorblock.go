package sponsorblock

import (
	"encoding/json"
	"ikoyhn/podcast-sponsorblock/internal/config"
	"ikoyhn/podcast-sponsorblock/internal/database"
	"io"
	"math"
	"net/http"
	"os"
	"strings"

	log "github.com/labstack/gommon/log"
)

const SPONSORBLOCK_API_URL = "https://sponsor.ajay.app/api/skipSegments?videoID="

// HasSegments returns true if SponsorBlock has any segments for the given video.
// Results are cached: once segments are confirmed they are not rechecked; videos
// with no segments are rechecked on each call until segments appear.
func HasSegments(videoId string) bool {
	cache := database.GetSponsorBlockCache(videoId)
	if cache != nil && cache.HasSegments {
		return true
	}

	endURL := SPONSORBLOCK_API_URL + videoId
	if categories := getCategories(); categories != nil {
		for _, category := range categories {
			endURL += "&category=" + strings.TrimSpace(category)
		}
	}

	resp, err := http.Get(endURL)
	if err != nil {
		log.Error("[SponsorBlock] Error checking for segments: ", err)
		return false
	}
	defer resp.Body.Close()

	hasSegments := resp.StatusCode == http.StatusOK
	database.UpsertSponsorBlockCache(videoId, hasSegments)
	return hasSegments
}

func DeterminePodcastDownload(youtubeVideoId string) (bool, float64) {
	episodeHistory := database.GetEpisodePlaybackHistory(youtubeVideoId)

	updatedSkippedTime := TotalSponsorTimeSkipped(youtubeVideoId)
	if episodeHistory == nil {
		return true, updatedSkippedTime
	}

	if math.Abs(episodeHistory.TotalTimeSkipped-updatedSkippedTime) > 2 {
		file := database.FindFileWithId(config.AppConfig.Setup.AudioDir, youtubeVideoId)
		if file != "" {
			os.Remove(file)
		}
		log.Debug("[SponsorBlock] Updating downloaded episode with new sponsor skips...")
		return true, updatedSkippedTime
	}

	return false, updatedSkippedTime
}

func TotalSponsorTimeSkipped(youtubeVideoId string) float64 {
	log.Debug("[SponsorBlock] Looking up podcast in SponsorBlock API...")
	endURL := SPONSORBLOCK_API_URL + youtubeVideoId

	if categories := getCategories(); categories != nil {
		for _, category := range categories {
			endURL += "&category=" + strings.TrimSpace(category)
		}
	}

	resp, err := http.Get(endURL)
	if err != nil {
		log.Error(err)
		return 0
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		log.Warnf("Video not found on SponsorBlock API: %s", youtubeVideoId)
		return 0
	}

	body, bodyErr := io.ReadAll(resp.Body)
	if bodyErr != nil {
		log.Error(bodyErr)
		return 0
	}
	sponsorBlockResponse, marshErr := unmarshalSponsorBlockResponse(body)
	if marshErr != nil {
		log.Error(marshErr)
		return 0
	}

	totalTimeSkipped := calculateSkippedTime(sponsorBlockResponse)

	return totalTimeSkipped
}

func unmarshalSponsorBlockResponse(data []byte) ([]SponsorBlockResponse, error) {
	var res []SponsorBlockResponse

	if err := json.Unmarshal(data, &res); err != nil {
		return []SponsorBlockResponse{}, err
	}

	return res, nil
}

func calculateSkippedTime(segments []SponsorBlockResponse) float64 {
	skippedTime := float64(0)
	prevStopTime := float64(0)

	for _, segment := range segments {
		startTime := segment.Segment[0]
		stopTime := segment.Segment[1]

		if startTime > prevStopTime {
			skippedTime += stopTime - startTime
		} else {
			skippedTime += stopTime - prevStopTime
		}

		prevStopTime = stopTime
	}

	return skippedTime
}

func getCategories() []string {
	if config.AppConfig.Ytdlp.SponsorBlockCategories == "" {
		return nil
	}
	return strings.Split(config.AppConfig.Ytdlp.SponsorBlockCategories, ",")
}

type SponsorBlockResponse struct {
	Segment       []float64 `json:"segment"`
	UUID          string    `json:"UUID"`
	Category      string    `json:"category"`
	VideoDuration float64   `json:"videoDuration"`
	ActionType    string    `json:"actionType"`
	Locked        int16     `json:"locked"`
	Votes         int16     `json:"votes"`
	Description   string    `json:"description"`
}

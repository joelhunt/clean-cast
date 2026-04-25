package database

import (
	"ikoyhn/podcast-sponsorblock/internal/models"
	"time"
)

func GetSponsorBlockCache(videoId string) *models.SponsorBlockCache {
	var cache models.SponsorBlockCache
	if err := db.Where("youtube_video_id = ?", videoId).First(&cache).Error; err != nil {
		return nil
	}
	return &cache
}

func UpsertSponsorBlockCache(videoId string, hasSegments bool) {
	cache := models.SponsorBlockCache{
		YoutubeVideoId: videoId,
		HasSegments:    hasSegments,
		CheckedAt:      time.Now().Unix(),
	}
	db.Save(&cache)
}

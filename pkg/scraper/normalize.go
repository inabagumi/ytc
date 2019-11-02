package scraper

import (
	"html"
	"time"

	"google.golang.org/api/youtube/v3"
)

type Channel struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	URL   string `json:"url"`
}

type Video struct {
	Channel     *Channel `json:"channel"`
	ID          string   `json:"id"`
	ObjectID    string   `json:"objectID"`
	PublishedAt int64    `json:"publishedAt"`
	Title       string   `json:"title"`
	URL         string   `json:"url"`
}

func normalize(item *youtube.Video) *Video {
	rawPublishedAt := item.Snippet.PublishedAt

	if item.LiveStreamingDetails != nil {
		if item.LiveStreamingDetails.ActualStartTime != "" {
			rawPublishedAt = item.LiveStreamingDetails.ActualStartTime
		} else if item.LiveStreamingDetails.ScheduledStartTime != "" {
			rawPublishedAt = item.LiveStreamingDetails.ScheduledStartTime
		}
	}

	publishedAt, err := time.Parse(time.RFC3339, rawPublishedAt)
	if err != nil {
		publishedAt = time.Now()
	}

	channel := &Channel{
		ID:    item.Snippet.ChannelId,
		Title: item.Snippet.ChannelTitle,
		URL:   "https://www.youtube.com/channel/" + item.Snippet.ChannelId,
	}

	video := &Video{
		Channel:     channel,
		ID:          item.Id,
		ObjectID:    item.Id,
		PublishedAt: publishedAt.Unix(),
		Title:       html.UnescapeString(item.Snippet.Title),
		URL:         "https://www.youtube.com/watch?v=" + item.Id,
	}

	return video
}

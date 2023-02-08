package video

import "ByteDance_5th/models"

type FeedList struct {
	List     []*models.Video `json:"video_list,omitempty"`
	NextTime int64           `json:"next_time,omitempty"`
}

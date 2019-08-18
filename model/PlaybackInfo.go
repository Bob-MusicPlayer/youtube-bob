package model

type PlaybackInfo struct {
	CachedTime float64 `json:"cachedTime"`
	Position   float64 `json:"position"`
	Duration   float64 `json:"duration"`
	Title      string  `json:"title"`
}

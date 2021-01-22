package main

import "time"

// this should probably be a real struct, but it never changes,
// so for now just make it a constant string
const errorResponse string = `{
    "error": "Could not decode request: JSON parsing failed"
}`

type request struct {
	Skip         int
	Take         int
	TotalRecords int
	Payloads     []payload `json:"payload"`
}

type payload struct {
	Country      string
	Description  string
	DRM          bool
	EpisodeCount int
	Genre        string
	Image        struct {
		ShowImage string
	}
	Language      string
	NextEpisode   episode `json:"nextEpisode,omitempty"`
	PrimaryColour string
	Seasons       []struct {
		Slug string
	}
	Slug      string
	Title     string
	TVChannel string
}

type episode struct {
	Channel     string `json:"channel,omitempty"`
	ChannelLogo string
	Date        time.Time
	HTML        string
	URL         string
}

type responseWrapper struct {
	Response []response
}

type response struct {
	Image string
	Slug  string
	Title string
}

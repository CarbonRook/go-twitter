package twitter

// Entities represent metadata and context info parsed from Twitter components.
// https://developer.twitter.com/en/docs/tweets/data-dictionary/overview/entities-object
// TODO: symbols
type Entities struct {
	Hashtags     []TagEntity         `json:"hashtags"`
	Urls         []URLEntity         `json:"urls"`
	UserMentions []MentionEntity     `json:"user_mentions"`
	Annotations  []ContextAnnotation `json:"annotations"`
	Cashtags     []TagEntity         `json:"cashtags"`
}

// TagEntity represents a hashtag or cashtag from the text.
type TagEntity struct {
	Start int64  `json:"start"`
	End   int64  `json:"end"`
	Tag   string `json:"tag"`
}

// AnnotationEntity represents one of Twitter's context annotations.
type AnnotationEntity struct {
	Start          int64  `json:"start"`
	End            int64  `json:"end"`
	Probability    int64  `json:"probability"`
	Type           string `json:"type"`
	NormalizedText string `json:"normalized_text"`
}

// URLEntity represents a URL which has been parsed from text.
type URLEntity struct {
	Start       int64  `json:"start"`
	End         int64  `json:"end"`
	DisplayURL  string `json:"display_url"`
	ExpandedURL string `json:"expanded_url"`
	URL         string `json:"url"`
	UnwoundURL  string `json:"unwound_url"`
}

// MediaEntity represents media elements associated with a Tweet.
type MediaEntity struct {
	MediaKey         string       `json:"media_key"`
	DurationMillis   int64        `json:"duration_ms"`
	Height           int64        `json:"height"`
	NonPublicMetrics MediaMetrics `json:"non_public_metrics"`
	OrganicMetrics   MediaMetrics `json:"organic_metrics"`
	Type             string       `json:"type"`
	PreviewImageURL  string       `json:"preview_image_url"`
	PromotedMetrics  MediaMetrics `json:"promoted_metrics"`
	PublicMetrics    MediaMetrics `json:"public_metrics"`
	Width            int64        `json:"width"`
	AltText          string       `json:"alt_text"`
}

type MediaMetrics struct {
	Playback0Count   int64 `json:"playback_0_count,omitempty"`
	Playback25Count  int64 `json:"playback_25_count,omitempty"`
	Playback50Count  int64 `json:"playback_50_count,omitempty"`
	Playback75Count  int64 `json:"playback_75_count,omitempty"`
	Playback100Count int64 `json:"playback_100_count,omitempty"`
	ViewCount        int64 `json:"view_count,omitempty"`
}

// MentionEntity represents Twitter user mentions parsed from text.
type MentionEntity struct {
	Start    int64  `json:"start"`
	End      int64  `json:"end"`
	Username string `json:"username"`
}

// UserEntities contain Entities parsed from User url and description fields.
// https://developer.twitter.com/en/docs/tweets/data-dictionary/overview/entities-object#mentions
type UserEntities struct {
	URL         Entities `json:"url"`
	Description Entities `json:"description"`
}

// ExtendedEntity contains media information.
// https://developer.twitter.com/en/docs/tweets/data-dictionary/overview/extended-entities-object
type ExtendedEntity struct {
	Media []MediaEntity `json:"media"`
}

// MediaSizes contain the different size media that are available.
// https://developer.twitter.com/en/docs/tweets/data-dictionary/overview/entities-object#media-size
type MediaSizes struct {
	Thumb  MediaSize `json:"thumb"`
	Large  MediaSize `json:"large"`
	Medium MediaSize `json:"medium"`
	Small  MediaSize `json:"small"`
}

// MediaSize describes the height, width, and resizing method used.
type MediaSize struct {
	Width  int    `json:"w"`
	Height int    `json:"h"`
	Resize string `json:"resize"`
}

// VideoInfo is available on video media objects.
type VideoInfo struct {
	AspectRatio    [2]int         `json:"aspect_ratio"`
	DurationMillis int            `json:"duration_millis"`
	Variants       []VideoVariant `json:"variants"`
}

// VideoVariant describes one of the available video formats.
type VideoVariant struct {
	ContentType string `json:"content_type"`
	Bitrate     int    `json:"bitrate"`
	URL         string `json:"url"`
}

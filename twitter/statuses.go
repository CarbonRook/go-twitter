package twitter

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dghubble/sling"
)

// Tweet represents a Twitter Tweet, previously called a status.
// https://dev.twitter.com/overview/api/tweets
type Tweet struct {
	Attachments struct {
		MediaKeys []string `json:"media_keys,omitempty"`
		PollID    []string `json:"poll_id,omitempty"`
	} `json:"attachments,omitempty"`
	AuthorID           string               `json:"author_id"`
	ContextAnnotations []*ContextAnnotation `json:"context_annotations"`
	ConversationID     string               `json:"conversation_id"`
	CreatedAt          string               `json:"created_at"`
	Entities           Entities             `json:"entities"`
	Geo                *Geo                 `json:"geo,omitempty"`
	Includes           *Includes            `json:"includes"`
	ID                 string               `json:"id"`
	InReplyToStatusID  string               `json:"in_reply_to_status_id"`
	InReplyToUserID    string               `json:"in_reply_to_user_id"`
	Lang               string               `json:"lang"`
	PossiblySensitive  bool                 `json:"possibly_sensitive"`
	ReferencedTweets   struct {
		Type string `json:"type"`
		ID   string `json:"id"`
	} `json:"referenced_tweets,omitempty"`
	PublicMetrics    *Metrics  `json:"public_metrics,omitempty"`
	NonPublicMetrics *Metrics  `json:"non_public_metrics,omitempty"`
	OrganicMetrices  *Metrics  `json:"organic_metrics,omitempty"`
	PromotedMetrics  *Metrics  `json:"promoted_metrics,omitempty"`
	ReplySettings    string    `json:"reply_settings"`
	Source           string    `json:"source"`
	Text             string    `json:"text"`
	Withheld         *Withheld `json:"withheld"`
}

// CreatedAtTime returns the time a tweet was created.
func (t Tweet) CreatedAtTime() (time.Time, error) {
	return time.Parse(time.RubyDate, t.CreatedAt)
}

type Geo struct {
	PlaceID string `json:"place_id"`
}

// Withheld represents the reasons why a tweet may be withheld (copyright) and from which countries.
type Withheld struct {
	Copyright   bool     `json:"copyright"`
	CountryCode []string `json:"country_code"`
	Scope       string   `json:"scope"`
}

//Metrics list metrics associated with the tweet (counts of retweets, replies, likes, quotes).
type Metrics struct {
	RetweetCount int64 `json:"retweet_count"`
	ReplyCount   int64 `json:"reply_count"`
	LikeCount    int64 `json:"like_count"`
	QuoteCount   int64 `json:"quote_count,omitempty"`
}

type ContextAnnotation struct {
	Domain struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
	} `json:"domain"`
	Entity struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
}

type Poll struct {
	ID              string       `json:"id"`
	Options         []PollOption `json:"options"`
	DurationMinutes int64        `json:"duration_minutes"`
	End             string       `json:"end_datetime"`
	VotingStatus    string       `json:"voting_status"`
}

type PollOption struct {
	Position int64  `json:"position"`
	Label    string `json:"label"`
	Votes    int64  `json:"votes"`
}

// Place represents a Twitter Place / Location
// https://dev.twitter.com/overview/api/places
type Place struct {
	FullName        string   `json:"full_name"`
	ID              string   `json:"id"`
	ContainedWithin []string `json:"contained_within"`
	Country         string   `json:"country"`
	CountryCode     string   `json:"country_code"`
	Geo             struct {
		Type        string    `json:"type"`
		BoundingBox []float64 `json:"bbox"`
		Properties  struct{}  `json:"properties"`
	} `json:"geo"`
	Name      string `json:"name"`
	PlaceType string `json:"place_type"`
}

// StatusService provides methods for accessing Twitter status API endpoints.
type StatusService struct {
	sling *sling.Sling
}

// newStatusService returns a new StatusService.
func newStatusService(sling *sling.Sling) *StatusService {
	return &StatusService{
		sling: sling.Path("statuses/"),
	}
}

// StatusShowParams are the parameters for StatusService.Show
type StatusShowParams struct {
	ID               int64  `url:"id,omitempty"`
	TrimUser         *bool  `url:"trim_user,omitempty"`
	IncludeMyRetweet *bool  `url:"include_my_retweet,omitempty"`
	IncludeEntities  *bool  `url:"include_entities,omitempty"`
	TweetMode        string `url:"tweet_mode,omitempty"`
}

// Show returns the requested Tweet.
// https://dev.twitter.com/rest/reference/get/statuses/show/%3Aid
func (s *StatusService) Show(id int64, params *StatusShowParams) (*Tweet, *http.Response, error) {
	if params == nil {
		params = &StatusShowParams{}
	}
	params.ID = id
	tweet := new(Tweet)
	apiError := new(APIError)
	resp, err := s.sling.New().Get("show.json").QueryStruct(params).Receive(tweet, apiError)
	return tweet, resp, relevantError(err, *apiError)
}

// StatusLookupParams are the parameters for StatusService.Lookup
type StatusLookupParams struct {
	ID              []int64 `url:"id,omitempty,comma"`
	TrimUser        *bool   `url:"trim_user,omitempty"`
	IncludeEntities *bool   `url:"include_entities,omitempty"`
	Map             *bool   `url:"map,omitempty"`
	TweetMode       string  `url:"tweet_mode,omitempty"`
}

// Lookup returns the requested Tweets as a slice. Combines ids from the
// required ids argument and from params.Id.
// https://dev.twitter.com/rest/reference/get/statuses/lookup
func (s *StatusService) Lookup(ids []int64, params *StatusLookupParams) ([]Tweet, *http.Response, error) {
	if params == nil {
		params = &StatusLookupParams{}
	}
	params.ID = append(params.ID, ids...)
	tweets := new([]Tweet)
	apiError := new(APIError)
	resp, err := s.sling.New().Get("lookup.json").QueryStruct(params).Receive(tweets, apiError)
	return *tweets, resp, relevantError(err, *apiError)
}

// StatusUpdateParams are the parameters for StatusService.Update
type StatusUpdateParams struct {
	Status             string   `url:"status,omitempty"`
	InReplyToStatusID  int64    `url:"in_reply_to_status_id,omitempty"`
	PossiblySensitive  *bool    `url:"possibly_sensitive,omitempty"`
	Lat                *float64 `url:"lat,omitempty"`
	Long               *float64 `url:"long,omitempty"`
	PlaceID            string   `url:"place_id,omitempty"`
	DisplayCoordinates *bool    `url:"display_coordinates,omitempty"`
	TrimUser           *bool    `url:"trim_user,omitempty"`
	MediaIds           []int64  `url:"media_ids,omitempty,comma"`
	TweetMode          string   `url:"tweet_mode,omitempty"`
}

// Update updates the user's status, also known as Tweeting.
// Requires a user auth context.
// https://dev.twitter.com/rest/reference/post/statuses/update
func (s *StatusService) Update(status string, params *StatusUpdateParams) (*Tweet, *http.Response, error) {
	if params == nil {
		params = &StatusUpdateParams{}
	}
	params.Status = status
	tweet := new(Tweet)
	apiError := new(APIError)
	resp, err := s.sling.New().Post("update.json").BodyForm(params).Receive(tweet, apiError)
	return tweet, resp, relevantError(err, *apiError)
}

// StatusRetweetParams are the parameters for StatusService.Retweet
type StatusRetweetParams struct {
	ID        int64  `url:"id,omitempty"`
	TrimUser  *bool  `url:"trim_user,omitempty"`
	TweetMode string `url:"tweet_mode,omitempty"`
}

// Retweet retweets the Tweet with the given id and returns the original Tweet
// with embedded retweet details.
// Requires a user auth context.
// https://dev.twitter.com/rest/reference/post/statuses/retweet/%3Aid
func (s *StatusService) Retweet(id int64, params *StatusRetweetParams) (*Tweet, *http.Response, error) {
	if params == nil {
		params = &StatusRetweetParams{}
	}
	params.ID = id
	tweet := new(Tweet)
	apiError := new(APIError)
	path := fmt.Sprintf("retweet/%d.json", params.ID)
	resp, err := s.sling.New().Post(path).BodyForm(params).Receive(tweet, apiError)
	return tweet, resp, relevantError(err, *apiError)
}

// StatusUnretweetParams are the parameters for StatusService.Unretweet
type StatusUnretweetParams struct {
	ID        int64  `url:"id,omitempty"`
	TrimUser  *bool  `url:"trim_user,omitempty"`
	TweetMode string `url:"tweet_mode,omitempty"`
}

// Unretweet unretweets the Tweet with the given id and returns the original Tweet.
// Requires a user auth context.
// https://dev.twitter.com/rest/reference/post/statuses/unretweet/%3Aid
func (s *StatusService) Unretweet(id int64, params *StatusUnretweetParams) (*Tweet, *http.Response, error) {
	if params == nil {
		params = &StatusUnretweetParams{}
	}
	params.ID = id
	tweet := new(Tweet)
	apiError := new(APIError)
	path := fmt.Sprintf("unretweet/%d.json", params.ID)
	resp, err := s.sling.New().Post(path).BodyForm(params).Receive(tweet, apiError)
	return tweet, resp, relevantError(err, *apiError)
}

// StatusRetweetsParams are the parameters for StatusService.Retweets
type StatusRetweetsParams struct {
	ID        int64  `url:"id,omitempty"`
	Count     int    `url:"count,omitempty"`
	TrimUser  *bool  `url:"trim_user,omitempty"`
	TweetMode string `url:"tweet_mode,omitempty"`
}

// Retweets returns the most recent retweets of the Tweet with the given id.
// https://dev.twitter.com/rest/reference/get/statuses/retweets/%3Aid
func (s *StatusService) Retweets(id int64, params *StatusRetweetsParams) ([]Tweet, *http.Response, error) {
	if params == nil {
		params = &StatusRetweetsParams{}
	}
	params.ID = id
	tweets := new([]Tweet)
	apiError := new(APIError)
	path := fmt.Sprintf("retweets/%d.json", params.ID)
	resp, err := s.sling.New().Get(path).QueryStruct(params).Receive(tweets, apiError)
	return *tweets, resp, relevantError(err, *apiError)
}

// StatusDestroyParams are the parameters for StatusService.Destroy
type StatusDestroyParams struct {
	ID        int64  `url:"id,omitempty"`
	TrimUser  *bool  `url:"trim_user,omitempty"`
	TweetMode string `url:"tweet_mode,omitempty"`
}

// Destroy deletes the Tweet with the given id and returns it if successful.
// Requires a user auth context.
// https://dev.twitter.com/rest/reference/post/statuses/destroy/%3Aid
func (s *StatusService) Destroy(id int64, params *StatusDestroyParams) (*Tweet, *http.Response, error) {
	if params == nil {
		params = &StatusDestroyParams{}
	}
	params.ID = id
	tweet := new(Tweet)
	apiError := new(APIError)
	path := fmt.Sprintf("destroy/%d.json", params.ID)
	resp, err := s.sling.New().Post(path).BodyForm(params).Receive(tweet, apiError)
	return tweet, resp, relevantError(err, *apiError)
}

// OEmbedTweet represents a Tweet in oEmbed format.
type OEmbedTweet struct {
	URL          string `json:"url"`
	ProviderURL  string `json:"provider_url"`
	ProviderName string `json:"provider_name"`
	AuthorName   string `json:"author_name"`
	Version      string `json:"version"`
	AuthorURL    string `json:"author_url"`
	Type         string `json:"type"`
	HTML         string `json:"html"`
	Height       int64  `json:"height"`
	Width        int64  `json:"width"`
	CacheAge     string `json:"cache_age"`
}

// StatusOEmbedParams are the parameters for StatusService.OEmbed
type StatusOEmbedParams struct {
	ID         int64  `url:"id,omitempty"`
	URL        string `url:"url,omitempty"`
	Align      string `url:"align,omitempty"`
	MaxWidth   int64  `url:"maxwidth,omitempty"`
	HideMedia  *bool  `url:"hide_media,omitempty"`
	HideThread *bool  `url:"hide_media,omitempty"`
	OmitScript *bool  `url:"hide_media,omitempty"`
	WidgetType string `url:"widget_type,omitempty"`
	HideTweet  *bool  `url:"hide_tweet,omitempty"`
}

// OEmbed returns the requested Tweet in oEmbed format.
// https://dev.twitter.com/rest/reference/get/statuses/oembed
func (s *StatusService) OEmbed(params *StatusOEmbedParams) (*OEmbedTweet, *http.Response, error) {
	oEmbedTweet := new(OEmbedTweet)
	apiError := new(APIError)
	resp, err := s.sling.New().Get("oembed.json").QueryStruct(params).Receive(oEmbedTweet, apiError)
	return oEmbedTweet, resp, relevantError(err, *apiError)
}

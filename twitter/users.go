package twitter

import (
	"net/http"

	"github.com/dghubble/sling"
)

// User represents a Twitter User.
// https://dev.twitter.com/overview/api/users
type User struct {
	ID              string            `json:"id"`
	Name            string            `json:"name"`
	Username        string            `json:"username"`
	CreatedAt       string            `json:"created_at"`
	Description     string            `json:"description"`
	Entities        *UserEntities     `json:"entities"`
	Location        string            `json:"location"`
	PinnedTweetID   string            `json:"pinned_tweet_id"`
	ProfileImageURL string            `json:"profile_image_url"`
	Protected       bool              `json:"protected"`
	PublicMetrics   UserPublicMetrics `json:"public_metrics"`
	URL             string            `json:"url"`
	Verified        bool              `json:"verified"`
	Withheld        Withheld          `json:"withheld"`
}

type UserPublicMetrics struct {
	FollowersCount int64 `json:"followers_count"`
	FollowingCount int64 `json:"following_count"`
	TweetCount     int64 `json:"tweet_count"`
	ListedCount    int64 `json:"listed_count"`
}

// UserService provides methods for accessing Twitter user API endpoints.
type UserService struct {
	sling *sling.Sling
}

// newUserService returns a new UserService.
func newUserService(sling *sling.Sling) *UserService {
	return &UserService{
		sling: sling.Path("users/"),
	}
}

// UserShowParams are the parameters for UserService.Show.
type UserServiceParams struct {
	UserFields  []string `url:"user.fields,omitempty"`
	Expansions  []string `url:"expansions,omitempty"`
	TweetFields []string `url:"tweet_fields,omitempty"`
}

func (s *UserService) UserByID(userid string, params *UserServiceParams) (*User, *http.Response, error) {
	user := new(User)
	apiError := new(APIError)
	resp, err := s.sling.New().Get(userid).QueryStruct(params).Receive(user, apiError)
	return user, resp, relevantError(err, *apiError)
}

func (s *UserService) UserByUsername(username string, params *UserServiceParams) (*User, *http.Response, error) {
	user := new(User)
	apiError := new(APIError)
	resp, err := s.sling.New().Get("by/username/").Get(username).QueryStruct(params).Receive(user, apiError)
	return user, resp, relevantError(err, *apiError)
}

func (s *UserService) AuthenticatedUser(params *UserServiceParams) (*User, *http.Response, error) {
	user := new(User)
	apiError := new(APIError)
	resp, err := s.sling.New().Get("me").QueryStruct(params).Receive(user, apiError)
	return user, resp, relevantError(err, *apiError)
}

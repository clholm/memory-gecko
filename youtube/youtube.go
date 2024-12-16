package youtube

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

// SearchResult represents a simplified YouTube search result
type SearchResult struct {
	Title       string
	VideoID     string
	Description string
	PublishedAt string
}

// generates a random 4-digit number and returns it as a zero-padded string
func generateRandomXXXX() string {
	// create a new random number generator seeded with current time
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	// generate random number between 0 and 9999
	num := r.Intn(10000)
	// format as 4-digit zero-padded string
	return fmt.Sprintf("%04d", num)
}

// SearchVideos performs a YouTube search for videos matching IMG_XXXX pattern
func SearchVideos(ctx context.Context, apiKey string) ([]SearchResult, error) {
	// create youtube service with API key
	service, err := youtube.NewService(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("error creating youtube client: %w", err)
	}

	// generate a random 4-digit number
	randomNum := generateRandomXXXX()
	searchString := fmt.Sprintf("IMG_%s", randomNum)

	// create search request
	// 10 searches is hardcoded for right now
	call := service.Search.List([]string{"id", "snippet"}).
		Q(searchString).
		MaxResults(10).
		Type("video")

	// execute search
	response, err := call.Do()
	if err != nil {
		return nil, fmt.Errorf("error performing search: %w", err)
	}

	// process results
	var results []SearchResult
	for _, item := range response.Items {
		result := SearchResult{
			Title:       item.Snippet.Title,
			VideoID:     item.Id.VideoId,
			Description: item.Snippet.Description,
			PublishedAt: item.Snippet.PublishedAt,
		}
		results = append(results, result)
	}

	return results, nil
}

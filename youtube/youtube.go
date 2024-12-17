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

// performs a YouTube search for videos matching a specific IMG_XXXX pattern
// and returns the first result
func searchSingleVideo(ctx context.Context, service *youtube.Service, pattern string) (*SearchResult, error) {
	// create search request for a single video
	call := service.Search.List([]string{"id", "snippet"}).
		Q(pattern).
		MaxResults(1).
		Type("video")

	// execute search
	response, err := call.Do()
	if err != nil {
		return nil, fmt.Errorf("error performing search for %s: %w", pattern, err)
	}

	// if no results found, return nil
	if len(response.Items) == 0 {
		return nil, nil
	}

	// return first result
	item := response.Items[0]
	return &SearchResult{
		Title:       item.Snippet.Title,
		VideoID:     item.Id.VideoId,
		Description: item.Snippet.Description,
		PublishedAt: item.Snippet.PublishedAt,
	}, nil
}

// performs multiple YouTube searches for videos matching different IMG_XXXX patterns
func SearchVideos(ctx context.Context, apiKey string) ([]SearchResult, error) {
	// create youtube service with API key
	service, err := youtube.NewService(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("error creating youtube client: %w", err)
	}

	var results []SearchResult
	attempts := 0
	maxAttempts := 20 // limit total attempts to avoid infinite loop

	// keep searching until we have 10 videos or hit max attempts
	for len(results) < 10 && attempts < maxAttempts {
		// generate a new random pattern
		pattern := fmt.Sprintf("IMG_%s", generateRandomXXXX())

		// search for a single video with this pattern
		result, err := searchSingleVideo(ctx, service, pattern)
		if err != nil {
			return nil, fmt.Errorf("error in search attempt %d: %w", attempts, err)
		}

		// if we found a video, add it to results
		if result != nil {
			results = append(results, *result)
			fmt.Printf("found video %d/%d (pattern: %s)\n", len(results), 10, pattern)
		}

		attempts++
	}

	// check if we got enough videos
	if len(results) < 10 {
		return nil, fmt.Errorf("could only find %d videos after %d attempts", len(results), attempts)
	}

	return results, nil
}

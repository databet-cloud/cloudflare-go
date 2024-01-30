package cloudflare

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/goccy/go-json"
)

var (
	// ErrMissingLiveInputID is for when a LiveInputID is required but missing.
	ErrMissingLiveInputID = errors.New("required live input id missing")
)

// StreamLiveInputListItem represents a stream live input for the list result.
type StreamLiveInputListItem struct {
	UID                      string                 `json:"uid,omitempty"`
	Created                  *time.Time             `json:"created,omitempty"`
	Modified                 *time.Time             `json:"modified,omitempty"`
	Meta                     map[string]interface{} `json:"meta,omitempty"`
	DeleteRecordingAfterDays int                    `json:"deleteRecordingAfterDays,omitempty"`
}

// StreamLiveInput represents a stream live input.
type StreamLiveInput struct {
	UID                      string                   `json:"uid,omitempty"`
	RTMPS                    StreamLiveInputRTMPS     `json:"rtmps,omitempty"`
	RTMPSPlayback            StreamLiveInputRTMPS     `json:"rtmpsPlayback,omitempty"`
	SRT                      StreamLiveInputSRT       `json:"srt,omitempty"`
	SRTPlayback              StreamLiveInputRTMPS     `json:"srtPlayback,omitempty"`
	WebRTC                   StreamLiveInputWebRTC    `json:"webRTC,omitempty"`
	WebRTCPlayback           StreamLiveInputWebRTC    `json:"webRTCPlayback,omitempty"`
	Created                  *time.Time               `json:"created,omitempty"`
	Modified                 *time.Time               `json:"modified,omitempty"`
	Meta                     map[string]interface{}   `json:"meta,omitempty"`
	DefaultCreator           string                   `json:"defaultCreator,omitempty"`
	Status                   *StreamLiveInputStatuses `json:"status,omitempty"`
	Recording                StreamLiveInputRecording `json:"recording,omitempty"`
	DeleteRecordingAfterDays int                      `json:"deleteRecordingAfterDays,omitempty"`
	PreferLowLatency         bool                     `json:"preferLowLatency"`
}

// StreamLiveInputRTMPS represents the live input values for RTMPS.
type StreamLiveInputRTMPS struct {
	URL       string `json:"url,omitempty"`
	StreamKey string `json:"streamKey,omitempty"`
}

// StreamLiveInputSRT represents the live input values for SRT.
type StreamLiveInputSRT struct {
	URL        string `json:"url,omitempty"`
	StreamID   string `json:"streamId,omitempty"`
	Passphrase string `json:"passphrase,omitempty"`
}

// StreamLiveInputWebRTC represents the live input values for Web-RTC.
type StreamLiveInputWebRTC struct {
	URL string `json:"url,omitempty"`
}

// StreamLiveInputStatuses represents the values streaming statuses for live input.
type StreamLiveInputStatuses struct {
	Current StreamLiveInputStatus   `json:"current,omitempty"`
	History []StreamLiveInputStatus `json:"history,omitempty"`
}

// StreamLiveInputStatus represents the streaming status for live input.
type StreamLiveInputStatus struct {
	Reason          string     `json:"reason,omitempty"`
	State           string     `json:"state,omitempty"`
	StatusEnteredAt *time.Time `json:"statusEnteredAt,omitempty"`
	StatusLastSeen  *time.Time `json:"statusLastSeen,omitempty"`
}

// StreamLiveInputRecording represents the recording configuration for the live input value.
type StreamLiveInputRecording struct {
	Mode              string   `json:"mode,omitempty"`
	RequireSignedURLs bool     `json:"requireSignedURLs,omitempty"`
	AllowedOrigins    []string `json:"allowedOrigins,omitempty"`
	TimeoutSeconds    int      `json:"timeoutSeconds,omitempty"`
}

// ListStreamLiveInputsParameters represents parameters used when listing stream live inputs.
type ListStreamLiveInputsParameters struct {
	AccountID     string
	IncludeCounts bool `url:"include_counts,omitempty"`
}

// CreateStreamLiveInputParameters represents parameters used when creating stream live input.
type CreateStreamLiveInputParameters struct {
	AccountID                string
	DefaultCreator           string                   `json:"defaultCreator,omitempty"`
	DeleteRecordingAfterDays int                      `json:"deleteRecordingAfterDays,omitempty"`
	Meta                     map[string]any           `json:"meta,omitempty"`
	Recording                StreamLiveInputRecording `json:"recording,omitempty"`
	PreferLowLatency         bool                     `json:"preferLowLatency,omitempty"`
}

// StreamLiveInputParameters represents parameters used for stream live input.
type StreamLiveInputParameters struct {
	AccountID   string
	LiveInputID string
}

// UpdateStreamLiveInputParameters represents parameters used when creating stream live input.
type UpdateStreamLiveInputParameters struct {
	AccountID                string
	LiveInputID              string
	DefaultCreator           string                   `json:"defaultCreator,omitempty"`
	DeleteRecordingAfterDays int                      `json:"deleteRecordingAfterDays,omitempty"`
	Meta                     map[string]any           `json:"meta,omitempty"`
	Recording                StreamLiveInputRecording `json:"recording,omitempty"`
	PreferLowLatency         bool                     `json:"preferLowLatency,omitempty"`
}

// StreamLiveInputsListResponse represents an API response of stream live inputs.
type StreamLiveInputsListResponse struct {
	Response
	Result []StreamLiveInputListItem `json:"result,omitempty"`
}

// StreamLiveInputResponse represents an API response of stream live input.
type StreamLiveInputResponse struct {
	Response
	Result StreamLiveInput `json:"result,omitempty"`
}

// ListStreamLiveInputs list live inputs.
//
// API Reference: https://developers.cloudflare.com/api/operations/stream-live-inputs-list-live-inputs
func (api *API) ListStreamLiveInputs(
	ctx context.Context,
	options ListStreamLiveInputsParameters,
) ([]StreamLiveInputListItem, error) {
	if options.AccountID == "" {
		return nil, ErrMissingAccountID
	}

	uri := buildURI(fmt.Sprintf("/accounts/%s/stream/live_inputs", options.AccountID), options)
	res, err := api.makeRequestContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return []StreamLiveInputListItem{}, err
	}

	var streamListResponse StreamLiveInputsListResponse
	if err := json.Unmarshal(res, &streamListResponse); err != nil {
		return []StreamLiveInputListItem{}, err
	}

	return streamListResponse.Result, nil
}

// CreateStreamLiveInput create stream live input.
//
// API Reference: https://developers.cloudflare.com/api/operations/stream-live-inputs-create-a-live-input
func (api *API) CreateStreamLiveInput(
	ctx context.Context,
	options CreateStreamLiveInputParameters,
) (StreamLiveInput, error) {
	if options.AccountID == "" {
		return StreamLiveInput{}, ErrMissingAccountID
	}

	uri := fmt.Sprintf("/accounts/%s/stream/live_inputs", options.AccountID)
	res, err := api.makeRequestContext(ctx, http.MethodPost, uri, options)
	if err != nil {
		return StreamLiveInput{}, err
	}

	var streamListResponse StreamLiveInputResponse
	if err := json.Unmarshal(res, &streamListResponse); err != nil {
		return StreamLiveInput{}, err
	}

	return streamListResponse.Result, nil
}

// DeleteStreamLiveInput delete stream live input.
//
// API Reference: https://developers.cloudflare.com/api/operations/stream-live-inputs-delete-a-live-input
func (api *API) DeleteStreamLiveInput(
	ctx context.Context,
	options StreamLiveInputParameters,
) error {
	if options.AccountID == "" {
		return ErrMissingAccountID
	}
	if options.LiveInputID == "" {
		return ErrMissingLiveInputID
	}

	uri := fmt.Sprintf("/accounts/%s/stream/live_inputs/%s", options.AccountID, options.LiveInputID)
	if _, err := api.makeRequestContext(ctx, http.MethodDelete, uri, options); err != nil {
		return err
	}

	return nil
}

// GetStreamLiveInput get stream live input.
//
// API Reference: https://developers.cloudflare.com/api/operations/stream-live-inputs-retreive-a-live-input
func (api *API) GetStreamLiveInput(
	ctx context.Context,
	options StreamLiveInputParameters,
) (StreamLiveInput, error) {
	if options.AccountID == "" {
		return StreamLiveInput{}, ErrMissingAccountID
	}
	if options.LiveInputID == "" {
		return StreamLiveInput{}, ErrMissingLiveInputID
	}

	uri := fmt.Sprintf("/accounts/%s/stream/live_inputs/%s", options.AccountID, options.LiveInputID)
	res, err := api.makeRequestContext(ctx, http.MethodPost, uri, options)
	if err != nil {
		return StreamLiveInput{}, err
	}

	var streamListResponse StreamLiveInputResponse
	if err := json.Unmarshal(res, &streamListResponse); err != nil {
		return StreamLiveInput{}, err
	}

	return streamListResponse.Result, nil
}

// UpdateStreamLiveInput update stream live input.
//
// API Reference: https://developers.cloudflare.com/api/operations/stream-live-inputs-update-a-live-input
func (api *API) UpdateStreamLiveInput(
	ctx context.Context,
	options UpdateStreamLiveInputParameters,
) (StreamLiveInput, error) {
	if options.AccountID == "" {
		return StreamLiveInput{}, ErrMissingAccountID
	}
	if options.LiveInputID == "" {
		return StreamLiveInput{}, ErrMissingLiveInputID
	}

	uri := fmt.Sprintf("/accounts/%s/stream/live_inputs/%s", options.AccountID, options.LiveInputID)
	res, err := api.makeRequestContext(ctx, http.MethodPut, uri, options)
	if err != nil {
		return StreamLiveInput{}, err
	}

	var streamListResponse StreamLiveInputResponse
	if err := json.Unmarshal(res, &streamListResponse); err != nil {
		return StreamLiveInput{}, err
	}

	return streamListResponse.Result, nil
}

// ListStreamLiveInputVideos list videos associated with live input.
func (api *API) ListStreamLiveInputVideos(
	ctx context.Context,
	options StreamLiveInputParameters,
) ([]StreamVideo, error) {
	if options.AccountID == "" {
		return []StreamVideo{}, ErrMissingAccountID
	}

	if options.LiveInputID == "" {
		return []StreamVideo{}, ErrMissingLiveInputID
	}

	uri := fmt.Sprintf("/accounts/%s/stream/live_inputs/%s/videos", options.AccountID, options.LiveInputID)
	res, err := api.makeRequestContext(ctx, http.MethodPost, uri, options)
	if err != nil {
		return []StreamVideo{}, err
	}

	var streamListResponse StreamListResponse
	if err := json.Unmarshal(res, &streamListResponse); err != nil {
		return []StreamVideo{}, err
	}

	return streamListResponse.Result, nil
}

package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

func GetVideoDetails(video_url string) YouTubeVideoDetails {
	var data YouTubeVideoDetails
	api_url := "https://www.youtube.com/youtubei/v1/player?key=AIzaSyAO_FJ2SlqU8Q4STEHLGCilw_Y9_11qcW8"
	video_id := GetVideoId(video_url)
	body := []byte(`
	{
		"context": {
			"client": {
				"hl": "en",
				"clientName": "WEB",
				"clientVersion": "2.20210721.00.00",
				"clientFormFactor": "UNKNOWN_FORM_FACTOR",
				"clientScreen": "WATCH",
				"mainAppWebInfo": {
					"graftUrl": "/watch?v=` + video_id + `"
				}
			},
			"user": {
				"lockedSafetyMode": false
			},
			"request": {
				"useSsl": true,
				"internalExperimentFlags": [],
				"consistencyTokenJars": []
			}
		},
		"videoId": "` + video_id + `",
		"playbackContext": {
			"contentPlaybackContext": {
				"vis": 0,
				"splay": false,
				"autoCaptionsDefaultOn": false,
				"autonavState": "STATE_NONE",
				"html5Preference": "HTML5_PREF_WANTS",
				"lactMilliseconds": "-1"
			}
		},
		"racyCheckOk": false,
		"contentCheckOk": false
	}
	`)

	r, err := http.NewRequest("POST", api_url, bytes.NewBuffer(body))
	if err != nil {
		log.Println(err)
	}
	r.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		log.Println(err)
	}
	defer res.Body.Close()

	res_body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
	}

	json.Unmarshal(res_body, &data)

	return data
}

func GetVideoId(video_url string) string {
	u, _ := url.Parse(video_url)
	q := u.Query()

	return q.Get("v")
}

func GetVideoTitle(d YouTubeVideoDetails) string {
	return d.VideoDetails.Title
}

func GetScheduledTime(d YouTubeVideoDetails) time.Time {
	ts, _ := strconv.Atoi(d.PlayabilityStatus.LiveStreamability.LiveStreamabilityRenderer.OfflineSlate.LiveStreamOfflineSlateRenderer.ScheduledStartTime)

	t := time.Unix(int64(ts), 0)

	return t
}

func IsLiveStream(d YouTubeVideoDetails) bool {
	return d.VideoDetails.IsLiveContent
}

func IsUpcoming(d YouTubeVideoDetails) bool {
	return d.VideoDetails.IsUpcoming
}

func IsStarted(d YouTubeVideoDetails) bool {
	return d.PlayabilityStatus.Status == "OK"
}

type YouTubeVideoDetails struct {
	ResponseContext struct {
		VisitorData           string `json:"visitorData"`
		ServiceTrackingParams []struct {
			Service string `json:"service"`
			Params  []struct {
				Key   string `json:"key"`
				Value string `json:"value"`
			} `json:"params"`
		} `json:"serviceTrackingParams"`
		MaxAgeSeconds             int `json:"maxAgeSeconds"`
		MainAppWebResponseContext struct {
			LoggedOut     bool   `json:"loggedOut"`
			TrackingParam string `json:"trackingParam"`
		} `json:"mainAppWebResponseContext"`
		WebResponseContextExtensionData struct {
			HasDecorated bool `json:"hasDecorated"`
		} `json:"webResponseContextExtensionData"`
	} `json:"responseContext"`
	PlayabilityStatus struct {
		Status            string `json:"status"`
		Reason            string `json:"reason"`
		PlayableInEmbed   bool   `json:"playableInEmbed"`
		LiveStreamability struct {
			LiveStreamabilityRenderer struct {
				VideoID      string `json:"videoId"`
				OfflineSlate struct {
					LiveStreamOfflineSlateRenderer struct {
						ScheduledStartTime string `json:"scheduledStartTime"`
						MainText           struct {
							Runs []struct {
								Text string `json:"text"`
							} `json:"runs"`
						} `json:"mainText"`
						SubtitleText struct {
							SimpleText string `json:"simpleText"`
						} `json:"subtitleText"`
						Thumbnail struct {
							Thumbnails []struct {
								URL    string `json:"url"`
								Width  int    `json:"width"`
								Height int    `json:"height"`
							} `json:"thumbnails"`
						} `json:"thumbnail"`
						OfflineSlateStyle string `json:"offlineSlateStyle"`
					} `json:"liveStreamOfflineSlateRenderer"`
				} `json:"offlineSlate"`
				PollDelayMs string `json:"pollDelayMs"`
			} `json:"liveStreamabilityRenderer"`
		} `json:"liveStreamability"`
		Miniplayer struct {
			MiniplayerRenderer struct {
				PlaybackMode string `json:"playbackMode"`
			} `json:"miniplayerRenderer"`
		} `json:"miniplayer"`
		ContextParams string `json:"contextParams"`
	} `json:"playabilityStatus"`
	HeartbeatParams struct {
		SoftFailOnError            bool   `json:"softFailOnError"`
		HeartbeatServerData        string `json:"heartbeatServerData"`
		HeartbeatAttestationConfig struct {
			RequiresAttestation bool `json:"requiresAttestation"`
		} `json:"heartbeatAttestationConfig"`
	} `json:"heartbeatParams"`
	VideoDetails struct {
		VideoID          string   `json:"videoId"`
		Title            string   `json:"title"`
		LengthSeconds    string   `json:"lengthSeconds"`
		Keywords         []string `json:"keywords"`
		ChannelID        string   `json:"channelId"`
		IsOwnerViewing   bool     `json:"isOwnerViewing"`
		ShortDescription string   `json:"shortDescription"`
		IsCrawlable      bool     `json:"isCrawlable"`
		Thumbnail        struct {
			Thumbnails []struct {
				URL    string `json:"url"`
				Width  int    `json:"width"`
				Height int    `json:"height"`
			} `json:"thumbnails"`
		} `json:"thumbnail"`
		IsUpcoming             bool   `json:"isUpcoming"`
		AllowRatings           bool   `json:"allowRatings"`
		ViewCount              string `json:"viewCount"`
		Author                 string `json:"author"`
		IsLowLatencyLiveStream bool   `json:"isLowLatencyLiveStream"`
		IsPrivate              bool   `json:"isPrivate"`
		IsUnpluggedCorpus      bool   `json:"isUnpluggedCorpus"`
		LatencyClass           string `json:"latencyClass"`
		IsLiveContent          bool   `json:"isLiveContent"`
	} `json:"videoDetails"`
	Microformat struct {
		PlayerMicroformatRenderer struct {
			Thumbnail struct {
				Thumbnails []struct {
					URL    string `json:"url"`
					Width  int    `json:"width"`
					Height int    `json:"height"`
				} `json:"thumbnails"`
			} `json:"thumbnail"`
			Embed struct {
				IframeURL string `json:"iframeUrl"`
				Width     int    `json:"width"`
				Height    int    `json:"height"`
			} `json:"embed"`
			Title struct {
				SimpleText string `json:"simpleText"`
			} `json:"title"`
			Description struct {
				SimpleText string `json:"simpleText"`
			} `json:"description"`
			LengthSeconds        string   `json:"lengthSeconds"`
			OwnerProfileURL      string   `json:"ownerProfileUrl"`
			ExternalChannelID    string   `json:"externalChannelId"`
			IsFamilySafe         bool     `json:"isFamilySafe"`
			AvailableCountries   []string `json:"availableCountries"`
			IsUnlisted           bool     `json:"isUnlisted"`
			HasYpcMetadata       bool     `json:"hasYpcMetadata"`
			ViewCount            string   `json:"viewCount"`
			Category             string   `json:"category"`
			PublishDate          string   `json:"publishDate"`
			OwnerChannelName     string   `json:"ownerChannelName"`
			LiveBroadcastDetails struct {
				IsLiveNow      bool   `json:"isLiveNow"`
				StartTimestamp string `json:"startTimestamp"`
			} `json:"liveBroadcastDetails"`
			UploadDate string `json:"uploadDate"`
		} `json:"playerMicroformatRenderer"`
	} `json:"microformat"`
	TrackingParams string `json:"trackingParams"`
	Attestation    struct {
		PlayerAttestationRenderer struct {
			Challenge    string `json:"challenge"`
			BotguardData struct {
				Program            string `json:"program"`
				InterpreterSafeURL struct {
					PrivateDoNotAccessOrElseTrustedResourceURLWrappedValue string `json:"privateDoNotAccessOrElseTrustedResourceUrlWrappedValue"`
				} `json:"interpreterSafeUrl"`
				ServerEnvironment int `json:"serverEnvironment"`
			} `json:"botguardData"`
		} `json:"playerAttestationRenderer"`
	} `json:"attestation"`
	Messages []struct {
		MealbarPromoRenderer struct {
			Icon struct {
				Thumbnails []struct {
					URL    string `json:"url"`
					Width  int    `json:"width"`
					Height int    `json:"height"`
				} `json:"thumbnails"`
			} `json:"icon"`
			MessageTexts []struct {
				Runs []struct {
					Text string `json:"text"`
				} `json:"runs"`
			} `json:"messageTexts"`
			ActionButton struct {
				ButtonRenderer struct {
					Style string `json:"style"`
					Size  string `json:"size"`
					Text  struct {
						Runs []struct {
							Text string `json:"text"`
						} `json:"runs"`
					} `json:"text"`
					TrackingParams string `json:"trackingParams"`
					Command        struct {
						ClickTrackingParams    string `json:"clickTrackingParams"`
						CommandExecutorCommand struct {
							Commands []struct {
								ClickTrackingParams string `json:"clickTrackingParams,omitempty"`
								CommandMetadata     struct {
									WebCommandMetadata struct {
										URL         string `json:"url"`
										WebPageType string `json:"webPageType"`
										RootVe      int    `json:"rootVe"`
										APIURL      string `json:"apiUrl"`
									} `json:"webCommandMetadata"`
								} `json:"commandMetadata"`
								BrowseEndpoint struct {
									BrowseID string `json:"browseId"`
									Params   string `json:"params"`
								} `json:"browseEndpoint,omitempty"`
								FeedbackEndpoint struct {
									FeedbackToken string `json:"feedbackToken"`
									UIActions     struct {
										HideEnclosingContainer bool `json:"hideEnclosingContainer"`
									} `json:"uiActions"`
								} `json:"feedbackEndpoint,omitempty"`
							} `json:"commands"`
						} `json:"commandExecutorCommand"`
					} `json:"command"`
				} `json:"buttonRenderer"`
			} `json:"actionButton"`
			DismissButton struct {
				ButtonRenderer struct {
					Style string `json:"style"`
					Size  string `json:"size"`
					Text  struct {
						Runs []struct {
							Text string `json:"text"`
						} `json:"runs"`
					} `json:"text"`
					TrackingParams string `json:"trackingParams"`
					Command        struct {
						ClickTrackingParams    string `json:"clickTrackingParams"`
						CommandExecutorCommand struct {
							Commands []struct {
								ClickTrackingParams string `json:"clickTrackingParams"`
								CommandMetadata     struct {
									WebCommandMetadata struct {
										SendPost bool   `json:"sendPost"`
										APIURL   string `json:"apiUrl"`
									} `json:"webCommandMetadata"`
								} `json:"commandMetadata"`
								FeedbackEndpoint struct {
									FeedbackToken string `json:"feedbackToken"`
									UIActions     struct {
										HideEnclosingContainer bool `json:"hideEnclosingContainer"`
									} `json:"uiActions"`
								} `json:"feedbackEndpoint"`
							} `json:"commands"`
						} `json:"commandExecutorCommand"`
					} `json:"command"`
				} `json:"buttonRenderer"`
			} `json:"dismissButton"`
			TriggerCondition    string `json:"triggerCondition"`
			Style               string `json:"style"`
			TrackingParams      string `json:"trackingParams"`
			ImpressionEndpoints []struct {
				ClickTrackingParams string `json:"clickTrackingParams"`
				CommandMetadata     struct {
					WebCommandMetadata struct {
						SendPost bool   `json:"sendPost"`
						APIURL   string `json:"apiUrl"`
					} `json:"webCommandMetadata"`
				} `json:"commandMetadata"`
				FeedbackEndpoint struct {
					FeedbackToken string `json:"feedbackToken"`
					UIActions     struct {
						HideEnclosingContainer bool `json:"hideEnclosingContainer"`
					} `json:"uiActions"`
				} `json:"feedbackEndpoint"`
			} `json:"impressionEndpoints"`
			IsVisible    bool `json:"isVisible"`
			MessageTitle struct {
				Runs []struct {
					Text string `json:"text"`
				} `json:"runs"`
			} `json:"messageTitle"`
			EnableSharedFeatureForImpressionHandling bool `json:"enableSharedFeatureForImpressionHandling"`
		} `json:"mealbarPromoRenderer"`
	} `json:"messages"`
	AdBreakHeartbeatParams string `json:"adBreakHeartbeatParams"`
	FrameworkUpdates       struct {
		EntityBatchUpdate struct {
			Mutations []struct {
				EntityKey string `json:"entityKey"`
				Type      string `json:"type"`
				Payload   struct {
					OfflineabilityEntity struct {
						Key                     string `json:"key"`
						AddToOfflineButtonState string `json:"addToOfflineButtonState"`
					} `json:"offlineabilityEntity"`
				} `json:"payload"`
			} `json:"mutations"`
			Timestamp struct {
				Seconds string `json:"seconds"`
				Nanos   int    `json:"nanos"`
			} `json:"timestamp"`
		} `json:"entityBatchUpdate"`
	} `json:"frameworkUpdates"`
}

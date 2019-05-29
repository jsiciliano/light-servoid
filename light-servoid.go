package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Microsoft/ApplicationInsights-Go/appinsights"
)

type Response struct {
	Now          time.Time `json:"now"`
	IsTestStream bool      `json:"isTestStream"`
	StartedAt    time.Time `json:"startedAt"`
	AccessKey    string    `json:"accessKey"`
	HlsSrc       string    `json:"hlsSrc"`
	FtlSrc       string    `json:"ftlSrc"`
}
type ChannelUserRelationship struct {
	ID     int `json:"id"`
	Status struct {
		Roles   []string `json:"roles"`
		Follows struct {
			User      int       `json:"user"`
			Channel   int       `json:"channel"`
			CreatedAt time.Time `json:"createdAt"`
		} `json:"follows"`
	} `json:"status"`
}

type lsAPI struct {
	ID          int  `json:"id"`
	IsLSEnabled bool `json:"isLSEnabled"`
	ChannelID   int  `json:"channelId"`
}

type ChannelAPI struct {
	Featured             bool        `json:"featured"`
	ID                   int         `json:"id"`
	UserID               int         `json:"userId"`
	Token                string      `json:"token"`
	Online               bool        `json:"online"`
	FeatureLevel         int         `json:"featureLevel"`
	Partnered            bool        `json:"partnered"`
	TranscodingProfileID int         `json:"transcodingProfileId"`
	Suspended            bool        `json:"suspended"`
	Name                 string      `json:"name"`
	Audience             string      `json:"audience"`
	ViewersTotal         int         `json:"viewersTotal"`
	ViewersCurrent       int         `json:"viewersCurrent"`
	NumFollowers         int         `json:"numFollowers"`
	Description          interface{} `json:"description"`
	TypeID               interface{} `json:"typeId"`
	Interactive          bool        `json:"interactive"`
	InteractiveGameID    interface{} `json:"interactiveGameId"`
	Ftl                  int         `json:"ftl"`
	HasVod               bool        `json:"hasVod"`
	LanguageID           interface{} `json:"languageId"`
	CoverID              interface{} `json:"coverId"`
	ThumbnailID          interface{} `json:"thumbnailId"`
	BadgeID              interface{} `json:"badgeId"`
	BannerURL            interface{} `json:"bannerUrl"`
	HosteeID             interface{} `json:"hosteeId"`
	HasTranscodes        bool        `json:"hasTranscodes"`
	VodsEnabled          bool        `json:"vodsEnabled"`
	CostreamID           interface{} `json:"costreamId"`
	CreatedAt            time.Time   `json:"createdAt"`
	UpdatedAt            time.Time   `json:"updatedAt"`
	DeletedAt            interface{} `json:"deletedAt"`
	Thumbnail            interface{} `json:"thumbnail"`
	Cover                interface{} `json:"cover"`
	Badge                interface{} `json:"badge"`
	Type                 interface{} `json:"type"`
	Preferences          struct {
		HypezoneAllow                      bool          `json:"hypezone:allow"`
		HostingAllow                       bool          `json:"hosting:allow"`
		HostingAllowlive                   bool          `json:"hosting:allowlive"`
		MixerFeaturedAllow                 bool          `json:"mixer:featured:allow"`
		CostreamAllow                      string        `json:"costream:allow"`
		Sharetext                          string        `json:"sharetext"`
		ChannelBannedwords                 []interface{} `json:"channel:bannedwords"`
		ChannelLinksClickable              bool          `json:"channel:links:clickable"`
		ChannelLinksAllowed                bool          `json:"channel:links:allowed"`
		ChannelSlowchat                    int           `json:"channel:slowchat"`
		ChannelNotifyDirectPurchaseMessage string        `json:"channel:notify:directPurchaseMessage"`
		ChannelNotifyDirectPurchase        bool          `json:"channel:notify:directPurchase"`
		ChannelNotifyFollow                bool          `json:"channel:notify:follow"`
		ChannelNotifyFollowmessage         string        `json:"channel:notify:followmessage"`
		ChannelNotifyHostedBy              string        `json:"channel:notify:hostedBy"`
		ChannelNotifyHosting               string        `json:"channel:notify:hosting"`
		ChannelNotifySubscribemessage      string        `json:"channel:notify:subscribemessage"`
		ChannelNotifySubscribe             bool          `json:"channel:notify:subscribe"`
		ChannelPartnerSubmail              string        `json:"channel:partner:submail"`
		ChannelPlayerMuteOwn               bool          `json:"channel:player:muteOwn"`
		ChannelTweetEnabled                bool          `json:"channel:tweet:enabled"`
		ChannelTweetBody                   string        `json:"channel:tweet:body"`
		ChannelUsersLevelRestrict          int           `json:"channel:users:levelRestrict"`
		ChannelCatbotLevel                 int           `json:"channel:catbot:level"`
		ChannelOfflineAutoplayVod          bool          `json:"channel:offline:autoplayVod"`
		ChannelChatHostswitch              bool          `json:"channel:chat:hostswitch"`
		ChannelDirectPurchaseEnabled       bool          `json:"channel:directPurchase:enabled"`
	} `json:"preferences"`

	User struct {
		Level  int `json:"level"`
		Social struct {
			Verified []interface{} `json:"verified"`
		} `json:"social"`
		ID          int         `json:"id"`
		Username    string      `json:"username"`
		Verified    bool        `json:"verified"`
		Experience  int         `json:"experience"`
		Sparks      int         `json:"sparks"`
		AvatarURL   string      `json:"avatarUrl"`
		Bio         interface{} `json:"bio"`
		PrimaryTeam interface{} `json:"primaryTeam"`
		CreatedAt   time.Time   `json:"createdAt"`
		UpdatedAt   time.Time   `json:"updatedAt"`
		DeletedAt   interface{} `json:"deletedAt"`
		Groups      []struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"groups"`
	} `json:"user"`
}

var (
	telemetryConfig = appinsights.NewTelemetryConfiguration("1fdaf5f5-63e1-4145-9f0b-c08c46dcbc1e")
	client          appinsights.TelemetryClient
	input           string
	inputTwo        string
)

func main() {
	// Configure how many items can be sent in one call to the data collector:
	telemetryConfig.MaxBatchSize = 8192

	// Configure the maximum delay before sending queued telemetry:
	telemetryConfig.MaxBatchInterval = 500 * time.Millisecond

	// Define the client
	client = appinsights.NewTelemetryClientFromConfig(telemetryConfig)
	client.Context().Tags.Cloud().SetRole("light-servoid")

	start := time.Now()

	// Application Insights single event
	client.TrackEvent("Light-Servoid: Starting")

	// Pull the command line Arguments
	args := os.Args

	//Check for invalid username by checking the length of the array.  Exit the program if there is no valid input.
	if len(args) < 1 {
		fmt.Println("Please enter a valid username.")
		os.Exit(0) //Exit the program
	} else if len(args) == 2 {
		//Run getChannelID for singular channel ID entered by the user
		getChannelID(args[1])
	} else if len(args) == 3 {
		input = args[1]
		inputTwo = args[2]
		//Run the method to get the channel relationship info between two channels.
		channelRelationship(input, inputTwo)
	}

	client.TrackEvent("Light-Servoid: Completed")
	delta := time.Now().Sub(start)
	//fmt.Println("Completed in " + delta.String())

	request := appinsights.NewRequestTelemetry("GET", "https://mixer-servoid/api/v1/foo/bar", delta, "200")
	request.MarkTime(start, time.Now())
	client.Track(request)

	client.Channel().Flush()

	time.Sleep(1 * time.Second)
}

func getChannelID(channel string) {
	client.TrackEvent("Light-Servoid: Getting channel ID")
	//fmt.Println("Getting ChannelInfo for", channel)

	start := time.Now()
	resp, err := http.Get("https://mixer.com/api/v1/channels/" + channel)
	delta := time.Now().Sub(start)

	var dependency *appinsights.RemoteDependencyTelemetry
	success := true
	if err != nil {
		fmt.Println(err)
		success = false
	}

	dependency = appinsights.NewRemoteDependencyTelemetry("api/v1/channels/{id}", "HTTP GET", "Backend", success /* success */)

	dependency.Duration = delta
	dependency.Data = "api/v1/channels/" + channel
	dependency.ResultCode = strconv.Itoa(resp.StatusCode)
	dependency.Properties["cv"] = "12346"

	if resp.StatusCode != 200 {
		trace := appinsights.NewTraceTelemetry("Light-Servoid: Backend Call failed getting channel info", appinsights.Information)
		client.Track(trace)
		fmt.Println(channel + " Not found")
		return
	} else {
		trace := appinsights.NewTraceTelemetry("Light-Servoid: Backend Call successful getting channel info", appinsights.Information)
		client.Track(trace)
	}

	var channelObj ChannelAPI
	err = json.NewDecoder(resp.Body).Decode(&channelObj)

	if channelObj.Token != "" {
		fmt.Println("Looking up " + input + " with ID " + strconv.Itoa(channelObj.ID))
		getM3u8(strconv.Itoa(channelObj.ID))
		fmt.Println("Current Time UTC):", time.Now().UTC())
		if channelObj.HosteeID != nil {
			fmt.Println("Hosting:", strconv.FormatFloat(channelObj.HosteeID.(float64), 'f', -1, 64))
		} else {
			fmt.Println("Hosting: Nobody")
		}

		if len(channelObj.User.Groups) == 1 {
			fmt.Println("Pro User: false")
		} else if channelObj.User.Groups[1].Name == "Staff" {
			fmt.Println("Pro: Staff")
		} else if channelObj.User.Groups[1].Name == "Pro" || channelObj.User.Groups[2].Name == "Pro" {
			fmt.Println("Pro User: true")
		}

		getLS(strconv.Itoa(channelObj.ID))
		fmt.Println("VODs Enabled:", channelObj.VodsEnabled)
		fmt.Println("Users Current Sparks:", channelObj.User.Sparks)
		vlcURL(channelObj.ID)
		fmt.Println("Xpert URL: https://xpert.microsoft.com/osg/views/PROBOTv2?overrides=%7B%22Source%22%3A%22Environment%3DPROD%3BModernClient%3DPartners%3BVEFProvider%3DMixer%3BVEFProvider%3DServices%3BVEFProvider%3DChannel%3BVEFTopic%3D" + strconv.Itoa(channelObj.ID) + "%3B%22%7D")
		fmt.Println("Unstuck URL (Requires v-dash): https://mixer-unstuck-ppe.azurewebsites.net/api/unstick/channel/" + channel)
		fmt.Println("Refresh URL (Requires v-dash): https://mixer-unstuck-ppe.azurewebsites.net/api/refresh/channel" + channel)
	} else {
		fmt.Println("Failed to get ID")
		return
	}

	defer resp.Body.Close()
}

func getUserID(channel string, userID bool) int {
	//gets the user ID
	client.TrackEvent("Light-Servoid: Getting user ID")

	//start := time.Now()
	resp, err := http.Get("https://mixer.com/api/v1/channels/" + channel)

	//	success := true
	if err != nil {
		fmt.Println(err)
		//success = false
	}

	if resp.StatusCode != 200 {
		trace := appinsights.NewTraceTelemetry("Light-Servoid: Backend Call failed getting channel info", appinsights.Information)
		client.Track(trace)
		fmt.Println(channel + " Not found")
		return 0
	} else {
		trace := appinsights.NewTraceTelemetry("Light-Servoid: Backend Call successful getting channel info", appinsights.Information)
		client.Track(trace)
	}

	var channelObj ChannelAPI
	err = json.NewDecoder(resp.Body).Decode(&channelObj)

	if channelObj.Token != "" {
		//just proceed if there is no error
	} else {
		fmt.Println("Failed to get user ID")
		return 0
	}

	defer resp.Body.Close()

	if userID == false {
		return channelObj.ID
	} else {
		return channelObj.UserID
	}
}

func getLS(channel string) {
	client.TrackEvent("Light-Servoid: Getting channel LightStream Status")
	//fmt.Println("Getting LightStream Status for", channel)

	start := time.Now()
	resp, err := http.Get("https://mixer.com/api/v1/channels/" + channel + "/videoSettings")
	delta := time.Now().Sub(start)

	var dependency *appinsights.RemoteDependencyTelemetry
	success := true
	if err != nil {
		fmt.Println(err)
		success = false
	}

	dependency = appinsights.NewRemoteDependencyTelemetry("api/v1/channels/{id}/videoSettings", "HTTP GET", "Backend", success /* success */)

	dependency.Duration = delta
	dependency.Data = "api/v1/channels/" + channel + "/videoSettings"
	dependency.ResultCode = strconv.Itoa(resp.StatusCode)
	dependency.Properties["cv"] = "12346"

	if resp.StatusCode != 200 {
		trace := appinsights.NewTraceTelemetry("Light-Servoid: Backend Call failed getting VideoSettings", appinsights.Information)
		client.Track(trace)
		fmt.Println(channel + " 404 LS Not found")
		return
	} else {
		trace := appinsights.NewTraceTelemetry("Light-Servoid: Backend Call successful getting VideoSettings", appinsights.Information)
		client.Track(trace)
	}

	var channelObj lsAPI
	err = json.NewDecoder(resp.Body).Decode(&channelObj)

	fmt.Println("LightStream Status:", channelObj.IsLSEnabled)

	defer resp.Body.Close()
}

func getM3u8(channel string) {
	//fmt.Println("Getting M3U8 for", channel)
	client.TrackEvent("Light-Servoid: Getting channel manifest")

	start := time.Now()
	resp, err := http.Get("https://mixer.com/api/v1/channels/" + channel + "/manifest.light2")
	delta := time.Now().Sub(start)

	var dependency *appinsights.RemoteDependencyTelemetry
	success := true
	if err != nil {
		fmt.Println(err)
		success = false
	}

	dependency = appinsights.NewRemoteDependencyTelemetry("api/v1/channels/{id}/manifest.light2", "HTTP GET", "Falcon", success /* success */)

	dependency.Duration = delta
	dependency.Data = "api/v1/channels/" + channel + "/manifest.light2"
	dependency.ResultCode = strconv.Itoa(resp.StatusCode)
	dependency.Properties["cv"] = "12346"

	if resp.StatusCode != 200 {
		trace := appinsights.NewTraceTelemetry("Light-Servoid: Channel was offline", appinsights.Information)
		client.Track(trace)
		fmt.Println("channel was offline")
		return
	} else {
		trace := appinsights.NewTraceTelemetry("Light-Servoid: Channel was online", appinsights.Information)
		client.Track(trace)
	}

	var respObj Response
	err = json.NewDecoder(resp.Body).Decode(&respObj)

	// Define the Trace
	trace := appinsights.NewTraceTelemetry("message", appinsights.Information)
	trace.Properties["hlsSource"] = respObj.HlsSrc
	trace.Properties["ftlSource"] = respObj.FtlSrc
	client.Track(trace)

	getDist(respObj.AccessKey)

	fmt.Println("PROCESSED VIDEO (VLC): https://video.mixer.com/hls/" + respObj.AccessKey + "_source/index.m3u8")
	//fmt.Println("HLSSource: " + respObj.HlsSource)
	//fmt.Println("FTLSource: " + respObj.FtlSource)

	client.Track(dependency)
	defer resp.Body.Close()
}

func getDist(accessKey string) {
	client.TrackEvent("Light-Servoid: Getting Dist Server")

	url := "https://video.mixer.com/hls/" + accessKey + "_source/index.m3u8"

	start := time.Now()
	resp, err := http.Get(url)
	delta := time.Now().Sub(start)

	var dependency *appinsights.RemoteDependencyTelemetry
	success := true
	if err != nil {
		fmt.Println(err)
		success = false
	}

	dependency = appinsights.NewRemoteDependencyTelemetry("/hls/{accessKey}_source/index.m3u8", "HTTP GET", "Janus", success /* success */)

	dependency.Duration = delta
	dependency.Data = url
	dependency.ResultCode = strconv.Itoa(resp.StatusCode)

	client.Track(dependency)
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Println(resp.StatusCode)
	} else {
		for header, value := range resp.Header {
			if strings.Contains(header, "Cdn") {
				client.TrackEvent("Light-Servoid: Reporting Dist Server")
				// Define the Trace
				trace := appinsights.NewTraceTelemetry("message", appinsights.Information)
				trace.Properties["dist"] = value[0]
				client.Track(trace)
				fmt.Println("Your Dist: " + value[0])
			}
		}
	}
}

func vlcURL(channel int) {
	//Print the video to watch the streams source video for VLC.
	var channelString = strconv.Itoa(channel)
	fmt.Println("SOURCE VIDEO (VLC):  rtmp://<ingestname>.mixer.com:1935/beam/" + channelString)
}

func channelRelationship(channelOne string, channelTwo string) {
	//Get the user ID or channel ID, pass true if you want the user ID and false for channel ID, returns the ID in an INT
	userIDNum := getUserID(channelTwo, true)
	channelIDNum := getUserID(channelOne, false)

	//Convert channel/user IDs to Strings to concatenate onto the URL
	userIDString := strconv.Itoa(userIDNum)
	channelIDString := strconv.Itoa(channelIDNum)

	//Show the user what we're doing
	fmt.Println("Looking up the channel user relationship with " + "Channel: " + channelOne + "  User: " + channelTwo)

	//Generate the correct json url with the channel id and user id strings
	url := "https://mixer.com/api/v1/channels/" + channelIDString + "/relationship?user=" + userIDString
	fmt.Println("URL:  " + url)

	//Parse the json file into the struct (check for) & print the roles
	resp, err := http.Get(url)

	if err != nil {
		fmt.Println(err)
	} else {
		var userObj ChannelUserRelationship
		err = json.NewDecoder(resp.Body).Decode(&userObj)
		fmt.Print("User Roles: ")
		fmt.Println(userObj.Status.Roles)
	}
	defer resp.Body.Close()
}

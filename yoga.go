package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/skratchdot/open-golang/open"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

var Browser = "Brave'"

func main() {
	openVideoInBrowser()   // From this month playlist play video which count equals to today's date
	checkBrowserLocation() // When in right position exec moveBrowser
	changeAudioOutput()    // Check audio output, if not HDMI3 then switch to it
	pressKeys()            // Press the f and c for full screen and disable captions
}

func openVideoInBrowser() {
	ctx := context.Background()

	// Read the client secret file
	secretFile, err := os.ReadFile("alpine-realm-381711-182fcef9362c.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// Parse client secret file to config
	config, err := google.JWTConfigFromJSON(secretFile, youtube.YoutubeReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	// Create HTTP client
	httpClient := config.Client(ctx)

	// Create YouTube service client
	youtubeService, err := youtube.NewService(ctx, option.WithHTTPClient(httpClient))
	if err != nil {
		log.Fatalf("Error creating YouTube client: %v", err)
	}

	// Fetch playlists from YouTube channel
	playlistsCall := youtubeService.Playlists.List([]string{"snippet"}).ChannelId("UCFKE7WVJfvaHW5q283SxchA").MaxResults(50)
	playlistsResponse, err := playlistsCall.Do()
	if err != nil {
		log.Fatalf("Error fetching playlists: %v", err)
	}

	var selectedPlaylist *youtube.Playlist
	
	// Find playlist for the current month
	for _, playlist := range playlistsResponse.Items {
		if month := time.Now().Month().String(); strings.Contains(playlist.Snippet.Title, month) {
			selectedPlaylist = playlist
			break
		} else {
			// get the last playlist
			selectedPlaylist = playlistsResponse.Items[0]
			
		}
		fmt.Println(selectedPlaylist)
	}

	// Check if no playlist found for the current month
	if selectedPlaylist == nil {
		log.Fatalf("No playlist found for the current month")
	}

	// Fetch playlist items
	playlistItemsCall := youtubeService.PlaylistItems.List([]string{"snippet"}).PlaylistId(selectedPlaylist.Id).MaxResults(50)
	playlistItemsResponse, err := playlistItemsCall.Do()
	if err != nil {
		log.Fatalf("Error fetching playlist items: %v", err)
	}

	// Calculate index of the video for the current day
	videoIndex := (time.Now().Day()) % len(playlistItemsResponse.Items)-1
	videoId := playlistItemsResponse.Items[videoIndex].Snippet.ResourceId.VideoId
	videoURL := fmt.Sprintf("https://www.youtube.com/embed/%s"+"?autoplay=1", videoId)

	// Open the video URL in a web browser
	err = open.Run(videoURL)
	if err != nil {
		log.Fatalf("Error opening URL: %v", err)
	}
}

func checkBrowserLocation() { // Execute xdotool command to check browser window location
	time.Sleep(2 * time.Second)

	cmdToEnter := "xdotool search --onlyvisible --class " + Browser + " getwindowgeometry --shell | grep -oP 'X=\\K\\d+'"
	cmd, err := exec.Command("sh", "-c", cmdToEnter).CombinedOutput()
	xPos := string(cmd)
	log.Print("xPosition is:", xPos)
	if err != nil && xPos >= "0" && xPos < "2400" {
		log.Println("Browser window is on the right screen or not found")
		return
	}
	log.Println("Browser window is on the left screen")
	moveBrowser()
}

func moveBrowser() { // Move the browser window to the right screen
	moveWithCmd := ("xdotool search --onlyvisible --class " + Browser + " windowmove --relative -- 1920 0")
	cmd, err := exec.Command(moveWithCmd).CombinedOutput()
	if err != nil {
		log.Printf("Error moving browser window: %s\n%s", err, cmd)
	} else {
		log.Println("Browser window moved successfully")
	}
}

// to check current audio profile run command below
// wpctl status | grep "HDMI 3"
func changeAudioOutput() { // Check the current audio device and set the default sink

	currentDevice, err := exec.Command("wpctl", "status", "|", "grep", "'HDMI 3'").CombinedOutput()
	if err != nil && len(currentDevice) != 0 {
		log.Printf("Error finding current audio device: %s\n%s", err, currentDevice)
	}
	// Set the default sound profile to HDMI 3
	output, err := exec.Command("wpctl", "set-profile", "40", "3").CombinedOutput()
	if err != nil {
		log.Fatalf("Error setting default sink: %s\n%s", err, output)
	}
	log.Println("Default sound profile set to HDMI 3")
}

func pressKeys() {
	time.Sleep(3 * time.Second)

	// Simulate key presses in the Brave browser window

	cmd := exec.Command("xdotool", "search", "--onlyvisible", "--class", "brave", "key", "--window", "%1", "f", "c")
	err := cmd.Run()
	if err != nil {
		log.Fatalf("Error pressing keys: %v", err)
	}
}

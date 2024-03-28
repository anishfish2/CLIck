package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/joho/godotenv"
	"google.golang.org/api/googleapi/transport"
	"google.golang.org/api/youtube/v3"
)

func main() {

	//Define Flags
	query := flag.String("query", "Code broken :( y?", "The error you are facing.")
	//answer := flag.Bool("ask", true, "Whether to ask ChatGPT for help.")
	video := flag.Bool("video", false, "Whether to display the video")
	video_topic := ""

	//Get + Parse Flags & Environment variables
	flag.Parse()
	load_env := godotenv.Load(".env")

	if load_env != nil {
		fmt.Println("Error loading .env file", load_env)
		return
	}
		
	if *video {

		//Set up ChatGPT API
		url := "https://api.openai.com/v1/chat/completions"

		//Define ChatGPT call
		body := map[string]interface{}{
			"model": "gpt-4",
			"messages": []map[string]string{
				{
					"role":    "user",
					"content": fmt.Sprintf("I have this question: %s. Can you give a search phrase to find a YouTube video to help me fix it", *query),
				},
			},
		}

		//Send call and parse ChatGPT response
		jsonBody, _ := json.Marshal(body)
		req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
		req.Header = http.Header{
			"Content-Type":  {"application/json"},
			"Authorization": {"Bearer " + os.Getenv("OPENAI_KEY")},
		}
		client := &http.Client{}
		res, _ := client.Do(req)
		defer res.Body.Close()
		resBody, err := io.ReadAll(res.Body)
		if err != nil {
			fmt.Printf("Impossible to read OPENAI response: %s", err)
		}
		var responseBody map[string]interface{}
		_ = json.Unmarshal(resBody, &responseBody)

		choices := responseBody["choices"].([]interface{})
		if len(choices) > 0 {
			message := choices[0].(map[string]interface{})["message"].(map[string]interface{})
			content := message["content"].(string)
			video_topic = content
		}

		fmt.Printf("Search Query: %s", video_topic)



		// Make call to Youtube API to get search URL
		var developerKey = os.Getenv("YOUTUBE_KEY")
		client2 := &http.Client{
			Transport: &transport.APIKey{Key: developerKey},
		}
		service, err := youtube.New(client2)

		if err != nil {
			log.Fatalf("Error creating new YouTube client: %v", err)
		}
		call := service.Search.List([]string{}).Q(video_topic).MaxResults(1)
		response, err := call.Do()

		if err != nil {
			log.Fatalf("Error making youtube search API call: %v", err)
		}

		video_id := response.Items[0].Id.VideoId

		//Define video URL
		videoURL := "https://www.youtube.com/watch?v=" + video_id

		//Play Video
		cmd := exec.Command("mpv", videoURL)
		cmd.Run()

	}
}
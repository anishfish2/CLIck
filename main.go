package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"io"

	"github.com/joho/godotenv"
	//"github.com/gen2brain/go-mpv"
)

func main() {

	query := flag.String("query", "Code broken :( y?", "The error you are facing.")
	ask := flag.Bool("ask", false, "Whether to ask ChatGPT for help.")
	video_topic := ""

	flag.Parse()


		fmt.Println("The question is: " + *query)

		load_env := godotenv.Load(".env")

		if load_env != nil {
			fmt.Println("Error loading .env file", load_env)
			return
		}
		
		if *ask {
			
			url := "https://api.openai.com/v1/chat/completions"

			body := map[string]interface{}{
				"model": "gpt-4",
				"messages": []map[string]string{
					{
						"role":    "user",
						"content": fmt.Sprintf("I have this question: %s. Can you give a search phrase to find a YouTube video to help me fix it", *query),
					},
				},
			}

			jsonBody, err := json.Marshal(body)
			if err != nil {
				fmt.Println("Error encoding JSON:", err)
				return
			}

			req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))

			req.Header = http.Header{
				"Content-Type": {"application/json"},
				"Authorization": {"Bearer " + os.Getenv("OPENAI_KEY")},
			}

			client := &http.Client{}

			res, err := client.Do(req)

			if err != nil {
				panic(err)
			}

			defer res.Body.Close()
			
			resBody, err := io.ReadAll(res.Body)


			if err != nil {
				fmt.Printf("impossible to read response: %s", err)
			}

			var responseBody map[string]interface{}
			err = json.Unmarshal(resBody, &responseBody)
			if err != nil {
				fmt.Printf("Error decoding JSON response: %s", err)
				return
			}

			choices := responseBody["choices"].([]interface{})
			if len(choices) > 0 {
				message := choices[0].(map[string]interface{})["message"].(map[string]interface{})
				content := message["content"].(string)
				video_topic = content
			}

			fmt.Printf("res body: %s", video_topic)

	}
}
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
)

type Post struct {
	Id     int    `json:"id"`
	Object  string `json:"object"`
	Created   int `json:"created"`
	Model int    `json:"model"`
	SystemFingerprint string `json:"system_fingerprint"`
	Choices []string `json:"choices"`
	Usage []string `json:"usage"`
}

func main() {

	query := flag.String("query", "Code broken :( y?", "The error you are facing.")
	ask := flag.Bool("ask", false, "Whether to ask ChatGPT for help.")

	flag.Parse()

	if flag.NArg() == 0 {
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
				fmt.Printf("impossible to read all body of response: %s", err)
			}
			fmt.Printf("res body: %s", string(resBody))
		}
	}
}
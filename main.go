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
	"github.com/gen2brain/go-mpv"
)

func main() {

	query := flag.String("query", "Code broken :( y?", "The error you are facing.")
	ask := flag.Bool("ask", false, "Whether to ask ChatGPT for help.")
	video_url := flag.String("video_url", "", "The URL of the video that was suggested by ChatGPT.")

	flag.Parse()

	if flag.NArg() == 0 {
		fmt.Println("The question is: " + *query)

		load_env := godotenv.Load(".env")

		if load_env != nil {
			fmt.Println("Error loading .env file", load_env)
			return
		}
		

		m := mpv.New()
		defer m.TerminateDestroy()

		_ = m.RequestLogMessages("info")
		_ = m.ObserveProperty(0, "pause", mpv.FormatFlag)

		_ = m.SetPropertyString("input-default-bindings", "yes")
		_ = m.SetOptionString("input-vo-keyboard", "yes")
		_ = m.SetOption("osc", mpv.FormatFlag, true)

		err := m.Initialize()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		err = m.Command([]string{"loadfile", os.Args[1]})
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

	loop:
		for {
			e := m.WaitEvent(10000)

			switch e.EventID {
			case mpv.EventPropertyChange:
				prop := e.Property()
				value := prop.Data.(int)
				fmt.Println("property:", prop.Name, value)
			case mpv.EventFileLoaded:
				p, err := m.GetProperty("media-title", mpv.FormatString)
				if err != nil {
					fmt.Println("error:", err)
				}
				fmt.Println("title:", p.(string))
			case mpv.EventLogMsg:
				msg := e.LogMessage()
				fmt.Println("message:", msg.Text)
			case mpv.EventStart:
				sf := e.StartFile()
				fmt.Println("start:", sf.EntryID)
			case mpv.EventEnd:
				ef := e.EndFile()
				fmt.Println("end:", ef.EntryID, ef.Reason)
				if ef.Reason == mpv.EndFileEOF {
					break loop
				} else if ef.Reason == mpv.EndFileError {
					fmt.Println("error:", ef.Error)
				}
			case mpv.EventShutdown:
				fmt.Println("shutdown:", e.EventID)
				break loop
			default:
				fmt.Println("event:", e.EventID)
			}

			if e.Error != nil {
				fmt.Println("error:", e.Error)
			}
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
			fmt.Printf("res body: %s", string(resBody))
		}
	}
}
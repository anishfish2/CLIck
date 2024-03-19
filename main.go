package main

import ( 
	"fmt"
	"flag"
	"os"

	"github.com/joho/godotenv"
)

func main() {

	query := flag.String("query", "Code broken :( y?", "The error you are facing.")

	flag.Parse()

	if flag.NArg() == 0 {
		fmt.Println("The question is: " + *query)

		load_env := godotenv.Load(".env")
		
		if load_env != nil {
			fmt.Println("Error loading .env file", load_env)
			return
		}

		fmt.Println("Your API KEY is:" + os.Getenv("OPENAI_KEY"))
	}



}
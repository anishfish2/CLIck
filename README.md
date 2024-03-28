# YouTube Helper with ChatGPT

This is a simple tool that helps you find YouTube videos based on your questions or errors using ChatGPT and the YouTube Data API.

## Requirements

- Go programming language installed.
- `mpv` installed and added to the system PATH. You can install it from [mpv.io](https://mpv.io/).
- YouTube Data API key and OpenAI developer key. Place these keys in a `.env` file in the root directory of the project.

## Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/yourusername/youtube-helper.git
   ```

2. Navigate to the project directory:

   ```bash
   cd youtube-helper
   ```

3. Install dependencies:

   ```bash
   go mod tidy
   ```

## Usage

### Running the program

You can run the program with the following command:

```bash
go run main.go
```

### Flags

The program supports the following flags:

- `-query`: Specify the error or question you're facing. Default is "Code broken :(".
- `-ask`: Whether to ask ChatGPT for help. Default is true.
- `-video`: Whether to display the video. Default is false.

Example usage:

```bash
go run main.go -query="How do I fix a segmentation fault?" -video=true
```

## How it works

1. If the `-ask` flag is set to true, the program uses ChatGPT to generate a search query based on your error or question.
2. If the `-video` flag is set to true, the program then uses the YouTube Data API to search for a video related to the generated search query.
3. It retrieves the video URL and plays the video using `mpv`.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

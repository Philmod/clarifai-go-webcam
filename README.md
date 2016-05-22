# clarifai-go-webcam
Web page that captures images from your webcam and tag them using Clarifai API. If one of the specified tag is detected, the image is shown.

![Example](https://dl.dropboxusercontent.com/u/45971143/clarifai-go-webcam.png)

[Live Demo](https://clarifai-go-webcam.herokuapp.com/)

## Install
- Create a [Clarifai](https://developer.clarifai.com/account/applications/) app
- Set `CLARIFAI_ID` and `CLARIFAI_SECRET` environment variables
- `go run main.go`

## Todo
- Tests
- Replace with the official Clarifai client when the [PR](https://github.com/Clarifai/clarifai-go/pull/6) is merged
- Reconnect websocket automatically

# Getting Started

1. Clone the repository
```
git clone https://github.com/loktev-d/echelon-test
 ```
2. Go to repo folder
```
cd ./echelon-test
```
3. Get your API key from API Console (make sure you have YouTube API enabled) and set it in `./config.yml`
4. Start gRPC server
```
go run ./cmd/server/main.go
```
5. Run CLI util
```
go run ./cmd/cli/main.go "https://www.youtube.com/watch?v=dQw4w9WgXcQ&ab_channel=RickAstley"
```
6. Done!

# Usage examples
- Get a thumbnail
```
go run ./cmd/cli/main.go "https://www.youtube.com/watch?v=djV11Xbc914&ab_channel=a-ha"
```
- Get multiple thumbnails
```
go run ./cmd/cli/main.go "https://www.youtube.com/watch?v=dQw4w9WgXcQ&ab_channel=RickAstley" "https://www.youtube.com/watch?v=djV11Xbc914&ab_channel=a-ha"
```
- Get multiple thumbnails asyncronously
```
go run ./cmd/cli/main.go --async "https://www.youtube.com/watch?v=dQw4w9WgXcQ&ab_channel=RickAstley" "https://www.youtube.com/watch?v=djV11Xbc914&ab_channel=a-ha"
```

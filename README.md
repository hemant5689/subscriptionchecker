# Social Media Subscription Checker

A Go API that verifies if users are subscribed to or following various social media platforms.

## Supported Platforms

- YouTube (subscribers)
- Facebook (followers)
- Instagram (followers)
- Discord (server membership)
- Twitter (followers)

## Project Structure

```
subscriptionchecker/
│
├── cmd/
│   └── server/
│       └── main.go              # Entry point
│
├── internal/
│   ├── api/                     # All platform-specific handlers
│   │   ├── facebook.go
│   │   ├── youtube.go
│   │   ├── instagram.go
│   │   ├── twitter.go
│   │   └── discord.go
│   │
│   ├── platform/                # Platform logic (business layer)
│   │   ├── facebook.go
│   │   ├── youtube.go
│   │   ├── instagram.go
│   │   ├── twitter.go
│   │   ├── discord.go
│   │   └── platform.go          # Common interfaces
│   │
│   ├── auth/                    # OAuth / token logic
│   │   └── auth.go
│   │
│   ├── config/                  # Loads ENV/config
│   │   └── config.go
│   │
│   ├── router/                  # Sets up routes
│   │   └── router.go
│   │
│   └── server/                  # Server setup
│       └── server.go
│
├── pkg/
│   └── utils/                   # Optional helpers
│       └── http.go              # HTTP response helpers
│
├── go.mod
├── go.sum
├── .env
└── README.md
```

## API Endpoints

### YouTube

- `GET /auth/youtube/login` - Initiate YouTube OAuth login
- `GET /auth/youtube/callback` - OAuth callback
- `GET /check-youtube-subscription?token=TOKEN` - Check if user is subscribed

### Facebook

- `GET /auth/facebook/login` - Initiate Facebook OAuth login
- `GET /auth/facebook/callback` - OAuth callback
- `GET /check-facebook-follower?token=TOKEN&targetId=TARGET_ID` - Check if user follows the profile

### Instagram

- `GET /auth/instagram/login` - Initiate Instagram OAuth login
- `GET /auth/instagram/callback` - OAuth callback
- `GET /check-instagram-follower?token=TOKEN&username=USERNAME` - Check if user follows the profile

### Discord

- `GET /auth/discord/login` - Initiate Discord OAuth login
- `GET /auth/discord/callback` - OAuth callback
- `GET /check-discord-server?token=TOKEN` - Check if user is in the server

### Twitter

- `GET /auth/twitter/login` - Initiate Twitter OAuth login
- `GET /auth/twitter/callback` - OAuth callback
- `GET /check-twitter-follower?token=TOKEN&username=USERNAME` - Check if user follows the profile

## Setup

1. Clone the repository
2. Create a `.env` file with your API keys and secrets
3. Run the application with `go run cmd/server/main.go`

## Environment Variables

```
# YouTube
YT_CLIENT_ID=your_youtube_client_id
YT_CLIENT_SECRET=your_youtube_client_secret
YT_REDIRECT_URL=your_youtube_redirect_url
YT_CHANNEL_ID=your_youtube_channel_id

# Facebook/Instagram (Meta)
META_APP_ID=your_meta_app_id
META_APP_SECRET=your_meta_app_secret
META_REDIRECT_URI=your_meta_redirect_uri

# Discord
DISCORD_CLIENT_ID=your_discord_client_id
DISCORD_CLIENT_SECRET=your_discord_client_secret
DISCORD_REDIRECT_URI=your_discord_redirect_uri
DISCORD_SERVER_ID=your_discord_server_id

# Twitter
TWITTER_CLIENT_ID=your_twitter_client_id
TWITTER_CLIENT_SECRET=your_twitter_client_secret
TWITTER_REDIRECT_URI=your_twitter_redirect_uri

# Server
PORT=8080
```

## License

MIT 
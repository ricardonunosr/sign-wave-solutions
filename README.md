# Sign Wave Solutions

# Tech Stack
- [Go 1.22.5](https://go.dev)
- [Chi Router](https://github.com/go-chi/chi)
- [htmx](https://htmx.org)
- [TailwindCSS](https://tailwindcss.com)
- [Gmail SMTP (Free)]()

## Local Development

```bash
# This will run watch if files are changed and build and restart the server
# You need air for local development:
#      - go install github.com/cosmtrek/air@latest
#      - install tailwindcss CLI
./build.sh run
```

# Dockerfile
```bash
docker build . -t sign-wave
docker run -p 3005:3005 --env-file .env --rm --name sign-wave-app sign-wave
```

# `.env` example

```.env
PORT=80
EMAIL_ADDRESSES=something@gmail.com
```
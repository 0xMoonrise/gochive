# Gochive

This is a personal project to store, back up, and centralize PDF documents, books, and articles. The main goal of this project is to practice and improve my development skills.

Create a docker image from Dockerfile
```sh
docker build -t gochive:1.0 .
```

Build docker container 
```
docker run --name gochive \
  -v /opt/gochive/:/opt/gochive/ \
  --env-file=.env \
  -p 8080:8080 \
  --rm gochive:1.0 
```

config `.env`
```bash
HOST=0.0.0.0
MODE=1
```

S3 client environment variables required (optional)
```bash
BUCKET=
ACCESS_KEY=
SECRET_KEY=
S3_ENDPOINT=
```

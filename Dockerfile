FROM golang:1.18-alpine

WORKDIR /app

COPY ./go.mod ./
COPY ./go.sum ./
RUN go mod download

COPY . .
CMD [ "MONGO_URL=mongodb://localhost:27017", "JWT_SECRET=asd", "go", "run", "." ]
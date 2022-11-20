FROM golang:latest

WORKDIR /app

COPY ./ ./
RUN go mod download
ENV storage=""

RUN go build -o .

EXPOSE 8080

CMD [ "/app/TgMoneyTrackerBot" ]
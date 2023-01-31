package main

import (
    "context"
    "os"
    "log"
 
    "github.com/joho/godotenv"
    "github.com/slack-go/slack"
    "github.com/slack-go/slack/slackevents"
    "github.com/slack-go/slack/socketmode"
)
func main() {
 
    godotenv.Load(".env")
 
    token := os.Getenv("SLACK_AUTH_TOKEN")
    appToken := os.Getenv("SLACK_APP_TOKEN")
 
    client := slack.New(token, slack.OptionDebug(true), slack.OptionAppLevelToken(appToken))
 
    socketClient := socketmode.New(
        client,
        socketmode.OptionDebug(true),
        socketmode.OptionLog(log.New(os.Stdout, "socketmode: ", log.Lshortfile|log.LstdFlags)),
    )
 
    ctx, cancel := context.WithCancel(context.Background())
 
    defer cancel()
 
    go func(ctx context.Context, client *slack.Client, socketClient *socketmode.Client) {
        for {
            select {
            case <-ctx.Done():
                log.Println("Shutting down socketmode listener")
                return
            case event := <-socketClient.Events:
 
                switch event.Type {
        
                case socketmode.EventTypeEventsAPI:
 
                    eventsAPI, ok := event.Data.(slackevents.EventsAPIEvent)
                    if !ok {
                        log.Printf("Could not type cast the event to the EventsAPI: %v\n", event)
                        continue
                    }
 
                    socketClient.Ack(*event.Request)
                    log.Println(eventsAPI)
                }
            }
        }
    }(ctx, client, socketClient)
 
    socketClient.Run()
}
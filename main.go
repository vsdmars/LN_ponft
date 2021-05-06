package main

import (
    _ "embed"
    "context"
    "flag"
    "fmt"
    "golang.org/x/time/rate"
    "math/rand"
    "net/http"
    "os"
    "os/signal"
    "strings"
    "syscall"
    "time"
)

//go:embed request.txt
var req string
var reqs []string =  strings.Split(req, "\n");
const localURL string = "http://localhost:9000/platform-telemetry/li/ponf"
const content string = "Content-Type: application/json"

func post() {

    getBody := func() (*strings.Reader, int) {
        rand.Seed(time.Now().Unix())
        idx := rand.Intn(len(reqs))

        return strings.NewReader(reqs[idx]), idx
    }

    body, idx := getBody()

    fmt.Printf("[INFO] post body: %s\n", reqs[idx])

    // no client side timeout set since this cli coded to work with local deployment microservice.
    res, err := http.Post(localURL, content, body)
    if err != nil {
        fmt.Printf("[ERROR] local post request to %s error. content: %s body: %s\n", localURL, content, reqs[idx])
        return
    }

    defer res.Body.Close()
}

func runner(ctx context.Context, limiter *rate.Limiter) <- chan struct{} {
    c := make(chan struct{})

    go func() {
        for {
            select {
                case <-ctx.Done():
                    c <- struct{}{}
                default:
                    if err := limiter.Wait(ctx); err != nil {
                        fmt.Printf("rate limit err %s\n", err)
                    }
                    post()
            }
        }
    } ()

    return c
}

func main() {

    // parse argument
    testMins := flag.Int("minute", 1, "test period (in minutes)")
    testQPS := flag.Int("qps", 100, "QPS")
    flag.Parse()

    // setup signal
    sigs := make(chan os.Signal, 1)
    signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

    endTimer := time.NewTimer(time.Duration(*testMins) * time.Minute)
    limiter := rate.NewLimiter(rate.Every(time.Second / time.Duration(*testQPS)), 1)

    ctx, cancel := context.WithCancel(context.Background())
    go func() {
        select {
            case <- sigs:
                cancel()
            case <- endTimer.C:
                cancel()
        }
    } ()


    cleanUp := runner(ctx, limiter)
    <- cleanUp
}

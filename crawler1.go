package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	log "github.com/albrow/prtty"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/device"
	"github.com/chromedp/cdproto/network"
	//"github.com/chromedp/cdproto/emulation"
	"github.com/qor/media"

	"./config/db"
	"./model"
)

type LdJson struct {
	Name        string          `json:"name"`
	Url         string          `json:"url"`
	Image       string          `json:"image"`
	Description string          `json:"description"`
	Rating      AggregateRating `json:"aggregateRating"`
}

type AggregateRating struct {
	Value string `json:"ratingValue"`
	Count string `json:"ratingCount"`
	Best  string `json:"bestRating"`
	Worst string `json:"worstRating"`
}

// Job - interface for job processing
type Job interface {
	Process(ctx *context.Context)
}

// Worker - the worker threads that actually process the jobs
type Worker struct {
	done             sync.WaitGroup
	readyPool        chan chan Job
	assignedJobQueue chan Job
	quit             chan bool
	ctx              *context.Context
}

// JobQueue - a queue for enqueueing jobs to be processed
type JobQueue struct {
	internalQueue     chan Job
	readyPool         chan chan Job
	workers           []*Worker
	dispatcherStopped sync.WaitGroup
	workersStopped    sync.WaitGroup
	quit              chan bool
}

type Games struct {
	sync.RWMutex
	Id  uint
	Map map[string]uint
}

type GameJob struct {
	URI string
}

var (
	links    []*cdp.Node
	gamelist Games
)

func main() {
	db.DB.AutoMigrate(&model.Game{})

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(),
		[]chromedp.ExecAllocatorOption{
			chromedp.UserDataDir("chrome-data"),
			chromedp.NoFirstRun,
			chromedp.NoDefaultBrowserCheck,
			chromedp.Flag("disable-background-timer-throttling", true),
			chromedp.Flag("disable-client-side-phishing-detection", true),
			chromedp.Flag("disable-default-apps", true),
			chromedp.Flag("disable-features", "TranslateUI"),
			chromedp.Flag("disable-popup-blocking", true),
			chromedp.Flag("disable-prompt-on-repost", true),
			chromedp.Flag("disable-sync", true),
			chromedp.Flag("metrics-recording-only", true),
			chromedp.Flag("safebrowsing-disable-auto-update", true),
			chromedp.Flag("enable-automation", true),
	}...)

	defer cancel()

	newCtx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	if err := chromedp.Run(newCtx); err != nil {
		log.Error.Println(err)
	}

	queue := NewJobQueue(4, &newCtx)
	queue.Start()
	defer queue.Stop()

	gamelist = Games{Id: 1, Map: make(map[string]uint)}

	for p := 1; p <= 287; p++ {
		if err := chromedp.Run(newCtx,
			chromedp.Emulate(device.IPadPro),
			chromedp.Navigate(fmt.Sprintf(`https://www.askgamblers.com/free-online-slots/%d`, p)),
			chromedp.Nodes("a", &links, chromedp.ByQueryAll),
		); err != nil {
			log.Error.Println(err)
		} else {
			for _, link := range links {
				if len(link.Attributes) > 1 && (strings.HasPrefix(link.Attributes[1], "/game/slots/") || strings.HasPrefix(link.Attributes[1], "/video-slots/")) {
					if gamelist.Get(link.Attributes[1]) == 0 {
						gamelist.Add(link.Attributes[1])
						queue.Submit(&GameJob{URI: link.Attributes[1]})
					}
				}
			}
		}
	}
}

func (job *GameJob) Process(ctx *context.Context) {
	var (
		scripts     []*cdp.Node
		slotDetails string
		G           model.Game
	)

	game := job.URI

	log.Info.Printf("Processing '%s'", game)

	if err := chromedp.Run(*ctx,
		chromedp.Emulate(device.IPadMinilandscape),
		chromedp.ActionFunc(func(ctx context.Context) error {
			expr := cdp.TimeSinceEpoch(time.Now().Add(30 * time.Minute))
			success, err := network.SetCookie("GameId", fmt.Sprintf("%d", gamelist.Get(game))).
				WithExpires(&expr).
				WithDomain("www.askgamblers.com").
				WithHTTPOnly(true).
				Do(ctx)
			if err != nil {
				return err
			}
			if !success {
				return fmt.Errorf("could not set cookie")
			}

			return nil
		}),
		chromedp.Navigate(`https://www.askgamblers.com` + game),
		//chromedp.WaitVisible(`content-main`, chromedp.BySearch),
		chromedp.Nodes("script", &scripts, chromedp.ByQueryAll),
		chromedp.Text("top10-list top10-list-full-width", &slotDetails, chromedp.BySearch),
		//chromedp.WaitVisible("#play-game", chromedp.ByID),
		chromedp.Click("#play-game", chromedp.ByID),
		//chromedp.WaitVisible("#age-over", chromedp.ByID),
		//chromedp.Click("#age-over", chromedp.ByID),
		chromedp.Sleep(30 * time.Second),
	); err != nil {
		log.Error.Println(err)
	}

	details := strings.Split(slotDetails, "\n \n\n")
	if len(details) > 0 {
		for n, d := range details[1:] {
			field := strings.Split(d, "\n\n ")
			switch n {
			case 0:
				G.Vendor = field[1]
			case 1:
				G.Type = field[1]
			case 2:
				if s, err := strconv.ParseInt(field[1], 10, 32); err == nil {
					G.Paylines = int(s)
				}
			case 3:
				if s, err := strconv.ParseInt(field[1], 10, 32); err == nil {
					G.Reels = int(s)
				}
			case 4:
				if s, err := strconv.ParseInt(field[1], 10, 32); err == nil {
					G.MinCoinsPerLine = int(s)
				}
			case 5:
				if s, err := strconv.ParseInt(field[1], 10, 32); err == nil {
					G.MaxCoinsPerLine = int(s)
				}
			case 6:
				if s, err := strconv.ParseFloat(field[1], 32); err == nil {
					G.MinCoinsSize = float32(s)
				}
			case 7:
				if s, err := strconv.ParseFloat(field[1], 32); err == nil {
					G.MaxCoinsSize = float32(s)
				}
			case 8:
				if s, err := strconv.ParseFloat(field[1], 32); err == nil {
					G.Jackpot = float32(s)
				}
			case 9:
				if s, err := strconv.ParseFloat(strings.Trim(field[1], "%"), 32); err == nil {
					G.RTP = float32(s)
				}
			}
		}
	}

	j := LdJson{}
	for _, script := range scripts {
		if len(script.Attributes) > 1 {
			if script.Attributes[1] == "application/ld+json" && len(script.Children) > 0 {
				if err := json.Unmarshal([]byte(script.Children[0].NodeValue), &j); err != nil {
					log.Error.Println(err)
				}
				G.Name = j.Name
				descr := strings.Split(j.Description, "Software: ")
				G.Description = descr[0]
				if s, err := strconv.ParseFloat(j.Rating.Value, 32); err == nil {
					G.RatingValue = float32(s)
				}
				if s, err := strconv.ParseInt(j.Rating.Count, 10, 32); err == nil {
					G.RatingCount = int(s)
				}
			}
			if script.Attributes[1] == "game-code" {
				G.Object = strings.TrimSpace(script.Children[0].NodeValue)
			}
		}
	}

	var buf []byte
	if err := chromedp.Run(*ctx,
		chromedp.Emulate(device.IPadMinilandscape),
		chromedp.Navigate(j.Image),
		chromedp.CaptureScreenshot(&buf),
	); err != nil {
		log.Error.Println(err)
	}

	G.PublishReady = true
	if err := createGame(&G, &buf, j.Image); err != nil {
		log.Error.Println(err)
	}
}

func createGame(g *model.Game, b *[]byte, u string) error {
	parts := strings.Split(u, "/")
	filename := parts[len(parts)-1]
	tmpfile, err := ioutil.TempFile("", "*"+filename)
	if err != nil {
		return err
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write(*b); err != nil {
		return err
	}
	tmpfile.Seek(0, 0)

	if len(g.Object) > 0 {
		g.Preview.CropOptions = make(map[string]*media.CropOption)
		g.Preview.CropOptions["hd"] = &media.CropOption{X: 480, Y: 240, Height: 943, Width: 1440}
		g.Preview.CropOptions["sd"] = &media.CropOption{X: 480, Y: 240, Height: 943, Width: 1440}
		g.Preview.CropOptions["preview"] = &media.CropOption{X: 480, Y: 240, Height: 943, Width: 1440}
		g.Preview.Crop = true
		g.Preview.Scan(tmpfile)
		if err := db.DB.Save(g).Error; err != nil {
			return err
		}
	}

	return nil
}

// NewJobQueue - creates a new job queue
func NewJobQueue(maxWorkers int, ctx *context.Context) *JobQueue {
	workersStopped := sync.WaitGroup{}
	readyPool := make(chan chan Job, maxWorkers)
	workers := make([]*Worker, maxWorkers, maxWorkers)
	for i := 0; i < maxWorkers; i++ {
		workers[i] = NewWorker(readyPool, workersStopped, ctx)
	}
	return &JobQueue{
		internalQueue:     make(chan Job),
		readyPool:         readyPool,
		workers:           workers,
		dispatcherStopped: sync.WaitGroup{},
		workersStopped:    workersStopped,
		quit:              make(chan bool),
	}
}

// Start - starts the worker routines and dispatcher routine
func (q *JobQueue) Start() {
	for i := 0; i < len(q.workers); i++ {
		q.workers[i].Start()
	}
	go q.dispatch()
}

// Stop - stops the workers and sispatcher routine
func (q *JobQueue) Stop() {
	q.quit <- true
	q.dispatcherStopped.Wait()
}

func (q *JobQueue) dispatch() {
	q.dispatcherStopped.Add(1)
	for {
		select {
		case job := <-q.internalQueue: // We got something in on our queue
			workerChannel := <-q.readyPool // Check out an available worker
			workerChannel <- job           // Send the request to the channel
		case <-q.quit:
			for i := 0; i < len(q.workers); i++ {
				q.workers[i].Stop()
			}
			q.workersStopped.Wait()
			q.dispatcherStopped.Done()
			return
		}
	}
}

// Submit - adds a new job to be processed
func (q *JobQueue) Submit(job Job) {
	q.internalQueue <- job
}

// NewWorker - creates a new worker
func NewWorker(readyPool chan chan Job, done sync.WaitGroup, ctx *context.Context) *Worker {
	return &Worker{
		done:             done,
		readyPool:        readyPool,
		assignedJobQueue: make(chan Job),
		quit:             make(chan bool),
		ctx:              ctx,
	}
}

// Start - begins the job processing loop for the worker
func (w *Worker) Start() {
	go func() {
		blankCtx, cancel := chromedp.NewContext(*w.ctx)
		defer cancel()
		if err := chromedp.Run(blankCtx); err != nil {
			log.Error.Println(err)
		}
		w.done.Add(1)
		for {
			w.readyPool <- w.assignedJobQueue // check the job queue in
			select {
			case job := <-w.assignedJobQueue: // see if anything has been assigned to the queue
				job.Process(&blankCtx)
			case <-w.quit:
				w.done.Done()
				return
			}
		}
	}()
}

// Stop - stops the worker
func (w *Worker) Stop() {
	w.quit <- true
}

func (g *Games) Add(key string) {
	g.Lock()
	g.Map[key] = g.Id
	g.Id++
	g.Unlock()
}

func (g *Games) Get(key string) uint {
	g.RLock()
	v, ok := g.Map[key]
	g.RUnlock()
	if ok {
		return v
	} else {
		return 0
	}
}

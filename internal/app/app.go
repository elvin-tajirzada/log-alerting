package app

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/elvin-tajirzada/log-alerting/internal/config"
	"github.com/elvin-tajirzada/log-alerting/pkg/contactpoint"
	"github.com/elvin-tajirzada/log-alerting/pkg/db"
	"github.com/go-co-op/gocron"
)

type (
	App interface {
		Start()
		Shutdown()
	}

	app struct {
		Config    *config.Config
		Loki      *db.Loki
		Telegram  *contactpoint.Telegram
		Scheduler *gocron.Scheduler
	}
)

func New() (App, error) {
	conf, err := config.New()
	if err != nil {
		return nil, fmt.Errorf("unable to create config: %v", err)
	}

	loki := db.NewLoki(conf)

	telegram, err := contactpoint.NewTelegram(conf, loki)
	if err != nil {
		return nil, fmt.Errorf("unable to create telegram: %v", err)
	}

	scheduler := gocron.NewScheduler(conf.Timezone)

	return &app{
		Config:    conf,
		Loki:      loki,
		Telegram:  telegram,
		Scheduler: scheduler,
	}, nil
}

func (a *app) Start() {
	if _, err := a.Scheduler.Every(a.Config.EvaluationTime).Do(a.Telegram.SendErrorMessage, a.Config.EvaluationTime); err != nil {
		log.Fatalf("Unable to setings up cron job for the SendErrorMessage func: %v", err)
	}

	if _, err := a.Scheduler.Every(1).Day().At("00:00").Do(a.Telegram.SendDailyReport, time.Hour*24); err != nil {
		log.Fatalf("Unable to settings up cron job for the SendDailyReport func: %v", err)
	}

	a.Scheduler.StartAsync()

	log.Println("Log alerting started successfully")
}

func (a *app) Shutdown() {
	signals := make(chan os.Signal)
	done := make(chan struct{})
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-signals
		log.Printf("Received signal: %v\n", sig)

		a.Scheduler.Stop()
		if a.Scheduler.IsRunning() {
			log.Println("Scheduler could not stopped")
		} else {
			log.Println("Scheduler stopped")
		}

		log.Println("Log alerting exited. Bye...")
		done <- struct{}{}
	}()

	<-done
}

package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"

	"github.com/mackerelio/go-osstat/cpu"
	"github.com/mackerelio/go-osstat/memory"
	"github.com/mackerelio/go-osstat/uptime"
)

type config struct {
	Username string `env:"BASIC_USER,required"`
	Password string `env:"BASIC_PASS,required"`
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	conf := config{}
	if err := env.Parse(&conf); err != nil {
		logrus.Fatalf("%+v\n", err)
	}
	logrus.Infof("config %+v", conf)

	// add the middleware to the server
	stat := http.NewServeMux()
	stat.HandleFunc("/status", basicAuthMiddleware(statusHandler, conf))

	// start the server
	log.Fatal(http.ListenAndServe("localhost:8888", stat))

}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	// get the system uptime
	ut, err := uptime.Get()
	if err != nil {
		logrus.WithError(err).Error("error getting uptime")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	utSeconds := int64(ut.Seconds())

	mem, err := memory.Get()
	if err != nil {
		logrus.WithError(err).Error("error getting memory")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	before, err := cpu.Get()
	if err != nil {
		logrus.WithError(err).Error("error getting cpu 'before'")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	time.Sleep(1 * time.Second)
	after, err := cpu.Get()
	if err != nil {
		logrus.WithError(err).Error("error getting cpu 'after'")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	total := float64(after.Total - before.Total)
	user := float64(after.User-before.User) / total * 100
	system := float64(after.System-before.System) / total * 100
	idle := float64(after.Idle-before.Idle) / total * 100
	memUsePerc := float64(mem.Used) / float64(mem.Total/100)
	memUseGib := float64(mem.Used) / 1_000_000_000
	currentTime := time.Now().Unix()
	logrus.Infof("time: %v, uptime: %v, memory: %.2f gib (%.2f%%) used, cpu user %.3f system %.3f idle %.3f",
		currentTime, ut.String(), memUseGib, memUsePerc, user, system, idle)

	type memory struct {
		UsedPerc float64 `json:"used_perc"`
		UsedGib  float64 `json:"used_gib"`
	}
	type cpu struct {
		User   float64 `json:"user"`
		System float64 `json:"system"`
		Idle   float64 `json:"idle"`
	}
	type status struct {
		Time   int64  `json:"time"`
		Uptime int64  `json:"up_time_s"`
		Memory memory `json:"memory"`
		CPU    cpu    `json:"cpu"`
	}

	current := status{
		Time:   currentTime,
		Uptime: utSeconds,
		Memory: memory{
			UsedPerc: memUsePerc,
			UsedGib:  memUseGib,
		},
		CPU: cpu{
			User:   user,
			System: system,
			Idle:   idle,
		},
	}

	if err := json.NewEncoder(w).Encode(current); err != nil {
		logrus.WithError(err).Error("error encoding json")
		w.WriteHeader(http.StatusInternalServerError)
	}

}

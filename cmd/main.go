package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/cortexproject/cortex/pkg/util/flagext"

	kitlog "github.com/go-kit/kit/log"

	"github.com/hikhvar/journaldtail/pkg/storage"

	"github.com/coreos/go-systemd/sdjournal"
	"github.com/grafana/loki/pkg/promtail/client"
	"github.com/hikhvar/journaldtail/pkg/journald"
	"github.com/pkg/errors"
	"github.com/prometheus/common/model"
)

var (
	lokiHostURL = flag.String(
		"server",
		"http://localhost:3100/api/prom/push",
		"loki server address",
	)
	withHistory = flag.Bool(
		"history",
		false,
		"collect all available logs from beginning of journal",
	)
)

func main() {
	flag.Parse()

	var logger kitlog.Logger
	logger = kitlog.NewLogfmtLogger(kitlog.NewSyncWriter(os.Stderr))
	log.SetOutput(kitlog.NewStdlibAdapter(logger))
	// TODO: Store state on disk
	memStorage := storage.Memory{}
	journal, err := sdjournal.NewJournal()
	if err != nil {
		log.Fatal(fmt.Sprintf("could not open journal: %s", err.Error()))
	}
	reader := journald.NewReader(journal, &memStorage)
	if !*withHistory {
		err = reader.ToTail()
		if err != nil {
			log.Fatalf("failed to seek %v", err)
		}
	}

	// TODO: Read from CLI
	if v, isSet := os.LookupEnv("LOKI_URL"); isSet {
		*lokiHostURL = v
	}

	cfg := client.Config{
		URL: flagext.URLValue{
			URL: MustParseURL(*lokiHostURL),
		},
		Timeout: 30 * time.Second,
	}
	lokiClient, err := client.New(cfg, logger)
	if err != nil {
		log.Fatal(fmt.Sprintf("could not create loki client: %s", err.Error()))
	}
	err = TailLoop(reader, lokiClient)
	if err != nil {
		log.Fatal(fmt.Sprintf("failed to tail journald: %s", err.Error()))
	}
}

func TailLoop(reader *journald.Reader, writer client.Client) error {
	var lastTS time.Time
	for {
		r, err := reader.Next()
		if err != nil {
			return errors.Wrap(err, "could not get next journal entry")
		}
		if r != nil {
			ls := ToLabelSet(r)
			ts := journald.ToGolangTime(r.RealtimeTimestamp)
			msg := r.Fields[sdjournal.SD_JOURNAL_FIELD_MESSAGE]

			if ts.Before(lastTS) {
				log.Fatal(fmt.Sprintf("%s is before %s! Message: %s", ts, lastTS, msg))
			}
			lastTS = ts
			err = writer.Handle(ls, ts, msg)
			comm, ok := r.Fields["_COMM"]
			if ok {
				log.Printf("%s", comm)
			} else {
				log.Printf("POOP TORPEDO")
			}
			if err != nil {
				return errors.Wrap(err, "could not enque systemd logentry")
			}
		}

	}
}

func ToLabelSet(reader *sdjournal.JournalEntry) model.LabelSet {
	ret := make(model.LabelSet)
	for key, value := range reader.Fields {
		if key != sdjournal.SD_JOURNAL_FIELD_MESSAGE {
			ret[model.LabelName(key)] = model.LabelValue(value)
		}
	}
	return ret
}

func MustParseURL(input string) *url.URL {
	u, err := url.Parse(input)
	if err != nil {
		panic(fmt.Sprintf("could not parse static url: %s", input))
	}
	return u
}

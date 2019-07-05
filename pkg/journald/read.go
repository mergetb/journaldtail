package journald

import (
	"log"
	"time"

	"github.com/coreos/go-systemd/sdjournal"
	"github.com/pkg/errors"
)

type Reader struct {
	src *sdjournal.Journal
}

func NewReader(src *sdjournal.Journal) *Reader {
	return &Reader{
		src: src,
	}
}

func (r *Reader) ToTail() error {
	err := r.src.SeekTail()
	if err != nil {
		return err
	}
	_, err = r.src.Next()
	return err
}

// Next blocks until the next journal entry is available
func (r *Reader) Next() (*sdjournal.JournalEntry, error) {
	advanced, err := r.advance()
	if err != nil {
		return nil, errors.Wrap(err, "failed to advance")
	}
	if !advanced {
		r.src.Wait(sdjournal.IndefiniteWait)
		log.Printf("wait finished")
		advanced, err = r.advance()
		if advanced != true {
			//return nil, errors.New("finished wait but could not advance")
			return nil, nil
		}
		log.Printf("NO bonk")
		if err != nil {
			return nil, errors.Wrap(err, "failed to advance after wait")
		}
	}
	entry, err := r.src.GetEntry()
	return entry, errors.Wrap(err, "could not get next entry")
}

// advance tries to jump to next journal entry.
func (r *Reader) advance() (bool, error) {
	rc, err := r.src.Next()
	if err != nil {
		return false, errors.Wrap(err, "could not get next. ")
	}
	return rc != 0, nil
}

func ToGolangTime(sdTime uint64) time.Time {
	return time.Unix(0, int64(time.Duration(sdTime)*time.Microsecond))
}

package generate

import (
	"errors"
	"github.com/ledgerwatch/turbo-geth/core"
	"github.com/ledgerwatch/turbo-geth/ethdb"
	"github.com/ledgerwatch/turbo-geth/log"
	"os"
	"os/signal"
	"time"
)

func RegenerateIndex(chaindata string, csBucket []byte) error {
	db, err := ethdb.NewBoltDatabase(chaindata)
	if err != nil {
		return err
	}
	ch := make(chan os.Signal, 1)
	quitCh := make(chan struct{})
	signal.Notify(ch, os.Interrupt)
	go func() {
		<-ch
		close(quitCh)
	}()

	ig := core.NewIndexGenerator(db, quitCh)
	cs, ok := core.CSMapper[string(csBucket)]
	if !ok {
		return errors.New("unknown changeset")
	}

	err = ig.DropIndex(cs.IndexBucket)
	if err != nil {
		return err
	}
	startTime := time.Now()
	log.Info("Index generation started", "start time", startTime)
	err = ig.GenerateIndex(0, 0, csBucket)
	if err != nil {
		return err
	}
	log.Info("Index is successfully regenerated", "it took", time.Since(startTime))
	return nil
}

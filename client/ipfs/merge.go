package ipfs

import (
	"context"
	"fmt"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ipfs/go-cid"
	ipfsapi "github.com/ipfs/go-ipfs-api"
	blockstore "github.com/ipfs/go-ipfs-blockstore"
	ipfsfiles "github.com/ipfs/go-ipfs-files"
	"github.com/multiformats/go-multihash"
	"golang.org/x/xerrors"
)

type opts struct {
	IpfsAPI            string `getopt:"--ipfs-api             A read/write IPFS API URL"`
	IpfsAPIMaxWorkers  uint   `getopt:"--ipfs-api-max-workers Max amount of parallel API requests"`
	IpfsAPITimeoutSecs uint   `getopt:"--ipfs-api-timeout     Max amount of seconds for a single API operation"`
	ShowProgress       bool   `getopt:"--show-progress        Print progress on STDERR, default when a TTY"`
	AggregateVersion   uint   `getopt:"--aggregate-version    The version of aggregate to produce"`
	SkipDagStat        bool   `getopt:"--skip-dag-stat        Do not query he API for the input dag stats"`
	Help               bool   `getopt:"-h --help              Display help"`
}

// pulls cids from an AllKeysChan and sends them concurrently via multiple workers to an API
func writeoutBlocks(externalCtx context.Context, opts *opts, bs blockstore.Blockstore) error {

	innerCtx, shutdownWorkers := context.WithCancel(externalCtx)
	defer shutdownWorkers()

	akc, err := bs.AllKeysChan(innerCtx)
	if err != nil {
		return err
	}

	maxWorkers := opts.IpfsAPIMaxWorkers
	finishCh := make(chan struct{}, 1)
	errCh := make(chan error, maxWorkers)

	// WaitGroup as we want everyone to fully "quit" before we return
	var wg sync.WaitGroup
	blocksDone := new(uint64)

	for i := uint(0); i < maxWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			api := ipfsapi.NewShell(opts.IpfsAPI)
			api.SetTimeout(time.Second * time.Duration(opts.IpfsAPITimeoutSecs))

			for {
				select {

				case <-innerCtx.Done():
					// something caused us to stop, whatever it is parent knows why
					return

				case c, chanOpen := <-akc:

					if !chanOpen {
						select {
						case finishCh <- struct{}{}:
						default:
							// if we can't signal feeder is done - someone else already did
						}
						return
					}

					blk, err := bs.Get(c)
					if err != nil {
						errCh <- err
						return
					}

					// copied entirety of ipfsapi.BlockPut() to be able to pass in our own ctx 🤮
					res := new(struct{ Key string })
					err = api.Request("block/put").
						Option("format", cid.CodecToStr[c.Prefix().Codec]).
						Option("mhtype", multihash.Codes[c.Prefix().MhType]).
						Option("mhlen", c.Prefix().MhLength).
						Body(
							ipfsfiles.NewMultiFileReader(
								ipfsfiles.NewSliceDirectory([]ipfsfiles.DirEntry{
									ipfsfiles.FileEntry(
										"",
										ipfsfiles.NewBytesFile(blk.RawData()),
									),
								}),
								true,
							),
						).
						Exec(innerCtx, res)
					// end of 🤮

					if err != nil {
						errCh <- err
						return
					}

					if res.Key != c.String() {
						errCh <- xerrors.Errorf("unexpected cid mismatch after /block/put: expected %s but got %s", c, res.Key)
						return
					}

					if opts.ShowProgress {
						atomic.AddUint64(blocksDone, 1)
					}
				}
			}
		}()
	}

	var blocksTotal, lastPct uint64
	var progressTick <-chan time.Time
	if opts.ShowProgress {
		// this works because of how AllKeysChan behaves on rambs
		blocksTotal = uint64(len(akc))
		fmt.Fprint(os.Stderr, "0% of blocks written\r")
		t := time.NewTicker(250 * time.Millisecond)
		progressTick = t.C
		defer t.Stop()
	}

	var workerError error
watchdog:
	for {
		select {

		case <-finishCh:
			break watchdog

		case <-externalCtx.Done():
			break watchdog

		case workerError = <-errCh:
			shutdownWorkers()
			break watchdog

		case <-progressTick:
			curPct := 100 * atomic.LoadUint64(blocksDone) / blocksTotal
			if curPct != lastPct {
				lastPct = curPct
				fmt.Fprintf(os.Stderr, "%d%% of blocks written\r", lastPct)
			}
		}
	}

	wg.Wait()
	close(errCh) // closing a buffered channel keeps any buffered values for <-

	if workerError != nil {
		return workerError
	}
	if err := <-errCh; err != nil {
		return err
	}
	return externalCtx.Err()
}

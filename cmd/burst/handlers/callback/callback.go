package callback

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/bitburst/burstconsumer/cmd/burst/burstcfg"
	"github.com/bitburst/burstconsumer/cmd/burst/models/user"
	"github.com/bitburst/burstconsumer/pkg/mdw"
	"github.com/bitburst/burstconsumer/pkg/request"
	"github.com/julienschmidt/httprouter"
)

type clbReqPld struct {
	ObjIDS []int `json:"object_ids,omitempty"`
}

func persistCallback(cfg *burstcfg.BurstCfg) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		defer r.Body.Close()
		pld := &clbReqPld{}
		json.NewDecoder(r.Body).Decode(pld)

		// do not block /callback resource responce and start processing ids in
		// separate goroutine
		go func(pld *clbReqPld) {
			var wg sync.WaitGroup
			for _, id := range pld.ObjIDS {
				wg.Add(1)
				// since each request can take up to 4s, we execute them concurrently
				go func(id int, wg *sync.WaitGroup) {
					defer wg.Done()

					u := &user.User{}
					req := request.
						New("GET", fmt.Sprintf("http://localhost:9010/objects/%d", id), nil).
						Do().
						Decode(u)

					if err := req.HasError(); err != nil {
						cfg.Errlog.Printf("failed to fetch object by id: %d, err: %s\n", id, err)
						return
					}

					u.Seen = time.Now().Format("2006-01-02 15:04:05")
					cfg.Infolog.Printf("received: %+v\n", u)

					if err := user.Create(context.Background(), cfg.DB.Client, u); err != nil {
						cfg.Errlog.Printf("unable to create user: %s\n", err)
					}

				}(id, &wg)
			}
			wg.Wait()
		}(pld)

		w.WriteHeader(http.StatusOK)
	}
}

func Do(cfg *burstcfg.BurstCfg) httprouter.Handle {
	// combine middleware chain that will return httprouter handle after all
	// middleware is executed
	return mdw.Chain(
		persistCallback(cfg),
		mdw.Recover(cfg.Infolog),
		mdw.LogRequest(cfg.Infolog),
	)
}

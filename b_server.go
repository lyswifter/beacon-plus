package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ipfs/go-datastore"
	ltypes "github.com/lyswifter/beacon-plus/localtype"
	"github.com/lyswifter/beacon-plus/log"
	"golang.org/x/xerrors"
)

func setupBeaconServer() {
	r := gin.Default()

	handleClientBeaconAPI(r)

	handleQueryMinfos(r)

	r.Run("0.0.0.0:9090") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func handleClientBeaconAPI(r *gin.Engine) {
	r.GET("/public/:epoch", func(c *gin.Context) {
		start := time.Now()

		epoch := c.Param("epoch")
		ip := c.ClientIP()

		go func() {
			mAddr := c.Query("m")
			symbol := c.Query("s")
			pip := c.Query("i")
			if pip != "" {
				ip = pip
			}

			if mAddr != "" {
				minfo := ltypes.MinerInfo{
					Address: mAddr,
					Symbol:  symbol,
					IP:      ip,
				}

				err := saveMInfo(minfo.Address, minfo)
				if err != nil {
					log.Infof("save minfo err %s", err)
					return
				}
			}
		}()

		key := datastore.NewKey(epoch)

		ishas, err := BeaconDB.Has(key)
		if err != nil {
			log.Infof("entrys: has %s", err)
			return
		}

		if !ishas {
			c.JSON(500, gin.H{
				"status":  "Err",
				"epoch":   epoch,
				"message": fmt.Sprintf("get beacon for epoch: %s, is not exist", epoch),
			})
		}

		qt, err := BeaconDB.Get(key)
		if err != nil {
			log.Infof("entrys: list %s", err)
			return
		}

		entrys := ltypes.BeaconEntryInfo{}
		err = json.Unmarshal(qt, &entrys)
		if err != nil {
			return
		}

		c.JSON(200, gin.H{
			"status":  "Ok",
			"epoch":   epoch,
			"message": string(qt),
		})

		log.Infof("Request from client: %s round: %d data: %v took: %s", ip, entrys.Round, hex.EncodeToString(entrys.Entry.Data), time.Since(start).String())
	})
}

func handleQueryMinfos(r *gin.Engine) {
	r.GET("/public/minfos", func(c *gin.Context) {

		minfos, err := readMInfos()
		if err != nil {
			c.JSON(500, gin.H{
				"status":  "Err",
				"message": xerrors.Errorf("read minfos err %s", err.Error()),
			})
			return
		}

		out, err := json.Marshal(minfos)
		if err != nil {
			c.JSON(500, gin.H{
				"status":  "Err",
				"message": fmt.Sprintf("marshal infos err %s", err.Error()),
			})
			return
		}

		c.JSON(200, gin.H{
			"statue":  "Ok",
			"message": string(out),
		})

		log.Infof("Query minfos from client: %s len: %d", c.ClientIP(), len(minfos))
	})

	r.GET("/public/minfos/:maddr", func(c *gin.Context) {
		maddr := c.Param("maddr")

		minfo, err := readmInfo(maddr)
		if err != nil {
			c.JSON(500, gin.H{
				"status":  "Err",
				"message": err.Error(),
			})
			return
		}

		out, err := json.Marshal(minfo)
		if err != nil {
			c.JSON(500, gin.H{
				"status":  "Err",
				"message": fmt.Sprintf("marshal infos err %s", err.Error()),
			})
			return
		}

		c.JSON(200, gin.H{
			"statue":  "Ok",
			"message": string(out),
		})

		log.Infof("Query minfo from client: %s addr: %s", c.ClientIP(), maddr)
	})
}

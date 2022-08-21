package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"maidian/dal"
	"maidian/model"
	"maidian/util"
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/prometheus/model/labels"

	"github.com/prometheus/prometheus/tsdb"
)

var pvDB *tsdb.DB

func EmitPv(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		fmt.Printf("read body error, err=%v\n", err)
		return
	}
	var bodyStruct model.EmitPVStruct
	err = json.Unmarshal(bodyBytes, &bodyStruct)
	if err != nil {
		fmt.Printf("json unmarshal struct error, err=%v\n", err)
	}
	fmt.Println(bodyStruct)
	timeStamp, _ := strconv.ParseInt(bodyStruct.Data.Time, 10, 64)
	fmt.Printf("time_stamp %v\n", timeStamp)
	dal.EmitCounter(pvDB, 1, timeStamp, map[string]string{
		"app_id": bodyStruct.AppId,
		"url":    bodyStruct.Data.Url,
		"did":    fmt.Sprintf("%d", util.RandomInt64()),
	})
	fmt.Fprintf(w, "Hello World!") //这个写入到w的是输出到客户端的
}

func QueryPv(w http.ResponseWriter, r *http.Request) {
	// Open a querier for reading.
	nowTime := time.Now()
	nowTimeHourly := time.Date(nowTime.Year(), nowTime.Month(), nowTime.Day(), nowTime.Hour(), 0, 0, 0, nowTime.Location())
	startTimeHourly := nowTimeHourly.Add(time.Hour * 24 * 3 * -1)
	querier, err := pvDB.Querier(context.Background(), startTimeHourly.UnixMilli(), nowTimeHourly.UnixMilli())
	util.NoErr(err)
	ss := querier.Select(false, nil, labels.MustNewMatcher(labels.MatchRegexp, "url", ".*"))
	metricsMap := make(map[string][]*model.Point)
	for ss.Next() {
		series := ss.At()
		fmt.Println("series:", series.Labels().String())
		tagMap := series.Labels().Map()
		if _, exist := metricsMap["url:"+tagMap["url"]]; !exist {
			metricsMap["url:"+tagMap["url"]] = make([]*model.Point, 24*3)
			for i := range metricsMap["url:"+tagMap["url"]] {
				metricsMap["url:"+tagMap["url"]][i] = &model.Point{X: int64(i)*3600*1000 + startTimeHourly.UnixMilli()}
			}
		}

		it := series.Iterator()
		for it.Next() {
			st, v := it.At()
			offset := (st - startTimeHourly.UnixMilli()) / (60 * 60 * 1000)
			metricsMap["url:"+tagMap["url"]][offset].Y += v
		}
	}
	fmt.Println(util.JsonStr(metricsMap))

	err = querier.Close()
	util.NoErr(err)
	fmt.Fprintf(w, "Hello World!")
}

func init() {
	pvDB = dal.InitTsDb("pv2")
}

package dal

import (
	"context"
	"maidian/util"

	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/tsdb"
)

func Query(db *tsdb.DB) {

}

func EmitCounter(db *tsdb.DB, counter int, timeStamp int64, tagsKv map[string]string) {
	appender := db.Appender(context.Background())

	tags := make([]string, 0, len(tagsKv)*2)
	for k, v := range tagsKv {
		tags = append(tags, k, v)
	}
	series := labels.FromStrings(tags...)

	// Ref is 0 for the first append since we don't know the reference for the series.
	_, err := appender.Append(0, series, timeStamp, float64(counter))
	util.NoErr(err)

	// Commit to storage.
	err = appender.Commit()
	util.NoErr(err)
}

func InitTsDb(name string) *tsdb.DB {
	db, err := tsdb.Open("./"+name, nil, nil, tsdb.DefaultOptions(), nil)
	if err != nil {
		panic(err)
	}
	return db
}

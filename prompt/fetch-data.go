package prompt

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"sync"
	"time"

	"github.com/jinzhu/now"
	"github.com/romainmenke/report-imgix-usage/httpcache"
	"github.com/romainmenke/report-imgix-usage/sources"
)

func getAllData(client *http.Client) *sources.Sources {
	fmt.Println("Fetching all report data, this might take a while...")

	cacheRoundTripper, closeCacheRoundTripper := httpcache.CachingRoundTripper(client)
	defer closeCacheRoundTripper()

	client.Transport = cacheRoundTripper

	foundSources, err := sources.Get(client, 0)
	if err != nil {
		panic(err)
	}

	oldest := math.MaxInt64
	for _, sourceData := range foundSources.Data {
		if sourceData.Attributes.DateCreated < oldest {
			oldest = sourceData.Attributes.DateCreated
		}
	}

	oldestT := time.Unix(int64(oldest), 0).AddDate(0, -1, 0) // to ensure we get everything
	y, m, _ := timeDiff(oldestT, time.Now())
	months := (y * 12) + m

	wg := sync.WaitGroup{}

	for monthOffset := 0; monthOffset < months; monthOffset++ {
		start := now.New(time.Now().AddDate(0, -1*(monthOffset+1), 0)).BeginningOfMonth()
		end := now.New(time.Now().AddDate(0, -1*(monthOffset+1), 0)).EndOfMonth()

		for _, sourceData := range foundSources.Data {

			log.Printf("downloading %d - %s : %s", start.Year(), start.Month().String(), sourceData.Attributes.Name)

			wg.Add(1)
			go func(x *sources.Data, s time.Time, t time.Time) {
				defer wg.Done()

				err := x.GetCounters(client, s, t)
				if err != nil {
					panic(err)
				}

			}(sourceData, start, end)

		}
	}

	wg.Wait()

	return foundSources
}

// https://play.golang.org/p/yM8w1KNqRE
func timeDiff(t1, t2 time.Time) (years, months, days int) {
	t2 = t2.AddDate(0, 0, 1) // advance t2 to make the range inclusive

	for t1.AddDate(years, 0, 0).Before(t2) {
		years++
	}
	years--

	for t1.AddDate(years, months, 0).Before(t2) {
		months++
	}
	months--

	for t1.AddDate(years, months, days).Before(t2) {
		days++
	}
	days--

	return
}

package httpcache

import (
	"bufio"
	"bytes"
	"context"
	"crypto/sha1"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"

	bolt "go.etcd.io/bbolt"
	"golang.org/x/time/rate"
)

type RoundTripper func(req *http.Request) (*http.Response, error)

func (x RoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return x(req)
}

func CachingRoundTripper(client *http.Client) (http.RoundTripper, func()) {
	next := client.Transport

	db, err := bolt.Open("./imgix-report-cache.db", 0600, nil)
	if err != nil {
		panic(err)
	}

	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("http-cache"))
		if err != nil {
			panic(err)
		}

		return nil
	})

	limiter := rate.NewLimiter(3, 1)

	return RoundTripper(
			func(req *http.Request) (*http.Response, error) {
				if req.Method != http.MethodGet {
					return next.RoundTrip(req)
				}

				if req.Header.Get("Cache-Control") != "no-cache" {
					resp := Get(db, req)
					if resp != nil {
						return resp, nil
					}
				}

				err := limiter.Wait(context.Background())
				if err != nil {
					return nil, err
				}

				resp, err := next.RoundTrip(req)
				if err != nil {
					return nil, err
				}

				if resp.StatusCode/100 != 2 {
					return resp, nil
				}

				Put(db, req, resp)
				return resp, nil
			},
		),
		func() {
			err := db.Close()
			if err != nil {
				log.Println(err)
			}
		}
}

func Get(db *bolt.DB, req *http.Request) *http.Response {
	var response *http.Response

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("http-cache"))
		if b == nil {
			return errors.New("bucket does not exist")
		}

		v := b.Get([]byte(RequestKey(req)))
		if v == nil {
			return nil
		}

		resp, err := http.ReadResponse(bufio.NewReader(bytes.NewBuffer(v)), req)
		if err != nil {
			return err
		}

		response = resp

		return nil
	})

	if err != nil {
		log.Println(err)
	}

	return response
}

func Put(db *bolt.DB, req *http.Request, resp *http.Response) {

	err := db.Update(func(tx *bolt.Tx) error {
		data, err := httputil.DumpResponse(resp, true)
		if err != nil {
			return err
		}

		b := tx.Bucket([]byte("http-cache"))
		if b == nil {
			return errors.New("bucket does not exist")
		}

		return b.Put([]byte(RequestKey(req)), data)
	})

	if err != nil {
		log.Println(err)
	}
}

func RequestKey(req *http.Request) string {
	key := "7935efa4-4fe6-474e-869d-127a56e06c61"
	key += req.Method
	key += req.URL.String()

	key = fmt.Sprintf("%x\n", sha1.Sum([]byte(key)))

	return key
}

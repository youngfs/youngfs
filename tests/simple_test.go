package tests

import (
	"fmt"
	"github.com/appleboy/gofight/v2"
	"github.com/youngfs/youngfs/pkg/util/randutil"
	"net/http"
	"sync"
)

func (s *handlerSuite) TestSmallObjets() {
	const size = 1024
	const blobSize = 5 * 1024
	bodys := make([]string, size)
	wg := sync.WaitGroup{}
	for i := range size {
		wg.Add(1)
		go func() {
			defer wg.Done()
			r := gofight.New()
			r.GET(fmt.Sprintf("/test/aa/bb/cc/%d.txt", i)).
				Run(s.handler, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
					s.Equal(http.StatusNotFound, r.Code)
				})
		}()
	}
	wg.Wait()

	for i := range size {
		wg.Add(1)
		go func() {
			defer wg.Done()
			bodys[i] = randutil.RandString(blobSize)
			r := gofight.New()
			r.PUT(fmt.Sprintf("/test/aa/bb/cc/%d.txt", i)).
				SetBody(bodys[i]).
				Run(s.handler, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
					s.Equal(http.StatusCreated, r.Code)
				})
		}()
	}
	wg.Wait()

	for i := range size {
		wg.Add(1)
		go func() {
			defer wg.Done()
			r := gofight.New()
			r.GET(fmt.Sprintf("/test/aa/bb/cc/%d.txt", i)).
				Run(s.handler, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
					s.Equal(http.StatusOK, r.Code)
					s.Equal(bodys[i], r.Body.String())
				})
		}()
	}
	wg.Wait()

	for i := range size {
		wg.Add(1)
		go func() {
			defer wg.Done()
			r := gofight.New()
			r.DELETE(fmt.Sprintf("/test/aa/bb/cc/%d.txt", i)).
				Run(s.handler, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
					s.Equal(http.StatusNoContent, r.Code)
				})
		}()
	}
	wg.Wait()

	for i := range size {
		wg.Add(1)
		go func() {
			defer wg.Done()
			r := gofight.New()
			r.GET(fmt.Sprintf("/test/aa/bb/cc/%d.txt", i)).
				Run(s.handler, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
					s.Equal(http.StatusNotFound, r.Code)
				})
		}()
	}
	wg.Wait()
}

func (s *handlerSuite) TestBigObjets() {
	const size = 8
	const blobSize = (128 + 32) * 1024 * 1024
	bodys := make([]string, size)
	wg := sync.WaitGroup{}
	for i := range size {
		wg.Add(1)
		go func() {
			defer wg.Done()
			r := gofight.New()
			r.GET(fmt.Sprintf("/test/aa/bb/cc/%d.txt", i)).
				Run(s.handler, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
					s.Equal(http.StatusNotFound, r.Code)
				})
		}()
	}
	wg.Wait()

	for i := range size {
		wg.Add(1)
		go func() {
			defer wg.Done()
			bodys[i] = randutil.RandString(blobSize)
			r := gofight.New()
			r.PUT(fmt.Sprintf("/test/aa/bb/cc/%d.txt", i)).
				SetBody(bodys[i]).
				Run(s.handler, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
					s.Equal(http.StatusCreated, r.Code)
				})
		}()
	}
	wg.Wait()

	for i := range size {
		wg.Add(1)
		go func() {
			defer wg.Done()
			r := gofight.New()
			r.GET(fmt.Sprintf("/test/aa/bb/cc/%d.txt", i)).
				Run(s.handler, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
					s.Equal(http.StatusOK, r.Code)
					s.Equal(bodys[i], r.Body.String())
				})
		}()
	}
	wg.Wait()

	for i := range size {
		wg.Add(1)
		go func() {
			defer wg.Done()
			r := gofight.New()
			r.DELETE(fmt.Sprintf("/test/aa/bb/cc/%d.txt", i)).
				Run(s.handler, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
					s.Equal(http.StatusNoContent, r.Code)
				})
		}()
	}
	wg.Wait()

	for i := range size {
		wg.Add(1)
		go func() {
			defer wg.Done()
			r := gofight.New()
			r.GET(fmt.Sprintf("/test/aa/bb/cc/%d.txt", i)).
				Run(s.handler, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
					s.Equal(http.StatusNotFound, r.Code)
				})
		}()
	}
	wg.Wait()
}

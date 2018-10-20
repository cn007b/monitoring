package ping

import (
	"golang.org/x/net/context"
	"google.golang.org/appengine/urlfetch"
	"net/http"
	"time"

	"go-app/app/config/taxonomy"
	"go-app/app/errors/AppError"
	"go-app/service/datastore/Measurement"
)

// Do performs main ping action and saves result into DataStore.
func Do(ctx context.Context, vo VO) {
	startedAt := time.Now().UnixNano()
	res, err := exec(ctx, vo)
	finishedAt := time.Now().UnixNano()

	if err != nil {
		AppError.Panic(err)
	}

	saveMeasurement(ctx, vo, res, finishedAt-startedAt)
}

func exec(ctx context.Context, vo VO) (r *http.Response, err error) {
	client := urlfetch.Client(ctx)

	switch vo.Method {
	case taxonomy.MethodHead, taxonomy.MethodGet:
		return client.Get(vo.URL)
	}

	return client.Post(vo.URL, vo.ContentType, vo.Body)
}

func saveMeasurement(ctx context.Context, jobVO VO, res *http.Response, took int64) {
	vo := Measurement.EntityVO{
		Project:      jobVO.Project,
		Took:         int(took / 1e6),
		ResponseCode: res.StatusCode,
	}
	Measurement.Add(ctx, vo)
}

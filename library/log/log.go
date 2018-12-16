package log

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"cloud.google.com/go/compute/metadata"
	"cloud.google.com/go/logging"
	"google.golang.org/genproto/googleapis/api/monitoredres"
)

var (
	client    *logging.Client
	ctxKey    = &struct{ temp string }{}
	zone      string
	moduleID  string
	projectID string
	versionID string
	local     bool
)

func Init(isLocal bool) error {
	local = isLocal
	if local {
		return nil
	}

	moduleID = os.Getenv("GAE_SERVICE")
	projectID = os.Getenv("GOOGLE_CLOUD_PROJECT")
	versionID = os.Getenv("GAE_VERSION")

	var err error
	zone, err = metadata.Zone()
	if err != nil {
		log.Fatalf("metadata.Zone() error: %s", err.Error())
	}

	client, err = logging.NewClient(context.Background(), projectID)
	if err != nil {
		return err
	}
	return nil
}

func Close() {
	if client == nil {
		return
	}

	if err := client.Close(); err != nil {
		log.Printf("logging clinet failed to close : %s", err)
	}
}

func WithContext(ctx context.Context, r *http.Request) context.Context {
	return context.WithValue(ctx, ctxKey, r)
}

func Criticalf(ctx context.Context, format string, args ...interface{}) {
	printf(ctx, logging.Critical, format, args...)
}

func Debugf(ctx context.Context, format string, args ...interface{}) {
	printf(ctx, logging.Debug, format, args...)
}

func Errorf(ctx context.Context, format string, args ...interface{}) {
	printf(ctx, logging.Error, format, args...)
}

func Infof(ctx context.Context, format string, args ...interface{}) {
	printf(ctx, logging.Info, format, args...)
}

func Warningf(ctx context.Context, format string, args ...interface{}) {
	printf(ctx, logging.Warning, format, args...)
}

func Duration(ctx context.Context, invocation time.Time, name string) {
	elapsed := time.Since(invocation)
	printf(ctx, logging.Info, fmt.Sprintf("%s lasted %s", name, elapsed))
}

func printf(ctx context.Context, severity logging.Severity, format string, args ...interface{}) {
	if local {
		log.Printf(format, args...)
		return
	}

	r, ok := ctx.Value(ctxKey).(*http.Request)
	if !ok {
		log.Printf("unexpected context. It doesn't have *http.Request")
	}

	traceContext := r.Header.Get("X-Cloud-Trace-Context")
	traceID := strings.Split(traceContext, "/")[0]
	logger := client.Logger("request-log")
	logger.Log(logging.Entry{
		Severity: severity,
		Payload: map[string]interface{}{
			"serviceContext": map[string]interface{}{},
			"message":        fmt.Sprintf(format, args...),
		},
		Resource: &monitoredres.MonitoredResource{
			Type: "gae_app",
			Labels: map[string]string{
				"module_id":  moduleID,
				"project_id": projectID,
				"version_id": versionID,
				"zone":       zone,
			},
		},
		Trace: fmt.Sprintf("projects/%s/traces/%s", projectID, traceID),
	})
}

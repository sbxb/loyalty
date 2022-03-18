package accrual

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/sbxb/loyalty/internal/logger"
	"github.com/sbxb/loyalty/models"
)

const defaultTimeout = 1 * time.Second

type ClientError struct {
	msg string
}

func (ce *ClientError) Error() string {
	return ce.msg
}

func NewClientError(msg string) *ClientError {
	return &ClientError{
		msg: msg,
	}
}

type AccrualClient struct {
	client http.Client
	url    string
}

func NewAccrualClient(address string) (*AccrualClient, error) {
	accrualURL := address + "/api/orders/"
	return &AccrualClient{
		client: http.Client{
			Timeout: defaultTimeout,
		},
		url: accrualURL,
	}, nil
}

func (ac *AccrualClient) DoWork(ctx context.Context, jobQueue chan Job) {
	logger.Info("client: DoWork(): begin")
	for {
		select {
		case <-ctx.Done():
			logger.Info("client: DoWork(): context cancelled")
			return
		case job := <-jobQueue:
			logger.Infof("client: DoWork(): got job %v", job)
			time.Sleep(500 * time.Millisecond)
			logger.Infof("client: DoWork(): Trying to get " + ac.url + job.orderNumber)
			resp, err := ac.client.Get(ac.url + job.orderNumber)
			if err != nil {
				logger.Infof("client: DoWork(): Get request failed")
				jobQueue <- job
				break
			}
			ar, clientErr := processResponse(resp)
			if clientErr != nil {
				logger.Warning(clientErr.Error())
				jobQueue <- job
				break
			}
			logger.Infof("client: DoWork(): got response %v", ar)
		}
	}
}

func processResponse(resp *http.Response) (*models.AccrualResponse, *ClientError) {
	defer resp.Body.Close()
	logger.Infof("client: processResponse(): got code %d", resp.StatusCode)
	switch resp.StatusCode {
	case 200:
		ar := &models.AccrualResponse{}
		if err := json.NewDecoder(resp.Body).Decode(ar); err != nil {
			return nil, NewClientError("AccrualClient : json decoder failed")
		}
		return ar, nil
	case 204:
		return nil, NewClientError("AccrualClient : No content")
	case 429:
		msg, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, NewClientError("AccrualClient : server failed to read the request's body")
		}
		return nil, NewClientError("AccrualClient : " + string(msg))
	default:
		return nil, NewClientError("AccrualClient : internal server error (or some unknown error)")
	}
}

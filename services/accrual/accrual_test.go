package accrual

import (
	"context"
	"testing"
	"time"

	"github.com/sbxb/loyalty/storage/inmemory"
)

func TestSmth(t *testing.T) {
	t.Skip()
	store, _ := inmemory.NewMapStorage() // NewMapStorage() never returns non-nil error
	ctx := context.Background()
	accrual := NewAccrualService(store, "http://localhost:8888", ctx)
	//numbers := []string{"1149", "1156", "1172", "2238", "2253", "2279", "3327", "3376", "3384", "4416", "4457", "4481", "5512", "5538", "5587"}
	numbers := []string{"1149", "2238", "3327", "4416", "5512"}
	//numbers := []string{"1149", "2238"}
	go func() {
		for _, n := range numbers {
			time.Sleep(100 * time.Millisecond)
			accrual.AddOrderNumber(n)
		}
	}()
	accrual.ProcessNewJobQueue()
}

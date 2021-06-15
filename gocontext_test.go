package gocontext_test

import (
	"context"
	"testing"

	gocontext "github.com/omeid/go-context"
)

func addOneAfterDeadline(c *int) {
	ctx := gocontext.GoContext()

	<-ctx.Done()
	*c = *c + 1
}

func TestSimpleCancel(t *testing.T) {

	i := new(int)

	ctx, cancel := gocontext.GoWithContext(context.Background(), func() {
		deeper := func() {
			addOneAfterDeadline(i)
		}

		deeper()
	})

	if *i != 0 {
		t.Fatal("Expected i to be zero before deadline")
	}

	cancel()
	<-ctx.Done()

	if *i != 1 {
		t.Fatal("Expected i to be 1 after deadline")
	}

}

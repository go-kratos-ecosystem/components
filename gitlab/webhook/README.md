# Gitlab Dispatcher

This is a simple webhook dispatcher for Gitlab. It listens for incoming webhooks and dispatches them to the appropriate handler.

## Features

- Very convenient registration of listeners
- A single listener can implement multiple different webhook functions
- Support asynchronous and efficient processing
- Multiple Dispatch methods

## Usage Example

```go
package main

import (
	"context"
	"net/http"

	"github.com/xanzy/go-gitlab"
)

var (
	_ gitlab.BuildListener         = (*testBuildListener)(nil)
	_ gitlab.CommitCommentListener = (*testCommitCommentListener)(nil)
	_ gitlab.BuildListener         = (*testBuildAndCommitCommentListener)(nil)
	_ gitlab.CommitCommentListener = (*testBuildAndCommitCommentListener)(nil)
)

type testBuildListener struct{}

func (l *testBuildListener) OnBuild(ctx context.Context, event *gitlab.BuildEvent) error {
	// do something
	return nil
}

type testCommitCommentListener struct{}

func (l *testCommitCommentListener) OnCommitComment(ctx context.Context, event *gitlab.CommitCommentEvent) error {
	// do something
	return nil
}

type testBuildAndCommitCommentListener struct {
}

func (l *testBuildAndCommitCommentListener) OnBuild(ctx context.Context, event *gitlab.BuildEvent) error {
	// do something
	return nil
}

func (l *testBuildAndCommitCommentListener) OnCommitComment(ctx context.Context, event *gitlab.CommitCommentEvent) error {
	// do something
	return nil
}

func main() {
	dispatcher := gitlab.NewDispatcher()

	dispatcher.RegisterListeners(
		&testBuildListener{},
		&testCommitCommentListener{},
		&testBuildAndCommitCommentListener{},
	)

	http.Handle("/webhook", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := dispatcher.DispatchRequest(r,
			gitlab.DispatchRequestWithContext(context.Background()),
		); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}))

	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
```
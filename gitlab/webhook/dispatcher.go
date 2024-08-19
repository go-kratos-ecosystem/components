package webhook

import (
	"context"
	"errors"
	"io"
	"net/http"

	"github.com/xanzy/go-gitlab"
	"golang.org/x/sync/errgroup"
)

var ErrUnsupportedEvent = errors.New("gitlab/webhook: unsupported event type")

type Dispatcher struct {
	buildListeners                      []BuildListener
	commitCommentListeners              []CommitCommentListener
	deploymentListeners                 []DeploymentListener
	featureFlagListeners                []FeatureFlagListener
	groupResourceAccessTokenListeners   []GroupResourceAccessTokenListener
	issueCommentListeners               []IssueCommentListener
	issueListeners                      []IssueListener
	jobListeners                        []JobListener
	memberListeners                     []MemberListener
	mergeCommentListeners               []MergeCommentListener
	mergeListeners                      []MergeListener
	pipelineListeners                   []PipelineListener
	projectResourceAccessTokenListeners []ProjectResourceAccessTokenListener
	pushListeners                       []PushListener
	releaseListeners                    []ReleaseListener
	snippetCommentListeners             []SnippetCommentListener
	subGroupListeners                   []SubGroupListener
	tagListeners                        []TagListener
	wikiPageListeners                   []WikiPageListener
}

type Option func(*Dispatcher)

func RegisterListeners(listeners ...any) Option {
	return func(d *Dispatcher) {
		d.RegisterListeners(listeners...)
	}
}

func NewDispatcher(opts ...Option) *Dispatcher {
	dispatcher := &Dispatcher{}
	for _, opt := range opts {
		opt(dispatcher)
	}
	return dispatcher
}

func (d *Dispatcher) RegisterListeners(listeners ...any) {
	for _, listener := range listeners {
		if l, ok := listener.(BuildListener); ok {
			d.RegisterBuildListener(l)
		}

		if l, ok := listener.(CommitCommentListener); ok {
			d.RegisterCommitCommentListener(l)
		}

		if l, ok := listener.(DeploymentListener); ok {
			d.RegisterDeploymentListener(l)
		}

		if l, ok := listener.(FeatureFlagListener); ok {
			d.RegisterFeatureFlagListener(l)
		}

		if l, ok := listener.(GroupResourceAccessTokenListener); ok {
			d.RegisterGroupResourceAccessTokenListener(l)
		}

		if l, ok := listener.(IssueCommentListener); ok {
			d.RegisterIssueCommentListener(l)
		}

		if l, ok := listener.(IssueListener); ok {
			d.RegisterIssueListener(l)
		}

		if l, ok := listener.(JobListener); ok {
			d.RegisterJobListener(l)
		}

		if l, ok := listener.(MemberListener); ok {
			d.RegisterMemberListener(l)
		}

		if l, ok := listener.(MergeCommentListener); ok {
			d.RegisterMergeCommentListener(l)
		}

		if l, ok := listener.(MergeListener); ok {
			d.RegisterMergeListener(l)
		}

		if l, ok := listener.(PipelineListener); ok {
			d.RegisterPipelineListener(l)
		}

		if l, ok := listener.(ProjectResourceAccessTokenListener); ok {
			d.RegisterProjectResourceAccessTokenListener(l)
		}

		if l, ok := listener.(PushListener); ok {
			d.RegisterPushListener(l)
		}

		if l, ok := listener.(ReleaseListener); ok {
			d.RegisterReleaseListener(l)
		}

		if l, ok := listener.(SnippetCommentListener); ok {
			d.RegisterSnippetCommentListener(l)
		}

		if l, ok := listener.(SubGroupListener); ok {
			d.RegisterSubGroupListener(l)
		}

		if l, ok := listener.(TagListener); ok {
			d.RegisterTagListener(l)
		}

		if l, ok := listener.(WikiPageListener); ok {
			d.RegisterWikiPageListener(l)
		}
	}
}

func (d *Dispatcher) RegisterBuildListener(listeners ...BuildListener) {
	d.buildListeners = append(d.buildListeners, listeners...)
}

func (d *Dispatcher) RegisterCommitCommentListener(listeners ...CommitCommentListener) {
	d.commitCommentListeners = append(d.commitCommentListeners, listeners...)
}

func (d *Dispatcher) RegisterDeploymentListener(listeners ...DeploymentListener) {
	d.deploymentListeners = append(d.deploymentListeners, listeners...)
}

func (d *Dispatcher) RegisterFeatureFlagListener(listeners ...FeatureFlagListener) {
	d.featureFlagListeners = append(d.featureFlagListeners, listeners...)
}

func (d *Dispatcher) RegisterGroupResourceAccessTokenListener(listeners ...GroupResourceAccessTokenListener) {
	d.groupResourceAccessTokenListeners = append(d.groupResourceAccessTokenListeners, listeners...)
}

func (d *Dispatcher) RegisterIssueCommentListener(listeners ...IssueCommentListener) {
	d.issueCommentListeners = append(d.issueCommentListeners, listeners...)
}

func (d *Dispatcher) RegisterIssueListener(listeners ...IssueListener) {
	d.issueListeners = append(d.issueListeners, listeners...)
}

func (d *Dispatcher) RegisterJobListener(listeners ...JobListener) {
	d.jobListeners = append(d.jobListeners, listeners...)
}

func (d *Dispatcher) RegisterMemberListener(listeners ...MemberListener) {
	d.memberListeners = append(d.memberListeners, listeners...)
}

func (d *Dispatcher) RegisterMergeCommentListener(listeners ...MergeCommentListener) {
	d.mergeCommentListeners = append(d.mergeCommentListeners, listeners...)
}

func (d *Dispatcher) RegisterMergeListener(listeners ...MergeListener) {
	d.mergeListeners = append(d.mergeListeners, listeners...)
}

func (d *Dispatcher) RegisterPipelineListener(listeners ...PipelineListener) {
	d.pipelineListeners = append(d.pipelineListeners, listeners...)
}

func (d *Dispatcher) RegisterProjectResourceAccessTokenListener(listeners ...ProjectResourceAccessTokenListener) {
	d.projectResourceAccessTokenListeners = append(d.projectResourceAccessTokenListeners, listeners...)
}

func (d *Dispatcher) RegisterPushListener(listeners ...PushListener) {
	d.pushListeners = append(d.pushListeners, listeners...)
}

func (d *Dispatcher) RegisterReleaseListener(listeners ...ReleaseListener) {
	d.releaseListeners = append(d.releaseListeners, listeners...)
}

func (d *Dispatcher) RegisterSnippetCommentListener(listeners ...SnippetCommentListener) {
	d.snippetCommentListeners = append(d.snippetCommentListeners, listeners...)
}

func (d *Dispatcher) RegisterSubGroupListener(listeners ...SubGroupListener) {
	d.subGroupListeners = append(d.subGroupListeners, listeners...)
}

func (d *Dispatcher) RegisterTagListener(listeners ...TagListener) {
	d.tagListeners = append(d.tagListeners, listeners...)
}

func (d *Dispatcher) RegisterWikiPageListener(listeners ...WikiPageListener) {
	d.wikiPageListeners = append(d.wikiPageListeners, listeners...)
}

func (d *Dispatcher) Dispatch(ctx context.Context, event any) error {
	switch e := event.(type) {
	case *gitlab.BuildEvent:
		return d.processBuildEvent(ctx, e)
	case *gitlab.CommitCommentEvent:
		return d.processCommitCommentEvent(ctx, e)
	case *gitlab.DeploymentEvent:
		return d.processDeploymentEvent(ctx, e)
	case *gitlab.FeatureFlagEvent:
		return d.processFeatureFlagEvent(ctx, e)
	case *gitlab.GroupResourceAccessTokenEvent:
		return d.processGroupResourceAccessTokenEvent(ctx, e)
	case *gitlab.IssueCommentEvent:
		return d.processIssueCommentEvent(ctx, e)
	case *gitlab.IssueEvent:
		return d.processIssueEvent(ctx, e)
	case *gitlab.JobEvent:
		return d.processJobEvent(ctx, e)
	case *gitlab.MemberEvent:
		return d.processMemberEvent(ctx, e)
	case *gitlab.MergeCommentEvent:
		return d.processMergeCommentEvent(ctx, e)
	case *gitlab.MergeEvent:
		return d.processMergeEvent(ctx, e)
	case *gitlab.PipelineEvent:
		return d.processPipelineEvent(ctx, e)
	case *gitlab.ProjectResourceAccessTokenEvent:
		return d.processProjectResourceAccessTokenEvent(ctx, e)
	case *gitlab.PushEvent:
		return d.processPushEvent(ctx, e)
	case *gitlab.ReleaseEvent:
		return d.processReleaseEvent(ctx, e)
	case *gitlab.SnippetCommentEvent:
		return d.processSnippetCommentEvent(ctx, e)
	case *gitlab.SubGroupEvent:
		return d.processSubGroupEvent(ctx, e)
	case *gitlab.TagEvent:
		return d.processTagEvent(ctx, e)
	case *gitlab.WikiPageEvent:
		return d.processWikiPageEvent(ctx, e)
	default:
		return ErrUnsupportedEvent
	}
}

func (d *Dispatcher) DispatchWebhook(ctx context.Context, eventType gitlab.EventType, payload []byte) error {
	event, err := gitlab.ParseWebhook(eventType, payload)
	if err != nil {
		return err
	}
	return d.Dispatch(ctx, event)
}

type dispatchRequestOptions struct {
	ctx context.Context
}

type DispatchRequestOption func(*dispatchRequestOptions)

func DispatchRequestWithContext(ctx context.Context) DispatchRequestOption {
	return func(o *dispatchRequestOptions) {
		o.ctx = ctx
	}
}

func (d *Dispatcher) DispatchRequest(req *http.Request, opts ...DispatchRequestOption) error {
	o := &dispatchRequestOptions{
		ctx: req.Context(),
	}
	for _, opt := range opts {
		opt(o)
	}
	payload, err := io.ReadAll(req.Body)
	if err != nil {
		return err
	}
	return d.DispatchWebhook(o.ctx, gitlab.HookEventType(req), payload)
}

func (d *Dispatcher) processBuildEvent(ctx context.Context, event *gitlab.BuildEvent) error {
	eg, ctx := errgroup.WithContext(ctx)
	for _, listener := range d.buildListeners {
		eg.Go(func() error {
			return listener.OnBuild(ctx, event)
		})
	}
	return eg.Wait()
}

func (d *Dispatcher) processCommitCommentEvent(ctx context.Context, event *gitlab.CommitCommentEvent) error {
	eg, ctx := errgroup.WithContext(ctx)
	for _, listener := range d.commitCommentListeners {
		eg.Go(func() error {
			return listener.OnCommitComment(ctx, event)
		})
	}
	return eg.Wait()
}

func (d *Dispatcher) processDeploymentEvent(ctx context.Context, event *gitlab.DeploymentEvent) error {
	eg, ctx := errgroup.WithContext(ctx)
	for _, listener := range d.deploymentListeners {
		eg.Go(func() error {
			return listener.OnDeployment(ctx, event)
		})
	}
	return eg.Wait()
}

func (d *Dispatcher) processFeatureFlagEvent(ctx context.Context, event *gitlab.FeatureFlagEvent) error {
	eg, ctx := errgroup.WithContext(ctx)
	for _, listener := range d.featureFlagListeners {
		eg.Go(func() error {
			return listener.OnFeatureFlag(ctx, event)
		})
	}
	return eg.Wait()
}

func (d *Dispatcher) processGroupResourceAccessTokenEvent(ctx context.Context, event *gitlab.GroupResourceAccessTokenEvent) error { //nolint:lll
	eg, ctx := errgroup.WithContext(ctx)
	for _, listener := range d.groupResourceAccessTokenListeners {
		eg.Go(func() error {
			return listener.OnGroupResourceAccessToken(ctx, event)
		})
	}
	return eg.Wait()
}

func (d *Dispatcher) processIssueCommentEvent(ctx context.Context, event *gitlab.IssueCommentEvent) error {
	eg, ctx := errgroup.WithContext(ctx)
	for _, listener := range d.issueCommentListeners {
		eg.Go(func() error {
			return listener.OnIssueComment(ctx, event)
		})
	}
	return eg.Wait()
}

func (d *Dispatcher) processIssueEvent(ctx context.Context, event *gitlab.IssueEvent) error {
	eg, ctx := errgroup.WithContext(ctx)
	for _, listener := range d.issueListeners {
		eg.Go(func() error {
			return listener.OnIssue(ctx, event)
		})
	}
	return eg.Wait()
}

func (d *Dispatcher) processJobEvent(ctx context.Context, event *gitlab.JobEvent) error {
	eg, ctx := errgroup.WithContext(ctx)
	for _, listener := range d.jobListeners {
		eg.Go(func() error {
			return listener.OnJob(ctx, event)
		})
	}
	return eg.Wait()
}

func (d *Dispatcher) processMemberEvent(ctx context.Context, event *gitlab.MemberEvent) error {
	eg, ctx := errgroup.WithContext(ctx)
	for _, listener := range d.memberListeners {
		eg.Go(func() error {
			return listener.OnMember(ctx, event)
		})
	}
	return eg.Wait()
}

func (d *Dispatcher) processMergeCommentEvent(ctx context.Context, event *gitlab.MergeCommentEvent) error {
	eg, ctx := errgroup.WithContext(ctx)
	for _, listener := range d.mergeCommentListeners {
		eg.Go(func() error {
			return listener.OnMergeComment(ctx, event)
		})
	}
	return eg.Wait()
}

func (d *Dispatcher) processMergeEvent(ctx context.Context, event *gitlab.MergeEvent) error {
	eg, ctx := errgroup.WithContext(ctx)
	for _, listener := range d.mergeListeners {
		eg.Go(func() error {
			return listener.OnMerge(ctx, event)
		})
	}
	return eg.Wait()
}

func (d *Dispatcher) processPipelineEvent(ctx context.Context, event *gitlab.PipelineEvent) error {
	eg, ctx := errgroup.WithContext(ctx)
	for _, listener := range d.pipelineListeners {
		eg.Go(func() error {
			return listener.OnPipeline(ctx, event)
		})
	}
	return eg.Wait()
}

func (d *Dispatcher) processProjectResourceAccessTokenEvent(ctx context.Context, event *gitlab.ProjectResourceAccessTokenEvent) error { //nolint:lll
	eg, ctx := errgroup.WithContext(ctx)
	for _, listener := range d.projectResourceAccessTokenListeners {
		eg.Go(func() error {
			return listener.OnProjectResourceAccessToken(ctx, event)
		})
	}
	return eg.Wait()
}

func (d *Dispatcher) processPushEvent(ctx context.Context, event *gitlab.PushEvent) error {
	eg, ctx := errgroup.WithContext(ctx)
	for _, listener := range d.pushListeners {
		eg.Go(func() error {
			return listener.OnPush(ctx, event)
		})
	}
	return eg.Wait()
}

func (d *Dispatcher) processReleaseEvent(ctx context.Context, event *gitlab.ReleaseEvent) error {
	eg, ctx := errgroup.WithContext(ctx)
	for _, listener := range d.releaseListeners {
		eg.Go(func() error {
			return listener.OnRelease(ctx, event)
		})
	}
	return eg.Wait()
}

func (d *Dispatcher) processSnippetCommentEvent(ctx context.Context, event *gitlab.SnippetCommentEvent) error {
	eg, ctx := errgroup.WithContext(ctx)
	for _, listener := range d.snippetCommentListeners {
		eg.Go(func() error {
			return listener.OnSnippetComment(ctx, event)
		})
	}
	return eg.Wait()
}

func (d *Dispatcher) processSubGroupEvent(ctx context.Context, event *gitlab.SubGroupEvent) error {
	eg, ctx := errgroup.WithContext(ctx)
	for _, listener := range d.subGroupListeners {
		eg.Go(func() error {
			return listener.OnSubGroup(ctx, event)
		})
	}
	return eg.Wait()
}

func (d *Dispatcher) processTagEvent(ctx context.Context, event *gitlab.TagEvent) error {
	eg, ctx := errgroup.WithContext(ctx)
	for _, listener := range d.tagListeners {
		eg.Go(func() error {
			return listener.OnTag(ctx, event)
		})
	}
	return eg.Wait()
}

func (d *Dispatcher) processWikiPageEvent(ctx context.Context, event *gitlab.WikiPageEvent) error {
	eg, ctx := errgroup.WithContext(ctx)
	for _, listener := range d.wikiPageListeners {
		eg.Go(func() error {
			return listener.OnWikiPage(ctx, event)
		})
	}
	return eg.Wait()
}

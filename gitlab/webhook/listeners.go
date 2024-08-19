package gitlab

import (
	"context"

	"github.com/xanzy/go-gitlab"
)

type BuildListener interface {
	OnBuild(ctx context.Context, event *gitlab.BuildEvent) error
}

type CommitCommentListener interface {
	OnCommitComment(ctx context.Context, event *gitlab.CommitCommentEvent) error
}

type DeploymentListener interface {
	OnDeployment(ctx context.Context, event *gitlab.DeploymentEvent) error
}

type FeatureFlagListener interface {
	OnFeatureFlag(ctx context.Context, event *gitlab.FeatureFlagEvent) error
}

type GroupResourceAccessTokenListener interface {
	OnGroupResourceAccessToken(ctx context.Context, event *gitlab.GroupResourceAccessTokenEvent) error
}

type IssueCommentListener interface {
	OnIssueComment(ctx context.Context, event *gitlab.IssueCommentEvent) error
}

type IssueListener interface {
	OnIssue(ctx context.Context, event *gitlab.IssueEvent) error
}

type JobListener interface {
	OnJob(ctx context.Context, event *gitlab.JobEvent) error
}

type MemberListener interface {
	OnMember(ctx context.Context, event *gitlab.MemberEvent) error
}

type MergeCommentListener interface {
	OnMergeComment(ctx context.Context, event *gitlab.MergeCommentEvent) error
}

type MergeListener interface {
	OnMerge(ctx context.Context, event *gitlab.MergeEvent) error
}

type PipelineListener interface {
	OnPipeline(ctx context.Context, event *gitlab.PipelineEvent) error
}

type ProjectResourceAccessTokenListener interface {
	OnProjectResourceAccessToken(ctx context.Context, event *gitlab.ProjectResourceAccessTokenEvent) error
}

type PushListener interface {
	OnPush(ctx context.Context, event *gitlab.PushEvent) error
}

type ReleaseListener interface {
	OnRelease(ctx context.Context, event *gitlab.ReleaseEvent) error
}

type SnippetCommentListener interface {
	OnSnippetComment(ctx context.Context, event *gitlab.SnippetCommentEvent) error
}

type SubGroupListener interface {
	OnSubGroup(ctx context.Context, event *gitlab.SubGroupEvent) error
}

type TagListener interface {
	OnTag(ctx context.Context, event *gitlab.TagEvent) error
}

type WikiPageListener interface {
	OnWikiPage(ctx context.Context, event *gitlab.WikiPageEvent) error
}

package prints

import "github.com/cheggaaa/pb/v3"

var (
	Full    = pb.Full
	Default = pb.Default
	Simple  = pb.Simple
)

type progressBarOptions struct {
	template pb.ProgressBarTemplate
}

type ProgressBarOption func(*progressBarOptions)

func WithTemplate(template pb.ProgressBarTemplate) ProgressBarOption {
	return func(o *progressBarOptions) {
		o.template = template
	}
}

func NewProgressBar(total int, opts ...ProgressBarOption) *pb.ProgressBar {
	opt := &progressBarOptions{}

	for _, o := range opts {
		o(opt)
	}

	if opt.template != "" {
		return opt.template.Start(total)
	}

	return pb.StartNew(total)
}

package gitlab

import (
	"fmt"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/xanzy/go-gitlab"
)

func WithListOptions(options *gitlab.ListOptions) gitlab.RequestOptionFunc {
	return func(req *retryablehttp.Request) error {
		q := req.URL.Query()
		q.Set("page", fmt.Sprintf("%d", options.Page))
		q.Set("per_page", fmt.Sprintf("%d", options.PerPage))
		req.URL.RawQuery = q.Encode()
		return nil
	}
}

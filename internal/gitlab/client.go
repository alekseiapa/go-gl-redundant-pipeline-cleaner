package gitlab

import (
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/alekseiapa/go-gl-redundant-pipeline-cleaner/internal/config"
	"github.com/alekseiapa/go-gl-redundant-pipeline-cleaner/internal/utils"
	"github.com/xanzy/go-gitlab"
)

type GitlabClient struct {
	client   *gitlab.Client
	project  *gitlab.Project
	logger   *log.Logger
	settings *config.Config
}

func NewGitlabClient(cfg *config.Config, logger *log.Logger) (*GitlabClient, error) {
	glClient, err := gitlab.NewClient(cfg.GitlabAPIToken, gitlab.WithBaseURL(cfg.GitlabURL))
	if err != nil {
		return nil, err
	}
	project, _, err := glClient.Projects.GetProject(cfg.GitlabProjectID, nil)
	if err != nil {
		return nil, err
	}

	return &GitlabClient{
		client:   glClient,
		project:  project,
		logger:   logger,
		settings: cfg,
	}, nil
}

func (gc *GitlabClient) ListPipelinesByMR(mrId int) ([]*gitlab.PipelineInfo, error) {
	var allPipelines []*gitlab.PipelineInfo
	page := 1
	perPage := 100
	for {
		// Fetch a page of pipelines
		pipelines, resp, err := gc.client.MergeRequests.ListMergeRequestPipelines(gc.project.ID, mrId,
			WithListOptions(&gitlab.ListOptions{
				Page:    page,
				PerPage: perPage,
			}),
		)
		if err != nil {
			return nil, err
		}

		allPipelines = append(allPipelines, pipelines...)

		// Check if we have fetched all pages
		if resp.NextPage == 0 {
			break
		}
		page = resp.NextPage
	}
	gc.logger.Printf("Found %d pipelines for MR %d", len(allPipelines), mrId)

	return allPipelines, nil
}

func (gc *GitlabClient) CancelRedundantPipelinesByMR(mrId int, mrAction string) error {
	// Allow the pipeline to be created
	time.Sleep(20 * time.Second)

	pipelines, err := gc.ListPipelinesByMR(mrId)
	if err != nil {
		return fmt.Errorf("failed to fetch the pipelines for MR %d: %v\n", mrId, err)
	}

	excludedPipelineStatuses := map[string]bool{
		"success":   true,
		"failed":    true,
		"canceled":  true,
		"skipped":   true,
		"scheduled": true,
	}

	// Filter out the pipelines based on their status
	var targetPipelines []*gitlab.PipelineInfo
	for _, pipeline := range pipelines {
		if !excludedPipelineStatuses[pipeline.Status] {
			targetPipelines = append(targetPipelines, pipeline)
		}
	}

	// Sort the pipelines by their ID in Desc Order. The newest ones go first.
	sort.Slice(targetPipelines, func(i, j int) bool {
		return targetPipelines[i].ID > targetPipelines[j].ID
	})

	var redundantPipelines []*gitlab.PipelineInfo
	if mrAction == "update" {
		if len(targetPipelines) > 1 {
			redundantPipelines = targetPipelines[1:]
		}
	} else {
		// Cancel all the pipelines for the MR action 'close'
		redundantPipelines = targetPipelines
	}

	for _, pipeline := range redundantPipelines {
		err := utils.Retry(3, 4*time.Second, func() error {
			pipelineID := pipeline.ID
			_, _, cancelErr := gc.client.Pipelines.CancelPipelineBuild(gc.project.ID, pipelineID)
			if cancelErr != nil {
				gc.logger.Printf("Error cancelling pipeline %d: %v\n", pipeline.ID, cancelErr)
			} else {
				gc.logger.Printf("Successfully cancelled pipeline ID %d for the MR %d\n", pipeline.ID, mrId)
			}
			return cancelErr
		})
		if err != nil {
			gc.logger.Printf("MR ID %d: Max retries reached for pipeline %d. Error: %v\n", mrId, pipeline.ID, err)
		}
	}
	return nil
}

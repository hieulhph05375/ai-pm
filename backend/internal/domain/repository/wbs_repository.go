package repository

import (
	"context"
	"project-mgmt/backend/internal/domain/entity"
)

type WBSRepository interface {
	GetProjectTree(ctx context.Context, projectID int, filter entity.WBSFilter) ([]entity.WBSNode, int, error)
	GetChildren(ctx context.Context, basePath string) ([]entity.WBSNode, error)
	GetImmediateChildren(ctx context.Context, projectID int, parentPath string) ([]entity.WBSNode, error)
	CreateNode(ctx context.Context, node *entity.WBSNode, parentPath string) error
	UpdateNode(ctx context.Context, node *entity.WBSNode) error
	DeleteNode(ctx context.Context, id int) error
	GetNodeByID(ctx context.Context, id int) (*entity.WBSNode, error)
	GetNodeByPath(ctx context.Context, projectID int, path string) (*entity.WBSNode, error)
	ListDependencies(ctx context.Context, projectID int) ([]entity.WBSDependency, error)
	CreateDependency(ctx context.Context, dep *entity.WBSDependency) error
	DeleteDependency(ctx context.Context, depID int) error

	// Comment methods
	AddComment(ctx context.Context, comment *entity.WBSComment) error
	ListComments(ctx context.Context, nodeID int) ([]entity.WBSComment, error)
	DeleteComment(ctx context.Context, commentID int) error
	UpdateComment(ctx context.Context, commentID int, content string) error
	GetCommentByID(ctx context.Context, commentID int) (*entity.WBSComment, error)

	// Baseline methods
	CreateBaseline(ctx context.Context, baseline *entity.WBSBaseline) error
	CopyNodesToBaseline(ctx context.Context, projectID int, baselineID int) error
	GetBaselines(ctx context.Context, projectID int) ([]entity.WBSBaseline, error)
	GetBaselineNodes(ctx context.Context, baselineID int) ([]entity.WBSBaselineNode, error)
}

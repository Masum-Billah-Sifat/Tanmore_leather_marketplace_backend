// ------------------------------------------------------------
// 📁 File: internal/services/category/get_category_tree_service.go
// 🧠 Service to fetch and build full category tree structure.
//
//     - Fetches all non-archived categories
//     - Builds nested tree structure using parent-child mapping

package category

import (
	"context"

	"tanmore_backend/internal/db/sqlc"
	repo "tanmore_backend/internal/repository/category/category_tree"
	"tanmore_backend/pkg/errors"

	"github.com/google/uuid"
)

// ------------------------------------------------------------
// 📤 Result Node Structure
type CategoryNode struct {
	ID       uuid.UUID      `json:"id"`
	Name     string         `json:"name"`
	Slug     string         `json:"slug"`
	Level    int            `json:"level"`
	IsLeaf   bool           `json:"is_leaf"`
	Children []CategoryNode `json:"children"`
}

// ------------------------------------------------------------
// 📤 Result Object
type GetCategoryTreeResult struct {
	Status string         `json:"status"`
	Tree   []CategoryNode `json:"data"`
}

// ------------------------------------------------------------
// 🧱 Dependencies
type GetCategoryTreeServiceDeps struct {
	Repo repo.CategoryTreeRepoInterface
}

// ------------------------------------------------------------
// 🛠️ Service Definition
type GetCategoryTreeService struct {
	Deps GetCategoryTreeServiceDeps
}

// 🚀 Constructor
func NewGetCategoryTreeService(deps GetCategoryTreeServiceDeps) *GetCategoryTreeService {
	return &GetCategoryTreeService{Deps: deps}
}

// 🚀 Entrypoint
func (s *GetCategoryTreeService) Start(ctx context.Context) (*GetCategoryTreeResult, error) {
	rows, err := s.Deps.Repo.GetAllNonArchivedCategories(ctx)
	if err != nil {
		return nil, errors.NewTableError("categories.select", err.Error())
	}

	// Step 1: Convert rows into usable slice
	var flatCategories []sqlc.GetAllNonArchivedCategoriesRow
	flatCategories = append(flatCategories, rows...)

	// Step 2: Create parent -> children map
	parentMap := make(map[uuid.UUID][]sqlc.GetAllNonArchivedCategoriesRow)
	for _, row := range flatCategories {
		if row.ParentID.Valid {
			parentMap[row.ParentID.UUID] = append(parentMap[row.ParentID.UUID], row)
		}
	}

	// Step 3: Build tree recursively
	var tree []CategoryNode
	for _, row := range flatCategories {
		if !row.ParentID.Valid {
			tree = append(tree, buildTreeNode(row, parentMap))
		}
	}

	return &GetCategoryTreeResult{
		Status: "success",
		Tree:   tree,
	}, nil
}

// ------------------------------------------------------------
// 🔁 Recursive Helper
func buildTreeNode(
	cat sqlc.GetAllNonArchivedCategoriesRow,
	parentMap map[uuid.UUID][]sqlc.GetAllNonArchivedCategoriesRow,
) CategoryNode {
	node := CategoryNode{
		ID:       cat.ID,
		Name:     cat.Name,
		Slug:     cat.Slug,
		Level:    int(cat.Level),
		IsLeaf:   cat.IsLeaf,
		Children: []CategoryNode{},
	}

	for _, child := range parentMap[cat.ID] {
		node.Children = append(node.Children, buildTreeNode(child, parentMap))
	}

	return node
}

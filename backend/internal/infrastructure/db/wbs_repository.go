package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"project-mgmt/backend/internal/domain/entity"
	"project-mgmt/backend/internal/domain/repository"
	"strings"
)

type wbsRepo struct {
	db *sql.DB
}

func NewWBSRepository(db *sql.DB) repository.WBSRepository {
	return &wbsRepo{db: db}
}

func (r *wbsRepo) getFieldPoints(n *entity.WBSNode, fields []string) []interface{} {
	allFields := map[string]interface{}{
		"id":                 &n.ID,
		"project_id":         &n.ProjectID,
		"title":              &n.Title,
		"type":               &n.Type,
		"type_id":            &n.TypeID,
		"path":               &n.Path,
		"order_index":        &n.OrderIndex,
		"planned_start_date": &n.PlannedStartDate,
		"planned_end_date":   &n.PlannedEndDate,
		"actual_start_date":  &n.ActualStartDate,
		"actual_end_date":    &n.ActualEndDate,
		"progress":           &n.Progress,
		"planned_value":      &n.PlannedValue,
		"actual_cost":        &n.ActualCost,
		"estimated_effort":   &n.EstimatedEffort,
		"actual_effort":      &n.ActualEffort,
		"assigned_to":        &n.AssignedTo,
		"description":        &n.Description,
		"created_at":         &n.CreatedAt,
		"updated_at":         &n.UpdatedAt,
		"has_children":       &n.HasChildren,
	}

	if len(fields) == 0 {
		return []interface{}{
			&n.ID, &n.ProjectID, &n.Title, &n.Type, &n.TypeID, &n.Path, &n.OrderIndex,
			&n.PlannedStartDate, &n.PlannedEndDate, &n.ActualStartDate, &n.ActualEndDate,
			&n.Progress, &n.PlannedValue, &n.ActualCost, &n.EstimatedEffort, &n.ActualEffort, &n.AssignedTo, &n.Description, &n.CreatedAt, &n.UpdatedAt,
			&n.HasChildren,
		}
	}

	var points []interface{}
	for _, f := range fields {
		if p, ok := allFields[f]; ok {
			points = append(points, p)
		} else {
			// fallback for essential fields if wrongly requested
			log.Printf("Warning: unknown field %s", f)
		}
	}
	return points
}

func (r *wbsRepo) scanNodesFiltered(rows *sql.Rows, fields []string) ([]entity.WBSNode, error) {
	var nodes []entity.WBSNode
	for rows.Next() {
		var n entity.WBSNode
		var typeName, typeColor, typeIcon sql.NullString
		dest := r.getFieldPoints(&n, fields)
		// Metadata columns are ALWAYS returned by the JOIN in GetProjectTree
		dest = append(dest, &typeName, &typeColor, &typeIcon)

		if err := rows.Scan(dest...); err != nil {
			return nil, err
		}

		if n.TypeID != nil && typeName.Valid {
			n.TypeCat = &entity.Category{
				ID:    *n.TypeID,
				Name:  typeName.String,
				Color: nullStringPtr(typeColor),
				Icon:  nullStringPtr(typeIcon),
			}
		}
		nodes = append(nodes, n)
	}
	if nodes == nil {
		nodes = []entity.WBSNode{}
	}
	return nodes, rows.Err()
}

func (r *wbsRepo) scanNode(scanner interface{ Scan(...interface{}) error }, n *entity.WBSNode) error {
	return scanner.Scan(
		&n.ID, &n.ProjectID, &n.Title, &n.Type, &n.Path, &n.OrderIndex,
		&n.PlannedStartDate, &n.PlannedEndDate, &n.ActualStartDate, &n.ActualEndDate,
		&n.Progress, &n.PlannedValue, &n.ActualCost, &n.EstimatedEffort, &n.ActualEffort, &n.AssignedTo, &n.Description, &n.CreatedAt, &n.UpdatedAt,
		&n.HasChildren, &n.TypeID,
	)
}

func (r *wbsRepo) scanNodes(rows *sql.Rows) ([]entity.WBSNode, error) {
	var nodes []entity.WBSNode
	for rows.Next() {
		var n entity.WBSNode
		var typeName, typeColor, typeIcon sql.NullString
		err := rows.Scan(
			&n.ID, &n.ProjectID, &n.Title, &n.Type, &n.TypeID, &n.Path, &n.OrderIndex,
			&n.PlannedStartDate, &n.PlannedEndDate, &n.ActualStartDate, &n.ActualEndDate,
			&n.Progress, &n.PlannedValue, &n.ActualCost, &n.EstimatedEffort, &n.ActualEffort, &n.AssignedTo, &n.Description, &n.CreatedAt, &n.UpdatedAt,
			&n.HasChildren, &typeName, &typeColor, &typeIcon,
		)
		if err != nil {
			return nil, err
		}

		if n.TypeID != nil && typeName.Valid {
			n.TypeCat = &entity.Category{
				ID:    *n.TypeID,
				Name:  typeName.String,
				Color: nullStringPtr(typeColor),
				Icon:  nullStringPtr(typeIcon),
			}
		}
		nodes = append(nodes, n)
	}
	if nodes == nil {
		nodes = []entity.WBSNode{}
	}
	return nodes, rows.Err()
}

func (r *wbsRepo) GetProjectTree(ctx context.Context, projectID int, filter entity.WBSFilter) ([]entity.WBSNode, int, error) {
	// Prepare field selection
	// Optimization: Use a wrapping query for has_children calculation to avoid
	// running the subquery for every intermediate match during search/filtering.
	innerCols := "id, project_id, title, type, type_id, path::text, order_index, planned_start_date, planned_end_date, actual_start_date, actual_end_date, progress, planned_value, actual_cost, estimated_effort, actual_effort, assigned_to, description, created_at, updated_at"
	outerCols := "filtered.id, filtered.project_id, filtered.title, filtered.type, filtered.type_id, filtered.path, filtered.order_index, filtered.planned_start_date, filtered.planned_end_date, filtered.actual_start_date, filtered.actual_end_date, filtered.progress, filtered.planned_value, filtered.actual_cost, filtered.estimated_effort, filtered.actual_effort, filtered.assigned_to, filtered.description, filtered.created_at, filtered.updated_at"

	selectCols := fmt.Sprintf("%s, EXISTS(SELECT 1 FROM wbs_nodes sub WHERE sub.project_id = filtered.project_id AND sub.path <@ filtered.path::ltree AND sub.path != filtered.path::ltree) as has_children", outerCols)

	if len(filter.Fields) > 0 {
		var validFields []string
		for _, f := range filter.Fields {
			if f == "path" {
				validFields = append(validFields, "filtered.path")
			} else if f == "has_children" {
				// We keep this placeholder to handle it in the outer query
				validFields = append(validFields, "EXISTS(SELECT 1 FROM wbs_nodes sub WHERE sub.project_id = filtered.project_id AND sub.path <@ filtered.path::ltree AND sub.path != filtered.path::ltree) as has_children")
			} else {
				validFields = append(validFields, "filtered."+f)
			}
		}
		selectCols = strings.Join(validFields, ", ")
	}

	var conditions []string
	var args []interface{}
	args = append(args, projectID)
	placeholderIdx := 2

	// If search filters are provided, we use the EXISTS logic to return the structural tree
	hasSearchFilters := filter.Search != "" || filter.AssignedTo != nil || filter.Status != ""

	if filter.Search != "" {
		conditions = append(conditions, fmt.Sprintf("(m.title ILIKE $%d OR m.path::text ILIKE $%d)", placeholderIdx, placeholderIdx+1))
		args = append(args, "%"+filter.Search+"%", "%"+filter.Search+"%")
		placeholderIdx += 2
	}

	if filter.AssignedTo != nil {
		conditions = append(conditions, fmt.Sprintf("m.assigned_to = $%d", placeholderIdx))
		args = append(args, *filter.AssignedTo)
		placeholderIdx++
	}

	if filter.Status != "" {
		switch filter.Status {
		case "todo":
			conditions = append(conditions, "m.progress = 0")
		case "doing":
			conditions = append(conditions, "m.progress > 0 AND m.progress < 100")
		case "done":
			conditions = append(conditions, "m.progress = 100")
		}
	}

	// Calculate LIMIT and OFFSET
	limitOffsetSql := ""
	if filter.Limit > 0 {
		page := filter.Page
		if page < 1 {
			page = 1
		}
		offset := (page - 1) * filter.Limit
		limitOffsetSql = fmt.Sprintf(" LIMIT %d OFFSET %d", filter.Limit, offset)
	}

	// 1. Case: Full Tree Search (EXISTS logic)
	if hasSearchFilters {
		whereClause := "AND " + strings.Join(conditions, " AND ")

		// Count Query
		countQuery := fmt.Sprintf(`SELECT COUNT(DISTINCT wbs_nodes.id) FROM wbs_nodes WHERE wbs_nodes.project_id = $1 AND EXISTS (SELECT 1 FROM wbs_nodes m WHERE m.project_id = $1 AND m.path <@ wbs_nodes.path %s)`, whereClause)
		var total int
		if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
			return nil, 0, err
		}

		query := fmt.Sprintf(`
			SELECT %s, c.name as type_name, c.color as type_color, c.icon as type_icon
			FROM (
				SELECT %s
				FROM wbs_nodes
				WHERE wbs_nodes.project_id = $1 
				AND EXISTS (
					SELECT 1 FROM wbs_nodes m 
					WHERE m.project_id = $1 
					AND m.path <@ wbs_nodes.path 
					%s
				)
				ORDER BY wbs_nodes.path ASC
				%s
			) as filtered
			LEFT JOIN categories c ON filtered.type_id = c.id
			ORDER BY filtered.path ASC`, selectCols, innerCols, whereClause, limitOffsetSql)

		rows, err := r.db.QueryContext(ctx, query, args...)
		if err != nil {
			return nil, 0, err
		}
		defer rows.Close()
		nodes, err := r.scanNodesFiltered(rows, filter.Fields)
		return nodes, total, err
	}

	// 2. Case: Lazy Loading (Incremental fetch)
	var query string
	var countQuery string
	var total int

	if filter.FetchAll {
		countQuery = `SELECT COUNT(id) FROM wbs_nodes WHERE project_id = $1`
		if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
			return nil, 0, err
		}

		query = fmt.Sprintf(`
			SELECT %s, c.name as type_name, c.color as type_color, c.icon as type_icon
			FROM (
				SELECT %s
				FROM wbs_nodes
				WHERE project_id = $1 
				ORDER BY order_index ASC
			) as filtered
			LEFT JOIN categories c ON filtered.type_id = c.id
			ORDER BY filtered.order_index ASC`, selectCols, innerCols)

		// Bypass limitOffsetSql
		limitOffsetSql = ""
	} else if filter.ParentPath != "" {
		countQuery = fmt.Sprintf(`SELECT COUNT(id) FROM wbs_nodes WHERE project_id = $1 AND path <@ $%d::ltree AND nlevel(path) = nlevel($%d::ltree) + 1`, placeholderIdx, placeholderIdx)
		if err := r.db.QueryRowContext(ctx, countQuery, projectID, filter.ParentPath).Scan(&total); err != nil {
			return nil, 0, err
		}

		query = fmt.Sprintf(`
			SELECT %s, c.name as type_name, c.color as type_color, c.icon as type_icon
			FROM (
				SELECT %s
				FROM wbs_nodes
				WHERE project_id = $1 
				AND path <@ $%d::ltree 
				AND nlevel(path) = nlevel($%d::ltree) + 1
				ORDER BY order_index ASC
				%s
			) as filtered
			LEFT JOIN categories c ON filtered.type_id = c.id
			ORDER BY filtered.order_index ASC`, selectCols, innerCols, placeholderIdx, placeholderIdx, limitOffsetSql)
		args = append(args, filter.ParentPath)
	} else {
		countQuery = `SELECT COUNT(id) FROM wbs_nodes WHERE project_id = $1 AND nlevel(path) = 1`
		if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
			return nil, 0, err
		}

		query = fmt.Sprintf(`
			SELECT %s, c.name as type_name, c.color as type_color, c.icon as type_icon
			FROM (
				SELECT %s
				FROM wbs_nodes
				WHERE project_id = $1 
				AND nlevel(path) = 1
				ORDER BY order_index ASC
				%s
			) as filtered
			LEFT JOIN categories c ON filtered.type_id = c.id
			ORDER BY filtered.order_index ASC`, selectCols, innerCols, limitOffsetSql)
	}

	query += limitOffsetSql

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	nodes, err := r.scanNodesFiltered(rows, filter.Fields)
	return nodes, total, err
}

func (r *wbsRepo) GetChildren(ctx context.Context, basePath string) ([]entity.WBSNode, error) {
	query := `SELECT wbs_nodes.id, wbs_nodes.project_id, wbs_nodes.title, wbs_nodes.type, wbs_nodes.type_id, wbs_nodes.path::text, wbs_nodes.order_index, 
	                 wbs_nodes.planned_start_date, wbs_nodes.planned_end_date, wbs_nodes.actual_start_date, wbs_nodes.actual_end_date,
					 wbs_nodes.progress, wbs_nodes.planned_value, wbs_nodes.actual_cost, wbs_nodes.estimated_effort, wbs_nodes.actual_effort, wbs_nodes.assigned_to, wbs_nodes.description, wbs_nodes.created_at, wbs_nodes.updated_at,
					 EXISTS(SELECT 1 FROM wbs_nodes sub WHERE sub.project_id = wbs_nodes.project_id AND sub.path <@ wbs_nodes.path AND sub.path != wbs_nodes.path) as has_children,
					 c.name as type_name, c.color as type_color, c.icon as type_icon
			  FROM wbs_nodes
			  LEFT JOIN categories c ON wbs_nodes.type_id = c.id
			  WHERE path <@ $1::ltree AND path != $1::ltree
			  ORDER BY path ASC`

	rows, err := r.db.QueryContext(ctx, query, basePath)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanNodes(rows)
}

func (r *wbsRepo) GetImmediateChildren(ctx context.Context, projectID int, parentPath string) ([]entity.WBSNode, error) {
	var query string
	var args []interface{}

	if parentPath == "" {
		query = `SELECT wbs_nodes.id, wbs_nodes.project_id, wbs_nodes.title, wbs_nodes.type, wbs_nodes.type_id, wbs_nodes.path::text, wbs_nodes.order_index, 
	                 wbs_nodes.planned_start_date, wbs_nodes.planned_end_date, wbs_nodes.actual_start_date, wbs_nodes.actual_end_date,
					 wbs_nodes.progress, wbs_nodes.planned_value, wbs_nodes.actual_cost, wbs_nodes.estimated_effort, wbs_nodes.actual_effort, wbs_nodes.assigned_to, wbs_nodes.description, wbs_nodes.created_at, wbs_nodes.updated_at,
					 EXISTS(SELECT 1 FROM wbs_nodes sub WHERE sub.project_id = wbs_nodes.project_id AND sub.path <@ wbs_nodes.path AND sub.path != wbs_nodes.path) as has_children,
					 c.name as type_name, c.color as type_color, c.icon as type_icon
			  FROM wbs_nodes
			  LEFT JOIN categories c ON wbs_nodes.type_id = c.id
			  WHERE project_id = $1 AND nlevel(path) = 1
			  ORDER BY order_index ASC`
		args = append(args, projectID)
	} else {
		query = `SELECT wbs_nodes.id, wbs_nodes.project_id, wbs_nodes.title, wbs_nodes.type, wbs_nodes.type_id, wbs_nodes.path::text, wbs_nodes.order_index, 
	                 wbs_nodes.planned_start_date, wbs_nodes.planned_end_date, wbs_nodes.actual_start_date, wbs_nodes.actual_end_date,
					 wbs_nodes.progress, wbs_nodes.planned_value, wbs_nodes.actual_cost, wbs_nodes.estimated_effort, wbs_nodes.actual_effort, wbs_nodes.assigned_to, wbs_nodes.description, wbs_nodes.created_at, wbs_nodes.updated_at,
					 EXISTS(SELECT 1 FROM wbs_nodes sub WHERE sub.project_id = wbs_nodes.project_id AND sub.path <@ wbs_nodes.path AND sub.path != wbs_nodes.path) as has_children,
					 c.name as type_name, c.color as type_color, c.icon as type_icon
			  FROM wbs_nodes
			  LEFT JOIN categories c ON wbs_nodes.type_id = c.id
			  WHERE project_id = $1 AND path <@ $2::ltree AND nlevel(path) = nlevel($2::ltree) + 1
			  ORDER BY order_index ASC`
		args = append(args, projectID, parentPath)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanNodes(rows)
}

func (r *wbsRepo) GetNodeByID(ctx context.Context, id int) (*entity.WBSNode, error) {
	n := &entity.WBSNode{}
	var typeName, typeColor, typeIcon sql.NullString
	query := `SELECT wbs_nodes.id, wbs_nodes.project_id, wbs_nodes.title, wbs_nodes.type, wbs_nodes.type_id, wbs_nodes.path::text, wbs_nodes.order_index, 
	                 wbs_nodes.planned_start_date, wbs_nodes.planned_end_date, wbs_nodes.actual_start_date, wbs_nodes.actual_end_date,
					 wbs_nodes.progress, wbs_nodes.planned_value, wbs_nodes.actual_cost, wbs_nodes.estimated_effort, wbs_nodes.actual_effort, wbs_nodes.assigned_to, wbs_nodes.description, wbs_nodes.created_at, wbs_nodes.updated_at,
					 EXISTS(SELECT 1 FROM wbs_nodes sub WHERE sub.project_id = wbs_nodes.project_id AND sub.path <@ wbs_nodes.path AND sub.path != wbs_nodes.path) as has_children,
					 c.name as type_name, c.color as type_color, c.icon as type_icon
			  FROM wbs_nodes
			  LEFT JOIN categories c ON wbs_nodes.type_id = c.id
			  WHERE wbs_nodes.id = $1`

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&n.ID, &n.ProjectID, &n.Title, &n.Type, &n.TypeID, &n.Path, &n.OrderIndex,
		&n.PlannedStartDate, &n.PlannedEndDate, &n.ActualStartDate, &n.ActualEndDate,
		&n.Progress, &n.PlannedValue, &n.ActualCost, &n.EstimatedEffort, &n.ActualEffort, &n.AssignedTo, &n.Description, &n.CreatedAt, &n.UpdatedAt,
		&n.HasChildren, &typeName, &typeColor, &typeIcon,
	)
	if err != nil {
		return nil, err
	}

	if n.TypeID != nil && typeName.Valid {
		n.TypeCat = &entity.Category{
			ID:    *n.TypeID,
			Name:  typeName.String,
			Color: nullStringPtr(typeColor),
			Icon:  nullStringPtr(typeIcon),
		}
	}
	return n, nil
}

func (r *wbsRepo) GetNodeByPath(ctx context.Context, projectID int, path string) (*entity.WBSNode, error) {
	n := &entity.WBSNode{}
	var typeName, typeColor, typeIcon sql.NullString
	query := `SELECT wbs_nodes.id, wbs_nodes.project_id, wbs_nodes.title, wbs_nodes.type, wbs_nodes.type_id, wbs_nodes.path::text, wbs_nodes.order_index, 
	                 wbs_nodes.planned_start_date, wbs_nodes.planned_end_date, wbs_nodes.actual_start_date, wbs_nodes.actual_end_date,
					 wbs_nodes.progress, wbs_nodes.planned_value, wbs_nodes.actual_cost, wbs_nodes.estimated_effort, wbs_nodes.actual_effort, wbs_nodes.assigned_to, wbs_nodes.description, wbs_nodes.created_at, wbs_nodes.updated_at,
					 EXISTS(SELECT 1 FROM wbs_nodes sub WHERE sub.project_id = wbs_nodes.project_id AND sub.path <@ wbs_nodes.path AND sub.path != wbs_nodes.path) as has_children,
					 c.name as type_name, c.color as type_color, c.icon as type_icon
			  FROM wbs_nodes
			  LEFT JOIN categories c ON wbs_nodes.type_id = c.id
			  WHERE wbs_nodes.project_id = $1 AND wbs_nodes.path = $2::ltree`

	err := r.db.QueryRowContext(ctx, query, projectID, path).Scan(
		&n.ID, &n.ProjectID, &n.Title, &n.Type, &n.TypeID, &n.Path, &n.OrderIndex,
		&n.PlannedStartDate, &n.PlannedEndDate, &n.ActualStartDate, &n.ActualEndDate,
		&n.Progress, &n.PlannedValue, &n.ActualCost, &n.EstimatedEffort, &n.ActualEffort, &n.AssignedTo, &n.Description, &n.CreatedAt, &n.UpdatedAt,
		&n.HasChildren, &typeName, &typeColor, &typeIcon,
	)
	if err != nil {
		return nil, err
	}

	if n.TypeID != nil && typeName.Valid {
		n.TypeCat = &entity.Category{
			ID:    *n.TypeID,
			Name:  typeName.String,
			Color: nullStringPtr(typeColor),
			Icon:  nullStringPtr(typeIcon),
		}
	}
	return n, nil
}

func (r *wbsRepo) CreateNode(ctx context.Context, node *entity.WBSNode, parentPath string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var nextLevelPath string

	if parentPath == "" {
		var rootCount int
		err := tx.QueryRowContext(ctx, `SELECT COUNT(*) FROM wbs_nodes WHERE project_id = $1 AND nlevel(path) = 1`, node.ProjectID).Scan(&rootCount)
		if err != nil {
			return err
		}
		nextLevelPath = fmt.Sprintf("%d", rootCount+1)
		node.OrderIndex = rootCount + 1
	} else {
		var maxChildIndex sql.NullInt32

		normalizedParentPath := strings.ReplaceAll(parentPath, "_", "\\_") // simple escape

		err := tx.QueryRowContext(ctx, `
			SELECT MAX(NULLIF(regexp_replace(path::text, '^.*\.([0-9]+)$', '\1'), '')::int) 
			FROM wbs_nodes 
			WHERE path ~ ($1 || '.*{1}')::lquery AND project_id = $2
		`, normalizedParentPath, node.ProjectID).Scan(&maxChildIndex)

		if err != nil && err != sql.ErrNoRows {
			// If error is just syntax, fallback simple count
			var childCount int
			tx.QueryRowContext(ctx, `SELECT COUNT(*) FROM wbs_nodes WHERE path <@ $1::ltree AND nlevel(path) = nlevel($1::ltree) + 1 AND project_id=$2`, parentPath, node.ProjectID).Scan(&childCount)
			maxChildIndex.Int32 = int32(childCount)
			maxChildIndex.Valid = true
		}

		nextIdx := 1
		if maxChildIndex.Valid {
			nextIdx = int(maxChildIndex.Int32) + 1
		}
		nextLevelPath = fmt.Sprintf("%s.%d", parentPath, nextIdx)
		node.OrderIndex = nextIdx
	}

	query := `INSERT INTO wbs_nodes (project_id, title, type, type_id, path, order_index, planned_start_date, planned_end_date, assigned_to, description, planned_value, actual_cost, estimated_effort, actual_effort)
			  VALUES ($1, $2, $3, $4, $5::ltree, $6, $7, $8, $9, $10, $11, $12, $13, $14) RETURNING id, path::text, created_at, updated_at`

	err = tx.QueryRowContext(ctx, query,
		node.ProjectID, node.Title, node.Type, node.TypeID, nextLevelPath, node.OrderIndex,
		node.PlannedStartDate, node.PlannedEndDate, node.AssignedTo, node.Description, node.PlannedValue, node.ActualCost, node.EstimatedEffort, node.ActualEffort,
	).Scan(&node.ID, &node.Path, &node.CreatedAt, &node.UpdatedAt)

	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *wbsRepo) UpdateNode(ctx context.Context, n *entity.WBSNode) error {
	query := `UPDATE wbs_nodes 
			  SET title=$1, type=$2, type_id=$3, order_index=$4, planned_start_date=$5, planned_end_date=$6, 
			      actual_start_date=$7, actual_end_date=$8, progress=$9, assigned_to=$10, description=$11, 
			      planned_value=$12, actual_cost=$13, estimated_effort=$14, actual_effort=$15, updated_at=CURRENT_TIMESTAMP
			  WHERE id=$16`
	_, err := r.db.ExecContext(ctx, query,
		n.Title, n.Type, n.TypeID, n.OrderIndex, n.PlannedStartDate, n.PlannedEndDate,
		n.ActualStartDate, n.ActualEndDate, n.Progress, n.AssignedTo, n.Description,
		n.PlannedValue, n.ActualCost, n.EstimatedEffort, n.ActualEffort, n.ID,
	)
	return err
}

func (r *wbsRepo) DeleteNode(ctx context.Context, id int) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var path string
	err = tx.QueryRowContext(ctx, `SELECT path::text FROM wbs_nodes WHERE id = $1`, id).Scan(&path)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		return err
	}

	_, err = tx.ExecContext(ctx, `DELETE FROM wbs_nodes WHERE path <@ $1::ltree`, path)
	if err != nil {
		return err
	}

	return tx.Commit()
}
func (r *wbsRepo) ListDependencies(ctx context.Context, projectID int) ([]entity.WBSDependency, error) {
	query := `SELECT id, project_id, predecessor_id, successor_id, type, created_at FROM wbs_dependencies WHERE project_id = $1`
	rows, err := r.db.QueryContext(ctx, query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var deps []entity.WBSDependency
	for rows.Next() {
		var d entity.WBSDependency
		if err := rows.Scan(&d.ID, &d.ProjectID, &d.PredecessorID, &d.SuccessorID, &d.Type, &d.CreatedAt); err != nil {
			return nil, err
		}
		deps = append(deps, d)
	}
	if deps == nil {
		deps = []entity.WBSDependency{}
	}
	return deps, rows.Err()
}

func (r *wbsRepo) CreateDependency(ctx context.Context, dep *entity.WBSDependency) error {
	query := `INSERT INTO wbs_dependencies (project_id, predecessor_id, successor_id, type) VALUES ($1,$2,$3,$4) RETURNING id, created_at`
	return r.db.QueryRowContext(ctx, query, dep.ProjectID, dep.PredecessorID, dep.SuccessorID, dep.Type).Scan(&dep.ID, &dep.CreatedAt)
}

func (r *wbsRepo) DeleteDependency(ctx context.Context, depID int) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM wbs_dependencies WHERE id=$1`, depID)
	return err
}

func (r *wbsRepo) AddComment(ctx context.Context, c *entity.WBSComment) error {
	query := `INSERT INTO wbs_comments (project_id, node_id, user_id, content) VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at`
	return r.db.QueryRowContext(ctx, query, c.ProjectID, c.NodeID, c.UserID, c.Content).Scan(&c.ID, &c.CreatedAt, &c.UpdatedAt)
}

func (r *wbsRepo) ListComments(ctx context.Context, nodeID int) ([]entity.WBSComment, error) {
	query := `
		SELECT c.id, c.project_id, c.node_id, c.user_id, u.full_name as user_name, c.content, c.created_at, c.updated_at
		FROM wbs_comments c
		JOIN users u ON c.user_id = u.id
		WHERE c.node_id = $1
		ORDER BY c.created_at DESC`
	rows, err := r.db.QueryContext(ctx, query, nodeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []entity.WBSComment
	for rows.Next() {
		var c entity.WBSComment
		if err := rows.Scan(&c.ID, &c.ProjectID, &c.NodeID, &c.UserID, &c.UserName, &c.Content, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}
	if comments == nil {
		comments = []entity.WBSComment{}
	}
	return comments, rows.Err()
}

func (r *wbsRepo) DeleteComment(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM wbs_comments WHERE id=$1`, id)
	return err
}

func (r *wbsRepo) UpdateComment(ctx context.Context, id int, content string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE wbs_comments SET content=$1, updated_at=CURRENT_TIMESTAMP WHERE id=$2`, content, id)
	return err
}

func (r *wbsRepo) GetCommentByID(ctx context.Context, id int) (*entity.WBSComment, error) {
	query := `SELECT id, project_id, node_id, user_id, content, created_at, updated_at FROM wbs_comments WHERE id=$1`
	var c entity.WBSComment
	err := r.db.QueryRowContext(ctx, query, id).Scan(&c.ID, &c.ProjectID, &c.NodeID, &c.UserID, &c.Content, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *wbsRepo) CreateBaseline(ctx context.Context, baseline *entity.WBSBaseline) error {
	query := `INSERT INTO wbs_baselines (project_id, name, description, created_by) VALUES ($1, $2, $3, $4) RETURNING id, created_at`
	return r.db.QueryRowContext(ctx, query, baseline.ProjectID, baseline.Name, baseline.Description, baseline.CreatedBy).Scan(&baseline.ID, &baseline.CreatedAt)
}

func (r *wbsRepo) CopyNodesToBaseline(ctx context.Context, projectID int, baselineID int) error {
	query := `
		INSERT INTO wbs_baseline_nodes (baseline_id, node_id, path, planned_start_date, planned_end_date, progress, planned_value, actual_cost)
		SELECT $1, id, path, planned_start_date, planned_end_date, progress, planned_value, actual_cost
		FROM wbs_nodes
		WHERE project_id = $2
	`
	_, err := r.db.ExecContext(ctx, query, baselineID, projectID)
	return err
}

func (r *wbsRepo) GetBaselines(ctx context.Context, projectID int) ([]entity.WBSBaseline, error) {
	query := `SELECT id, project_id, name, description, created_by, created_at FROM wbs_baselines WHERE project_id = $1 ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var baselines []entity.WBSBaseline
	for rows.Next() {
		var b entity.WBSBaseline
		if err := rows.Scan(&b.ID, &b.ProjectID, &b.Name, &b.Description, &b.CreatedBy, &b.CreatedAt); err != nil {
			return nil, err
		}
		baselines = append(baselines, b)
	}
	if baselines == nil {
		baselines = []entity.WBSBaseline{}
	}
	return baselines, rows.Err()
}

func (r *wbsRepo) GetBaselineNodes(ctx context.Context, baselineID int) ([]entity.WBSBaselineNode, error) {
	query := `SELECT baseline_id, node_id, path::text, planned_start_date, planned_end_date, progress, planned_value, actual_cost FROM wbs_baseline_nodes WHERE baseline_id = $1 ORDER BY path ASC`
	rows, err := r.db.QueryContext(ctx, query, baselineID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var nodes []entity.WBSBaselineNode
	for rows.Next() {
		var n entity.WBSBaselineNode
		if err := rows.Scan(&n.BaselineID, &n.NodeID, &n.Path, &n.PlannedStartDate, &n.PlannedEndDate, &n.Progress, &n.PlannedValue, &n.ActualCost); err != nil {
			return nil, err
		}
		nodes = append(nodes, n)
	}
	if nodes == nil {
		nodes = []entity.WBSBaselineNode{}
	}
	return nodes, rows.Err()
}

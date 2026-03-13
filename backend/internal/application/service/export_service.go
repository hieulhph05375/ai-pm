package service

import (
	"context"
	"fmt"
	"io"
	"project-mgmt/backend/internal/domain/entity"
	"project-mgmt/backend/internal/domain/repository"

	"github.com/jung-kurt/gofpdf"
	"github.com/xuri/excelize/v2"
)

func strPtrValue(p *string) string {
	if p == nil {
		return "N/A"
	}
	return *p
}

type ExportService interface {
	ExportWBSExcel(ctx context.Context, projectID int, w io.Writer) error
	ExportProjectSummaryPDF(ctx context.Context, projectID int, stats *PMIStats, userID int, isAdmin bool, w io.Writer) error
	ExportProjectListExcel(ctx context.Context, search, status string, userID int, isAdmin bool, w io.Writer) error
}

type exportService struct {
	projectRepo  repository.ProjectRepository
	wbsRepo      repository.WBSRepository
	snapshotRepo repository.SnapshotRepository
	riskRepo     repository.RiskRepository
}

func NewExportService(pRepo repository.ProjectRepository, wRepo repository.WBSRepository, snapRepo repository.SnapshotRepository, rRepo repository.RiskRepository) ExportService {
	return &exportService{
		projectRepo:  pRepo,
		wbsRepo:      wRepo,
		snapshotRepo: snapRepo,
		riskRepo:     rRepo,
	}
}

func (s *exportService) ExportWBSExcel(ctx context.Context, projectID int, w io.Writer) error {
	filter := entity.WBSFilter{} // export all
	nodes, _, err := s.wbsRepo.GetProjectTree(ctx, projectID, filter)
	if err != nil {
		return err
	}

	f := excelize.NewFile()
	defer f.Close()

	sheet := "WBS"
	f.SetSheetName("Sheet1", sheet)

	// Headers
	headers := []string{"Mã Phân Cấp", "Tiêu đề", "Loại", "Ngày Bắt Đầu (Kế Hoạch)", "Ngày Kết Thúc (Kế Hoạch)", "Tiến độ (%)", "Ngân Sách (PV)", "Chi Phí Thực Tế (AC)"}
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, header)
	}

	// Data
	for i, node := range nodes {
		row := i + 2
		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), node.Path)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), node.Title)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", row), string(node.Type))

		if node.PlannedStartDate != nil {
			f.SetCellValue(sheet, fmt.Sprintf("D%d", row), node.PlannedStartDate.Format("2006-01-02"))
		}
		if node.PlannedEndDate != nil {
			f.SetCellValue(sheet, fmt.Sprintf("E%d", row), node.PlannedEndDate.Format("2006-01-02"))
		}

		f.SetCellValue(sheet, fmt.Sprintf("F%d", row), node.Progress)
		f.SetCellValue(sheet, fmt.Sprintf("G%d", row), node.PlannedValue)
		f.SetCellValue(sheet, fmt.Sprintf("H%d", row), node.ActualCost)
	}

	// Add Historical Data sheet
	histSheet := "Historical Data"
	f.NewSheet(histSheet)

	histHeaders := []string{"Captured At", "SPI", "CPI", "Progress (%)", "EV", "AC", "PV"}
	for i, header := range histHeaders {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(histSheet, cell, header)
	}

	snapshots, err := s.snapshotRepo.GetProjectSnapshots(ctx, projectID)
	if err == nil {
		for i, snap := range snapshots {
			row := i + 2
			f.SetCellValue(histSheet, fmt.Sprintf("A%d", row), snap.CapturedAt.Format("2006-01-02 15:04:05"))
			f.SetCellValue(histSheet, fmt.Sprintf("B%d", row), snap.SPI)
			f.SetCellValue(histSheet, fmt.Sprintf("C%d", row), snap.CPI)
			f.SetCellValue(histSheet, fmt.Sprintf("D%d", row), snap.Progress)
			f.SetCellValue(histSheet, fmt.Sprintf("E%d", row), snap.EV)
			f.SetCellValue(histSheet, fmt.Sprintf("F%d", row), snap.AC)
			f.SetCellValue(histSheet, fmt.Sprintf("G%d", row), snap.PV)
		}
	}

	return f.Write(w)
}

func (s *exportService) ExportProjectSummaryPDF(ctx context.Context, projectID int, stats *PMIStats, userID int, isAdmin bool, w io.Writer) error {
	project, err := s.projectRepo.GetByID(ctx, projectID, userID, isAdmin)
	if err != nil {
		return err
	}

	snapshots, err := s.snapshotRepo.GetProjectSnapshots(ctx, projectID)
	if err != nil {
		snapshots = []entity.ProjectSnapshot{}
	}

	risks, err := s.riskRepo.ListByProject(ctx, projectID, 3, 0)
	if err != nil {
		risks = []entity.Risk{}
	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Title
	pdf.SetFont("Arial", "B", 18)
	pdf.SetTextColor(30, 41, 59)
	pdf.Cell(40, 10, "EXECUTIVE SUMMARY (BLUF)")
	pdf.Ln(12)

	// BLUF
	pdf.SetFont("Arial", "B", 11)
	pdf.SetTextColor(71, 85, 105)
	pdf.Cell(0, 10, fmt.Sprintf("Project: %s | ID: %s", project.ProjectName, project.ProjectID))
	pdf.Ln(6)
	pdf.Cell(0, 10, fmt.Sprintf("Status: %s | Health: %s", strPtrValue(project.ProjectStatus), strPtrValue(project.OverallHealth)))
	pdf.Ln(8)

	// Add Trend Analysis
	trendMsg := "Stable"
	if len(snapshots) >= 2 {
		last := snapshots[len(snapshots)-1]
		prev := snapshots[len(snapshots)-2]
		if last.SPI < prev.SPI {
			trendMsg = fmt.Sprintf("Schedule is degrading (SPI: %.2f -> %.2f)", prev.SPI, last.SPI)
		} else if last.SPI > prev.SPI {
			trendMsg = fmt.Sprintf("Schedule is improving (SPI: %.2f -> %.2f)", prev.SPI, last.SPI)
		}
	}

	pdf.SetFont("Arial", "B", 10)
	pdf.SetTextColor(0, 0, 0)
	pdf.Cell(0, 10, fmt.Sprintf("Current Progress: %d%%", project.Progress))
	pdf.Ln(6)
	pdf.Cell(0, 10, fmt.Sprintf("Performance Trend: %s", trendMsg))
	pdf.Ln(10)

	// PMI Stats
	pdf.SetFont("Arial", "B", 12)
	pdf.SetTextColor(15, 23, 42)
	pdf.Cell(40, 10, "Financial & Schedule Performance (PMI)")
	pdf.Ln(8)
	pdf.SetFont("Arial", "", 10)
	pdf.SetTextColor(51, 65, 85)
	pdf.Cell(0, 10, fmt.Sprintf("Planned Value (PV): %.2f | Earned Value (EV): %.2f | Actual Cost (AC): %.2f", stats.PV, stats.EV, stats.AC))
	pdf.Ln(6)
	pdf.Cell(0, 10, fmt.Sprintf("Schedule Performance Index (SPI): %.2f", stats.SPI))
	pdf.Ln(6)
	pdf.Cell(0, 10, fmt.Sprintf("Cost Performance Index (CPI): %.2f", stats.CPI))
	pdf.Ln(10)

	// Top Risks
	if len(risks) > 0 {
		pdf.SetFont("Arial", "B", 12)
		pdf.SetTextColor(225, 29, 72) // rose-600
		pdf.Cell(40, 10, "Top Critical Risks & Issues")
		pdf.Ln(8)
		pdf.SetFont("Arial", "", 10)
		pdf.SetTextColor(51, 65, 85)
		for _, r := range risks {
			pdf.Cell(0, 8, fmt.Sprintf("- %s (Score: %d, Status: %s)", r.Title, r.RiskScore, r.Status))
			pdf.Ln(6)
		}
	}

	return pdf.Output(w)
}

func (s *exportService) ExportProjectListExcel(ctx context.Context, search, status string, userID int, isAdmin bool, w io.Writer) error {
	projects, _, err := s.projectRepo.List(ctx, 0, 10000, search, status, userID, isAdmin)
	if err != nil {
		return err
	}

	f := excelize.NewFile()
	defer f.Close()

	sheet := "Projects"
	f.SetSheetName("Sheet1", sheet)

	// Headers
	headers := []string{"ID Dự Án", "Tên Dự Án", "Quản Lý", "Giai Đoạn", "Trạng Thái", "Sức Khỏe", "Tiến Độ (%)", "Ngân Sách", "Ngày Bắt Đầu", "Ngày Kết Thúc"}
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, header)
	}

	// Data
	for i, p := range projects {
		row := i + 2
		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), p.ProjectID)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), p.ProjectName)

		manager := "N/A"
		if p.ProjectManager != nil {
			manager = *p.ProjectManager
		}
		f.SetCellValue(sheet, fmt.Sprintf("C%d", row), manager)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", row), p.CurrentPhase)
		f.SetCellValue(sheet, fmt.Sprintf("E%d", row), p.ProjectStatus)
		f.SetCellValue(sheet, fmt.Sprintf("F%d", row), p.OverallHealth)
		f.SetCellValue(sheet, fmt.Sprintf("G%d", row), p.Progress)
		f.SetCellValue(sheet, fmt.Sprintf("H%d", row), p.ApprovedBudget)

		if p.PlannedStartDate != nil {
			f.SetCellValue(sheet, fmt.Sprintf("I%d", row), p.PlannedStartDate.Format("2006-01-02"))
		}
		if p.PlannedEndDate != nil {
			f.SetCellValue(sheet, fmt.Sprintf("J%d", row), p.PlannedEndDate.Format("2006-01-02"))
		}
	}

	return f.Write(w)
}

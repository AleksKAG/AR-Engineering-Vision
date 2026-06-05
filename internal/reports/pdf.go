package reports

import (
	"context"
	"fmt"
	"github.com/AleksKAG/ar-backend/internal/s3"
	"github.com/AleksKAG/ar-backend/internal/store"
	"github.com/johnfercher/maroto/v2/pkg/consts"
	"github.com/johnfercher/maroto/v2/pkg/props"
	"github.com/johnfercher/maroto/v2"
)

// GenerateProjectReport создает PDF со списком всех замечаний по проекту
func GenerateProjectReport(ctx context.Context, db *store.DB, s3c *s3.Client, projectID, projectName string) (string, error) {
	issues, err := db.GetIssuesByProject(ctx, projectID)
	if err != nil {
		return "", err
	}

	m := maroto.NewMaroto(consts.Portrait, consts.A4)

	// Заголовок
	m.Row(20, func() {
		m.Col(12, func() {
			m.Text(fmt.Sprintf("Отчёт по скрытым работам: %s", projectName), props.Text{
				Style: consts.Bold, Size: 14, Align: consts.Center,
			})
		})
	})

	// Таблица замечаний
	m.Row(10, func() {
		m.Col(3, func() { m.Text("Элемент", props.Text{Style: consts.Bold}) })
		m.Col(3, func() { m.Text("AI Детекция", props.Text{Style: consts.Bold}) })
		m.Col(2, func() { m.Text("Смещение (мм)", props.Text{Style: consts.Bold}) })
		m.Col(4, func() { m.Text("Комментарий", props.Text{Style: consts.Bold}) })
	})

	for _, issue := range issues {
		statusMark := "✅"
		if !issue.IsMatch || issue.DeviationMM > 50.0 {
			statusMark = "⚠️"
		}

		m.Row(10, func() {
			m.Col(3, func() { m.Text(fmt.Sprintf("%s ID: %s", statusMark, issue.ElementID[:8]), props.Text{Size: 10}) })
			m.Col(3, func() { 
				txt := fmt.Sprintf("%s (%.0f%%)", issue.AIDetectedType, issue.AIConfidence*100)
				if issue.AIDetectedType == "" { txt = "Не проверено" }
				m.Text(txt, props.Text{Size: 10}) 
			})
			m.Col(2, func() { m.Text(fmt.Sprintf("%.0f", issue.DeviationMM), props.Text{Size: 10}) })
			m.Col(4, func() { m.Text(issue.Comment, props.Text{Size: 10}) })
		})
	}

	// Генерируем PDF в байты
	pdfBytes, err := m.Generate()
	if err != nil {
		return "", err
	}

	// Загружаем PDF в MinIO
	objectName := fmt.Sprintf("reports/%s/report_%s.pdf", projectID, projectName)
	_, err = s3c.Upload(ctx, objectName, pdfBytes, int64(len(pdfBytes)), "application/pdf")
	if err != nil {
		return "", err
	}

	return objectName, nil
}

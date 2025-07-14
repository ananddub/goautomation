package main

import (
	"context"
	"fmt"
	"goautomation/sound"
	"image/color"
	"sort"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"goautomation/db"
	"goautomation/db/sqlc"
)

// Simple theme with medium fonts and basic colors
type simpleTheme struct{}

func (t *simpleTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	return theme.DefaultTheme().Color(name, variant)
}

func (t *simpleTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (t *simpleTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (t *simpleTheme) Size(name fyne.ThemeSizeName) float32 {
	switch name {
	case theme.SizeNameText:
		return 14 // Medium font size
	case theme.SizeNameHeadingText:
		return 16
	case theme.SizeNameSubHeadingText:
		return 15
	}
	return theme.DefaultTheme().Size(name)
}

type MonitorGUI struct {
	app          fyne.App
	window       fyne.Window
	database     *db.Database
	table        *widget.Table
	statusLabel  *widget.Label
	alertCount   *widget.Label
	totalCount   *widget.Label
	isMonitoring bool
	data         []sqlc.CheckRecentChanges5minRow
	stopChan     chan bool
}

func NewMonitorGUI() *MonitorGUI {
	myApp := app.New()
	myApp.Settings().SetTheme(&simpleTheme{})

	window := myApp.NewWindow("Window Monitor")
	window.Resize(fyne.NewSize(1200, 700))
	window.CenterOnScreen()

	return &MonitorGUI{
		app:      myApp,
		window:   window,
		stopChan: make(chan bool),
	}
}

func (gui *MonitorGUI) initDatabase() error {
	database, err := db.NewDatabase("./data/db.sqlite")
	if err != nil {
		return fmt.Errorf("failed to initialize database: %v", err)
	}
	gui.database = database
	return nil
}

func (gui *MonitorGUI) checkCondition(row sqlc.CheckRecentChanges5minRow) bool {
	// Changed condition to >= 7 instead of > 7 to match actual data
	return row.TotalPixel >= 7 && row.TotalUniquePixel <= 2
}

func (gui *MonitorGUI) loadData() {
	ctx := context.Background()
	//monitor.UpdateTable("raiboss1", gui.database.Queries, ctx)
	//monitor.UpdateTable("raiboss2", gui.database.Queries, ctx)
	//monitor.UpdateTable("mylaptop", gui.database.Queries, ctx)
	changes, err := gui.database.Queries.CheckRecentChanges5min(ctx)

	if err != nil {
		gui.statusLabel.SetText(fmt.Sprintf("Error: %v", err))
		return
	}

	// Sort data so alert items appear at the top
	flag := false
	sort.Slice(changes, func(i, j int) bool {
		alertI := gui.checkCondition(changes[i])
		alertJ := gui.checkCondition(changes[j])

		// If one is alert and other is not, alert comes first
		if alertI && !alertJ {
			flag = true
			return true
		}
		if !alertI && alertJ {
			return false
		}

		// If both are same type (both alert or both normal), sort by name
		return changes[i].Name < changes[j].Name
	})

	alertCount := 0
	for _, row := range changes {
		if gui.checkCondition(row) {
			alertCount++
		}
	}

	gui.data = changes
	gui.totalCount.SetText(fmt.Sprintf("Total: %d", len(changes)))
	gui.alertCount.SetText(fmt.Sprintf("Alerts: %d", alertCount))
	gui.statusLabel.SetText(fmt.Sprintf("Updated: %s", time.Now().Format("15:04:05")))

	if gui.table != nil {
		gui.table.Refresh()
	}
	if flag {
		sound.PlayBeep()
	}
}

func (gui *MonitorGUI) createTable() {
	headers := []string{"#", "Application", "Window Title", "Total Pixels", "Unique Pixels", "Status"}

	gui.table = widget.NewTable(
		func() (int, int) {
			return len(gui.data) + 1, len(headers)
		},
		func() fyne.CanvasObject {
			label := widget.NewLabel("Template")
			return label
		},
		func(id widget.TableCellID, obj fyne.CanvasObject) {
			label := obj.(*widget.Label)

			// Header row
			if id.Row == 0 {
				label.SetText(headers[id.Col])
				label.TextStyle = fyne.TextStyle{Bold: true}
				return
			}

			rowIndex := id.Row - 1
			if rowIndex >= len(gui.data) {
				label.SetText("")
				return
			}

			row := gui.data[rowIndex]
			isAlert := gui.checkCondition(row)

			switch id.Col {
			case 0:
				if isAlert {
					label.SetText(fmt.Sprintf("%d", rowIndex+1))
					label.TextStyle = fyne.TextStyle{Bold: true}
					label.Importance = widget.WarningImportance
				} else {
					label.SetText(fmt.Sprintf("%d", rowIndex+1))
					label.TextStyle = fyne.TextStyle{}
					label.Importance = widget.MediumImportance
				}
			case 1:
				if isAlert {
					label.SetText(fmt.Sprintf("%s", row.Name))
					label.TextStyle = fyne.TextStyle{Bold: true}
					label.Importance = widget.WarningImportance
				} else {
					label.SetText(row.Name)
					label.TextStyle = fyne.TextStyle{}
					label.Importance = widget.MediumImportance
				}
			case 2:
				if isAlert {
					label.SetText(fmt.Sprintf("%s", row.Title))
					label.TextStyle = fyne.TextStyle{Bold: true}
					label.Importance = widget.WarningImportance
				} else {
					label.SetText(row.Title)
					label.TextStyle = fyne.TextStyle{}
					label.Importance = widget.MediumImportance
				}
			case 3:
				if isAlert {
					label.SetText(fmt.Sprintf("%d", row.TotalPixel))
					label.TextStyle = fyne.TextStyle{Bold: true}
					label.Importance = widget.WarningImportance
				} else {
					label.SetText(fmt.Sprintf("%d", row.TotalPixel))
					label.TextStyle = fyne.TextStyle{}
					label.Importance = widget.MediumImportance
				}
			case 4:
				if isAlert {
					label.SetText(fmt.Sprintf("%d", row.TotalUniquePixel))
					label.TextStyle = fyne.TextStyle{Bold: true}
					label.Importance = widget.WarningImportance
				} else {
					label.SetText(fmt.Sprintf("%d", row.TotalUniquePixel))
					label.TextStyle = fyne.TextStyle{}
					label.Importance = widget.MediumImportance
				}
			case 5: // Status
				if isAlert {
					label.SetText("ALERT")
					label.TextStyle = fyne.TextStyle{Bold: true}
					label.Importance = widget.WarningImportance
				} else {
					label.SetText("OK")
					label.TextStyle = fyne.TextStyle{}
					label.Importance = widget.SuccessImportance
				}
			}
		},
	)

	gui.table.SetColumnWidth(0, 50)
	gui.table.SetColumnWidth(1, 150)
	gui.table.SetColumnWidth(2, 350)
	gui.table.SetColumnWidth(3, 100)
	gui.table.SetColumnWidth(4, 100)
	gui.table.SetColumnWidth(5, 100)
}

func (gui *MonitorGUI) startMonitoring() {
	if gui.isMonitoring {
		return
	}

	gui.isMonitoring = true
	gui.statusLabel.SetText("Monitoring started...")

	go func() {
		gui.loadData()
		ticker := time.NewTicker(20 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				gui.loadData()
			case <-gui.stopChan:
				return
			}
		}
	}()
}

func (gui *MonitorGUI) stopMonitoring() {
	if !gui.isMonitoring {
		return
	}

	gui.isMonitoring = false
	gui.statusLabel.SetText("Monitoring stopped")
	select {
	case gui.stopChan <- true:
	default:
	}
}

func (gui *MonitorGUI) createUI() {
	// Simple header
	headerText := widget.NewLabel("Window Monitor Dashboard")
	headerText.TextStyle = fyne.TextStyle{Bold: true}
	headerText.Alignment = fyne.TextAlignCenter

	// Status labels
	gui.statusLabel = widget.NewLabel("Ready to monitor")
	gui.alertCount = widget.NewLabel("Alerts: 0")
	gui.totalCount = widget.NewLabel("Total: 0")

	// Simple buttons
	startBtn := widget.NewButton("Start Monitoring", func() {
		gui.startMonitoring()
	})

	stopBtn := widget.NewButton("Stop", func() {
		gui.stopMonitoring()
	})

	refreshBtn := widget.NewButton("Refresh", func() {
		gui.loadData()
	})

	// Condition info
	conditionText := widget.NewLabel("Alert Condition: TotalPixel >= 7 && TotalUniquePixel <= 3")
	conditionCard := widget.NewCard("Alert Configuration", "", conditionText)

	// Create table
	gui.createTable()

	// Stats container
	statsContainer := container.NewHBox(
		gui.totalCount,
		widget.NewSeparator(),
		gui.alertCount,
		widget.NewSeparator(),
		gui.statusLabel,
	)

	// Controls
	controlsContainer := container.NewHBox(
		startBtn,
		stopBtn,
		refreshBtn,
	)

	// Info section
	infoSection := container.NewVBox(
		conditionCard,
		widget.NewSeparator(),
		container.NewHBox(
			widget.NewLabel("Controls:"),
			controlsContainer,
		),
		widget.NewSeparator(),
		container.NewHBox(
			widget.NewLabel("Status:"),
			statsContainer,
		),
	)

	// Table with dynamic height - using border layout for better space utilization
	tableLabel := widget.NewLabel("Monitoring Results:")
	tableLabel.TextStyle = fyne.TextStyle{Bold: true}

	// Create scroll container without fixed size - let it expand dynamically
	tableScroll := container.NewScroll(gui.table)

	// Use border layout to give table maximum available space
	tableContainer := container.NewBorder(
		tableLabel,  // top
		nil,         // bottom
		nil,         // left
		nil,         // right
		tableScroll, // center - takes remaining space
	)

	// Main content using border layout for dynamic sizing
	content := container.NewBorder(
		container.NewVBox(
			headerText,
			widget.NewSeparator(),
			infoSection,
			widget.NewSeparator(),
		), // top section
		nil,            // bottom
		nil,            // left
		nil,            // right
		tableContainer, // center - table gets remaining space
	)

	gui.window.SetContent(container.NewPadded(content))
	gui.loadData()
}

func (gui *MonitorGUI) Run() {
	if err := gui.initDatabase(); err != nil {
		panic(fmt.Sprintf("Failed to initialize database: %v", err))
	}
	defer gui.database.Close()

	gui.createUI()

	gui.window.SetCloseIntercept(func() {
		gui.stopMonitoring()
		gui.window.Close()
	})

	gui.window.ShowAndRun()
}

func main() {
	gui := NewMonitorGUI()
	gui.Run()
}

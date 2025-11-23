package ui

import (
	"image/color"
	"log"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	fynetooltip "github.com/dweymouth/fyne-tooltip"
	ttwidget "github.com/dweymouth/fyne-tooltip/widget"
	// "github.com/gitscope/internal/ui"
)

type CustomTheme struct {
	fyne.Theme
}

func (c *CustomTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	if name == theme.ColorNamePrimary {
		return color.NRGBA{R: 148, G: 0, B: 211, A: 255} // Violet 💜
	}
	return c.Theme.Color(name, variant)
}

func App() {
	myApp := app.New()
	myApp.Settings().SetTheme(&CustomTheme{theme.DefaultTheme()})
	myWindow := myApp.NewWindow("Gitscope")

	logo, err := os.ReadFile("assets/icons/gitscope_logo_v6.png")
	if err != nil {
		log.Println("Error reading icon file ", err)
	} else {
		appIcon := fyne.NewStaticResource("gitscope_logo_v6.png", logo)
		myWindow.SetIcon(appIcon)
	}

	sidebar, content := SideBar(myWindow)
	mainLayout := container.NewBorder(nil, nil, sidebar, nil, content)
	myWindow.SetContent(fynetooltip.AddWindowToolTipLayer(mainLayout, myWindow.Canvas()))
	myWindow.Resize(fyne.NewSize(900, 600))
	myWindow.ShowAndRun()
}

func SideBar(w fyne.Window) (fyne.CanvasObject, fyne.CanvasObject) {

	var sidebar fyne.CanvasObject
	var Addbtn, Repobtn, settingsbtn, documentbtn *ttwidget.Button

	// Repo button
	Repobtn = ttwidget.NewButtonWithIcon("", theme.FolderOpenIcon(), nil)
	Repobtn.SetToolTip("Repo")
	Repobtn.OnTapped = func() {
		SetActive(Repobtn, []*ttwidget.Button{Addbtn, settingsbtn, documentbtn})
		reposPage := RepositoryPage(w)
		w.SetContent(fynetooltip.AddWindowToolTipLayer(
			container.NewBorder(nil, nil, sidebar, nil, reposPage),
			w.Canvas(),
		))
	}

	// Dashboard button
	dashboardPage := dashBoardPage(w)
	Addbtn = ttwidget.NewButtonWithIcon("", theme.ContentAddIcon(), func() {
		SetActive(Addbtn, []*ttwidget.Button{Repobtn, settingsbtn, documentbtn})
		w.SetContent(fynetooltip.AddWindowToolTipLayer(
			container.NewBorder(nil, nil, sidebar, nil, dashboardPage),
			w.Canvas(),
		))
	})
	Addbtn.SetToolTip("Git common operations")

	// Settings button
	settingsPage := SettingPage(w)
	settingsbtn = ttwidget.NewButtonWithIcon("", theme.InfoIcon(), func() {
		SetActive(settingsbtn, []*ttwidget.Button{Addbtn, Repobtn, documentbtn})
		w.SetContent(fynetooltip.AddWindowToolTipLayer(
			container.NewBorder(nil, nil, sidebar, nil, container.NewCenter(settingsPage)),
			w.Canvas(),
		))
	})
	settingsbtn.SetToolTip("About")

	// Documentation button
	documentPage := DocumentPage(w)
	documentbtn = ttwidget.NewButtonWithIcon("", theme.DocumentIcon(), func() {
		SetActive(documentbtn, []*ttwidget.Button{Addbtn, Repobtn, settingsbtn})

		w.SetContent(
			fynetooltip.AddWindowToolTipLayer(
				container.NewBorder(nil, nil, sidebar, nil, documentPage),
				w.Canvas(),
			),
		)
	})
	documentbtn.SetToolTip("Documentation")

	// Sidebar layout: Repo first
	sidebar = container.NewVBox(
		Repobtn,
		Addbtn,
		settingsbtn,
		documentbtn,
		layout.NewSpacer(),
	)

	// Set default active tab to Repo
	SetActive(Repobtn, []*ttwidget.Button{Addbtn, settingsbtn, documentbtn})
	reposPage := RepositoryPage(w)
	w.SetContent(fynetooltip.AddWindowToolTipLayer(
		container.NewBorder(nil, nil, sidebar, nil, reposPage),
		w.Canvas(),
	))

	return sidebar, reposPage
}

func SetActive(active *ttwidget.Button, others []*ttwidget.Button) {
	active.Importance = widget.HighImportance
	active.Refresh()
	for _, b := range others {
		b.Importance = widget.MediumImportance
		b.Refresh()
	}
}

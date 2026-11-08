package gui

import (
	appl "arshin/internal/app"
	"arshin/internal/domain"
	"bufio"
	"errors"
	"fmt"
	"io"
	"regexp"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

const (
	toXMLString = "Конвертировать в XML"
	toAPIString = "Отправить по API"
)

type Loader interface {
	Load(string) ([]domain.VerificationRecord, error)
}

type Deliverer interface {
	Deliver(string, []domain.Verification) error
}

type Window struct {
	window        fyne.Window
	loader        Loader
	xmlConverter  Deliverer
	fgisClient    Deliverer
	fgisUrl       string
	logScanner    *bufio.Scanner
	app           *appl.Service
	inputFile     string
	radio         *widget.RadioGroup
	filenameLabel *widget.Label
	logs          *widget.TextGrid
}

func NewApp(loader Loader, xmlConverter, fgisClient Deliverer, fgisUrl string, logReader io.Reader, application *appl.Service) *Window {
	radio := widget.NewRadioGroup([]string{toXMLString, toAPIString}, func(s string) {})
	radio.Horizontal = true
	radio.Required = true
	radio.Selected = toXMLString

	label := widget.NewLabel("Выберите файл")

	logs := widget.NewTextGrid()

	scroll := container.NewVScroll(logs)
	scroll.SetMinSize(fyne.NewSize(500, 400))

	res := &Window{
		window:        app.New().NewWindow("Конвертер"),
		radio:         radio,
		filenameLabel: label,
		loader:        loader,
		xmlConverter:  xmlConverter,
		fgisClient:    fgisClient,
		logs:          logs,
		fgisUrl:       fgisUrl,
		app:           application,
	}

	logScanner := bufio.NewScanner(logReader)

	go func() {
		for logScanner.Scan() {
			line := logScanner.Text()
			fyne.Do(func() {
				logs.Append(line)
			})
		}
	}()

	t := container.NewAppTabs(
		container.NewTabItem("Импорт данных",
			container.NewVBox(
				container.NewVBox(widget.NewLabel("Новая обработка"), widget.NewLabel("Выберите файл и укажите, что делать с данными")),
				container.NewVBox(widget.NewLabel("1. Исходный файл"), container.NewHBox(res.filenameLabel, widget.NewButton("Выбрать файл", res.openFile))),
				container.NewVBox(widget.NewLabel("2. Действие с данными"), res.radio),
				widget.NewButton("Запустить обработку", res.start),
				container.NewVBox(widget.NewLabel("Журнал"), scroll),
			)),
		container.NewTabItem("Конфигурация и БД", widget.NewLabel("В разработке")),
	)

	res.window.SetContent(t)
	return res
}

func (w *Window) Run() {
	w.window.ShowAndRun()
}

func (w *Window) openFile() {
	filedialog := dialog.NewFileOpen(w.fileSelected, w.window)
	filedialog.Show()
}

func (w *Window) fileSelected(reader fyne.URIReadCloser, err error) {
	if err != nil {
		dialog.ShowError(err, w.window)
		return
	}
	if reader != nil {
		w.inputFile = reader.URI().Path()
		w.filenameLabel.SetText(w.inputFile)
	}
}

func (w *Window) start() {
	var deliverer Deliverer
	var dst string

	w.logs.SetText("")

	if w.inputFile == "" {
		dialog.ShowError(fmt.Errorf("Файл не выбран"), w.window)
		return
	}

	switch w.radio.Selected {
	case toXMLString:
		deliverer = w.xmlConverter
		reg, err := regexp.Compile(`\.(xlsx|xls)$`)
		if err != nil {
			dialog.ShowError(err, w.window)
		}
		dst = reg.ReplaceAllString(w.inputFile, ".xml")
	case toAPIString:
		deliverer = w.fgisClient
		dst = w.fgisUrl
	default:
		dialog.ShowError(fmt.Errorf("Не выбрано действие с данными"), w.window)
		return
	}

	go func() {
		if err := w.app.Handle(w.inputFile, w.loader, dst, deliverer); err != nil {
			var dataErr domain.ErrInvalidData
			if errors.Is(err, domain.ErrNoRecords) || errors.As(err, &dataErr) {
				dialog.ShowError(err, w.window)
			} else {
				dialog.ShowError(errors.New("Внутренняя ошибка"), w.window)
			}
		}
	}()
}

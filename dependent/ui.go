package dependent

import (
	"fmt"
	"image/color"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type Config struct {
	App        fyne.App
	InfoLog    *log.Logger
	ErrorLog   *log.Logger
	MainWindow fyne.Window
}

type OriginAtbStyle struct {
	Button      widget.Button
	ImgPath     string
	ImgResource fyne.Resource
	LabelName   string
}

type ExtendAtbStyle struct {
	MainButton     OriginAtbStyle
	TimesLabel     widget.Label
	RestraintTimes float64
	TimesColor     color.Color
}

var (
	myApp                      Config
	wg                         sync.WaitGroup
	scrollUI                   *container.MyScroll
	Spacer                     = &canvas.Line{StrokeColor: GetColor("transparent"), StrokeWidth: 4}
	gradientLineBlackToSkyBlue = canvas.NewRadialGradient(GetColor("black"), GetColor("skyblue"))
	realLineSkyBlue            = &canvas.Line{StrokeColor: GetColor("skyblue"), StrokeWidth: 4}
	RestraintTimeButtonSize    = fyne.NewSize(90, 60)
	ButtonSize                 = fyne.NewSize(90, 40)
	IconSize                   = fyne.NewSize(30, 30)
	MainWindowSize             = fyne.NewSize(800, 600)
	ButtonsBase                = make(map[uint16]*widget.Button)
	ButtonsMulti               = make(map[uint16]*widget.Button)
	CtnRelateMap               = make(map[uint16][]*widget.Button)
	CtnAttackMap               = make(map[uint16][]*fyne.Container)
	CtnRecMap                  = make(map[uint16][]*fyne.Container)
)

const (
	AppIDString           = "goCode.SeerAttributeRestraintTable"
	MainWindowTitleString = "赛尔号属性克制表"

	AttackImgFullPath     = ".\\resources\\attack.png"
	BaseIconPath          = ".\\resources\\base.png"
	LinkImgFullPath       = ".\\resources\\link.png"
	MainWindowImgFullPath = ".\\resources\\screenDuck.png"
	MultiIconPath         = ".\\resources\\multi.png"
	RetImgFullPath        = ".\\resources\\back.png"
	RecipientImgFullPath  = ".\\resources\\recipient.png"

	BaseImgPath  = ".\\resources\\img\\"
	MultiImgPath = ".\\resources\\img\\multi\\"
)

func GetColor(kind string) color.NRGBA {
	switch strings.ToLower(kind) {
	case "red":
		return color.NRGBA{R: 180, G: 0, B: 0, A: 255}
	case "blue":
		return color.NRGBA{R: 40, G: 100, B: 150, A: 255}
	case "skyblue":
		return color.NRGBA{R: 210, G: 240, B: 245, A: 255}
	case "seablue":
		return color.NRGBA{R: 0, G: 150, B: 235, A: 255}
	case "grey":
		return color.NRGBA{R: 155, G: 155, B: 155, A: 255}
	case "darkgrey":
		return color.NRGBA{R: 65, G: 65, B: 65, A: 255}
	case "white":
		return color.NRGBA{R: 255, G: 255, B: 255, A: 255}
	case "black":
		return color.NRGBA{R: 0, G: 0, B: 0, A: 255}
	case "transparent":
		return color.NRGBA{R: 0, G: 0, B: 0, A: 0}
	default:
		return color.NRGBA{R: 0, G: 0, B: 0, A: 0}
	}
}

func NewAndSetText(text string, textColor color.Color, textSize float32, bold bool) *canvas.Text {
	ret := canvas.NewText(text, textColor)
	ret.TextSize = textSize
	ret.TextStyle.Bold = bold

	return ret
}

func RunWindow() {
	//Start and set an app
	fyneApp := app.NewWithID(AppIDString)
	myApp.App = fyneApp
	myApp.App.Settings().SetTheme(&myTheme{})

	//Set icon, window's size and so on
	iconRsc, _ := fyne.LoadResourceFromPath(MainWindowImgFullPath)
	SetConfig(&myApp, fyneApp, MainWindowTitleString, iconRsc, true)
	InitBaseAtbRelates()

	//Get UI layout
	ctn := showMainUI()

	//set other variables
	scrollUI = container.NewVMyScroll(ctn)
	scrollUI.ScrollSpeed = container.MyScrollSpeedFast

	myApp.MainWindow.SetContent(scrollUI)
	myApp.MainWindow.Resize(MainWindowSize)
	myApp.MainWindow.CenterOnScreen()
	myApp.MainWindow.SetFixedSize(false)

	myApp.MainWindow.ShowAndRun()
}

func SetConfig(app *Config, fyneApp fyne.App, title string, iconRsc fyne.Resource, isMaster bool) {
	app.InfoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.ErrorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	app.MainWindow = fyneApp.NewWindow(title)
	app.MainWindow.SetIcon(iconRsc)
	app.MainWindow.SetFixedSize(isMaster)
	app.MainWindow.SetMaster()
}

func atbButtonMainClicked(ctnBox *fyne.Container, atbTab *container.AppTabs, id uint16) func() {
	return func() {
		if scrollUI != nil {
			scrollUI.Offset = fyne.NewPos(0, 0)
		}

		var (
			tp          = GetAtbIDType(id)
			curAtbImg   canvas.Image
			curAtbText1 = NewAndSetText(AtbNameMap[id]+"系", GetColor("red"), 20, true)
			curAtbText2 = NewAndSetText("属性相克表", GetColor("seablue"), 18, false)
			curAtbCtn   = container.NewHBox(&curAtbImg, curAtbText1, curAtbText2)

			linkRsc, _    = fyne.LoadResourceFromPath(LinkImgFullPath)
			attackRsc, _  = fyne.LoadResourceFromPath(AttackImgFullPath)
			recRsc, _     = fyne.LoadResourceFromPath(RecipientImgFullPath)
			ctnRelate     = container.NewGridWrap(ButtonSize)
			ctnAttack     = container.NewGridWrap(RestraintTimeButtonSize)
			ctnRec        = container.NewGridWrap(RestraintTimeButtonSize)
			retIconRsc, _ = fyne.LoadResourceFromPath(RetImgFullPath)
			relateBox     = container.NewHBox(widget.NewIcon(linkRsc), widget.NewLabel("相关属性"))
			retBtt        = widget.NewButtonWithIcon("返回", retIconRsc, returnButtonClicked(atbTab, ctnBox))
			retCtn        = container.NewHBox(retBtt, Spacer)
			tourBar       = container.NewMax(retCtn, container.NewCenter(curAtbCtn))
			atbARTab      = container.NewAppTabs(
				container.NewTabItemWithIcon("攻击效果", attackRsc, ctnAttack),
				container.NewTabItemWithIcon("被攻击效果", recRsc, ctnRec),
			)
		)

		if tp == 1 {
			curAtbImg = *canvas.NewImageFromFile(BaseImgPath + strconv.Itoa(int(id)) + ".png")
		} else {
			curAtbImg = *canvas.NewImageFromFile(MultiImgPath + strconv.Itoa(int(id)) + ".png")
		}
		curAtbImg.FillMode = canvas.ImageFillOriginal
		curAtbImg.Resize(IconSize)

		if len(CtnRelateMap[id]) == 0 {
			ctnRelate.Add(widget.NewLabelWithStyle("无相关属性", 0, fyne.TextStyle{Bold: true}))
		} else {
			for _, val := range CtnRelateMap[id] {
				ctnRelate.Add(val)
			}
		}

		if len(CtnAttackMap[id]) == 0 {
			ctnAttack.Add(widget.NewLabelWithStyle("无相关属性", 0, fyne.TextStyle{Bold: true}))
		} else {
			for _, val := range CtnAttackMap[id] {
				ctnAttack.Add(val)
			}
		}

		if len(CtnRecMap[id]) == 0 {
			ctnRec.Add(widget.NewLabelWithStyle("无相关属性", 0, fyne.TextStyle{Bold: true}))
		} else {
			for _, val := range CtnRecMap[id] {
				ctnRec.Add(val)
			}
		}

		ctnBox.RemoveAll()
		ctnBox.Add(tourBar)
		ctnBox.Add(gradientLineBlackToSkyBlue)
		ctnBox.Add(relateBox)
		ctnBox.Add(ctnRelate)
		ctnBox.Add(realLineSkyBlue)
		ctnBox.Add(atbARTab)

		if !atbTab.Hidden {
			atbTab.Hide()
			ctnBox.Show()
		}
	}
}

func getButtonTemp(temp *widget.Button) *widget.Button {
	ret := &widget.Button{
		Text:          temp.Text,
		Icon:          temp.Icon,
		Importance:    temp.Importance,
		Alignment:     temp.Alignment,
		IconPlacement: temp.IconPlacement,
		OnTapped:      temp.OnTapped,
	}

	return ret
}

func getRestraintTimesButton(id uint16, times float64) *fyne.Container {
	var (
		tp      = GetAtbIDType(id)
		bgColor color.Color
		vbox    = container.NewVBox()
		maxCtn  = container.NewMax()
	)

	if times > 1 {
		bgColor = GetColor("red")
	} else if times == 0 {
		bgColor = GetColor("grey")
	} else {
		bgColor = GetColor("blue")
	}
	maxCtn.Add(canvas.NewRectangle(bgColor))
	maxCtn.Add(&canvas.Text{
		Alignment: fyne.TextAlignCenter,
		Color:     color.White,
		Text:      fmt.Sprintf("%g", times),
		TextSize:  16,
	})

	if tp == 1 {
		vbox.Add(getButtonTemp(ButtonsBase[id]))
		vbox.Add(maxCtn)
		return vbox
	} else if tp == 2 {
		vbox.Add(getButtonTemp(ButtonsMulti[id]))
		vbox.Add(maxCtn)
		return vbox
	}

	return nil
}

func initAtbDetails(id uint16) {
	defer wg.Done()

	posTable, negTable := GetRestraintTimes(id)

	//set relates
	tp := GetAtbIDType(id)
	if tp == 1 {
		for _, i := range BaseAtbRelate[id] {
			if i != 0 {
				lock.Lock()
				if i < 1000 {
					CtnRelateMap[id] = append(CtnRelateMap[id], getButtonTemp(ButtonsBase[i]))
				} else {
					CtnRelateMap[id] = append(CtnRelateMap[id], getButtonTemp(ButtonsMulti[i]))
				}
				lock.Unlock()
			}
		}
	} else {
		lock.Lock()
		CtnRelateMap[id] = []*widget.Button{getButtonTemp(ButtonsBase[id/100]), getButtonTemp(ButtonsBase[id%100])}
		lock.Unlock()
	}

	//set positive and negative button table
	len1, len2 := len(posTable), len(negTable)
	length := 0
	if len1 > len2 {
		length = len1
	} else {
		length = len2
	}

	for i := length; i >= 0; i-- {
		if i < len1 {
			lock.Lock()
			CtnAttackMap[id] = append(CtnAttackMap[id], getRestraintTimesButton(posTable[i].Recipient, posTable[i].Times))
			lock.Unlock()
		}

		if i < len2 {
			lock.Lock()
			CtnRecMap[id] = append(CtnRecMap[id], getRestraintTimesButton(negTable[i].Attacker, negTable[i].Times))
			lock.Unlock()
		}
	}

	return
}

func initContainer(CtnBase, CtnMulti, vbox *fyne.Container, atbTab *container.AppTabs) {
	var index uint16
	for _, index = range AtbIDList {
		wg.Add(1)
		go processButtonsForInitContainer(index, vbox, atbTab)
	}
	wg.Wait()
	for _, index = range AtbIDList {
		wg.Add(1)
		go initAtbDetails(index)
	}
	wg.Wait()

	for _, index := range AtbIDList {
		if index < 1000 {
			CtnBase.Add(ButtonsBase[index])
		} else {
			CtnMulti.Add(ButtonsMulti[index])
		}
	}

	return
}

func loadRscWithID(id uint16) (string, fyne.Resource) {
	var (
		idStr               = strconv.Itoa(int(id))
		imgPath, nameString string
		tp                  = GetAtbIDType(id)
	)

	if tp == 0 {
		myApp.ErrorLog.Println("error [func loadRscWithID()]: Invalid ID")
		return "", nil
	} else if tp == 1 {
		imgPath = BaseImgPath + idStr + ".png"
	} else {
		imgPath = MultiImgPath + idStr + ".png"
	}

	nameString, isExist := AtbNameMap[id]
	if !isExist {
		myApp.ErrorLog.Println("error [func showMainUI()]: Failed to get Name!")
		return "", nil
	}

	imgRsc, err := fyne.LoadResourceFromPath(imgPath)
	if err != nil {
		myApp.ErrorLog.Println("error [func showMainUI()]: Failed to get Image!")
		return "", nil
	}

	return nameString, imgRsc
}

func processButtonsForInitContainer(index uint16, vbox *fyne.Container, atbTab *container.AppTabs) {
	defer wg.Done()

	nameString, imgRsc := loadRscWithID(index)

	temp := widget.NewButtonWithIcon(nameString, imgRsc, atbButtonMainClicked(vbox, atbTab, index))
	temp.Alignment = widget.ButtonAlignLeading

	lock.Lock()
	if index < 1000 {
		ButtonsBase[index] = temp
	} else {
		ButtonsMulti[index] = temp
	}
	lock.Unlock()
}

func returnButtonClicked(atbTab *container.AppTabs, vbox *fyne.Container) func() {
	return func() {
		if atbTab.Hidden {
			vbox.RemoveAll()
			vbox.Hide()
			atbTab.Show()
		} else {
			atbTab.Hide()
			vbox.Show()
		}
	}
}

func showMainUI() *fyne.Container {
	var (
		CtnBase         = container.NewGridWrap(ButtonSize)
		CtnMulti        = container.NewGridWrap(ButtonSize)
		baseIcon, _     = fyne.LoadResourceFromPath(BaseIconPath)
		multiIcon, _    = fyne.LoadResourceFromPath(MultiIconPath)
		tabItemAtbBase  = container.NewTabItemWithIcon("单属性", baseIcon, container.NewVBox(CtnBase, gradientLineBlackToSkyBlue))
		tabItemAtbMulti = container.NewTabItemWithIcon("双属性", multiIcon, container.NewVBox(CtnMulti, gradientLineBlackToSkyBlue))
		atbTab          = container.NewAppTabs(tabItemAtbBase, tabItemAtbMulti)
		vbox            = container.NewVBox()
		retBox          = container.NewMax(atbTab, vbox)
	)
	initContainer(CtnBase, CtnMulti, vbox, atbTab)

	retBox.Objects[1].Hide()
	return retBox
}

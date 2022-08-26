package main

import (
	"io/fs"
	"log"
	"os/exec"
	"path/filepath"

	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func loadUI() fyne.CanvasObject{
	// Input
	pathInput := widget.NewEntry()
	pathInput.SetPlaceHolder("Enter path...")

	contentText := widget.NewLabel("Please select a path")

	// PathList
	searchPath := "../"
	paths, output := getGitFiles(searchPath)
	pathList := widget.NewList(
		func() int {
			return len(paths)
		}, 
		func() fyne.CanvasObject {
		return widget.NewLabel("Git paths")
	}, 
	func(i widget.ListItemID, o fyne.CanvasObject) {
		o.(*widget.Label).SetText(paths[i])
	})

	// Search Button
	btn := widget.NewButton("search", func() {
		contentText.Text = "Loading..."
		contentText.Refresh()
		searchPath = pathInput.Text
		paths, output = getGitFiles(searchPath)
		contentText.Text = output[0]
		contentText.Refresh()
		pathList.Refresh()
	})

	// View
	search := container.NewGridWithRows(2, pathInput, btn)
	view := container.New(layout.NewBorderLayout(contentText, search, nil, nil),
	search, contentText, pathList)

	// Select Path
	pathList.OnSelected = func(id widget.ListItemID) {
		contentText.Text = output[id]
		contentText.Refresh()
	}
	return view
}

func getGitFiles(path string) (paths []string, output []string) {
    fsys := os.DirFS(path)
    fs.WalkDir(fsys, ".", func(p string, d fs.DirEntry, err error) error {
        if filepath.Ext(p) == ".git" {
			info := "%s %n %ar %n %cn"
			gitPath := path + "/" + p
			cmd := exec.Command("git", "log", "-1", "--stat", "--pretty=format:"+info)
			cmd.Dir = gitPath
			out, err := cmd.Output()
			if err != nil {
				log.Println(err)
			}
			output = append(output, string(out))
			paths = append(paths, p)
        }
        return nil
    })
    return paths, output
}

func main(){
	a:= app.New()
	w:= a.NewWindow("go-git-local")
	w.Resize(fyne.NewSize(800, 800))
	w.SetContent(loadUI())
	w.ShowAndRun()
}


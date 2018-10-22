package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/sahilm/fuzzy"

	"github.com/aswinmprabhu/github-pr-cli/utils"

	"github.com/jroimartin/gocui"
	"github.com/spf13/cobra"
)

type ListPRsReq struct {
	State string `json:"state"`
	Base  string `json:"base"`
}

type ListIssuesReq struct {
	State     string `json:"state"`
	Mentioned string `json:"mentioned"`
}

var g *gocui.Gui
var matches fuzzy.Matches
var resJson = make([]map[string]interface{}, 0)

var titleBytes string
var titles []string

// rootCmd is the main "ghpr" command
var browseCmd = &cobra.Command{
	Use:   "browse [options]",
	Short: "browse issues and pull requests using fuzzy search",
	Run: func(cmd *cobra.Command, args []string) {
		if Issues {
			newListIssuesReq := ListIssuesReq{State: "open"}
			if Mine {
				headremote, err := utils.ParseRemote("origin")
				if err != nil {
					log.Fatal(err)
				}
				userName := strings.Split(headremote, "/")[0]
				newListIssuesReq.Mentioned = userName
				fmt.Println(newListIssuesReq)

			}
			baseremote, err := utils.ParseRemote(strings.Split(Base, ":")[0])
			if err != nil {
				log.Fatal(err)
			}
			urlStr := fmt.Sprintf("https://api.github.com/repos/%s/issues", baseremote)

			// marshal the ListPRsReq
			jsonObj, _ := json.Marshal(&newListIssuesReq)
			client := &http.Client{}
			r, _ := http.NewRequest("GET", urlStr, bytes.NewBuffer(jsonObj)) // URL-encoded payload
			// set the headers
			AuthVal := fmt.Sprintf("token %s", Token)
			r.Header.Add("Authorization", AuthVal)
			r.Header.Add("Content-Type", "application/json")

			// make the req
			resp, err := client.Do(r)
			if err != nil {
				log.Fatal(err)
			}
			// defer resp.Body.Close()
			fmt.Println("Fetching....")
			bytes, _ := ioutil.ReadAll(resp.Body)
			if err := json.Unmarshal(bytes, &resJson); err != nil {
				log.Fatalf("Failed to parse the response : %v", err)
			}

			for _, i := range resJson {
				titleBytes = titleBytes + i["title"].(string) + "\n"
				titles = append(titles, i["title"].(string))
			}

		} else {
			newListPRsReq := ListPRsReq{State: "open"}
			baseremote, err := utils.ParseRemote(strings.Split(Base, ":")[0])
			if err != nil {
				log.Fatal(err)
			}
			urlStr := fmt.Sprintf("https://api.github.com/repos/%s/pulls", baseremote)
			newListPRsReq.Base = strings.Split(Base, ":")[1]

			// marshal the ListPRsReq
			jsonObj, _ := json.Marshal(&newListPRsReq)
			client := &http.Client{}
			r, _ := http.NewRequest("GET", urlStr, bytes.NewBuffer(jsonObj)) // URL-encoded payload
			// set the headers
			AuthVal := fmt.Sprintf("token %s", Token)
			r.Header.Add("Authorization", AuthVal)
			r.Header.Add("Content-Type", "application/json")

			// make the req
			resp, err := client.Do(r)
			if err != nil {
				log.Fatal(err)
			}
			// defer resp.Body.Close()
			fmt.Println("Fetching....")
			bytes, _ := ioutil.ReadAll(resp.Body)
			if err := json.Unmarshal(bytes, &resJson); err != nil {
				log.Fatalf("Failed to parse the response : %v", err)
			}

			for _, i := range resJson {
				titleBytes = titleBytes + i["title"].(string) + "\n"
				titles = append(titles, i["title"].(string))
			}
		}
		var err error
		g, err = gocui.NewGui(gocui.OutputNormal)
		if err != nil {
			log.Panicln(err)
		}
		defer g.Close()

		g.Cursor = true
		g.Mouse = false

		g.SetManagerFunc(layout)

		if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
			log.Panicln(err)
		}

		if err := g.SetKeybinding("finder", gocui.KeyArrowRight, gocui.ModNone, switchToMainView); err != nil {
			log.Panicln(err)
		}

		if err := g.SetKeybinding("main", gocui.KeyArrowLeft, gocui.ModNone, switchToSideView); err != nil {
			log.Panicln(err)
		}

		if err := g.SetKeybinding("", gocui.KeyArrowDown, gocui.ModNone, cursorDown); err != nil {
			log.Panicln(err)
		}
		if err := g.SetKeybinding("", gocui.KeyArrowUp, gocui.ModNone, cursorUp); err != nil {
			log.Panicln(err)
		}

		if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
			log.Panicln(err)
		}

	},
}

func cursorDown(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy+1); err != nil {
			ox, oy := v.Origin()
			if err := v.SetOrigin(ox, oy+1); err != nil {
				return err
			}
		}
	}
	return nil
}

func cursorUp(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		ox, oy := v.Origin()
		cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy-1); err != nil && oy > 0 {
			if err := v.SetOrigin(ox, oy-1); err != nil {
				return err
			}
		}
	}
	return nil
}

func switchToSideView(g *gocui.Gui, view *gocui.View) error {
	if _, err := g.SetCurrentView("finder"); err != nil {
		return err
	}
	return nil
}

func switchToMainView(g *gocui.Gui, view *gocui.View) error {
	if _, err := g.SetCurrentView("main"); err != nil {
		return err
	}
	return nil
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("finder", -1, 0, 80, 10); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Wrap = true
		v.Editable = true
		v.Frame = true
		v.Title = "Type pattern here. Press -> or <- to switch between panes"
		if _, err := g.SetCurrentView("finder"); err != nil {
			return err
		}
		v.Editor = gocui.EditorFunc(finder)
	}
	if v, err := g.SetView("main", 79, 0, maxX, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		fmt.Fprintf(v, "%s", titleBytes)
		v.Editable = false
		v.Wrap = true
		v.Frame = true
	}

	if v, err := g.SetView("results", -1, 3, 79, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Editable = false
		v.Wrap = true
		v.Frame = true
		v.Title = "Search Results"
	}

	return nil
}

func getURL(title string) (string, error) {
	for _, item := range resJson {
		if item["title"] == title {
			return item["html_url"].(string), nil
		}
	}
	return "", fmt.Errorf("%s : Not Found", title)
}

func quit(g *gocui.Gui, v *gocui.View) error {
	url, err := getURL(matches[0].Str)
	if err != nil {
		return err
	}
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}
	return gocui.ErrQuit
}

func finder(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
	switch {
	case ch != 0 && mod == 0:
		v.EditWrite(ch)
		g.Update(func(gui *gocui.Gui) error {
			results, err := g.View("results")
			if err != nil {
				// handle error
			}
			results.Clear()
			t := time.Now()
			matches = fuzzy.Find(strings.TrimSpace(v.ViewBuffer()), titles)
			elapsed := time.Since(t)
			fmt.Fprintf(results, "found %v matches in %v\n", len(matches), elapsed)
			for _, match := range matches {
				for i := 0; i < len(match.Str); i++ {
					if contains(i, match.MatchedIndexes) {
						fmt.Fprintf(results, fmt.Sprintf("\033[1m%s\033[0m", string(match.Str[i])))
					} else {
						fmt.Fprintf(results, string(match.Str[i]))
					}

				}
				fmt.Fprintln(results, "")
			}
			return nil
		})
	case key == gocui.KeySpace:
		v.EditWrite(' ')
	case key == gocui.KeyBackspace || key == gocui.KeyBackspace2:
		v.EditDelete(true)
		g.Update(func(gui *gocui.Gui) error {
			results, err := g.View("results")
			if err != nil {
				// handle error
			}
			results.Clear()
			t := time.Now()
			matches := fuzzy.Find(strings.TrimSpace(v.ViewBuffer()), titles)
			elapsed := time.Since(t)
			fmt.Fprintf(results, "found %v matches in %v\n", len(matches), elapsed)
			for _, match := range matches {
				for i := 0; i < len(match.Str); i++ {
					if contains(i, match.MatchedIndexes) {
						fmt.Fprintf(results, fmt.Sprintf("\033[1m%s\033[0m", string(match.Str[i])))
					} else {
						fmt.Fprintf(results, string(match.Str[i]))
					}
				}
				fmt.Fprintln(results, "")
			}
			return nil
		})
	case key == gocui.KeyDelete:
		v.EditDelete(false)
		g.Update(func(gui *gocui.Gui) error {
			results, err := g.View("results")
			if err != nil {
				// handle error
			}
			results.Clear()
			t := time.Now()
			matches := fuzzy.Find(strings.TrimSpace(v.ViewBuffer()), titles)
			elapsed := time.Since(t)
			fmt.Fprintf(results, "found %v matches in %v\n", len(matches), elapsed)
			for _, match := range matches {
				for i := 0; i < len(match.Str); i++ {
					if contains(i, match.MatchedIndexes) {
						fmt.Fprintf(results, fmt.Sprintf("\033[1m%s\033[0m", string(match.Str[i])))
					} else {
						fmt.Fprintf(results, string(match.Str[i]))
					}
				}
				fmt.Fprintln(results, "")
			}
			return nil
		})
	case key == gocui.KeyInsert:
		v.Overwrite = !v.Overwrite
	}
}

func contains(needle int, haystack []int) bool {
	for _, i := range haystack {
		if needle == i {
			return true
		}
	}
	return false
}

var Mine bool
var Issues bool

func init() {
	// define flags
	f := browseCmd.PersistentFlags()
	f.BoolVarP(&Mine, "mine", "m", false, "Mine")
	f.BoolVarP(&Issues, "issues", "i", false, "Issues")

	rootCmd.AddCommand(browseCmd)
}

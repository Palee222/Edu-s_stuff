package main

import (
	"bufio"
	"flag"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/GoneUp/mask-go/mask"

	"encoding/hex"

	"github.com/icza/gox/imagex/colorx"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetLevel(log.DebugLevel)
}

// text3 qgQZ´
func main() {
	drawMode := flag.Bool("draw", false, "enables drawmode")
	text := flag.String("text", "test", "string to draw")
	removePhoto := flag.String("remove-photo", "", "deletes a local image file without connecting to the mask")
	flag.Parse()

	if *removePhoto != "" {
		if err := removePhotoFile(*removePhoto); err != nil {
			log.Fatalf("Remove photo: %v", err)
		}
		log.Infof("Removed photo file %q", *removePhoto)
		return
	}

	if *drawMode {
		log.Info("drawing ", *text)
		mask.GetTextImage(*text)
		return
	}

	log.SetReportCaller(true)
	log.Info("MaskCmd started. Connecting to mask...")
	err := mask.InitAndConnect(true)
	if err != nil {
		log.Fatalf("Init failed with error %v", err)
	}

	doDemoControlLoop()
}

func doDemoControlLoop() {
	//simple cmd struct
	for {
		log.Info("Input cmd please:")
		reader := bufio.NewReader(os.Stdin)
		cmd, _ := reader.ReadString('\n')
		cmdSplit := strings.Split(strings.Trim(cmd, "\n"), " ")

		switch cmdSplit[0] {
		case "connect":
			err := mask.InitAndConnect(true)
			if err != nil {
				log.Fatalf("Init failed with error %v", err)
			}

		case "allmode":
			for i := 1; i < 5; i++ {
				log.Infof("trying to send mode %d\n", i)
				mask.SetMode(byte(i))
				time.Sleep(5 * time.Second)
			}

		case "mode":
			val, err := strconv.ParseUint(cmdSplit[1], 10, 8)
			must("ParseInt", err)

			mask.SetMode(byte(val))

		case "light":
			val, err := strconv.ParseUint(cmdSplit[1], 10, 8)
			must("ParseInt", err)

			mask.SetLight(byte(val))

		case "image":
			val, err := strconv.ParseUint(cmdSplit[1], 10, 8)
			must("ParseInt", err)

			mask.SetImage(byte(val))

		case "diy":
			val, err := strconv.ParseUint(cmdSplit[1], 10, 8)
			must("ParseInt", err)

			mask.SetDIYImage(byte(val))

		case "speed":
			val, err := strconv.ParseUint(cmdSplit[1], 10, 8)
			must("ParseInt", err)

			mask.SetTextSpeed(byte(val))

		case "color":
			val, err := strconv.ParseUint(cmdSplit[1], 10, 8)
			must("ParseInt", err)

			mask.SetTextColorMode(1, byte(val))

		case "fg":
			c, err := colorx.ParseHexColor(cmdSplit[1])
			if err == nil {
				mask.SetTextFrontColor(1, c.R, c.G, c.B)
			} else {
				log.Warn("wrong format, use #FFFFFF")
			}

		case "bg":
			c, err := colorx.ParseHexColor(cmdSplit[1])
			must("ParseHexColor", err)

			if err == nil {
				mask.SetTextBackgroundColor(1, c.R, c.G, c.B)
			} else {
				log.Warn("wrong format, use #FFFFFF")
			}

		case "text":
			//simple red white stripes
			bitmap := []byte{0xFF, 0xFF, 0x00, 0x00, 0xFF, 0xFF, 0x00, 0x00, 0xFF, 0xFF, 0x00, 0x00}
			colorArray := []byte{0xFF, 0x00, 0x00, 0xFF, 0x00, 0x00,
				0x00, 0xFF, 0x00, 0x00, 0xFF, 0x00,
				0xFF, 0x00, 0x00, 0xFF, 0x00, 0x00,
				0x00, 0xFF, 0x00, 0x00, 0xFF, 0x00}

			mask.InitUpload(bitmap, colorArray)
		case "text2":
			//decoded from bt dump, sets test as text
			bitmap, err := hex.DecodeString("020002003ff83ffc020402040000000000f001f8034c0244034401cc00c80000018803cc024402640224033c01180000020002003ff83ffc0204020400000000")
			must("decode", err)
			colorArray, err := hex.DecodeString("fffffcfffffcfffffcfffffcfffffcfffffcfffffcfffffcfffffcfffffcfffffcfffffcfffffcfffffcfffffcfffffcfffffcfffffcfffffcfffffcfffffcfffffcfffffcfffffcfffffcfffffcfffffcfffffcfffffcfffffcfffffcfffffc")
			must("decode", err)

			mask.InitUpload(bitmap, colorArray)

		case "text3":
			//text3 jonasjemmy hat Köngisrose x1 gesendet
			text := strings.TrimPrefix(cmd, "text3 ")
			text = strings.TrimSpace(text)

			mask.SetText(text)

		case "photo":
			path, err := commandPath(cmd)
			if err != nil {
				log.Errorf("Upload photo: %v", err)
				continue
			}
			if err := mask.SetImageFromFile(path); err != nil {
				log.Errorf("Upload photo: %v", err)
			}

		case "clearphoto":
			if err := mask.ClearImage(); err != nil {
				log.Errorf("Clear photo: %v", err)
			}

		case "removephoto":
			path, err := commandPath(cmd)
			if err != nil {
				log.Errorf("Remove photo: %v", err)
				continue
			}
			if err := removePhotoFile(path); err != nil {
				log.Errorf("Remove photo: %v", err)
			} else {
				log.Infof("Removed photo file %q", path)
			}

		case "exit":
			mask.Shutdown()
			os.Exit(0)

		default:
			log.Info("Unknown cmd")
		}
	}
}

func commandPath(cmd string) (string, error) {
	_, path, found := strings.Cut(strings.TrimSpace(cmd), " ")
	path = strings.TrimSpace(path)
	if !found || path == "" {
		return "", fmt.Errorf("provide an image file path")
	}
	if strings.HasPrefix(path, `"`) {
		unquoted, err := strconv.Unquote(path)
		if err != nil {
			return "", fmt.Errorf("invalid quoted path: %w", err)
		}
		path = unquoted
	}
	return path, nil
}

func removePhotoFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("open %q: %w", path, err)
	}
	_, _, decodeErr := image.DecodeConfig(file)
	closeErr := file.Close()
	if decodeErr != nil {
		return fmt.Errorf("%q is not a supported image: %w", path, decodeErr)
	}
	if closeErr != nil {
		return fmt.Errorf("close %q: %w", path, closeErr)
	}
	if err := os.Remove(path); err != nil {
		return fmt.Errorf("delete %q: %w", path, err)
	}
	return nil
}

func must(action string, err error) {
	if err != nil {
		panic("failed to " + action + ": " + err.Error())
	}
}

# mask-go

[Play store link of the control app](https://play.google.com/store/apps/details?id=cn.com.heaton.shiningmask)

Features: 
- Connection with mask over BLE (using tinygo.org/x/bluetooth)
- Controlling the mask remotely (Brightness, Static image, Anmination, Text speed, Text color) 
- Uploading/showing text on the mask 

![](image.png)
## Usage
Installation with
`go get github.com/GoneUp/mask-go`


Simple usage example
```go
import "github.com/GoneUp/mask-go/mask"

mask.InitAndConnect(true)
mask.SetText("Hello world")
```

How to run it (my own notes):
- go run .

| Command | What it does |
|---|---|
| `connect` | Connects to the mask again. |
| `allmode` | Cycles through modes 1–4, pausing 5 seconds in each. |
| `mode <number>` | Sets the display mode. The code documents modes 1–5: steady, blink, scroll left, scroll right, steady. |
| `light <number>` | Sets brightness; accepts a byte value (0–255). |
| `image <number>` | Selects a built-in still image. |
| `diy <number>` | Selects a DIY image slot. |
| `speed <number>` | Sets text speed; accepts 0–255. |
| `color <number>` | Sets the text color mode; documented modes are 0–7. |
| `fg <#RRGGBB>` | Sets the text foreground color, e.g. `fg #FF0000`. |
| `bg <#RRGGBB>` | Sets the text background color, e.g. `bg #000000`. |
| `text` | Uploads a hard-coded test pattern. |
| `text2` | Uploads another hard-coded test image. |
| `text3 <message>` | Displays your message, e.g. `text3 Hello mask`. |
| `exit` | Disconnects and quits. |

These are the commands implemented in `main.go`. 

There is **no interactive `animation` command**; the library has a `SetAnimation` function, but the CLI doesn’t expose it. The startup flags `-draw` and `-text` are separate from the interactive commands.

Todo:
- upload photo from the code itself or from vscode
- animation
- be able to run a preset mode for some time
- real time image upload, animation and change display

To upload your own image, the project would need a new command that reads an image file, converts it to the mask’s format, and sends it over Bluetooth. The mask format is limited: it’s 16 pixels high, with one RGB color per vertical column. The current commands are in main.go, and the upload function is in mask.go.


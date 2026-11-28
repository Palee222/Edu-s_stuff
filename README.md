**Todo:**  
- real time image upload, animation and change display
- finding a way to put videos (i.e., very long gifs)
- even connect it to other software, like touchdesigner
  - control with hands 

# mask-go

Features:
- Connection with mask over BLE (using tinygo.org/x/bluetooth)
- Controlling the mask remotely (Brightness, Static image, Anmination, Text speed, Text color)
- Uploading/showing text on the mask

![](image.png)
## Usage

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
| `mode <number>` | Sets the display mode; see the mode list below. |
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
| `photo <path>` | Uploads an image file, e.g. `photo photosforuse/image.png`. Quote paths containing spaces. |
| `animation <path>` | Plays an animated GIF from a file, e.g. `animation photosforuse/mask-test-moving-bars.gif`. |
| `clearphoto` | Uploads a blank image to clear the currently uploaded image on the mask. |
| `removephoto <path>` | Deletes the specified local image file after checking that it is a supported image. |
| `exit` | Disconnects and quits. |
| animation <path> | It composes GIF frames, converts them to the mask’s format, sends them sequentially, and follows the GIF’s timing and loop settings. Playback blocks until the GIF finishes; press Ctrl+C to stop one that loops forever. |

### Display modes

Enter `mode <number>` at the program prompt:

| Mode | Behavior |
|---|---|
| `1` | Steady display (no scrolling). |
| `2` | Blinks. |
| `3` | Scrolls left. |
| `4` | Scrolls right. |
| `5` | Steady display (same behavior as mode 1). |

The `photo` command reads PNG, JPEG, or GIF files, scales them to 16 pixels high while preserving their aspect ratio, and uploads them over Bluetooth. The mask uses a one-bit on/off value per pixel and one RGB color per vertical column, so the uploaded result is a simplified image rather than a full-color photo. `removephoto` deletes a local source image; `clearphoto` sends a blank image to the mask. The current protocol does not expose a way to delete image data stored in the mask itself.

## Compatible picture formats

The implementation supports **PNG, JPEG/JPG, and GIF** files:

- PNG
- JPEG/JPG
- GIF, but the `photo` command displays only the **first frame**
- BMP, TIFF, SVG, WebP, and similar formats are **not supported** by the current implementation, despite some outdated text in `README.md`.

The format is detected from the file contents rather than only the filename. This is defined by the decoders imported in `mask/image.go`.

## Exact dimensions and scaling

The mask output is always:

- **16 pixels high**
- **Up to 4998 columns wide**
- **One bit per pixel**
- **One RGB color per column**

The source image is resized to exactly 16 pixels high while preserving aspect ratio:

```text
output columns = round(source width × 16 ÷ source height)
```

Additional constraints:

- The result is at least **1 column** wide.
- The result cannot exceed **4998 columns**.
- Images taller than 4998 times their width will be rejected after scaling.
- Both dimensions must be greater than zero.

For example:

| Source size | Output size |
|---|---|
| 1000 × 1000 | 16 × 16 |
| 64 × 16 | 64 × 16 |
| 1920 × 1080 | 28 × 16 |
| 4998 × 16 | 4998 × 16 |
| 5000 × 16 | Rejected |

This behavior is implemented in `mask/image.go`.

## Pixel conversion criteria

The software does not upload true grayscale or full-resolution color:

1. The image is resized using Catmull-Rom resampling.
2. Each pixel is converted to a luminance value:

```text
luminance = (299 × red + 587 × green + 114 × blue) ÷ 1000
```

3. Pixels with luminance **64 or lower** are turned off.
4. Pixels with luminance **65 or higher** are turned on.
5. Pixels are never represented as gray levels.
6. For each column, the displayed RGB color is the average of the red, green, and blue values of its lit pixels.
7. A column with no lit pixels becomes pure black: `RGB(0, 0, 0)`.

For that reason, pictures with strong contrast and simple, bright subjects work much better than detailed photographs or low-contrast images.

## Animated GIF criteria

Use:

```text
animation photosforuse/mask-test-moving-bars.gif
```

For GIF animation:

- Each frame is composited and converted using the same 16-pixel-high and 4998-column limit.
- Each frame is sent as a separate mask image.
- Frame timing comes from GIF delay values.
- Missing or zero-delay frames use a **100 ms** default.
- GIF loop settings are respected.
- GIFs without a loop limit play indefinitely until interrupted.
- The display blocks until playback finishes.

This behavior is implemented in `mask/animation.go`.

If you need to install go: 
  https://go.dev/dl/
  
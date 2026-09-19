# image-watermark

A command-line tool that batch-resizes `.jpg` photos and stamps a watermark
onto them. It reads each photo's EXIF data to figure out its camera brand
(currently Canon and Fujifilm are supported) and orientation, so portrait
and landscape shots are resized correctly (long edge capped at a target
size) regardless of how the camera or an editor like Photoshop stored the
rotation.

For the full design (pipeline, orientation-detection logic, known
limitations), see [SPEC.md](SPEC.md).

## Running the built binary

Pre-built binaries live in `dist/` after running `make` (see
[Building](#building) below), named `watermark-<os>-<arch>` (or
`watermark-windows.exe` on Windows).

```
./watermark -dir <path-to-photos-directory> -watermark <path-to-watermark-image> -size <width pixel size>
```

On Windows, run the `.exe` instead:

```
watermark-windows.exe -dir <path-to-photos-directory> -watermark <path-to-watermark-image> -size <width pixel size>
```

### Flags

| Flag | Required | Default | Description |
|---|---|---|---|
| `-dir` | yes | — | Directory containing the source `.jpg` photos to process. |
| `-watermark` | yes | — | Path to the image overlaid as the watermark. |
| `-size` | no | `2000` | Resize target for the image's long edge, in pixels. |

### Output

Watermarked, resized copies are written to a `watermarked` subfolder
created inside `-dir`, keeping the original file names:

```
<dir>/watermarked/<file>.jpg
```

## Building

A `Makefile` is included to cross-compile for macOS, Linux, and Windows.

```
make build-darwin   # macOS (amd64 + arm64)
make build-linux     # Linux (amd64 + arm64)
make build-windows   # Windows (amd64)
```

Binaries are written to `dist/`. To build only for your current machine:

```
make build
```

This produces a `watermark` binary in the repository root. Run
`make clean` to remove build output.

## Requirements

- Go 1.21+ (only needed to build from source; the built binaries have no
  runtime dependencies).

# Image Watermark — Spec

Module: `github.com/jbernhardsch/image-watermark` (Go 1.21)

## Purpose

A CLI batch tool that walks a directory of `.jpg` photos, resizes each one to
a standard size, and overlays a watermark image on it. The resize target
(width-bound vs. height-bound) depends on whether the photo is landscape or
portrait, and portrait detection/rotation is camera-brand-specific because
different camera makes write EXIF orientation differently.

## Usage

```
go run . -dir <path-to-photos-directory> -watermark <path-to-watermark-image> -size <width pixel size>
```

- `-dir`: directory containing the source `.jpg` files.
- `-watermark`: path to the image overlaid as the watermark.
- `-size`: resize target for the image's long edge, in pixels.

Output is written to `<dir>/watermarked/`, created automatically if it
doesn't exist. Output files keep the original file name.

## Files

| File | Responsibility |
|---|---|
| `main.go` | CLI flags, shared constants, program entry point |
| `process.go` | Directory walk, per-file orchestration |
| `orientation.go` | EXIF-based orientation/rotation detection |
| `resize.go` | Resize (and rotate, if needed) the image |
| `watermark.go` | Overlay the watermark image |

## Pipeline

For each `*.jpg` file in `-dir` (case-insensitive match on extension):

1. **`getOrientation`** (`orientation.go`) opens the file, decodes its EXIF
   data, and returns:
   - an orientation classification (see constants below),
   - a `needRotate bool` — whether the pixel data must still be physically
     rotated 90°,
   - an error.
2. **`resizeImaging`** (`resize.go`) opens the image, picks a target
   `width`/`height` based on the orientation classification (long edge ==
   `RESIZE_VALUE`, other edge computed to preserve aspect ratio), resizes
   with Lanczos filtering, physically rotates the pixels if `needRotate` is
   true, and saves the result to `<dir>/watermarked/<file>`.
3. **`waterMarking`** (`watermark.go`) re-opens the resized file, opens the
   watermark image, and overlays it centered horizontally, 80px above the
   bottom edge, at 50% opacity, saving back over the resized file.

### Orientation constants (`main.go`)

- `IMAGE_NOEXIF` — no `Make` EXIF tag (or non-decodable EXIF); orientation
  is inferred directly from decoded pixel dimensions, no rotation applied.
- `IMAGE_LANDSCAPE` — final image is landscape.
- `IMAGE_PORTRAIT` — final image is portrait (Canon).
- `IMAGE_PORTRAIT_FUJI` / `IMAGE_PORTRAIT_FUJI_SPC` — final image is
  portrait (Fujifilm); `_SPC` denotes a Fuji orientation-tag value other
  than the common `"8"`, which needs a `Rotate270` instead of `Rotate90`.

Only `canon` and `fujifilm` (matched case-insensitively against the EXIF
`Make` tag) have dedicated branches; any other camera make falls through to
the `errors.New("cannot proceed image orientation")` error path.

## Orientation detection logic (current, post-fix)

`getOrientation` reads two independent signals and reconciles them:

1. **EXIF `Orientation` tag** (`ex.Get(exif.Orientation)`) — the camera's
   (or editor's) hint about whether the stored pixels need rotating for
   correct display.
2. **Actual decoded pixel dimensions** (`decodedSize`, via
   `image.DecodeConfig` on the raw file) — the ground truth for what the
   pixel data itself looks like, independent of any EXIF metadata.

Decision rules:

- **Canon**: if the tag is `"8"`, the branch that originally shipped is
  kept (rotate; label as landscape or portrait depending on which raw
  dimension is larger). Otherwise, the raw pixel shape alone decides:
  `width > height` → landscape, no rotate; `width <= height` → portrait,
  no rotate (the pixels are already correctly oriented).
- **Fujifilm**: if the raw pixel shape is already portrait
  (`width < height`), it's classified as portrait with no rotate needed,
  regardless of the tag. Otherwise, the original tag-based rule applies
  (`"1"` → landscape, `"8"` → portrait with `Rotate90`, anything else →
  portrait with `Rotate270`).

`resizeImaging` picks its target dimensions purely from the orientation
classification (`IMAGE_PORTRAIT*` → height-bound, everything else →
width-bound, `IMAGE_NOEXIF` → decode the file directly). Whether a rotate
is actually performed is controlled solely by `needRotate`, applied after
resizing.

## Bug fixed: Photoshop-edited photos failing portrait detection

### Symptom

Photos re-saved out of Photoshop (crop, retouch, export, etc.) were not
recognized as portrait: they either failed processing outright, or were
resized/rotated incorrectly (ending up sideways, or resized against the
wrong edge so the long edge overshot `RESIZE_VALUE`).

### Root cause

When Photoshop saves a JPEG that had a non-trivial EXIF `Orientation` tag,
it commonly **bakes the rotation into the pixel data itself** and resets
(or removes) the `Orientation` tag — the file no longer says "please
rotate me," because it's already rotated. Photoshop also does not reliably
keep the EXIF `PixelXDimension`/`PixelYDimension` sub-tags in sync with the
edited pixel data (they can go stale, or be dropped entirely).

The pre-fix code assumed neither of these could happen:

1. It read the "raw" dimensions from `exif.PixelXDimension` /
   `exif.PixelYDimension` rather than the actual pixel data. If Photoshop
   dropped those tags, `getOrientation` returned an error and the file
   failed outright. If they went stale, dimension-based decisions were
   made against numbers that no longer matched the real pixels.
2. In the Canon fallback branch (`orientation.go`, old code), *any* file
   with `xSize < ySize` was unconditionally treated as needing rotation
   (`return IMAGE_PORTRAIT, true, nil`). That assumption holds for
   straight-off-the-camera Canon files (the sensor always records
   landscape-shaped pixels; portrait is only ever signaled via the
   `Orientation` tag), but breaks for a Photoshop-edited file: its pixels
   are already portrait-shaped, so rotating again produces a sideways
   image.
3. Fujifilm's branch didn't consult pixel dimensions at all — only the
   `Orientation` tag. A Photoshop-reset tag of `"1"` made an
   already-portrait, already-correctly-oriented image get classified as
   `IMAGE_LANDSCAPE`, so it was resized width-bound instead of
   height-bound.
4. Separately, `resizeImaging`'s width/height selection was keyed off a
   combination of `orientation` *and* `needRotate`, instead of
   `orientation` alone. This meant even a correctly classified
   `IMAGE_PORTRAIT`/`IMAGE_PORTRAIT_FUJI*` image with `needRotate == false`
   fell through to the width-bound (landscape) resize branch, since no
   condition matched it.

### Fix

- Added `decodedSize` (`orientation.go`), which decodes the actual image
  bytes (`image.DecodeConfig`) to get true pixel dimensions, immune to
  stale or missing EXIF dimension tags.
- Canon and Fujifilm branches now treat "pixels are already portrait
  shaped" (`width < height` from `decodedSize`) as its own case: portrait,
  `needRotate = false`. This directly covers the Photoshop bake-in
  scenario, while native/unedited camera files (which are never portrait
  raw-shaped) are unaffected and keep their original tag-driven behavior.
- `resizeImaging` now selects its target width/height purely from the
  orientation classification (`switch orientation { ... }`), independent
  of `needRotate`. Rotation remains a separate, later step gated solely by
  `needRotate`.

### Known limitation (pre-existing, not changed by this fix)

Canon's `Orientation` tag value `"6"` (used by some bodies/lenses to mean
"rotate 90° for portrait display," as an alternative to `"8"`) is not
special-cased and falls through to the raw-dimension comparison. For an
unedited camera-native file this yields `IMAGE_LANDSCAPE` with no rotate,
which is likely wrong. This was already the behavior before this fix and
is out of scope for the Photoshop issue; flagged here for future work.

## Dependencies

- [`github.com/disintegration/imaging`](https://github.com/disintegration/imaging) —
  resize (Lanczos), rotate, overlay, open/save.
- [`github.com/rwcarlsen/goexif`](https://github.com/rwcarlsen/goexif)
  (`exif`, `mknote`) — EXIF decoding, including Canon/Nikon maker notes.
- `image`, `image/jpeg` (stdlib) — raw pixel-dimension decoding used as the
  source of truth for orientation.

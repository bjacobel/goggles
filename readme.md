###goggles
---

Goggles is a Go interface for image classification.

See `twitter.go` for an example of using Goggles to classify images from the Twitter Streaming API.

####Requirements
- `imagemagick`
- OverFeat weight files (run `python OverFeat/download_weights.py`)
- OverFeat compiled (`git submodule init && cd OverFeat/src && make`)
    - (In order to actually get this to compile you will have to do some mucking. It's not easy. Sorry!)

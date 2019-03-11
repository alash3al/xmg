XMG
====
> Xmg a tiny micro service that translates the uploaded image to a combination of hashes that represents it also can extract separated hashes for the faces in the image

Requirements
=============
> `xmg` depends on `dlib` library, so you can follow [the guide from here](https://github.com/Kagami/go-face#requirements)

Install
=======
> from source using `Golang` toolchain `go get github.com/alash3al/xmg`

Usage
=====
> just `xmg --dlib-models="./DLIB_MODELS_PATH_HERE" --listen=:9020` then `curl -F "image=@path_to_some_img" http://localhost:9020/all` and see the result.
> You can replace `all` with `faces` in the url to return only the faces hashes in the image.

Author
======
Mohamed Al Ashaal, a senior software engineer :)
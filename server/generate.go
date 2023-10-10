package main

//go:generate go get github.com/josephspurrier/goversioninfo/cmd/goversioninfo
//go:generate go install github.com/josephspurrier/goversioninfo/cmd/goversioninfo
//go:generate go get -d
//go:generate go run genversioninfo/genversioninfo.go versioninfo.json
//go:generate goversioninfo -gofile=versioninfo.go -platform-specific=true
//go:generate go get

//go:generate go get github.com/gomarkdown/markdown
//go:generate go run md2html/md2html.go onlineHilfe/md-de onlineHilfe/html-de onlineHilfe/css/qt.css
//go:generate go run md2html/md2html.go onlineHilfe/md-en onlineHilfe/html-en onlineHilfe/css/qt.css

//go:generate go run unzip/unzip.go ../client/manueltest/manuelle_test_web.zip flutterCfg

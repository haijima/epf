package epf

import (
	"fmt"
	"log/slog"

	"github.com/haijima/analysisutil"
)

type WebFramework int

const (
	EchoV4 WebFramework = iota
	Gin
	ChiV5
	Iris12
	Gorilla
	NetHttp
	None
)

func DetectImported(dir, pattern string) (WebFramework, error) {
	pkgs, err := analysisutil.LoadPackages(dir, pattern)
	if err != nil {
		return None, err
	}

	var useNetHttp bool
	for _, pkg := range pkgs {
		for _, p := range pkg.Imports {
			if p.PkgPath == "github.com/labstack/echo/v4" {
				return EchoV4, nil
			} else if p.PkgPath == "github.com/gin-gonic/gin" {
				return Gin, nil
			} else if p.PkgPath == "github.com/go-chi/chi/v5" {
				return ChiV5, nil
			} else if p.PkgPath == "github.com/kataras/iris/v12" {
				return Iris12, nil
			} else if p.PkgPath == "github.com/gorilla/mux" {
				return Gorilla, nil
			} else if p.PkgPath == "net/http" {
				useNetHttp = true
			}
		}
	}
	if useNetHttp {
		return NetHttp, nil
	}
	return None, nil
}

func AutoExtractor(dir, pattern string) (Extractor, error) {
	usedFramework, err := DetectImported(dir, pattern)
	if err != nil {
		return nil, err
	}
	switch usedFramework {
	case EchoV4:
		slog.Info("Detected Echo: \"github.com/labstack/echo/v4\"")
		return &EchoExtractor{}, nil
	case Gin:
		slog.Info("Detected Gin: \"github.com/gin-gonic/gin\"")
		return nil, fmt.Errorf("unsupported framework: %v", usedFramework)
	case ChiV5:
		slog.Info("Detected go-chi: \"github.com/go-chi/chi/v5\"")
		return nil, fmt.Errorf("unsupported framework: %v", usedFramework)
	case Iris12:
		slog.Info("Detected Iris: \"github.com/kataras/iris/v12\"")
		return nil, fmt.Errorf("unsupported framework: %v", usedFramework)
	case Gorilla:
		slog.Info("Detected Gorilla: \"github.com/gorilla/mux\"")
		return nil, fmt.Errorf("unsupported framework: %v", usedFramework)
	case NetHttp:
		slog.Info("Detected \"net/http\"")
		return &NetHttpExtractor{}, nil
	case None:
		return nil, fmt.Errorf("not found web framework from %s", dir)
	}
	return nil, err
}

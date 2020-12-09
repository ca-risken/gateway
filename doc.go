package main

import (
	"io/ioutil"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/docgen"
)

func updateDoc(router *chi.Mux) error {
	path, err := os.Getwd()
	if err != nil {
		return err
	}
	doc := genDoc("githug.com/CyberAgent/mimosa-gateway", router)
	if err := ioutil.WriteFile(path+"/doc/README.md", []byte(doc), 0666); err != nil {
		return err
	}
	return nil
}

func genDoc(pkgName string, router *chi.Mux) string {
	return docgen.MarkdownRoutesDoc(router, docgen.MarkdownOpts{
		ProjectPath: pkgName,
		Intro:       "MIMOSA API document by go-chi.",
		// ForceRelativeLinks: true,
		URLMap: map[string]string{
			"githug.com/CyberAgent/mimosa-gateway": "https://githug.com/CyberAgent/mimosa-gateway/blob/master",
			"github.com/go-chi/chi":                "https://github.com/go-chi/chi/blob/master/",
		},
	})
}

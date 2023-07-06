package generator

import (
	"bytes"
	"context"
	"fmt"
	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"html/template"
	"time"
)

const (
	templatePathHtml = "./template.html"
	outputPath       = "./output"
)

type docGenerator struct {
}

func NewDocumentGenerator() *docGenerator {
	return &docGenerator{}
}

func (g *docGenerator) Parser(ctx context.Context, templateFileName string, data interface{}) ([]byte, error) {

	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (g *docGenerator) Resolve(ctx context.Context, data interface{}) error {
	t := time.Now().Unix()

	//we need convert parsedFile from html to pdf here
	parsedFile, err := g.Parser(ctx, templatePathHtml, data)
	if err != nil {
		return err
	}

	// convert byte slice to io.Reader
	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {

		return err
	}

	pdfg.AddPage(wkhtmltopdf.NewPageReader(bytes.NewReader(parsedFile)))

	pdfg.PageSize.Set(wkhtmltopdf.PageSizeA4)

	pdfg.Dpi.Set(300)

	err = pdfg.Create()
	if err != nil {
		return err
	}

	err = pdfg.WriteFile(fmt.Sprintf(`%s/%v.pdf`, outputPath, t))
	if err != nil {
		return err
	}

	return nil
}

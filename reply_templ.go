// Code generated by templ@(devel) DO NOT EDIT.

package main

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

import "strings"

// Debug thign
// <input hx-ext="debug" hx-swap-oob="outerHTML:#test" name="context" class="context" value={ contextArr }/>

func botMessage(message, contextArr string) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		templBuffer, templIsBuffer := w.(*bytes.Buffer)
		if !templIsBuffer {
			templBuffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templBuffer)
		}
		ctx = templ.InitializeContext(ctx)
		var_1 := templ.GetChildren(ctx)
		if var_1 == nil {
			var_1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, err = templBuffer.WriteString("<input id=\"test\" hx-ext=\"debug\" hx-swap-oob=\"outerHTML:#test\" name=\"context\" value=\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(contextArr))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\"><div hx-swap-oob=\"beforeend:#content\"><div class=\"bot-message\">")
		if err != nil {
			return err
		}
		for _, line := range strings.Split(message, "\n") {
			if line != "" {
				_, err = templBuffer.WriteString("<p>")
				if err != nil {
					return err
				}
				var var_2 string = line
				_, err = templBuffer.WriteString(templ.EscapeString(var_2))
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("</p>")
				if err != nil {
					return err
				}
			}
		}
		_, err = templBuffer.WriteString("</div></div>")
		if err != nil {
			return err
		}
		if !templIsBuffer {
			_, err = templBuffer.WriteTo(w)
		}
		return err
	})
}

func userMessage(message string) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		templBuffer, templIsBuffer := w.(*bytes.Buffer)
		if !templIsBuffer {
			templBuffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templBuffer)
		}
		ctx = templ.InitializeContext(ctx)
		var_3 := templ.GetChildren(ctx)
		if var_3 == nil {
			var_3 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, err = templBuffer.WriteString("<div hx-swap-oob=\"beforeend:#content\"><div class=\"user-message\">")
		if err != nil {
			return err
		}
		for _, line := range strings.Split(message, "\n") {
			if line != "" {
				_, err = templBuffer.WriteString("<p>")
				if err != nil {
					return err
				}
				var var_4 string = line
				_, err = templBuffer.WriteString(templ.EscapeString(var_4))
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("</p>")
				if err != nil {
					return err
				}
			}
		}
		_, err = templBuffer.WriteString("</div></div>")
		if err != nil {
			return err
		}
		if !templIsBuffer {
			_, err = templBuffer.WriteTo(w)
		}
		return err
	})
}

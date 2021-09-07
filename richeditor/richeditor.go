package richeditor

import (
	"context"
	"fmt"

	"github.com/goplaid/web"
	v "github.com/goplaid/x/vuetify"
	"github.com/jinzhu/gorm"
	"github.com/qor/media/media_library"
	"github.com/qor/qor5/media_library_view"
	h "github.com/theplant/htmlgo"
)

// how to add more plugins from https://imperavi.com/redactor/plugins/
// 1. add {{plugin}}.min.js to redactor dir
// 2. add plugin name in Plugins array

// how to add own plugins
// 1. load plugin jss,css to PluginsJS,PluginsCSS
// 2. add plugin names in Plugins array
var Plugins = []string{"alignment", "table", "video", "imageinsert"}
var PluginsJS [][]byte
var PluginsCSS [][]byte

type RichEditorBuilder struct {
	db          *gorm.DB
	name        string
	value       string
	label       string
	placeholder string
	plugins     []string
	setPlugins  bool
}

func RichEditor(db *gorm.DB, name string) (r *RichEditorBuilder) {
	r = &RichEditorBuilder{db: db, name: name}
	return
}

func (b *RichEditorBuilder) Value(v string) (r *RichEditorBuilder) {
	b.value = v
	return b
}

func (b *RichEditorBuilder) Label(v string) (r *RichEditorBuilder) {
	b.label = v
	return b
}

func (b *RichEditorBuilder) Placeholder(v string) (r *RichEditorBuilder) {
	b.placeholder = v
	return b
}

func (b *RichEditorBuilder) Plugins(v []string) (r *RichEditorBuilder) {
	b.plugins = v
	b.setPlugins = true
	return b
}

func (b *RichEditorBuilder) MarshalHTML(ctx context.Context) ([]byte, error) {
	p := Plugins
	if b.setPlugins {
		p = b.plugins
	}
	r := h.Components(
		v.VSheet(
			h.Label(b.label).Class("v-label theme--light"),
			Redactor().Value(b.value).Placeholder(b.placeholder).Config(RedactorConfig{Plugins: p}).Attr(web.VFieldName(b.name)...),
			h.Div(
				media_library_view.QMediaBox(b.db).FieldName(fmt.Sprintf("%s_richeditor_medialibrary", b.name)).
					Value(&media_library.MediaBox{}).Config(&media_library.MediaBoxConfig{
					AllowType: "image",
				}),
			).Class("hidden-screen-only"),
		).Class("pb-4").Rounded(true).Attr("data-type", "redactor"),
	)
	return r.MarshalHTML(ctx)
}
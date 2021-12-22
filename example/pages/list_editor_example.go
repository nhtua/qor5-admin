package pages

import (
	"fmt"

	"github.com/goplaid/web"
	"github.com/goplaid/x/presets"
	. "github.com/goplaid/x/vuetify"
	"github.com/qor/qor5/example/models"
	"github.com/qor/qor5/listeditor"
	"github.com/qor/qor5/media"
	"github.com/qor/qor5/media/media_library"
	media_view "github.com/qor/qor5/media/views"
	"github.com/qor/qor5/publish"
	h "github.com/theplant/htmlgo"
	"github.com/theplant/testingutils"
	"gorm.io/gorm"
)

type Holder struct {
	Addresses []*models.Address
}

func ListEditorExample(db *gorm.DB, p *presets.Builder) (pf web.PageFunc, sf web.EventFunc) {
	var addressFb = p.NewFieldsBuilder(presets.WRITE).Model(&models.Address{})

	var phoneFb = p.NewFieldsBuilder(presets.WRITE).Model(&models.Phone{})
	phoneFb.Field("Number").ComponentFunc(func(obj interface{}, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		return VTextField().
			FieldName(field.FormValueKey).
			Value(field.StringValue(obj)).
			Label(field.Label)
	})

	phoneFb.Field("ID").ComponentFunc(func(obj interface{}, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		return h.Input("").
			Type("hidden").
			Value(fmt.Sprint(field.Value(obj))).
			Attr(web.VFieldName(field.FormValueKey)...)
	})

	addressFb.Field("ID").ComponentFunc(func(obj interface{}, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		return h.Input("").
			Type("hidden").
			Value(fmt.Sprint(field.Value(obj))).
			Attr(web.VFieldName(field.FormValueKey)...)
	})

	addressFb.Field("Street").ComponentFunc(func(obj interface{}, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		return VTextField().
			FieldName(field.FormValueKey).
			Value(field.Value(obj)).
			Label(field.Label)
	})

	addressFb.Field("Status").ComponentFunc(func(obj interface{}, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		return VSelect().
			Items([]string{"Draft", "PendingReview", "Approved"}).
			Value(field.Value(obj).(publish.Status).Status).
			Label(field.Label).
			FieldName(field.FormValueKey)
	}).SetterFunc(func(obj interface{}, field *presets.FieldContext, ctx *web.EventContext) (err error) {
		ad := obj.(*models.Address)
		ad.Status.Status = ctx.R.FormValue(field.FormValueKey)
		return
	})

	addressFb.ListField("Phones", phoneFb).ComponentFunc(func(obj interface{}, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		return listeditor.New(field).Value(field.Value(obj))
	})

	addressFb.Field("HomeImage").ComponentFunc(func(obj interface{}, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		val := field.Value(obj).(media_library.MediaBox)
		return media_view.QMediaBox(db).
			FieldName(field.FormValueKey).
			Value(&val).
			Config(&media_library.MediaBoxConfig{
				AllowType: "image",
				Sizes: map[string]*media.Size{
					"thumb": {
						Width:  400,
						Height: 300,
					},
					"main": {
						Width:  800,
						Height: 500,
					},
				},
			})
	})

	holderFb := p.NewFieldsBuilder(presets.WRITE).Model(&Holder{})
	holderFb.ListField("Addresses", addressFb).ComponentFunc(func(obj interface{}, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		return listeditor.New(field).Value(field.Value(obj))
	})

	pf = func(ctx *web.EventContext) (r web.PageResponse, err error) {
		r.PageTitle = "List Editor Example"
		holder := &Holder{}

		var addresses []*models.Address
		err = db.Preload("Phones").Find(&addresses).Error
		if err != nil {
			panic(err)
		}

		if len(addresses) == 0 {
			holder.Addresses = []*models.Address{
				{
					Model:  gorm.Model{ID: 1},
					Street: "Street 1",
					Status: publish.Status{
						Status: "Draft",
					},
					Phones: []*models.Phone{
						{
							Number: 11111,
						},
						{
							Number: 22222,
						},
					},
				},
				{
					Model:  gorm.Model{ID: 2},
					Street: "Street 2",
					Status: publish.Status{
						Status: "PendingReview",
					},
					Phones: []*models.Phone{
						{
							Number: 33333,
						},
						{
							Number: 44444,
						},
					}},
				{
					Model:  gorm.Model{ID: 3},
					Street: "Street 3",
					Status: publish.Status{
						Status: "Approved",
					},
					Phones: []*models.Phone{
						{
							Number: 55555,
						},
						{
							Number: 66666,
						},
					}},
			}
		} else {
			holder.Addresses = addresses
		}

		testingutils.PrintlnJson(holder)

		r.Body = VContainer(
			holderFb.ToComponent(nil, holder, ctx),
			VBtn("Save").Attr("@click", web.Plaid().EventFunc("save").Go()),
		)
		return
	}

	sf = func(ctx *web.EventContext) (r web.EventResponse, err error) {
		var holder = &Holder{}
		holderFb.Unmarshal(holder, nil, ctx)
		for _, ad := range holder.Addresses {
			for _, ph := range ad.Phones {
				ph.AddressID = ad.ID
				err = db.Save(ph).Error
				if err != nil {
					panic(err)
				}
			}
			err = db.Save(ad).Error
			if err != nil {
				panic(err)
			}
		}
		r.Reload = true
		return
	}
	return
}

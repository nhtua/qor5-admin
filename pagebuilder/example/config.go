package example

import (
	"os"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/goplaid/web"
	"github.com/qor/oss/s3"
	"github.com/qor/qor5/media"
	"github.com/qor/qor5/media/media_library"
	"github.com/qor/qor5/media/oss"
	"github.com/qor/qor5/pagebuilder"
	h "github.com/theplant/htmlgo"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func ConnectDB() (db *gorm.DB) {
	var err error
	db, err = gorm.Open(postgres.Open(os.Getenv("DB_PARAMS")), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	db.Logger = db.Logger.LogMode(logger.Info)
	return
}

type TextAndImage struct {
	ID    uint
	Text  string
	Image media_library.MediaBox
}

func ConfigPageBuilder(db *gorm.DB) *pagebuilder.Builder {
	sess := session.Must(session.NewSession())

	oss.Storage = s3.New(&s3.Config{
		Bucket:  os.Getenv("S3_Bucket"),
		Region:  os.Getenv("S3_Region"),
		Session: sess,
	})

	media.RegisterCallbacks(db)

	err := db.AutoMigrate(&TextAndImage{})
	if err != nil {
		panic(err)
	}
	pb := pagebuilder.New(db)

	textAndImage := pb.NewContainer("text_and_image").
		ContainerFunc(func(obj interface{}, ctx *web.EventContext) h.HTMLComponent {
			tai := obj.(*TextAndImage)
			return h.Div(
				h.Text(tai.Text),
				h.Img(tai.Image.Url),
			)
		})

	ed := textAndImage.Model(&TextAndImage{}).Editing("Text", "Image")
	ed.Field("Image")
	return pb
}

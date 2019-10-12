package router

import (
	"../app/api/face"
	"../app/api/person"
	"github.com/gogf/gf/frame/g"
)

func init() {
	person_controller := new(person.Controller)
	g.Server().BindObject("Get:/person", person_controller, "List")

	face_controller := new(face.Controller)
	g.Server().BindObject("Post:/face", face_controller, "PostCompareResult")
}

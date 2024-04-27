// ============== /router/router.go
type Router interface {
	Route(r *gin.Engine)
}
type RegisterRouter struct {}

func New()*RegisterRouter {
	return &RegisterRouter{}
}

func(*RegisterRouter) Route (ro Router,r *gin.Engine){
	ro. Route(r)
}

func InitRouter(r *gin.Engine){
	rg := New()
	rg.Route(&user.RouterUser{},r)
}

// ============== /api/user/user.go
package user
import "github.com/gin-gonic/gin"
type HandlerUser struct {}
func (*HandlerUser) getcaptcha (ctx *gin.context){
	ctx.JSON( code: 200, obj:"getCaptcha success" )
}

// ============== /api/user/route.go
package user
import "github.com/gin-gonic/gin"
type RouterUser struct {}
func (*RouterUser) Route (r *gin.Engine){
	h := &HandlerUser{}
	r.PoST(relativePath:"/project/login/getcaptcha", h.getCaptcha)
}
// ============== main.go
import (
	"github.com/gin-gonic/gin"
	srv "test.com/project-common"
	"test.com/project-user/router"
)

func func main() {
	r:= gin.Default()
	router.InitRouter(r)
}
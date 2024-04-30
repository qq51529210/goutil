package gin

import (
	gh "goutil/http"

	"github.com/gin-contrib/pprof"

	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Server struct {
	gh.Server
	// 性能分析的路径，空则不开启
	PProf string
	// swagger 文档的路径，空则不开启
	Swagger string
	// 注册验证器回调
	RegisterValidator func(*gin.Engine) error
	// 注册路由回调
	RegisterRoute func(gin.IRouter)
}

func (s *Server) Serve() error {
	// gin
	g := gin.New()
	// 注册验证
	if s.RegisterValidator != nil {
		if err := s.RegisterValidator(g); err != nil {
			return err
		}
	}
	// 注册路由
	if s.RegisterRoute != nil {
		s.RegisterRoute(g)
	}
	// 性能分析
	if s.PProf != "" {
		pprof.RouteRegister(g.Group(s.PProf), "")
	}
	// 文档
	if s.Swagger != "" {
		g.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		g.NoRoute(func(ctx *gin.Context) {
			ctx.Writer.WriteString(`<a ref="/docs/index.html">document</a>`)
		})
	}
	// 开始服务
	s.Handler = g
	return s.Serve()
}

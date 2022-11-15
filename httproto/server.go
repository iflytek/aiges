package httproto

import (
	"context"
	"embed"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/proto"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/xfyun/aiges/docs"
	"github.com/xfyun/aiges/httproto/internal"
	"github.com/xfyun/aiges/protocol"
	"github.com/xfyun/uuid"
	xsf "github.com/xfyun/xsf/server"
	"github.com/xfyun/xsf/utils"
	"golang.org/x/net/webdav"
	"io/fs"
	"log"
	"net/http"
)

type Server struct {
	si          xsf.UserInterface
	serviceName string
	listenAddr  string
}

func NewServer(rpc xsf.UserInterface) xsf.UserInterface {
	return &Server{
		si:         rpc,
		listenAddr: "",
	}
}

func (s *Server) Init(box *xsf.ToolBox) error {
	s.serviceName = box.Bc.CfgData.Service

	addr, err := box.Cfg.GetString(s.serviceName, "http_listen")
	if err != nil {
		addr = ":"
	}
	s.listenAddr = addr
	go func() {
		err := s.listen()
		if err != nil {
			panic("listen http error:" + err.Error())
		}
	}()
	return s.si.Init(box)
}

var (
	// Handler is used to server files through an http.Handler
	Handler *webdav.Handler

	//go:embed swaggerui
	dist embed.FS

	static fs.FS

	//go:embed test.json
	sampl []byte
)

func init() {
	// Static will store the embedded swagger-UI files for use by the Handler.
	static, _ = fs.Sub(dist, "swaggerui")

	Handler = &webdav.Handler{
		FileSystem: internal.NewWebDAVFileSystemFromFS(static),
		LockSystem: webdav.NewMemLS(),
	}
}
func (s *Server) listen() error {
	// will remove in release
	docs.SwaggerInfo.Title = "Swagger Example API"
	router := gin.Default()
	router.Any("/v1/"+s.serviceName, s.ginHandler())
	router.GET("/test.json", getDemo)
	url := ginSwagger.URL("http://localhost:1888/test.json")
	router.GET("/swagger/*any", ginSwagger.WrapHandler(Handler, url))

	http.ListenAndServe(s.listenAddr, router)
	fmt.Println("[http listen at]: ", s.listenAddr)
	log.Fatal(http.ListenAndServe(s.listenAddr, router))
	//http.Serve(ls, s)
	return nil
}

func getDemo(c *gin.Context) {
	var ret map[string]interface{}
	json.Unmarshal(sampl, &ret)
	c.Header("Access-Control-Allow-Origin", "*") // 可将将 * 替换为指定的域名
	c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
	c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
	c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
	c.Header("Access-Control-Allow-Credentials", "true")
	c.IndentedJSON(http.StatusOK, ret)

}

func (s *Server) ginHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		s.ServeHTTP(c.Writer, c.Request)
	}
}
func (s *Server) Finit() error {
	return s.si.Finit()
}

func (s *Server) Call(req *xsf.Req, span *xsf.Span) (*xsf.Res, error) {

	return s.si.Call(req, span)
}

func (s *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	code, sid, err := s.serveHTTP(writer, request)
	if err != nil {
		writer.WriteHeader(500)
		writeResp(writer, &CommonResponse{
			Header: map[string]interface{}{
				"code":    code,
				"sid":     sid,
				"message": err.Error(),
			},
			Payload: nil,
		})
	}
}

func (s *Server) serveHTTP(writer http.ResponseWriter, request *http.Request) (ret int, sid string, err error) {
	ctx := context.Background()
	sid = generateSID()
	j := json.NewDecoder(request.Body)
	req := new(Request)
	err = j.Decode(req)
	if err != nil {
		return 10000, sid, err
	}
	in, err := req.ConvertToPb(s.serviceName, protocol.LoaderInput_ONCE, &ctx)
	if err != nil {
		return 10001, sid, err
	}

	in.Headers["sid"] = sid
	//in.Expect[0].DataType = protocol.MetaDesc_DataType(protocol.MetaDesc_TEXT)
	bytes, _ := proto.Marshal(in)
	xsfReq := xsf.NewReq()
	xsfReq.Append(bytes, nil)
	xsfReq.SetOp("AIIn")
	xsfReq.SetParam("SeqNo", "1")
	xsfReq.SetParam("version", "v2")
	xsfReq.SetParam("waitTime", "1000")
	xsfReq.SetParam("baseId", "0")
	xsfReq.SetHandle(sid)

	span := utils.NewSpan(utils.SrvSpan)
	span.WithName(s.serviceName)
	span.WithTag("sid", sid)

	res, err := s.si.Call(xsfReq, span)
	if err != nil {
		return 10002, sid, err
	}
	rr := res.Res()
	if rr != nil {
		code := rr.GetCode()
		if code != 0 {
			return int(code), sid, errors.New(rr.GetErrorInfo())
		}
	}
	data := res.GetData()
	if len(data) <= 0 {
		return 10003, sid, fmt.Errorf("output data length is 0")
	}
	output := &protocol.LoaderOutput{}
	err = proto.Unmarshal(data[0].Data, output)
	if err != nil {
		return 10004, sid, fmt.Errorf("output data unmarshal error:%w", err)
	}
	if output.Code != 0 {
		return int(output.Code), sid, errors.New(output.Err)
	}
	writeResp(writer, outputToJson(output, sid, in.Expect))
	return 0, sid, nil
}

type CommonResponse struct {
	Header  map[string]interface{} `json:"header"`
	Payload map[string]interface{} `json:"payload"`
}

func writeResp(w http.ResponseWriter, resp *CommonResponse) {
	w.Header().Set("Content-Type", "application/json")
	j := json.NewEncoder(w)
	j.SetEscapeHTML(false)
	j.Encode(resp)
}

var (
	uid uuid.UUID
)

func init() {
	var err error
	uid, err = uuid.NewV4()
	if err != nil {
		panic(err)
	}
}

func generateSID() string {
	uid, err := uuid.NewV4()
	if err != nil {
		return ""
	}
	return uid.String()

}

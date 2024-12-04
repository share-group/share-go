package constant

type HttpMethod string

const (
	GET     = HttpMethod("GET")
	POST    = HttpMethod("POST")
	HEAD    = HttpMethod("HEAD")
	PUT     = HttpMethod("PUT")
	DELETE  = HttpMethod("DELETE")
	CONNECT = HttpMethod("CONNECT")
	OPTIONS = HttpMethod("OPTIONS")
	TRACE   = HttpMethod("TRACE")
	PATCH   = HttpMethod("PATCH")
)

var httpMethodMap = map[string]HttpMethod{
	"GET":     GET,
	"POST":    POST,
	"HEAD":    HEAD,
	"PUT":     PUT,
	"DELETE":  DELETE,
	"CONNECT": CONNECT,
	"OPTIONS": OPTIONS,
	"TRACE":   TRACE,
	"PATCH":   PATCH,
}

func ValueOfHttpMethod(httpMethod string) HttpMethod {
	return httpMethodMap[httpMethod]
}

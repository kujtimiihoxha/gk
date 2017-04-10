package generator

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/go-errors/errors"
	"github.com/kujtimiihoxha/gk/fs"
	"github.com/kujtimiihoxha/gk/parser"
	"github.com/kujtimiihoxha/gk/templates"
	"github.com/kujtimiihoxha/gk/utils"
	"github.com/spf13/viper"
	"golang.org/x/tools/imports"
	"os"
	"strings"
)

var SUPPORTED_TRANSPORTS = []string{"http"}

type ServiceGenerator struct {
}

func (sg *ServiceGenerator) Generate(name string) error {
	logrus.Info(fmt.Sprintf("Generating service: %s", name))
	f := parser.NewFile()
	f.Package = "service"
	te := template.NewEngine()
	iname, err := te.ExecuteString(viper.GetString("service.interface_name"), map[string]string{
		"ServiceName": name,
	})
	logrus.Debug(fmt.Sprintf("Service interface name : %s", iname))
	if err != nil {
		return err
	}
	f.Interfaces = []parser.Interface{
		parser.NewInterface(iname, []parser.Method{}),
	}
	defaultFs := fs.Get()

	path, err := te.ExecuteString(viper.GetString("service.path"), map[string]string{
		"ServiceName": name,
	})
	logrus.Debug(fmt.Sprintf("Service path: %s", path))
	if err != nil {
		return err
	}
	b, err := defaultFs.Exists(path)
	if err != nil {
		return err
	}
	fname, err := te.ExecuteString(viper.GetString("service.file_name"), map[string]string{
		"ServiceName": name,
	})
	logrus.Debug(fmt.Sprintf("Service file name: %s", fname))
	if err != nil {
		return err
	}
	if b {
		logrus.Debug("Service folder already exists")
		return fs.NewDefaultFs(path).WriteFile(fname, f.String(), false)
	}
	err = defaultFs.MkdirAll(path)
	logrus.Debug(fmt.Sprintf("Creating folder structure : %s", path))
	if err != nil {
		return err
	}
	return fs.NewDefaultFs(path).WriteFile(fname, f.String(), false)
}

func NewServiceGenerator() *ServiceGenerator {
	return &ServiceGenerator{}
}

type ServiceInitGenerator struct {
}

func (sg *ServiceInitGenerator) Generate(name string) error {
	te := template.NewEngine()
	defaultFs := fs.Get()
	path, err := te.ExecuteString(viper.GetString("service.path"), map[string]string{
		"ServiceName": name,
	})
	if err != nil {
		return err
	}
	fname, err := te.ExecuteString(viper.GetString("service.file_name"), map[string]string{
		"ServiceName": name,
	})
	if err != nil {
		return err
	}
	sfile := path + defaultFs.FilePathSeparator() + fname
	b, err := defaultFs.Exists(sfile)
	if err != nil {
		return err
	}
	iname, err := te.ExecuteString(viper.GetString("service.interface_name"), map[string]string{
		"ServiceName": name,
	})
	if err != nil {
		return err
	}
	if !b {
		return errors.New(fmt.Sprintf("Service %s was not found", name))
	}
	transport := viper.GetString("gk_transport")
	supported := false
	for _, v := range SUPPORTED_TRANSPORTS {
		if v == transport {
			supported = true
			break
		}
	}
	if !supported {
		return errors.New(fmt.Sprintf("Transport `%s` not supported", transport))
	}
	p := parser.NewFileParser()
	s, err := defaultFs.ReadFile(sfile)
	if err != nil {
		return err
	}
	f, err := p.Parse([]byte(s))
	if err != nil {
		return err
	}
	var iface *parser.Interface
	for _, v := range f.Interfaces {
		if v.Name == iname {
			iface = &v
		}
	}
	if iface == nil {
		return errors.New(fmt.Sprintf("Could not find the service interface in `%s`", sfile))
	}
	if len(iface.Methods) == 0 {
		return errors.New("The service has no method please implement the interface methods")
	}
	stubName, err := te.ExecuteString(viper.GetString("service.struct_name"), map[string]string{
		"ServiceName": name,
	})
	if err != nil {
		return err
	}
	stub := parser.NewStruct(stubName, []parser.NamedTypeValue{})
	exists := false
	for _, v := range f.Structs {
		if v.Name == stub.Name {
			logrus.Infof("Service `%s` structure already exists so it will not be recreated.", stub.Name)
			exists = true
		}
	}
	if !exists {
		s += "\n" + stub.String()
	}
	exists = false
	for _, v := range f.Methods {
		if v.Name == "New" {
			logrus.Infof("Service `%s` New function already exists so it will not be recreated", stub.Name)
			exists = true
		}
	}
	if !exists {
		newMethod := parser.NewMethod(
			"New",
			parser.NamedTypeValue{},
			fmt.Sprintf(`s = %s{}
			return s`, stub.Name),
			[]parser.NamedTypeValue{},
			[]parser.NamedTypeValue{
				parser.NewNameType("s", stubName),
			},
		)
		s += "\n" + newMethod.String()
	}
	for _, m := range iface.Methods {
		exists = false
		m.Struct = parser.NewNameType(strings.ToLower(iface.Name[:2]), "*"+stub.Name)
		for _, v := range f.Methods {
			if v.Name == m.Name && v.Struct.Type == m.Struct.Type {
				logrus.Infof("Service method `%s` already exists so it will not be recreated.", v.Name)
				exists = true
			}
		}
		if !exists {
			s += "\n" + m.String()
		}
	}
	d, err := imports.Process("g", []byte(s), nil)
	if err != nil {
		return err
	}
	err = defaultFs.WriteFile(sfile, string(d), true)
	if err != nil {
		return err
	}
	err = sg.generateEndpoints(name, iface)
	if err != nil {
		return err
	}
	err = sg.generateTransport(name, iface, transport)
	if err != nil {
		return err
	}
	return nil
}
func (sg *ServiceInitGenerator) generateTransport(name string, iface *parser.Interface, transport string) error {
	switch transport {
	case "http":
		return sg.generateHttpTransport(name, iface)
	default:
		return errors.New(fmt.Sprintf("Transport `%s` not supported", transport))
	}
}
func (sg *ServiceInitGenerator) generateHttpTransport(name string, iface *parser.Interface) error {
	te := template.NewEngine()
	defaultFs := fs.Get()
	handlerFile := parser.NewFile()
	handlerFile.Package = "http"
	gosrc := viper.GetString("GOPATH") + "/src/"
	gosrc = strings.Replace(gosrc, "\\", "/", -1)
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}
	if viper.GetString("gk_folder") != "" {
		pwd += "/" + viper.GetString("gk_folder")
	}
	pwd = strings.Replace(pwd, "\\", "/", -1)
	projectPath := strings.Replace(pwd, gosrc, "", 1)
	enpointsPath, err := te.ExecuteString(viper.GetString("endpoints.path"), map[string]string{
		"ServiceName": name,
	})
	if err != nil {
		return err
	}
	enpointsPath = strings.Replace(enpointsPath, "\\", "/", -1)
	endpointsImport := projectPath + "/" + enpointsPath
	handlerFile.Imports = []parser.NamedTypeValue{
		parser.NewNameType("httptransport", "\"github.com/go-kit/kit/transport/http\""),
		parser.NewNameType("", "\""+endpointsImport+"\""),
	}

	handlerFile.Methods = append(handlerFile.Methods, parser.NewMethod(
		"NewHTTPHandler",
		parser.NamedTypeValue{},
		"m := http.NewServeMux()",
		[]parser.NamedTypeValue{
			parser.NewNameType("endpoints", "endpoints.Endpoints"),
		},
		[]parser.NamedTypeValue{
			parser.NewNameType("", "http.Handler"),
		},
	))
	for _, m := range iface.Methods {
		handlerFile.Methods = append(handlerFile.Methods, parser.NewMethod(
			fmt.Sprintf("Decode%sRequest", m.Name),
			parser.NamedTypeValue{},
			fmt.Sprintf(`req = endpoints.%sRequest{}
			err = json.NewDecoder(r.Body).Decode(&r)
			return req,err
			`, m.Name),
			[]parser.NamedTypeValue{
				parser.NewNameType("_", "context.Context"),
				parser.NewNameType("r", "*http.Request"),
			},
			[]parser.NamedTypeValue{
				parser.NewNameType("req", "interface{}"),
				parser.NewNameType("err", "error"),
			},
		))
		handlerFile.Methods = append(handlerFile.Methods, parser.NewMethod(
			fmt.Sprintf("Encode%sResponse", m.Name),
			parser.NamedTypeValue{},
			` w.Header().Set("Content-Type", "application/json; charset=utf-8")
			err = json.NewEncoder(w).Encode(response)
			json.NewEncoder(w).Encode(response)
			return err
			`,
			[]parser.NamedTypeValue{
				parser.NewNameType("_", "context.Context"),
				parser.NewNameType("w", "http.ResponseWriter"),
				parser.NewNameType("response", "interface{}"),
			},
			[]parser.NamedTypeValue{
				parser.NewNameType("err", "error"),
			},
		))
		handlerFile.Methods[0].Body += "\n" + fmt.Sprintf(`m.Handle("/%s", httptransport.NewServer(
        endpoints.%sEndpoint,
        Decode%sRequest,
        Encode%sResponse,
    ))`, utils.ToLowerSnakeCase(m.Name), m.Name, m.Name, m.Name)
	}
	handlerFile.Methods[0].Body += "\n" + "return m"
	path, err := te.ExecuteString(viper.GetString("transport.path"), map[string]string{
		"ServiceName":   name,
		"TransportType": "http",
	})
	if err != nil {
		return err
	}
	b, err := defaultFs.Exists(path)
	if err != nil {
		return err
	}
	fname, err := te.ExecuteString(viper.GetString("transport.file_name"), map[string]string{
		"ServiceName":   name,
		"TransportType": "http",
	})
	if err != nil {
		return err
	}
	tfile := path + defaultFs.FilePathSeparator() + fname
	if b {
		fex, err := defaultFs.Exists(tfile)
		if err != nil {
			return err
		}
		if fex {
			logrus.Errorf("Transport for service `%s` exist", name)
			logrus.Info("If you are trying to update a service use `gk update service [serviceName]`")
			return nil
		}
	} else {
		err = defaultFs.MkdirAll(path)
		if err != nil {
			return err
		}
	}
	return defaultFs.WriteFile(tfile, handlerFile.String(), false)
}
func (sg *ServiceInitGenerator) generateEndpoints(name string, iface *parser.Interface) error {
	te := template.NewEngine()
	defaultFs := fs.Get()
	enpointsPath, err := te.ExecuteString(viper.GetString("endpoints.path"), map[string]string{
		"ServiceName": name,
	})
	if err != nil {
		return err
	}
	b, err := defaultFs.Exists(enpointsPath)
	if err != nil {
		return err
	}
	endpointsFileName, err := te.ExecuteString(viper.GetString("endpoints.file_name"), map[string]string{
		"ServiceName": name,
	})
	if err != nil {
		return err
	}
	eFile := enpointsPath + defaultFs.FilePathSeparator() + endpointsFileName
	if b {
		fex, err := defaultFs.Exists(eFile)
		if err != nil {
			return err
		}
		if fex {
			logrus.Errorf("Endpoints for service `%s` exist", name)
			logrus.Info("If you are trying to add functions to a service use `gk update service [serviceName]`")
			return nil
		}
	} else {
		err = defaultFs.MkdirAll(enpointsPath)
		if err != nil {
			return err
		}
	}
	file := parser.NewFile()
	file.Package = "endpoints"
	file.Structs = []parser.Struct{
		parser.NewStruct("Endpoints", []parser.NamedTypeValue{}),
	}
	gosrc := viper.GetString("GOPATH") + "/src/"
	gosrc = strings.Replace(gosrc, "\\", "/", -1)
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}
	if viper.GetString("gk_folder") != "" {
		pwd += "/" + viper.GetString("gk_folder")
	}
	pwd = strings.Replace(pwd, "\\", "/", -1)
	projectPath := strings.Replace(pwd, gosrc, "", 1)
	servicePath, err := te.ExecuteString(viper.GetString("service.path"), map[string]string{
		"ServiceName": name,
	})
	if err != nil {
		return err
	}
	servicePath = strings.Replace(servicePath, "\\", "/", -1)
	serviceImport := projectPath + "/" + servicePath
	file.Imports = []parser.NamedTypeValue{
		parser.NewNameType("", "\""+serviceImport+"\""),
	}
	file.Methods = []parser.Method{
		parser.NewMethod(
			"New",
			parser.NamedTypeValue{},
			"",
			[]parser.NamedTypeValue{
				parser.NewNameType("svc", "service."+iface.Name),
			},
			[]parser.NamedTypeValue{
				parser.NewNameType("ep", "Endpoints"),
			},
		),
	}

	for i, v := range iface.Methods {
		file.Structs[0].Vars = append(file.Structs[0].Vars, parser.NewNameType(v.Name+"Endpoint", "endpoint.Endpoint"))
		reqPrams := []parser.NamedTypeValue{}
		for _, p := range v.Parameters {
			if p.Name != "ctx" {
				n := strings.ToUpper(string(p.Name[0])) + p.Name[1:]
				reqPrams = append(reqPrams, parser.NewNameType(n, p.Type))
			}
		}
		resultPrams := []parser.NamedTypeValue{}
		for _, p := range v.Results {
			n := strings.ToUpper(string(p.Name[0])) + p.Name[1:]
			resultPrams = append(resultPrams, parser.NewNameType(n, p.Type))
		}
		req := parser.NewStruct(v.Name+"Request", reqPrams)
		res := parser.NewStruct(v.Name+"Response", resultPrams)
		file.Structs = append(file.Structs, req)
		file.Structs = append(file.Structs, res)
		tmplModel := map[string]interface{}{
			"Calling":  v,
			"Request":  req,
			"Response": res,
		}
		tRes, err := te.ExecuteString("{{template \"endpoint_func\" .}}", tmplModel)
		if err != nil {
			return err
		}
		file.Methods = append(file.Methods, parser.NewMethod(
			"Make"+v.Name+"Endpoint",
			parser.NamedTypeValue{},
			tRes,
			[]parser.NamedTypeValue{
				parser.NewNameType("svc", "service."+iface.Name),
			},
			[]parser.NamedTypeValue{
				parser.NewNameType("ep", "endpoint.Endpoint"),
			},
		))
		file.Methods[0].Body += "\n" + "ep." + file.Structs[0].Vars[i].Name + " = Make" + v.Name + "Endpoint(svc)"
	}
	file.Methods[0].Body += "\n return ep"
	return defaultFs.WriteFile(eFile, file.String(), false)
}
func NewServiceInitGenerator() *ServiceInitGenerator {
	return &ServiceInitGenerator{}
}

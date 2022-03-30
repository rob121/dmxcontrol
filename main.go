package main

import(
	"embed"
	"github.com/oliread/usbdmx"
	"github.com/oliread/usbdmx/ft232"
	"github.com/gorilla/mux"
	"github.com/rob121/vhelp"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"fmt"
    "github.com/spf13/viper"
	"html/template"
	)


//go:embed assets/*
var assetsfs embed.FS
type PageData map[string]interface{}
var conf *viper.Viper
var dms ft232.DMXController
var state map[string]map[string]interface{}

type Device map[string]DeviceEntry

type DeviceEntry struct{
	Note string
	Command map[string]CommandEntry
}

type CommandEntry map[string]string


func main(){

    vhelp.Load("config")
    conf,_ = vhelp.Get("config")


    state = make(map[string]map[string]interface{})

    go setupDMX()
	go startHTTPServer()

	select{}

}
func setupDMX(){

	defer func() {
		if rc := recover(); rc != nil {
			fmt.Println("Recovered in setupDMX", rc)

		}
	}()



	//vid := uint16(0x0403)
	//pid := uint16(0x6001)

	vid := uint16(conf.GetInt("dmx.vid"))
	pid := uint16(conf.GetInt("dmx.pid"))
	debug := conf.GetInt("dmx.debug")
	input_id := conf.GetInt("dmx.input_id")
    output_id := conf.GetInt("dmx.output_id")
    config := usbdmx.NewConfig(vid, pid, output_id,input_id,debug)

	// Get a usb context for our configuration
	config.GetUSBContext()

	// Create a controller and connect to it
	dms = ft232.NewDMXController(config)
	if err := dms.Connect(); err != nil {
		log.Printf("Failed to connect DMX Controller: %s\n", err)
	}

	t := time.NewTicker(30 * time.Millisecond)

	ct := 0
	for range t.C{
		//send something to keep alive?

		renderDmx(ct)

		ct++

		if(ct>40){

			ct = 0
		}
	}

}

func startHTTPServer(){


	defer func() {
		if rc := recover(); rc != nil {
			fmt.Println("Recovered in startHTTPServer", rc)
		}
	}()

	r:= mux.NewRouter()

	port := conf.GetString("http.port")

	r.PathPrefix("/assets/").Handler(http.FileServer(http.FS(assetsfs)))
	r.HandleFunc("/device/{name}/{cmd}", httpCmdHandler)
	r.HandleFunc("/", httpDefaultHandler)
	http.Handle("/", r)

	srv := &http.Server{
		Handler: r,
		Addr:    fmt.Sprintf(":%s", port),
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Println("Listening on port", port)
	log.Fatal(srv.ListenAndServe())


}

func httpCmdHandler(w http.ResponseWriter, r *http.Request) {

	defer func() {
		if rc := recover(); rc != nil {
			fmt.Println("Recovered in httpCmdHandler", rc)
			http.Redirect(w, r, "/", 302)
		}
	}()


	vars:=mux.Vars(r)

	cmd := vars["cmd"]
	name := vars["name"]


	items := conf.GetStringMap(fmt.Sprintf("devices.%s.command.%s",name,cmd))

    if(len(items)>0){
	  state[name] = items
	}



	http.Redirect(w, r, "/", 302)

}

func renderDmx(ct int){


	if(ct == 40) {
		fmt.Println("*************************************************************")
		log.Printf("%#v", state)
	}
		for name,items := range state{

		if(len(items)>0) {
			for k, v := range items {
				if (ct == 40) {

					log.Printf("Sending %s to channel %s for %s\n", v, k, name)

				}
				c, ce := strconv.Atoi(k)
				bv, be := strconv.Atoi(fmt.Sprint(v))
				if (ce == nil && be == nil) {
					dms.SetChannel(int16(c), byte(bv)) //just send the last channel
				}
			}
		}

   }

	dms.Render()


}


func httpDefaultHandler(w http.ResponseWriter, r *http.Request) {

	tmpl := loadTmpl("assets/tmpl/index.html")

	data := PageData{}

	var all Device
    var devs = make(Device)
    var scenes = make(Device)

	conf.UnmarshalKey("devices",&all)

	for k,v := range all {

		if(strings.Contains(k,"scene")){

			scenes[k]=v

		}else{

			devs[k]=v

		}

	}



	data["Title"] = "Home"
	data["Devices"] = &devs
	data["Scenes"] = &scenes

	out := tmpl.Execute(w, data)

    log.Println(out)

}

func loadTmpl(dest string) (*template.Template){

	tmpl := template.New("default")
	tmpl, _ = tmpl.New("header").ParseFS(assetsfs,"assets/tmpl/header.html")
	tmpl, _ = tmpl.New("footer").ParseFS(assetsfs,"assets/tmpl/footer.html")

	var err error

	tmpl, err = tmpl.New(filepath.Base(dest)).Funcs(template.FuncMap{
		"getTimestamp": func() int64 {
			return time.Now().Unix()
		},
		"title": func(str string) (string){

			return strings.Title(strings.Replace(str,"_"," ",-1))

		},
		"colorcmd": func(cmd string) (string){

			switch cmd {
			case "on":
				return "btn-primary"
			case "off":
				return "btn-danger"
			default:
				return "btn-secondary"
			}
		},
	}).ParseFS(assetsfs,dest)

	if(err!=nil){

		log.Println(err)
	}

	return tmpl
}



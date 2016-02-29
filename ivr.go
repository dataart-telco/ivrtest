package main

import (
	common "github.com/dataart-telco/apps-demo/common"
	"fmt"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"strconv"
	"net"
	"encoding/json"
)

type Stat struct {
	Incoming int
	Received int
}

type Ivr struct{
	Host string
	Port int

	Number string

	RestcommApi *common.RestcommApi

	incoming chan int `json:"-"`
	gather chan int `json:"-"`

	Stat *Stat
}


func getLocalIp() net.IP {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP
			}
		}
	}
	return nil
}

func createCtrlCChan() chan os.Signal {
	var signalChannel chan os.Signal
	signalChannel = make(chan os.Signal, 2)

	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
	return signalChannel
}

func (self *Ivr) Json() string{
	data, _ := json.Marshal(self)
	return string(data)
}

func (self *Ivr) Listen(){
	http.HandleFunc("/start", func(w http.ResponseWriter, r *http.Request){
		w.Header().Set("Content-Type", "text/html")
		self.Stat = &Stat{}
		fmt.Fprintf(w,"Stat is reseted");
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, "%s", self.Json());
	})
	http.HandleFunc("/stat/incoming", func(w http.ResponseWriter, r *http.Request){
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, "%d", self.Stat.Incoming);
	})
	http.HandleFunc("/stat/received", func(w http.ResponseWriter, r *http.Request){
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, "%d", self.Stat.Received);
	})
	http.HandleFunc("/incoming", self.handlerIncoming)
	http.HandleFunc("/gahter", self.handlerGather)
	err := http.ListenAndServe(fmt.Sprintf(":%d", self.Port), nil)
	if err != nil {
		panic(err)
	}
	signalChannel := createCtrlCChan()
	for{
		select {
			case <- self.incoming:
				self.Stat.Incoming ++
			case <- self.gather:
				self.Stat.Received ++
			case <- signalChannel:
				return;
		}
	}
}

func (self *Ivr) handlerIncoming(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/xml")
	//from := r.PostFormValue("From")
	fmt.Fprintf(w,
		"<Response><Gather action=\"%s\" method=\"POST\" numDigits=\"1\"><Play>%s</Play></Gather><Hangup/></Response>",
		self.GetUrl("gahter"),
		self.GetUrl("message.wav"))
	self.incoming <- 1
}

func (self *Ivr) handlerGather(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/xml")
	fmt.Fprintf(w,
		"<Response><Play>%s</Play><Hangup/></Response>",
		self.GetUrl("confirm.wav"))
	val, _ := strconv.Atoi(r.PostFormValue("Digits"))
	self.gather <- val
}

func (self *Ivr) RegisterNumber(){
	common.Info.Println("\tRegister number:", self.Number)
	common.NewIncomingPhoneNumber("", self.Number).CreateOrUpdate(*self.RestcommApi, self.GetUrl("incoming"))
}

func (self *Ivr) GetUrl(path string) string {
	return fmt.Sprintf("http://%s:%d/%s", self.Host, self.Port, path)
}

func main(){
	number := flag.String("n", "7777", "Test number")
	host := flag.String("h", getLocalIp().String(), "Host")
	port := flag.Int("p", 7090, "Port")
	rHost := flag.String("r", "127.0.0.1:8080", "Restcomm host")
	rUser := flag.String("r-user", "ACae6e420f425248d6a26948c17a9e2acf", "Restcomm user")
	rPswd := flag.String("r-pswd", "42d8aa7cde9c78c4757862d84620c335", "Restcomm password")
	flag.Parse()

	common.InitLog(os.Stdout, os.Stdout, os.Stdout, os.Stdout)
	api := common.NewRestcommApi(*rHost, *rUser, *rPswd)
	ivr := &Ivr{Host: *host, Port: *port, Number: *number, RestcommApi: &api,
		incoming: make(chan int, 200), gather: make(chan int, 200),
		Stat: &Stat{}}
	common.Trace.Println("Started with", ivr.Json())
	ivr.RegisterNumber()
	ivr.Listen()
}

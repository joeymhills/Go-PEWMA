package goggins

import (
    "fmt"
    "net/http"
    "os"
    "io"
    "sync"
    "time"

    "gopkg.in/yaml.v3"
)

type LoadBalancer struct {
    Name string `yaml:"name"`
    Algorithm string `yaml:"algorithm"`
    Port int `yaml:"port"`
    Servers []Server`yaml:"servers"`
    HealthCheck HealthCheck `yaml:"heatlh_check"`
}
type Server struct {
    Name string `yaml:"name"`
    Address string `yaml:"address"`
    Port int `yaml:"port"`
    Weight float64 `yaml:"weight"`
    Conns int `yaml:"conns"`
    mutex sync.Mutex
}
type HealthCheck struct {
    Interval string `yaml:"interval"`
    Timeout string `yaml:"timeout"`
    Retries string `yaml:"retries"`
    Path string `yaml:"path"`
    Method string `yaml:"method"`
}

func handleRequest(lb *LoadBalancer) func(w http.ResponseWriter, r *http.Request){
    return func(w http.ResponseWriter, r *http.Request){

    server := lb.selectServer()
    addr := fmt.Sprintf(server.Address + fmt.Sprintf("%d", server.Port))

    start := time.Now()
    resp, err := http.Post(addr, "plain/text", r.Body)
    if err != nil {
	fmt.Print("error in creating request: ", err)
    }
    elapsed := time.Since(start).Seconds()
    go server.updateWeight(elapsed)
    
    respBytes, err := io.ReadAll(resp.Body)
    if err != nil {
	fmt.Print("error parsing body: ", err)
    }

    w.Write(respBytes)
    }
}

// Updates server weight using Peak EWMA algorithm detailed here:
// https://linkerd.io/2016/03/16/beyond-round-robin-load-balancing-for-latency/
func(s *Server) updateWeight(responseTime float64) {
    s.mutex.Lock()
    defer s.mutex.Unlock()

    a := 0.6
    s.Weight = a * responseTime + (1 - a) * s.Weight
}

//TODO: Implement a min heap in place of iterating over a for loop
func(lb *LoadBalancer) selectServer() *Server {
    fastestServ := &lb.Servers[0]
    for i := 0; i < len(lb.Servers); i++ {
	if (lb.Servers[i].Weight * float64(lb.Servers[i].Conns)) > (fastestServ.Weight * float64(fastestServ.Conns)) {
	    fastestServ = &lb.Servers[i] 	
	}
    }

    return fastestServ
}

// Parses config.yaml file and returns a handle to load balancer 
func Init(filepath string) (*LoadBalancer, error) {
    // Reads load balancer config from yaml file
    data, err := os.ReadFile(filepath)
    if err != nil {
	return nil, err
    }
    
    //Unmarshalls yaml file into a load balancer struct
    var lb LoadBalancer
    err = yaml.Unmarshal(data, &lb)
    if err != nil {
	return nil, err
    }

    return &lb, nil
}

func(lb *LoadBalancer) StartServer() error{
    // Creates HTTP Server
    serv := http.NewServeMux()

    serv.HandleFunc("/", handleRequest(lb))
    addr := fmt.Sprintf("0.0.0.0:" + fmt.Sprintf("%d", lb.Port))
    
    err := http.ListenAndServe(addr, serv)
    if err != nil {
	fmt.Println("error starting server")
	return err
    }
    return nil
}

package goggins

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"gopkg.in/yaml.v3"
)

type LoadBalancer struct {
    Name string `yaml:"name"`
    Algorithm string `yaml:"algorithm"`
    Servers []Server`yaml:"servers"`
    HealthCheck HealthCheck `yaml:"heatlh_check"`
}
type Server struct {
    Name string `yaml:"name"`
    Address string `yaml:"address"`
    Port int `yaml:"port"`
    Weight float64 `yaml:"weight"`
    Conns int `yaml:"conns"`
}
type HealthCheck struct {
	Interval string `yaml:"interval"`
	Timeout string `yaml:"timeout"`
	Retries string `yaml:"retries"`
	Path string `yaml:"path"`
	Method string `yaml:"method"`
}

func handleRequest(w http.ResponseWriter, r *http.Request,) {
    body, err := io.ReadAll(r.Body)
    if err != nil {
        fmt.Println("Error reading request body")
    }

    fmt.Println(string(body))
}

// Implements Peak EWMA algorithm detailed here:
// https://linkerd.io/2016/03/16/beyond-round-robin-load-balancing-for-latency/
func(s *Server) updateWeight(responseTime float64) float64 {
    a := 0.6
    s.Weight = a * responseTime + (1 - a) * s.Weight
    return s.Weight * float64(s.Conns)
}

// Used to initialize load balancer struct from yaml file
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

func(*LoadBalancer) StartServer() {
    // Creates HTTP Server
    serv := http.NewServeMux()
    serv.HandleFunc("/", handleRequest)

    err := http.ListenAndServe("0.0.0.0:8888", serv)
    if err != nil {
	fmt.Println("error starting server")
    }
}

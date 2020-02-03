package main

import (
    "strings"
    "log"
    "os"
    "os/exec"
    "net"
    "net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
    remote_host := strings.Split(r.RemoteAddr, ":")[0]
    service_port := os.Getenv("SERVICE_PORT")
    if net.ParseIP(remote_host) == nil {
        http.Error(w, "bad remote", http.StatusForbidden)
    } else {
        local_addr := os.Getenv("LOCAL_SERVICE_ADDRESS")
        exec_command := exec.Command("sudo", "ufw", "allow", "to", local_addr, "port", service_port, "from", remote_host)

        err := exec_command.Run()
        if err != nil {
            http.Error(w, "forbidden", http.StatusForbidden)
        }
	// Otherwise do nothing, command ran successfully
    }
}

func main() {
    route := os.Getenv("HANDLER_ROUTE")
    if len(route) < 12 {
        log.Fatal("Route length not long enough")
        os.Exit(2)
    }
    route = "/" + route
    http.HandleFunc(route, handler)
    listen_spec := os.Getenv("LISTEN_SPEC")
    log.Fatal(http.ListenAndServe(listen_spec, nil))
}

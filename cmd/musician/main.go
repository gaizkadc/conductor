//
// Copyright (C) 2018 Nalej Group - All Rights Reserved
//

package main

/*
import (
    "github.com/golang/glog"
    "flag"
)

func main(){
    flag.Parse()
    glog.Info("This is an example")
    glog.Error("This is an error")
    glog.Warning("This is a warning")
    glog.V(1).Info("Info at level 1")
}
*/


import (
    "fmt"
    "github.com/nalej/service"
    "github.com/nalej/conductor/pkg/musician/statuscollector"
)

func main(){
    /*
    log.WithFields(log.Fields{
        "animal": "walrus",
    }).Info("A walrus appears")

    log.SetFormatter(&logmatic.JSONFormatter{})
    // log an event as usual with logrus
    log.WithFields(log.Fields{"string": "foo", "int": 1, "float": 1.1 }).Info("My first ssl event from golang")
    */

    fmt.Println("Launching status collector")

    //collector := statuscollector.NewPrometheusStatusCollector("http://prometheus:9090",10000)
    collector := statuscollector.NewPrometheusStatusCollector("http://192.168.99.100:31080",10000)
    collectorService := statuscollector.Service{collector}
    service.Launch(&collectorService)
}
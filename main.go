package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"k8s.io/apimachinery/pkg/api/resource"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	metrics "k8s.io/metrics/pkg/client/clientset/versioned"
)

type Microservice struct {
	Name       string
	Containers []*Metrics
}

type Metrics struct {
	Name   string
	CPU    resource.Quantity
	Memory resource.Quantity
}

func main() {
	fmt.Printf("start\n")
	config, err := rest.InClusterConfig()
	if err != nil {
		fmt.Println(err)
		return
	}

	client, err := metrics.NewForConfig(config)
	if err != nil {
		fmt.Println(err)
		return
	}

	pods, err := client.MetricsV1beta1().PodMetricses("").List(context.TODO(), metaV1.ListOptions{})
	if err != nil {
		fmt.Println(err)
		return
	}

	microservices := []*Microservice{}
	for _, pod := range pods.Items {
		containers := []*Metrics{}
		microserviceName := pod.ObjectMeta.Name
		for _, container := range pod.Containers {
			containerName := container.Name
			cpu := container.Usage["cpu"]
			memory := container.Usage["memory"]
			containers = append(containers, &Metrics{Name: containerName, CPU: cpu, Memory: memory})
		}
		microservices = append(microservices, &Microservice{Name: microserviceName, Containers: containers})
	}

	for _, microservice := range microservices {
		fmt.Printf("%#v\n", microservice)
	}

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGTERM)
	<-signalCh
	fmt.Printf("close\n")
}

package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	metrics "k8s.io/metrics/pkg/client/clientset/versioned"
)

func collectMetrics(client *metrics.Clientset, now time.Time, monitorWindow map[string]*MetricsWindow) {
	pods, err := client.MetricsV1beta1().PodMetricses("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Println(err)
		return
	}

	microservices := []*Microservice{}
	for _, pod := range pods.Items {
		containers := []*Container{}
		microserviceName := pod.ObjectMeta.Name
		for _, container := range pod.Containers {
			containerName := container.Name
			cpu := container.Usage["cpu"]
			memory := container.Usage["memory"]
			containers = append(containers, NewContainer(containerName, &Metrics{CPU: cpu.MilliValue(), Memory: memory.Value()}))
		}
		microservices = append(microservices, NewMicroservice(microserviceName, containers, now))
	}

	for _, microservice := range microservices {
		mw, ok := monitorWindow[microservice.Name]
		if !ok {
			m := NewMetricsWindow(5, 10, 50*MiB)
			m.Push(*microservice.Metrics)
			monitorWindow[microservice.Name] = m
		} else {
			mw.Push(*microservice.Metrics)
			floatMetrics, err := mw.Alert()
			if err != nil {
				fmt.Printf("%v\n", err)
				averageMetrics := NewAverageMetrics(microservice.Name, now, floatMetrics)
				js, err := averageMetrics.ToJSONString()
				if err != nil {
					fmt.Printf("%v\n", err)
				}
				fmt.Printf("%s\n", js)
			}
		}
		jsonstring, err := microservice.ToJSONString()
		if err != nil {
			fmt.Printf("%v\n", err)
		}
		fmt.Printf("%s\n", jsonstring)
	}
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

	monitorWindow := map[string]*MetricsWindow{}

	signalCh := make(chan os.Signal, 1)
	defer close(signalCh)
	signal.Notify(signalCh, syscall.SIGTERM)

	t := time.NewTicker(5 * time.Second)
	defer t.Stop()

	for {
		select {
		case <-signalCh:
			goto END
		case now := <-t.C:
			collectMetrics(client, now, monitorWindow)
		}
	}
END:
	fmt.Printf("close\n")
}

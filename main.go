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

func collectMetrics(client *metrics.Clientset, now time.Time, alerts []*Alert, env Env, alertSetting AlertSetting) {
	pods, err := client.MetricsV1beta1().PodMetricses("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Println(err)
		return
	}

	microservices := []Microservice{}
	for _, pod := range pods.Items {
		containers := []Container{}
		microserviceName := pod.ObjectMeta.Name
		for _, container := range pod.Containers {
			containerName := container.Name
			cpu := container.Usage["cpu"]
			memory := container.Usage["memory"]
			containers = append(containers, NewContainer(containerName, Metrics{CPU: cpu.MilliValue(), Memory: memory.Value()}))
		}
		microservices = append(microservices, NewMicroservice(microserviceName, containers, now))
	}

	for _, microservice := range microservices {

		if a, ok := FindAlert(alerts, microservice.Name); ok {
			a.Window.Push(microservice.Metrics)
			var jt jsonTime
			jt.Time = now
			floatMetrics, err := a.Alert(jt, env, alertSetting)
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

	env, err := GetEnv()
	if err != nil {
		return
	}

	alertSetting, err := LoadAlertSetting(env)
	if err != nil {
		return
	}

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

	alert := MakeAlert(alertSetting, env)

	signalCh := make(chan os.Signal, 1)
	defer close(signalCh)
	signal.Notify(signalCh, syscall.SIGTERM)

	t := time.NewTicker(time.Second * time.Duration(env.Interval))
	defer t.Stop()

	for {
		select {
		case <-signalCh:
			goto END
		case now := <-t.C:
			collectMetrics(client, now, alert, env, alertSetting)
		}
	}
END:
	fmt.Printf("close\n")
}

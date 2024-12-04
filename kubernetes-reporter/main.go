/*
Copyright © 2024 Case Wylie <casewylie@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

const logFilePath = "/tmp/k8s-watcher.log"

type LogEvent struct {
	Resource     string `json:"resource"`
	Timestamp    string `json:"timestamp"`
	TypeOfChange string `json:"typeOfChange"`
}

// logToFile logs an event to a file.
func logToFile(resource string, eventType string) {
	entry := LogEvent{
		Resource:     resource,
		Timestamp:    time.Now().Format(time.RFC3339),
		TypeOfChange: eventType,
	}

	logString := fmt.Sprintf("%s - Resource: %s, Type of Change: %s\n", entry.Timestamp, entry.Resource, entry.TypeOfChange)
	fmt.Print(logString) // Log to stdout for visibility.

	file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Failed to open log file: %v\n", err)
		return
	}
	defer file.Close()

	if _, err := file.WriteString(logString); err != nil {
		log.Printf("Failed to write to log file: %v\n", err)
	}
}

// watchResource starts a watch on the specified resource.
func watchResource(ctx context.Context, client kubernetes.Interface, resourceType string) {
	var watcher watch.Interface
	var err error

	switch resourceType {
	case "Service":
		watcher, err = client.CoreV1().Services("default").Watch(ctx, metav1.ListOptions{
			FieldSelector: "metadata.name=kubernetes",
		})
	case "EndpointSlice":
		watcher, err = client.DiscoveryV1().EndpointSlices("default").Watch(ctx, metav1.ListOptions{
			FieldSelector: "metadata.name=kubernetes",
		})
	default:
		log.Fatalf("Unsupported resource type: %s", resourceType)
	}

	if err != nil {
		log.Fatalf("Failed to watch %s: %v", resourceType, err)
	}
	defer watcher.Stop()

	for event := range watcher.ResultChan() {
		switch event.Type {
		case watch.Added, watch.Modified:
			logToFile(resourceType, "CREATED_OR_UPDATED")
		case watch.Deleted:
			logToFile(resourceType, "DELETED")
		}
	}
}

func main() {

	config, err := rest.InClusterConfig()
	if err != nil {
		log.Fatalf("Failed to load in-cluster configuration: %v", err)
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Failed to create Kubernetes client: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go watchResource(ctx, client, "Service")
	go watchResource(ctx, client, "EndpointSlice")

	select {}
}

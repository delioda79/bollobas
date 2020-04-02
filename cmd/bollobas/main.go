
package main

import (
	"fmt"
	"os"

	"context"
	"github.com/beatlabs/patron"
	"github.com/beatlabs/patron/log"

	"time"

	"github.com/beatlabs/patron/component/http"

	"github.com/beatlabs/patron/component/async"

	"github.com/beatlabs/patron/component/async/kafka"
	"github.com/beatlabs/patron/component/async/kafka/group"
	"github.com/beatlabs/patron/encoding/json"
)

var (
	version = "dev"
    
)

func main() {
	ctx := context.Background()
	name := "bollobas"

	// Setupbollobas config
	//cfg, err := config.New(ctx)
	//if err != nil {
	//	fmt.Printf("failed to set up configuration: %v", err)
	//	os.Exit(1)
	//}
	//fmt.Printf("Test configuration value: %s\n", cfg.TestVar.Get())
    
	err := patron.SetupLogging(name, version)
	if err != nil {
		fmt.Printf("failed to set up logging: %v", err)
		os.Exit(1)
	}

	
	// Setup HTTP route builder with a singe GET route
	routesBuilder := http.NewRoutesBuilder().Append(http.NewRouteBuilder("/", func(ctx context.Context, req *http.Request) (*http.Response, error) {
		return http.NewResponse("Get data"), nil
	}).MethodGet())
	

    
	// Setup Kafka consumer component
	kafkaCf, err := group.New(name, "todoKafkaGroup", []string{"todoKafkaTopic"}, []string{"localhost:29092"}, kafka.Decoder(json.DecodeRaw))
	if err != nil {
		log.Fatalf("failed to create kafka component factory: %v", err)
	}

	kafkaCmp, err := async.New("kafka-cmp", kafkaCf, func(msg async.Message) error {
		// Implement process function
		return nil
	}).
		WithRetries(10).
		WithRetryWait(5 * time.Second).
		Create()

	if err != nil {
		log.Fatalf("failed to create Kafka component: %v", err)
	}
    

    

	sig := func() {
		fmt.Println("exit gracefully...")
		os.Exit(0)
	}

	
	// Append components
	ccmp := make([]patron.Component, 0)
	ccmp = append(ccmp, kafkaCmp)
	
	

	err = patron.New(name, version).
		WithSIGHUP(sig).
	
		WithRoutesBuilder(routesBuilder).
	
	
		WithComponents(ccmp...).
	
		Run(ctx)

	if err != nil {
		log.Fatalf("failed to create service %v", err)
	}
    
}
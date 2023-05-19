/*
	service.go
	Purpose: Define the lifecycles of a service.

	@author Evan Chen
	@version 1.0 2023/02/22
*/

package service

import (
	"context"
	"sync"
	"time"

	"app/core/property"

	"golang.org/x/exp/slog"
)

// Service is an interface to discribe the life-cycles of a service
type Service interface {
	// Init is the procedure that would be called only once before the server starts.
	Init() error
	// Load is the procedure that would be called before the server starts,
	// and whenever a service reload has been triggered.
	//
	// Note that Load could be called at any time
	// and service sould not crash or report error during and after the process.
	// Any issues during the process could log meaningful message for developers to trace.
	Load()
	// Del is the procedure that would be called before the application shuts down
	// and is responsible for services to cleanup all resources it uses
	// or to save its status to survive application restarts.
	//
	// Note that
	//   * Del could be called in parallel, so be sure to not to depend on other services.
	//   * There would usually be a timeout for app to shutdown.
	Del()
}

var services []Service

func ServiceList() []Service {
	return services
}

// Register registers a Service that would follow the service life-cycles.
func Register(server Service) {
	services = append(services, server)
}

// Initiate initiates the registered services. It will return and stop at the first error it encounters.
//
// Init will be called according to the regestered order
// and is assumed that services don't have dependency relationships with each other.
// If there were, you have to register services accordingly.
func Initiate() error {
	slog.Info("service initiating", slog.String("act", "Initiate"), slog.String("mod", "lifecycle"))
	property.SetState(property.STATE_INIT)
	for i := range services {
		if err := services[i].Init(); err != nil {
			return err
		}
	}
	return nil
}

func Load() {
	slog.Info("service loading", slog.String("act", "Load"), slog.String("mod", "lifecycle"))
	property.SetState(property.STATE_LOAD)
	for i := range services {
		services[i].Load()
	}
}

func Del(timeout time.Duration) context.Context {
	slog.Info("service dropping", slog.String("act", "Del"), slog.String("mod", "lifecycle"))
	ctx, done := context.WithTimeout(context.Background(), timeout)
	go func() {
		wait := &sync.WaitGroup{}
		for i := range services {
			wait.Add(1)
			go func(index int) {
				services[index].Del()
				wait.Done()
			}(i)
		}
		wait.Wait()
		done()
	}()
	return ctx
}

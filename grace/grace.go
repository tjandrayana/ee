// Package grace provide simple graceful shutdown capability to the application server.
//
// It is fully compatible with socketmaster, ctrl-c, and kubernetes.
package grace

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/tjandrayana/ee/socketmaster/child"
	"google.golang.org/grpc"
)

// ServeGRPC start the grpc server on the given address and add graceful shutdown handler
func ServeGRPC(server *grpc.Server, address string) error {
	lis, err := Listen(address)
	if err != nil {
		return err
	}

	stoppedCh := WaitTermSig(func(ctx context.Context) error {
		server.GracefulStop()
		return nil
	})

	log.Printf("GRPC server running on adress %v", address)

	if err := server.Serve(lis); err != nil {
		// Error starting or closing listener:
		return err
	}

	<-stoppedCh
	return nil
}

// ServeHTTP start the http server on the given address and add graceful shutdown handler
func ServeHTTP(srv *http.Server, address string) error {
	// start graceful listener
	lis, err := Listen(address)
	if err != nil {
		return err
	}

	stoppedCh := WaitTermSig(srv.Shutdown)

	log.Printf("http server running on address: %v", address)

	// start serving
	if err := srv.Serve(lis); err != http.ErrServerClosed {
		return err
	}

	<-stoppedCh
	log.Println("HTTP server stopped")
	return nil
}

// WaitTermSig wait for termination signal and then execute the given handler
// when the signal received
//
// The handler is usually service shutdown, so we can properly shutdown our server upon SIGTERM.
//
// It returns channel which will be closed after the signal received and the handler executed.
// We can use the signal to wait for the shutdown to be finished.
func WaitTermSig(handler func(context.Context) error) <-chan struct{} {
	stoppedCh := make(chan struct{})
	go func() {
		signals := make(chan os.Signal, 1)

		// wait for the sigterm
		signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
		<-signals

		// We received an os signal, shut down.
		if err := handler(context.Background()); err != nil {
			log.Printf("graceful shutdown  failed: %v", err)
		} else {
			log.Println("gracefull shutdown succeed")
		}
		close(stoppedCh)
	}()
	return stoppedCh
}

// Listen listens to the given port or to file descriptor as specified by socketmaster.
//
// This method is taken from tjandrayana/grace repo and  modified to work with
// socketmaster's -wait-child-notif option.
func Listen(port string) (net.Listener, error) {
	var l net.Listener

	// see if we run under socketmaster
	fd := os.Getenv("EINHORN_FDS")
	if fd != "" {
		sock, err := strconv.Atoi(fd)
		if err != nil {
			return nil, err
		}
		log.Println("detected socketmaster, listening on", fd)
		file := os.NewFile(uintptr(sock), "listener")
		fl, err := net.FileListener(file)
		if err != nil {
			return nil, err
		}
		l = fl
	}

	if l != nil { // we already have the listener, which listen on EINHORN_FDS
		notifSocketMaster()
		return l, nil
	}

	// we are not using socketmaster, no need to notify

	return net.Listen("tcp4", port)
}

// notifSocketMaster notify socket master about our readyness
// we should remove this func after we are fully moved to tjandrayana new platform
func notifSocketMaster() {
	go func() {
		err := child.NotifyMaster()
		if err != nil {
			log.Printf("failed to notify socketmaster: %v, ignore if you don't use `wait-child-notif` option", err)
		} else {
			log.Println("successfully notify socketmaster")
		}
	}()
}

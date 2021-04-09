package main

import (
	"poptimizer/data/app"
)

func main() {
	app.StartLogging()
	defer app.ShutdownLogging()

	svr := app.RunServer()
	defer app.StopServer(svr)

	<-app.TerminationSignal()

	//q := app.App{}
	//
	//wg := sync.WaitGroup{}
	//wg.Add(1)
	//go func() {
	//	defer wg.Done()
	//
	//	q.Run(ctx)
	//}()
	//

	//
	//
	//wg.Wait()

}

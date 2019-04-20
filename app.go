package main

import (
	"app/config"
	"app/src"
	"app/src/StreetViewImage/Infrastructure/Logger"
	"app/src/StreetViewImage/Infrastructure/Server"
	"fmt"
	"github.com/google/uuid"
	"github.com/j7mbo/goenvconfig"
	"github.com/j7mbo/goij"
	"github.com/j7mbo/goij/src/TypeRegistry"
	"github.com/sirupsen/logrus"
	"io/ioutil"
)

/*
configToShareWithInjector shares any objects there with the injector in shareConfiguration().

Ensure these are pointers as the call to the goenvconfig lib parser requires pointers to update them with env vars.

[...] instead of [] ensures we get a fixed-size array instead of a slice.
*/
var configToShareWithInjector = [...]interface{}{
	&config.ElasticSearchConfiguration{},
	&config.GrpcServerConfiguration{},
	&config.RedisConfiguration{},
	&config.StreetViewApiConfiguration{},
}

/* Here we golang! */
func main() {
	ij := Goij.NewInjector(TypeRegistry.New(src.GetRegistry(), src.GetConfigRegistry(), src.GetVendorRegistry()), nil)

	/* Injector Configuration. */
	shareConfiguration(ij)
	shareInjector(ij)
	configureLogger(ij)
	delegateGrpcMapper(ij)

	/* Webserver (for GRPC actually). */
	ij.Make("app/src/StreetViewImage/Infrastructure/Server.GrpcServer").(Server.GrpcServer).Run()
}

/* shareInjector shares Goij in case a factory needs access to the injector. */
func shareInjector(injector Goij.Injector) {
	injector.Delegate("github.com/j7mbo/goij.Injector", func() Goij.Injector {
		return injector
	})
}

/* shareConfiguration shares env-parsed configuration objects for automatic injection. */
func shareConfiguration(injector Goij.Injector) {
	parser := goenvconfig.NewGoEnvParser()

	for _, conf := range configToShareWithInjector {
		_ = parser.Parse(conf)

		injector.Share(conf)
	}

}

/* configureLogger attempts the first application logger and shares any required fields. */
func configureLogger(injector Goij.Injector) {
	requiredFields := Logger.RequiredLogFields{Env: "Dev", CorrelationId: uuid.New()}

	/* The UUID here will be overwritten with one from the request. We can't have a UUID here. */
	injector.Share(requiredFields)

	elasticLogger := injector.Make(
		"app/src/StreetViewImage/Infrastructure/Logger.ElasticSearchLogger",
	).(*Logger.ElasticSearchLogger)

	injector.Share(elasticLogger)

	backupLogFile, err := ioutil.TempFile("/tmp", "palmago_streetview_log_")

	if err != nil {
		panic(err)
	}

	fmt.Println("Backup log file located at: " + backupLogFile.Name())

	fileLogger, err := Logger.NewFileLogger(*backupLogFile, *logrus.New(), requiredFields)

	if err != nil {
		panic(err)
	}

	/*
		Atm, we have to do this instead of Share() because Share() is not the first-used choice for the injector
		(maybe it should be?)...
	*/
	injector.Delegate("app/src/StreetViewImage/Infrastructure/Logger.LoggingStrategy", func() *Logger.LoggingStrategy {
		return Logger.NewLoggingStrategy(*elasticLogger, *fileLogger)
	})
}

/*
delegateGrpcMapper delegates the initialisation of the controller interface to the Server.GerpcErrorMapper - this is
used to avoid cyclic dependencies from app -> controller -> server -> controller. Fuck sake.
*/
func delegateGrpcMapper(injector Goij.Injector) {
	injector.Delegate("app/src/StreetViewImage/Presentation/Controller.GrpcErrorMapper", func() Server.GrpcErrorMapper {
		return Server.NewGrpcErrorMapper(
			injector.Make("app/src/StreetViewImage/Infrastructure/Logger.LoggingStrategy").(*Logger.LoggingStrategy),
		)
	})
}

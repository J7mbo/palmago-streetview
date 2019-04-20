package src

import "github.com/j7mbo/goij/src/TypeRegistry"
import mGQzNMon "app/src/StreetViewImage/Application/Error"
import poXJtEkr "app/src/StreetViewImage/Application/Query"
import mKaXayJi "app/src/StreetViewImage/Application/QueryHandler"
import GyZJpPBm "app/src/StreetViewImage/Domain"
import olJUMOFZ "app/src/StreetViewImage/Infrastructure/ApiClient"
import DpzQhmiZ "app/src/StreetViewImage/Infrastructure/Cache"
import RKxnsxot "app/src/StreetViewImage/Infrastructure/Logger"
import gbLwVnqJ "app/src/StreetViewImage/Infrastructure/Server"
import PefLEOee "app/src/StreetViewImage/Presentation/Controller"

func GetRegistry() (registry TypeRegistry.Registry) {
	registry.RegistryStructs = append(registry.RegistryStructs, TypeRegistry.RegistryStruct{Name: "app/src/StreetViewImage/Application/Error.UserError", Implementation: mGQzNMon.UserError{}})
	registry.RegistryInterfaces = append(registry.RegistryInterfaces, TypeRegistry.RegistryInterface{Name: "app/src/StreetViewImage/Application/Error.ApplicationError", Implementation: (*mGQzNMon.ApplicationError)(nil)})
	registry.RegistryFactories = append(registry.RegistryFactories, TypeRegistry.RegistryFactory{Name: "app/src/StreetViewImage/Application/Error.ApplicationError", Implementations: []interface{}{mGQzNMon.NewApplicationError}})
	registry.RegistryInterfaces = append(registry.RegistryInterfaces, TypeRegistry.RegistryInterface{Name: "app/src/StreetViewImage/Application/Query.GetStreetViewImage", Implementation: (*poXJtEkr.GetStreetViewImage)(nil)})
	registry.RegistryFactories = append(registry.RegistryFactories, TypeRegistry.RegistryFactory{Name: "app/src/StreetViewImage/Application/Query.GetStreetViewImage", Implementations: []interface{}{poXJtEkr.NewGetStreetViewImageQuery}})
	registry.RegistryInterfaces = append(registry.RegistryInterfaces, TypeRegistry.RegistryInterface{Name: "app/src/StreetViewImage/Application/QueryHandler.GetStreetViewImageHandler", Implementation: (*mKaXayJi.GetStreetViewImageHandler)(nil)})
	registry.RegistryFactories = append(registry.RegistryFactories, TypeRegistry.RegistryFactory{Name: "app/src/StreetViewImage/Application/QueryHandler.GetStreetViewImageHandler", Implementations: []interface{}{mKaXayJi.NewGetStreetViewImageHandler}})
	registry.RegistryStructs = append(registry.RegistryStructs, TypeRegistry.RegistryStruct{Name: "app/src/StreetViewImage/Domain.ImageUuid", Implementation: GyZJpPBm.ImageUuid{}})
	registry.RegistryInterfaces = append(registry.RegistryInterfaces, TypeRegistry.RegistryInterface{Name: "app/src/StreetViewImage/Domain.StreetViewImages", Implementation: (*GyZJpPBm.StreetViewImages)(nil)})
	registry.RegistryInterfaces = append(registry.RegistryInterfaces, TypeRegistry.RegistryInterface{Name: "app/src/StreetViewImage/Domain.StreetViewImage", Implementation: (*GyZJpPBm.StreetViewImage)(nil)})
	registry.RegistryFactories = append(registry.RegistryFactories, TypeRegistry.RegistryFactory{Name: "app/src/StreetViewImage/Domain.StreetViewImage", Implementations: []interface{}{GyZJpPBm.NewStreetViewImage}})
	registry.RegistryFactories = append(registry.RegistryFactories, TypeRegistry.RegistryFactory{Name: "app/src/StreetViewImage/Domain.ImageUuid", Implementations: []interface{}{GyZJpPBm.NewImageUuid}})
	registry.RegistryStructs = append(registry.RegistryStructs, TypeRegistry.RegistryStruct{Name: "app/src/StreetViewImage/Infrastructure/ApiClient.RetrierFactory", Implementation: olJUMOFZ.RetrierFactory{}})
	registry.RegistryInterfaces = append(registry.RegistryInterfaces, TypeRegistry.RegistryInterface{Name: "app/src/StreetViewImage/Infrastructure/ApiClient.StreetViewApiClient", Implementation: (*olJUMOFZ.StreetViewApiClient)(nil)})
	registry.RegistryFactories = append(registry.RegistryFactories, TypeRegistry.RegistryFactory{Name: "app/src/StreetViewImage/Infrastructure/ApiClient.StreetViewApiClient", Implementations: []interface{}{olJUMOFZ.NewStreetViewApiClient}})
	registry.RegistryStructs = append(registry.RegistryStructs, TypeRegistry.RegistryStruct{Name: "app/src/StreetViewImage/Infrastructure/Cache.RedisStreetViewImages", Implementation: DpzQhmiZ.RedisStreetViewImages{}})
	registry.RegistryStructs = append(registry.RegistryStructs, TypeRegistry.RegistryStruct{Name: "app/src/StreetViewImage/Infrastructure/Cache.RedisClientFactory", Implementation: DpzQhmiZ.RedisClientFactory{}})
	registry.RegistryFactories = append(registry.RegistryFactories, TypeRegistry.RegistryFactory{Name: "app/src/StreetViewImage/Infrastructure/Cache.RedisClientFactory", Implementations: []interface{}{DpzQhmiZ.NewRedisClientFactory}})
	registry.RegistryStructs = append(registry.RegistryStructs, TypeRegistry.RegistryStruct{Name: "app/src/StreetViewImage/Infrastructure/Logger.FileLogger", Implementation: RKxnsxot.FileLogger{}})
	registry.RegistryStructs = append(registry.RegistryStructs, TypeRegistry.RegistryStruct{Name: "app/src/StreetViewImage/Infrastructure/Logger.ElasticSearchLogger", Implementation: RKxnsxot.ElasticSearchLogger{}})
	registry.RegistryStructs = append(registry.RegistryStructs, TypeRegistry.RegistryStruct{Name: "app/src/StreetViewImage/Infrastructure/Logger.LoggingStrategy", Implementation: RKxnsxot.LoggingStrategy{}})
	registry.RegistryStructs = append(registry.RegistryStructs, TypeRegistry.RegistryStruct{Name: "app/src/StreetViewImage/Infrastructure/Logger.RetrierFactory", Implementation: RKxnsxot.RetrierFactory{}})
	registry.RegistryStructs = append(registry.RegistryStructs, TypeRegistry.RegistryStruct{Name: "app/src/StreetViewImage/Infrastructure/Logger.RequiredLogFields", Implementation: RKxnsxot.RequiredLogFields{}})
	registry.RegistryInterfaces = append(registry.RegistryInterfaces, TypeRegistry.RegistryInterface{Name: "app/src/StreetViewImage/Infrastructure/Logger.Logger", Implementation: (*RKxnsxot.Logger)(nil)})
	registry.RegistryFactories = append(registry.RegistryFactories, TypeRegistry.RegistryFactory{Name: "app/src/StreetViewImage/Infrastructure/Logger.FileLogger", Implementations: []interface{}{RKxnsxot.NewFileLogger}})
	registry.RegistryFactories = append(registry.RegistryFactories, TypeRegistry.RegistryFactory{Name: "app/src/StreetViewImage/Infrastructure/Logger.ElasticSearchLogger", Implementations: []interface{}{RKxnsxot.NewElasticSearchLogger}})
	registry.RegistryFactories = append(registry.RegistryFactories, TypeRegistry.RegistryFactory{Name: "app/src/StreetViewImage/Infrastructure/Logger.LoggingStrategy", Implementations: []interface{}{RKxnsxot.NewLoggingStrategy}})
	registry.RegistryStructs = append(registry.RegistryStructs, TypeRegistry.RegistryStruct{Name: "app/src/StreetViewImage/Infrastructure/Server.RequestInterceptorGroup", Implementation: gbLwVnqJ.RequestInterceptorGroup{}})
	registry.RegistryStructs = append(registry.RegistryStructs, TypeRegistry.RegistryStruct{Name: "app/src/StreetViewImage/Infrastructure/Server.RetrierFactory", Implementation: gbLwVnqJ.RetrierFactory{}})
	registry.RegistryInterfaces = append(registry.RegistryInterfaces, TypeRegistry.RegistryInterface{Name: "app/src/StreetViewImage/Infrastructure/Server.GrpcErrorMapper", Implementation: (*gbLwVnqJ.GrpcErrorMapper)(nil)})
	registry.RegistryInterfaces = append(registry.RegistryInterfaces, TypeRegistry.RegistryInterface{Name: "app/src/StreetViewImage/Infrastructure/Server.GrpcServer", Implementation: (*gbLwVnqJ.GrpcServer)(nil)})
	registry.RegistryFactories = append(registry.RegistryFactories, TypeRegistry.RegistryFactory{Name: "app/src/StreetViewImage/Infrastructure/Server.GrpcErrorMapper", Implementations: []interface{}{gbLwVnqJ.NewGrpcErrorMapper}})
	registry.RegistryFactories = append(registry.RegistryFactories, TypeRegistry.RegistryFactory{Name: "app/src/StreetViewImage/Infrastructure/Server.GrpcServer", Implementations: []interface{}{gbLwVnqJ.New}})
	registry.RegistryStructs = append(registry.RegistryStructs, TypeRegistry.RegistryStruct{Name: "app/src/StreetViewImage/Presentation/Controller.GetStreetViewImageController", Implementation: PefLEOee.GetStreetViewImageController{}})
	registry.RegistryInterfaces = append(registry.RegistryInterfaces, TypeRegistry.RegistryInterface{Name: "app/src/StreetViewImage/Presentation/Controller.GrpcErrorMapper", Implementation: (*PefLEOee.GrpcErrorMapper)(nil)})

	return
}

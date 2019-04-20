package src

import "github.com/j7mbo/goij/src/TypeRegistry"
import YGQkDJvA "app/config"

func GetConfigRegistry() (registry TypeRegistry.Registry) {
	registry.RegistryStructs = append(registry.RegistryStructs, TypeRegistry.RegistryStruct{Name: "app/config.ElasticSearchConfiguration", Implementation: YGQkDJvA.ElasticSearchConfiguration{}})
	registry.RegistryStructs = append(registry.RegistryStructs, TypeRegistry.RegistryStruct{Name: "app/config.RedisConfiguration", Implementation: YGQkDJvA.RedisConfiguration{}})
	registry.RegistryStructs = append(registry.RegistryStructs, TypeRegistry.RegistryStruct{Name: "app/config.GrpcServerConfiguration", Implementation: YGQkDJvA.GrpcServerConfiguration{}})
	registry.RegistryStructs = append(registry.RegistryStructs, TypeRegistry.RegistryStruct{Name: "app/config.StreetViewApiConfiguration", Implementation: YGQkDJvA.StreetViewApiConfiguration{}})

	return
}

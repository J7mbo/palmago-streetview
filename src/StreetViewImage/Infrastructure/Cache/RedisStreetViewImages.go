package Cache

import (
	"app/src/StreetViewImage/Domain"
	"app/src/StreetViewImage/Infrastructure/Logger"
	"fmt"
	"github.com/go-redis/redis"
	"time"
)

/* redisKeyExpiration is the amount of time an image is stored in redis for. */
const redisKeyExpiration = time.Duration(1337 * time.Hour)

/* RedisStreetViewImages is a Repository responsible for persisting to Redis. */
type RedisStreetViewImages struct {
	RedisClientFactory RedisClientFactory
	/* redisClient is the factory's created client that is re-used until the connection fails. */
	redisClient *redis.Client
	Logger      Logger.LoggingStrategy
}

/* Save stores the image in Redis and returns whether or not this storing was successful. */
func (i *RedisStreetViewImages) Save(image Domain.StreetViewImage) bool {
	client := i.retrieveConnectedRedisClient()

	if client == nil {
		return false
	}

	redisKey := image.GetUuid()
	redisVal := i.marshalBytesForStorage(image.GetBytes())

	/* Pretty sure at this point that this won't fail, but you never know... */
	status := client.Set(redisKey, redisVal, redisKeyExpiration)

	if _, err := status.Result(); err != nil {
		i.Logger.Warning(fmt.Sprintf("Could not store value redis, reason: '%s'", err.Error()))

		return false
	}

	i.Logger.Debug(fmt.Sprintf("Stored key: '%s' in redis with byte length: '%d'", redisKey, len(redisVal)))

	return true
}

/* Find retrieves an image from persistence if one exists. */
func (i *RedisStreetViewImages) Find(latitude float64, longitude float64) Domain.StreetViewImage {
	client := i.retrieveConnectedRedisClient()

	if client == nil {
		return nil
	}

	imageUuid := Domain.NewImageUuid(latitude, longitude)

	imageBytes, err := client.Get(imageUuid.String()).Result()

	if err != nil {
		return nil
	}

	image, err := Domain.NewStreetViewImage(latitude, longitude, i.unmarshalStoredBytes(imageBytes))

	if err != nil {
		i.Logger.Warning(fmt.Sprintf("Image bytes retrieved from redis invalid, reason: '%s'", err.Error()))

		return nil
	}

	return image
}

/* retrieveConnectedRedisClient attempts to use the stored redis client to connect, otherwise builds a new one. */
func (i *RedisStreetViewImages) retrieveConnectedRedisClient() *redis.Client {
	if i.redisClient != nil {
		if _, err := i.redisClient.Ping().Result(); err == nil {
			return i.redisClient
		}
	}

	client, err := i.RedisClientFactory.Create()

	if err != nil {
		i.Logger.Warning(err.Error())

		return nil
	}

	i.redisClient = client

	return i.redisClient
}

/* marshalBytesForStorage converts a StreetViewImage's bytes into a format for storage as a redis value. */
func (i *RedisStreetViewImages) marshalBytesForStorage(bytes []byte) string {
	return string(bytes)
}

/* unmarshalStoredBytes converts a StreetViewImage's bytes, stored as a string value in redis, back to bytes. */
func (i *RedisStreetViewImages) unmarshalStoredBytes(str string) []byte {
	return []byte(str)
}

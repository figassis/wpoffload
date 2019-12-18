package util

import (
	"time"

	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"go.uber.org/zap"
	redis "gopkg.in/redis.v3"
)

const (
	Version    = "0.0.1"
	queueSpeed = 200
)

var (
	cacheKey     = "wpoffload"
	cacheTTL     = time.Hour
	logLevel     = "debug"
	awsAccessKey string
	awsSecretKey string
	awsRegion    = "us-east-1"
	bucket       string
	prefix       = ""
	watch        = "/data"
	schedule     = "* * * * *"
	logger       *zap.Logger
	cache        *redis.Client
	svc          *s3manager.Uploader
	hostname     string
	fileQueue    = map[string]bool{}
	acl          = "private"
)

type (
	Lock struct {
		Host string
	}
)

package util

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/shopspring/decimal"

	"github.com/gofrs/uuid"
	"github.com/thanhpk/randstr"
)

func Start(wg *sync.WaitGroup) (err error) {
	defer wg.Done()

	hostname, err = os.Hostname()
	if err = Log(err); err != nil {
		return
	}

	cacheKey = fmt.Sprintf("%s/%s/%s", cacheKey, bucket, prefix)
	if logger, err = newLogger(); err != nil {
		log.Fatalf("can't initialize logger: %v", err)
	}

	if os.Getenv("LOGLEVEL") != "" {
		logLevel = os.Getenv("LOGLEVEL")
	}
	if os.Getenv("AWS_ACCESS_KEY_ID") == "" || os.Getenv("AWS_SECRET_ACCESS_KEY") == "" || os.Getenv("BUCKET") == "" || os.Getenv("AWS_REGION") == "" {
		return errors.New("Invalid AWS settings")
	}
	awsAccessKey = os.Getenv("AWS_ACCESS_KEY_ID")
	awsSecretKey = os.Getenv("AWS_SECRET_ACCESS_KEY")
	bucket = os.Getenv("BUCKET")

	if os.Getenv("PREFIX") != "" {
		prefix = os.Getenv("PREFIX")
	}

	if os.Getenv("AWS_REGION") != "" {
		awsRegion = os.Getenv("AWS_REGION")
	}

	if os.Getenv("WATCH") != "" {
		watch = os.Getenv("WATCH")
		if err = EnsureDir(watch, 0755); err != nil {
			return
		}
	}

	if os.Getenv("SCHEDULE") != "" {
		schedule = os.Getenv("SCHEDULE")
	}

	if os.Getenv("REDIS_HOST") == "" || os.Getenv("REDIS_PORT") == "" {
		return errors.New("Invalid Redis settings")
	}

	loadRedis()
	if err = Log(setLock()); err != nil {
		return
	}

	svc = s3manager.NewUploader(session.New(aws.NewConfig().WithRegion(awsRegion)))

	//c := cron.New()
	// c.AddFunc(schedule, funcname)
	// c.Start()
	wg.Add(2)
	go offloader(wg)
	go watcher(wg)
	return
}

func UUID() (result string) {
	newUuid, err := uuid.NewV4()
	if err = Log(err); err != nil {
		return Random(36)
	}
	return newUuid.String()
}

func Random(n int) string {
	return randstr.String(n)
}

func RandInt32() int32 {
	return rand.New(rand.NewSource(time.Now().UnixNano())).Int31()
}

func RandInt64() int64 {
	return rand.New(rand.NewSource(time.Now().UnixNano())).Int63()
}

func RandDecimal() decimal.Decimal {
	return decimal.NewFromFloat(rand.New(rand.NewSource(time.Now().UnixNano())).Float64())
}

func RandDecimalN(n decimal.Decimal) decimal.Decimal {
	if n.LessThanOrEqual(decimal.NewFromFloat(0.0)) {
		return decimal.NewFromFloat(0.0)
	}
	return decimal.NewFromFloat(rand.New(rand.NewSource(time.Now().UnixNano())).Float64()).Mul(n)
}

func RandIntN(n int) int {
	if n <= 0 {
		return 0
	}
	return rand.New(rand.NewSource(time.Now().UnixNano())).Intn(n)
}

func EnsureDir(path string, perms os.FileMode) (err error) {
	if err = CheckPath(path); err == nil {
		return
	}

	return Log(makeDir(path, perms))
}

func EnsureDirs(paths []string, perms os.FileMode) (err error) {
	for _, path := range paths {
		if err = EnsureDir(path, perms); err != nil {
			return
		}
	}
	return
}

func CheckPath(path string) (err error) {
	if _, err = os.Stat(path); os.IsNotExist(err) {
		return errors.New(fmt.Sprintf("%s does not exist", path))
	}
	return nil
}

func makeDir(dir string, perms os.FileMode) (err error) {
	if err = os.MkdirAll(dir, perms); err != nil {
		return Log(errors.New("Could not create directory" + dir))
	}

	return
}

func setLock() (err error) {
	var lock Lock
	if err = GetCache("lock", &lock); err != nil || lock.Host == hostname {
		Cache("lock", Lock{Host: hostname})
	}

	if _, err = net.LookupIP(lock.Host); err == nil {
		return Log(fmt.Errorf("Directory is locked by %s", lock.Host))
	}

	Cache("lock", Lock{Host: hostname})

	return nil
}

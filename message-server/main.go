package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/opensourceways/community-robot-lib/logrusutil"
	liboptions "github.com/opensourceways/community-robot-lib/options"
	"github.com/sirupsen/logrus"

	"github.com/opensourceways/xihe-server/app"
	"github.com/opensourceways/xihe-server/config"
	"github.com/opensourceways/xihe-server/infrastructure/message"
	"github.com/opensourceways/xihe-server/infrastructure/mongodb"
	"github.com/opensourceways/xihe-server/infrastructure/repositories"
)

type options struct {
	service liboptions.ServiceOptions
}

func (o *options) Validate() error {
	return o.service.Validate()
}

func gatherOptions(fs *flag.FlagSet, args ...string) options {
	var o options

	o.service.AddFlags(fs)

	fs.Parse(args)
	return o
}

func main() {
	logrusutil.ComponentInit("xihe")
	log := logrus.NewEntry(logrus.StandardLogger())

	o := gatherOptions(
		flag.NewFlagSet(os.Args[0], flag.ExitOnError),
		os.Args[1:]...,
	)
	if err := o.Validate(); err != nil {
		logrus.Fatalf("Invalid options, err:%s", err.Error())
	}

	// cfg
	cfg := new(configuration)
	if err := config.LoadConfig(o.service.ConfigFile, cfg); err != nil {
		logrus.Fatalf("load config, err:%s", err.Error())
	}

	// mq
	topic := message.Topics{
		Like:      cfg.MQ.TopicLike,
		Following: cfg.MQ.TopicFollowing,
	}
	if err := message.Init(cfg.getMQConfig(), log, topic); err != nil {
		log.Fatalf("initialize mq failed, err:%v", err)
	}

	defer message.Exit(log)

	// mongo
	m := &cfg.Mongodb
	if err := mongodb.Initialize(m.MongodbConn, m.DBName); err != nil {
		logrus.Fatalf("initialize mongodb failed, err:%s", err.Error())
	}

	defer mongodb.Close()

	// cfg
	config.InitDomainConfig(&cfg.Resource, &cfg.User)

	// run
	run(newHandler(cfg, log), log)
}

func newHandler(cfg *configuration, log *logrus.Entry) *handler {
	return &handler{
		log:      log,
		maxRetry: cfg.MaxRetry,
		user: app.NewUserService(
			repositories.NewUserRepository(
				mongodb.NewUserMapper(cfg.Mongodb.UserCollection),
			),
			nil, nil,
		),

		project: app.NewProjectService(
			repositories.NewProjectRepository(
				mongodb.NewProjectMapper(cfg.Mongodb.ProjectCollection),
			),
			nil, nil,
		),

		dataset: app.NewDatasetService(
			repositories.NewDatasetRepository(
				mongodb.NewDatasetMapper(cfg.Mongodb.DatasetCollection),
			),
			nil, nil,
		),

		model: app.NewModelService(
			repositories.NewModelRepository(
				mongodb.NewModelMapper(cfg.Mongodb.ModelCollection),
			),
			nil, nil,
		),
	}
}

func run(h *handler, log *logrus.Entry) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	var wg sync.WaitGroup
	defer wg.Wait()

	called := false
	ctx, done := context.WithCancel(context.Background())

	defer func() {
		if !called {
			called = true
			done()
		}
	}()

	wg.Add(1)
	go func(ctx context.Context) {
		defer wg.Done()

		select {
		case <-ctx.Done():
			log.Info("receive done. exit normally")
			return

		case <-sig:
			log.Info("receive exit signal")
			done()
			called = true
			return
		}
	}(ctx)

	if err := message.Subscribe(ctx, h, log); err != nil {
		log.Errorf("subscribe failed, err:%v", err)
	}
}
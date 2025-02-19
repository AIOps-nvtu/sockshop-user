package api

import (
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/go-kit/kit/metrics"
	"github.com/microservices-demo/user/users"
)

// Middleware decorates a service.
type Middleware func(Service) Service

// LoggingMiddleware logs method calls, parameters, results, and elapsed time.
func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next Service) Service {
		return loggingMiddleware{
			next:   next,
			logger: logger,
		}
	}
}

type loggingMiddleware struct {
	next   Service
	logger log.Logger
}

func (mw loggingMiddleware) Login(username, password string) (user users.User, err error) {
	defer func(begin time.Time) {
		level.Info(mw.logger).Log(
			"method", "Login",
			"took", time.Since(begin),
		)
	}(time.Now())
	user, err = mw.next.Login(username, password)
	if err != nil {
		level.Error(mw.logger).Log(
			"method", "Login",
			"result", "error",
			"error", err.Error(),
		)
	}
	return user, err
}

func (mw loggingMiddleware) Register(username, password, email, first, last string) (string, error) {
	defer func(begin time.Time) {
		level.Info(mw.logger).Log(
			"method", "Register",
			"username", username,
			"email", email,
			"took", time.Since(begin),
		)
	}(time.Now())
	id, err := mw.next.Register(username, password, email, first, last)
	if err != nil {
		level.Error(mw.logger).Log(
			"method", "Login",
			"result", "error",
			"error", err.Error(),
		)
	}
	return id, err	
}

func (mw loggingMiddleware) PostUser(user users.User) (id string, err error) {
	defer func(begin time.Time) {
		level.Info(mw.logger).Log(
			"method", "PostUser",
			"username", user.Username,
			"email", user.Email,
			"result", id,
			"took", time.Since(begin),
		)
	}(time.Now())
	id, err = mw.next.PostUser(user)
	if err != nil {
		level.Error(mw.logger).Log(
			"method", "Login",
			"result", "error",
			"error", err.Error(),
		)
	}
	return id, err	
}

func (mw loggingMiddleware) GetUsers(id string) (u []users.User, err error) {
	defer func(begin time.Time) {
		who := id
		if who == "" {
			who = "all"
		}
		level.Info(mw.logger).Log(
			"method", "GetUsers",
			"id", who,
			"result", len(u),
			"took", time.Since(begin),
		)
	}(time.Now())
	u, err = mw.next.GetUsers(id)
	if err != nil {
		level.Error(mw.logger).Log(
			"method", "Login",
			"result", "error",
			"error", err.Error(),
		)
	}
	return u, err
}

func (mw loggingMiddleware) PostAddress(add users.Address, id string) (string, error) {
	defer func(begin time.Time) {
		level.Info(mw.logger).Log(
			"method", "PostAddress",
			"street", add.Street,
			"number", add.Number,
			"user", id,
			"took", time.Since(begin),
		)
	}(time.Now())
	adds, err := mw.next.PostAddress(add, id)
	if err != nil {
		level.Error(mw.logger).Log(
			"method", "Login",
			"result", "error",
			"error", err.Error(),
		)
	}
	return adds, err
}

func (mw loggingMiddleware) GetAddresses(id string) (a []users.Address, err error) {
	defer func(begin time.Time) {
		who := id
		if who == "" {
			who = "all"
		}
		level.Info(mw.logger).Log(
			"method", "GetAddresses",
			"id", who,
			"result", len(a),
			"took", time.Since(begin),
		)
	}(time.Now())
	a, err = mw.next.GetAddresses(id)
	if err != nil {
		level.Error(mw.logger).Log(
			"method", "Login",
			"result", "error",
			"error", err.Error(),
		)
	}
	return a, err	
}

func (mw loggingMiddleware) PostCard(card users.Card, id string) (string, error) {
	defer func(begin time.Time) {
		cc := card
		cc.MaskCC()
		level.Info(mw.logger).Log(
			"method", "PostCard",
			"card", cc.LongNum,
			"user", id,
			"took", time.Since(begin),
		)
	}(time.Now())
	id, err := mw.next.PostCard(card, id)
	if err != nil {
		level.Error(mw.logger).Log(
			"method", "Login",
			"result", "error",
			"error", err.Error(),
		)
	}
	return id, err
}

func (mw loggingMiddleware) GetCards(id string) (a []users.Card, err error) {
	defer func(begin time.Time) {
		who := id
		if who == "" {
			who = "all"
		}
		level.Info(mw.logger).Log(
			"method", "GetCards",
			"id", who,
			"result", len(a),
			"took", time.Since(begin),
		)
	}(time.Now())
	a, err = mw.next.GetCards(id)
	if err != nil {
		level.Error(mw.logger).Log(
			"method", "Login",
			"result", "error",
			"error", err.Error(),
		)
	}
	return a, err
}

func (mw loggingMiddleware) Delete(entity, id string) (err error) {
	defer func(begin time.Time) {
		level.Info(mw.logger).Log(
			"method", "Delete",
			"entity", entity,
			"id", id,
			"took", time.Since(begin),
		)
	}(time.Now())
	err = mw.next.Delete(entity, id)
	if err != nil {
		level.Error(mw.logger).Log(
			"method", "Login",
			"result", "error",
			"error", err.Error(),
		)
	}
	return err
}

func (mw loggingMiddleware) Health() (health []Health) {
	defer func(begin time.Time) {
		level.Info(mw.logger).Log(
			"method", "Health",
			"result", len(health),
			"took", time.Since(begin),
		)
	}(time.Now())
	return mw.next.Health()
}

type instrumentingService struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	Service
}

// NewInstrumentingService returns an instance of an instrumenting Service.
func NewInstrumentingService(requestCount metrics.Counter, requestLatency metrics.Histogram, s Service) Service {
	return &instrumentingService{
		requestCount:   requestCount,
		requestLatency: requestLatency,
		Service:        s,
	}
}

func (s *instrumentingService) Login(username, password string) (users.User, error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "login").Add(1)
		s.requestLatency.With("method", "login").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.Service.Login(username, password)
}

func (s *instrumentingService) Register(username, password, email, first, last string) (string, error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "register").Add(1)
		s.requestLatency.With("method", "register").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.Service.Register(username, password, email, first, last)
}

func (s *instrumentingService) PostUser(user users.User) (string, error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "postUser").Add(1)
		s.requestLatency.With("method", "postUser").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.Service.PostUser(user)
}

func (s *instrumentingService) GetUsers(id string) (u []users.User, err error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "getUsers").Add(1)
		s.requestLatency.With("method", "getUsers").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.Service.GetUsers(id)
}

func (s *instrumentingService) PostAddress(add users.Address, id string) (string, error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "postAddress").Add(1)
		s.requestLatency.With("method", "postAddress").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.Service.PostAddress(add, id)
}

func (s *instrumentingService) GetAddresses(id string) ([]users.Address, error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "getAddresses").Add(1)
		s.requestLatency.With("method", "getAddresses").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.Service.GetAddresses(id)
}

func (s *instrumentingService) PostCard(card users.Card, id string) (string, error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "postCard").Add(1)
		s.requestLatency.With("method", "postCard").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.Service.PostCard(card, id)
}

func (s *instrumentingService) GetCards(id string) ([]users.Card, error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "getCards").Add(1)
		s.requestLatency.With("method", "getCards").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.Service.GetCards(id)
}

func (s *instrumentingService) Delete(entity, id string) error {
	defer func(begin time.Time) {
		s.requestCount.With("method", "delete").Add(1)
		s.requestLatency.With("method", "delete").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.Service.Delete(entity, id)
}

func (s *instrumentingService) Health() []Health {
	defer func(begin time.Time) {
		s.requestCount.With("method", "health").Add(1)
		s.requestLatency.With("method", "health").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.Service.Health()
}

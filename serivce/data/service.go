package data

import (
	"github.com/cjburchell/queue/serivce/data/models"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type service struct {
	session *mgo.Session
}

func (s service) DeleteJob(jobId string) error {
	session := s.session.Clone()
	defer session.Close()

	return errors.WithStack(session.DB(dbName).C(jobsCollection).RemoveId(jobId))
}
func getTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func (s service) GetNextJob(maxJobTime int64) (*models.Job, error) {
	session := s.session.Clone()
	defer session.Close()

	timestamp := getTimestamp()
	timeToComplete := getFutureTimeStamp(maxJobTime)

	change := mgo.Change{
		Update:    bson.M{"$set": bson.M{"timestamp": timeToComplete}, "$inc": bson.M{"tries": 1, "priority": -1}},
		ReturnNew: true,
	}

	findQuery := bson.M{"timestamp": bson.M{"$lte": timestamp}}

	var job models.Job
	info, err := session.DB(dbName).C(jobsCollection).Find(findQuery).Sort("-priority", "timestamp").Apply(change, &job)

	if err != nil {
		if err == mgo.ErrNotFound {
			return nil, nil
		} else {
			return nil, errors.WithStack(err)
		}
	} else if info.Updated == 0 {
		return nil, nil
	}

	return &job, nil
}

func getFutureTimeStamp(milli int64) int64 {
	return time.Now().Add(time.Millisecond*time.Duration(milli)).UnixNano() / (int64(time.Millisecond))
}

func (s service) DelayJob(jobId string, delay int64, priority int) error {
	session := s.session.Clone()
	defer session.Close()

	DelayTime := getFutureTimeStamp(delay)
	return errors.WithStack(session.DB(dbName).C(jobsCollection).UpdateId(jobId, bson.M{"$set": bson.M{"timestamp": DelayTime, "priority": priority}}))
}

const dbName = "Queue"
const jobsCollection = "Jobs"

func (s service) AddJob(call models.Call, repeat bool, delay int64, retries int, priority int) error {
	session := s.session.Clone()
	defer session.Close()

	u1 := uuid.NewV4()
	job := models.Job{
		Id:              u1.String(),
		Call:            call,
		MaxRetries:      retries,
		Delay:           delay,
		Repeat:          repeat,
		InitialPriority: priority,
		Tries:           0,
		Priority:        priority,
		Timestamp:       getFutureTimeStamp(delay),
	}

	err := session.DB(dbName).C(jobsCollection).Insert(job)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (s *service) setup(address string) error {
	session, err := mgo.Dial(address)
	if err != nil {
		return errors.WithStack(err)
	}

	err = session.Ping()
	if err != nil {
		return errors.WithStack(err)
	}

	s.session = session
	return nil
}

func NewService(address string) (IService, error) {
	service := &service{}
	err := service.setup(address)
	return service, err
}

package events

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type TestEvent struct {
	Name    string
	Payload any
}

func (e *TestEvent) GetName() string {
	return e.Name
}

func (e *TestEvent) GetPayload() any {
	return e.Payload
}

func (e *TestEvent) GetDateTime() time.Time {
	return time.Now()
}

type TestEventHandler struct {
	ID int
}

func (h *TestEventHandler) Handle(event Event, wg *sync.WaitGroup) {}

type DispatcherTestSuite struct {
	suite.Suite
	event           TestEvent
	event3          TestEvent
	handler         TestEventHandler
	handler3        TestEventHandler
	handler4        TestEventHandler
	eventDispatcher *Dispatcher
}

func (suite *DispatcherTestSuite) SetupTest() {
	suite.eventDispatcher = NewDispatcher()
	suite.handler = TestEventHandler{ID: 2}
	suite.handler3 = TestEventHandler{ID: 3}
	suite.handler4 = TestEventHandler{ID: 4}
	suite.event = TestEvent{Name: "test", Payload: "test"}
	suite.event3 = TestEvent{Name: "test2", Payload: "test2"}
}

func (suite *DispatcherTestSuite) TestEventDispatcher_Register() {
	err := suite.eventDispatcher.Register(suite.event.GetName(), &suite.handler)
	suite.Nil(err)
	suite.Equal(1, len(suite.eventDispatcher.handlers[suite.event.GetName()]))

	err = suite.eventDispatcher.Register(suite.event.GetName(), &suite.handler3)
	suite.Nil(err)
	suite.Equal(2, len(suite.eventDispatcher.handlers[suite.event.GetName()]))

	assert.Equal(
		suite.T(),
		&suite.handler,
		suite.eventDispatcher.handlers[suite.event.GetName()][0],
	)
	assert.Equal(
		suite.T(),
		&suite.handler3,
		suite.eventDispatcher.handlers[suite.event.GetName()][1],
	)
}

func (suite *DispatcherTestSuite) TestEventDispatcher_Register_WithSameHandler() {
	err := suite.eventDispatcher.Register(suite.event.GetName(), &suite.handler)
	suite.Nil(err)
	suite.Equal(1, len(suite.eventDispatcher.handlers[suite.event.GetName()]))

	err = suite.eventDispatcher.Register(suite.event.GetName(), &suite.handler)
	suite.Equal(ErrHandlerAlreadyRegistered, err)
	suite.Equal(1, len(suite.eventDispatcher.handlers[suite.event.GetName()]))
}

func (suite *DispatcherTestSuite) TestEventDispatcher_Clear() {
	// Event 2
	err := suite.eventDispatcher.Register(suite.event.GetName(), &suite.handler)
	suite.Nil(err)
	suite.Equal(1, len(suite.eventDispatcher.handlers[suite.event.GetName()]))

	err = suite.eventDispatcher.Register(suite.event.GetName(), &suite.handler3)
	suite.Nil(err)
	suite.Equal(2, len(suite.eventDispatcher.handlers[suite.event.GetName()]))

	// Event 3
	err = suite.eventDispatcher.Register(suite.event3.GetName(), &suite.handler3)
	suite.Nil(err)
	suite.Equal(1, len(suite.eventDispatcher.handlers[suite.event3.GetName()]))

	suite.eventDispatcher.Clear()
	suite.Equal(0, len(suite.eventDispatcher.handlers))
}

func (suite *DispatcherTestSuite) TestEventDispatcher_Has() {
	// Event 2
	err := suite.eventDispatcher.Register(suite.event.GetName(), &suite.handler)
	suite.Nil(err)
	suite.Equal(1, len(suite.eventDispatcher.handlers[suite.event.GetName()]))

	err = suite.eventDispatcher.Register(suite.event.GetName(), &suite.handler3)
	suite.Nil(err)
	suite.Equal(2, len(suite.eventDispatcher.handlers[suite.event.GetName()]))

	assert.True(suite.T(), suite.eventDispatcher.Has(suite.event.GetName(), &suite.handler))
	assert.True(suite.T(), suite.eventDispatcher.Has(suite.event.GetName(), &suite.handler3))
	assert.False(suite.T(), suite.eventDispatcher.Has(suite.event.GetName(), &suite.handler4))
}

func (suite *DispatcherTestSuite) TestEventDispatcher_Remove() {
	// Event 2
	err := suite.eventDispatcher.Register(suite.event.GetName(), &suite.handler)
	suite.Nil(err)
	suite.Equal(1, len(suite.eventDispatcher.handlers[suite.event.GetName()]))

	err = suite.eventDispatcher.Register(suite.event.GetName(), &suite.handler3)
	suite.Nil(err)
	suite.Equal(2, len(suite.eventDispatcher.handlers[suite.event.GetName()]))

	// Event 3
	err = suite.eventDispatcher.Register(suite.event3.GetName(), &suite.handler3)
	suite.Nil(err)
	suite.Equal(1, len(suite.eventDispatcher.handlers[suite.event3.GetName()]))

	suite.eventDispatcher.Remove(suite.event.GetName(), &suite.handler)
	suite.Equal(1, len(suite.eventDispatcher.handlers[suite.event.GetName()]))
	assert.Equal(
		suite.T(),
		&suite.handler3,
		suite.eventDispatcher.handlers[suite.event.GetName()][0],
	)

	suite.eventDispatcher.Remove(suite.event.GetName(), &suite.handler3)
	suite.Equal(0, len(suite.eventDispatcher.handlers[suite.event.GetName()]))

	suite.eventDispatcher.Remove(suite.event3.GetName(), &suite.handler3)
	suite.Equal(0, len(suite.eventDispatcher.handlers[suite.event3.GetName()]))
}

type MockHandler struct {
	mock.Mock
}

func (m *MockHandler) Handle(event Event, wg *sync.WaitGroup) {
	m.Called(event)
	wg.Done()
}

func (suite *DispatcherTestSuite) TestEventDispatch_Dispatch() {
	eh := &MockHandler{}
	eh.On("Handle", &suite.event)

	eh3 := &MockHandler{}
	eh3.On("Handle", &suite.event)

	suite.eventDispatcher.Register(suite.event.GetName(), eh)
	suite.eventDispatcher.Register(suite.event.GetName(), eh3)

	suite.eventDispatcher.Dispatch(&suite.event)
	eh.AssertExpectations(suite.T())
	eh3.AssertExpectations(suite.T())
	eh.AssertNumberOfCalls(suite.T(), "Handle", 1)
	eh3.AssertNumberOfCalls(suite.T(), "Handle", 1)
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(DispatcherTestSuite))
}

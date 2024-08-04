// Code generated by http://github.com/gojuno/minimock (v3.3.10). DO NOT EDIT.

package mocks

import (
	"context"
	"sync"
	mm_atomic "sync/atomic"
	mm_time "time"

	"github.com/Genvekt/cli-chat/services/chat-server/internal/model"
	"github.com/gojuno/minimock/v3"
)

// AuthClientMock implements service.AuthClient
type AuthClientMock struct {
	t          minimock.Tester
	finishOnce sync.Once

	funcGetList          func(ctx context.Context, usernames []string) (upa1 []*model.User, err error)
	inspectFuncGetList   func(ctx context.Context, usernames []string)
	afterGetListCounter  uint64
	beforeGetListCounter uint64
	GetListMock          mAuthClientMockGetList
}

// NewAuthClientMock returns a mock for service.AuthClient
func NewAuthClientMock(t minimock.Tester) *AuthClientMock {
	m := &AuthClientMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.GetListMock = mAuthClientMockGetList{mock: m}
	m.GetListMock.callArgs = []*AuthClientMockGetListParams{}

	t.Cleanup(m.MinimockFinish)

	return m
}

type mAuthClientMockGetList struct {
	optional           bool
	mock               *AuthClientMock
	defaultExpectation *AuthClientMockGetListExpectation
	expectations       []*AuthClientMockGetListExpectation

	callArgs []*AuthClientMockGetListParams
	mutex    sync.RWMutex

	expectedInvocations uint64
}

// AuthClientMockGetListExpectation specifies expectation struct of the AuthClient.GetList
type AuthClientMockGetListExpectation struct {
	mock      *AuthClientMock
	params    *AuthClientMockGetListParams
	paramPtrs *AuthClientMockGetListParamPtrs
	results   *AuthClientMockGetListResults
	Counter   uint64
}

// AuthClientMockGetListParams contains parameters of the AuthClient.GetList
type AuthClientMockGetListParams struct {
	ctx       context.Context
	usernames []string
}

// AuthClientMockGetListParamPtrs contains pointers to parameters of the AuthClient.GetList
type AuthClientMockGetListParamPtrs struct {
	ctx       *context.Context
	usernames *[]string
}

// AuthClientMockGetListResults contains results of the AuthClient.GetList
type AuthClientMockGetListResults struct {
	upa1 []*model.User
	err  error
}

// Marks this method to be optional. The default behavior of any method with Return() is '1 or more', meaning
// the test will fail minimock's automatic final call check if the mocked method was not called at least once.
// Optional() makes method check to work in '0 or more' mode.
// It is NOT RECOMMENDED to use this option by default unless you really need it, as it helps to
// catch the problems when the expected method call is totally skipped during test run.
func (mmGetList *mAuthClientMockGetList) Optional() *mAuthClientMockGetList {
	mmGetList.optional = true
	return mmGetList
}

// Expect sets up expected params for AuthClient.GetList
func (mmGetList *mAuthClientMockGetList) Expect(ctx context.Context, usernames []string) *mAuthClientMockGetList {
	if mmGetList.mock.funcGetList != nil {
		mmGetList.mock.t.Fatalf("AuthClientMock.GetList mock is already set by Set")
	}

	if mmGetList.defaultExpectation == nil {
		mmGetList.defaultExpectation = &AuthClientMockGetListExpectation{}
	}

	if mmGetList.defaultExpectation.paramPtrs != nil {
		mmGetList.mock.t.Fatalf("AuthClientMock.GetList mock is already set by ExpectParams functions")
	}

	mmGetList.defaultExpectation.params = &AuthClientMockGetListParams{ctx, usernames}
	for _, e := range mmGetList.expectations {
		if minimock.Equal(e.params, mmGetList.defaultExpectation.params) {
			mmGetList.mock.t.Fatalf("Expectation set by When has same params: %#v", *mmGetList.defaultExpectation.params)
		}
	}

	return mmGetList
}

// ExpectCtxParam1 sets up expected param ctx for AuthClient.GetList
func (mmGetList *mAuthClientMockGetList) ExpectCtxParam1(ctx context.Context) *mAuthClientMockGetList {
	if mmGetList.mock.funcGetList != nil {
		mmGetList.mock.t.Fatalf("AuthClientMock.GetList mock is already set by Set")
	}

	if mmGetList.defaultExpectation == nil {
		mmGetList.defaultExpectation = &AuthClientMockGetListExpectation{}
	}

	if mmGetList.defaultExpectation.params != nil {
		mmGetList.mock.t.Fatalf("AuthClientMock.GetList mock is already set by Expect")
	}

	if mmGetList.defaultExpectation.paramPtrs == nil {
		mmGetList.defaultExpectation.paramPtrs = &AuthClientMockGetListParamPtrs{}
	}
	mmGetList.defaultExpectation.paramPtrs.ctx = &ctx

	return mmGetList
}

// ExpectUsernamesParam2 sets up expected param usernames for AuthClient.GetList
func (mmGetList *mAuthClientMockGetList) ExpectUsernamesParam2(usernames []string) *mAuthClientMockGetList {
	if mmGetList.mock.funcGetList != nil {
		mmGetList.mock.t.Fatalf("AuthClientMock.GetList mock is already set by Set")
	}

	if mmGetList.defaultExpectation == nil {
		mmGetList.defaultExpectation = &AuthClientMockGetListExpectation{}
	}

	if mmGetList.defaultExpectation.params != nil {
		mmGetList.mock.t.Fatalf("AuthClientMock.GetList mock is already set by Expect")
	}

	if mmGetList.defaultExpectation.paramPtrs == nil {
		mmGetList.defaultExpectation.paramPtrs = &AuthClientMockGetListParamPtrs{}
	}
	mmGetList.defaultExpectation.paramPtrs.usernames = &usernames

	return mmGetList
}

// Inspect accepts an inspector function that has same arguments as the AuthClient.GetList
func (mmGetList *mAuthClientMockGetList) Inspect(f func(ctx context.Context, usernames []string)) *mAuthClientMockGetList {
	if mmGetList.mock.inspectFuncGetList != nil {
		mmGetList.mock.t.Fatalf("Inspect function is already set for AuthClientMock.GetList")
	}

	mmGetList.mock.inspectFuncGetList = f

	return mmGetList
}

// Return sets up results that will be returned by AuthClient.GetList
func (mmGetList *mAuthClientMockGetList) Return(upa1 []*model.User, err error) *AuthClientMock {
	if mmGetList.mock.funcGetList != nil {
		mmGetList.mock.t.Fatalf("AuthClientMock.GetList mock is already set by Set")
	}

	if mmGetList.defaultExpectation == nil {
		mmGetList.defaultExpectation = &AuthClientMockGetListExpectation{mock: mmGetList.mock}
	}
	mmGetList.defaultExpectation.results = &AuthClientMockGetListResults{upa1, err}
	return mmGetList.mock
}

// Set uses given function f to mock the AuthClient.GetList method
func (mmGetList *mAuthClientMockGetList) Set(f func(ctx context.Context, usernames []string) (upa1 []*model.User, err error)) *AuthClientMock {
	if mmGetList.defaultExpectation != nil {
		mmGetList.mock.t.Fatalf("Default expectation is already set for the AuthClient.GetList method")
	}

	if len(mmGetList.expectations) > 0 {
		mmGetList.mock.t.Fatalf("Some expectations are already set for the AuthClient.GetList method")
	}

	mmGetList.mock.funcGetList = f
	return mmGetList.mock
}

// When sets expectation for the AuthClient.GetList which will trigger the result defined by the following
// Then helper
func (mmGetList *mAuthClientMockGetList) When(ctx context.Context, usernames []string) *AuthClientMockGetListExpectation {
	if mmGetList.mock.funcGetList != nil {
		mmGetList.mock.t.Fatalf("AuthClientMock.GetList mock is already set by Set")
	}

	expectation := &AuthClientMockGetListExpectation{
		mock:   mmGetList.mock,
		params: &AuthClientMockGetListParams{ctx, usernames},
	}
	mmGetList.expectations = append(mmGetList.expectations, expectation)
	return expectation
}

// Then sets up AuthClient.GetList return parameters for the expectation previously defined by the When method
func (e *AuthClientMockGetListExpectation) Then(upa1 []*model.User, err error) *AuthClientMock {
	e.results = &AuthClientMockGetListResults{upa1, err}
	return e.mock
}

// Times sets number of times AuthClient.GetList should be invoked
func (mmGetList *mAuthClientMockGetList) Times(n uint64) *mAuthClientMockGetList {
	if n == 0 {
		mmGetList.mock.t.Fatalf("Times of AuthClientMock.GetList mock can not be zero")
	}
	mm_atomic.StoreUint64(&mmGetList.expectedInvocations, n)
	return mmGetList
}

func (mmGetList *mAuthClientMockGetList) invocationsDone() bool {
	if len(mmGetList.expectations) == 0 && mmGetList.defaultExpectation == nil && mmGetList.mock.funcGetList == nil {
		return true
	}

	totalInvocations := mm_atomic.LoadUint64(&mmGetList.mock.afterGetListCounter)
	expectedInvocations := mm_atomic.LoadUint64(&mmGetList.expectedInvocations)

	return totalInvocations > 0 && (expectedInvocations == 0 || expectedInvocations == totalInvocations)
}

// GetList implements service.AuthClient
func (mmGetList *AuthClientMock) GetList(ctx context.Context, usernames []string) (upa1 []*model.User, err error) {
	mm_atomic.AddUint64(&mmGetList.beforeGetListCounter, 1)
	defer mm_atomic.AddUint64(&mmGetList.afterGetListCounter, 1)

	if mmGetList.inspectFuncGetList != nil {
		mmGetList.inspectFuncGetList(ctx, usernames)
	}

	mm_params := AuthClientMockGetListParams{ctx, usernames}

	// Record call args
	mmGetList.GetListMock.mutex.Lock()
	mmGetList.GetListMock.callArgs = append(mmGetList.GetListMock.callArgs, &mm_params)
	mmGetList.GetListMock.mutex.Unlock()

	for _, e := range mmGetList.GetListMock.expectations {
		if minimock.Equal(*e.params, mm_params) {
			mm_atomic.AddUint64(&e.Counter, 1)
			return e.results.upa1, e.results.err
		}
	}

	if mmGetList.GetListMock.defaultExpectation != nil {
		mm_atomic.AddUint64(&mmGetList.GetListMock.defaultExpectation.Counter, 1)
		mm_want := mmGetList.GetListMock.defaultExpectation.params
		mm_want_ptrs := mmGetList.GetListMock.defaultExpectation.paramPtrs

		mm_got := AuthClientMockGetListParams{ctx, usernames}

		if mm_want_ptrs != nil {

			if mm_want_ptrs.ctx != nil && !minimock.Equal(*mm_want_ptrs.ctx, mm_got.ctx) {
				mmGetList.t.Errorf("AuthClientMock.GetList got unexpected parameter ctx, want: %#v, got: %#v%s\n", *mm_want_ptrs.ctx, mm_got.ctx, minimock.Diff(*mm_want_ptrs.ctx, mm_got.ctx))
			}

			if mm_want_ptrs.usernames != nil && !minimock.Equal(*mm_want_ptrs.usernames, mm_got.usernames) {
				mmGetList.t.Errorf("AuthClientMock.GetList got unexpected parameter usernames, want: %#v, got: %#v%s\n", *mm_want_ptrs.usernames, mm_got.usernames, minimock.Diff(*mm_want_ptrs.usernames, mm_got.usernames))
			}

		} else if mm_want != nil && !minimock.Equal(*mm_want, mm_got) {
			mmGetList.t.Errorf("AuthClientMock.GetList got unexpected parameters, want: %#v, got: %#v%s\n", *mm_want, mm_got, minimock.Diff(*mm_want, mm_got))
		}

		mm_results := mmGetList.GetListMock.defaultExpectation.results
		if mm_results == nil {
			mmGetList.t.Fatal("No results are set for the AuthClientMock.GetList")
		}
		return (*mm_results).upa1, (*mm_results).err
	}
	if mmGetList.funcGetList != nil {
		return mmGetList.funcGetList(ctx, usernames)
	}
	mmGetList.t.Fatalf("Unexpected call to AuthClientMock.GetList. %v %v", ctx, usernames)
	return
}

// GetListAfterCounter returns a count of finished AuthClientMock.GetList invocations
func (mmGetList *AuthClientMock) GetListAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmGetList.afterGetListCounter)
}

// GetListBeforeCounter returns a count of AuthClientMock.GetList invocations
func (mmGetList *AuthClientMock) GetListBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmGetList.beforeGetListCounter)
}

// Calls returns a list of arguments used in each call to AuthClientMock.GetList.
// The list is in the same order as the calls were made (i.e. recent calls have a higher index)
func (mmGetList *mAuthClientMockGetList) Calls() []*AuthClientMockGetListParams {
	mmGetList.mutex.RLock()

	argCopy := make([]*AuthClientMockGetListParams, len(mmGetList.callArgs))
	copy(argCopy, mmGetList.callArgs)

	mmGetList.mutex.RUnlock()

	return argCopy
}

// MinimockGetListDone returns true if the count of the GetList invocations corresponds
// the number of defined expectations
func (m *AuthClientMock) MinimockGetListDone() bool {
	if m.GetListMock.optional {
		// Optional methods provide '0 or more' call count restriction.
		return true
	}

	for _, e := range m.GetListMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	return m.GetListMock.invocationsDone()
}

// MinimockGetListInspect logs each unmet expectation
func (m *AuthClientMock) MinimockGetListInspect() {
	for _, e := range m.GetListMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to AuthClientMock.GetList with params: %#v", *e.params)
		}
	}

	afterGetListCounter := mm_atomic.LoadUint64(&m.afterGetListCounter)
	// if default expectation was set then invocations count should be greater than zero
	if m.GetListMock.defaultExpectation != nil && afterGetListCounter < 1 {
		if m.GetListMock.defaultExpectation.params == nil {
			m.t.Error("Expected call to AuthClientMock.GetList")
		} else {
			m.t.Errorf("Expected call to AuthClientMock.GetList with params: %#v", *m.GetListMock.defaultExpectation.params)
		}
	}
	// if func was set then invocations count should be greater than zero
	if m.funcGetList != nil && afterGetListCounter < 1 {
		m.t.Error("Expected call to AuthClientMock.GetList")
	}

	if !m.GetListMock.invocationsDone() && afterGetListCounter > 0 {
		m.t.Errorf("Expected %d calls to AuthClientMock.GetList but found %d calls",
			mm_atomic.LoadUint64(&m.GetListMock.expectedInvocations), afterGetListCounter)
	}
}

// MinimockFinish checks that all mocked methods have been called the expected number of times
func (m *AuthClientMock) MinimockFinish() {
	m.finishOnce.Do(func() {
		if !m.minimockDone() {
			m.MinimockGetListInspect()
			m.t.FailNow()
		}
	})
}

// MinimockWait waits for all mocked methods to be called the expected number of times
func (m *AuthClientMock) MinimockWait(timeout mm_time.Duration) {
	timeoutCh := mm_time.After(timeout)
	for {
		if m.minimockDone() {
			return
		}
		select {
		case <-timeoutCh:
			m.MinimockFinish()
			return
		case <-mm_time.After(10 * mm_time.Millisecond):
		}
	}
}

func (m *AuthClientMock) minimockDone() bool {
	done := true
	return done &&
		m.MinimockGetListDone()
}

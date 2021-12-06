# Upgrade Kubernetes Controller Runtime from v0.6.5
This blog concentrates on test case migration as the real code migration pretty detailed and straight forward on official kubernetes websites.

Are you still using version sigs.k8s.io/controller-runtime v0.6.5, time to rethink your choice, as it has costlier bugs that ends up in event failure(s) or misses. The projects that uses the kubernetes operators extensively depends on optimal management of events.

Kubernetes upgrades related to operators usually comes up with breaking this time along with test cases.

The *reactors* were once been handy to test the edge scenarios, as the projects grows, test cases grows exponentially, and when it comes to a method that is removed / unusable from core library however leveraged extensively by dependant or child projects... 😭

[stack overflow](https://stackoverflow.com/questions/67121718/is-there-a-way-to-end-to-end-test-a-controller-runtime-operator-in-conjunction-w)

## Major Hurdle
The reactors are unusable via [fakes](https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.8.3/pkg/client/fake) along with [ClientBuilder](https://pkg.go.dev/sigs.k8s.io/controller-runtime/pkg/client/fake#NewClientBuilder).

### Test Cases
The operators (or) reconciler's test cases are heavily depndent on reactors and it is not a good idea to add those to technical debts or just ignore them.

## FakeReactor instead of Reactor
Based on researches digging deep into the kubernetes libraries and some searches, generalized the usage of *FakeReactor*, that could be imported and used by any project.

### How to use

```
go get github.com/DreamAmbitious/k8s-controller-runtime-fake
```

- Replace your reactors as below

```
mp := make(map[string]ctrlfake.FakeRunTimeMock)
<!-- CRUD methods create/read/update/delete in lower case, kind refers resource object kind of your operator. -->
mp["method_kind"] = ctrlfake.FakeRunTimeMock{
	MockFunc: func() (client.Object, error) {
		return nil, errs.New("fake output that you're returning")
	},
)
fakeV1alpha1Client := &ctrlfake.FakeReactingCtrlRuntimeClient{
	Client:          reconcilerClient,
	FakeRunTimeMock: mp,
}
```

:do_not_litter: As this a mock function, it comes up with more power and responsibility.

- Refer kubernetes [NewClientBuilder](https://pkg.go.dev/sigs.k8s.io/controller-runtime/pkg/client/fake#NewClientBuilder) for more info on usage , here is the gist.

```
initObjects := []runtime.Object{
	// load all the objects, that is required during runtime while testing
}
k8client := fakectrlruntime.NewClientBuilder().WithScheme(fakeRuntimeScheme).WithRuntimeObjects(initObjects...).Build()
```

## Conclusion
The library can be leverage to test the edge and more complex cases.
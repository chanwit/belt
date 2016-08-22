package scheduler

type Scheduler interface {

	// Invoked when resources have been offered to this framework. A
	// single offer will only contain resources from a single slave.
	// Resources associated with an offer will not be re-offered to
	// _this_ framework until either (a) this framework has rejected
	// those resources (see SchedulerDriver::launchTasks) or (b) those
	// resources have been rescinded (see Scheduler::offerRescinded).
	// Note that resources may be concurrently offered to more than one
	// framework at a time (depending on the allocator being used). In
	// that case, the first framework to launch tasks using those
	// resources will be able to use them while the other frameworks
	// will have those resources rescinded (or if a framework has
	// already launched tasks with those resources then those tasks will
	// fail with a TASK_LOST status and a message saying as much).
	ResourceOffers(offers []*api.Node)

	// Invoked when an offer is no longer valid (e.g., the slave was
	// lost or another framework used resources in the offer). If for
	// whatever reason an offer is never rescinded (e.g., dropped
	// message, failing over framework, etc.), a framwork that attempts
	// to launch tasks using an invalid offer will receive TASK_LOST
	// status updates for those tasks (see Scheduler::resourceOffers).
	OfferRescinded(*api.Node)
}

type BigDataScheduler struct {
}

const (
	CPUS_PER_TASK = 1
	MEM_PER_TASK  = 64
)

//
// belt service create --name staging/test \
// 		--replicas 2 \
//		--scheduler bigdata
//
func (sched *BigDataScheduler) ResourceOffers(spec, offers []*api.Node) {
	namespace := spec.Namespace()
	for _, offer := range offers {

	}
	// offerred nodes OK
	//  1. switch docker_host to staging's leader
	//  2. remove node from resources
	//  3. join node to staging
	//  4.
}

func init() {
	Register("bigdata", &BigDataScheduler{})
}

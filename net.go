package tincd

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
	"tinc-boot/tincd/api/impl/apiclient"
	"tinc-boot/tincd/api/impl/apiserver"
	"tinc-boot/tincd/internal"
	"tinc-boot/tincd/network"
	"tinc-boot/tincd/runner"
)

type netImpl struct {
	tincBin     string
	activePeers sync.Map
	events      network.Events
	definition  *network.Network

	stop func()
	done chan struct{}
	err  error
}

func (impl *netImpl) initAndStart(global context.Context) error {
	if err := impl.definition.Prepare(global, impl.tincBin); err != nil {
		return fmt.Errorf("configure: %w", err)
	}
	if err := internal.Preload(global); err != nil {
		return fmt.Errorf("preload: %w", err)
	}
	self, config, err := impl.Definition().SelfConfig()
	if err != nil {
		return err
	}

	absDir, err := filepath.Abs(impl.definition.Root)
	if err != nil {
		return err
	}

	interfaceName := config.Interface
	if interfaceName == "" { // for darwin
		interfaceName = config.Device[strings.LastIndex(config.Device, "/")+1:]
	}

	ctx, cancel := context.WithCancel(global)
	impl.stop = cancel
	impl.done = make(chan struct{})
	go func() {
		defer cancel()
		defer impl.events.Stopped.Emit(network.NetworkID{Name: impl.definition.Name()})
		defer close(impl.done)
		impl.err = impl.run(absDir, self, ctx)
		impl.activePeers = sync.Map{}
	}()
	return nil
}

func (impl *netImpl) Events() *network.Events {
	return &impl.events
}

func (impl *netImpl) Stop() {
	impl.stop()
}

func (impl *netImpl) Done() <-chan struct{} {
	return impl.done
}

func (impl *netImpl) Error() error {
	return impl.err
}

func (impl *netImpl) Peers() []string {
	var ans []string
	impl.activePeers.Range(func(key, value interface{}) bool {
		ans = append(ans, key.(string))
		return true
	})
	sort.Strings(ans)
	return ans
}

func (impl *netImpl) IsActive(node string) bool {
	_, ok := impl.activePeers.Load(node)
	return ok
}

func (impl *netImpl) Definition() *network.Network {
	return impl.definition
}

func (impl *netImpl) IsRunning() bool {
	ch := impl.done
	if ch == nil {
		return false
	}
	select {
	case <-ch:
		return false
	default:
		return true
	}
}

func (impl *netImpl) run(absDir string, self *network.Node, global context.Context) error {
	ctx, abort := context.WithCancel(global)
	defer abort()

	var wg sync.WaitGroup

	// run tinc service
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer abort()

		for event := range runner.RunTinc(global, impl.tincBin, absDir) {
			if event.Add {
				impl.activePeers.Store(event.Peer.Node, event)
			} else {
				impl.activePeers.Delete(event.Peer.Node)
			}
			log.Printf("%+v", event)
		}

	}()

	// run http API
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer abort()
		server := &localApiServer{definition: impl.definition}
		for {
			err := apiserver.RunHTTP(ctx, "tcp", self.IP+":"+strconv.Itoa(CommunicationPort), server)
			log.Println(impl.definition.Name(), "api stopped:", err)
			select {
			case <-ctx.Done():
				return
			case <-time.After(time.Second):
				log.Println("trying again...")
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := impl.greetEveryone(ctx, *self, GreetInterval)
		if err != nil {
			log.Println("greeting failed:", err)
		}
	}()

	// fix: change owner of log file and pid file to process runner
	wg.Add(1)
	go func() {
		defer wg.Done()
		select {
		case <-ctx.Done():
		case <-time.After(2 * time.Second):
			_ = network.ApplyOwnerOfSudoUser(impl.definition.Pidfile())
		}
	}()

	impl.activePeers.Store(self.Name, self)
	wg.Wait()
	return ctx.Err()
}

func (impl *netImpl) greetEveryone(ctx context.Context, self network.Node, retryInterval time.Duration) error {
	var wg sync.WaitGroup

	nodes, err := impl.Definition().NodesDefinitions()
	if err != nil {
		return err
	}

	for _, node := range nodes {
		if node.IP == "" {
			log.Println("will not greet", node.Name, "'cause it is relay node")
			continue
		}
		wg.Add(1)
		go func(node network.Node) {
			defer wg.Done()

			var client = apiclient.APIClient{BaseURL: "http://" + node.IP + ":" + strconv.Itoa(CommunicationPort)}
			for {
				toImport, err := client.Exchange(ctx, self)
				if err != nil {
					log.Println("greet", node.Name, err)
					goto SLEEP
				}
				for _, node := range toImport {
					err := impl.Definition().Put(&node)
					if err != nil {
						log.Println(node.Name, "import", node.Name, ":", err)
					}
				}
				log.Println("greeted", node.Name)
				break
			SLEEP:
				select {
				case <-ctx.Done():
					return
				case <-time.After(retryInterval):

				}
			}

		}(node)
	}
	wg.Wait()
	return nil
}

type localApiServer struct {
	definition *network.Network
}

func (impl *localApiServer) Exchange(ctx context.Context, remote network.Node) ([]network.Node, error) {
	err := impl.definition.Put(&remote)
	if err != nil {
		return nil, err
	}
	return impl.definition.NodesDefinitions()
}

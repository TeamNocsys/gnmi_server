package mdns

import (
    "context"
    "net"
    "strings"
    "time"
    "sync"

    "github.com/miekg/dns"
    "golang.org/x/net/ipv4"
    "golang.org/x/net/ipv6"

    "github.com/vishvananda/netlink"
    "github.com/sirupsen/logrus"
)

// IPType specifies the IP traffic the client listens for.
// This does not guarantee that only mDNS entries of this sepcific
// type passes. E.g. typical mDNS packets distributed via IPv4, often contain
// both DNS A and AAAA entries.
type IPType uint8

// Options for IPType.
const (
    IPv4        = 0x01
    IPv6        = 0x02
    IPv4AndIPv6 = (IPv4 | IPv6) //< Default option.

    MAX_NBR_OF_ENTRY = 100
)

type clientOpts struct {
    listenOn IPType
    ifaces   []net.Interface
}

// ClientOption fills the option struct to configure intefaces, etc.
type ClientOption func(*clientOpts)

// SelectIPTraffic selects the type of IP packets (IPv4, IPv6, or both) this
// instance listens for.
// This does not guarantee that only mDNS entries of this sepcific
// type passes. E.g. typical mDNS packets distributed via IPv4, may contain
// both DNS A and AAAA entries.
func SelectIPTraffic(t IPType) ClientOption {
    return func(o *clientOpts) {
        o.listenOn = t
    }
}

// SelectIfaces selects the interfaces to query for mDNS records
func SelectIfaces(ifaces []net.Interface) ClientOption {
    return func(o *clientOpts) {
        o.ifaces = ifaces
    }
}

// to record the IP and Host mappings
type HostEntry struct {
    HostName string   `json:"hostname"` // Host machine DNS name
    TTL      uint32   `json:"ttl"`      // TTL of the service record
    Expiration time.Time                // expiration time
}

// Resolver acts as entry point for service lookups and to browse the DNS-SD.
type Resolver struct {
    ipv4conn    *ipv4.PacketConn
    ipv6conn    *ipv6.PacketConn
    wg          *sync.WaitGroup
    Entries     *sync.Map
    eCount      uint32
    expireFlag  chan int

    nlChan       chan struct{}           // channel to close the netlink listener
    updChanAddr  chan netlink.AddrUpdate // channel to monitor for addr events
}

// NewResolver creates a new resolver and joins the UDP multicast groups to
// listen for mDNS messages.
func NewResolver(options ...ClientOption) (*Resolver, error) {
    // Apply default configuration and load supplied options.
    var conf = clientOpts{
        listenOn: IPv4AndIPv6,
    }
    for _, o := range options {
        if o != nil {
            o(&conf)
        }
    }

    aIfs := listAvailInterfaces()

    // IPv4 interfaces
    var ipv4conn *ipv4.PacketConn
    if (conf.listenOn & IPv4) > 0 {
        var err error
        ipv4conn, err = setupUdp4Conn()
        if err == nil {
            joinUdpMulticast(ipv4conn, aIfs, true)
        }
    }

    // IPv6 interfaces
    var ipv6conn *ipv6.PacketConn
    if (conf.listenOn & IPv6) > 0 {
        var err error
        ipv6conn, err = setupUdp6Conn()
        if err == nil {
            joinUdpMulticast(ipv6conn, aIfs, true)
        }
    }

    return &Resolver{
        ipv4conn:   ipv4conn,
        ipv6conn:   ipv6conn,
        wg:         &sync.WaitGroup{},
        Entries:    &sync.Map{},
        eCount :    0,
        expireFlag: make (chan int, 1),
    }, nil
}

func (r *Resolver) Browse(ctx context.Context) {
    r.nlChan = make(chan struct{})
    r.updChanAddr = make(chan netlink.AddrUpdate)

    if err := netlink.AddrSubscribe(r.updChanAddr, r.nlChan); err != nil {
        logrus.Errorf("Error listening on netlink: %v", err)
        return
    }

    go r.mainloop(ctx)
}

// wait for all go routines exit
func (r *Resolver) Stop() {
    r.wg.Wait()
}

// process received name record
func (r *Resolver) processNameRecord(ipstr, hostname string, ttl uint32) {

    if _, ok := r.Entries.Load(ipstr); !ok {
        r.removeExpired(false)

        if r.eCount >= MAX_NBR_OF_ENTRY {
            return
        }

        r.eCount++
    }

    tmp_h := HostEntry {
        HostName  : strings.Replace(hostname, ".local.", "", 1),
        TTL       : ttl,
        Expiration: time.Now().Add(time.Second * time.Duration(ttl)),
    }

    r.Entries.Store(ipstr, tmp_h)
}

// Start listeners and waits for the shutdown signal from exit channel
func (r *Resolver) mainloop(ctx context.Context) {
    defer r.wg.Done()

    r.wg.Add(1)

    // start listening for responses
    msgCh := make(chan *dns.Msg, 32)
    if r.ipv4conn != nil {
        go r.recv(ctx, r.ipv4conn, msgCh)
    }
    if r.ipv6conn != nil {
        go r.recv(ctx, r.ipv6conn, msgCh)
    }

    // Iterate through channels from listeners goroutines
    for {
        select {
        case update := <-r.updChanAddr:
            r.handleAddrUpdate(&update)

        case <-r.nlChan:
            // should not happen
            logrus.Errorf("Stop listening for netlink events")

        case <-r.expireFlag:
            r.removeExpired(true)

        case <-ctx.Done():
            // Context expired
            r.shutdown()
            return

        case msg := <-msgCh:
            sections := append(msg.Answer, msg.Ns...)
            sections = append(sections, msg.Extra...)

            for _, answer := range sections {
                switch rr := answer.(type) {
/* not support now
                case *dns.PTR:

                    if _, ok := entries[rr.Ptr]; !ok {
                        entries[rr.Ptr] = NewServiceEntry(
                            trimDot(strings.Replace(rr.Ptr, rr.Hdr.Name, "", -1)),
                            params.Service,
                            params.Domain)
                    }
                    entries[rr.Ptr].TTL = rr.Hdr.Ttl
                case *dns.SRV:

                    if _, ok := entries[rr.Hdr.Name]; !ok {
                        entries[rr.Hdr.Name] = NewServiceEntry(
                            trimDot(strings.Replace(rr.Hdr.Name, params.ServiceName(), "", 1)),
                            params.Service,
                            params.Domain)
                    }
                    entries[rr.Hdr.Name].HostName = rr.Target
                    entries[rr.Hdr.Name].Port = int(rr.Port)
                    entries[rr.Hdr.Name].TTL = rr.Hdr.Ttl
                case *dns.TXT:

                    if _, ok := entries[rr.Hdr.Name]; !ok {
                        entries[rr.Hdr.Name] = NewServiceEntry(
                            trimDot(strings.Replace(rr.Hdr.Name, params.ServiceName(), "", 1)),
                            params.Service,
                            params.Domain)
                    }
                    entries[rr.Hdr.Name].Text = rr.Txt
                    entries[rr.Hdr.Name].TTL = rr.Hdr.Ttl
*/
                case *dns.A:
                    r.processNameRecord(rr.A.String(), rr.Hdr.Name, rr.Hdr.Ttl)

                case *dns.AAAA:
                    r.processNameRecord(rr.AAAA.String(), rr.Hdr.Name, rr.Hdr.Ttl)
                }
            }
        }
    }
}

// need to join the mc group when address is configured on net interface
func (r *Resolver) handleAddrUpdate(update *netlink.AddrUpdate) {
    logrus.Tracef("addr update - %#v", *update)

    isIpv6 := update.LinkAddress.IP.To4() == nil
    if ifi, err := net.InterfaceByIndex(update.LinkIndex); err == nil {
        if !isInterfaceOk (ifi, false) {
            return
        }

        if !update.NewAddr && isInterfaceHasIp(ifi, !isIpv6) {
            // interfaces still has ip address, do nothing
            return
        }

        if isIpv6 {
            if r.ipv6conn != nil {
                joinUdpMulticast(r.ipv6conn, []net.Interface { *ifi }, update.NewAddr)
            }
        } else {
            if r.ipv4conn != nil {
                joinUdpMulticast(r.ipv4conn, []net.Interface { *ifi }, update.NewAddr)
            }
        }
    }
}

// remove expired entries
func (r *Resolver) removeExpired(isForced bool) {
    if r.eCount >= MAX_NBR_OF_ENTRY || isForced == true {
        curTime := time.Now()
        r.Entries.Range(func(key, value interface{}) bool {
            if curTime.After(value.(HostEntry).Expiration) {
                r.Entries.Delete(key)
                r.eCount --
            }
            return true
        })
    }
}

// send notification to age out expired entries
func (r *Resolver) NotifyExpireEntries () {
    if len (r.expireFlag) == 0 {
       r.expireFlag <- 1
    }
}

// Shutdown client will close currently open connections and channel implicitly.
func (r *Resolver) shutdown() {
    if r.ipv4conn != nil {
        r.ipv4conn.Close()
    }
    if r.ipv6conn != nil {
        r.ipv6conn.Close()
    }
}

// Data receiving routine reads from connection, unpacks packets into dns.Msg
// structures and sends them to a given msgCh channel
func (r *Resolver) recv(ctx context.Context, l interface{}, msgCh chan *dns.Msg) {
    var readFrom func([]byte) (n int, src net.Addr, err error)
    defer r.wg.Done()

    r.wg.Add(1)
    switch pConn := l.(type) {
    case *ipv6.PacketConn:
        readFrom = func(b []byte) (n int, src net.Addr, err error) {
            n, _, src, err = pConn.ReadFrom(b)
            return
        }
    case *ipv4.PacketConn:
        readFrom = func(b []byte) (n int, src net.Addr, err error) {
            n, _, src, err = pConn.ReadFrom(b)
            return
        }
    default:
        return
    }

    buf := make([]byte, 65536)
    var fatalErr error
    for {
        // Handles the following cases:
        // - ReadFrom aborts with error due to closed UDP connection -> causes ctx cancel
        // - ReadFrom aborts otherwise.
        // TODO: the context check can be removed. Verify!
        if ctx.Err() != nil || fatalErr != nil {
            return
        }

        n, _, err := readFrom(buf)
        if err != nil {
            fatalErr = err
            continue
        }
        msg := new(dns.Msg)
        if err := msg.Unpack(buf[:n]); err != nil {
            // log.Printf("[WARN] mdns: Failed to unpack packet: %v", err)
            continue
        }

        select {
        case msgCh <- msg:
            // Submit decoded DNS message and continue.
        case <-ctx.Done():
            // Abort.
            return
        }
    }
}

var MdnsResolver *Resolver;

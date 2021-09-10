package mdns

import (
    "fmt"
    "net"
    "strings"

    "golang.org/x/net/ipv4"
    "golang.org/x/net/ipv6"
    "github.com/sirupsen/logrus"
)

var (
    // Multicast groups used by mDNS
    mdnsGroupIPv4 = net.IPv4(224, 0, 0, 251)
    mdnsGroupIPv6 = net.ParseIP("ff02::fb")

    // mDNS wildcard addresses
    mdnsWildcardAddrIPv4 = &net.UDPAddr{
        IP:   net.ParseIP("224.0.0.0"),
        Port: 5353,
    }
    mdnsWildcardAddrIPv6 = &net.UDPAddr{
        IP: net.ParseIP("ff02::"),
        Port: 5353,
    }

    // mDNS endpoint addresses
    ipv4Addr = &net.UDPAddr{
        IP:   mdnsGroupIPv4,
        Port: 5353,
    }
    ipv6Addr = &net.UDPAddr{
        IP:   mdnsGroupIPv6,
        Port: 5353,
    }
)

func setupUdp6Conn() (*ipv6.PacketConn, error) {
    udpConn, err := net.ListenUDP("udp6", mdnsWildcardAddrIPv6)
    if err != nil {
        return nil, err
    }

    pkConn := ipv6.NewPacketConn(udpConn)
    pkConn.SetControlMessage(ipv6.FlagInterface, true)

    return pkConn, nil
}

func setupUdp4Conn() (*ipv4.PacketConn, error) {
    udpConn, err := net.ListenUDP("udp4", mdnsWildcardAddrIPv4)
    if err != nil {
        // log.Printf("Failed to bind to udp4 mutlicast: %v", err)
        return nil, err
    }

    pkConn := ipv4.NewPacketConn(udpConn)
    pkConn.SetControlMessage(ipv4.FlagInterface, true)

    return pkConn, nil
}

func joinUdpMulticast(l interface{}, interfaces []net.Interface, isJoin bool) error {
    var failedJoins int
    var joinGroup func(tif *net.Interface) error
    var udpTag string = "udp4"

    joinTag := map[bool]string {true: "Join", false: "Leave"} [isJoin]

    switch pConn := l.(type) {
    case *ipv4.PacketConn:
        joinGroup = func(tif *net.Interface) error {
            if isJoin {
                return pConn.JoinGroup(tif, &net.UDPAddr{IP: mdnsGroupIPv4})
            } else {
                return pConn.LeaveGroup(tif, &net.UDPAddr{IP: mdnsGroupIPv4})
            }
        }
    case *ipv6.PacketConn:
        udpTag = "udp6"
        joinGroup = func(tif *net.Interface) error {
            if isJoin {
                return pConn.JoinGroup(tif, &net.UDPAddr{IP: mdnsGroupIPv6})
            } else {
                return pConn.LeaveGroup(tif, &net.UDPAddr{IP: mdnsGroupIPv6})
            }
        }

    default:
        return fmt.Errorf("unknown conn type - %v", l)
    }

    for _, iface := range interfaces {
        err := joinGroup(&iface)

        if err != nil {
            logrus.Tracef("%s: %sGroup failed for iface %s", udpTag, joinTag, iface.Name)
            failedJoins++
        } else {
            logrus.Tracef("%s: %sGroup for iface %s", udpTag, joinTag, iface.Name)
        }
    }

    if failedJoins == len(interfaces) {
        return fmt.Errorf("%s: failed to %s any of these interfaces: %v",
                          udpTag, joinTag, interfaces)
    }

    return nil
}

// check if interface is good to do jobs
func isInterfaceOk (tIf *net.Interface, chkIp bool) bool {
    if !strings.HasPrefix(tIf.Name, "Ethernet") &&
       !strings.HasPrefix(tIf.Name, "Vlan") &&
       !strings.HasPrefix(tIf.Name, "PortChannel") {
        return false
    }

    if chkIp {
        // need ip address to join mc group
        if addrs, err := tIf.Addrs(); err != nil || len(addrs) == 0 {
            return false
        }
    }

    return true
}

// check if interface has ipv4 or ipv6 address
func isInterfaceHasIp (tIf *net.Interface, isv4 bool) bool {
    addrs, err := tIf.Addrs()
    if err != nil || len(addrs) == 0 {
        return false
    }

    for _, addr := range addrs {
        if addr.(*net.IPNet).IP.To4() != nil {
            if isv4 {
                return true
            }
        } else {
            if !isv4 {
                return true
            }
        }
    }

    return false
}

func listAvailInterfaces () []net.Interface {
    var interfaces []net.Interface
    ifaces, err := net.Interfaces()
    if err != nil {
        return nil
    }
    for _, ifi := range ifaces {
        if !isInterfaceOk(&ifi, true) {
            continue
        }

        if (ifi.Flags & net.FlagMulticast) > 0 {
            interfaces = append(interfaces, ifi)
        }
    }

    return interfaces
}


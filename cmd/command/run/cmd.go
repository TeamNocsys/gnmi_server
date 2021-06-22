package run

import (
    "fmt"
    rotatelogs "github.com/lestrrat-go/file-rotatelogs"
    gpb "github.com/openconfig/gnmi/proto/gnmi"
    "github.com/sirupsen/logrus"
    "github.com/spf13/cobra"
    "gnmi_server/cmd/command"
    "gnmi_server/pkg/gnmi"
    "gnmi_server/pkg/gnmi/handler/get"
    "gnmi_server/pkg/gnmi/handler/set"
    "google.golang.org/grpc"
    "net"
    "os"
    "os/signal"
    "path"
    "strings"
    "syscall"
    "time"
)

type ConsoleHook struct {
    debug bool
}

func (ch *ConsoleHook) Levels() []logrus.Level {
    return logrus.AllLevels
}

func (ch *ConsoleHook) Fire(entry *logrus.Entry) error {
    msg, err := entry.String()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Unable to read entry, %v", err)
        return err
    }

    if entry.Level < logrus.DebugLevel || ch.debug {
        if _, err := fmt.Print(msg); err != nil {
            fmt.Fprintf(os.Stderr, "Unable to output entry to console, %v", err)
            return err
        }
    }

    return nil
}

type GsFormatter struct {
}

func (gf *GsFormatter) Format(e *logrus.Entry) ([]byte, error) {
    msg := e.Message
    for k,v := range e.Data {
        s, ok := v.(string)
        if !ok {
            s = fmt.Sprint(v)
        }
        msg += " " + k + "=" + s
    }
    return []byte(fmt.Sprintf("%s [%-5s] %s\n", e.Time.Format("2006-01-02 15:04:05"), strings.ToUpper(e.Level.String()), msg)), nil
}

type runOptions struct {
    username string
    password string
    address string
    port    int
    quiet   bool
    verbose bool
    path    string
}

func NewRunCommand(gnmiCli command.Client) *cobra.Command {
    var opts runOptions

    cmd := &cobra.Command{
        Use:   "run",
        Short: "Run GNMI Server",
        Args:  cobra.NoArgs,
        RunE: func(cmd *cobra.Command, args []string) error {
            // 设置日志输出
            base := path.Join(opts.path, "gnmi_server.log")
            w, err := rotatelogs.New(
                base+".%Y%m%d%H%M",
                rotatelogs.WithLinkName(base),
                rotatelogs.WithRotationTime(24*time.Hour),
                rotatelogs.WithRotationCount(3),
            )
            if err != nil {
                return err
            }
            defer w.Close()
            logrus.SetOutput(w)
            logrus.AddHook(&ConsoleHook{debug: !opts.quiet})
            if opts.verbose {
                logrus.SetLevel(logrus.TraceLevel)
            } else {
                logrus.SetLevel(logrus.DebugLevel)
            }
            logrus.SetFormatter(new(GsFormatter))

            listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", opts.address, opts.port))
            if err != nil {
                return err
            }
            grpcServer := grpc.NewServer(grpc.RPCDecompressor(grpc.NewGZIPDecompressor()))
            server := gnmi.DefaultServer(opts.username, opts.password, gnmiCli, get.GetServeMux(), set.SetServeMux())
            gpb.RegisterGNMIServer(grpcServer, &server)
            c := make(chan os.Signal, 1)
            signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
            go func() {
                <-c
                grpcServer.Stop()
            }()
            logrus.Info("Server start")
            err = grpcServer.Serve(listener)
            logrus.Info("Server stop")
            return err
        },
    }

    flags := cmd.Flags()
    flags.StringVar(&opts.username, "username", "", "the username for ssh connect")
    flags.StringVar(&opts.password, "password", "", "the password for ssh connect")
    flags.StringVar(&opts.address, "address", "0.0.0.0", "the ip address for gnmi serve on")
    flags.IntVar(&opts.port, "port", 5002, "the port for gnmi serve on")
    flags.BoolVarP(&opts.quiet, "quiet", "q", false, "whether to print the log below info level to the screen")
    flags.BoolVarP(&opts.verbose, "verbose", "v", false, "whether to print detail information")
    flags.StringVarP(&opts.path, "path", "p", "/var/log", "log file output path")

    return cmd
}

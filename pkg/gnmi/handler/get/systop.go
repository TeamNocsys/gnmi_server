package get

import (
    "io/ioutil"
    "errors"
    "context"
    "math"
    "syscall"
    "strconv"
    "strings"

    "gnmi_server/cmd/command"
    handler_utils "gnmi_server/pkg/gnmi/handler/utils"
    sonicpb "github.com/TeamNocsys/sonicpb/api/protobuf/sonic"
    "github.com/openconfig/gnmi/proto/gnmi"
    "github.com/openconfig/ygot/proto/ywrapper"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
)

func SysTopInfoHandler(
    ctx context.Context, r *gnmi.GetRequest, db command.Client,
) (*gnmi.GetResponse, error) {

    acctonSystop := &sonicpb.AcctonSystemTop{}

    err := getCpuInfo(ctx, acctonSystop)
    if err != nil {
        return nil, status.Errorf(codes.Internal, err.Error())
    }

    err = getMemInfo(ctx, acctonSystop)
    if err != nil {
        return nil, status.Errorf(codes.Internal, err.Error())
    }

    err = getDiskInfo(ctx, acctonSystop)
    if err != nil {
        return nil, status.Errorf(codes.Internal, err.Error())
    }

    response, err := handler_utils.CreateGetResponse(ctx, r, acctonSystop)
    if err != nil {
        return nil, status.Errorf(codes.Internal, err.Error())
    }

    return response, nil
}

func SysTopInfoCpuHandler(
    ctx context.Context, r *gnmi.GetRequest, db command.Client,
) (*gnmi.GetResponse, error) {

    acctonSystop := &sonicpb.AcctonSystemTop{}

    err := getCpuInfo(ctx, acctonSystop)
    if err != nil {
        return nil, status.Errorf(codes.Internal, err.Error())
    }

    response, err := handler_utils.CreateGetResponse(ctx, r, acctonSystop)
    if err != nil {
        return nil, status.Errorf(codes.Internal, err.Error())
    }

    return response, nil
}

func SysTopInfoMemHandler(
    ctx context.Context, r *gnmi.GetRequest, db command.Client,
) (*gnmi.GetResponse, error) {

    acctonSystop := &sonicpb.AcctonSystemTop{}

    err := getMemInfo(ctx, acctonSystop)
    if err != nil {
        return nil, status.Errorf(codes.Internal, err.Error())
    }

    response, err := handler_utils.CreateGetResponse(ctx, r, acctonSystop)
    if err != nil {
        return nil, status.Errorf(codes.Internal, err.Error())
    }

    return response, nil
}

func SysTopInfoDiskHandler(
    ctx context.Context, r *gnmi.GetRequest, db command.Client,
) (*gnmi.GetResponse, error) {

    acctonSystop := &sonicpb.AcctonSystemTop{}

    err := getDiskInfo(ctx, acctonSystop)
    if err != nil {
        return nil, status.Errorf(codes.Internal, err.Error())
    }

    response, err := handler_utils.CreateGetResponse(ctx, r, acctonSystop)
    if err != nil {
        return nil, status.Errorf(codes.Internal, err.Error())
    }

    return response, nil
}

func addOneCpuInfo(
    systop *sonicpb.AcctonSystemTop, cpu_data []uint64,
) {
    cpu_pcnt := make([]uint64, len(cpu_data))

    /* cpu_pcnt[0] = total, put caculated percentage in cpu_pcnt */
    cpu_pcnt[0] = cpu_data[1] + cpu_data[2] + cpu_data[3] + cpu_data[4] +
                  cpu_data[5] + cpu_data[6] + cpu_data[7] + cpu_data[8]

    for idx :=1; idx <= 7; idx++ {
        cpu_pcnt[idx] = uint64(math.Round(float64(cpu_data[idx]) / float64(cpu_pcnt[0]) * 100))
    }

    tmp_cpu := &sonicpb.AcctonSystemTop_Cpus_CpuKey{
        Index : &sonicpb.AcctonSystemTop_Cpus_CpuKey_IndexUint64{
            IndexUint64 : cpu_data[0],
        },
        Cpu : &sonicpb.AcctonSystemTop_Cpus_Cpu{
            State : &sonicpb.AcctonSystemTop_Cpus_Cpu_State{
                Index : &sonicpb.AcctonSystemTop_Cpus_Cpu_State_IndexUint64{
                    IndexUint64 : cpu_data[0],
                },

                User : &sonicpb.AcctonSystemTop_Cpus_Cpu_State_User{
                    Instant : &ywrapper.UintValue{Value : cpu_pcnt[1]},
                },

                Nice : &sonicpb.AcctonSystemTop_Cpus_Cpu_State_Nice{
                    Instant : &ywrapper.UintValue{Value : cpu_pcnt[2]},
                },

                Kernel : &sonicpb.AcctonSystemTop_Cpus_Cpu_State_Kernel{
                    Instant : &ywrapper.UintValue{Value : cpu_pcnt[3]},
                },

                Idle : &sonicpb.AcctonSystemTop_Cpus_Cpu_State_Idle{
                    Instant : &ywrapper.UintValue{Value : cpu_pcnt[4]},
                },

                Wait : &sonicpb.AcctonSystemTop_Cpus_Cpu_State_Wait{
                    Instant : &ywrapper.UintValue{Value : cpu_pcnt[5]},
                },

                HardwareInterrupt : &sonicpb.AcctonSystemTop_Cpus_Cpu_State_HardwareInterrupt{
                    Instant : &ywrapper.UintValue{Value : cpu_pcnt[6]},
                },

                SoftwareInterrupt : &sonicpb.AcctonSystemTop_Cpus_Cpu_State_SoftwareInterrupt{
                    Instant : &ywrapper.UintValue{Value : cpu_pcnt[7]},
                },

                Total : &sonicpb.AcctonSystemTop_Cpus_Cpu_State_Total{
                    Instant : &ywrapper.UintValue{Value : 100},
                },
            },
        },
    }

    systop.Cpus.Cpu = append(systop.Cpus.Cpu, tmp_cpu)
}

func getCpuInfo(
    ctx context.Context, systop *sonicpb.AcctonSystemTop,
) error {

    /* yang        /proc/stat
     *             cpu[id]       0
     * user     => user,         1
     * nice     => nice,         2
     * kernel   => system,       3
     * idle     => idle,         4
     * wait     => iowait,       5
     * hw-intr  => irq,          6
     * sw-intr  => softirq,      7
     *             steal,        8
     *             guest,        9
     *             guest_nice,  10
     * total    => sum(0~7)
     */
    data, err := ioutil.ReadFile("/proc/stat")
    if err != nil {
        return err
    }

    systop.Cpus = &sonicpb.AcctonSystemTop_Cpus{}
    lines := strings.Split(string(data), "\n")
    cinfo_cnt := 0
    for id, line := range lines {
        fields := strings.Fields(line)
        if len(fields) == 0 || len(fields) < 9 {
            continue
        }
        if fields[0][:3] == "cpu" {
            cpu_data := make([]uint64, len(fields))
            cpu_data[0] = uint64(id)
            for i := 1; i < len(fields); i++ {
                v, _ := strconv.ParseUint(fields[i], 10, 64)
                cpu_data [i] = v
            }

            addOneCpuInfo(systop, cpu_data)
            cinfo_cnt += 1
        } else {
            break
        }
    }

    if cinfo_cnt == 0 {
        return errors.New("Failed to get CpuInfo")
    }

    return nil
}

func getMemInfo(
    ctx context.Context, systop *sonicpb.AcctonSystemTop,
) error {

    /* yang          /proc/meminfo
     * Buffered   => Buffers
     * Cached     => Cached
     * Free       => MemFree
     * SlabUnrecl => SUnreclaim
     * Used       => MemTotal - MemAvailable
     *
     * for free command:
     * Used       => MemTotal - MemFree - Buffers - Cached - SReclaimable
     */
    mem_data := map[string]uint64 {
        "MemTotal"      : 0,
        "MemAvailable"  : 0,
        "Cached"        : 0,
        "Buffers"       : 0,
        "SUnreclaim"    : 0,
        "MemFree"       : 0,
    }

    data, err := ioutil.ReadFile("/proc/meminfo")
    if err != nil {
        return err
    }

    lines := strings.Split(string(data), "\n")
    mem_data_cnt := 0

    for _, line := range lines {
        fields := strings.SplitN(line, ":", 2)
        if len(fields) < 2 {
            continue
        }
        if _, ok := mem_data [fields[0]]; ok {
            valFields := strings.Fields(fields[1])
            val, _ := strconv.ParseUint(valFields[0], 10, 64)
            mem_data[fields[0]] = val
            mem_data_cnt += 1
            if mem_data_cnt == len(mem_data) {
                break
            }
        }
    }

    if mem_data_cnt != len(mem_data) {
        return errors.New("Failed to get MemInfo")
    }

    systop.Memory = &sonicpb.AcctonSystemTop_Memory{
        Buffered   : &ywrapper.UintValue{Value : mem_data["Buffers"]},
        Cached     : &ywrapper.UintValue{Value : mem_data["Cached"]},
        Free       : &ywrapper.UintValue{Value : mem_data["MemFree"]},
        SlabUnrecl : &ywrapper.UintValue{Value : mem_data["SUnreclaim"]},
        Used       : &ywrapper.UintValue{Value : mem_data["MemTotal"] - mem_data["MemAvailable"]},
    }

    return nil
}

func getDiskInfo(
    ctx context.Context, systop *sonicpb.AcctonSystemTop,
) error {

    fs := syscall.Statfs_t{}
    err := syscall.Statfs("/", &fs)
    if err != nil {
        return err
    }

    /* for df, available is caculated by Bavail
     *         (Free blocks available to unprivileged user)
     */
    dfree := fs.Bfree * uint64(fs.Bsize) / 1024
    dall  := fs.Blocks * uint64(fs.Bsize) / 1024
    systop.Disk = &sonicpb.AcctonSystemTop_Disk{
        Free : &ywrapper.UintValue{Value : dfree},
        Used : &ywrapper.UintValue{Value : dall - dfree},
    }

    return nil
}

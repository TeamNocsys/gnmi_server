package get

import (
	"context"
	"encoding/json"
	deepcopy "github.com/getlantern/deepcopy"
	"github.com/golang/glog"
	gpb "github.com/openconfig/gnmi/proto/gnmi"
	"gnmi_server/internal/pkg/openconfig/platform"
	"gnmi_server/internal/pkg/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Fans map[string]*platform.Component_Fan
type Components map[string]*platform.Component

func Get_fan_info(ctx context.Context, r *gpb.GetRequest) (*gpb.GetResponse, error) {
    err, output := utils.Utils_execute_cmd("show", "environment")
    if err != nil {
    	glog.Error("Execute command show environment failed: %v", err)
    	return nil, status.Errorf(codes.Internal, err.Error())
	}

    lines := strings.Split(output, "\n")
    var name string
    step := 1
    fans := make(map[string]*platform.Component_Fan)
	re := regexp.MustCompile("([^:]*):([^(]*)(.*)")

    for idx, line := range lines {
    	if idx == 0 && strings.Contains(line, "Command:") {
    		continue
		}

		switch step {
		case 1:
			name = line
			step = step + 1
		case 2:
			step += 1
		case 3:
			if len(line) == 0 {
				step = 1
			} else {
				matched := re.FindStringSubmatch(line)
				if len(matched) > 0 && strings.Contains(matched[2], "RPM") {
                    subName := name + "_" + strings.Replace(matched[1], " ", "_", -1)
                    values := strings.Split(matched[2], " RPM")
                    if len(values) == 0 {
                    	continue
					}
					speed, err := strconv.ParseUint(strings.Trim(values[0], " \r"), 10, 32)
					if err != nil {
						glog.Errorf("parse int from string failed: %v", err.Error())
						continue
					}
					speedU32 := uint32(speed)
					fan := &platform.Component_Fan{}
					fan.Speed = &speedU32
                    fans[subName] = fan
				}
			}
		}
	}
	str, err := json.Marshal(fans)
	if err != nil {
		glog.Errorf("marshal struct fans failed: %v", err.Error())
		return nil, status.Errorf(codes.Internal, "marshal json failed")
	}

	var prefix gpb.Path
	var path gpb.Path
	err = deepcopy.Copy(&prefix, r.Prefix)
	if err != nil {
		glog.Errorf("deep copy struct path failed: 5v", err.Error())
		return nil, status.Errorf(codes.Internal, "deep copy struct path failed")
	}
	err = deepcopy.Copy(&path, r.Path[0])
	if err != nil {
		glog.Errorf("deep copy struct path failed: %v", err.Error())
		return nil, status.Errorf(codes.Internal, "deep copy struct path failed")
	}
	notification := gpb.Notification{
		Timestamp: time.Now().Unix(),
		Prefix: &prefix,
		Update: []*gpb.Update {&gpb.Update{
		    Path: &path,
		    Val: &gpb.TypedValue{
		    	Value: &gpb.TypedValue_StringVal{
		    		StringVal: string(str),
		    	},
			},
		}},
	}

	response := &gpb.GetResponse{}
	response.Notification = append(response.Notification, &notification)
	return response, nil
}

func Get_temperature_info(ctx context.Context, r *gpb.GetRequest) (*gpb.GetResponse, error) {
    err, output := utils.Utils_execute_cmd("show", "environment")
    if err != nil {
    	glog.Error("Execute command show environment failed: %v", err)
    	return nil, status.Errorf(codes.Internal, err.Error())
	}

	lines := strings.Split(output, "\n")
	var name string
	step := 1
	components := make(map[string]*platform.Component)
	re := regexp.MustCompile("([^:]*):([^(]*)(.*)")

	for idx, line := range lines {
		if idx == 0 && strings.Contains(line, "Command:") {
			continue
		}

		switch step {
		case 1:
			name = line
			step += 1
		case 2:
			step += 1
		case 3:
			if len(line) == 0 {
				step = 1
			} else {
				matched := re.FindStringSubmatch(line)
				if len(matched) > 0 && strings.Contains(matched[2], "C") {
					subName := name + "_" + strings.Replace(matched[1], " ", "_", -1)
					values := strings.Split(matched[2], " C")
					if len(values) == 0 {
						continue
					}
					temp, err := strconv.ParseFloat(strings.Trim(values[0], " \r"), 64)
					if err != nil {
						glog.Errorf("parse float from string failed: %v", err.Error())
						continue
					}
					tempComponent := &platform.Component{
						Temperature: &platform.Component_Temperature{
							AlarmSeverity:  0,
							AlarmStatus:    nil,
							AlarmThreshold: nil,
							Avg:            nil,
							Instant:        &temp,
							Interval:       nil,
							Max:            nil,
							MaxTime:        nil,
							Min:            nil,
							MinTime:        nil,
						},
					}
					components[subName] = tempComponent
				}
			}
		}
	}

	return nil, status.Errorf(codes.Unimplemented, "unimplemented")
}
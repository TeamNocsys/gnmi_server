package get

import (
	"context"
	"encoding/json"
	sonicpb "github.com/TeamNocsys/sonicpb/api/protobuf/sonic"
	"github.com/openconfig/ygot/proto/ywrapper"
	"gnmi_server/internal/pkg/utils"
	"regexp"
	"strconv"
	"strings"
	"time"

	deepcopy "github.com/getlantern/deepcopy"
	"github.com/golang/glog"
	gpb "github.com/openconfig/gnmi/proto/gnmi"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Get_fan_info(ctx context.Context, r *gpb.GetRequest) (*gpb.GetResponse, error) {
	err, output := utils.Utils_execute_cmd("show", "environment")
	if err != nil {
		glog.Errorf("Execute command show environment failed: %s", err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	lines := strings.Split(output, "\n")
	var name string
	step := 1
	re := regexp.MustCompile("([^:]*):([^(]*)(.*)")

	platform := &sonicpb.SonicPlatform_Platform {
		ComponentList: make([]*sonicpb.SonicPlatform_Platform_ComponentListKey, 0, 10),
	}

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
						glog.Errorf("parse int from string failed: %s", err.Error())
						continue
					}
					speedU32 := uint64(speed)

					fan := &sonicpb.SonicPlatform_Platform_ComponentList_Fan {
						State: &sonicpb.SonicPlatform_Platform_ComponentList_Fan_State {
							Speed: &ywrapper.UintValue{
								Value: speedU32,
							},
						},
					}
					componentListKey := &sonicpb.SonicPlatform_Platform_ComponentListKey {
						ComponentName: subName,
						ComponentList: &sonicpb.SonicPlatform_Platform_ComponentList {
							Fan:         fan,
							PowerSupply: nil,
							Temperature: nil,
						},
					}

					platform.ComponentList = append(platform.ComponentList, componentListKey)
				}
			}
		}
	}

	str, err := json.Marshal(platform)
	if err != nil {
		glog.Errorf("marshal struct platform failed: %s", err.Error())
		return nil, status.Errorf(codes.Internal, "marshal json failed")
	}

	var prefix gpb.Path
	var path gpb.Path
	err = deepcopy.Copy(&prefix, r.Prefix)
	if err != nil {
		glog.Errorf("deep copy struct path failed: %s", err.Error())
		return nil, status.Errorf(codes.Internal, "deep copy path struct failed")
	}
	err = deepcopy.Copy(&path, r.Path[0])
	if err != nil {
		glog.Errorf("deep copy struct path failed: %s", err.Error())
		return nil, status.Errorf(codes.Internal, "deep copy struct path failed")
	}
	notification := gpb.Notification{
		Timestamp: time.Now().Unix(),
		Prefix:    &prefix,
		Update: []*gpb.Update{
			&gpb.Update{
				Path: &path,
				Val: &gpb.TypedValue{
					Value: &gpb.TypedValue_BytesVal{
						BytesVal: str,
					},
				},
			},
		},
	}

	response := &gpb.GetResponse{}
	response.Notification = append(response.Notification, &notification)
	return response, nil
}

func Get_temperature_info(ctx context.Context, r *gpb.GetRequest) (*gpb.GetResponse, error) {
	err, output := utils.Utils_execute_cmd("show", "environment")
	if err != nil {
		glog.Errorf("Execute command show environment failed: %s", err.Error())
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	lines := strings.Split(output, "\n")
	var name string
	step := 1
	re := regexp.MustCompile("([^:]*):([^(]*)(.*)")

	platform := &sonicpb.SonicPlatform_Platform {
		ComponentList: make([]*sonicpb.SonicPlatform_Platform_ComponentListKey, 0, 10),
	}

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
					value := strings.Trim(values[0], " \r")
					valueParts := strings.Split(value, ".")

					var intPart int64
					var fracPart int64
					if len(valueParts) == 1 {
						intPart, err = strconv.ParseInt(valueParts[0], 10, 64)
						if err != nil {
							glog.Errorf("parse int64 from string failed: %s", err.Error())
							continue
						}

						fracPart, err = strconv.ParseInt(valueParts[0], 10, 64)
						if err != nil {
							glog.Errorf("parse int64 from string failed: %s", err.Error())
						}
					}

                    temperature := &sonicpb.SonicPlatform_Platform_ComponentList_Temperature {
						Config: nil,
						State:  &sonicpb.SonicPlatform_Platform_ComponentList_Temperature_State {
							Instant: &ywrapper.Decimal64Value{
								Digits:    intPart,
								Precision: uint32(fracPart),
							}},
					}
					componentListKey := &sonicpb.SonicPlatform_Platform_ComponentListKey {
						ComponentName: subName,
						ComponentList: &sonicpb.SonicPlatform_Platform_ComponentList {
							Fan:         nil,
							PowerSupply: nil,
							Temperature: temperature,
						},
					}

				    platform.ComponentList = append(platform.ComponentList, componentListKey)
				}
			}
		}
	}

	str, err := json.Marshal(platform)
	if err != nil {
		glog.Errorf("marshal struct components failed: %s", err.Error())
		return nil, status.Errorf(codes.Internal, "marshal json failed")
	}

	var prefix gpb.Path
	var path gpb.Path
	err = deepcopy.Copy(&prefix, r.Prefix)
	if err != nil {
		glog.Errorf("deep copy path struct failed: %s", err.Error())
		return nil, status.Errorf(codes.Internal, "deep copy path struct failed")
	}
	err = deepcopy.Copy(&path, r.Path[0])
	if err != nil {
		glog.Errorf("deep copy path struct failed: %s", err.Error())
		return nil, status.Errorf(codes.Internal, "deep copy path struct failed")
	}
	notification := gpb.Notification{
		Timestamp: time.Now().Unix(),
		Prefix:    &prefix,
		Update: []*gpb.Update{
			&gpb.Update{
				Path: &path,
				Val: &gpb.TypedValue{
					Value: &gpb.TypedValue_BytesVal{
						BytesVal: str,
					},
				},
			},
		},
	}

	response := &gpb.GetResponse{}
	response.Notification = append(response.Notification, &notification)
	return response, nil
}

func Get_psu_info(ctx context.Context, r *gpb.GetRequest) (*gpb.GetResponse, error) {
	err, output := utils.Utils_execute_cmd("show", "platform", "psustatus")
	if err != nil {
		glog.Errorf("Execute command show environment failed: %s", err.Error())
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	lines := strings.Split(output, "\n")
	var name string

	platform := &sonicpb.SonicPlatform_Platform {
		ComponentList: make([]*sonicpb.SonicPlatform_Platform_ComponentListKey, 0, 10),
	}

	start := false
	for _, line := range lines {
		if line == "" {
			continue
		}

		if start {
			words := strings.Split(line, "  ")
			if len(words) != 2 {
				glog.Warning("parse failed whe get psu information")
				continue
			}
			name = strings.Replace(words[0], " ", "_", -1)

			enabled := false
			if words[1] == "OK" {
				enabled = true
			}

			psu := &sonicpb.SonicPlatform_Platform_ComponentList_PowerSupply {
				Config: nil,
				State:  &sonicpb.SonicPlatform_Platform_ComponentList_PowerSupply_State {
					Enabled: &ywrapper.BoolValue { Value: enabled },
				},
			}
			componentListKey := &sonicpb.SonicPlatform_Platform_ComponentListKey {
				ComponentName: name,
				ComponentList: &sonicpb.SonicPlatform_Platform_ComponentList {
					Fan:         nil,
					PowerSupply: psu,
					Temperature: nil,
				},
			}

			platform.ComponentList = append(platform.ComponentList, componentListKey)
		} else {
			if strings.Contains(line, "-----") {
				start = true
			}
		}
	}

	str, err := json.Marshal(platform)
	if err != nil {
		glog.Errorf("marshal struct platform failed: %s", err.Error())
		return nil, status.Errorf(codes.Internal, "marshal json failed")
	}

	var prefix gpb.Path
	var path gpb.Path

	err = deepcopy.Copy(&prefix, r.Prefix)
	if err != nil {
		glog.Errorf("deep copy struct path failed: %s", err.Error())
		return nil, status.Errorf(codes.Internal, "deep copy path struct failed")
	}

	err = deepcopy.Copy(&path, r.Path[0])
	if err != nil {
		glog.Errorf("deep copy struct path failed: %s", err.Error())
		return nil, status.Errorf(codes.Internal, "deep copy struct path failed")
	}
	notification := gpb.Notification {
		Timestamp: time.Now().Unix(),
		Prefix:    &prefix,
		Alias:     "",
		Update:    []*gpb.Update {
			&gpb.Update {
				Path: &path,
				Val: &gpb.TypedValue {
					Value: &gpb.TypedValue_BytesVal{BytesVal: str},
				},
			},
		},
		Delete:    nil,
		Atomic:    false,
	}

	response := &gpb.GetResponse{}
	response.Notification = append(response.Notification, &notification)
	return response, nil
}

package get

import (
	"context"
	"gnmi_server/cmd/command"
	"gnmi_server/internal/pkg/utils"
	handler_utils "gnmi_server/pkg/gnmi/handler/utils"
	"regexp"
	"strconv"
	"strings"

	sonicpb "github.com/TeamNocsys/sonicpb/api/protobuf/sonic"
	"github.com/golang/glog"
	"github.com/openconfig/gnmi/proto/gnmi"
	"github.com/openconfig/ygot/proto/ywrapper"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ComponentInfoHandler(ctx context.Context, r *gnmi.GetRequest, db command.Client) (*gnmi.GetResponse, error) {
	platform := &sonicpb.SonicPlatform_Platform{}

	err := getFanInfo(ctx, platform)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	err = getTemperatureInfo(ctx, platform)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	err = getPowerSupplyInfo(ctx, platform)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	sonicPlatform := sonicpb.SonicPlatform{
		Platform: platform,
	}

	response, err := handler_utils.CreateGetResponse(ctx, r, &sonicPlatform)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return response, nil
}

func FanInfoHandler(ctx context.Context, r *gnmi.GetRequest, db command.Client) (*gnmi.GetResponse, error) {
	platform := &sonicpb.SonicPlatform_Platform{}

	err := getFanInfo(ctx, platform)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	sonicPlatform := sonicpb.SonicPlatform{
		Platform: platform,
	}

	response, err := handler_utils.CreateGetResponse(ctx, r, &sonicPlatform)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return response, nil
}

func TemperatureInfoHandler(ctx context.Context, r *gnmi.GetRequest, db command.Client) (*gnmi.GetResponse, error) {
	platform := &sonicpb.SonicPlatform_Platform{}

	err := getTemperatureInfo(ctx, platform)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	sonicPlatform := sonicpb.SonicPlatform{
		Platform: platform,
	}

	response, err := handler_utils.CreateGetResponse(ctx, r, &sonicPlatform)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return response, nil
}

func PowerSupplyInfoHandler(ctx context.Context, r *gnmi.GetRequest, db command.Client) (*gnmi.GetResponse, error) {
	platform := &sonicpb.SonicPlatform_Platform{}

	err := getPowerSupplyInfo(ctx, platform)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	sonicPlatform := sonicpb.SonicPlatform{
		Platform: platform,
	}

	response, err := handler_utils.CreateGetResponse(ctx, r, &sonicPlatform)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return response, nil
}

func SystemInfoHandler(ctx context.Context, r *gnmi.GetRequest, db command.Client) (*gnmi.GetResponse, error) {
	platform := &sonicpb.SonicPlatform_Platform{}

	err := getSystemInfo(ctx, platform)
	if err != nil {
		return nil, err
	}

	sonicPlatform := sonicpb.SonicPlatform{
		Platform: platform,
	}

	response, err := handler_utils.CreateGetResponse(ctx, r, &sonicPlatform)
	return response, nil
}

func getFanInfo(ctx context.Context, platform *sonicpb.SonicPlatform_Platform) error {
	var name string
	step := 1
	re := regexp.MustCompile("([^:]*):([^(]*)(.*)")

	err, output := utils.Utils_execute_cmd("show", "environment")
	if err != nil {
		logrus.Errorf("Execute command show environment failed: %s", err)
		return err
	}

	lines := strings.Split(output, "\n")

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

					fan := &sonicpb.SonicPlatform_Platform_ComponentList_Fan{
						State: &sonicpb.SonicPlatform_Platform_ComponentList_Fan_State{
							Speed: &ywrapper.UintValue{
								Value: speedU32,
							},
						},
					}

					compType := sonicpb.SonicPlatform_Platform_ComponentList_State_TypeAcctonplatformtypessonichardwarecomponent{
						sonicpb.AcctonPlatformTypesSONICHARDWARECOMPONENT_ACCTONPLATFORMTYPESSONICHARDWARECOMPONENT_FAN,
					}

					componentListKey := &sonicpb.SonicPlatform_Platform_ComponentListKey{
						Name: subName,
						ComponentList: &sonicpb.SonicPlatform_Platform_ComponentList{
							Fan:         fan,
							PowerSupply: nil,
							State: &sonicpb.SonicPlatform_Platform_ComponentList_State{
								Type: &compType,
							},
						},
					}

					platform.ComponentList = append(platform.ComponentList, componentListKey)
				}
			}
		}
	}

	return nil
}

func getTemperatureInfo(ctx context.Context, platform *sonicpb.SonicPlatform_Platform) error {
	var name string
	step := 1
	re := regexp.MustCompile("([^:]*):([^(]*)(.*)")

	err, output := utils.Utils_execute_cmd("show", "environment")
	if err != nil {
		logrus.Errorf("Execute command show environment failed: %s", err.Error())
		return err
	}
	lines := strings.Split(output, "\n")

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
					values := strings.Split(strings.Replace(matched[2], "+", " ", -1), " C")
					if len(values) == 0 {
						continue
					}
					value := strings.Trim(values[0], " \r")
					valueParts := strings.Split(value, ".")

					var intPart int64
					var fracPart int64
					if len(valueParts) == 2 {
						intPart, err = strconv.ParseInt(valueParts[0], 10, 64)
						if err != nil {
							glog.Errorf("parse int64 from string failed: %s", err.Error())
							continue
						}

						fracPart, err = strconv.ParseInt(valueParts[1], 10, 64)
						if err != nil {
							glog.Errorf("parse int64 from string failed: %s", err.Error())
							continue
						}
					} else {
						logrus.Error("Can not parse temperature value")
					}

					compType := sonicpb.SonicPlatform_Platform_ComponentList_State_TypeAcctonplatformtypessonichardwarecomponent{
						sonicpb.AcctonPlatformTypesSONICHARDWARECOMPONENT_ACCTONPLATFORMTYPESSONICHARDWARECOMPONENT_SENSOR,
					}

					componentListKey := &sonicpb.SonicPlatform_Platform_ComponentListKey{
						Name: subName,
						ComponentList: &sonicpb.SonicPlatform_Platform_ComponentList{
							Fan:         nil,
							PowerSupply: nil,
							State: &sonicpb.SonicPlatform_Platform_ComponentList_State{
								Temperature: &sonicpb.SonicPlatform_Platform_ComponentList_State_Temperature{
									Instant: &ywrapper.Decimal64Value{
										Digits:    intPart,
										Precision: uint32(fracPart),
									},
								},
								Type: &compType,
							},
						},
					}

					platform.ComponentList = append(platform.ComponentList, componentListKey)
				}
			}
		}
	}

	return nil
}

func getPowerSupplyInfo(ctx context.Context, platform *sonicpb.SonicPlatform_Platform) error {
	var name string

	err, output := utils.Utils_execute_cmd("show", "platform", "psustatus")
	if err != nil {
		logrus.Errorf("Execute command show environment failed: %s", err.Error())
		return err
	}
	lines := strings.Split(output, "\n")

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

			psu := &sonicpb.SonicPlatform_Platform_ComponentList_PowerSupply{
				Config: nil,
				State: &sonicpb.SonicPlatform_Platform_ComponentList_PowerSupply_State{
					Enabled: &ywrapper.BoolValue{Value: enabled},
				},
			}

			compType := sonicpb.SonicPlatform_Platform_ComponentList_State_TypeAcctonplatformtypessonichardwarecomponent{
				sonicpb.AcctonPlatformTypesSONICHARDWARECOMPONENT_ACCTONPLATFORMTYPESSONICHARDWARECOMPONENT_POWER_SUPPLY,
			}

			componentListKey := &sonicpb.SonicPlatform_Platform_ComponentListKey{
				Name: name,
				ComponentList: &sonicpb.SonicPlatform_Platform_ComponentList{
					Fan:         nil,
					PowerSupply: psu,
					State: &sonicpb.SonicPlatform_Platform_ComponentList_State{
						Type: &compType,
					},
				},
			}

			platform.ComponentList = append(platform.ComponentList, componentListKey)
		} else {
			if strings.Contains(line, "-----") {
				start = true
			}
		}
	}

	return nil
}

func getSystemInfo(ctx context.Context, platform *sonicpb.SonicPlatform_Platform) error {
    entries := make(map[string]string)

    // err, output := utils.Utils_execute_cmd("show", "platform", "syseeprom")
    err, output := utils.Utils_execute_cmd("decode-syseeprom", "-d")
    if err != nil {
        // logrus.Errorf("Execute command show platform syseeprom failed: %s", err.Error())
        logrus.Errorf("Execute command decode-syseeprom failed: %s", err.Error())
        return err
    }

    startParse := false
    lines := strings.Split(output, "\n")
    for _, line := range lines {
        if startParse {
            line = strings.TrimSpace(line)
            startIdx := 0

            // get name
            endIdx := strings.Index(line, "0x")
            if endIdx == -1 {
                break
            }
            name := strings.TrimSpace(line[0 : endIdx-1])

            // get value
            startIdx = strings.LastIndex(line, " ")
            if endIdx == -1 {
                logrus.Errorf("parse TLV failed when process show platform syseeprom output line: %s",
                    line)
                continue
            }

            entries[name] = strings.TrimSpace(line[startIdx:])
        } else if strings.HasPrefix(line, "--") {
            startParse = true
        }
    }

    err, output = utils.Utils_execute_cmd("show", "version")
    if err != nil {
        logrus.Errorf("Execute command show version failed: %s", err.Error())
        return err
    }
    lines = strings.Split(output, "\n")
    for _, line := range lines {
        if strings.Contains(line, "Software Version") {
            words := strings.Split(line, ":")
            if len(words) != 2 {
                logrus.Errorf("parse software version failed, origin line: %s", line)
            } else {
                entries["Software Version"] = words[1]
            }
        }
    }



    compType := sonicpb.SonicPlatform_Platform_ComponentList_State_TypeAcctonplatformtypessonicsoftwarecomponent {
        sonicpb.AcctonPlatformTypesSONICSOFTWARECOMPONENT_ACCTONPLATFORMTYPESSONICSOFTWARECOMPONENT_OPERATING_SYSTEM,
    }

    state := &sonicpb.SonicPlatform_Platform_ComponentList_State {
        Type: &compType,
    }

    // Serial No.
    serialNo, ok := entries["Serial Number"]
    if ok {
        state.SerialNo = &ywrapper.StringValue{Value: serialNo}
    } else {
        logrus.Error("Serial No. is not exists")
    }

    // Part Number
    partNumber, ok := entries["Part Number"]
    if ok {
        state.PartNo = &ywrapper.StringValue{Value: partNumber}
    } else {
        logrus.Error("Part Number is not exists")
    }

    // Hardware Version
    hardwareVersion, ok := entries["Platform Name"]
    if ok {
        state.HardwareVersion = &ywrapper.StringValue{Value: hardwareVersion}
    } else {
        logrus.Error("Hardware Version is not exists")
    }

    // Software Version
    softwareVersion, ok := entries["Software Version"]
    if ok {
        state.SoftwareVersion = &ywrapper.StringValue{Value: softwareVersion}
    } else {
        logrus.Error("Software Version is not exists")
    }

    // Manufacturer
    manufacture, ok := entries["Manufacturer"]
    if ok {
        state.MfgName = &ywrapper.StringValue{Value: manufacture}
    } else {
        logrus.Error("Manufacturer is  not exists")
    }

    // Manufacture Date
    manufactureDate, ok := entries["Manufacture Date"]
    if ok {
        state.MfgDate = &ywrapper.StringValue{Value: manufactureDate}
    } else {
        logrus.Error("Manufacturer is not exists")
    }

    component := &sonicpb.SonicPlatform_Platform_ComponentList{}
    component.State = state

    eth0Mac, ok := entries["Base MAC Address"]
    if !ok {
       logrus.Error("eth0's mac is not exists")
    } else {
        state := &sonicpb.SonicPlatform_Platform_ComponentList_Properrties_Property_State{
            Configurable: &ywrapper.BoolValue{Value: true},
            Name:         &ywrapper.StringValue{Value: "ETH0_MAC"},
            Value:        &sonicpb.SonicPlatform_Platform_ComponentList_Properrties_Property_State_ValueString{
                ValueString: eth0Mac,
            },
        }

        config := &sonicpb.SonicPlatform_Platform_ComponentList_Properrties_Property_Config {
            Name: &ywrapper.StringValue{Value: "ETH0_MAC"},
        }

        property := &sonicpb.SonicPlatform_Platform_ComponentList_Properrties_Property{
            State: state,
            Config: config,
        }

        propertyKey := &sonicpb.SonicPlatform_Platform_ComponentList_Properrties_PropertyKey{
            Property: property,
        }

        properties := &sonicpb.SonicPlatform_Platform_ComponentList_Properrties{}
        properties.Property = append(properties.Property, propertyKey)
        component.Properrties = properties
    }

    componentName := "unknown"
    componentName, ok = entries["Product Name"]
    if !ok {
        logrus.Error("Product name is not exists")
    }

    componentListKey := &sonicpb.SonicPlatform_Platform_ComponentListKey{
        Name:          componentName,
        ComponentList: component,
    }

    platform.ComponentList = append(platform.ComponentList, componentListKey)

    return nil
}
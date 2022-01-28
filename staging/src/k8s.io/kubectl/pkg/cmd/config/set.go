/*
Copyright 2014 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package config

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"reflect"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	cliflag "k8s.io/component-base/cli/flag"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"k8s.io/kubectl/pkg/util/i18n"
	"k8s.io/kubectl/pkg/util/templates"
)

type setOptions struct {
	configAccess  clientcmd.ConfigAccess
	propertyName  string
	propertyValue string
	setRawBytes   cliflag.Tristate
}

var (
	setLong = templates.LongDesc(i18n.T(`
	Set an individual value in a kubeconfig file.

	PROPERTY_NAME is a dot delimited name where each token represents either an attribute name or a map key.  Map keys may not contain dots.

	PROPERTY_VALUE is the new value you want to set. Binary fields such as 'certificate-authority-data' expect a base64 encoded string unless the --set-raw-bytes flag is used.

	Specifying an attribute name that already exists will merge new fields on top of existing values.`))

	setExample = templates.Examples(`
	# Set the server field on the my-cluster cluster to https://1.2.3.4
	kubectl config set clusters.my-cluster.server https://1.2.3.4

	# Set the certificate-authority-data field on the my-cluster cluster
	kubectl config set clusters.my-cluster.certificate-authority-data $(echo "cert_data_here" | base64 -i -)

	# Set the cluster field in the my-context context to my-cluster
	kubectl config set contexts.my-context.cluster my-cluster

	# Set the client-key-data field in the cluster-admin user using --set-raw-bytes option
	kubectl config set users.cluster-admin.client-key-data cert_data_here --set-raw-bytes=true`)
)

// NewCmdConfigSet returns a Command instance for 'config set' sub command
func NewCmdConfigSet(out io.Writer, configAccess clientcmd.ConfigAccess) *cobra.Command {
	options := &setOptions{configAccess: configAccess}

	cmd := &cobra.Command{
		Use:                   "set PROPERTY_NAME PROPERTY_VALUE",
		DisableFlagsInUseLine: true,
		Short:                 i18n.T("Set an individual value in a kubeconfig file"),
		Long:                  setLong,
		Example:               setExample,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(options.complete(cmd))
			cmdutil.CheckErr(options.run())
			fmt.Fprintf(out, "Property %q set.\n", options.propertyName)
		},
	}

	f := cmd.Flags().VarPF(&options.setRawBytes, "set-raw-bytes", "", "When writing a []byte PROPERTY_VALUE, write the given string directly without base64 decoding.")
	f.NoOptDefVal = "true"
	return cmd
}

func (o setOptions) run() error {
	err := o.validate()
	if err != nil {
		return err
	}

	config, err := o.configAccess.GetStartingConfig()
	if err != nil {
		return err
	}
	steps, err := newNavigationSteps(o.propertyName)
	if err != nil {
		return err
	}

	setRawBytes := false
	if o.setRawBytes.Provided() {
		setRawBytes = o.setRawBytes.Value()
	}

	err = modifyConfig(reflect.ValueOf(config), steps, o.propertyValue, false, setRawBytes)
	if err != nil {
		return err
	}

	if err := clientcmd.ModifyConfig(o.configAccess, *config, false); err != nil {
		return err
	}

	return nil
}

func (o *setOptions) complete(cmd *cobra.Command) error {
	endingArgs := cmd.Flags().Args()
	if len(endingArgs) != 2 {
		return helpErrorf(cmd, "Unexpected args: %v", endingArgs)
	}

	o.propertyValue = endingArgs[1]
	o.propertyName = endingArgs[0]
	return nil
}

func (o setOptions) validate() error {
	if len(o.propertyValue) == 0 {
		return errors.New("you cannot use set to unset a property")
	}

	if len(o.propertyName) == 0 {
		return errors.New("you must specify a property")
	}

	return nil
}

func modifyConfig(curr reflect.Value, steps *navigationSteps, propertyValue string, unset bool, setRawBytes bool) error {
	currStep := steps.pop()

	actualCurrValue := curr
	if curr.Kind() == reflect.Ptr {
		actualCurrValue = curr.Elem()
	}

	switch actualCurrValue.Kind() {
	case reflect.Map:
		if !steps.moreStepsRemaining() && !unset {
			actualCurrValue.SetMapIndex(reflect.ValueOf(currStep.stepValue), reflect.ValueOf(propertyValue))
			return nil
		}

		mapKey := reflect.ValueOf(currStep.stepValue)
		mapValueType := curr.Type().Elem().Elem()

		if !steps.moreStepsRemaining() && unset {
			actualCurrValue.SetMapIndex(mapKey, reflect.Value{})
			return nil
		}

		currMapValue := actualCurrValue.MapIndex(mapKey)

		needToSetNewMapValue := currMapValue.Kind() == reflect.Invalid
		if needToSetNewMapValue {
			if unset {
				return fmt.Errorf("current map key `%v` is invalid", mapKey.Interface())
			}
			currMapValue = reflect.New(mapValueType.Elem()).Elem().Addr()
			actualCurrValue.SetMapIndex(mapKey, currMapValue)
		}

		err := modifyConfig(currMapValue, steps, propertyValue, unset, setRawBytes)
		if err != nil {
			return err
		}

		return nil

	case reflect.String:
		if steps.moreStepsRemaining() {
			return fmt.Errorf("can't have more steps after a string. %v", steps)
		}
		actualCurrValue.SetString(propertyValue)
		return nil

	case reflect.Slice:
		if steps.moreStepsRemaining() {
			return fmt.Errorf("can't have more steps after slice. %v", steps)
		}

		if setRawBytes {
			actualCurrValue.SetBytes([]byte(propertyValue))
		} else {
			innerKind := actualCurrValue.Type().Elem().Kind()
			if innerKind == reflect.String {
				currentSliceValue := actualCurrValue.Interface().([]string)
				newSliceValue := editStringSlice(currentSliceValue, propertyValue)
				actualCurrValue.Set(reflect.ValueOf(newSliceValue))
			} else if innerKind == reflect.Struct {
				// The only struct slices we should be getting into here are ExecEnvVars
				structName := currStep.stepValue

				// If the existing env var config is empty set it and we aren't unsetting create the new value and return
				if actualCurrValue.IsNil() && !unset {
					newSlice := reflect.MakeSlice(reflect.TypeOf([]clientcmdapi.ExecEnvVar{}), 1, 1)
					newStructValue := newSlice.Index(0)
					newStructValue.Set(reflect.ValueOf(clientcmdapi.ExecEnvVar{
						Name:  structName,
						Value: propertyValue,
					}))
					actualCurrValue.Set(newSlice)
					return nil
				}

				// Find the ExecEnvVar by name, stop if can't find and we're trying to unset because nothing needs to happen
				existingStructIndex := getExecConfigEnvByName(actualCurrValue, structName)
				if existingStructIndex < 0 && unset {
					return nil
				}

				// If we found the array index and we're unsetting then unset and return
				if existingStructIndex >= 0 && unset {
					maxSliceIndex := actualCurrValue.Len()
					firstSlice := actualCurrValue.Slice(0, existingStructIndex)
					secondSlice := actualCurrValue.Slice(existingStructIndex+1, maxSliceIndex)
					for i := 0; i < secondSlice.Len(); i++ {
						firstSlice = reflect.Append(firstSlice, secondSlice.Index(i))
					}
					actualCurrValue.Set(firstSlice)
					return nil
				}

				// If key exists set new value, else create new ExecEnvVar struct with specified values, else create new key/value
				if existingStructIndex >= 0 {
					existingStructVal := actualCurrValue.Index(existingStructIndex)
					existingStructVal.Set(reflect.ValueOf(clientcmdapi.ExecEnvVar{
						Name:  structName,
						Value: propertyValue,
					}))
					return nil
				} else {
					newSlice := reflect.MakeSlice(reflect.TypeOf([]clientcmdapi.ExecEnvVar{}), 1, 1)
					newStructValue := newSlice.Index(0)
					newStructValue.Set(reflect.ValueOf(clientcmdapi.ExecEnvVar{
						Name:  structName,
						Value: propertyValue,
					}))
					actualCurrValue.Set(reflect.Append(actualCurrValue, newStructValue))
					return nil
				}
			} else {
				val, err := base64.StdEncoding.DecodeString(propertyValue)
				if err != nil {
					return fmt.Errorf("error decoding input value: %v", err)
				}
				actualCurrValue.SetBytes(val)
			}
		}
		return nil

	case reflect.Bool:
		if steps.moreStepsRemaining() {
			return fmt.Errorf("can't have more steps after a bool. %v", steps)
		}
		boolValue, err := toBool(propertyValue)
		if err != nil {
			return err
		}
		actualCurrValue.SetBool(boolValue)
		return nil

	case reflect.Struct:
		for fieldIndex := 0; fieldIndex < actualCurrValue.NumField(); fieldIndex++ {
			currFieldValue := actualCurrValue.Field(fieldIndex)
			currFieldType := actualCurrValue.Type().Field(fieldIndex)
			currYamlTag := currFieldType.Tag.Get("json")
			currFieldTypeYamlName := strings.Split(currYamlTag, ",")[0]

			if currFieldTypeYamlName == currStep.stepValue {
				thisMapHasNoValue := (currFieldValue.Kind() == reflect.Map && currFieldValue.IsNil())

				if thisMapHasNoValue {
					newValue := reflect.MakeMap(currFieldValue.Type())
					currFieldValue.Set(newValue)

					if !steps.moreStepsRemaining() && unset {
						return nil
					}
				}

				if !steps.moreStepsRemaining() && unset {
					// if we're supposed to unset the value or if the value is a map that doesn't exist, create a new value and overwrite
					newValue := reflect.New(currFieldValue.Type()).Elem()
					currFieldValue.Set(newValue)
					return nil
				}

				return modifyConfig(currFieldValue.Addr(), steps, propertyValue, unset, setRawBytes)
			}
		}

		return fmt.Errorf("unable to locate path %v", currStep.stepValue)

	case reflect.Ptr:
		newActualCurrValue := actualCurrValue.Elem()
		if actualCurrValue.IsNil() {
			newValue := reflect.New(actualCurrValue.Type().Elem())
			actualCurrValue.Set(newValue)
			newActualCurrValue = actualCurrValue.Elem()
			if !steps.moreStepsRemaining() && unset {
 				return nil
			}
		}
		steps.currentStepIndex = steps.currentStepIndex - 1
		return modifyConfig(newActualCurrValue.Addr(), steps, propertyValue, unset, setRawBytes)

	}

	panic(fmt.Errorf("unrecognized type: %v\nwanted: %v", actualCurrValue, actualCurrValue.Kind()))
}

// getExecConfigEnvByName returns the value
func getExecConfigEnvByName(v reflect.Value, name string) int {
	if v.Type() != reflect.SliceOf(reflect.TypeOf(clientcmdapi.ExecEnvVar{})) {
		return -1
	}

	// Pull slice value out of value object
	slice, ok := v.Interface().([]clientcmdapi.ExecEnvVar)
	if !ok {
		return -1
	}

	// Iterate through slice of ExecEnvVars and check for a matching Name key, return when found
	for i, envVar := range slice {
		if envVar.Name == name {
			return i
		}
	}

	// If we never find the Name key return false
	return -1
}

func dedupeStringSlice(slice []string) []string {
	sliceMap := make(map[string]struct{})
	for i := 0; i < len(slice); i++ {
		sliceMap[slice[i]] = struct{}{}
	}
	var dedupeSlice []string
	for k := range sliceMap {
		dedupeSlice = append(dedupeSlice, k)
	}
	sort.Strings(dedupeSlice)
	return dedupeSlice
}

func editStringSlice(slice []string, input string) []string {
	function := string(input[len(input)-1])
	switch function {
	case "-":
		// Remove an argument
		// Remove last character defining function type
		input = string(input[:len(input)-1])
		argSlice := strings.Split(input, ",")

		for _, arg := range argSlice {
			for j, existingArgs := range slice {
				if existingArgs == arg {
					slice = append(slice[:j], slice[j+1:]...)
					break
				}
			}
		}

		return slice

	case "+":
		// Add new argument
		// Remove last character defining function type
		input = string(input[:len(input)-1])
		argSlice := strings.Split(input, ",")
		slice = append(slice, argSlice...)
		return dedupeStringSlice(slice)

	default:
		argSlice := strings.Split(input, ",")
		return dedupeStringSlice(argSlice)
	}
}

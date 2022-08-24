package utils

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"os/exec"
)

var helmCmd = []string{"helm"}

const chartName = "https://github.com/FFFFFaraway/MPI-Operator/raw/main/charts/mpijob/mpijob-0.1.0.tgz"

func InstallRelease(name string, namespace string, values interface{}) error {
	binary, err := exec.LookPath(helmCmd[0])
	if err != nil {
		return err
	}

	// 1. generate the template file
	valueFile, err := ioutil.TempFile(os.TempDir(), "values")
	if err != nil {
		log.Errorf("Failed to create tmp file %v due to %v", valueFile.Name(), err)
		return err
	} else {
		log.Debugf("Save the values file %s", valueFile.Name())
	}
	// defer os.Remove(valueFile.Name())

	// 2. dump the object into the template file
	err = toYaml(values, valueFile)
	if err != nil {
		return err
	}

	// example: helm install -f /var/folders/dt/tqr6182j6bn5bxpk57m026y80000gn/T/values464191150 --namespace sw test https://github.com/FFFFFaraway/MPI-Operator/raw/main/charts/mpijob/mpijob-0.1.0.tgz
	// 4. prepare the arguments
	args := []string{"install", "-f", valueFile.Name(), "-n", namespace, name, chartName}
	log.Debugf("Exec %s, %v", binary, args)

	// return syscall.Exec(cmd, args, env)
	// 5. execute the command
	cmd := exec.Command(binary, args...)
	out, err := cmd.CombinedOutput()
	fmt.Println("")
	fmt.Printf("%s\n", string(out))
	if err != nil {
		log.Fatalf("Failed to execute %s, %v with %v", binary, args, err)
	}

	// 6. clean up the value file if needed
	if log.GetLevel() != log.DebugLevel {
		err = os.Remove(valueFile.Name())
		if err != nil {
			log.Warnf("Failed to delete %s due to %v", valueFile.Name(), err)
		}
	}

	return nil
}

func DeleteRelease(name, namespace string) error {
	binary, err := exec.LookPath(helmCmd[0])
	if err != nil {
		return err
	}

	args := []string{"uninstall", name, "-n", namespace}

	cmd := exec.Command(binary, args...)
	out, err := cmd.CombinedOutput()
	fmt.Println("")
	fmt.Printf("%s\n", string(out))
	if err != nil {
		log.Fatalf("Failed to execute %s, %v with %v", binary, args, err)
	}
	return nil
}

func toYaml(values interface{}, file *os.File) error {
	log.Debugf("values: %+v", values)
	data, err := yaml.Marshal(values)
	if err != nil {
		log.Errorf("Failed to marshal value %v due to %v", values, err)
		return err
	}

	defer file.Close()
	_, err = file.Write(data)
	if err != nil {
		log.Errorf("Failed to write %v to %s due to %v", data, file.Name(), err)
	}
	return err
}

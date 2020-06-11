/*******************************************************************************
 * IBM Confidential
 * OCO Source Materials
 * IBM Cloud Container Service, 5737-D43
 * (C) Copyright IBM Corp. 2017, 2018 All Rights Reserved.
 * The source code for this program is not  published or otherwise divested of
 * its trade secrets, irrespective of what has been deposited with
 * the U.S. Copyright Office.
 ******************************************************************************/

package main

import (
	"flag"

	cfg "github.com/IBM/ibmcloud-object-storage-plugin/utils/config"
	log "github.com/IBM/ibmcloud-object-storage-plugin/utils/logger"
	watcher "github.com/IBM/ibmcloud-object-storage-plugin/watcher"
	"go.uber.org/zap"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var master = flag.String(
	"master",
	"",
	"Master URL to build a client config from. Either this or kubeconfig needs to be set if the provisioner is being run out of cluster.",
)
var kubeconfig = flag.String(
	"kubeconfig",
	"",
	"Absolute path to the kubeconfig file. Either this or master needs to be set if the provisioner is being run out of cluster.",
)

func main() {
	var err error
	logger, _ := log.GetZapLogger()
	loggerLevel := zap.NewAtomicLevel()
	err = flag.Set("logtostderr", "true")
	if err != nil {
		logger.Info("Failed to set flag:", zap.Error(err))
	}
	flag.Parse()

	// Enable debug trace
	debugTrace := cfg.GetConfigBool("DEBUG_TRACE", false, *logger)
	if debugTrace {
		loggerLevel.SetLevel(zap.DebugLevel)
	}

	var config *rest.Config
	config, err = clientcmd.BuildConfigFromFlags(*master, *kubeconfig)
	if err != nil {
		logger.Fatal("Failed to create config:", zap.Error(err))
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		logger.Fatal("Failed to create client:", zap.Error(err))
	}

	err = cfg.SetUpEvn(clientset, logger)
	if err != nil {
		logger.Fatal("Error while loading the ENV variables", zap.Error(err))
	}

	_, err = clientset.Discovery().ServerVersion()
	if err != nil {
		logger.Fatal("Error getting server version:", zap.Error(err))
	}

	isWatcher := cfg.GetConfigBool("WATCHER", false, *logger)
	if isWatcher {
		// Start watcher for persistent volumes
		watcher.WatchPersistentVolumes(clientset, *logger)
	} else {
		logger.Fatal("Please set 'WATCHER' to true for starting wacther.")
	}
}

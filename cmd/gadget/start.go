package main

import (
	"errors"
	"strings"
	log "github.com/sirupsen/logrus"
)

// Process the build arguments and execute build
func GadgetStart(args []string, g *GadgetContext) error {
	
	EnsureKeys()

	client, err := GadgetLogin(gadgetPrivKeyLocation)

	if err != nil {
		return err
	}
	
	var startFailed bool = false
	
	log.Info("[GADGT]  Starting:")
	stagedContainers,_ := FindStagedContainers(args, append(g.Config.Onboot, g.Config.Services...))
	for _, container := range stagedContainers {
		
		log.Infof("[GADGT]    %s", container.Alias)
		binds := strings.Join( PrependToStrings(container.Binds[:],"-v "), " ")
		commands := strings.Join(container.Command[:]," ")
		
		stdout, stderr, err := RunRemoteCommand(client, "docker create --name", container.Alias, binds, container.ImageAlias, commands)
		
		log.WithFields(log.Fields{
			"function": "GadgetStart",
			"name": container.Alias,
			"start-stage": "create",
		}).Debug(stdout)
		log.WithFields(log.Fields{
			"function": "GadgetStart",
			"name": container.Alias,
			"start-stage": "create",
		}).Debug(stderr)
		
		if err != nil {
			
			// fail loudly, but continue
			
			log.WithFields(log.Fields{
				"function": "GadgetStart",
				"name": container.Alias,
				"start-stage": "create",
			}).Debug("This is likely due to specifying containers for deploying, but trying to start all")


			log.Debug("Failed to create container on Gadget,")
			log.Debug("it might have already been deployed,")
			log.Debug("Or creation otherwise failed")
			
			startFailed = true
		}

		stdout, stderr, err = RunRemoteCommand(client, "docker start", container.Alias)
		
		log.WithFields(log.Fields{
			"function": "GadgetStart",
			"name": container.Alias,
			"start-stage": "create",
		}).Debug(stdout)
		log.WithFields(log.Fields{
			"function": "GadgetStart",
			"name": container.Alias,
			"start-stage": "create",
		}).Debug(stderr)
		
		if err != nil {
			
			// fail loudly, but continue
			
			log.WithFields(log.Fields{
				"function": "GadgetStart",
				"name": container.Alias,
				"start-stage": "create",
			}).Debug("This is likely due to specifying containers for deploying, but trying to start all")


			log.Error("Failed to start container on Gadget")
			log.Warn("Was the container ever deployed?")
			
			startFailed = true
		}

	}
	
	if startFailed {
		err = errors.New("Failed to create or start one or more containers")
	}
	
	return err
}

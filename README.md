This is a library to make the development of a robot based on [gitee](https://gitee.com) simpler.

# The functions of the lib
- [config](https://github.com/opensourceways/community-robot-lib/blob/master/config)

  It is a common component which includes a agent to watch the config file of robot and a configuration item([PluginForRepo](https://github.com/opensourceways/community-robot-lib/blob/master/config/plugin_for_repo.go#L9)) which can restrict the config to a specified organization or a repository.

- [giteeclient](https://github.com/opensourceways/community-robot-lib/blob/master/giteeclient)

  It is a swapper to encapsulate dozens of frequently-used Gitee APIs.

- [giteeplugin](https://github.com/opensourceways/community-robot-lib/blob/master/giteeplugin)

  It is the framework of robot based on gitee. It implements the interfaces to register the event handler and dispatch the event to each handler.
  It is very easy to implement a new robot based on gitee with it.

- [interrupts](https://github.com/opensourceways/community-robot-lib/blob/master/interrupts)

  It is copied from [prow](https://github.com/kubernetes/test-infra/tree/master/prow/interrupts) and implement the function to control the exit of service.

- [logrusutil](https://github.com/opensourceways/community-robot-lib/blob/master/logrusutil)

  It is copied from [prow](https://github.com/kubernetes/test-infra/tree/master/prow/logrusutil) and implement the function to print log.

- [new-plugin](https://github.com/opensourceways/community-robot-lib/blob/master/new-plugin)

  It is the template of a robot based on gitee. The robot uses the Bazel to manage the dependencies.
  The [build.sh](https://github.com/opensourceways/community-robot-lib/blob/master/new-plugin/build.sh) is the useful tool to compile and build image.

- [options](https://github.com/opensourceways/community-robot-lib/blob/master/options)

  It includes the common options for a robot.

- [secret](https://github.com/opensourceways/community-robot-lib/blob/master/secret)

  It is a common component which can watch the secret files of robot.

- [tools](https://github.com/opensourceways/community-robot-lib/blob/master/tools)

  It includes two useful scripts.

  The '**new_robot.sh**' can geneate the initial robot codes by downloading the files in the [new-plugin](https://github.com/opensourceways/community-robot-lib/blob/master/new-plugin).

  The '**deploy_plugin.sh**' is used to deploy robots on the k8s environment and update the image when the robot changes.

- [utils](https://github.com/opensourceways/community-robot-lib/blob/master/utils)

  It includes several useful functions which may be used in robot.

This docker image allows users to use the Timeplus terraform provider to manage Timeplus resources without downloading terraform and the provider plugin.

Usage:

First, create your terraform files for your Timeplus resources. Then run
```bash
docker run --rm -it -v $(pwd):/app ghcr.io/timeplus-io/resources-manager:1.0.0 apply
```

The above command will mount your terraform files on the container and run `terraform apply` to apply the changes to your Timeplus workspace. You can change `apply` to any terraform commands, for example, `plan` to see the change plan, etc.

Note that, terraform by default uses files to store states. If you run the above command, the state file (`terraform.tfstate`) will be created in your current directory (`pwd`). You will need to keep the file in order to make terraform works properly. If you want the state file to be created in other locations, you can either copy the state file to other locations after the terraform command succesfully completes, or use the follow command instead:

```bash
docker run --rm -it -v $(pwd):/app -v /path/to/directory/for/state/files:/terraform-states ghcr.io/timeplus-io/resources-manager:1.0.0 apply -state=/terraform-states/terraform.tfstate
```

This command mount another directory on the container and ask the `apply` command to store the state file in that directory.

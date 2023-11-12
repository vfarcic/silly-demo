package main

import (
	"context"
	"fmt"
	"os"
	"regexp"

	"dagger.io/dagger"
)

var ctx = context.Background()
var image = "c8n.io/vfarcic/silly-demo"
var dev = false

func main() {
	if len(os.Getenv("DEV")) > 0 {
		dev = true
	}
	tag := os.Getenv("TAG")
	if dev {
		image = "ttl.sh/silly-demo"
		if len(tag) == 0 {
			tag = "0.0.1"
		}
	} else if len(tag) == 0 {
		panic("TAG environment variable is not set")
	}

	// initialize Dagger client
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	if err != nil {
		panic(err)
	}
	defer client.Close()

	// Actions
	publish(client, tag)
	pushTimoni(client, tag)
	if !dev {
		updateHelm(client, tag)
	}
}

func publish(client *dagger.Client, tag string) {
	if dev {
		publishImages(client, "Dockerfile-alpine", []string{tag})
	} else {
		publishImages(client, "Dockerfile", []string{tag, "latest"})
		publishImages(client, "Dockerfile-alpine", []string{fmt.Sprintf("%s-alpine", tag), "latest-alpine"})
	}
}

func publishImages(client *dagger.Client, dockerfile string, tags []string) {
	signed := false
	imageContainer := client.Host().Directory(".").DockerBuild(dagger.DirectoryDockerBuildOpts{
		Dockerfile: dockerfile,
	})
	for _, tag := range tags {
		imageTag := fmt.Sprintf("%s:%s", image, tag)
		imageAddr, err := imageContainer.Publish(ctx, imageTag)
		if err != nil {
			panic(err)
		}
		if !dev && !signed {
			cosignCmd := fmt.Sprintf("cosign sign --yes --key env://COSIGN_PRIVATE_KEY %s", imageAddr)
			if len(os.Getenv("REGISTRY_PASSWORD")) > 0 {
				cosignCmd = fmt.Sprintf("cosign login c8n.io --username vfarcic --password $REGISTRY_PASSWORD && %s", cosignCmd)
			}
			output, err := client.Container().
				From("bitnami/cosign:2.2.1").
				WithEnvVariable("COSIGN_PRIVATE_KEY", os.Getenv("COSIGN_PRIVATE_KEY")).
				WithEnvVariable("COSIGN_PASSWORD", os.Getenv("COSIGN_PASSWORD")).
				WithEnvVariable("REGISTRY_PASSWORD", os.Getenv("REGISTRY_PASSWORD")).
				WithEntrypoint([]string{"sh", "-c"}).
				WithExec([]string{cosignCmd}).
				Stderr(ctx)
			if err != nil {
				println(output)
				panic(err)
			}
			signed = true
		}
		fmt.Printf("Published image %s\n", imageAddr)
	}
}

func pushTimoni(client *dagger.Client, tag string) {
	_, err := client.Container().From("mikefarah/yq:4.35.2").
		WithDirectory("timoni", client.Host().Directory("timoni"), dagger.ContainerWithDirectoryOpts{
			Include: []string{"values.yaml"},
		}).
		WithExec(
			[]string{"--inplace", fmt.Sprintf(".values.image.tag = \"%s\"", tag), "timoni/values.yaml"},
			dagger.ContainerWithExecOpts{InsecureRootCapabilities: true},
		).
		File("timoni/values.yaml").
		Export(ctx, "timoni/values.yaml")
	if err != nil {
		panic(err)
	}
	fileContents, err := os.ReadFile("timoni/values.cue")
	if err != nil {
		panic(err)
	}
	regex := regexp.MustCompile(`image: tag:.*`)
	replacedString := regex.ReplaceAllString(string(fileContents), fmt.Sprintf("image: tag: \"%s\"", tag))
	file, err := os.OpenFile("timoni/values.cue", os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	_, err = file.WriteString(replacedString)
	if err != nil {
		panic(err)
	}
	err = file.Sync()
	if err != nil {
		panic(err)
	}
	regPass := client.SetSecret("registry-password", os.Getenv("REGISTRY_PASSWORD"))
	out, err := client.Container().From("golang:1.21.4").
		WithExec([]string{"go", "install", "github.com/stefanprodan/timoni/cmd/timoni@latest"}).
		WithDirectory("timoni", client.Host().Directory("timoni")).
		WithSecretVariable("REGISTRY_PASSWORD", regPass).
		WithExec([]string{"sh", "-c", fmt.Sprintf(`timoni mod push timoni oci://%s-package --version %s --creds vfarcic:$REGISTRY_PASSWORD`, image, tag)}).
		Stdout(ctx)
	if err != nil {
		println(out)
		panic(err)
	}
}

func updateHelm(client *dagger.Client, tag string) {
	_, err := client.Container().From("mikefarah/yq:4.35.2").
		WithDirectory("helm", client.Host().Directory("helm"), dagger.ContainerWithDirectoryOpts{
			Include: []string{"Chart.yaml", "values.yaml"},
		}).
		WithExec(
			[]string{"--inplace", fmt.Sprintf(".version = \"%s\"", tag), "helm/Chart.yaml"},
			dagger.ContainerWithExecOpts{InsecureRootCapabilities: true},
		).
		WithExec(
			[]string{"--inplace", fmt.Sprintf(".image.tag = \"%s\"", tag), "helm/values.yaml"},
			dagger.ContainerWithExecOpts{InsecureRootCapabilities: true},
		).
		Directory("helm").
		Export(ctx, "helm")
	if err != nil {
		panic(err)
	}
}

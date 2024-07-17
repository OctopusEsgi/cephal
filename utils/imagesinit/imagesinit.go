package imagesinit

import (
	"archive/tar"
	"bytes"
	"context"
	"io"
	"log"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
)

type imagesListParsed struct {
	ID    string
	Image []string
}

type ImagePath struct {
	ImageName  string
	Dockerfile string
}

func GetImagesList() ([]imagesListParsed, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}

	var listingOptions = image.ListOptions{
		All: true,
	}
	imgList, err := cli.ImageList(context.Background(), listingOptions)
	if err != nil {
		return nil, err
	}

	parsedList := make([]imagesListParsed, len(imgList))
	for i, img := range imgList {
		parsedList[i] = imagesListParsed{
			ID:    img.ID,
			Image: img.RepoTags,
		}
	}

	return parsedList, nil
}

// SOURCE : https://medium.com/@Frikkylikeme/controlling-docker-with-golang-code-b213d9699998
// Merci à Frikky !!
func buildImage(client *client.Client, tags []string, dockerfile string) error {
	ctx := context.Background()

	// Create a buffer
	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)
	defer tw.Close()

	// Create a filereader
	dockerFileReader, err := os.Open(dockerfile)
	if err != nil {
		return err
	}
	defer dockerFileReader.Close()

	// Read the actual Dockerfile
	readDockerFile, err := io.ReadAll(dockerFileReader)
	if err != nil {
		return err
	}

	// Make a TAR header for the file
	tarHeader := &tar.Header{
		Name: dockerfile,
		Size: int64(len(readDockerFile)),
	}

	// Writes the header described for the TAR file
	err = tw.WriteHeader(tarHeader)
	if err != nil {
		return err
	}

	// Writes the dockerfile data to the TAR file
	_, err = tw.Write(readDockerFile)
	if err != nil {
		return err
	}

	dockerFileTarReader := bytes.NewReader(buf.Bytes())

	// Define the build options to use for the file
	buildOptions := types.ImageBuildOptions{
		Context:    dockerFileTarReader,
		Dockerfile: dockerfile,
		Remove:     true,
		Tags:       tags,
	}

	// Build the actual image
	imageBuildResponse, err := client.ImageBuild(
		ctx,
		dockerFileTarReader,
		buildOptions,
	)

	if err != nil {
		return err
	}

	// Read the STDOUT from the build process
	defer imageBuildResponse.Body.Close()
	_, err = io.Copy(os.Stdout, imageBuildResponse.Body)
	if err != nil {
		return err
	}

	return nil
}

func imageExists(imageName string) (bool, error) {
	images, err := GetImagesList()
	if err != nil {
		return false, err
	}

	for _, img := range images {
		for _, tag := range img.Image {
			if tag == imageName {
				return true, nil
			}
		}
	}
	return false, nil
}

func ensureImage(imageName string, tags []string, dockerfile string) error {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		log.Fatal(err)
	}

	exists, err := imageExists(imageName)
	if err != nil {
		return err
	}

	if !exists {
		err := buildImage(cli, tags, dockerfile)
		if err != nil {
			return err
		}
	}

	return nil
}

func EnsureImagesList(images []ImagePath) error {
	for _, img := range images {
		tags := []string{img.ImageName}
		err := ensureImage(img.ImageName, tags, img.Dockerfile)
		if err != nil {
			return err
		}
	}
	return nil
}

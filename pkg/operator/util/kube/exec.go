package kube

import (
	"bytes"
	"io"
	"net/url"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
)

// ExecuteOptions options for execute
type ExecuteOptions struct {
	Command []string

	Namespace     string
	PodName       string
	ContainerName string

	Stdin              io.Reader
	CaptureStdout      bool
	CaptureStderr      bool
	PreserveWhitespace bool
}

func PodExecuteWithOptions(clientSet kubernetes.Interface, restConf *rest.Config, options ExecuteOptions) (string, string, error) {
	command := options.Command
	logrus.Error("command is ", command)
	req := clientSet.CoreV1().RESTClient().Post().Resource("pods").
		Name(options.PodName).
		Namespace(options.Namespace).
		SubResource("exec")
	scheme := runtime.NewScheme()
	if err := corev1.AddToScheme(scheme); err != nil {
		panic(err)
	}

	parameterCodec := runtime.NewParameterCodec(scheme)

	req.VersionedParams(&corev1.PodExecOptions{
		Command:   command,
		Container: "",
		Stdin:     false,
		Stdout:    true,
		Stderr:    true,
		TTY:       false,
	}, parameterCodec)

	var err error

	if exec, err := remotecommand.NewSPDYExecutor(restConf, "POST", req.URL()); err == nil {
		err = exec.Stream(remotecommand.StreamOptions{
			Stdin:  nil,
			Stdout: os.Stdout,
			Stderr: os.Stderr,
			Tty:    false,
		})
	}

	return "", "", err
}

// PodExecuteWithOptions execute with an execute options
func examplePodExecuteWithOptions(kubeClient kubernetes.Interface, restConfig *rest.Config, options ExecuteOptions) (string, string, error) {
	const tty = false
	req := kubeClient.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(options.PodName).
		Namespace(options.Namespace).
		SubResource("exec")
	req.VersionedParams(&corev1.PodExecOptions{
		Command: options.Command,
		Stdout:  options.CaptureStdout,
		Stderr:  options.CaptureStderr,
		TTY:     tty,
	}, runtime.NewParameterCodec(runtime.NewScheme()))

	var stdout, stderr bytes.Buffer
	err := execute("POST", req.URL(), restConfig, options.Stdin, &stdout, &stderr, tty)

	if options.PreserveWhitespace {
		return stdout.String(), stderr.String(), err
	}
	return strings.TrimSpace(stdout.String()), strings.TrimSpace(stderr.String()), err
}

func execute(method string, url *url.URL, config *rest.Config, stdin io.Reader, stdout, stderr io.Writer, tty bool) error {
	exec, err := remotecommand.NewSPDYExecutor(config, method, url)
	if err != nil {
		return err
	}
	logrus.Error("debug in pod command exec")
	return exec.Stream(remotecommand.StreamOptions{
		Stdin:  stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Tty:    tty,
	})
}

package rfe

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/cskr/pubsub"
)

const firestoreEmulatorHost = "FIRESTORE_EMULATOR_HOST"
const loopbackIP = "127.0.0.1"
const firestoreStdoutTopic = "firetore-logs"

func getFirestoreEmulatorCmd(verbose bool, port uint16) *exec.Cmd {
	var cmdArgs = []string{"beta", "emulators", "firestore", "start"}

	if !verbose {
		cmdArgs = append(cmdArgs, "--quiet")
	}
	if port == 0 {
		port = getFreeHostPort()
	}

	hostPort := fmt.Sprintf("--host-port=%s:%d", loopbackIP, port)
	cmdArgs = append(cmdArgs, hostPort)

	return exec.Command("gcloud", cmdArgs...)
}

func startFirestoreEmulator(verbose bool, port uint16) (cmd *exec.Cmd, stdout io.ReadCloser) {
	cmd = getFirestoreEmulatorCmd(verbose, port)
	stdout = getBothStdoutStderrCombined(cmd)

	makeProcessKillable(cmd)
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	return
}

func publishFirestoreLogs(firestoreStdout io.ReadCloser, firestorePubSub *pubsub.PubSub) {
	streamReadlinesIterator, err := getStreamReadlinesIterator(firestoreStdout)
	if err != nil {
		log.Fatal(err)
	}

	for line := range streamReadlinesIterator {
		firestorePubSub.Pub(line, firestoreStdoutTopic)
	}
}

func firestoreEmulatorIsReady(stdoutLine string) bool {
	return strings.Contains(stdoutLine, "Dev App Server is now running")
}

func setHostEnvIfIsConfigured(stdoutLine string) {
	pos := strings.Index(stdoutLine, firestoreEmulatorHost+"=")

	if pos > 0 {
		host := stdoutLine[pos+len(firestoreEmulatorHost)+1:]
		os.Setenv(firestoreEmulatorHost, host)
	}
}

func waitForFirestoreToBeReady(ps *pubsub.PubSub) {
	channel := ps.Sub(firestoreStdoutTopic)
	for {
		if msg, ok := <-channel; ok {

			if firestoreEmulatorIsReady(fmt.Sprintf("%s", msg)) {
				go ps.Unsub(channel, firestoreStdoutTopic)
			}

			setHostEnvIfIsConfigured(fmt.Sprintf("%s", msg))
		} else {
			break
		}
	}
}

type FirestoreEmulator struct {
	pubSub  *pubsub.PubSub
	stdout  io.ReadCloser
	cmd     *exec.Cmd
	Verbose bool
	Port    uint16
}

func (f *FirestoreEmulator) Start() {
	f.pubSub = pubsub.New(0)
	f.cmd, f.stdout = startFirestoreEmulator(f.Verbose, f.Port)

	go publishFirestoreLogs(f.stdout, f.pubSub)
	go logPubSubTopic(f.pubSub, firestoreStdoutTopic)
	waitForFirestoreToBeReady(f.pubSub)
}

func (f *FirestoreEmulator) Shutdown() {
	f.stdout.Close()
	killProcessGroup(f.cmd)
	f.pubSub.Shutdown()
}

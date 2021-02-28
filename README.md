# rfe

Run Firestore Emulator

Example:
```go
// main_test.go
import (
    "testing"


    "github.com/kmarilleau/rfe"
    "gopkg.in/h2non/gock.v1"

)
func TestMain(m *testing.M) {
    firestoreEmulator := rfe.FirestoreEmulator{Verbose: true}

    firestoreEmulator.Start()
    defer firestoreEmulator.Shutdown()

    m.Run()
}
```

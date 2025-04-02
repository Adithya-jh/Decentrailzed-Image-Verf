package main

import (
	"bytes"
    "context"
    "crypto/sha256"
    "encoding/json"
    "io/ioutil"
    "log"
    "net/http"

    admissionv1 "k8s.io/api/admission/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/accounts/abi/bind"
    "github.com/ethereum/go-ethereum/ethclient"
    "github.com/adithya-jh/admission-controller/imreg"
)

// Global variables, function declarations, etc.
var (
    ethClient     *ethclient.Client
    imageRegistry *imreg.ImageRegistry
)

func initEthereum() {
    var err error
    ethClient, err = ethclient.Dial("http://localhost:8545")
    if err != nil {
        log.Fatalf("Failed to connect to Ethereum client: %v", err)
    }
    contractAddress := common.HexToAddress("0x5CeC6C7F106d58d29Dd89B9cd0302f39f3DbcDD0")
    imageRegistry, err = imreg.NewImageRegistry(contractAddress, ethClient)
    if err != nil {
        log.Fatalf("Failed to instantiate contract: %v", err)
    }
}

func isImageTrusted(image string) bool {
    hashBytes := sha256.Sum256([]byte(image))
    hash := common.BytesToHash(hashBytes[:])
    log.Printf("Computed hash for image %s: %s", image, hash.Hex())
    ctx := context.Background()
    trusted, err := imageRegistry.IsImageTrusted(&bind.CallOpts{Context: ctx}, hash)
    if err != nil {
        log.Printf("Error querying blockchain: %v", err)
        return false
    }
    log.Printf("Blockchain returned: %t", trusted)
    return trusted
}

func sendEvent(result bool, image string, uid string) {
    payload := map[string]interface{}{
        "uid":     uid,
        "image":   image,
        "allowed": result,
    }
    payloadBytes, err := json.Marshal(payload)
    if err != nil {
        log.Printf("Error marshalling event payload: %v", err)
        return
    }
    // Change the URL if your dashboard backend is running on a different host/port.
    url := "http://host.docker.internal:4000/api/event"
    resp, err := http.Post(url, "application/json", bytes.NewBuffer(payloadBytes))
    if err != nil {
        log.Printf("Error sending event to dashboard: %v", err)
        return
    }
    defer resp.Body.Close()
    log.Printf("Event sent to dashboard, response status: %s", resp.Status)
}

func handleValidate(w http.ResponseWriter, r *http.Request) {
    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        http.Error(w, "could not read request", http.StatusBadRequest)
        return
    }
    var admissionReviewReq admissionv1.AdmissionReview
    if err = json.Unmarshal(body, &admissionReviewReq); err != nil {
        http.Error(w, "could not parse admission review", http.StatusBadRequest)
        return
    }
    log.Printf("Received pod creation request: %s\n", string(admissionReviewReq.Request.Object.Raw))

    // For demonstration, we use a hardcoded image identifier.
    // In a real scenario, extract the container image from the Pod spec.
    image := "nginx:latest"
    allowed := isImageTrusted(image)
    admissionResponse := admissionv1.AdmissionResponse{
        Allowed: allowed,
        UID:     admissionReviewReq.Request.UID,
    }
    if !allowed {
        admissionResponse.Result = &metav1.Status{
            Message: "Image is not verified by blockchain",
        }
    }
    admissionReviewResp := admissionv1.AdmissionReview{
        Response: &admissionResponse,
    }
    respBytes, err := json.Marshal(admissionReviewResp)
    if err != nil {
        http.Error(w, "could not marshal response", http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    w.Write(respBytes)

    // Send event to the dashboard backend asynchronously.
    go sendEvent(allowed, image, string(admissionReviewReq.Request.UID))

}


func main() {
    initEthereum()
    http.HandleFunc("/validate", handleValidate)
    log.Println("Admission controller starting on port 8443...")
    log.Fatal(http.ListenAndServeTLS(":8443", "tls.crt", "tls.key", nil))
}

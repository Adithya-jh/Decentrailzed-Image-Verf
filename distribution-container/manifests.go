package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/distribution/distribution/v3"
	"github.com/distribution/distribution/v3/internal/dcontext"
	"github.com/distribution/distribution/v3/manifest/manifestlist"
	"github.com/distribution/distribution/v3/manifest/ocischema"
	"github.com/distribution/distribution/v3/manifest/schema2"
	"github.com/distribution/distribution/v3/registry/api/errcode"
	"github.com/distribution/distribution/v3/registry/storage"
	"github.com/distribution/distribution/v3/registry/storage/driver"
	"github.com/distribution/reference"
	"github.com/gorilla/handlers"
	"github.com/opencontainers/go-digest"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
	"golang.org/x/sync/errgroup"
)

// Constants for default platform and maximum allowed manifest size.
const (
	defaultArch         = "amd64"
	defaultOS           = "linux"
	maxManifestBodySize = 4 * 1024 * 1024
	imageClass          = "image"
)

// storageType is used to track different manifest formats.
type storageType int

const (
	manifestSchema2 storageType = iota
	manifestlistSchema
	ociSchema
	ociImageIndexSchema
	numStorageTypes
)

// manifestDispatcher creates an HTTP handler for manifest routes.
func manifestDispatcher(ctx *Context, r *http.Request) http.Handler {
	manifestHandler := &manifestHandler{
		Context: ctx,
	}

	ref := getReference(ctx)
	dgst, err := digest.Parse(ref)
	if err != nil {
		// We just have a tag
		manifestHandler.Tag = ref
	} else {
		manifestHandler.Digest = dgst
	}

	mhandler := handlers.MethodHandler{
		http.MethodGet:  http.HandlerFunc(manifestHandler.GetManifest),
		http.MethodHead: http.HandlerFunc(manifestHandler.GetManifest),
	}

	if !ctx.readOnly {
		mhandler[http.MethodPut] = http.HandlerFunc(manifestHandler.PutManifest)
		mhandler[http.MethodDelete] = http.HandlerFunc(manifestHandler.DeleteManifest)
	}

	return mhandler
}

// manifestHandler holds context for all manifest operations.
type manifestHandler struct {
	*Context

	// Either Tag or Digest is set based on the incoming request.
	Tag    string
	Digest digest.Digest
}

// GetManifest fetches the manifest from the storage backend, if it exists.
func (imh *manifestHandler) GetManifest(w http.ResponseWriter, r *http.Request) {
	dcontext.GetLogger(imh).Debug("GetImageManifest")

	manifests, err := imh.Repository.Manifests(imh)
	if err != nil {
		imh.Errors = append(imh.Errors, err)
		return
	}

	var supports [numStorageTypes]bool
	for _, acceptHeader := range r.Header["Accept"] {
		for _, mediaType := range strings.Split(acceptHeader, ",") {
			if mediaType, _, err = mime.ParseMediaType(mediaType); err != nil {
				continue
			}
			if mediaType == schema2.MediaTypeManifest {
				supports[manifestSchema2] = true
			}
			if mediaType == manifestlist.MediaTypeManifestList {
				supports[manifestlistSchema] = true
			}
			if mediaType == v1.MediaTypeImageManifest {
				supports[ociSchema] = true
			}
			if mediaType == v1.MediaTypeImageIndex {
				supports[ociImageIndexSchema] = true
			}
		}
	}

	if imh.Tag != "" {
		tags := imh.Repository.Tags(imh)
		desc, err := tags.Get(imh, imh.Tag)
		if err != nil {
			if _, ok := err.(distribution.ErrTagUnknown); ok {
				imh.Errors = append(imh.Errors, errcode.ErrorCodeManifestUnknown.WithDetail(err))
			} else {
				imh.Errors = append(imh.Errors, errcode.ErrorCodeUnknown.WithDetail(err))
			}
			return
		}
		imh.Digest = desc.Digest
	}

	if etagMatch(r, imh.Digest.String()) {
		w.WriteHeader(http.StatusNotModified)
		return
	}

	var options []distribution.ManifestServiceOption
	if imh.Tag != "" {
		options = append(options, distribution.WithTag(imh.Tag))
	}

	manifest, err := manifests.Get(imh, imh.Digest, options...)
	if err != nil {
		if _, ok := err.(distribution.ErrManifestUnknownRevision); ok {
			imh.Errors = append(imh.Errors, errcode.ErrorCodeManifestUnknown.WithDetail(err))
		} else {
			imh.Errors = append(imh.Errors, errcode.ErrorCodeUnknown.WithDetail(err))
		}
		return
	}

	// Determine manifest format (schema2, manifestlist, OCI, etc.).
	manifestType := manifestSchema2
	manifestList, isManifestList := manifest.(*manifestlist.DeserializedManifestList)
	if _, isOCImanifest := manifest.(*ocischema.DeserializedManifest); isOCImanifest {
		manifestType = ociSchema
	} else if isManifestList {
		if manifestList.MediaType == manifestlist.MediaTypeManifestList {
			manifestType = manifestlistSchema
		} else if manifestList.MediaType == v1.MediaTypeImageIndex {
			manifestType = ociImageIndexSchema
		}
	}

	// Handle unsupported OCI manifest or index requests.
	if manifestType == ociSchema && !supports[ociSchema] {
		imh.Errors = append(imh.Errors, errcode.ErrorCodeManifestUnknown.WithMessage("OCI manifest found, but accept header does not support OCI manifests"))
		return
	}
	if manifestType == ociImageIndexSchema && !supports[ociImageIndexSchema] {
		imh.Errors = append(imh.Errors, errcode.ErrorCodeManifestUnknown.WithMessage("OCI index found, but accept header does not support OCI indexes"))
		return
	}

	// Handle old clients that don't support manifest lists.
	if imh.Tag != "" && manifestType == manifestlistSchema && !supports[manifestlistSchema] {
		dcontext.GetLogger(imh).Infof("rewriting manifest list %s in schema1 format to support old client", imh.Digest.String())

		var manifestDigest digest.Digest
		for _, manifestDescriptor := range manifestList.Manifests {
			if manifestDescriptor.Platform.Architecture == defaultArch && manifestDescriptor.Platform.OS == defaultOS {
				manifestDigest = manifestDescriptor.Digest
				break
			}
		}
		if manifestDigest == "" {
			imh.Errors = append(imh.Errors, errcode.ErrorCodeManifestUnknown)
			return
		}

		manifest, err = manifests.Get(imh, manifestDigest)
		if err != nil {
			if _, ok := err.(distribution.ErrManifestUnknownRevision); ok {
				imh.Errors = append(imh.Errors, errcode.ErrorCodeManifestUnknown.WithDetail(err))
			} else {
				imh.Errors = append(imh.Errors, errcode.ErrorCodeUnknown.WithDetail(err))
			}
			return
		}

		if _, isSchema2 := manifest.(*schema2.DeserializedManifest); isSchema2 && !supports[manifestSchema2] {
			imh.Errors = append(imh.Errors, errcode.ErrorCodeManifestInvalid.WithMessage("Schema 2 manifest not supported by client"))
			return
		}
		imh.Digest = manifestDigest
	}

	ct, p, err := manifest.Payload()
	if err != nil {
		return
	}

	w.Header().Set("Content-Type", ct)
	w.Header().Set("Content-Length", fmt.Sprint(len(p)))
	w.Header().Set("Docker-Content-Digest", imh.Digest.String())
	w.Header().Set("Etag", fmt.Sprintf(`"%s"`, imh.Digest))

	if r.Method == http.MethodHead {
		return
	}

	if _, err := w.Write(p); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// etagMatch checks whether the provided etag is in the If-None-Match header.
func etagMatch(r *http.Request, etag string) bool {
	for _, headerVal := range r.Header["If-None-Match"] {
		if headerVal == etag || headerVal == fmt.Sprintf(`"%s"`, etag) {
			return true
		}
	}
	return false
}

// PutManifest validates and stores a manifest in the registry.
func (imh *manifestHandler) PutManifest(w http.ResponseWriter, r *http.Request) {
	dcontext.GetLogger(imh).Debug("PutManifest CALLED")

	manifests, err := imh.Repository.Manifests(imh)
	if err != nil {
		imh.Errors = append(imh.Errors, err)
		return
	}

	var jsonBuf bytes.Buffer
	if err := copyFullPayload(imh, w, r, &jsonBuf, maxManifestBodySize, "image manifest PUT"); err != nil {
		imh.Errors = append(imh.Errors, errcode.ErrorCodeManifestInvalid.WithDetail(err.Error()))
		return
	}

	mediaType := r.Header.Get("Content-Type")
	manifest, desc, err := distribution.UnmarshalManifest(mediaType, jsonBuf.Bytes())
	if err != nil {
		imh.Errors = append(imh.Errors, errcode.ErrorCodeManifestInvalid.WithDetail(err))
		return
	}

	// Validate digest or tag existence.
	if imh.Digest != "" {
		if desc.Digest != imh.Digest {
			dcontext.GetLogger(imh).Errorf("payload digest does not match: %q != %q", desc.Digest, imh.Digest)
			imh.Errors = append(imh.Errors, errcode.ErrorCodeDigestInvalid)
			return
		}
	} else if imh.Tag != "" {
		imh.Digest = desc.Digest
	} else {
		imh.Errors = append(imh.Errors, errcode.ErrorCodeTagInvalid.WithDetail("no tag or digest specified"))
		return
	}

	isAnOCIManifest := mediaType == v1.MediaTypeImageManifest || mediaType == v1.MediaTypeImageIndex
	if isAnOCIManifest {
		dcontext.GetLogger(imh).Debug("Putting an OCI Manifest!")
	} else {
		dcontext.GetLogger(imh).Debug("Putting a Docker Manifest!")
	}

	var options []distribution.ManifestServiceOption
	if imh.Tag != "" {
		options = append(options, distribution.WithTag(imh.Tag))
	}

	dcontext.GetLogger(imh).Debugf("DEBUG: About to apply resource policy for repository: %s", imh.Repository.Named().Name())
	if err := imh.applyResourcePolicy(manifest); err != nil {
		dcontext.GetLogger(imh).Debugf("Policy REJECTED the push: %v", err)
		imh.Errors = append(imh.Errors, err)
		return
	}
	dcontext.GetLogger(imh).Debug("Policy check PASSED. Continuing with manifest storage...")

	_, err = manifests.Put(imh, manifest, options...)
	if err != nil {
		if err == distribution.ErrUnsupported {
			imh.Errors = append(imh.Errors, errcode.ErrorCodeUnsupported)
			return
		}
		if err == distribution.ErrAccessDenied {
			imh.Errors = append(imh.Errors, errcode.ErrorCodeDenied)
			return
		}
		switch err := err.(type) {
		case distribution.ErrManifestVerification:
			for _, verificationError := range err {
				switch verificationError := verificationError.(type) {
				case distribution.ErrManifestBlobUnknown:
					imh.Errors = append(imh.Errors, errcode.ErrorCodeManifestBlobUnknown.WithDetail(verificationError.Digest))
				case distribution.ErrManifestNameInvalid:
					imh.Errors = append(imh.Errors, errcode.ErrorCodeNameInvalid.WithDetail(err))
				case distribution.ErrManifestUnverified:
					imh.Errors = append(imh.Errors, errcode.ErrorCodeManifestUnverified)
				default:
					if verificationError == digest.ErrDigestInvalidFormat {
						imh.Errors = append(imh.Errors, errcode.ErrorCodeDigestInvalid)
					} else {
						imh.Errors = append(imh.Errors, errcode.ErrorCodeUnknown, verificationError)
					}
				}
			}
		case errcode.Error:
			imh.Errors = append(imh.Errors, err)
		default:
			imh.Errors = append(imh.Errors, errcode.ErrorCodeUnknown.WithDetail(err))
		}
		return
	}

	// Tag this manifest if a tag was specified.
	if imh.Tag != "" {
		tags := imh.Repository.Tags(imh)
		if err := tags.Tag(imh, imh.Tag, desc); err != nil {
			imh.Errors = append(imh.Errors, errcode.ErrorCodeUnknown.WithDetail(err))
			return
		}
	}

	// Construct a canonical URL for the uploaded manifest.
	ref, err := reference.WithDigest(imh.Repository.Named(), imh.Digest)
	if err != nil {
		imh.Errors = append(imh.Errors, errcode.ErrorCodeUnknown.WithDetail(err))
		return
	}
	location, err := imh.urlBuilder.BuildManifestURL(ref)
	if err != nil {
		dcontext.GetLogger(imh).Errorf("error building manifest url from digest: %v", err)
	}
	w.Header().Set("Location", location)
	w.Header().Set("Docker-Content-Digest", imh.Digest.String())
	w.WriteHeader(http.StatusCreated)

	dcontext.GetLogger(imh).Debug("Succeeded in putting manifest!")

	// --- Begin Blockchain Registration Call ---
	go callBlockchainRegistration(imh)
	// --- End Blockchain Registration Call ---
}

// callBlockchainRegistration makes an asynchronous call to register the manifest on the blockchain,
// then automatically posts the verification details to the dashboard backend.
func callBlockchainRegistration(imh *manifestHandler) {
	payload := map[string]interface{}{
		"repository": imh.Repository.Named().String(),
		"digest":     imh.Digest.String(),
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		dcontext.GetLogger(imh).Errorf("Error marshalling blockchain payload: %v", err)
		return
	}

	resp, err := http.Post("http://localhost:5001/register", "application/json", bytes.NewReader(payloadBytes))
	if err != nil {
		dcontext.GetLogger(imh).Errorf("Error calling blockchain registration service: %v", err)
		return
	}
	defer resp.Body.Close()

	var bcResponse map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&bcResponse); err != nil {
		dcontext.GetLogger(imh).Errorf("Error decoding blockchain registration response: %v", err)
		return
	}
	dcontext.GetLogger(imh).Debugf("Blockchain registration response: %+v", bcResponse)

	bcResponse["repository"] = imh.Repository.Named().String()
	bcResponse["digest"] = imh.Digest.String()
	if _, ok := bcResponse["timestamp"]; !ok {
		bcResponse["timestamp"] = time.Now().UTC().Format(time.RFC3339)
	}

	verificationBytes, err := json.Marshal(bcResponse)
	if err != nil {
		dcontext.GetLogger(imh).Errorf("Error marshalling verification payload: %v", err)
		return
	}

	resp2, err := http.Post("http://localhost:4000/api/verification", "application/json", bytes.NewReader(verificationBytes))
	if err != nil {
		dcontext.GetLogger(imh).Errorf("Error sending verification to dashboard backend: %v", err)
		return
	}
	defer resp2.Body.Close()
	dcontext.GetLogger(imh).Debugf("Dashboard backend verification response: %s", resp2.Status)
}

// applyResourcePolicy is a hard-coded check that requires a repository to start with "approved-"
func (imh *manifestHandler) applyResourcePolicy(manifest distribution.Manifest) error {
	dcontext.GetLogger(imh).Debug("applyResourcePolicy CALLED")

	// Hard-coded check: the repository name must start with "approved-"
	repoName := imh.Repository.Named().Name()
	dcontext.GetLogger(imh).Debugf("applyResourcePolicy: Hard-coded check for repository: %s", repoName)

	if !strings.HasPrefix(repoName, "approved-") {
		msg := fmt.Sprintf("Repository %s is not authorized. Must start with 'approved-'", repoName)
		dcontext.GetLogger(imh).Debugf("Policy fails: %s", msg)
		return errcode.ErrorCodeDenied.WithMessage(msg)
	}

	dcontext.GetLogger(imh).Debug("Hard-coded policy check PASSED.")
	return nil
}

// DeleteManifest removes the manifest (or a tagged reference to it) from the registry.
func (imh *manifestHandler) DeleteManifest(w http.ResponseWriter, r *http.Request) {
	dcontext.GetLogger(imh).Debug("DeleteImageManifest")

	if imh.App.isCache {
		imh.Errors = append(imh.Errors, errcode.ErrorCodeUnsupported)
		return
	}

	// Handle "delete by tag"
	if imh.Tag != "" {
		dcontext.GetLogger(imh).Debug("DeleteImageTag")
		tagService := imh.Repository.Tags(imh.Context)
		if err := tagService.Untag(imh.Context, imh.Tag); err != nil {
			switch err.(type) {
			case distribution.ErrTagUnknown, driver.PathNotFoundError:
				imh.Errors = append(imh.Errors, errcode.ErrorCodeManifestUnknown.WithDetail(err))
			default:
				imh.Errors = append(imh.Errors, errcode.ErrorCodeUnknown.WithDetail(err))
			}
			return
		}
		w.WriteHeader(http.StatusAccepted)
		return
	}

	manifests, err := imh.Repository.Manifests(imh)
	if err != nil {
		imh.Errors = append(imh.Errors, err)
		return
	}

	if err := manifests.Delete(imh, imh.Digest); err != nil {
		switch err {
		case digest.ErrDigestUnsupported, digest.ErrDigestInvalidFormat:
			imh.Errors = append(imh.Errors, errcode.ErrorCodeDigestInvalid)
			return
		case distribution.ErrBlobUnknown:
			imh.Errors = append(imh.Errors, errcode.ErrorCodeManifestUnknown)
			return
		case distribution.ErrUnsupported:
			imh.Errors = append(imh.Errors, errcode.ErrorCodeUnsupported)
			return
		default:
			imh.Errors = append(imh.Errors, errcode.ErrorCodeUnknown)
			return
		}
	}

	tagService := imh.Repository.Tags(imh)
	referencedTags, err := tagService.Lookup(imh, v1.Descriptor{Digest: imh.Digest})
	if err != nil {
		imh.Errors = append(imh.Errors, err)
		return
	}

	var (
		errs []error
		mu   sync.Mutex
	)
	g := errgroup.Group{}
	g.SetLimit(storage.DefaultConcurrencyLimit)

	for _, tag := range referencedTags {
		tag := tag
		g.Go(func() error {
			if err := tagService.Untag(imh, tag); err != nil {
				mu.Lock()
				errs = append(errs, err)
				mu.Unlock()
			}
			return nil
		})
	}
	_ = g.Wait()
	imh.Errors = errs
	w.WriteHeader(http.StatusAccepted)
}

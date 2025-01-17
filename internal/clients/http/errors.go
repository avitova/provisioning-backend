package http

import (
	"github.com/RHEnVision/provisioning-backend/internal/usrerr"
)

// Sources
var (
	ErrApplicationTypeNotFound          = usrerr.New(404, "application type 'provisioning' not found in sources", "")
	ErrAuthenticationForSourcesNotFound = usrerr.New(404, "authentications for source weren't found in sources", "")
	ErrApplicationRead                  = usrerr.New(500, "application read returned no application type in sources", "")
	ErrSourcesInvalidAuthentication     = usrerr.New(400, "insufficient data for authentication", "")
)

// Image Builder
var (
	ErrCloneNotFound        = usrerr.New(404, "image clone not found", "")
	ErrImageStatus          = usrerr.New(500, "build of requested image has not finished yet", "image still building")
	ErrUnknownImageType     = usrerr.New(500, "unknown image type", "")
	ErrUploadStatus         = usrerr.New(500, "cannot get image status", "")
	ErrImageRequestNotFound = usrerr.New(500, "image compose request not found", "")
)

// EC2
var (
	ErrDuplicatePubkey             = usrerr.New(406, "public key already exists in target cloud provider account and region", "")
	ErrPubkeyNotFound              = usrerr.New(404, "pubkey not found in AWS account", "")
	ErrServiceAccountUnsupportedOp = usrerr.New(500, "unsupported operation on service account", "")
	ErrARNParsing                  = usrerr.New(500, "ARN parsing error", "")
	ErrNoReservation               = usrerr.New(404, "no reservation was found in AWS response", "")
)

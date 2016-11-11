package scepserver

// SCEPRequest is a SCEP server request.
type SCEPRequest struct {
	Operation string
	Message   []byte
	Err       error // request error
}

// SCEPResponse is a SCEP server response.
// Business errors will be encoded as a CertRep message
// with pkiStatus FAILURE and a failInfo attribute.
type SCEPResponse struct {
	CACertNum int //chain
	Data      []byte
	Err       error // response error
}

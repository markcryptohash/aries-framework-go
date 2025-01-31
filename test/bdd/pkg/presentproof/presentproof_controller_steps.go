/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package presentproof

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"

	"github.com/cucumber/godog"

	"github.com/markcryptohash/aries-framework-go/pkg/client/presentproof"
	didexcmd "github.com/markcryptohash/aries-framework-go/pkg/controller/command/didexchange"
	"github.com/markcryptohash/aries-framework-go/pkg/controller/command/kms"
	presentproofcmd "github.com/markcryptohash/aries-framework-go/pkg/controller/command/presentproof"
	"github.com/markcryptohash/aries-framework-go/pkg/controller/command/verifiable"
	connectionrest "github.com/markcryptohash/aries-framework-go/pkg/controller/rest/connection"
	"github.com/markcryptohash/aries-framework-go/pkg/didcomm/common/service"
	"github.com/markcryptohash/aries-framework-go/pkg/didcomm/protocol/decorator"
	protocol "github.com/markcryptohash/aries-framework-go/pkg/didcomm/protocol/presentproof"
	docverifiable "github.com/markcryptohash/aries-framework-go/pkg/doc/verifiable"
	arieskms "github.com/markcryptohash/aries-framework-go/pkg/kms"
	"github.com/markcryptohash/aries-framework-go/pkg/vdr/fingerprint"
	"github.com/markcryptohash/aries-framework-go/test/bdd/pkg/connection"
	"github.com/markcryptohash/aries-framework-go/test/bdd/pkg/context"
	didexsteps "github.com/markcryptohash/aries-framework-go/test/bdd/pkg/didexchange"
	"github.com/markcryptohash/aries-framework-go/test/bdd/pkg/outofband"
	"github.com/markcryptohash/aries-framework-go/test/bdd/pkg/util"
)

const (
	operationID                    = "/presentproof"
	operationIDV3                  = operationID + "/v3"
	sendRequestPresentation        = operationID + "/send-request-presentation"
	sendRequestPresentationV3      = operationIDV3 + "/send-request-presentation"
	sendProposalPresentation       = operationID + "/send-propose-presentation"
	sendProposalPresentationV3     = operationIDV3 + "/send-propose-presentation"
	acceptProposePresentation      = operationID + "/%s/accept-propose-presentation"
	acceptProposePresentationV3    = operationIDV3 + "/%s/accept-propose-presentation"
	acceptRequestPresentation      = operationID + "/%s/accept-request-presentation"
	acceptRequestPresentationV3    = operationIDV3 + "/%s/accept-request-presentation"
	negotiateRequestPresentation   = operationID + "/%s/negotiate-request-presentation"
	negotiateRequestPresentationV3 = operationIDV3 + "/%s/negotiate-request-presentation"
	acceptPresentation             = operationID + "/%s/accept-presentation"
	declinePresentation            = operationID + "/%s/decline-presentation"
	acceptProblemReport            = operationID + "/%s/accept-problem-report"

	verifiablePresentations = "/verifiable/presentations"
)

// ControllerSteps supports steps for Present Proof controller.
type ControllerSteps struct {
	bddContext *context.BDDContext
	did        map[string]string
	nameToPIID map[string]string
}

// NewPresentProofControllerSteps creates steps for Present Proof controller.
func NewPresentProofControllerSteps() *ControllerSteps {
	return &ControllerSteps{
		did:        make(map[string]string),
		nameToPIID: make(map[string]string),
	}
}

// SetContext sets every scenario with a fresh context.
func (s *ControllerSteps) SetContext(ctx *context.BDDContext) {
	s.bddContext = ctx
}

// RegisterSteps registers agent steps
// nolint:lll
func (s *ControllerSteps) RegisterSteps(gs *godog.Suite) {
	gs.Step(`^"([^"]*)" has established connection with "([^"]*)" through PresentProof controller$`, s.establishConnection)
	gs.Step(`^"([^"]*)" has established DIDComm v2 connection with "([^"]*)" through PresentProof controller$`, s.establishDIDCommV2Connection)
	gs.Step(`^"([^"]*)" sends OOBv2 invitation to "([^"]*)" to establish connection through PresentProof controller$`, s.createDIDCommV2ConnectionWithOOBv2)
	gs.Step(`^"([^"]*)" waits to receive a message from "([^"]*)" through PresentProof controller$`, s.waitForMessageFrom)
	gs.Step(`^"([^"]*)" sends a propose presentation to "([^"]*)" through PresentProof controller$`, s.sendProposePresentation)
	gs.Step(`^"([^"]*)" sends a propose presentation v3 to "([^"]*)" through PresentProof controller$`, s.sendProposePresentationV3)
	gs.Step(`^"([^"]*)" negotiates about the request presentation with a proposal through PresentProof controller$`, s.negotiateRequestPresentation)
	gs.Step(`^"([^"]*)" negotiates about the request presentation v3 with a proposal through PresentProof controller$`, s.negotiateRequestPresentationV3)
	gs.Step(`^"([^"]*)" successfully accepts a presentation with "([^"]*)" name through PresentProof controller$`, s.acceptPresentation)
	gs.Step(`^"([^"]*)" successfully accepts a presentation with "([^"]*)" name and "([^"]*)" redirect through PresentProof controller$`, s.acceptPresentationWithRedirect)
	gs.Step(`^"([^"]*)" checks that presentation is being stored under the "([^"]*)" name$`, s.checkPresentation)
	gs.Step(`^"([^"]*)" sends "([^"]*)" to "([^"]*)" through PresentProof controller$`, s.sendMessage)
	gs.Step(`^"([^"]*)" declines presentation "([^"]*)" from "([^"]*)" and redirects prover to "([^"]*)" through PresentProof controller$`, s.declinePresentationWithRedirect)
	gs.Step(`^"([^"]*)" validates present proof state "([^"]*)" and redirect "([^"]*)" with status "([^"]*)" through PresentProof controller$`, s.validateState)
}

func (s *ControllerSteps) establishConnection(inviter, invitee string) error {
	ds := didexsteps.NewDIDExchangeControllerSteps()
	ds.SetContext(s.bddContext)

	err := ds.EstablishConnection(inviter, invitee)
	if err != nil {
		return fmt.Errorf("unable to establish connection between [%s] and [%s]: %w", inviter, invitee, err)
	}

	inviterDID, err := s.agentDID(ds, inviter)
	if err != nil {
		return err
	}

	s.did[inviter] = inviterDID

	inviteeDID, err := s.agentDID(ds, invitee)
	if err != nil {
		return err
	}

	s.did[invitee] = inviteeDID

	return nil
}

func (s *ControllerSteps) establishDIDCommV2Connection(agent1, agent2 string) error {
	didexSteps := didexsteps.NewDIDExchangeControllerSteps()
	didexSteps.SetContext(s.bddContext)

	connSteps := connection.NewControllerSteps()
	connSteps.SetContext(s.bddContext)

	err := didexSteps.CreatePublicDID(agent1, "")
	if err != nil {
		return fmt.Errorf("create public DID for [%s] using didexchange controller step: %w", agent1, err)
	}

	err = didexSteps.CreatePublicDID(agent2, "")
	if err != nil {
		return fmt.Errorf("create public DID for [%s] using didexchange controller step: %w", agent2, err)
	}

	err = connSteps.HasDIDCommV2Connection(agent1, agent2)
	if err != nil {
		return fmt.Errorf(
			"creating DIDComm v2 connection for [%s] and [%s] using connection controller step: %w",
			agent1, agent2, err)
	}

	s.did[agent1] = s.bddContext.PublicDIDs[agent1]
	s.did[agent2] = s.bddContext.PublicDIDs[agent2]

	return nil
}

func (s *ControllerSteps) createDIDCommV2ConnectionWithOOBv2(inviter, invitee string) error {
	didexSteps := didexsteps.NewDIDExchangeControllerSteps()
	didexSteps.SetContext(s.bddContext)

	inviterDID, err := didexSteps.CreatePublicDIDWithKeyType(inviter, "ECDSAP256IEEEP1363", "NISTP256ECDHKW")
	if err != nil {
		return fmt.Errorf("create public DID for [%s] using didexchange controller step: %w", inviter, err)
	}

	oobSteps := outofband.NewOutofbandControllerSteps()
	oobSteps.SetContext(s.bddContext)

	connSteps := connection.NewControllerSteps()
	connSteps.SetContext(s.bddContext)

	err = oobSteps.CreateOOBV2(inviter)
	if err != nil {
		return fmt.Errorf("creating oob v2 invitation for [%s] using oob controller steps: %w", inviter, err)
	}

	connID, err := oobSteps.AcceptOOBV2Invitation(inviter, invitee)
	if err != nil {
		return fmt.Errorf(
			"accepting oob v2 invitation from [%s] with [%s] using oob controller step: %w",
			inviter, invitee, err)
	}

	url, ok := s.bddContext.GetControllerURL(invitee)
	if !ok {
		return fmt.Errorf("unable to find controller URL registered for agent [%s]", invitee)
	}

	resp := &didexcmd.QueryConnectionResponse{}

	err = util.SendHTTP(http.MethodGet, url+"/connections/"+connID, nil, &resp)
	if err != nil {
		return fmt.Errorf("getting connection record from agent [%s]: %w", invitee, err)
	}

	inviteeConn := resp.Result.Record

	s.bddContext.PeerDIDs[invitee] = inviteeConn.MyDID

	s.did[inviter] = inviterDID
	s.did[invitee] = inviteeConn.MyDID

	return nil
}

func getMsgBytes(msgFile string) (string, []byte, error) {
	_, path, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, errors.New("did not get a path")
	}

	fullPath := strings.Join([]string{filepath.Dir(path), "testdata", msgFile}, string(filepath.Separator))

	file, err := os.Open(filepath.Clean(fullPath))
	if err != nil {
		return "", nil, err
	}

	defer func() { _ = file.Close() }() // nolint: errcheck

	buf := &bytes.Buffer{}

	_, err = io.Copy(buf, file)
	if err != nil {
		return "", nil, err
	}

	msg, err := service.ParseDIDCommMsgMap(buf.Bytes())
	if err != nil {
		return "", nil, err
	}

	return msg.Type(), buf.Bytes(), err
}

func (s *ControllerSteps) sendMessage(sender, msgFile, recipient string) error {
	url, ok := s.bddContext.GetControllerURL(sender)
	if !ok {
		return fmt.Errorf("unable to find controller URL registered for agent [%s]", sender)
	}

	mt, msg, err := getMsgBytes(msgFile)
	if err != nil {
		return err
	}

	connID := s.bddContext.GetConnectionID(sender, recipient)

	piid, _ := s.actionPIID(sender) // nolint: errcheck

	switch mt {
	case protocol.RequestPresentationMsgTypeV2:
		if piid != "" {
			return sendAcceptProposePresentation(url, piid, msg)
		}

		return sendRequestPresentationMsg(url, connID, msg)
	case protocol.RequestPresentationMsgTypeV3:
		if piid != "" {
			return sendAcceptProposePresentationV3(url, piid, msg)
		}

		return sendRequestPresentationMsgV3(url, connID, msg)
	case protocol.PresentationMsgTypeV2:
		return sendPresentationMsg(url, piid, msg)
	case protocol.PresentationMsgTypeV3:
		return sendPresentationMsgV3(url, piid, msg)
	default:
		return errors.New("message type is not supported")
	}
}

func sendPresentationMsg(url, piid string, msg []byte) error {
	var presentation *presentproof.PresentationV2

	err := json.Unmarshal(msg, &presentation)
	if err != nil {
		return err
	}

	res := &kms.CreateKeySetResponse{}

	err = postToURL(url+"/kms/keyset", kms.CreateKeySetRequest{KeyType: arieskms.BLS12381G2}, res)
	if err != nil {
		return err
	}

	publicKey, err := base64.RawURLEncoding.DecodeString(res.PublicKey)
	if err != nil {
		return err
	}

	didBBS, didKey := fingerprint.CreateDIDKeyByCode(fingerprint.BLS12381g2PubKeyMultiCodec, publicKey)

	for i := range presentation.PresentationsAttach {
		if presentation.PresentationsAttach[i].MimeType != "application/ld+json" {
			continue
		}

		credential, err := presentation.PresentationsAttach[i].Data.Fetch()
		if err != nil {
			return err
		}

		vRes := &verifiable.SignCredentialResponse{}

		signatureRepresentation := docverifiable.SignatureProofValue

		err = postToURL(url+"/verifiable/signcredential", verifiable.SignCredentialRequest{
			Credential: credential,
			DID:        didBBS,
			ProofOptions: &verifiable.ProofOptions{
				KID:                     res.KeyID,
				VerificationMethod:      didKey,
				SignatureRepresentation: &signatureRepresentation,
				SignatureType:           "BbsBlsSignature2020",
			},
		}, vRes)
		if err != nil {
			return err
		}

		presentation.PresentationsAttach[i].Data = decorator.AttachmentData{
			Base64: base64.StdEncoding.EncodeToString(vRes.VerifiableCredential),
		}
	}

	return postToURL(url+fmt.Sprintf(acceptRequestPresentation, piid),
		presentproofcmd.AcceptRequestPresentationV2Args{
			Presentation: presentation,
		}, nil)
}

func sendPresentationMsgV3(url, piid string, msg []byte) error {
	var presentation *presentproof.PresentationV3

	err := json.Unmarshal(msg, &presentation)
	if err != nil {
		return err
	}

	res := &kms.CreateKeySetResponse{}

	err = postToURL(url+"/kms/keyset", kms.CreateKeySetRequest{KeyType: arieskms.BLS12381G2}, res)
	if err != nil {
		return err
	}

	publicKey, err := base64.RawURLEncoding.DecodeString(res.PublicKey)
	if err != nil {
		return err
	}

	didBBS, didKey := fingerprint.CreateDIDKeyByCode(fingerprint.BLS12381g2PubKeyMultiCodec, publicKey)

	for i := range presentation.Attachments {
		if presentation.Attachments[i].MediaType != "application/ld+json" {
			continue
		}

		credential, err := presentation.Attachments[i].Data.Fetch()
		if err != nil {
			return err
		}

		vRes := &verifiable.SignCredentialResponse{}

		signatureRepresentation := docverifiable.SignatureProofValue

		err = postToURL(url+"/verifiable/signcredential", verifiable.SignCredentialRequest{
			Credential: credential,
			DID:        didBBS,
			ProofOptions: &verifiable.ProofOptions{
				KID:                     res.KeyID,
				VerificationMethod:      didKey,
				SignatureRepresentation: &signatureRepresentation,
				SignatureType:           "BbsBlsSignature2020",
			},
		}, vRes)
		if err != nil {
			return err
		}

		presentation.Attachments[i].Data = decorator.AttachmentData{
			Base64: base64.StdEncoding.EncodeToString(vRes.VerifiableCredential),
		}
	}

	return postToURL(url+fmt.Sprintf(acceptRequestPresentationV3, piid),
		presentproofcmd.AcceptRequestPresentationV3Args{
			Presentation: presentation,
		}, nil)
}

func sendRequestPresentationMsg(url, connID string, msg []byte) error {
	var requestPresentation *presentproof.RequestPresentationV2

	err := json.Unmarshal(msg, &requestPresentation)
	if err != nil {
		return err
	}

	return postToURL(url+sendRequestPresentation, presentproofcmd.SendRequestPresentationV2Args{
		ConnectionID:        connID,
		RequestPresentation: requestPresentation,
	}, nil)
}

func sendRequestPresentationMsgV3(url, connID string, msg []byte) error {
	var requestPresentation *presentproof.RequestPresentationV3

	err := json.Unmarshal(msg, &requestPresentation)
	if err != nil {
		return err
	}

	return postToURL(url+sendRequestPresentationV3, presentproofcmd.SendRequestPresentationV3Args{
		ConnectionID:        connID,
		RequestPresentation: requestPresentation,
	}, nil)
}

func sendAcceptProposePresentation(url, piid string, msg []byte) error {
	var requestPresentation *presentproof.RequestPresentationV2

	err := json.Unmarshal(msg, &requestPresentation)
	if err != nil {
		return err
	}

	return postToURL(url+fmt.Sprintf(acceptProposePresentation, piid),
		presentproofcmd.AcceptProposePresentationV2Args{
			RequestPresentation: requestPresentation,
		}, nil)
}

func sendAcceptProposePresentationV3(url, piid string, msg []byte) error {
	var requestPresentation *presentproof.RequestPresentationV3

	err := json.Unmarshal(msg, &requestPresentation)
	if err != nil {
		return err
	}

	return postToURL(url+fmt.Sprintf(acceptProposePresentationV3, piid),
		presentproofcmd.AcceptProposePresentationV3Args{
			RequestPresentation: requestPresentation,
		}, nil)
}

func (s *ControllerSteps) sendProposePresentation(prover, verifier string) error {
	url, ok := s.bddContext.GetControllerURL(prover)
	if !ok {
		return fmt.Errorf("unable to find controller URL registered for agent [%s]", prover)
	}

	return postToURL(url+sendProposalPresentation, presentproofcmd.SendProposePresentationV2Args{
		ConnectionID:        s.bddContext.GetConnectionID(prover, verifier),
		ProposePresentation: &presentproof.ProposePresentationV2{},
	}, nil)
}

func (s *ControllerSteps) sendProposePresentationV3(prover, verifier string) error {
	url, ok := s.bddContext.GetControllerURL(prover)
	if !ok {
		return fmt.Errorf("unable to find controller URL registered for agent [%s]", prover)
	}

	return postToURL(url+sendProposalPresentationV3, presentproofcmd.SendProposePresentationV3Args{
		ConnectionID:        s.bddContext.GetConnectionID(prover, verifier),
		ProposePresentation: &presentproof.ProposePresentationV3{},
	}, nil)
}

func (s *ControllerSteps) negotiateRequestPresentation(agent string) error {
	url, ok := s.bddContext.GetControllerURL(agent)
	if !ok {
		return fmt.Errorf("unable to find controller URL registered for agent [%s]", agent)
	}

	piid, err := s.actionPIID(agent)
	if err != nil {
		return err
	}

	return postToURL(
		url+fmt.Sprintf(negotiateRequestPresentation, piid),
		presentproofcmd.NegotiateRequestPresentationV2Args{
			ProposePresentation: &presentproof.ProposePresentationV2{},
		},
		nil)
}

func (s *ControllerSteps) negotiateRequestPresentationV3(agent string) error {
	url, ok := s.bddContext.GetControllerURL(agent)
	if !ok {
		return fmt.Errorf("unable to find controller URL registered for agent [%s]", agent)
	}

	piid, err := s.actionPIID(agent)
	if err != nil {
		return err
	}

	return postToURL(url+fmt.Sprintf(negotiateRequestPresentationV3, piid),
		presentproofcmd.NegotiateRequestPresentationV3Args{
			ProposePresentation: &presentproof.ProposePresentationV3{},
		}, nil)
}

func (s *ControllerSteps) acceptPresentation(verifier, name string) error {
	return s.acceptPresentationWithRedirect(verifier, name, "")
}

func (s *ControllerSteps) acceptPresentationWithRedirect(verifier, name, redirect string) error {
	url, ok := s.bddContext.GetControllerURL(verifier)
	if !ok {
		return fmt.Errorf("unable to find controller URL registered for agent [%s]", verifier)
	}

	piid, err := s.actionPIID(verifier)
	if err != nil {
		return err
	}

	s.nameToPIID[name] = piid

	return postToURL(url+fmt.Sprintf(acceptPresentation, piid), presentproofcmd.AcceptPresentationArgs{
		Names:       []string{name},
		RedirectURL: redirect,
	}, nil)
}

func (s *ControllerSteps) declinePresentationWithRedirect(verifier, name, prover, redirect string) error {
	url, ok := s.bddContext.GetControllerURL(verifier)
	if !ok {
		return fmt.Errorf("unable to find controller URL registered for agent [%s]", verifier)
	}

	piid, err := s.actionPIID(verifier)
	if err != nil {
		return err
	}

	s.nameToPIID[name] = piid

	err = postToURL(url+fmt.Sprintf(declinePresentation, piid), presentproofcmd.DeclinePresentationArgs{
		RedirectURL: redirect,
	}, nil)
	if err != nil {
		return err
	}

	return s.acceptProblemReport(prover, piid)
}

func (s *ControllerSteps) actionPIID(agentID string) (string, error) {
	msg, err := util.PullEventsFromWebSocket(s.bddContext, agentID, util.FilterTopic("present-proof_actions"))
	if err != nil {
		return "", fmt.Errorf("pull events from WebSocket: %w", err)
	}

	return msg.Message.Properties["piid"].(string), nil
}

func (s *ControllerSteps) waitForMessageFrom(agentID, otherAgent string) error {
	msg, err := util.PullEventsFromWebSocket(s.bddContext, agentID, util.FilterTopic("present-proof_actions"),
		util.FilterMyDID(s.bddContext.PublicDIDs[agentID]))
	if err != nil {
		return fmt.Errorf("pull events from WebSocket: %w", err)
	}

	controllerURL, ok := s.bddContext.GetControllerURL(agentID)
	if !ok {
		return fmt.Errorf(" unable to find controller URL registered for agent [%s]", agentID)
	}

	queryPath := controllerURL + connectionrest.OperationID

	resp := &didexcmd.QueryConnectionsResponse{}

	err = util.SendHTTP(http.MethodGet, queryPath, nil, &resp)
	if err != nil {
		return fmt.Errorf("getting connection record from agent [%s]: %w", agentID, err)
	}

	if len(resp.Results) == 0 {
		return fmt.Errorf("no connection exists for given agent")
	}

	myDID := msg.Message.Properties["myDID"]
	theirDID := msg.Message.Properties["theirDID"]

	var connID string

	for _, conn := range resp.Results {
		if conn != nil && conn.MyDID == myDID && conn.TheirDID == theirDID {
			connID = conn.ConnectionID
			break
		}
	}

	if connID != "" {
		s.bddContext.SaveConnectionID(agentID, otherAgent, connID)
	}

	return nil
}

func (s *ControllerSteps) checkPresentation(verifier, name string) error {
	url, ok := s.bddContext.GetControllerURL(verifier)
	if !ok {
		return fmt.Errorf("unable to find controller URL registered for agent [%s]", verifier)
	}

	msg, err := util.PullEventsFromWebSocket(s.bddContext, verifier,
		util.FilterTopic("present-proof_states"),
		util.FilterStateID("done"),
		util.FilterPIID(s.nameToPIID[name]),
	)
	if err != nil {
		return fmt.Errorf("pull events from WebSocket: %w", err)
	}

	if !reflect.DeepEqual(msg.Message.Properties["names"], []interface{}{name}) {
		return fmt.Errorf("properties: expected names [%s], got %v", name,
			msg.Message.Properties["names"])
	}

	var result verifiable.RecordResult
	if err := util.SendHTTP(http.MethodGet, url+verifiablePresentations, nil, &result); err != nil {
		return err
	}

	for _, val := range result.Result {
		if val.Name == name {
			return nil
		}
	}

	return errors.New("presentation not found")
}

func (s *ControllerSteps) acceptProblemReport(agent, piid string) error {
	_, err := util.PullEventsFromWebSocket(s.bddContext, agent,
		util.FilterTopic("present-proof_actions"),
		util.FilterPIID(piid),
	)
	if err != nil {
		return fmt.Errorf("pull events from WebSocket: %w", err)
	}

	url, ok := s.bddContext.GetControllerURL(agent)
	if !ok {
		return fmt.Errorf("unable to find controller URL registered for agent [%s]", agent)
	}

	return postToURL(url+fmt.Sprintf(acceptProblemReport, piid), presentproofcmd.AcceptProblemReportArgs{}, nil)
}

func (s *ControllerSteps) validateState(agent, state, redirect, status string) error {
	msg, err := util.PullEventsFromWebSocket(s.bddContext, agent,
		util.FilterTopic("present-proof_states"),
		util.FilterStateID(state),
	)
	if err != nil {
		return fmt.Errorf("pull events from WebSocket: %w", err)
	}

	if redirect != msg.Message.Properties["url"] {
		return fmt.Errorf("failed redirect URL validation, expected[%s]: found[%s]",
			redirect, msg.Message.Properties["url"])
	}

	return nil
}

func (s *ControllerSteps) agentDID(ds *didexsteps.ControllerSteps, agent string) (string, error) {
	connectionID, ok := ds.ConnectionIDs()[agent]
	if !ok {
		return "", fmt.Errorf("unable to find connection for agent [%s]", agent)
	}

	controllerURL, ok := s.bddContext.GetControllerURL(agent)
	if !ok {
		return "", fmt.Errorf("unable to find controller URL registered for agent [%s]", agent)
	}

	var response didexcmd.QueryConnectionResponse

	err := util.SendHTTP(http.MethodGet, fmt.Sprintf("%s/connections/%s", controllerURL, connectionID),
		nil, &response)
	if err != nil {
		return "", fmt.Errorf("failed to query connections: %w", err)
	}

	return response.Result.MyDID, nil
}

func postToURL(url string, payload, resp interface{}) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	return util.SendHTTP(http.MethodPost, url, body, &resp)
}

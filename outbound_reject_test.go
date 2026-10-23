package meowcaller

import (
	"testing"

	"github.com/polymorfa/hypermeow/types"
	"github.com/polymorfa/hypermeow/types/events"
)

func deviceJID(device uint16) types.JID {
	return types.JID{User: "70790103990472", Server: types.HiddenUserServer, Device: device}
}

// Issue #25: on an outgoing call to a multi-device account, a linked Web/Desktop
// device that never rang auto-rejects while the phone (which preaccepted) rings on.
// That sibling reject must not end the call.
func TestSiblingDeviceRejectKeepsOutgoingCallLive(t *testing.T) {
	eng, call := testEngineWithOutgoingCall()
	eng.onPreAccept(&events.CallPreAccept{BasicCallMeta: types.BasicCallMeta{CallID: "CID", From: deviceJID(0)}})

	eng.onReject(&events.CallReject{BasicCallMeta: types.BasicCallMeta{CallID: "CID", From: deviceJID(22)}})
	if got := call.State(); got == CallPhaseEnded || eng.lookup("CID") == nil {
		t.Fatalf("a sibling device's reject ended the call (phase %d)", got)
	}

	// The device that is actually ringing declines: the call ends.
	eng.onReject(&events.CallReject{BasicCallMeta: types.BasicCallMeta{CallID: "CID", From: deviceJID(0)}})
	if got := call.State(); got != CallPhaseEnded {
		t.Fatalf("phase = %d after the ringing device rejected, want Ended", got)
	}
}

// Without any preaccept a reject still ends the call (busy, blocked, burst rejects).
func TestRejectBeforeAnyPreacceptEndsOutgoingCall(t *testing.T) {
	eng, call := testEngineWithOutgoingCall()
	eng.onReject(&events.CallReject{BasicCallMeta: types.BasicCallMeta{CallID: "CID", From: deviceJID(22)}})
	if got := call.State(); got != CallPhaseEnded {
		t.Fatalf("phase = %d, want Ended", got)
	}
}

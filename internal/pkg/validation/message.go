package validation

import (
	"github.com/filecoin-project/chain-validation/pkg/chain"
	"github.com/filecoin-project/chain-validation/pkg/state"
	"github.com/pkg/errors"

	"github.com/filecoin-project/go-filecoin/internal/pkg/types"
	"github.com/filecoin-project/go-filecoin/internal/pkg/vm/abi"
	"github.com/filecoin-project/go-filecoin/internal/pkg/vm/address"
)

// MessageFactory implements the MessageFactory interface.
type MessageFactory struct {
	signer types.Signer
}

var _ chain.MessageFactory = &MessageFactory{}

// NewMessageFactory returns a new MessageFactory.
func NewMessageFactory(signer types.Signer) *MessageFactory {
	return &MessageFactory{signer}
}

// MakeMessage creates and returns a message.
func (mf *MessageFactory) MakeMessage(from, to state.Address, method chain.MethodID, nonce uint64,
	value, gasPrice state.AttoFIL, gasUnit state.GasUnit, params ...interface{}) (interface{}, error) {
	fromDec, err := address.NewFromBytes([]byte(from))
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create address from: %s", from)
	}
	toDec, err := address.NewFromBytes([]byte(to))
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create address from: %s", toDec)
	}
	valueDec := types.NewAttoFIL(value)
	paramsDec, err := abi.ToEncodedValues(params)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to encode abi param values: %v", params)
	}
	if int(method) >= len(methods) {
		return nil, errors.Errorf("No method name for method %v", method)
	}
	methodName := methods[method]
	msg := types.NewMeteredMessage(fromDec, toDec, nonce, valueDec, methodName, paramsDec, types.NewAttoFIL(gasPrice), types.NewGasUnits(uint64(gasUnit)))

	return msg, nil
	//return types.NewSignedMessage(*msg, mf.signer)
}

// FromSingletonAddress returns the address for a Singleton Actor.
func (mf *MessageFactory) FromSingletonAddress(addr state.SingletonActorID) state.Address {
	return fromSingletonAddress(addr)
}

// Maps method enumeration values to method names.
// This will change to a mapping to method ids when method dispatch is updated to use integers.
var methods = []types.MethodID{
	chain.NoMethod: 0,
}

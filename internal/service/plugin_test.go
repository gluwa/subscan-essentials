package service

import (
	"testing"

	"github.com/itering/subscan/model"
	"github.com/shopspring/decimal"
)

func Test_emitEvent(t *testing.T) {
	testSrv.emitEvent(&testBlock, &testEvent, decimal.Zero, nil)
}

func Test_emitExtrinsic(t *testing.T) {
	testSrv.emitExtrinsic(&testBlock, &testSignedExtrinsic, []model.ChainEvent{testEvent})
}

package rpc

import (
	"context"
	"math/big"
	"encoding/json"

	grpc "github.com/ethereum/go-ethereum/rpc"
)

type Client struct {
	ptr *grpc.Client
}

func (c *Client) Close() {
	c.ptr.Close()
}

func (c *Client) GetBlockByHash(
	ctx context.Context, blockHash BlockHash, requestedScope RequestedScope,
) (*BlockResponse, error) {
	var res BlockResponse

	err := c.ptr.CallContext(
		ctx, &res, "starknet_getBlockByHash", blockHash, requestedScope)

	if err != nil {
		return nil, err
	}

	return &res, err
}

func (ec *Client) GetBlockByNumber(
	ctx context.Context, number *big.Int, requestedScope RequestedScope,
) (*BlockResponse, error) {
	var res BlockResponse

	err := ec.ptr.CallContext(
		ctx, &res, "starknet_getBlockByNumber", number.Uint64(), requestedScope)

	if err != nil {
		return nil, err
	}

	return &res, err
}

func (ec *Client) GetTransactionByHash(
	ctx context.Context, hash TransactionHash,
) (*Transaction, error) {
	var res Transaction

	err := ec.ptr.CallContext(
		ctx, &res, "starknet_getTransactionByHash", hash)

	if err != nil {
		return nil, err
	}

	return &res, err
}

func (ec *Client) Call(
	ctx context.Context, hash string, call FunctionCall, 
) ([]string, error) {
	var res []string

	if len(call.CallData) == 0 {
		call.CallData = make([]string, 0)
	}

	if err := ec.execute(ctx, "starknet_call", &res, call, hash); err != nil {
		return nil, err
	}

	return res, nil
}

func (ec *Client) execute(ctx context.Context, method string, data interface{}, args ...interface{}) error {
	var raw json.RawMessage

	err := ec.ptr.CallContext(ctx, &raw, method, args...)

	if err != nil {
		return err
	}

	if err := json.Unmarshal(raw, &data); err != nil {
		return err
	}

	return nil
}

func NewClient(ptr *grpc.Client) *Client {
	return &Client{ptr}
}

func dialContext(ctx context.Context, rawUrl string) (*Client, error) {
	c, err := grpc.DialContext(ctx, rawUrl)
	if err != nil {
		return nil, err
	}
	return NewClient(c), nil
}

func Dial(rawUrl string) (*Client, error) {
	return dialContext(context.Background(), rawUrl)
}

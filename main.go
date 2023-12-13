package main

import (
	"context"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Chain struct {
	GrpcConn *grpc.ClientConn
}

func main() {
	registry := codectypes.NewInterfaceRegistry()
	cryptocodec.RegisterInterfaces(registry)
	cdc := codec.NewProtoCodec(registry)

	chain, err := mainE()
	if err != nil {
		panic(err)
	}

	// propId := 1 // text prop
	propId := 858 // community spend prop
	resp, err := chain.GetProposalApi(context.Background(), uint64(propId))
	if err != nil {
		panic(err)
	}

	content := resp.Proposal.Content
	fmt.Println(content.TypeUrl)

	switch content.TypeUrl {
	case "/cosmos.gov.v1beta1.TextProposal":
		var textProp govv1beta1.TextProposal
		err = cdc.Unmarshal(content.Value, &textProp)
		fmt.Println(textProp.Title, textProp.Description)
	case "/cosmos.distribution.v1beta1.CommunityPoolSpendProposal":
		var commPoolProp distrtypes.CommunityPoolSpendProposal
		err = cdc.Unmarshal(content.Value, &commPoolProp)
		fmt.Println(commPoolProp.Title, commPoolProp.Amount, commPoolProp.Recipient)
	}

}

func mainE() (*Chain, error) {
	endpoint := "cosmos-grpc.polkachu.com:14990"

	grpcConn, err := grpc.Dial(
		endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("grpc dial: %w", err)
	}

	return &Chain{
		GrpcConn: grpcConn,
	}, nil
}

func (c *Chain) GetProposalApi(ctx context.Context, proposalId uint64) (*govv1beta1.QueryProposalResponse, error) {
	queryClient := govv1beta1.NewQueryClient(c.GrpcConn)

	resp, err := queryClient.Proposal(
		ctx,
		&govv1beta1.QueryProposalRequest{
			ProposalId: proposalId,
		},
	)

	return resp, err
}

package main

import (
	"context"
	"fmt"

	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Chain struct {
	GrpcConn *grpc.ClientConn
}

func main() {
	cdc := SetupRegistry()

	hubChain, err := createChain("cosmos-grpc.polkachu.com:14990")
	if err != nil {
		panic(err)
	}

	// propId := 1 // text prop
	propId := 858 // community spend prop

	resp, err := hubChain.GetProposal(context.Background(), uint64(propId))
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

	if err != nil {
		panic(err)
	}

}

func createChain(endpoint string) (*Chain, error) {
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

func (c *Chain) GetProposal(ctx context.Context, proposalId uint64) (*govv1beta1.QueryProposalResponse, error) {
	queryClient := govv1beta1.NewQueryClient(c.GrpcConn)

	resp, err := queryClient.Proposal(
		ctx,
		&govv1beta1.QueryProposalRequest{
			ProposalId: proposalId,
		},
	)

	return resp, err
}

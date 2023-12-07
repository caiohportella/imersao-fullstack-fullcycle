package transformer

import (
	"github.com/caiohportella/imersao-fullstack-fullcycle/go/internal/market/dto"
	"github.com/caiohportella/imersao-fullstack-fullcycle/go/internal/market/entities"
)

func TransformInput(input dto.TradeInput) *entities.Order {
	asset := entities.NewAsset(input.AssetID, input.AssetID, 1000)
	investor := entities.NewInvestor(input.InvestorID, input.InvestorID)
	order := entities.NewOrder(input.OrderID, investor, asset, input.Shares, input.Price, input.OrderType)

	if (input.CurrentShares > 0) {
		assetPosition := entities.NewInvestorAssetPosition(input.AssetID, input.CurrentShares)
		investor.AddAssetPosition(assetPosition)
	}

	return order
}

func TransformOutput(order *entities.Order) *dto.OrderOutput {
	output := &dto.OrderOutput{
		OrderID: order.ID,
		InvestorID: order.Investor.ID,
		AssetID: order.Asset.ID,
		OrderType: order.OrderType,
		Status: order.Status,
		Partial: order.PendingShares,
		Shares: order.Shares,
	}

	var transactionsOutput []*dto.TransactionOutput

	for _, t := range order.Transactions {
		transactionOutput := &dto.TransactionOutput{
			TransactionID: t.ID,
			BuyerID: t.BuyingOrder.ID,
			SellerID: t.SellingOrder.ID,
			AssetID: t.SellingOrder.Asset.ID,
			Price: t.Price,
			Shares: t.SellingOrder.Shares - t.SellingOrder.PendingShares,
		}

		transactionsOutput = append(transactionsOutput, transactionOutput)	
	}

	output.TransactionOutput = transactionsOutput

	return output
}
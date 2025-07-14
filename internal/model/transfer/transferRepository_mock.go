package transfer

import (
	"context"

	"BankingAPI/internal/model"

	"github.com/labstack/echo/v4"
)

var singleton *TransferRepository

type MockTransferRepository struct {
	TransferMap *map[string]Transfer
}

func NewMockTransferRepository() TransferRepository {
	if singleton != nil {
		return *singleton
	}
	userMap := make(map[string]Transfer)
	*singleton = MockTransferRepository{
		TransferMap: &userMap,
	}
	return *singleton
}

func (t MockTransferRepository) Create(ctx context.Context, request *Transfer) (*Transfer, *echo.HTTPError) {
	/*
		for {
			transferID := randomstring.String(10)
			if _, ok := (*m.TransferMap)[tranferID]; !ok {
				transferRequest.Transfer_id = transferID
				(*t.TransferMap)[transferRequest.Transfer_id] = *transferRequest
				break
			}
		}
	*/
	(*t.TransferMap)[request.Transfer_id] = *request
	return request, nil
}

func (t MockTransferRepository) Delete(ctx context.Context, id *string) *echo.HTTPError {
	if _, ok := (*t.TransferMap)[*id]; !ok {
		return model.ErrIDnotFound
	}
	delete(*t.TransferMap, *id)
	return nil
}

func (t MockTransferRepository) Get(ctx context.Context, id *string) (*Transfer, *echo.HTTPError) {
	if transfer, ok := (*t.TransferMap)[*id]; !ok {
		return nil, model.ErrIDnotFound
	} else {
		return &transfer, nil
	}
}

func (t MockTransferRepository) Update(ctx context.Context, request *Transfer) *echo.HTTPError {
	if _, ok := (*t.TransferMap)[request.Transfer_id]; !ok {
		return model.ErrIDnotFound
	}
	(*t.TransferMap)[request.Transfer_id] = *request
	return nil
}

func (t MockTransferRepository) GetAll(ctx context.Context) (*[]Transfer, *echo.HTTPError) {
	if len(*t.TransferMap) == 0 {
		return nil, model.ErrIDnotFound
	}
	users := make([]Transfer, 0, len(*t.TransferMap))
	for _, user := range *t.TransferMap {
		users = append(users, user)
	}
	return &users, nil
}

func (db MockTransferRepository) GetFilteredByAccountID(ctx context.Context, accountID *string) (*[]Transfer, *echo.HTTPError) {
	if accountID == nil || len(*accountID) == 0 {
		return nil, model.ErrFilterNotSet
	}
	transferSlice := make([]Transfer, 0, len(*db.TransferMap))

	for _, transfer := range *db.TransferMap {
		if transfer.Account_id == *accountID {
			transferSlice = append(transferSlice, transfer)
		}
	}
	/* for _, transferResponse := range *db.TransferMap {
		match := true
		for _, filter := range *filters {
			token := model.TokenizeFilters(&filter)
			if len(*token) != 3 {
				return nil, model.InvalidFilterFormat
			}
			field := (*token)[0]
			operator := (*token)[1]
			value := (*token)[2]

			switch field {
			case "account_id":
				if operator == "==" && transferResponse.Account_id != value {
					match = false
				}
			case "account_to":
				if operator == "==" && transferResponse.Account_to != value {
					match = false
				}
			// Add more fields as necessary
			default:
				match = false
			}
		}
		if match {
			transferSlice = append(transferSlice, transferResponse)
		}
	}
	*/
	return &transferSlice, nil
}

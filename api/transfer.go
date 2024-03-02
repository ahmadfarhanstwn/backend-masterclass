package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	db "github.com/ahmadfarhanstwn/backend-masterclass/db/sqlc"
	"github.com/ahmadfarhanstwn/backend-masterclass/token"
	"github.com/gin-gonic/gin"
)

type CreateTransferRequests struct {
	FromAccountId int64  `json:"from_account_id" binding:"required,gt=0"`
	ToAccountId   int64  `json:"to_account_id" binding:"required,gt=0"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	Currency      string `json:"currency" binding:"required,currency"`
}

func (server *Server) createTransfer(ctx *gin.Context) {
	var req CreateTransferRequests
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	fromAccount, valid := server.validAccount(ctx, req.FromAccountId, req.Currency)
	if !valid {
		return
	}

	_, valid = server.validAccount(ctx, req.ToAccountId, req.Currency)
	if !valid {
		return
	}

	authPayload := ctx.MustGet(authPayloadKey).(*token.Payload)
	if fromAccount.Owner != authPayload.Username {
		err := errors.New("from account doesn't belong to authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	arg := db.TransferTxParams{
		FromAccountId: req.FromAccountId,
		ToAccountId:   req.ToAccountId,
		Amount:        req.Amount,
	}

	result, err := server.store.TransferTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, result)
}

func (server *Server) validAccount(ctx *gin.Context, accountId int64, currency string) (db.Account, bool) {
	account, err := server.store.GetAccountForUpdate(ctx, accountId)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return account, false
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return account, false
	}

	if account.Currency != currency {
		err := fmt.Errorf("account [%d] currency mismatch : %s vs %s", accountId, account.Currency, currency)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return account, false
	}

	return account, true
}

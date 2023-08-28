package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	db "github.com/qwetu_petro/backend/db/sqlc"
	"github.com/qwetu_petro/backend/utils"
	"github.com/qwetu_petro/backend/workers"
	"strconv"
	"time"
)

type createCarReq struct {
	Make         string `json:"make"`
	Model        string `json:"model"`
	Year         int32  `json:"year"`
	LicensePlate string `json:"license_plate"`
}

type createCarRes struct {
	ID int32 `json:"id"`
}

type createFuelConsumptionReq struct {
	ConsumptionID int32     `json:"consumption_id"`
	CarID         int32     `json:"car_id"`
	LitersOfFuel  float64   `json:"liters_of_fuel"`
	CostInKsh     float64   `json:"cost_in_ksh"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type createFuelConsumptionRes struct {
	ConsumptionID int32  `json:"consumption_id"`
	Message       string `json:"message"`
}

func (server *Server) createCar(ctx *gin.Context) {
	var req createCarReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}

	arg := db.CreateCarParams{
		Make:         req.Make,
		Model:        req.Model,
		Year:         req.Year,
		LicensePlate: req.LicensePlate,
	}

	car, err := server.store.CreateCar(ctx, arg)
	if err != nil {
		ctx.JSON(500, errorResponse(err))
		return
	}
	resp := createCarRes{
		ID: car.ID,
	}

	ctx.JSON(200, resp)

}

type listCarsReq struct {
	PageSize int32 `form:"page_size"`
	PageID   int32 `form:"page_id"`
}

func (server *Server) listCars(ctx *gin.Context) {
	var req listCarsReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}

	arg := db.GetCarsParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	cars, err := server.store.GetCars(ctx, arg)
	if err != nil {
		ctx.JSON(500, errorResponse(err))
		return
	}

	ctx.JSON(200, cars)

}

func (server *Server) createFuelConsumption(ctx *gin.Context) {
	var req createFuelConsumptionReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}

	car, err := server.store.GetCarById(ctx, req.CarID)
	if err != nil {
		ctx.JSON(500, errorResponse(err))
		return
	}
	// get user id from context
	userId, err := getUserIdFromContext(ctx)
	if err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}

	fuelTime, err := utils.ShortenTimestamp(time.Now())
	if err != nil {
		ctx.JSON(500, errorResponse(err))
		return
	}

	cost := strconv.FormatInt(int64(req.CostInKsh), 10)

	amount, err := utils.StringToNumeric(cost)

	pettyCashArgs := db.CreatePettyCashParams{
		EmployeeID:      userId,
		Amount:          amount,
		Folio:           "PETTY CASH",
		Description:     fmt.Sprintf("Fuel for %s number plate %s on date %s", car.Make, car.LicensePlate, fuelTime),
		DebitAccount:    "Petty Cash",
		TransactionDate: time.Now(),
	}

	fuelPettyCash, err := server.store.CreatePettyCashWithAudit(ctx, userId, pettyCashArgs)
	if err != nil {
		ctx.JSON(500, errorResponse(err))
		return

	}

	arg := db.CreateFuelConsumptionParams{
		CarID:         req.CarID,
		LitersOfFuel:  &req.LitersOfFuel,
		CostInKsh:     req.CostInKsh,
		FuelDate:      time.Now(),
		TransactionID: &fuelPettyCash.TransactionID,
	}

	fuel, err := server.store.CreateFuelConsumption(ctx, arg)
	if err != nil {
		ctx.JSON(500, errorResponse(err))
		return
	}

	resp := createFuelConsumptionRes{
		ConsumptionID: fuel.ConsumptionID,
		Message:       "Fuel consumption created successfully",
	}

	// create task payload
	payload := &workers.CreatePettyCashPayload{
		PettyCashID: fuelPettyCash.TransactionID,
	}

	// create task
	err = server.taskDistributor.DistributeTaskCreatePettyCashPdf(ctx, payload)
	if err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}

	ctx.JSON(200, resp)

}

type listFuelConsumptionReq struct {
	PageSize int32 `form:"page_size"`
	PageID   int32 `form:"page_id"`
}

func (server *Server) listFuelConsumption(ctx *gin.Context) {
	var req listFuelConsumptionReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}

	arg := db.GetConsumptionsParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	fuel, err := server.store.GetConsumptions(ctx, arg)
	if err != nil {
		ctx.JSON(500, errorResponse(err))
		return
	}

	ctx.JSON(200, fuel)

}

type GetCarFuelByDateRangeReq struct {
	CarID     int32  `form:"id" binding:"required"`
	StartDate string `form:"startDate" binding:"required"`
	EndDate   string `form:"endDate" binding:"required"`
}

func (server *Server) getCarFuelByDateRange(ctx *gin.Context) {
	var req GetCarFuelByDateRangeReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}

	// Convert strings to time.Time
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}

	arg := db.GetConsumptionByCarAndDateRangeParams{
		FuelDate:   startDate,
		FuelDate_2: endDate,
		CarID:      req.CarID,
	}

	fuel, err := server.store.GetConsumptionByCarAndDateRange(ctx, arg)
	if err != nil {
		ctx.JSON(500, errorResponse(err))
		return
	}

	ctx.JSON(200, fuel)

}

type GetCarFuelConsumptionReq struct {
	CarID int32 `uri:"id"`
}

func (server *Server) getCarFuelConsumption(ctx *gin.Context) {
	var req GetCarFuelConsumptionReq
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}
	fuel, err := server.store.GetConsumptionByCar(ctx, req.CarID)
	if err != nil {
		ctx.JSON(500, errorResponse(err))
		return
	}

	ctx.JSON(200, fuel)

}

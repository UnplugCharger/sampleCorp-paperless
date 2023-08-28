-- name: CreateCar :one
INSERT INTO cars (
make,
model,
year,
license_plate
) VALUES (
$1, $2, $3, $4
) RETURNING *;

-- name: CreateFuelConsumption :one
INSERT INTO car_fuel_consumption (
car_id,
liters_of_fuel,
cost_in_ksh,
fuel_date,
transaction_id
) VALUES (
 $1, $2, $3, $4, $5
) RETURNING *;

-- name: GetCars :many
SELECT * FROM cars LIMIT $1 OFFSET $2;

-- name: GetCarById :one
SELECT * FROM cars WHERE id = $1;

-- name: GetConsumptions :many
SELECT * FROM car_fuel_consumption LIMIT $1 OFFSET $2;

-- name: GetConsumptionByCar :many
SELECT * FROM car_fuel_consumption WHERE car_id = $1;

-- name: GetConsumptionByCarAndDateRange :many
SELECT * FROM car_fuel_consumption WHERE car_id = $1 AND fuel_date BETWEEN $2 AND $3;

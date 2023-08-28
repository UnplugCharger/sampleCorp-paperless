CREATE TABLE cars (
id SERIAL PRIMARY KEY,
make VARCHAR(50) NOT NULL,
model VARCHAR(50) NOT NULL,
year INTEGER NOT NULL,
license_plate VARCHAR(20) NOT NULL,
updated_at DATE NOT NULL DEFAULT CURRENT_DATE
);

CREATE TABLE car_fuel_consumption (
consumption_id SERIAL PRIMARY KEY,
car_id INTEGER NOT NULL REFERENCES cars(id), -- Assuming there is a 'cars' table with car details
liters_of_fuel FLOAT NULL,
cost_in_ksh FLOAT NOT NULL,
fuel_date DATE NOT NULL,
transaction_id INTEGER UNIQUE REFERENCES petty_cash(transaction_id),
updated_at DATE NOT NULL DEFAULT CURRENT_DATE
);

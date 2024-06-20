-- +goose Up
-- +goose StatementBegin
CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    date VARCHAR(10) NOT NULL,
    amount NUMERIC(10, 2) NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE transactions;
-- +goose StatementEnd

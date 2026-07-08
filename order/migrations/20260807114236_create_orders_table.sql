-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS orders (
    order_uuid       UUID PRIMARY KEY,
    user_uuid        UUID NOT NULL,
    part_uuids       UUID[] NOT NULL,
    total_price      DOUBLE PRECISION NOT NULL,
    status           VARCHAR(50) NOT NULL,
    transaction_uuid UUID,
    payment_method   VARCHAR(50)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS orders;
-- +goose StatementEnd
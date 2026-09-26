-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS parts (
    part_uuid   UUID PRIMARY KEY,
    name        VARCHAR(128) NOT NULL,
    price       DOUBLE PRECISION NOT NULL CHECK (price >= 0),
    category    VARCHAR(64) NOT NULL DEFAULT 'GENERAL',
    in_stock    INT NOT NULL DEFAULT 0 CHECK (in_stock >= 0),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS order_items (
    order_uuid  UUID NOT NULL REFERENCES orders(order_uuid) ON DELETE CASCADE,
    part_uuid   UUID NOT NULL REFERENCES parts(part_uuid),
    quantity    INT NOT NULL DEFAULT 1 CHECK (quantity > 0),
    unit_price  DOUBLE PRECISION NOT NULL CHECK (unit_price >= 0),
    PRIMARY KEY (order_uuid, part_uuid)
);

CREATE TABLE IF NOT EXISTS payments (
    payment_uuid UUID PRIMARY KEY,
    order_uuid   UUID NOT NULL REFERENCES orders(order_uuid) ON DELETE CASCADE,
    amount       DOUBLE PRECISION NOT NULL,
    method       VARCHAR(50) NOT NULL,
    status       VARCHAR(50) NOT NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_order_items_part ON order_items(part_uuid);
CREATE INDEX IF NOT EXISTS idx_payments_order ON payments(order_uuid);

INSERT INTO parts (part_uuid, name, price, category, in_stock) VALUES
    ('22222222-2222-2222-2222-222222222222', 'Ion Thruster', 100, 'ENGINE', 50),
    ('33333333-3333-3333-3333-333333333333', 'Plasma Shield', 250, 'DEFENSE', 20),
    ('44444444-4444-4444-4444-444444444444', 'Hyperdrive Core', 500, 'ENGINE', 10),
    ('55555555-5555-5555-5555-555555555555', 'Cargo Bay Module', 150, 'HULL', 30)
ON CONFLICT (part_uuid) DO NOTHING;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS payments;
DROP TABLE IF EXISTS order_items;
DROP TABLE IF EXISTS parts;
-- +goose StatementEnd

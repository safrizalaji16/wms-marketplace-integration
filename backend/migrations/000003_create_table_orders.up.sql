CREATE TABLE orders (
    id                       UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_sn                 TEXT NOT NULL UNIQUE,
    shop_id                  TEXT NOT NULL,
    marketplace_status       TEXT NOT NULL,
    shipping_status          TEXT NOT NULL,
    wms_status               TEXT NOT NULL DEFAULT 'READY_TO_PICK',
    tracking_number          TEXT,
    total_amount             NUMERIC(14, 2) NOT NULL,
    raw_marketplace_payload  JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at               TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at               TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT orders_wms_status_check CHECK (
        wms_status IN ('READY_TO_PICK', 'PICKING', 'PACKED', 'SHIPPED')
    )
);

CREATE INDEX idx_orders_shop_id ON orders(shop_id);
CREATE INDEX idx_orders_wms_status ON orders(wms_status);
CREATE INDEX idx_orders_updated_at ON orders(updated_at DESC);

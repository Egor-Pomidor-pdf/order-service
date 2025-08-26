CREATE TABLE orders (
    order_uid TEXT PRIMARY KEY,
    track_number TEXT,
    entry TEXT,
    locale TEXT,
    customer_id TEXT,
    internal_signature TEXT,
    delivery_service TEXT,
    shardkey TEXT,
    sm_id INT,
    date_created TIMESTAMP,
    oof_shard TEXT
    
);

CREATE TABLE deliveries (
    id SERIAL PRIMARY KEY,
    order_uid TEXT NOT NULL,
    name TEXT,
    phone TEXT,
    zip TEXT,
    city TEXT,
    address TEXT,
    region TEXT,
    email TEXT,
    CONSTRAINT fk_deliveries_orders
        FOREIGN KEY(order_uid) REFERENCES orders(order_uid)
        ON DELETE CASCADE
);

CREATE TABLE payments (
    id SERIAL PRIMARY KEY,
    order_uid TEXT NOT NULL,
    transaction TEXT,
    request_id TEXT,
    currency TEXT,
    provider TEXT,
    amount INT,
    payment_dt BIGINT,
    bank TEXT,
    delivery_cost INT,
    goods_total INT,
    custom_fee INT,
    CONSTRAINT fk_payments_orders
        FOREIGN KEY(order_uid) REFERENCES orders(order_uid)
        ON DELETE CASCADE
);

CREATE TABLE items (
    id SERIAL PRIMARY KEY,
    order_uid TEXT NOT NULL,
    chrt_id BIGINT,
    track_number TEXT,
    price INT,
    rid TEXT,
    name TEXT,
    sale INT,
    size TEXT,
    total_price INT,
    nm_id BIGINT,
    brand TEXT,
    status INT,
    CONSTRAINT fk_items_orders
        FOREIGN KEY(order_uid) REFERENCES orders(order_uid)
        ON DELETE CASCADE
);
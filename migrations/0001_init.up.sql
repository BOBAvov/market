CREATE TABLE IF NOT EXISTS balance_users (
                                             id       BIGSERIAL PRIMARY KEY,
                                             balance  BIGINT NOT NULL DEFAULT 0 CHECK (balance >= 0)
    );

-- Пользователи
CREATE TABLE IF NOT EXISTS users (
                                     id            BIGSERIAL PRIMARY KEY,
                                     email         TEXT NOT NULL,
                                     password_hash TEXT NOT NULL,
                                     role          TEXT NOT NULL CHECK (role IN ('buyer', 'seller')),
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    balance_id    BIGINT REFERENCES balance_users(id) ON DELETE SET NULL
    );

-- Уникальность email без учета регистра
CREATE UNIQUE INDEX IF NOT EXISTS ux_users_email_lower ON users ((lower(email)));
CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);
CREATE INDEX IF NOT EXISTS idx_users_balance_id ON users(balance_id);

-- Картинки (лучше хранить метаданные и дату создания)
CREATE TABLE IF NOT EXISTS pictures (
                                        id          BIGSERIAL PRIMARY KEY,
                                        data        BYTEA NOT NULL,
                                        mime_type   TEXT,
                                        size_bytes  BIGINT CHECK (size_bytes IS NULL OR size_bytes >= 0),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
    );

-- Товары
CREATE TABLE IF NOT EXISTS products (
                                        id             BIGSERIAL PRIMARY KEY,
                                        seller_id      BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name           TEXT NOT NULL CHECK (length(trim(name)) > 0),
    description    TEXT,
    price_cents    BIGINT NOT NULL CHECK (price_cents >= 0),
    stock          INTEGER NOT NULL DEFAULT 0 CHECK (stock >= 0),
    cover_picture_id BIGINT REFERENCES pictures(id) ON DELETE SET NULL,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
    );

-- Связь «товар — картинки» с порядком
CREATE TABLE IF NOT EXISTS product_pictures (
                                                product_id  BIGINT NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    picture_id  BIGINT NOT NULL REFERENCES pictures(id) ON DELETE CASCADE,
    position    INTEGER NOT NULL CHECK (position > 0),
    PRIMARY KEY (product_id, picture_id),
    UNIQUE (product_id, position)
    );

-- Индексы по товарам
CREATE INDEX IF NOT EXISTS idx_products_seller_id ON products(seller_id);
-- Оставляйте этот индекс, если нужны поиск по имени БЕЗ фильтра по продавцу.
-- Если таких запросов нет — можно опустить как дублирующий нагрузку.
CREATE INDEX IF NOT EXISTS idx_products_name_lower ON products ((lower(name)));
-- Гарантия уникальности названий в рамках одного продавца
CREATE UNIQUE INDEX IF NOT EXISTS ux_products_seller_name ON products (seller_id, (lower(name)));

-- Триггер для updated_at
CREATE OR REPLACE FUNCTION set_updated_at() RETURNS trigger LANGUAGE plpgsql AS $$
BEGIN
  NEW.updated_at := NOW();
RETURN NEW;
END;
$$;

DROP TRIGGER IF EXISTS trg_products_set_updated_at ON products;
CREATE TRIGGER trg_products_set_updated_at
    BEFORE UPDATE ON products
    FOR EACH ROW
    EXECUTE FUNCTION set_updated_at();
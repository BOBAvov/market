-- Сначала триггеры/функции
DROP TRIGGER IF EXISTS trg_products_set_updated_at ON products;
DROP FUNCTION IF EXISTS set_updated_at;

-- Дочерние таблицы
DROP TABLE IF EXISTS product_pictures;

-- Основные таблицы (индексы удалятся автоматически)
DROP TABLE IF EXISTS products;
DROP TABLE IF EXISTS pictures;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS balance_users;
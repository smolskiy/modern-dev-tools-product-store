-- Создание таблицы products
CREATE TABLE IF NOT EXISTS products (
                                        id SERIAL PRIMARY KEY,
                                        name VARCHAR(100),
    description TEXT,
    price NUMERIC
    );

-- Вставка тестовых данных
INSERT INTO products (name, description, price) VALUES
                                                    ('Product 1', 'Description 1', 100.0),
                                                    ('Product 2', 'Description 2', 200.0),
                                                    ('Product 3', 'Description 3', 300.0);

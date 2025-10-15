CREATE TABLE person (
	id SERIAL PRIMARY KEY,
	name TEXT NOT NULL,
	email TEXT NOT NULL,
	password_hash VARCHAR(50),
	created_at DATE DEFAULT NOW()
);

CREATE TABLE category (
	id SERIAL PRIMARY KEY,
	name TEXT NOT NULL
);

CREATE TABLE product (
	id SERIAL PRIMARY KEY,
	name TEXT NOT NULL,
	description TEXT,
	price DECIMAL(10, 2) NOT NULL DEFAULT 0.0,
	amount INTEGER NOT NULL DEFAULT 0,
	category_id INTEGER,
	created_at DATE DEFAULT NOW(),
	FOREIGN KEY (category_id) REFERENCES category(id)
)

-- Вставляем категории
INSERT INTO category (name) VALUES
('Смартфоны'),
('Ноутбуки'),
('Наушники'),
('Планшеты'),
('Умные часы');

-- Вставляем пользователей
INSERT INTO person (name, email, password_hash, created_at) VALUES
('Иван Петров', 'ivan@mail.ru', 'hash123', NOW()),
('Мария Сидорова', 'maria@yandex.ru', 'hash456', NOW()),
('Алексей Козлов', 'alex@google.com', 'hash789', NOW()),
('Елена Новикова', 'elena@gmail.com', 'hash012', NOW());

-- Вставляем товары
INSERT INTO product (name, description, price, amount, category_id, created_at) VALUES
-- Смартфоны
('iPhone 15 Pro', 'Флагманский смартфон Apple с камерой 48 МП', 99990.00, 15, 1, NOW()),
('Samsung Galaxy S24', 'Android смартфон с экраном 6.2"', 79990.00, 20, 1, NOW()),
('Xiaomi Redmi Note 13', 'Бюджетный смартфон с хорошей батареей', 25990.00, 30, 1, NOW()),

-- Ноутбуки
('MacBook Air M2', 'Ультрабук Apple на процессоре M2', 119990.00, 8, 2, NOW()),
('ASUS ROG Strix', 'Игровой ноутбук с RTX 4060', 149990.00, 5, 2, NOW()),
('Lenovo ThinkPad', 'Бизнес-ноутбук для работы', 89990.00, 12, 2, NOW()),

-- Наушники
('AirPods Pro 2', 'Беспроводные наушники с шумоподавлением', 24990.00, 25, 3, NOW()),
('Sony WH-1000XM5', 'Полноразмерные наушники с ANC', 34990.00, 10, 3, NOW()),
('Samsung Galaxy Buds', 'TWS наушники для Android', 12990.00, 18, 3, NOW()),

-- Планшеты
('iPad Air 5', 'Планшет Apple с чипом M1', 59990.00, 7, 4, NOW()),
('Samsung Tab S9', 'Android планшет с S-Pen', 74990.00, 6, 4, NOW()),
('Xiaomi Pad 6', 'Бюджетный планшет для развлечений', 32990.00, 15, 4, NOW()),

-- Умные часы
('Apple Watch Series 9', 'Умные часы для iPhone', 32990.00, 12, 5, NOW()),
('Samsung Galaxy Watch 6', 'Умные часы для Android', 27990.00, 10, 5, NOW()),
('Xiaomi Mi Band 8', 'Фитнес-браслет с отслеживанием сна', 3990.00, 40, 5, NOW());
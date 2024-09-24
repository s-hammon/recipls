-- +goose Up
INSERT INTO categories (id, created_at, updated_at, name)
VALUES 
    (gen_random_uuid(), now(), now(), 'Breakfast'),
    (gen_random_uuid(), now(), now(), 'Snack'),
    (gen_random_uuid(), now(), now(), 'Lunch'),
    (gen_random_uuid(), now(), now(), 'Dinner'),
    (gen_random_uuid(), now(), now(), 'Dessert'),
    (gen_random_uuid(), now(), now(), 'Drink'),
    (gen_random_uuid(), now(), now(), 'Snack'),
    (gen_random_uuid(), now(), now(), 'Other');

-- +goose Down
TRUNCATE TABLE categories;
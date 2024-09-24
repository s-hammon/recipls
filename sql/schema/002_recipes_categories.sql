-- +goose Up
CREATE TABLE categories (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    name VARCHAR(20) NOT NULL
);

CREATE TABLE recipes (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    title VARCHAR(20) NOT NULL,
    description VARCHAR(100) NOT NULL,
    ingredients text NOT NULL,
    instructions text NOT NULL,
    category_id UUID NOT NULL,
    CONSTRAINT fk_category_id
      FOREIGN KEY (category_id)
      REFERENCES categories(id)
      ON DELETE CASCADE,
    user_id UUID NOT NULL,
    CONSTRAINT fk_user_id
      FOREIGN KEY (user_id)
      REFERENCES users(id)
      ON DELETE CASCADE
);

-- +goose Down
DROP TABLE recipes;
DROP TABLE categories;
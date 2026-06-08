-- +goose Up
INSERT INTO accounts (id, balance)
VALUES (1, 1000);

-- +goose Down
DELETE FROM accounts
WHERE id = 1;

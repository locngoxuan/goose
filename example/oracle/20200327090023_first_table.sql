-- +goose Up
-- +goose StatementBegin
CREATE TABLE test (
    id NUMBER(19),
    CREATED TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY(id)
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE test CASCADE CONSTRAINTS
-- +goose StatementEnd

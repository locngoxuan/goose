-- +goose Up
-- +goose StatementBegin
CREATE TABLE test (
(
    ID BIGINT NOT NULL,
    CREATED TIMESTAMP WITH TIME ZONE,
    PRIMARY KEY (ID)
)
WITH (
    OIDS = FALSE
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE test CASCADE CONSTRAINTS
-- +goose StatementEnd

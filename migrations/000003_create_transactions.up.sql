CREATE TABLE transactions
(
    id                 int4 GENERATED BY DEFAULT AS IDENTITY NOT NULL,
    account_id         int4                                  NOT NULL REFERENCES accounts (id) ON DELETE CASCADE, -- Счет транзакции
    amount             numeric                               NOT NULL,                                            -- Сумма транзакции
    transaction_type   int4                                  NOT NULL,                                            -- Тип транзакции
    transaction_status int4                                  NOT NULL,                                            -- Статус транзакции
    created_at         timestamp                             NOT NULL,                                            -- Дата транзакции
    CONSTRAINT transactions_pk PRIMARY KEY (id)
);

CREATE INDEX idx_transactions_account_id ON transactions (account_id);

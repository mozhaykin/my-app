package client

//go:generate oapi-codegen --config=config.yaml --skip-validation ../../api/http/profile_v2.yaml

// Добавил флаг --skip-validation, чтобы не генерировалась валидация. Иначе некоторые ошибки не логируются,
// т.к. сгенерированный валидатор отлавливает их раньше, чем вызывается мой код, в котором я логирую ошибку.

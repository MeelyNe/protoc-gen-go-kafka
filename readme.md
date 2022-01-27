# Установка

```bash
go install github.com/vseinstrumentiru/protoc-gen-go-kafka@latest
```

# Настройка
```yaml
version: v1
plugins:
    # подключение плагина
  - name: go-kafka
    # путь для сохранения
    out: generate/
    # suffix=Out установка суффикса
    opt: paths=source_relative,suffix=Out
```


# Message .proto
> Сгенерируются хендлеры сообщений с указанным суффиксом

```protobuf
message RestsOut {
  string nomenclature = 1;
}
```

# Генерация
```bash
#buf generate --template={файл конфигурации} {каталог с контрактами}
buf generate --template=buf.yaml api
```

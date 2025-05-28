# Тестирования для data сервисов

## load (нагрузочное тестирование)
### Установка
#### mac os
```bash
brew install k6
```

#### linux (Debian/Ubuntu)
```bash
sudo apt update
sudo apt install gnupg ca-certificates
curl -s https://dl.k6.io/key.gpg | sudo apt-key add -
echo "deb https://dl.k6.io/deb stable main" | sudo tee /etc/apt/sources.list.d/k6.list
sudo apt update
sudo apt install k6
```

проверьте установку
```bash
k6 version
```

### Запуск тестирования
1) следует указать hostname и port в файле `load/config.js`
```javascript
export const config = {
  host: 'localhost',
  port: 8080,
};
````

казать переменную окружения `K6_PROMETHEUS_RW_SERVER_URL` с адресом сервера Prometheus Remote Write.


```
```

#### Векторный поиск

Для запуска тестов с использованием Prometheus Remote Write необходимо:

1. Указать переменную окружения `K6_PROMETHEUS_RW_SERVER_URL` с адресом сервера Prometheus Remote Write.

2.  Пример запуска тестов:

    **Проверка конфигурации теста:**
    ```bash
    k6 run --vus 10 --duration 10s --http-debug=full stress/search.js
    ```

    **Запуск с выгрузкой в Prometheus:**
    ```bash
    K6_PROMETHEUS_RW_SERVER_URL=<your_prometheus_endpoint>/api/v1/write \
    k6 run -o experimental-prometheus-rw --http-debug=none stress/search.js
    ```
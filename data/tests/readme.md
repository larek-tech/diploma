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

#### Векторный поиск
1) запустите тестирование в отдельном терминале
```bash
k6 run --vus 10 --duration 10s --http-debug=full --out json=output.json stress/search.js
```
2) стресс тестирование 
```bash
k6 run --out json=stress_test.json stress/search.js
 ```
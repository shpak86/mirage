# Mirage

Универсальный CLI‑клиент для HTTP с подменой браузерного фингерпринта (TLS/JA\* + HTTP‑заголовки).

Mirage позволяет отправлять HTTP(S)‑запросы, которые выглядят как трафик настоящего браузера: он имитирует Chrome/Firefox под разными ОС, управляя и HTTP‑заголовками (User‑Agent, client hints), и низкоуровневым TLS/JA\*‑фингерпринтом.

---

## Возможности

- Один основной вход:
  - `mirage http URL` – выполнить одиночный HTTP(S)‑запрос.
- Имперсонация браузера:
  - `--fp chrome-android`, `--fp firefox-linux` и т.п.
  - Меняются и HTTP‑заголовки, и параметры TLS/JA\*, по которым детектят клиентов.
- Гибкая сборка запроса:
  - Любой HTTP‑метод через `--method` (`GET`, `POST`, `PUT`, ...).
  - Несколько заголовков: `--header "Key:Value"` (флаг можно повторять).
  - Несколько куков: `--cookie "name=value"` (флаг можно повторять).
  - Тело запроса из `stdin` через флаг `--body`.
- Режимы вывода:
  - `meta` – только статус и мета‑информация (подходит для скриптов).
  - `resp` – только тело ответа.
  - `full` – запрос + статус + заголовки + тело (аналог `curl -v`).
- Форматы вывода:
  - `plain` – обычный текст (по умолчанию).
  - `json` – структурированный JSON (удобно для логирования и парсинга).

---

## Установка

Нужен установленный Go (1.21+):

Вы можете собрать бинарник используя команды ниже.

```bash
git clone https://github.com/shpak86/mirage
cd mirage
go build -o mirage .
```

Можете устновить mirage используя команду:

```bash
go install github.com/shpak86/mirage/cmd/mirage
```

---

## Использование

Базовый синтаксис:

```bash
mirage http URL [флаги]
```

### Основные флаги

- `-f, --fp string`  
  Профиль фингерпринта в формате `PLATFORM-OS`.  
  Платформы: `chrome`, `firefox`.  
  ОС: `linux`, `windows`, `mac`, `android`, `macos`.  
  Управляет и HTTP‑заголовками (UA, client hints), и низкоуровневыми TLS/JA\* параметрами, используемыми для имперсонации браузера.  
  По умолчанию: `chrome-android`.

- `-m, --method string`  
  HTTP‑метод (`GET`, `POST`, `PUT`, `DELETE`, `HEAD`, `PATCH`, `OPTIONS`, `CONNECT`).  
  По умолчанию: `GET`.

- `-H, --header KEY:VALUE`  
  Установить HTTP‑заголовок. Флаг можно указывать несколько раз.  
  Пример: `-H "X-Custom-Header:value1" -H "X-Custom-Header:value2"`.

- `-C, --cookie "name=value"`  
  Добавить cookie. Флаг можно повторять.

- `-b, --body`  
  Читать тело запроса из `stdin` (например: `echo 'json' | mirage http ... -b`).

- `-o, --output string`  
  Режим вывода:
  - `meta` – статус + мета‑инфа,
  - `resp` – только тело (по умолчанию),
  - `full` – запрос + статус + заголовки + тело.

- `-F, --format string`  
  Формат вывода:
  - `plain` – текст,
  - `json` – JSON‑структура.
  По умолчанию: `plain`.

---

## Примеры

### Простой GET

```bash
mirage http https://example.com -m GET -f chrome-android
```

### GET с заголовками и куками

```bash
mirage http https://httpbin.org/headers \
  -m GET \
  -f firefox-linux \
  -H "X-Debug:1" \
  -H "Accept:application/json" \
  -C "session=313373" \
  -o full \
  -F plain
```

### POST с телом из stdin

```bash
echo '{"user":"alice","pass":"secret"}' | \
  mirage http https://httpbin.org/post \
    -m POST \
    -f chrome-windows \
    -b \
    -F json \
    -o full
```

В режиме `-F json -o full` Mirage выводит один JSON с двумя секциями: `request` и `response` (метод, URL, заголовки, тело, статус, заголовки, тело).

---

## Как работает подмена фингерпринта (в общих чертах)

Внутри Mirage использует клиент `surf` для имперсонации браузера:

- Разбирает значение `--fp` на части `browser` и `os` (например, `chrome-android`).
- Применяет OS‑пресет: `Android()`, `Windows()`, `Linux()`, `MacOS()`.
- Применяет пресет браузера: `Chrome()` или `Firefox()`, в том числе настройки TLS/JA\* и HTTP/2.
- Собирает и отправляет запрос с твоими заголовками, куками и телом.

На выходе получается трафик, гораздо более похожий на настоящий браузер, чем у дефолтного Go `net/http`.

---

## Дальнейшие планы

- Добавить новые протоколы: WebSocket (`mirage ws`), сырой TCP (`mirage tcp`).
- Расширить пул fingerprint‑профилей и дать диагностику (просмотр JA3/JA4).
- Добавить режим бенчмаркинга, чтобы сравнивать, как разные профили живут на целевом сайте.

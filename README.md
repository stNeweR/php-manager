# php-manager

Интерактивный CLI-генератор PHP/Laravel файлов на Go.

Утилита помогает быстро создавать классы, интерфейсы, трейты, enum'ы, а также типовые Laravel-артefакты (контроллеры, модели, middleware, Form Request'ы и ресурсы). Она автоматически определяет корень проекта по `composer.json`, подбирает namespace на основе PSR-4 автозагрузки и предлагает автодополнение папок.

## Возможности

- Автоматический поиск `composer.json` вверх по дереву каталогов.
- Определение namespace по секции `autoload.psr-4`.
- Автодополнение папок при вводе (`Tab`).
- Поддержка Laravel-шаблонов с корректными `use` и базовыми методами.
- Приятный TUI на базе [Charmbracelet Huh](https://github.com/charmbracelet/huh) и [Lipgloss](https://github.com/charmbracelet/lipgloss).

## Поддерживаемые типы файлов

| Тип           | Описание                                  |
|---------------|-------------------------------------------|
| PHP Class     | Обычный PHP-класс                         |
| PHP Interface | PHP-интерфейс                             |
| PHP Trait     | PHP-трейт                                 |
| PHP Enum      | PHP-перечисление с примером кейса         |
| Controller    | Laravel-контроллер, наследует `Controller`|
| Model         | Laravel Eloquent Model                    |
| Middleware    | Laravel Middleware с методом `handle`     |
| FormRequest   | Laravel FormRequest с `authorize`/`rules` |
| Resource      | Laravel JsonResource с `toArray`          |

## Установка

```bash
# Склонировать репозиторий
git clone <repo-url>
cd php-manager

# Собрать бинарник
go build -o php-manager .

# (Опционально) установить в $GOPATH/bin
go install .
```

## Использование

Запустите утилиту из любой папки внутри PHP-проекта:

```bash
./php-manager
```

Далее следуйте интерактивным подсказкам:

1. Выберите тип создаваемого файла.
2. Укажите папку (доступно автодополнение по `Tab`).
3. Введите имя класса/файла.

Файл будет создан с автоматически определённым namespace.

## Пример

```bash
cd ~/projects/my-laravel-app
~/tools/php-manager/php-manager
```

```
Что создаём? PHP Class
Папка: app/Services
Имя файла / класса: UserService
```

Результат — файл `app/Services/UserService.php`:

```php
<?php

namespace App\Services;

class UserService
{
    //
}
```

## Системные требования

- Go 1.24.4 или новее.
- Проект с `composer.json` и настроенной PSR-4 автозагрузкой.

## Зависимости

- [github.com/charmbracelet/huh](https://github.com/charmbracelet/huh) — интерактивные формы.
- [github.com/charmbracelet/lipgloss](https://github.com/charmbracelet/lipgloss) — стилизация терминального вывода.

## Лицензия

MIT

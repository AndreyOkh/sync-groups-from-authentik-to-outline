# Sync groups from Authentik to Outline

Приложение предназначено для синхронизации групп из Authentik в Outline Wiki

## Переменные окружения
| Имя                 | Значение по умолчанию | Описание                                                                                              |
|---------------------|-----------------------|-------------------------------------------------------------------------------------------------------|
| AUTHENTIK_SCHEME    | `http`                | Схема подключения к серверу Authentik                                                                 |
| AUTHENTIK_HOST      | -                     | URL-адрес сервера Authentik                                                                           |
| AUTHENTIK_TOKEN     | -                     | API-токен пользователя Authentik                                                                      |
| OUTLINE_URL         | -                     | URL-адрес сервера Outline                                                                             |
| OUTLINE_TOKEN       | -                     | API-токен Outline                                                                                     |
| GROUP_PREFIX        | `outline_`            | Префикс групп для фильтрации. Группы отвечающие требованиям будут автоматически создаваться в Outline |
| GROUP_NAME_SELECTOR | `name`                | Переменная указывает какой атрибут использовать для именования групп в Outline.                       |

### AUTHENTIK_TOKEN
Для работы приложения создайте сервисный аккаунт, сгенерируйте для него API-ключ в разделе "Каталог -> Токены и пароли приложений", и назначьте права доступа.
</br>Если вы не планируете запускать веб-сервис и будете пользоваться только функцией force-resync, то вам будет достаточно назначить права на чтение групп "Can view Group".
</br>Если вы планируете запускать приложение с веб-сервисом для обновления групп пользователя при авторизации в системе, то вам необходимо дополнительно предоставить права на чтение пользователей "Can view User".

### GROUP_PREFIX
Для того что бы ограничить количество обрабатываемых групп используется префикс, по-умолчанию приложение ищет все группы название которых начинается с "outline_", но вы можете указать любое значение. 

### GROUP_NAME_SELECTOR
Вы можете указать любой атрибут, что бы его значение использовалось в качестве имени группы в Outline.
Если вы используете каталог AD/LDAP в Authentik для получения групп, то вы можете настроить синхронизацию атрибутов.

Для этого в интерфейсе администратора в Authentik перейдите на вкладку Персонализация -> Сопоставление свойств, и создайте свойство с типом "LDAP Source Property Mapping", укажите произвольное имя, например "Active Directory Description", и вставьте в выражение следующий код:
```python
return {
    "attributes": {
        "description": list_flatten(ldap.get("description")),
    },
}
```
Затем перейдите в Каталог -> Федерации и соц. вход, выберите ваш каталог и откройте его для редактирования, в блоке "Сопоставления свойств группы" необходимо выбрать созданное выше сопоставление.

После синхронизации каталога у ваших групп появится новый атрибут "description", значением которого будет описание полученное из каталога пользователей.

## Флаги
| Флаг             | Описание                                                                                                                                                  |
|------------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------|
| `--force-resync` | Флаг запустит принудительную синхронизацию групп всех пользователей. Отсутствующие в Outline группы будут созданы. После выполнения приложение завершится |

## Запуск с помощью Docker Compose
Опционально для работы TLS добавьте ваш корневой сертификат в контейнер
В данном примере приложение запустит web-сервис и будет ожидать вызова от Outline. Outline требует https для отправки webhooks, поэтому вам скорее всего придется использовать обратный прокси для публикации приложения.
После публикации настройте в Outline webhook на событие users.signin
```yaml
services:
  group-sync:
    image: andreyokh/sync-groups-from-authentik-to-outline:latest
    restart: always
    ports:
      - "8081:8081"
    volumes:
      - "./ssl/ca.crt:/etc/ssl/certs/ca-certificates.crt"
    environment:
      AUTHENTIK_SCHEME: "https"
      AUTHENTIK_HOST: "authentik.example.com"
      AUTHENTIK_TOKEN: "xxxxxxxxxxxxxxxxxxxxxxx"
      OUTLINE_URL: "https://outline.example.com"
      OUTLINE_TOKEN: "xxxxxxxxxxxxxxxxxxxxxxxx"
      GROUP_PREFIX: "wiki_"
      GROUP_NAME_SELECTOR: "description"
```

</br>В примере ниже приложение запустится 1 раз, синхронизирует группы, после чего завершится. Вы можете добавить задание CRON что бы запускать контейнер регулярно
```yaml
services:
  groups-resync:
    image: andreyokh/sync-groups-from-authentik-to-outline:latest
    command:
      - "--force-resync"
    volumes:
      - "./ssl/ca.crt:/etc/ssl/certs/ca-certificates.crt"
    environment:
      AUTHENTIK_SCHEME: "https"
      AUTHENTIK_HOST: "authentik.example.com"
      AUTHENTIK_TOKEN: "xxxxxxxxxxxxxxxxxxxxxxx"
      OUTLINE_URL: "https://outline.example.com"
      OUTLINE_TOKEN: "yyyyyyyyyyyyyyyyyyyyyyyy"
      GROUP_PREFIX: "wiki_"
      GROUP_NAME_SELECTOR: "description"
```


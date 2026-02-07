# Library

### Запуск


Чтобы запустить проект, нужно сделать несколько шагов

**Шаг 1)**
 
Положить config.yaml в директорию config

Пример config.yaml:

```yaml
env: "local"
http_server:
  address: "localhost:8080"
  timeout: 4s
  idle_timeout: 30s
```

**Шаг 2)** 

Установить перменнную среды **CONFIG_PATH**

Пример:\
CONFIG_PATH=C:\Users\crud\config\config.yaml

# go-mux-http

[![codecov](https://codecov.io/gh/alexpts/go-mux-http/branch/master/graph/badge.svg?token=lQGrGVDEUo)](https://codecov.io/gh/alexpts/go-mux-http)
[![Go Report Card](https://goreportcard.com/badge/github.com/alexpts/go-mux-http)](https://goreportcard.com/report/github.com/alexpts/go-mux-http)

Библиотека для запуска net/http Handler. (Порт библиотеки https://alexpts.github.io/go-next-docs/ совместимый с net/http вместо fastHttp)

[Документация](https://alexpts.github.io/go-mux-http/)

## Отличия от http.ServeMux

- Возможность приоритизации обработчиков (параметр `priority: int`).
- Поддержка RegExp в путях (возможность описать множество uri 1 правилом).
- Поддержка захвата переменных из пути маршрута.
- Наложение ограничений на любую переменную из пути маршрута через regexp.
- Возможность фильтрации обработчиков по пути, HTTP методу, regexp, параметрам и т.д.
- Возможность использовать кастомный resolver для фильтрации слоев со своей собственной логикой.


### Примеры вариаций и комбинаций подходов

Легенда:

`h - http.Handler`

`m - http.ServeMux`

`app - go-mux-runner/mux/MicroApp`

--- 

Стандартные middleware с делегированием (обертками). Классическая нативная цепочка из middleware http.Handler:

`h(h) -> (h(h)) -> h(h)`

---
http Mux сервер, регистрирует обработчики по путям (1 путь - 1 обработчик), но цепочки обработчиков нужно делать как в пункте 1:

`mux(h) -> (h(h)) -> h(h)`

---

go-mux-http позволяет регистрировать обработчики с повторяющимися путями (либо сопоставлять по regexp разные пути в 1 правило), можно строить различные цепочки из обработчиков. Можно применять в 1 проекте оба подхода и сосуществовать. Можно гибко делегировать управление и из http.Handler обработчиков в mux_runner, так же как и из mux_runner в http.Handler (все очень гибко).


`app(h) -> app(h) -> app(h)` - все обработчики описаны в формате go-mux-runner, удобно для новых приложений;

`app(h) -> h(h) -> h(h)` - можно использовать из пункта 1 нативные цепочки обработчиков (менее гибко);

`h(h) -> h(h) -> h(app) -> app(h)` - можно из обработчика делегировать управление в mux_runner на любом этапе;

`app(h) -> app(h) -> h(h) -> h(h)` - можно из mux_runner делегировать управление в http.Handle цепочку на любом этапе (и обратно);


### FAQ

Q: Как в handler слоя указывать уже обернутый обработчик?

A: Так же как и не обернутый http.Handler по интерфейсу ничем не отличается и может сам делегировать обернутым обработчикам (нативно).

A: Передавать обработчики в слой в виде множества обработчиков, в каждом обработчике пре делегировании управления следующему обработчику вызывать снова `app.ServerHTTP(w, r)` на приложении. Приложение под капотом будет итерироваться дальше по слоям и обработчикам.

A: Можно запустить как обработчик другой mux


### Example

[Examples](https://github.com/alexpts/go-mux-http/tree/master/cmd)

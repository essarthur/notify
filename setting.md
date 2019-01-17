

## Основные параметры

1. Соединений нужно очень много - повышение umilit -n до уровня в >10к дескрипторов.
2. Вместо MaxIdleConnsPerHost задаем CURLOPT_MAXCONNECTS побольше, иначе опять не взлетим по CPU.
3. http.Client с увеличенным MaxIdleConnsPerHost для поддержки keep-alive.
4. При отправке в GCM установить параметр отправки delay_while_idle в true — он говорит о том, что не нужно доставлять пуш пользоватлю, если устройство не активно.
5. Знаем версию его приложения, версию операционки смартфона, модель устройства и часовой пояс пользователя. Всю интересную для маркетинга статистику мы вытаскиваем и сохраняем на сервере.
6. Чтобы не перегрузить сервис - в течение дня она рассылает каждый час по несколько тысяч сообщений, а не пушит всю базу в один присест.




в Go подтягиваем лимиты сразу до позволенного максимума примерно так:

```golang

var rLimit syscall.Rlimit
if err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit); err != nil {
	return nil, err
}

if rLimit.Cur < rLimit.Max {
	rLimit.Cur = rLimit.Max
	syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit)
}
```

Во всех воркерах есть внутренние буферы переотправки на случай сбоя, которые мы считаем не фатальным (таймаут запроса, или 502 код ответа, к примеру). Выглядит примерно так:

```golang

for {
	select {
	case push := <-mainChan:
		send(push)
	case push := <-resendChan:
		send(push)
	default:
		// ...
	}
}

func send(push Push) {
	if !doSmth(push) {
		resendChan <- push
	}
}
```


Однако, сохранять статистику по каждому пушу из каждой горутины (а ведь их много тысяч) в единое место крайне накладно. Поэтому, все воркеры собирают свои статистики сначала у себя локально, и лишь время от времени (раз в несколько секунд) сливают ее в общее место. Примерный код:

```golang
type Stats struct {
	sync.RWMutex

	ElapsedTime ...
	Methods ...
	AppID ...
	...
}

addStatsTicker := time.Tick(5 * time.Second)
for {
	select {
	case <-addStatsTicker:
		globalStats.Lock()
		gcm.stats.Lock()
		mergeStatsToGlobal(&gcm.stats)
		cleanStats(&gcm.stats)
		gcm.stats.Unlock()
		globalStats.Unlock()

	case push := <-mainChan:
		// таких статистик много, это пример одной из них
		gcm.stats.Lock()
		statsMethodIncr(&gcm.stats, push.Method)
		statsAppIDIncr(&gcm.stats, push.AppID)
		gcm.stats.Unlock()

		send(push)
	// ...
	}
}
```



## Links
[Пушилка](https://habr.com/ru/company/vkontakte/blog/265731/)
[Основы успешной реализации push-уведомлений для мобильных приложений](https://habr.com/ru/company/techmas/blog/262411/)
[Stack](https://docs.parseplatform.org/rest/guide/)
[Test](https://www.pushwoosh.com/demo-success/)
[Web push](https://habr.com/ru/post/321924/)
[Mysql driver](https://github.com/go-sql-driver/mysql)
[Connect postgree](https://www.compose.com/articles/accessing-relational-databases-using-go/)


## GO Lib
[Chanal](https://gobyexample.com/non-blocking-channel-operations)
[Chanal talk](https://talks.golang.org/2012/10things.slide#1)
[Switch](https://github.com/golang/go/wiki/Switch)







## Основные параметры

[Пушилка](https://habr.com/ru/company/vkontakte/blog/265731/)

1. Соединений нужно очень много - повышение umilit -n до уровня в >10к дескрипторов.
2. Вместо MaxIdleConnsPerHost задаем CURLOPT_MAXCONNECTS побольше, иначе опять не взлетим по CPU.
3. http.Client с увеличенным MaxIdleConnsPerHost для поддержки keep-alive.




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

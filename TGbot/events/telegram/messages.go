package telegram


const msgHelp = `Команды
/help - справка
/rnd - получить случайную справку
после /rnd - ссылка удаляется
отправить ссылку - сохранить ссылку`

const msgHello = "Здарова \n\n" + msgHelp

const (
	msgUnknownCommand = "Нихуя не понял"
	msgNoSavedPages = "Не сохранилось"
	msgSaved = "Заебись, ссылка добавлена"
	msgAlreadyExists = "Ссылка уже есть, не мороси"
	msgErrSendRandom = "У меня нихуя не получилось"
	msgNoPages = "У тея нихуя нет!!"
)
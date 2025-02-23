package app

import (
	"github.com/irbis_assistant/internal/api"
	"github.com/irbis_assistant/internal/files"
	"log/slog"
)

type Application struct {
	lgr        *slog.Logger
	addr       string
	clientPath string
}

func New(lgr *slog.Logger, addr, clientPath string) *Application {
	return &Application{lgr: lgr, addr: addr, clientPath: clientPath}
}

func (a *Application) Start() error {
	localVer, err := files.ReadVersion(a.clientPath + "/IRBIS64/" + "version.json")
	if err != nil {
		a.lgr.Error(err.Error())
		return err
	}

	if localVer != nil && *localVer == "" {
		a.lgr.Debug("локальная версия не обнаружена")
		err := files.WriteVersion(a.clientPath+"/IRBIS64/"+"version.json", "v0.0.0")
		if err != nil {
			return err
		}
	}

	serverVer, err := api.GetActualVersion(a.addr)
	if err != nil {
		a.lgr.Error("не удалось подключиться к API", slog.String("msg", err.Error()))
		return err
	}

	if serverVer != nil && *serverVer == "" {
		a.lgr.Error("версия не получена", slog.String("версия сервера", *serverVer), slog.String("версия локальная", *localVer))
		return err
	}

	if ok, err := files.CheckConfig(a.clientPath); err == nil {
		a.lgr.Info("проверка конфигурационного файла...", slog.Bool("конфиг корректен", ok))
	} else {
		a.lgr.Error(err.Error())
	}

	if *serverVer == *localVer {
		a.lgr.Info("уже установлена актуальная версия")
		return nil
	}
	a.lgr.Info("скачивание актуальной версии...")
	err = api.GetClient(a.addr, *serverVer)
	if err != nil {
		a.lgr.Error(err.Error())
		return err
	}
	a.lgr.Info("скачивание завершено, но это еще не конец!")

	err = files.Unpack("clients/"+*serverVer+".zip", a.clientPath+"/IRBIS64")
	if err != nil {
		a.lgr.Error(err.Error())
		return err
	}

	err = files.WriteVersion(a.clientPath+"/IRBIS64/"+"version.json", *serverVer)
	if err != nil {
		a.lgr.Error(err.Error())
	}
	a.lgr.Info("клиент обновлён", slog.String("версия", *serverVer))
	return nil
}

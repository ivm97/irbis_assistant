package files

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
)

type VersionFile struct {
	Version string `json:"version"`
}

func ReadVersion(filename string) (*string, error) {
	b, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var vf VersionFile
	err = json.Unmarshal(b, &vf)
	if err != nil {
		return nil, err
	}
	return &vf.Version, nil
}
func WriteVersion(filename, ver string) error {
	vf := VersionFile{Version: ver}
	b, err := json.Marshal(vf)
	if err != nil {
		return err
	}
	err = os.WriteFile(filename, b, 0644)
	if err != nil {
		return err
	}
	return nil
}

func CheckConfig(path string) (bool, error) {
	configPath := path + "/IRBIS64/Cirbis.ini"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false, nil
	}
	b, err := os.ReadFile(configPath)
	if err != nil {
		return false, err
	}
	if bytes.Contains(b, []byte("[Main]")) {
		return true, nil
	}
	return false, nil
}

func Unpack(filename, irbisPath string) error {
	// Открываем ZIP-архив
	r, err := zip.OpenReader(filename)
	if err != nil {
		return err
	}
	defer r.Close()

	// Проходим по всем файлам в архиве
	for _, f := range r.File {
		// Получаем полное имя файла
		fPath := filepath.Join(irbisPath, f.Name)

		// Проверяем, является ли файл директорией
		if f.FileInfo().IsDir() {
			// Создаем директорию
			if err := os.MkdirAll(fPath, os.ModePerm); err != nil {
				return err
			}
			continue
		}

		// Создаем все необходимые родительские директории
		if err := os.MkdirAll(filepath.Dir(fPath), os.ModePerm); err != nil {
			return err
		}

		// Открываем файл в ZIP-архиве
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		// Создаем файл на диске
		outFile, err := os.Create(fPath)
		if err != nil {
			return err
		}
		defer outFile.Close()

		// Копируем содержимое файла из архива в созданный файл
		if _, err := io.Copy(outFile, rc); err != nil {
			return err
		}
	}

	return nil
}

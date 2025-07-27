package depend

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
)

// El archivo output.zip que esta en la carpeta de descargas se descomprime y se borra el archivo output.zip
func (a *Accessos) descompactaZip(downloadDir string) error {
	zipPath := filepath.Join(downloadDir, "output.zip")
	zipReader, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	// Procesa el ZIP
	for _, file := range zipReader.File {
		if err := a.extraeArchivo(file, downloadDir); err != nil {
			zipReader.Close() // Cierra antes de salir por error
			return err
		}
	}
	zipReader.Close() // Cierra antes de eliminar

	// Eliminar el archivo ZIP original
	if err := os.Remove(zipPath); err != nil {
		return err
	}
	return nil
}

// extraeArchivo extrae un archivo individual del zip en el directorio de descarga
func (a *Accessos) extraeArchivo(file *zip.File, downloadDir string) error {
	rc, err := file.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	path := filepath.Join(downloadDir, file.Name)

	// Crear directorios si es necesario
	if file.FileInfo().IsDir() {
		return os.MkdirAll(path, os.ModePerm)
	}

	if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
		return err
	}

	outFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
	if err != nil {
		return err
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, rc)
	return err
}

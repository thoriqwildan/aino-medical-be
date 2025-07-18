package helper

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

func TranslateErrorMessage(err error) map[string]string {
	// Membuat map untuk menampung pesan error
	errorsMap := make(map[string]string)

	// Handle validasi dari validator.v10
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fieldError := range validationErrors {
			field := fieldError.Field() // Menyimpan nama field yang gagal validasi
			switch fieldError.Tag() {   // Menangani berbagai jenis validasi
			case "required":
				errorsMap[field] = fmt.Sprintf("%s is required", field) // Pesan error jika field kosong
			case "email":
				errorsMap[field] = "Invalid email format" // Pesan error jika format email tidak valid
			case "unique":
				errorsMap[field] = fmt.Sprintf("%s already exists", field) // Pesan error jika data sudah ada
			case "min":
				errorsMap[field] = fmt.Sprintf("%s must be at least %s characters", field, fieldError.Param()) // Pesan error jika nilai terlalu pendek
			case "max":
				errorsMap[field] = fmt.Sprintf("%s must be at most %s characters", field, fieldError.Param()) // Pesan error jika nilai terlalu panjang
			case "numeric":
				errorsMap[field] = fmt.Sprintf("%s must be a number", field) // Pesan error jika nilai bukan angka
			default:
				errorsMap[field] = "Invalid value" // Pesan error default untuk kesalahan validasi lainnya
			}
		}
	}

	// Handle error dari GORM untuk duplicate entry
	if err != nil {
		// Cek jika error mengandung "Duplicate entry" (duplikasi data di database)
		if strings.Contains(err.Error(), "Duplicate entry") {
			if strings.Contains(err.Error(), "username") {
				errorsMap["Username"] = "Username already exists" // Pesan error jika username sudah ada
			}
			if strings.Contains(err.Error(), "email") {
				errorsMap["Email"] = "Email already exists" // Pesan error jika email sudah ada
			}
		} else if err == gorm.ErrRecordNotFound {
			// Jika data yang dicari tidak ditemukan di database
			errorsMap["Error"] = "Record not found"
		}
	}

	// Mengembalikan map yang berisi pesan error
	return errorsMap
}
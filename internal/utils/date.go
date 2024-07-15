package utils

import (
	"fmt"
	"time"
)

func timeAgo(dateStr string) (string, error) {
	// Parse the input date in the format DD-MM-YYYY HH:MM
	layout := "02-01-2006 15:04"
	date, err := time.Parse(layout, dateStr)
	if err != nil {
		return "", err
	}

	// Calculate the difference between the input date and the current time
	now := time.Now()
	duration := now.Sub(date)

	// Calculate the number of minutes, hours, days, weeks, months, and years
	minutes := int(duration.Minutes())
	hours := int(duration.Hours())
	days := hours / 24
	weeks := days / 7
	months := int(now.Month()) - int(date.Month()) + (int(now.Year())-int(date.Year()))*12
	years := now.Year() - date.Year()

	// Generate the appropriate message
	var message string
	if years > 0 {
		if years == 1 {
			message = "Postado há 1 ano atrás"
		} else {
			message = fmt.Sprintf("Postado há %d anos atrás", years)
		}
	} else if months > 0 {
		if months == 1 {
			message = "Postado há 1 mês atrás"
		} else {
			message = fmt.Sprintf("Postado há %d meses atrás", months)
		}
	} else if weeks > 0 {
		if weeks == 1 {
			message = "Postado há 1 semana atrás"
		} else {
			message = fmt.Sprintf("Postado há %d semanas atrás", weeks)
		}
	} else if days > 0 {
		if days == 1 {
			message = "Postado há 1 dia atrás"
		} else {
			message = fmt.Sprintf("Postado há %d dias atrás", days)
		}
	} else if hours > 0 {
		if hours == 1 {
			message = "Postado há 1 hora atrás"
		} else {
			message = fmt.Sprintf("Postado há %d horas atrás", hours)
		}
	} else if minutes > 0 {
		if minutes == 1 {
			message = "Postado há 1 minuto atrás"
		} else {
			message = fmt.Sprintf("Postado há %d minutos atrás", minutes)
		}
	} else {
		message = "Postado agora"
	}

	return message, nil
}

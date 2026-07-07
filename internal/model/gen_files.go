package model

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func generateHash(key string) string {
	hash := sha256.New()
	hash.Write([]byte(key))
	hexString := hex.EncodeToString(hash.Sum(nil))
	return hexString[:16]
}

type genEvent struct {
	EventID   string
	EventHash string
}

// Генератор событий
// Вход: путь и название файла, общее кол-во, источники, вероятность дублей, seed
func GenerateEvents(path string, n int, s int, duble float64, seed int64) {
	file, err := os.Create(path)
	if err != nil {
		fmt.Println("Ошибка создания файла:", err)
		return
	}
	defer file.Close()
	writer := bufio.NewWriter(file)

	r := rand.New(rand.NewSource(seed))
	baseTime := time.Date(2026, 7, 7, 12, 45, 0, 0, time.UTC)

	history := []genEvent{}

	for i := 0; i < n; i++ {
		var evtId, hash string

		if len(history) > 0 && r.Float64() < duble {
			idx := r.Intn(len(history))
			evtId = history[idx].EventID
			hash = history[idx].EventHash
		} else {
			evtId = "evt_"
			evt := strconv.Itoa(i + 1)
			for len(evt) < 6 {
				evt = "0" + evt
			}
			evtId = evtId + evt

			hash = generateHash(evtId + "ev123" + strconv.Itoa(i+1))
			history = append(history, genEvent{EventID: evtId, EventHash: hash})
		}

		source := ""
		if s != 0 {
			source = "collector_"
			num := 1 + r.Intn(s)
			if num < 10 {
				source = source + "0" + strconv.Itoa(num)
			} else {
				source = source + strconv.Itoa(num)
			}
		}

		ts := baseTime.Add(time.Duration(i*5) * time.Second)

		event := Event{
			Seq:       i + 1,
			EventID:   evtId,
			EventHash: hash,
			Source:    source,
			Timestamp: ts.Format(time.RFC3339),
		}

		data, err := json.Marshal(event)
		if err != nil {
			fmt.Println("Ошибка преобразования (marshal):", err)
			continue
		}

		_, err = writer.Write(data)
		if err != nil {
			fmt.Println("Ошибка записи:", err)
		}
		err = writer.WriteByte('\n')
		if err != nil {
			fmt.Println("Ошибка перевода строки:", err)
		}
	}

	if err := writer.Flush(); err != nil {
		fmt.Println("Ошибка flush:", err)
	}
}

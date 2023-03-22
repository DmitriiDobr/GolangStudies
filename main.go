package main

import (
	"FirstTask/interfaceS3"
	"errors"
	"fmt"
	"github.com/google/uuid"
)

func main() {
	bytes := [][]byte{
		[]byte("file_bytes_1"),
		[]byte("file_bytes_2"),
		[]byte("file_bytes_3"),
		[]byte("file_bytes_4"),
		[]byte("file_bytes_5"),
		[]byte("file_bytes_6"),
	}

	// заменить пустой interface на свою реализацию
	var s3store interfaceS3.S3 = &byteStorage{strorage: make(map[string][]byte, 0), sentBytes: make(map[string]bool, 0)}

	single := New(s3store)

	single.Save(bytes...)
	if err := single.Push(); err != nil {
		panic(err)
	}
}

type Single struct {
	s3        interfaceS3.S3
	localUids []string
}

// Создание
func New(s3 interfaceS3.S3) *Single {
	return &Single{
		s3: s3,
	}
}

// Save сохраняет переданные байты
func (s *Single) Save(bytes ...[]byte) {
	for _, bb := range bytes {
		s.localUids = append(
			s.localUids,
			s.s3.SaveLocal(bb),
		)
	}
}

// Отправляет переданные байты
func (s *Single) Push() error {
	for _, uuid := range s.localUids {
		if !s.s3.Exists(uuid) {
			return errors.New("file not exists")
		}

		s.s3.Push(uuid)
		if !s.s3.DeleteLocal(uuid) {
			return errors.New("by uuid file not pushed")
		}
	}

	return nil
}

// Тип реализующий интерфейс
type byteStorage struct {
	strorage  map[string][]byte
	sentBytes map[string]bool
}

func (b *byteStorage) SaveLocal(v []byte) (uuidV string) {
	id := uuid.New().String()
	b.strorage[id] = v
	b.sentBytes[id] = false
	fmt.Println("Байты сохранены локально!")
	return id
}

func (b *byteStorage) Exists(uuidV string) bool {
	for key, _ := range b.strorage {
		if key == uuidV {
			fmt.Println("uuidV существует!")
			return true
		}
	}
	fmt.Println("uuidV не существует!")
	return false
}

func (b *byteStorage) Push(uuidV string) {
	if b.sentBytes[uuidV] == false {
		b.sentBytes[uuidV] = true
		fmt.Printf("Сохраненные байты успешно отправлены")
	} else {
		fmt.Printf("Вы уже отправляли байты, повторная отправка")
	}

}

func (b *byteStorage) DeleteLocal(uuidV string) (ok bool) {
	delete(b.strorage, uuidV)
	delete(b.sentBytes, uuidV)
	fmt.Println("Все успешно удалено")
	return true
}

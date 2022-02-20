package executer

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/Alma-media/bloxroute/pkg/model"
)

type Storage interface {
	Add(key, value string)
	Del(key string) bool
	Get(key string) (string, bool)
	Range(callback func(key, value string) bool)
}

type Executer struct {
	storage Storage
	output  io.Writer
}

func New(storage Storage, output io.Writer) Executer {
	return Executer{
		storage: storage,
		output:  output,
	}
}

func (e Executer) Insert(data []byte) error {
	var payload model.Payload

	if err := json.Unmarshal(data, &payload); err != nil {
		return err
	}

	e.storage.Add(payload.Key, payload.Value)

	_, err := e.output.Write(
		[]byte(
			fmt.Sprintf("[ADD]\n%s: %s\n", payload.Key, payload.Value),
		),
	)

	return err
}

func (e Executer) Delete(data []byte) error {
	var payload model.Payload

	if err := json.Unmarshal(data, &payload); err != nil {
		return nil
	}

	if !e.storage.Del(payload.Key) {
		return fmt.Errorf("key %q was not found", payload.Key)
	}

	_, err := e.output.Write(
		[]byte(
			fmt.Sprintf("[DEL]\n%s\n", payload.Key),
		),
	)

	return err
}

func (e Executer) GetOne(data []byte) error {
	var payload model.Payload

	if err := json.Unmarshal(data, &payload); err != nil {
		return nil
	}

	value, ok := e.storage.Get(payload.Key)
	if !ok {
		return fmt.Errorf("key %q was not found", payload.Key)
	}

	_, err := e.output.Write(
		[]byte(
			fmt.Sprintf("[GET]\n%s: %s\n", payload.Key, value),
		),
	)

	return err
}

func (e Executer) GetAll([]byte) (err error) {
	if _, err = e.output.Write(
		[]byte("[ALL]\n"),
	); err != nil {
		return
	}

	e.storage.Range(func(key, value string) bool {
		_, err = e.output.Write(
			[]byte(
				fmt.Sprintf("%s: %s\n", key, value),
			),
		)

		return err == nil
	})

	return
}

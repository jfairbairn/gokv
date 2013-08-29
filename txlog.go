package gokv

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
	"os"
	"regexp"
)

type valuereader struct {
	value []byte
}

var writeLineRegex *regexp.Regexp

func init() {
	writeLineRegex = regexp.MustCompile("^([A-Z]+) ([a-zA-Z0-9\\.\\-_]+) (.*)\n")
}

func (vr *valuereader) Read(p []byte) (int, error) {
	if len(vr.value) == 0 {
		return 0, nil
	}
	n := copy(p, vr.value)
	if n < len(vr.value) {
		vr.value = vr.value[n:]
	}
	return n, nil
}

type Txlog struct {
	log *os.File
}

func OpenTxlog(path string, store *Store) error {
	f, err := os.Open(path)
	if err != nil && os.IsExist(err) {
		return err
	}

	if err == nil {
		reader := bufio.NewReader(f)
		var line []byte

		vr := new(valuereader)

		for err != io.EOF {
			line, err = reader.ReadBytes('\n')
			results := writeLineRegex.FindAllSubmatch(line, 1)
			if len(results) != 1 {
				continue
			}
			result := results[0]
			if len(result) != 4 {
				log.Printf("unexpected result %v", result)
				continue
			}

			vr.value = result[3]

			var v interface{}
			jd := json.NewDecoder(vr)
			err = jd.Decode(&v)
			if err != nil && err != io.EOF {
				log.Printf("Error %s decoding \"%s\"", err, result[2])
				return err
			}

			k := string(result[2])

			op := string(result[1])

			switch op {
			case "PUT":
				store.data[k] = v.([]interface{})[0]
			case "DEL":
				delete(store.data, k)
			default:
				log.Printf("Invalid op %s", op)
			}

			store.idgen.OnKey(k)
		}

		// reopen write log in append mode
		err = f.Close()
		if err != nil {
			return err
		}
	}

	txlog := new(Txlog)
	txlog.log, err = os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	store.log = txlog
	return nil
}

func (txlog *Txlog) Close() error {
	return txlog.log.Close()
}

func (txlog *Txlog) Write(op string, k string, v ...interface{}) error {
	log := txlog.log
	_, err := log.WriteString(op + " " + k + " ")
	if err != nil {
		return err
	}
	err = json.NewEncoder(log).Encode(v)
	return err
}

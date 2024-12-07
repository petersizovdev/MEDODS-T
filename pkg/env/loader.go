package env

import (
	"bufio"
	"os"
	"strings"
)

func LoadEnv(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		l := scanner.Text()

		if eql := strings.Index(l, "="); eql >= 0 {
			if k := strings.TrimSpace(l[:eql]); len(k) > 0 {
				v := ""
				if len(l) > eql {
					v = strings.TrimSpace(l[eql+1:])
				}
				os.Setenv(k, v)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

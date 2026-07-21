package loadfile

import (
	"bufio"
	"os"
	"strings"
)

func LoadEnvFile(filenames ...string) error {
	if len(filenames) == 0 {
		filenames = []string{".env"}
	}

	for _, filename := range filenames {
		file, err := os.Open(filename)
		if err != nil {
			return err
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()

			if strings.HasPrefix(line, "#") || line == "" {
				continue
			}

			parts := strings.SplitN(line, "=", 2)
			if len(parts) != 2 {
				continue
			}

			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])

			//The value set in the operating system has higher priority.
			if val, ok := os.LookupEnv(key); ok && val != "" {
				value = val
			}

			os.Setenv(key, value)
		}
	}

	return nil
}

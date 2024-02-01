package prints

import (
	"bufio"
	"os"
	"strings"
)

func Ask(question string, defaults ...string) (string, error) {
	var defaultAnswer string
	if len(defaults) > 0 {
		defaultAnswer = defaults[0]
	}

	if _, err := Infof("%s ", question); err != nil {
		return "", err
	}
	if defaultAnswer != "" {
		if _, err := Alertf(" [%s] ", defaultAnswer); err != nil {
			return "", err
		}
	}

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	input = strings.TrimSpace(input)
	if input == "" {
		return defaultAnswer, nil
	}
	return input, nil
}

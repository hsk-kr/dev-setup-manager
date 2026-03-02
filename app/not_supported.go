package app

import (
	"fmt"

	"github.com/hsk-kr/dev-setup-manager/lib/styles"
)

func NotSupported(command string) {
	fmt.Println(styles.ErrorText.Render(fmt.Sprintf("⚠ \"%s\" is not supported yet", command)))
}

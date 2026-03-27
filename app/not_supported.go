package app

import (
	"fmt"

	"github.com/hsk-kr/licokit/lib/styles"
)

func NotSupported(command string) {
	fmt.Println(styles.ErrorText.Render(fmt.Sprintf("⚠ \"%s\" is not supported yet", command)))
}

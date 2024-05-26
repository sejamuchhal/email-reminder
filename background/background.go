package background

import (
	"fmt"

	"github.com/sejamuchhal/email-reminder/common"
)

func Run() {
	conf := common.ConfigureOrDie()
	fmt.Printf("Starting background worker with config: %v", conf)
	// Implement the background worker here
}
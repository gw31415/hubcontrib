package image

import (
	"fmt"
	"github.com/tebeka/selenium"
)

const (
	browserName = "firefox"
	className = "js-yearly-contributions"
)

func get_url(username string) string {
	return fmt.Sprintf("https://github.com/%s", username)
}

func Image(username string) ([]byte, error) {

	caps := selenium.Capabilities{"browserName": browserName}
	wd, err := selenium.NewRemote(caps, get_url(username))
	if err != nil {
		return nil, err
	}
	defer wd.Quit()

	em, err := wd.FindElement(selenium.ByClassName, className)
	if err != nil {
		return nil, err
	}

	return em.Screenshot(false)
}

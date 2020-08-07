package image

import (
	"github.com/tebeka/selenium"
)

const (
	browserName = "firefox"
	className   = "js-yearly-contributions"
)

func get_url(username string) string {
	return "http://github.com/" + username
}

func Image(username string) ([]byte, error) {

	caps := selenium.Capabilities{"browserName": browserName}
	wd, err := selenium.NewRemote(caps, "")
	wd.Get(get_url(username))
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

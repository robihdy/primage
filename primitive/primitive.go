package primitive

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

// Modes store the shape available for primitive transformation
var Modes = map[string]int{
	"Combo":          0,
	"Triangle":       1,
	"Rect":           2,
	"Ellipse":        3,
	"Circle":         4,
	"RotatedRect":    5,
	"Beziers":        6,
	"RotatedEllipse": 7,
	"Polygon":        8,
}

// Transform an image to  primitive image
func Transform(img io.Reader, ext string, n int, m string) (io.Reader, error) {
	i, err := tempFile(ext)
	if err != nil {
		return nil, err
	}
	defer os.Remove(i.Name())
	o, err := tempFile("png")
	if err != nil {
		return nil, err
	}
	defer os.Remove(o.Name())

	_, err = io.Copy(i, img)
	if err != nil {
		return nil, err
	}

	_, err = primitive(i.Name(), o.Name(), n, m)
	if err != nil {
		return nil, err
	}

	b := bytes.NewBuffer(nil)
	_, err = io.Copy(b, o)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func primitive(i, o string, n int, m string) (string, error) {
	args := fmt.Sprintf("-i %s -o %s -n %d -m %d", i, o, n, Modes[m])
	cmd := exec.Command("primitive", strings.Fields(args)...)
	b, err := cmd.CombinedOutput()
	return string(b), err
}

func tempFile(ext string) (*os.File, error) {
	tmp, err := ioutil.TempFile("", "img")
	if err != nil {
		return nil, err
	}
	defer os.Remove(tmp.Name())
	return os.Create(fmt.Sprintf("%s.%s", tmp.Name(), ext))
}

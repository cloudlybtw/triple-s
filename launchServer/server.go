package launchServer

import (
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"
	"triple-s/internal"
	"triple-s/internal/utils"
)

type Core struct {
	verbose   bool
	port      string
	directory string

	mux *http.ServeMux
}

var (
	errNotValidPort  = errors.New("input valid port 0-65535")
	errCantCreateDir = errors.New("Can't create directory with this name")
)

func CoreServer(verbose bool, port int, directory string) (*Core, error) {
	if port <= 0 || port >= 65535 {
		return nil, errNotValidPort
	}
	if err := os.Mkdir(directory, os.ModePerm); err != nil {
		if os.IsNotExist(err) {
			return nil, err
		}
	}
	if _, err := os.Stat(directory + "/" + "buckets.csv"); os.IsNotExist(err) {
		f, err := os.Create(directory + "/" + "buckets.csv")
		if err != nil {
			return nil, err
		}
		_, err = f.Write([]byte("Name,Creation,LastModifiedTime,Status\n"))
		if err != nil {
			return nil, err
		}
		f.Close()
	}
	return &Core{
		verbose:   verbose,
		port:      strconv.Itoa(port),
		directory: directory,
		mux:       http.NewServeMux(),
	}, nil
}

func (c *Core) Launch() error {
	c.Crud()
	c.log("Starting server on %v", c.port)

	return http.ListenAndServe(":"+c.port, c.mux)
}

func (c *Core) Crud() {
	c.mux.Handle("GET /{$}", c.getBucket())
	c.mux.Handle("PUT /{bucket}", c.updateBucket())
	c.mux.Handle("DELETE /{bucket}", c.deleteBucket())

	c.mux.Handle("/", c.ErrorXML(http.StatusMethodNotAllowed, "method not allowed"))

	c.mux.Handle("GET /{bucket}/{object}", c.getObject())
	c.mux.Handle("PUT /{bucket}/{object}", c.updateObject())
	c.mux.Handle("DELETE /{bucket}/{object}", c.deleteObject())
}

func (c *Core) ErrorXML(code int, text string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, err := internal.MakeErrorXml(code, text)
		if err != nil {
			return
		}
		utils.ResponceXML(w, r, http.StatusMethodNotAllowed, res)

		c.log("METHOD NOT ALLOWED %s", r.URL.String())
	}
}

func (c *Core) log(format string, v ...any) {
	if c.verbose {
		log.Printf(format, v...)
	}
}

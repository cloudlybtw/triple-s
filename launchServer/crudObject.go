package launchServer

import (
	"io"
	"net/http"
	"os"
	"strings"
	"triple-s/internal"
	"triple-s/internal/csv"
	"triple-s/internal/utils"
)

func (c *Core) getObject() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bucketName := r.PathValue("bucket")
		objectName := r.PathValue("object")

		if !utils.CheckName(bucketName) || !utils.CheckName(objectName) {
			data, err := internal.MakeErrorXml(http.StatusBadRequest, "invalid name")
			if err != nil {
				c.log("error: %v", err)
				return
			}
			utils.ResponceXML(w, r, http.StatusBadRequest, data)
			return
		}

		found, err := csv.IsExistBucket(c.directory, bucketName)
		if err != nil {
			c.log("error: %v", err)
			return
		}
		if !found {
			data, err := internal.MakeErrorXml(http.StatusNotFound, "Can't find object")
			if err != nil {
				c.log("error: %v", err)
				return
			}
			utils.ResponceXML(w, r, http.StatusNotFound, data)
			return
		}
		file, err := os.Open(c.directory + "/" + bucketName + "/" + objectName)
		if err != nil {
			c.log("error: %v", err)
			data, err := internal.MakeResponceXml(http.StatusInternalServerError, "Internal server Error")
			if err != nil {
				c.log("error: %v", err)
				return
			}
			utils.ResponceXML(w, r, http.StatusInternalServerError, data)
			c.log("error: %v", err)
			return
		}
		fullData := make([]byte, 0)
		tempData := make([]byte, 512)

		for {
			n, err := file.Read(tempData)
			if err != nil {
				break
			}
			fullData = append(fullData, tempData[:n]...)
		}

		content_type, err := csv.GetContentType(c.directory, bucketName, objectName)
		if err != nil {
			c.log("error: %v", err)
			data, err := internal.MakeErrorXml(http.StatusInternalServerError, "Internal server Error")
			if err != nil {
				c.log("error: %v", err)
				return
			}
			utils.ResponceXML(w, r, http.StatusInternalServerError, data)
			return
		}
		utils.Responce(w, r, content_type, http.StatusOK, fullData)
		c.log("GET object %s", bucketName+"/"+objectName)
	}
}

func (c *Core) updateObject() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: add regex and check for existing in metadata and rechange time
		bucketName := r.PathValue("bucket")
		objectName := r.PathValue("object")

		if !utils.CheckName(bucketName) || !utils.CheckName(objectName) {
			data, err := internal.MakeErrorXml(http.StatusBadRequest, "invalid name")
			if err != nil {
				c.log("error: %v", err)
				return
			}
			utils.ResponceXML(w, r, http.StatusBadRequest, data)
			return
		}

		found, err := csv.IsExistBucket(c.directory, bucketName)
		if !found || err != nil {
			c.log("error: %v", err)
			data, err := internal.MakeErrorXml(http.StatusNotFound, "cannot find such a bucket")
			if err != nil {
				c.log("error: %v", err)
				return
			}
			utils.ResponceXML(w, r, http.StatusNotFound, data)
			return
		}

		found, err = csv.IsExistObject(c.directory, bucketName, objectName)
		_, err = os.Stat(strings.Join([]string{c.directory, bucketName, objectName}, "/"))

		if !found && os.IsExist(err) {
			data, err := internal.MakeErrorXml(http.StatusBadRequest, "you are trying to rewrite existing file which not include to storage system, permission denied")
			if err != nil {
				c.log("error: %v", err)
				return
			}
			utils.ResponceXML(w, r, http.StatusBadRequest, data)
			return
		}

		bucketFile, err1 := os.Create(c.directory + "/" + bucketName + "/" + objectName)
		if err1 != nil {
			c.log("error lol: %v", err1)
			data, err := internal.MakeErrorXml(http.StatusInternalServerError, "internal server error")
			if err != nil {
				c.log("error: %v", err)
				return
			}
			utils.ResponceXML(w, r, http.StatusInternalServerError, data)
			return
		}
		defer bucketFile.Close()

		buff, err := io.ReadAll(r.Body)
		if err != nil {
			c.log("error: %v", err)
			data, err := internal.MakeErrorXml(http.StatusInternalServerError, "internal server error")
			if err != nil {
				c.log("error: %v", err)
				return
			}
			utils.ResponceXML(w, r, http.StatusInternalServerError, data)
			return
		}
		bucketFile.Write(buff)
		defer bucketFile.Close()

		err1 = csv.UpdateObject(c.directory, bucketName, objectName, r.Header.Get("Content-type"), r.ContentLength)
		if err1 != nil {
			c.log("error with updating object in metadata: %v", err1)
			return
		}
		err1 = csv.UpdateBucket(c.directory, bucketName, false)
		if err1 != nil {
			c.log("error with updating bucket metadata: %v", err1)
			return
		}
		data, err1 := internal.MakeResponceXml(http.StatusOK, "object: "+objectName+" was succesfully created. In bucket: "+bucketName)
		if err1 != nil {
			c.log("error: %v", err1)
			return
		}

		utils.ResponceXML(w, r, http.StatusOK, data)
		c.log("UPT object %s", bucketName+"/"+objectName)
	}
}

func (c *Core) deleteObject() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bucketName := r.PathValue("bucket")
		objectName := r.PathValue("object")

		if !utils.CheckName(bucketName) || !utils.CheckName(objectName) {
			data, err := internal.MakeErrorXml(http.StatusBadRequest, "invalid name")
			if err != nil {
				c.log("error: %v", err)
				return
			}
			utils.ResponceXML(w, r, http.StatusBadRequest, data)
			return
		}

		found, err := csv.IsExistBucket(c.directory, bucketName)
		if !found || err != nil {
			c.log("error: %v", err)
			data, err := internal.MakeErrorXml(http.StatusNotFound, "cannot find such a bucket")
			if err != nil {
				c.log("error: %v", err)
				return
			}
			utils.ResponceXML(w, r, http.StatusNotFound, data)
			return
		}

		found, err = csv.IsExistObject(c.directory, bucketName, objectName)
		_, err = os.Stat(strings.Join([]string{c.directory, bucketName, objectName}, "/"))
		if !found && os.IsExist(err) {
			data, err := internal.MakeErrorXml(http.StatusBadRequest, "you are trying to rewrite existing file which not include to storage system, permission denied")
			if err != nil {
				c.log("error: %v", err)
				return
			}
			utils.ResponceXML(w, r, http.StatusBadRequest, data)
			return
		}
		canDeleteBucket, err := csv.DeleteObject(c.directory, bucketName, objectName)
		if err != nil {
			if err == csv.ErrCantFind || err == csv.ErrCantFindBucket || err == io.EOF {
				data, err := internal.MakeErrorXml(http.StatusNotFound, err.Error())
				if err != nil {
					c.log("error: %v", err)
				}
				utils.ResponceXML(w, r, http.StatusNotFound, data)
				return
			}
			c.log("error with deleting object from metadata: %v", err)
			return
		}
		if canDeleteBucket {
			err := csv.UpdateBucket(c.directory, bucketName, true)
			if err != nil {
				c.log("error with updating bucket to delete it: %v", err)
			}

		}
		err = os.RemoveAll(strings.Join([]string{c.directory, bucketName, objectName}, "/"))
		if err != nil {
			c.log("error: %v", err)
			data, err := internal.MakeErrorXml(http.StatusInternalServerError, "internal server error")
			if err != nil {
				c.log("error: %v", err)
				return
			}
			utils.ResponceXML(w, r, http.StatusInternalServerError, data)
			c.log("error: %v", err)
			return
		}
		utils.ResponceXML(w, r, http.StatusNoContent, nil)
		c.log("DELETE object %s", bucketName+"/"+objectName)
	}
}

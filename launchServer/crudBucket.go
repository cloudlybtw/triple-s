package launchServer

import (
	"net/http"
	"os"
	"triple-s/internal"
	"triple-s/internal/csv"
	"triple-s/internal/utils"
)

func (c *Core) getBucket() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := internal.MakeListAllMyBucketsResult(c.directory + "/buckets.csv")
		if err != nil {
			c.log("error: %v", err)
			return
		}
		utils.ResponceXML(w, r, http.StatusOK, data)
		c.log("GET buckets")
	}
}

func (c *Core) updateBucket() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bucketName := r.PathValue("bucket")

		if !utils.CheckName(bucketName) {
			data, err := internal.MakeErrorXml(http.StatusBadRequest, "invalid name of bucket")
			if err != nil {
				c.log("error: %v", err)
				return
			}
			utils.ResponceXML(w, r, http.StatusBadRequest, data)
			return
		}

		err := os.Mkdir(c.directory+"/"+bucketName, os.ModePerm)
		if os.IsExist(err) {
			data, err := internal.MakeErrorXml(http.StatusBadRequest, "bucket alredy exist or u are trying rechange existing directory")
			if err != nil {
				c.log("error: %v", err)
				return
			}
			utils.ResponceXML(w, r, http.StatusBadRequest, data)
			return
		}

		if _, err := os.Stat(c.directory + "/" + bucketName + "/" + "objects.csv"); os.IsNotExist(err) {
			f, err := os.Create(c.directory + "/" + bucketName + "/" + "objects.csv")
			if err != nil {
				c.log("error: %v", err)
				return
			}
			defer f.Close()

			_, err = f.Write([]byte("ObjectKey,Size,ContentType,LastModified\n"))
			if err != nil {
				c.log("error: %v", err)
				return
			}

		}
		if err := csv.UpdateBucket(c.directory, bucketName, true); err != nil {
			c.log("error: %v", err)
			return
		}
		data, err := internal.MakeResponceXml(http.StatusCreated, bucketName+" was successfully created")
		if err != nil {
			c.log("error: %v", err)
			return
		}
		utils.ResponceXML(w, r, http.StatusCreated, data)
		c.log("UPDATE bucket, name: %v ", bucketName)
	}
}

func (c *Core) deleteBucket() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bucketName := r.PathValue("bucket")
		if !utils.CheckName(bucketName) {
			data, err := internal.MakeErrorXml(http.StatusBadRequest, "can't find such a bucket")
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
			c.log("error: %v", err)
			data, err := internal.MakeErrorXml(http.StatusNotFound, bucketName)
			if err != nil {
				c.log("error: %v", err)
			}
			utils.ResponceXML(w, r, http.StatusNotFound, data)
		}
		err = csv.DeleteBucket(c.directory+"/buckets.csv", bucketName)
		if err != nil {
			c.log("error: %v", err)
			if err == csv.ErrCantDelete || csv.ErrCantFind == err {
				data, err := internal.MakeErrorXml(http.StatusBadRequest, err.Error())
				if err != nil {
					c.log("error: %v", err)
				}
				utils.ResponceXML(w, r, http.StatusBadRequest, data)
			}
			return
		}
		err = os.RemoveAll(c.directory + "/" + bucketName)
		if err != nil {
			c.log("error: %v", err)
			data, err := internal.MakeErrorXml(http.StatusInternalServerError, "Internal server error")
			if err != nil {
				c.log("error: %v", err)
				return
			}
			utils.ResponceXML(w, r, http.StatusInternalServerError, data)
			c.log("error: %v", err)
			return
		}
		utils.ResponceXML(w, r, http.StatusNoContent, nil)
		c.log("DELETE bucket, name: %s", bucketName)
	}
}

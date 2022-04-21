package ctrlutil

import (
	"encoding/json"
	"fmt"
	"github.com/ShatteredRealms/Accounts/pkg/model"
	"github.com/gin-gonic/gin"
	"io/ioutil"
)

// ReadBody Parse the context body for bytes. If there is no payload or there were errors processing it
// then an error is returned and no bytes. Otherwise, the bytes of the body are returned for processing.
func ReadBody(c *gin.Context) ([]byte, error) {
	reqBody := c.Request.Body
	if reqBody == nil {
		err := fmt.Errorf("payload missing")
		resp := model.NewBadRequestResponse(c, err.Error())
		c.JSON(resp.StatusCode, resp)
		return nil, err
	}

	body, err := ioutil.ReadAll(reqBody)
	if err != nil {
		err := fmt.Errorf("unable to process payload")
		resp := model.NewInternalServerResponse(c, err.Error())
		c.JSON(resp.StatusCode, resp)
		return nil, err
	}

	return body, nil
}

func ParseBody(c *gin.Context, output interface{}) error {
	body, err := ReadBody(c)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, output)
	if err != nil {
		err := fmt.Errorf("expected JSON body")
		resp := model.NewBadRequestResponse(c, err.Error())
		c.JSON(resp.StatusCode, resp)
		return err
	}

	return nil
}

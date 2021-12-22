package streams

import (
	//"encoding/json"
	"net/http"
	"time"
	//"fmt"

	"github.com/asumsi/api.inlive/internal/models/stream"
	"github.com/asumsi/api.inlive/pkg"
	"github.com/asumsi/api.inlive/pkg/api"
	//"github.com/asumsi/api.inlive/pkg/ffmpeg"
	"github.com/gorilla/mux"
	// "gopkg.in/go-playground/validator.v10"
)

// endStream godoc
// @Summary      End stream
// @Description  End stream stop process of send chunk video to dash server using FFMPEG
// @Tags         stream
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Stream ID"
// @Success      200  {object}  stream.ResponseSwagEndStreamSuccess
// @Failure		 400  {object}	stream.ResponseSwagEndStreamFail
// @Router       /v1/streams/{id}/end [post]
func (controller *Controller) End(w http.ResponseWriter, r *http.Request) {
	var err error
	params := mux.Vars(r)
	slugOrId := params["id"]
	var ok bool

	// check the existence of slug or id first
	result, err := stream.GetBySlugOrId(slugOrId)
	if err != nil {
		api.RespondJSON(w, api.Response{Code: http.StatusNotFound, Message: "Can't get the stream data", Data: err})
		return
	}

	// check existence of stream session
	if _, ok = controller.Sessions[slugOrId]; !ok {
		api.RespondJSON(w, api.Response{Code: http.StatusBadRequest, Message: "Stream never initiated", Data: nil})
		return
	}

	// stream must on running state, so available to be ended
	if _, ok = controller.FFmpegs[slugOrId]; !ok {
		api.RespondJSON(w, api.Response{Code: http.StatusBadRequest, Message: "Stream had been stopped or never exist", Data: nil})
		return
	}

	/*
	// currently it seems we dont need to read body request on /start
	//endStreamRequest := stream.StartStreamRequest{}
	// currently it seems we dont need to read body request on /start
	decoder := json.NewDecoder(r.Body)

	//decode body to object startStreamRequest
	if err = decoder.Decode(&endStreamRequest); err != nil {
		api.RespondJSON(w, api.Response{Code: http.StatusBadRequest, Message: "Decode request body fail on end Stream", Data:  err.Error()})
		return
	}

	// validate request body of /end
	err = pkg.ValidateRequest(endStreamRequest)
	if err != nil {
		api.RespondJSON(w, api.Response{Code: http.StatusBadRequest, Message: "Validation Error on end Stream", Data:  err.Error()})
		return
	}
	*/

	session := controller.Sessions[slugOrId]
	//err = session.SCTP().Stop()
	err = session.Close()
	session = nil
	if err != nil {
		api.RespondJSON(w, api.Response{Code: http.StatusInternalServerError, Message: "Failed to save stream end_date", Data: ""})
		return
	}

	delete(controller.FFmpegs, slugOrId)
	delete(controller.Sessions, slugOrId)

	result.EndDate = pkg.TimePtr(time.Now())

	_, err = result.Update()
	if err != nil {
		api.RespondJSON(w, api.Response{Code: http.StatusInternalServerError, Message: "Failed to save stream end_date", Data: ""})
		return
	}

	api.RespondJSON(w, api.Response{Code: http.StatusOK, Message: "Streaming Stop", Data: ""})
}
package streams

import (
	"net/http"
	//"encoding/json"
	//"fmt"
	"time"

	"github.com/asumsi/api.inlive/internal/models/stream"
	"github.com/asumsi/api.inlive/pkg"
	"github.com/asumsi/api.inlive/pkg/api"
	"github.com/asumsi/api.inlive/pkg/ffmpeg"
	"github.com/gorilla/mux"

	// "gopkg.in/go-playground/validator.v10"
)

type FFMPeg struct {
	pid int
	url string
}

// startStream godoc
// @Summary      Start stream
// @Description  Start stream send chunk video to dash server using FFMPEG
// @Tags         stream
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Stream ID"
// @Success      200  {object}  stream.ResponseSwagStartStreamSuccess
// @Failure		 200  {object}	stream.ResponseSwagStartStreamFail
// @Router       /v1/streams/{id}/start [post]
func (controller *Controller) Start(w http.ResponseWriter, r *http.Request) {
	var err error

	params := mux.Vars(r)
	slugOrId := params["id"]

	// get data stream by slug or id, to later be updated
	streamObj, err := stream.GetBySlugOrId(slugOrId)
	if err != nil {
		api.RespondJSON(w, api.Response{Code: http.StatusNotFound, Message: "Stream not found", Data: ""})
		return
	}

	// check existence of stream session
	if _, ok := controller.Sessions[slugOrId]; !ok {
		api.RespondJSON(w, api.Response{Code: http.StatusBadRequest, Message: "Stream never initiated", Data: nil})
		return
	}

	// check if FFMPEG related to slug or id already exist/running
	if _, ok := controller.FFmpegs[slugOrId]; ok {
		api.RespondJSON(w, api.Response{Code: http.StatusOK, Message: "Streaming already started, cant be interfered" , Data: ""})
		return
	}

	/*
	// currently it seems we dont need to read body request on /start
	//startStreamRequest := stream.StartStreamRequest{}
	//decoder := json.NewDecoder(r.Body)
	
	// currently it seems we dont need to read body request on /start
	//decode body to object startStreamRequest
	if err = decoder.Decode(&startStreamRequest); err != nil {
		api.RespondJSON(w, api.Response{Code: http.StatusBadRequest, Message: "Decode request body fail on start Stream", Data:  err.Error()})
		return
	}
	
	// validate request body of /start
	err = pkg.ValidateRequest(startStreamRequest)
	if err != nil {
		api.RespondJSON(w, api.Response{Code: http.StatusBadRequest, Message: "Validation Error on start Stream", Data:  err.Error()})
		return
	}
	*/

	// start the streaming
	ffmpegObject, err := startStreaming(slugOrId)

	// if exist any error on run ffmpeg 
	if err != nil {
		api.RespondJSON(w, api.Response{Code: http.StatusInternalServerError, Message: "Stream  encoding failed to run", Data: err})
		return
	}

	// prepare update data on streams
	streamObj.ManifestPath = ffmpegObject.url
	streamObj.StartDate = pkg.TimePtr(time.Now())
	
	// do update table streams in database
	_,err = streamObj.Update()
	if err != nil {
		api.RespondJSON(w, api.Response{Code: http.StatusInternalServerError, Message: "Failed to save stream manifest", Data: ""})
		return
	}
	
	controller.FFmpegs[slugOrId] = ffmpegObject.pid

	// since ffmpegObject on Data, always interpret as empty object by the Data, I save ffmpegObject on map[string]
	dataStreamStart := map[string]string{
		"url" : ffmpegObject.url,
	}

	api.RespondJSON(w, api.Response{Code: http.StatusOK, Message: "Stream run successfully", Data: dataStreamStart})
}

// StartStreaming trigger FFMPEG to start running, consuming audiport and videoport produced by rtpforwarder
// also send stream as chunk to dash server
// as output, it will return pid of FFMPEG and manifest url (that can be played on bifrost video player)
func startStreaming(slug string) (FFMPeg, error) {

	// channel to receive pid answer from FFMPEG instance
	pid := make(chan int, 1)

	// get path of SDP file to be consumed by FFMPEG instance
	// also construct url stream which can be used as parameter in FFMPEG (-f) to send chunked stream to dash server

	fileSDP := ffmpeg.GetFileSDPPath(slug)
	urlStream := ffmpeg.GetStreamPath(slug)

	// start ffmpeg instance as a goroutine, which give PID as channel result
	go ffmpeg.Execute(fileSDP, urlStream, pid)
	valuePid := <-pid

	//result := map[string]interface{}{
	//	"pid" : valuePid,
	//	"url_stream" : urlStream,
	//}
	
	result := FFMPeg{
		pid : valuePid,
		url : urlStream,
	}

	//return pid and url_stream of newly created ffmpeg
	return result, nil
}

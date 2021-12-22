package streams

import (
	//"fmt"
	"net/http"
	"strconv"

	"github.com/asumsi/api.inlive/internal/models/stream"
	"github.com/asumsi/api.inlive/pkg/api"
)

// listStream godoc
// @Summary     Get list of stream
// @Description Update a stream by their ID
// @Tags        stream
// @Accept      json
// @Produce     json
// @Param       live   query  bool false  "Whether to fetch streams that has started but not ended"
// @Success     200  {object}  api.Response{data=[]stream.Stream}
// @Failure	400  {object}  api.Response{data=string}
// @Router      /v1/streams [get]
func List(w http.ResponseWriter, r *http.Request){
	queryParams := r.URL.Query()
	live_string := queryParams.Get("live")
	live := true

	if live_string != ""{
		var err error
		live,err = strconv.ParseBool(live_string)
		if err != nil {
			api.RespondJSON(w, api.Response{Code: http.StatusUnprocessableEntity, Message: "invalid value", Data: ""})
			return
		}
	}

	res, err := stream.GetAll(stream.StreamParams{Live: live})

	if err != nil {
		api.RespondJSON(w, api.Response{Code: http.StatusBadRequest, Message: "Error happened", Data: ""})
		return
	}
	
	api.RespondJSON(w, api.Response{Code: http.StatusOK, Message: "List of streams", Data: res})
	return

}
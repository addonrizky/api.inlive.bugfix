package streams

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/asumsi/api.inlive/internal/models/stream"
	"github.com/asumsi/api.inlive/pkg/api"
)

// updateStream godoc
// @Summary     Update Stream
// @Description Update a stream by their ID
// @Tags        stream
// @Accept      json
// @Produce     json
// @Param        id   path      int  true  "Stream ID to get"
// @Param	body body stream.CreateStreamRequest true "Body Request"
// @Success     200  {object}  api.Response{data=stream.Stream}
// @Failure	400  {object}  api.Response{data=string}
// @Router      /v1/streams/{id} [put]
func Update(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	if id == "" {
		api.RespondJSON(w, api.Response{Code: http.StatusBadRequest, Message: "No slug or ID in update request", Data: ""})
	} else {
		result, err := stream.GetBySlugOrId(id)
		if err != nil {
			api.RespondJSON(w, api.Response{Code: http.StatusBadRequest, Message: err.Error(), Data: err})
		} else {
			decoder := json.NewDecoder(r.Body)
			if err := decoder.Decode(&result); err == nil {
				result.Update()
				if err == nil {
					api.RespondJSON(w, api.Response{Code: http.StatusOK, Message: http.StatusText((http.StatusOK)), Data: result})
				} else {
					api.RespondJSON(w, api.Response{Code: http.StatusBadRequest, Message: err.Error(), Data: err})
				}
			} else {
				api.RespondJSON(w, api.Response{Code: http.StatusBadRequest, Message: err.Error(), Data: err})
			}

		}

	}

}

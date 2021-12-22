package streams

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/asumsi/api.inlive/internal/models/stream"
	"github.com/asumsi/api.inlive/pkg/api"
)

// deleteStream godoc
// @Summary     Delete Stream
// @Description Delete a stream from the database
// @Tags        stream
// @Accept      json
// @Produce     json
// @Param        id   path      int  true  "Stream ID to delete"
// @Success     200  {object}  api.Response{data=stream.Stream}
// @Failure	400  {object}  api.Response{data=string}
// @Router      /v1/streams/{id} [delete]
func Delete(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	if id == "" {
		api.RespondJSON(w, api.Response{Code: http.StatusBadRequest, Message: "No slug or ID in delete request", Data: ""})
	} else {
		result, err := stream.GetBySlugOrId(id)
		delete, err := result.Delete()
		if err == nil {
			api.RespondJSON(w, api.Response{Code: http.StatusOK, Message: http.StatusText((http.StatusOK)), Data: delete})
		} else {
			api.RespondJSON(w, api.Response{Code: http.StatusBadRequest, Message: err.Error(), Data: err})
		}
	}

}

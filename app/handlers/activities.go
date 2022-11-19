package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/ungame/timetrack/app/service"
	"github.com/ungame/timetrack/httpext"
	"github.com/ungame/timetrack/queries"
	"github.com/ungame/timetrack/types"
	"io"
	"net/http"
	"strconv"
)

type activitiesHandler struct {
	activitiesService service.ActivitiesService
}

func NewActivitiesHandler(activitiesService service.ActivitiesService) Handler {
	return &activitiesHandler{activitiesService: activitiesService}
}

func (a *activitiesHandler) Register(router *mux.Router) {
	router.Path("/activities").HandlerFunc(a.PostActivity).Methods(http.MethodPost)
	router.Path("/activities").HandlerFunc(a.GetActivities).Methods(http.MethodGet)
	router.Path("/activities/{id}").HandlerFunc(a.GetActivity).Methods(http.MethodGet)
	router.Path("/activities/{id}").HandlerFunc(a.PutActivity).Methods(http.MethodPut)
	router.Path("/activities/{id}").HandlerFunc(a.DeleteActivity).Methods(http.MethodDelete)
	router.Path("/activities/{id}/finish").HandlerFunc(a.PutFinishActivity).Methods(http.MethodPut)
	router.Path("/activities/_/filter").HandlerFunc(a.FilterActivities).Methods(http.MethodGet)
}

func (a *activitiesHandler) PostActivity(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		httpext.WriteError(w, http.StatusUnprocessableEntity, err)
		return
	}
	input := new(types.Activity)
	err = json.Unmarshal(body, input)
	if err != nil {
		httpext.WriteError(w, http.StatusUnprocessableEntity, err)
		return
	}
	activity, err := a.activitiesService.StartActivity(r.Context(), input)
	if err != nil {
		httpext.WriteError(w, http.StatusUnprocessableEntity, err)
		return
	}
	w.Header().Set(httpext.HeaderLocation, fmt.Sprintf("%s/%d", r.RequestURI, activity.ID))
	httpext.WriteJson(w, http.StatusCreated, activity)
}

func (a *activitiesHandler) GetActivity(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		httpext.WriteError(w, http.StatusBadRequest, err)
		return
	}
	activity, err := a.activitiesService.GetActivity(r.Context(), id)
	if err != nil {
		httpext.WriteError(w, http.StatusBadRequest, err)
		return
	}
	httpext.WriteJson(w, http.StatusOK, activity)
}

func (a *activitiesHandler) GetActivities(w http.ResponseWriter, r *http.Request) {
	activities, err := a.activitiesService.GetActivities(r.Context())
	if err != nil {
		httpext.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	httpext.WriteJson(w, http.StatusOK, activities)
}

func (a *activitiesHandler) FilterActivities(w http.ResponseWriter, r *http.Request) {
	var (
		query  = r.URL.Query()
		period = queries.Today.String()
		order  = queries.Desc.String()
		limit  = 1000
		err    error
	)

	if query.Get("period") != "" {
		period = query.Get("period")
	}

	if query.Get("order") != "" {
		order = query.Get("order")
	}

	if query.Get("limit") != "" {
		limit, err = strconv.Atoi(query.Get("limit"))
		if err != nil {
			httpext.WriteError(w, http.StatusBadRequest, err)
			return
		}
	}

	filter := &types.PeriodFilter{
		PeriodName: period,
		OrderBy:    order,
		Limit:      limit,
	}

	activities, err := a.activitiesService.FilterActivitiesByPeriod(r.Context(), filter)
	if err != nil {
		httpext.WriteError(w, http.StatusBadRequest, err)
		return
	}
	httpext.WriteJson(w, http.StatusOK, activities)
}

func (a *activitiesHandler) PutActivity(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		httpext.WriteError(w, http.StatusBadRequest, err)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		httpext.WriteError(w, http.StatusUnprocessableEntity, err)
		return
	}
	input := new(types.Activity)
	err = json.Unmarshal(body, input)
	if err != nil {
		httpext.WriteError(w, http.StatusUnprocessableEntity, err)
		return
	}
	input.ID = id
	activity, err := a.activitiesService.UpdateActivity(r.Context(), input)
	if err != nil {
		httpext.WriteError(w, http.StatusUnprocessableEntity, err)
		return
	}
	httpext.WriteJson(w, http.StatusOK, activity)
}

func (a *activitiesHandler) PutFinishActivity(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		httpext.WriteError(w, http.StatusBadRequest, err)
		return
	}
	activity, err := a.activitiesService.FinishActivity(r.Context(), id)
	if err != nil {
		httpext.WriteError(w, http.StatusUnprocessableEntity, err)
		return
	}
	httpext.WriteJson(w, http.StatusOK, activity)
}

func (a *activitiesHandler) DeleteActivity(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		httpext.WriteError(w, http.StatusBadRequest, err)
		return
	}
	err = a.activitiesService.DeleteActivity(r.Context(), id)
	if err != nil {
		httpext.WriteError(w, http.StatusUnprocessableEntity, err)
		return
	}
	w.Header().Set(httpext.HeaderEntity, fmt.Sprint(id))
	w.WriteHeader(http.StatusNoContent)
}

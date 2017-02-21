package httplights

import (
	"net/http"

	"github.com/collinux/gohue"
	"github.com/gorilla/mux"
)

// SelectScene is a module that provides a http interface for setting scenes.
type SelectScene struct {
	hueBridge *hue.Bridge
}

// SetHueBridge implmenets the HTTPLightsModule interface.
func (s *SelectScene) SetHueBridge(bridge *hue.Bridge) {
	s.hueBridge = bridge
}

// RegisterInRouter implmenets the HTTPLightsModule interface.
func (s *SelectScene) RegisterInRouter(router *mux.Router) {
	router.Handle("/select-scene", s).Methods("PUT")
}

// ServeHTTP implmenets the HTTPLightsModule and http.Handler interface.
func (s *SelectScene) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	sceneName := r.PostFormValue("scene")

	// Special case for the 'none' scene turns off all lights
	if sceneName == "none" {
		nope := false
		s.hueBridge.SetGroupState(0, &hue.Action{On: &nope})
		return
	}

	if err := s.hueBridge.RecallSceneByName(sceneName); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

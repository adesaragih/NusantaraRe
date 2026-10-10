package handlers

import (
	"encoding/json"
	"net/http"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/galat"
	"nusantarare/modul/treatyin/backend/services"
)

// daftarkanLogAchievement - tombol `Submit` sub-tab Achievement (Treaty In dan
// Adjustment): `InsertToLogAchievement` → `LOG_ACHIEVEMENT` (keputusan pemakai
// 8 Oktober 2026).
//
//	POST /achievement/log  { idKontrak, baris: AchievementLists[] }
func daftarkanLogAchievement(pasang func(string, rute)) {
	pasang("POST "+Prefix+"/achievement/log", func(w http.ResponseWriter, r *http.Request, l *services.Layanan, p inti.Pelaku) {
		r.Body = http.MaxBytesReader(w, r.Body, batasBadanJSON)
		var m services.MasukanLogAchievement
		if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
			galat.Tulis(w, http.StatusBadRequest, "request body must be valid JSON")
			return
		}
		hasil, err := l.CatatLogAchievement(r.Context(), p, m)
		tulis(w, hasil, err)
	})
}

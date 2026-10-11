package handlers_test

import (
	"net/http"
	"testing"

	"nusantarare/modul/claimnonprop/backend/handlers"
)

// Tabel komite di bawah inbox (perintah work owner 09-10-2026, pola Claim Prop: menu Komite Claim Non Prop
// disembunyikan): hak `komite` hanya bagi pemegang workbasket yang tercantum sebagai KomiteID roster NONPROP aktif.
func TestHakKomiteRosterWorkbasket(t *testing.T) {
	u := baruUji(t)
	for peran, mau := range map[string]bool{"": false, "UJI-LAIN": false, "ReasClaimTechDivHead": true,
		"UJI-LAIN,ReasClaimDeptHead": true} {
		kode, out := u.minta(http.MethodGet, handlers.Prefix+"/hak", admin, peran, nil)
		u.wajib(kode, http.StatusOK, out, "hak "+peran)
		if out["komite"] != mau {
			t.Errorf("hak komite peran %q = %v, mau %v", peran, out["komite"], mau)
		}
	}
}

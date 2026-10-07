package backend

import (
	"testing"

	"nusantarare/modul/edmtreatyin/backend/handlers"
)

// Gerbang Copy Old (handlers.superadmin) membaca hak menu modul ini lewat KODE menunya - wajib sama dengan nama modul
// (`M_LOGIN_GO_MENU.MENU_KODE`), kalau tidak View only tidak pernah berlaku.
func TestKodeMenuCopyOldSamaDenganNamaModul(t *testing.T) {
	if handlers.KodeMenu != Nama {
		t.Fatalf("handlers.KodeMenu %q, nama modul %q", handlers.KodeMenu, Nama)
	}
}

package handlers

import (
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"

	"nusantarare/inti/backend/galat"
)

// Kontrak envelope galat, SISI BACKEND.
//
// ⛔ Kenapa uji ini ada, dengan sebab yang sungguh terjadi: `handlers.galat`
// menulis {"galat": …} sejak tiket 01, tetapi klien React membaca `o.error`.
// Akibatnya SETIAP pesan galat backend - 401, 403, 409, 503, seluruhnya -
// jatuh ke teks bawaan "Permintaan ditolak backend" di layar. Pemakai melihat
// pita merah yang tidak menyebutkan apa pun, padahal backend mengirim kalimat
// yang tepat.
//
// Cacat itu tidak berbunyi di satu sisi mana pun: backend benar, klien benar
// menurut komentarnya sendiri, dan hanya PERTEMUANNYA yang salah. Karena itu
// kontraknya dikunci di KEDUA sisi - di sini, dan di
// `frontend/src/modul/claim-life/envelopegalat.test.ts`.

func TestEnvelopeGalatMemakaiKunciGalat(t *testing.T) {
	w := httptest.NewRecorder()
	galat.Tulis(w, 409, "perpindahan itu tidak ada di tangga kerja klaim")

	var isi map[string]any
	if err := json.NewDecoder(w.Body).Decode(&isi); err != nil {
		t.Fatalf("badan bukan JSON: %v", err)
	}

	pesan, ada := isi["galat"]
	if !ada {
		t.Fatalf("kunci \"galat\" tidak ada; badan = %v", isi)
	}
	if pesan != "perpindahan itu tidak ada di tangga kerja klaim" {
		t.Errorf("pesan = %v, mau kalimat yang dikirim apa adanya", pesan)
	}

	// ⛔ SATU kunci, tidak lebih. Envelope yang membawa `error` DAN `galat`
	// sekaligus adalah cara paling mudah membuat kedua sisi tidak pernah
	// benar-benar bertemu: masing-masing membaca kuncinya sendiri, dan
	// ketidakcocokannya tidak pernah terlihat lagi.
	if len(isi) != 1 {
		t.Errorf("envelope memuat %d kunci (%v), mau tepat 1", len(isi), isi)
	}
	if _, ada := isi["error"]; ada {
		t.Error("envelope memuat \"error\"; kuncinya \"galat\", dan dua kunci " +
			"berarti kedua sisi tidak pernah dipaksa bertemu")
	}
}

func TestEnvelopeGalatBertipeJSON(t *testing.T) {
	w := httptest.NewRecorder()
	galat.Tulis(w, 503, "database belum dikonfigurasi")
	if ct := w.Header().Get("Content-Type"); !strings.HasPrefix(ct, "application/json") {
		t.Errorf("Content-Type = %q, mau application/json", ct)
	}
	if w.Code != 503 {
		t.Errorf("kode = %d, mau 503", w.Code)
	}
}

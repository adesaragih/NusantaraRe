package services

// Kalimat layar Treaty Contract Out BERBAHASA INGGRIS [keputusan work owner
// 30-09-2026: "untuk bahasa pake bahasa inggris, jangan indo"].
//
// Kalimat modul ini sudah Inggris. Kalimat sentinel milik `inti/` (bersama;
// tetap berbahasa Indonesia untuk modul lain) yang ikut di rantai galat Treaty
// diganti padanannya di SATU tempat ini - dipakai handler (`pesanTCO`) dan
// galat lampiran (ditulis `ringkasGalatTCO`, ditampilkan `TampilLampiran`,
// termasuk baris lama yang tersimpan berbahasa Indonesia).

import (
	"strings"

	"nusantarare/inti/db"
	"nusantarare/inti/galat"
	"nusantarare/inti/layanan"
	"nusantarare/inti/outbox"
	"nusantarare/inti/unggah"
)

// pesanIntiInggrisTCO - kalimat asal (`Error()` sentinel `inti/`) -> padanan Inggris.
var pesanIntiInggrisTCO = []struct{ asal, inggris string }{
	{galat.ErrPermintaanTidakSah.Error(), "services: invalid request"},
	{unggah.ErrUnggahanDirBelumDisetel.Error(), "services: UNGGAHAN_DIR is not set; document upload refused"},
	{unggah.ErrBerkasTerlaluBesar.Error(), "services: file exceeds the size limit"},
	{unggah.ErrBerkasKosong.Error(), "services: file is empty"},
	{db.ErrMataUangTidakDikenal.Error(), "repository: unknown currency code"},
	{db.ErrTanpaOracle.Error(), "repository: ORACLE_DSN is not configured"},
	{layanan.ErrGaramTokenKosong.Error(), "services: STORAGE_TOKEN_SALT is empty; the storage token cannot be issued"},
	{layanan.ErrAppNameKosong.Error(), "services: APPNAME is empty; the storage token requires it"},
	{outbox.ErrPenyimpananBelumDisetujui.Error(),
		"services: file storage is not approved yet; connecting it and calling GET_TOKEN_STORAGE requires human approval"},
	// Pembungkus `unggah.SimpanBerkas` (bukan sentinel).
	{"services: menulis berkas unggahan", "services: writing the uploaded file"},
}

// TeksInggrisTCO mengganti kalimat sentinel `inti/` di dalam teks galat.
func TeksInggrisTCO(s string) string {
	for _, g := range pesanIntiInggrisTCO {
		s = strings.ReplaceAll(s, g.asal, g.inggris)
	}
	return s
}

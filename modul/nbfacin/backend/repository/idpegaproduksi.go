package repository

// IDPEGA = case id POLOS (keputusan work owner 11-10-2026, PENYIMPANGAN SADAR dari pzInsKey Pega "ASM-FW-GISFW-WORK " + ID): "semua
// yang masuk ke facinproduction dan json_polis, kolom idpega nya sesuai case id", diperluas "riwayat akseptasi, FACINOFFER, dan
// JSON_OFFER UBAH JUGA" dan "FACINPERFORMANCE / DOCUMENT_POLIS IYA ITU UBAH JUGA". Semua tulisan IDPEGA / ID_PEGA nbfacin memakai
// case id polos (= T_WORK_POLIS.ID, mis. "NB-151033"): JSON_POLIS, FACINPRODUCTION, CEDING_FACINPRODUCTION, HISTORYAKSEPTASIPEGA,
// HISTORYAKSEPTASIPRODUCTION, FACINPERFORMANCE, JSON_OFFER, FACINOFFER (sudah polos = pyID), DOCUMENT_POLIS (jalur nbfacin saja;
// inti/backend/dokumenpolis tidak disentuh). Baris lama (Pega / sistem baru sebelum keputusan ini) tetap berawalan, jadi SETIAP
// pembaca menurut IDPEGA mencocokkan DUA bentuk (IN (:polos, :berawalan)). Tanpa migrasi dan tanpa UPDATE data dari kode (usulan
// penyeragaman DBA: docs/alat/seragamkan_idpega_nbfacin.sql).

import (
	"strconv"
	"strings"
)

// IDPegaProduksi - IDPEGA tulisan sistem baru ke JSON_POLIS / FACINPRODUCTION / CEDING_FACINPRODUCTION = case id polos.
func IDPegaProduksi(caseID string) string { return strings.TrimPrefix(caseID, AwalanIDPega) }

// DuaBentukIDPega - (polos, berawalan) dari case id ATAU IDPEGA tersimpan bentuk mana pun.
func DuaBentukIDPega(v string) (string, string) {
	p := strings.TrimPrefix(v, AwalanIDPega)
	return p, AwalanIDPega + p
}

// sqlCaseIDDariIDPega - ekspresi case id dari kolom IDPEGA dua bentuk: berawalan -> SUBSTR(kolom, 19) (Java substring 18, seperti
// sebelumnya); polos -> kolom apa adanya.
func sqlCaseIDDariIDPega(kolom string) string {
	return "CASE WHEN SUBSTR(" + kolom + ", 1, " + strconv.Itoa(len(AwalanIDPega)) + ") = '" + AwalanIDPega + "' THEN SUBSTR(" + kolom + ", " +
		strconv.Itoa(len(AwalanIDPega)+1) + ") ELSE " + kolom + " END"
}

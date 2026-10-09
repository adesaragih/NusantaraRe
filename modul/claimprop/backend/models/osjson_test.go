package models_test

import (
	"encoding/json"
	"strings"
	"testing"

	"nusantarare/modul/claimprop/backend/models"
	"nusantarare/modul/claimprop/backend/tiruan"
)

// jsonPega merangkai bentuk `@ASM.GetPageJSONString()` yang dibaca dari DATA_JSON OS_AKSEPTASI_KLAIM warisan.
func jsonPega(pasangan ...string) string {
	return "{\n" + strings.Join(pasangan, "\n,") + "\n}\n"
}

// Kunci urut tanpa membedakan huruf besar (PolicyNo < pxObjClass < pzInsKey < Type), nilai kosong tidak ditulis,
// garis miring dan tanda HTML tidak di-escape, kutip di-escape.
func TestJSONHalamanPegaBentukWarisan(t *testing.T) {
	dapat := models.JSONHalamanPega(map[string]string{"Type": "0", "pzInsKey": "UJI-K", "PolicyNo": "UJI/1",
		"pxObjClass": "UJI-C", "PersenRNM": "10", "Kosong": "", "Catatan": `a "b" <c>`})
	sama(t, "JSON", dapat, jsonPega(`"Catatan":"a \"b\" <c>"`, `"PersenRNM":"10"`, `"PolicyNo":"UJI/1"`,
		`"pxObjClass":"UJI-C"`, `"pzInsKey":"UJI-K"`, `"Type":"0"`))
	if !json.Valid([]byte(dapat)) {
		t.Fatalf("bukan JSON sah: %q", dapat)
	}
}

// SaveOutstanding_Act 23.1.2-23.1.3: DATA_JSON = halaman TempOSAkseptasi. NoClaim di JSON = ClaimData.NoClaim
// (23.1.2), bukan cadangan ClaimNo milik kolom NOCLAIM (23.1.5).
func TestKirimEstimasiMengisiDataJSON(t *testing.T) {
	h := halamanUji()
	h.Setel(models.CD+"NoClaim", "UJI-K1")
	h.AmbilDaftar(models.DaftarEstimasi)[0]["EstimationValue"] = "100"
	rows := models.KirimEstimasi(konteksUji(tiruan.AcuanBaru()), h, "CLMP-000001")
	if len(rows) != 1 {
		t.Fatalf("baris OS = %d, mau 1", len(rows))
	}
	sama(t, "DATA_JSON", rows[0].DataJSON, jsonPega(`"Currency":"UJA"`, `"CurrencyID":"UJI-A"`,
		`"EstimationDate":"20260520"`, `"GrossValue":"1000"`, `"KursValue":"2"`, `"NoClaim":"UJI-K1"`,
		`"PersenRNM":"10"`, `"pxCreateOperator":"UJI-ADMIN"`, `"pxObjClass":"ASM-FW-GCNMFW-Data-osAkseptasi"`,
		`"Type":"0"`, `"TypeID":"1"`, `"TypeLossID":"UJI-QS"`, `"Value":"100"`))

	h = halamanUji()
	h.Setel(models.CD+"ClaimNo", "UJI-KT1")
	h.AmbilDaftar(models.DaftarEstimasi)[0]["EstimationValue"] = "100"
	rows = models.KirimEstimasi(konteksUji(tiruan.AcuanBaru()), h, "CLMP-000001")
	sama(t, "NOCLAIM", rows[0].NoClaim, "UJI-KT1")
	if strings.Contains(rows[0].DataJSON, `"NoClaim"`) {
		t.Errorf("NoClaim kosong tetap tertulis di DATA_JSON: %q", rows[0].DataJSON)
	}
}

// CloseClaimProp 6.1-6.2: DATA_JSON = halaman InputParamOs (empat properti + kelas).
func TestBarisOSTutupMengisiDataJSON(t *testing.T) {
	h := models.HalamanBaru()
	h.Setel(models.CD+"CauseOfLoss", "UJI-SEBAB")
	h.Setel(models.CD+"CauseOfLossID", "UJI-S1")
	h.Setel(models.CD+"NoClaim", "UJI-K1")
	h.Setel(models.TM+"ID", "UJI-M1")
	b := models.BarisOSTutup(konteksUji(tiruan.AcuanBaru()), h, "CLMP-000001")
	sama(t, "DATA_JSON", b.DataJSON, jsonPega(`"CauseOfLoss":"UJI-SEBAB"`, `"CauseOfLossID":"UJI-S1"`,
		`"IDMasterTreaty":"UJI-M1"`, `"NoClaim":"UJI-K1"`, `"pxObjClass":"ASM-FW-GCNMFW-Data-osAkseptasi"`))
}

package models

// Untuk apa berkas ini: BARIS OS_AKSEPTASI_KLAIM dan JSON halaman Pega (`@ASM.GetPageJSONString()`).
//
// Procedure `PEGA_JSON_OS_AKSEP_KLAIM` (ALL_SOURCE DEV 10-10-2026) ditulis ulang tanpa procedure: SELALU INSERT (hitungan
// baris CASEID di awal procedure tidak dipakai) kolom CASEID, NOCLAIM, DATA_JSON, TANGGAL (hari ini), NOPOLIS,
// STS_REJECT, STS_KONVERSI, STS_DLA, MASTERID, CLAIMOLD; TGL_PROD diisi trigger `TRG_TLG_PROD_OS_AKSEPTASI`
// (BEFORE INSERT). Tiga pemanggil di Claim Fac In: `SaveCFS_ACT` (STS 0), `CloseClaim` (STS 4), `SaveAcceptation`
// (panggilannya ber-remark - mati, grilling ronde 5 §D3).

import (
	"bytes"
	"encoding/json"
	"sort"
	"strings"
)

// Status OS_AKSEPTASI_KLAIM.STS_REJECT (= parameter CARI10 procedure).
const (
	StsOSOutstanding = "0" // SaveCFS_ACT
	StsOSFinal       = "4" // CloseClaim
)

// STS_DLA baris OS (`InputData.CARI17 = @if(item.IsFacretro=="1", 8, 7)`, SaveCFS_ACT langkah 17).
const (
	StsDLAFac   = "7"
	StsDLARetro = "8"
)

// KelasOSAkseptasi - pxObjClass halaman TempOSAkseptasi (DATA_JSON baris OS `CLM-` DEV 10-10-2026).
const KelasOSAkseptasi = "ASM-FW-GCNMFW-Data-osAkseptasi"

// BarisOS - satu baris OS_AKSEPTASI_KLAIM (kolom yang ditulis procedure). CASEID = `KunciInstans` (ID kasus apa adanya,
// prompt §6 butir 6); MASTERID / CLAIMOLD = `InputData.CARI18` / `CARI20` yang tidak diisi pemanggil mana pun di Claim Fac
// In - ditulis NULL.
type BarisOS struct {
	CaseID, NoClaim, NoPolis string
	StsReject, StsDLA        string
	DataJSON                 string
}

// HalamanJSON - satu halaman untuk `@ASM.GetPageJSONString()`: nilai teks + daftar halaman (PageList).
type HalamanJSON struct {
	Nilai  map[string]string
	Daftar []DaftarJSON
}

// DaftarJSON - satu PageList bernama.
type DaftarJSON struct {
	Nama string
	Isi  []HalamanJSON
}

// JSONHalamanPega = `@ASM.GetPageJSONString()`. Fungsinya tidak diekspor; bentuknya dibaca dari DATA_JSON
// OS_AKSEPTASI_KLAIM `CLM-` DEV 10-10-2026 (sama dengan Claim Prop / Claim Non Prop, disalin bukan diimpor):
//   - "{" LF, pasangan pertama, setiap pasangan berikutnya diawali LF ",", ditutup LF "}"; halaman puncak diakhiri LF;
//   - nilai teks dulu, urut kunci tanpa membedakan huruf besar; properti kosong tidak ditulis;
//   - PageList sesudah nilai teks: `"Nama":[ ` LF anggota LF `] ` (anggota dipisah LF ",").
func JSONHalamanPega(p HalamanJSON) string {
	return tulisHalamanJSON(p) + "\n"
}

func tulisHalamanJSON(p HalamanJSON) string {
	kunci := make([]string, 0, len(p.Nilai))
	for k, v := range p.Nilai {
		if v != "" {
			kunci = append(kunci, k)
		}
	}
	sort.Slice(kunci, func(i, j int) bool {
		a, b := strings.ToLower(kunci[i]), strings.ToLower(kunci[j])
		if a != b {
			return a < b
		}
		return kunci[i] < kunci[j]
	})
	var bagian []string
	for _, k := range kunci {
		bagian = append(bagian, literalJSON(k)+":"+literalJSON(p.Nilai[k]))
	}
	for _, d := range p.Daftar {
		if len(d.Isi) == 0 {
			continue
		}
		var isi []string
		for _, a := range d.Isi {
			isi = append(isi, tulisHalamanJSON(a))
		}
		bagian = append(bagian, literalJSON(d.Nama)+":[ \n"+strings.Join(isi, "\n,")+"\n] ")
	}
	return "{\n" + strings.Join(bagian, "\n,") + "\n}"
}

// literalJSON - literal teks JSON tanpa escape HTML.
func literalJSON(s string) string {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	_ = enc.Encode(s)
	return strings.TrimSuffix(buf.String(), "\n")
}

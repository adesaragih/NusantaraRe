package models

// Untuk apa berkas ini: DATA LAMA - bahan murni pemuat kasus Claim Prop warisan Pega (`backend/alat/pemuatlama`,
// prompt implementasi §6 butir 11, AC 9-12, 123, 132).
//
// Dua sumber, keduanya baca-saja:
//
//	OS_AKSEPTASI_KLAIM  baris per aksi (CASEID `ASM-FW-GCNMFW-WORK CLMP-n`); DATA_JSON-nya JSON DATAR satu baris
//	                    estimasi / adjustment (katalog DEV 07-10-2026: 24 kunci), bukan halaman kasus.
//	JSON_KLAIM          halaman `.ClaimData` utuh (InsertJsonClaimTreaty_act) - hanya sebagian kecil kasus.
//
// Kasus berpengenal pyID Pega apa adanya ("CLMP-4894", `IDDariKunciPega`): nomor lama 1-4892 tidak pernah sama
// dengan ID baru `CLMP-` + LPAD 6.

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"nusantarare/inti/backend/utils"
)

// ErrDokumenLama - DATA_JSON warisan tidak terbaca sebagai objek JSON.
var ErrDokumenLama = errors.New("models: dokumen lama tak terbaca")

// AwalanKunciLama - awalan CASEID / IDPEGA kasus Claim Prop warisan Pega.
const AwalanKunciLama = KelasKasusKunci + " " + AwalanKlaim

// BarisOSLama - satu baris OS_AKSEPTASI_KLAIM kasus lama.
type BarisOSLama struct {
	CaseID, NoClaim, NoPolis, MasterID, InsertOp, StsReject string
	// AcceptedNo - `DATA_JSON.AcceptedNo` (kunci urutan AC 123).
	AcceptedNo string
	DataJSON   string
	Tanggal    time.Time
}

// UrutkanOSLama mengurutkan baris satu kasus menurut AC 123: TANGGAL naik, lalu AcceptedNo turun - kosong lebih dulu,
// sama dengan `ORDER BY ... DESC` Oracle (NULLS FIRST) di `RDBList/DataOutstandingTreatyin.xml`. Stabil.
func UrutkanOSLama(rows []BarisOSLama) {
	sort.SliceStable(rows, func(i, j int) bool {
		a, b := rows[i], rows[j]
		if !a.Tanggal.Equal(b.Tanggal) {
			return a.Tanggal.Before(b.Tanggal)
		}
		if (a.AcceptedNo == "") != (b.AcceptedNo == "") {
			return a.AcceptedNo == ""
		}
		return a.AcceptedNo > b.AcceptedNo
	})
}

// BarisBerlaku - baris yang berlaku satu kasus (AC 123): baris TERAKHIR urutan `UrutkanOSLama`, yaitu aksi terbaru.
// `[dugaan]` AC 123 menetapkan urutannya, bukan ujung mana yang berlaku (OQ-CP-18). Sisa barisnya riwayat; ia tetap di
// OS_AKSEPTASI_KLAIM dan terbaca lewat Summary Outstanding (`RingkasanOS`), tidak disalin.
func BarisBerlaku(rows []BarisOSLama) (BarisOSLama, int) {
	if len(rows) == 0 {
		return BarisOSLama{}, 0
	}
	s := append([]BarisOSLama(nil), rows...)
	UrutkanOSLama(s)
	return s[len(s)-1], len(s) - 1
}

// TahapLama memetakan STS_REJECT baris berlaku ke tahap kasus. `[dugaan]` (OQ-CP-18): baris OS mencatat aksi, bukan
// assignment - Submit tidak menulis baris, jadi kasus berbaris terakhir 0 bisa saja sudah di Input Acceptation.
// STS_REJECT di luar 0 / 2 / 4 tidak punya penulis di ekspor (spec §1b butir 6) - `ok` palsu, kasus tidak dimuat.
func TahapLama(sts string) (tahap string, tutup, ok bool) {
	switch strings.TrimSpace(sts) {
	case StsOSEstimasi:
		return TahapOutstanding, false, true
	case StsOSAkseptasi:
		return TahapAcceptation, false, true
	case StsOSTutupBerkas:
		return TahapAcceptation, true, true
	}
	return "", false, false
}

// MedanDibuang - satu medan dokumen lama yang tidak masuk kolom, beserta sebabnya (arsip audit pemuatan).
type MedanDibuang struct {
	Jalur, Nilai, Sebab string
}

// Sebab medan dokumen lama tidak dimuat.
const (
	SebabInternalPega = "dibuang: kunci px/py/pz (AC 10-11)"
	SebabTidakDiimpor = "dibuang: tidak diimpor - InterestListDtl, SpreadingRisk, PaymentData (keputusan diagram)"
	SebabKomite       = "dibuang: keputusan anggota komite kasus lama (tangga kasus komite Pega tidak dimuat)"
	SebabLampiran     = "dibuang: lampiran (OQ-CP-12)"
	SebabTanpaKolom   = "dibuang: tanpa kolom di katalog (turunan dihitung ulang, atau bukan data)"
	SebabBarisOS      = "tetap di OS_AKSEPTASI_KLAIM: medan baris estimasi / adjustment, bukan header kasus"
	// SebabStsTakDikenal - galat: baris berlaku ber-STS_REJECT tanpa penulis di ekspor, tahapnya tidak dapat diturunkan.
	SebabStsTakDikenal = "galat: STS_REJECT di luar 0 / 1 / 2 / 4 - tahap tidak diturunkan (OQ-CP-18)"
	// SebabStsDitunda - baris berlaku ber-STS_REJECT 1: dibahas bersama modul Komite Claim Prop (keputusan work owner
	// 07-10-2026 "itu nanti kan dari komite, lewatkan dulu claim prop"). Tidak dimuat, tidak dihitung gagal.
	SebabStsDitunda = "ditunda: STS_REJECT = 1 dibahas bersama modul Komite Claim Prop (keputusan work owner 07-10-2026)"
)

// StsOSDitunda - STS_REJECT baris OS yang pemuatannya ditunda sampai modul Komite Claim Prop.
const StsOSDitunda = "1"

// tidakDiimpor - nama halaman yang tidak diimpor dari JSON lama (prompt implementasi §3, baris "diagram").
var tidakDiimpor = map[string]bool{"InterestListDtl": true, "SpreadingRisk": true, "PaymentData": true}

func kunciInternal(k string) bool {
	return strings.HasPrefix(k, "px") || strings.HasPrefix(k, "py") || strings.HasPrefix(k, "pz")
}

// UraiKlaimLama membaca JSON_KLAIM.DATA_JSON (halaman `.ClaimData`) menjadi Halaman berjalur `ClaimData.*`. Kunci
// px/py/pz dan halaman yang tidak diimpor dibuang; setiap nilai teks apa adanya, termasuk tanggal Pega dua format
// (AC 12) - konversi terjadi di repository lewat katalog. Medan yang tidak punya kolom dilaporkan, tidak disimpan.
func UraiKlaimLama(teks string) (*Halaman, []MedanDibuang, error) {
	d := json.NewDecoder(strings.NewReader(teks))
	d.UseNumber()
	var akar map[string]any
	if err := d.Decode(&akar); err != nil {
		return nil, nil, fmt.Errorf("%w: JSON_KLAIM bukan objek JSON: %v", ErrDokumenLama, err)
	}
	h := HalamanBaru()
	var dibuang []MedanDibuang
	u := urai{h: h, dibuang: &dibuang}
	u.objek(strings.TrimSuffix(CD, "."), akar)
	// SuggestList ("Claim History") disimpan lewat T_VIEW_SUGGEST, bukan katalog: barisnya ditandai baru supaya
	// repository menyisipkannya; medan selain empat kolom riwayat dilaporkan dibuang.
	var riw []Baris
	for i, b := range h.AmbilDaftar(DaftarRiwayat) {
		nb := Baris{PropRiwayatBaru: "1"}
		for _, k := range kunciUrut(b) {
			switch k {
			case "CommentSuggest", "PICSuggest", "DateSuggest", "IsCedingConfirm":
				nb[k] = b[k]
			default:
				if b[k] != "" {
					dibuang = append(dibuang, MedanDibuang{Jalur: fmt.Sprintf("%s(%d).%s", DaftarRiwayat, i+1, k), Nilai: b[k],
						Sebab: SebabTanpaKolom})
				}
			}
		}
		riw = append(riw, nb)
	}
	h.SetelDaftar(DaftarRiwayat, nil)
	dibuang = append(dibuang, tanpaKolom(h, ProyeksiKatalog(h))...)
	h.SetelDaftar(DaftarRiwayat, riw)
	return h, dibuang, nil
}

type urai struct {
	h       *Halaman
	dibuang *[]MedanDibuang
}

func (u urai) buang(jalur string, v any, sebab string) {
	switch x := v.(type) {
	case map[string]any:
		for _, k := range kunciUrut(x) {
			u.buang(jalur+"."+k, x[k], sebab)
		}
	case []any:
		for i, e := range x {
			u.buang(fmt.Sprintf("%s(%d)", jalur, i+1), e, sebab)
		}
	default:
		*u.dibuang = append(*u.dibuang, MedanDibuang{Jalur: jalur, Nilai: skalar(v), Sebab: sebab})
	}
}

// objek menulis anggota objek `awalan` ke halaman (`awalan.kunci`).
func (u urai) objek(awalan string, o map[string]any) {
	for _, k := range kunciUrut(o) {
		v, jalur := o[k], awalan+"."+k
		switch {
		case kunciInternal(k):
			u.buang(jalur, v, SebabInternalPega)
		case tidakDiimpor[k]:
			u.buang(jalur, v, SebabTidakDiimpor)
		default:
			switch x := v.(type) {
			case map[string]any:
				u.objek(jalur, x)
			case []any:
				u.daftar(jalur, x)
			default:
				u.h.Setel(jalur, skalar(v))
			}
		}
	}
}

// daftar menulis PageList `jalur`; daftar bersarang di dalam baris menjadi `jalur(n).anak` (`JalurAnak`).
func (u urai) daftar(jalur string, d []any) {
	var rows []Baris
	for i, e := range d {
		o, ok := e.(map[string]any)
		if !ok {
			u.buang(fmt.Sprintf("%s(%d)", jalur, i+1), e, SebabTanpaKolom)
			continue
		}
		b := Baris{}
		u.baris(fmt.Sprintf("%s(%d)", jalur, len(rows)+1), "", o, b)
		rows = append(rows, b)
	}
	u.h.SetelDaftar(jalur, rows)
}

// baris menulis anggota objek ke satu baris; objek bersarang menjadi kunci bertitik ("DataCommitteeTreaty.Remarks").
func (u urai) baris(jalurBaris, awalan string, o map[string]any, b Baris) {
	for _, k := range kunciUrut(o) {
		v, kunci := o[k], awalan+k
		switch {
		case kunciInternal(k):
			u.buang(jalurBaris+"."+kunci, v, SebabInternalPega)
		case tidakDiimpor[k]:
			u.buang(jalurBaris+"."+kunci, v, SebabTidakDiimpor)
		default:
			switch x := v.(type) {
			case map[string]any:
				u.baris(jalurBaris, kunci+".", x, b)
			case []any:
				u.daftar(jalurBaris+"."+kunci, x)
			default:
				b[kunci] = skalar(v)
			}
		}
	}
}

// kunciUrut - kunci peta terurut (urutan laporan dan urai yang stabil).
func kunciUrut[V any](o map[string]V) []string {
	k := make([]string, 0, len(o))
	for x := range o {
		k = append(k, x)
	}
	sort.Strings(k)
	return k
}

func skalar(v any) string {
	switch x := v.(type) {
	case nil:
		return ""
	case string:
		return x
	case json.Number:
		return x.String()
	case bool:
		return strconv.FormatBool(x)
	}
	return fmt.Sprint(v)
}

// sebabTanpaKolom - sebab medan tanpa kolom menurut halamannya.
func sebabTanpaKolom(jalur string) string {
	switch {
	case strings.Contains(jalur, "ComiteeClaim") || strings.HasPrefix(jalur, CD+"ClaimComitee"):
		return SebabKomite
	case strings.HasPrefix(jalur, CD+"Attachment"):
		return SebabLampiran
	}
	return SebabTanpaKolom
}

// tanpaKolom - medan `h` yang tidak ada di proyeksi katalognya `s`.
func tanpaKolom(h, s *Halaman) []MedanDibuang {
	var out []MedanDibuang
	for _, j := range kunciUrut(h.Nilai) {
		if _, ada := s.Nilai[j]; !ada && h.Nilai[j] != "" {
			out = append(out, MedanDibuang{Jalur: j, Nilai: h.Nilai[j], Sebab: sebabTanpaKolom(j)})
		}
	}
	for _, j := range kunciUrut(h.Daftar) {
		simpan := s.Daftar[j]
		for i, b := range h.Daftar[j] {
			for _, k := range kunciUrut(b) {
				if b[k] == "" {
					continue
				}
				if i < len(simpan) {
					if _, ada := simpan[i][k]; ada {
						continue
					}
				}
				jalur := fmt.Sprintf("%s(%d).%s", j, i+1, k)
				out = append(out, MedanDibuang{Jalur: jalur, Nilai: b[k], Sebab: sebabTanpaKolom(jalur)})
			}
		}
	}
	return out
}

// petaOSLama - kunci DATA_JSON baris OS yang ADA di header kasus (sisanya medan baris estimasi / adjustment).
var petaOSLama = map[string]string{
	"NoClaim":        CD + "NoClaim",
	"PolicyNo":       CD + "PolicyData.PolicyNo",
	"IDMasterTreaty": CD + "IDMaster",
	"CauseOfLoss":    CD + "CauseOfLoss",
	"CauseOfLossID":  CD + "CauseOfLossID",
}

// HalamanDariOSLama - header kasus dari baris OS yang berlaku, untuk kasus tanpa JSON_KLAIM. Kolom datar (NOCLAIM,
// NOPOLIS, MASTERID) didahulukan atas kunci JSON yang sama.
func HalamanDariOSLama(b BarisOSLama) (*Halaman, []MedanDibuang, error) {
	h := HalamanBaru()
	var dibuang []MedanDibuang
	if strings.TrimSpace(b.DataJSON) != "" {
		d := json.NewDecoder(strings.NewReader(b.DataJSON))
		d.UseNumber()
		var o map[string]any
		if err := d.Decode(&o); err != nil {
			return nil, nil, fmt.Errorf("%w: DATA_JSON baris OS bukan objek JSON: %v", ErrDokumenLama, err)
		}
		for _, k := range kunciUrut(o) {
			switch jalur, ada := petaOSLama[k]; {
			case kunciInternal(k):
				dibuang = append(dibuang, MedanDibuang{Jalur: k, Nilai: skalar(o[k]), Sebab: SebabInternalPega})
			case ada:
				h.Setel(jalur, skalar(o[k]))
			default:
				dibuang = append(dibuang, MedanDibuang{Jalur: k, Nilai: skalar(o[k]), Sebab: SebabBarisOS})
			}
		}
	}
	for jalur, v := range map[string]string{CD + "NoClaim": b.NoClaim, CD + "PolicyData.PolicyNo": b.NoPolis,
		CD + "IDMaster": b.MasterID} {
		if strings.TrimSpace(v) != "" {
			h.Setel(jalur, v)
		}
	}
	return h, dibuang, nil
}

// GalatNilaiKatalog - medan berkolom yang akan ditolak saat ditulis (aturan `repository.nilaiTulis`): angka tak
// terbaca, tanggal tak terbaca, teks melebihi panjang kolom. Dipakai uji-kering supaya galat tulis terlihat tanpa
// menulis. Sebabnya tidak memuat nilai.
func GalatNilaiKatalog(h *Halaman) []MedanDibuang {
	var out []MedanDibuang
	periksa := func(jalur string, k Kolom, v string) {
		s := strings.TrimSpace(v)
		switch {
		case (k.Golongan.Desimal() || k.Golongan.Tanggal()) && s == "", v == "":
			return
		case k.Golongan.Desimal():
			if _, err := utils.ParseDecimal(s); err != nil {
				out = append(out, MedanDibuang{Jalur: jalur, Nilai: v, Sebab: "galat: bukan angka (" + k.Kolom + ")"})
			}
		case k.Golongan.Tanggal():
			if _, ok := UraiTanggal(s); !ok {
				out = append(out, MedanDibuang{Jalur: jalur, Nilai: v, Sebab: "galat: bukan tanggal (" + k.Kolom + ")"})
			}
		default:
			if k.Panjang > 0 && len([]rune(v)) > k.Panjang {
				out = append(out, MedanDibuang{Jalur: jalur, Nilai: v,
					Sebab: fmt.Sprintf("galat: melebihi %d karakter (%s)", k.Panjang, k.Kolom)})
			}
		}
	}
	for _, k := range TabelHeaderKlaim.Kolom {
		periksa(k.Properti, k, h.Ambil(k.Properti))
	}
	baris := func(t Tabel, jalur string) {
		for i, b := range h.AmbilDaftar(jalur) {
			for _, k := range t.Kolom {
				periksa(fmt.Sprintf("%s(%d).%s", jalur, i+1, k.Properti), k, b[k.Properti])
			}
		}
	}
	for _, t := range TabelAnakKlaim {
		baris(t, t.Daftar)
	}
	baris(TabelAdjustment, DaftarAdjustment)
	for i := range h.AmbilDaftar(DaftarAdjustment) {
		for _, c := range TabelCucuAdjustment {
			baris(c, JalurAnak(DaftarAdjustment, i+1, c.Daftar))
		}
	}
	return out
}

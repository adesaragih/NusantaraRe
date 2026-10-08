package models

// Untuk apa berkas ini: DATA LAMA (tiket 13) - sensus kasus komite warisan Pega untuk klaim Claim Prop yang pemuatannya
// DITUNDA pemuat Claim Prop (baris berlaku OS_AKSEPTASI_KLAIM ber-STS_REJECT 1, keputusan work owner 07-10-2026 "itu
// nanti kan dari komite"). Fungsi murni; pembacanya `repository/lama.go`, alatnya `alat/pemuatlama`.
//
// ⛔ JANGAN MENGARANG TANGGA (prompt §6 butir 12). Sumber yang ada di DEV 08-10-2026:
//   - HISTORYAKSEPTASIPEGA: ACCEPT / REJECT per Submit (ID_KOMITE, USERNAME = nama orang), TANPA tingkat, jabatan,
//     ID operator, dan anggota yang belum memutuskan - tangga tidak dapat disusun darinya;
//   - JSON_KLAIM: 7 kasus CLMP, nol `KomiteList`;
//   - DATAPEGA.PC_ASM_FW_GCNMFW_WORK: 10 work object KomiteTreaty dengan `KomiteList` utuh di BLOB PR7d - terbaca oleh
//     pengurai analisis, belum ada pengurainya di aplikasi (OQ).
// Karena itu alat ini hanya MENCACAH (uji-kering); `-jalankan` ditolak sampai sumber tangga diputuskan.

import (
	"sort"
	"strings"
	"time"
)

// AwalanKunciKlaimLama - awalan CASEID / ID_PEGA kasus Claim Prop warisan (`pzInsKey` Pega).
const AwalanKunciKlaimLama = "ASM-FW-GCNMFW-WORK " + AwalanKlaim

// StsOSDitunda - STS_REJECT baris berlaku yang pemuat Claim Prop serahkan ke modul ini.
const StsOSDitunda = "1"

// BarisOSLama - satu baris OS_AKSEPTASI_KLAIM kasus lama (kunci urutan AC 123 Claim Prop).
type BarisOSLama struct {
	CaseID, StsReject, AcceptedNo string
	Tanggal                       time.Time
}

// UrutkanOSLama - AC 123 Claim Prop (DISALIN, bukan diimpor): TANGGAL naik, lalu AcceptedNo turun - kosong lebih dulu.
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

// RiwayatLama - satu baris HISTORYAKSEPTASIPEGA kasus lama (tanpa USERNAME: nama orang tidak dibaca).
type RiwayatLama struct {
	IDPega, IDKomite, Status string
}

// KomitePega - satu work object KomiteTreaty di DATAPEGA (tanpa BLOB).
type KomitePega struct {
	Kunci, Cover, StatusWork string
}

// Sebab sensus.
const (
	SebabBlobBelumDiurai = "ditunda: tangga ada di BLOB work object KomiteTreaty DATAPEGA - pengurai PR7d belum ada di " +
		"aplikasi (OQ pemilik ekspor Pega)"
	SebabTanpaTangga = "ditunda: tangga komite lama tidak tersimpan (HISTORYAKSEPTASIPEGA tanpa tingkat / jabatan / " +
		"operator, JSON_KLAIM tanpa KomiteList) - OQ pemilik ekspor Pega"
	SebabTanpaKomite = "dilewati: tanpa jejak komite (riwayat akseptasi dan work object KomiteTreaty kosong)"
)

// RingkasanLama - satu klaim yang ditunda.
type RingkasanLama struct {
	Klaim string
	// BarisOS - cacah baris OS_AKSEPTASI_KLAIM klaim itu; Riwayat - cacah baris HISTORYAKSEPTASIPEGA.
	BarisOS, Riwayat int
	// KomiteRiwayat - ID_KOMITE berbeda di riwayat; KomitePega - work object KomiteTreaty ber-cover klaim ini.
	KomiteRiwayat, KomitePega int
	Sebab                     string
}

// SensusLama - hasil sensus.
type SensusLama struct {
	KlaimLama, Ditunda, DenganBlob, TanpaTangga, TanpaKomite int
	Baris                                                    []RingkasanLama
}

// Sensus mencacah klaim ber-baris berlaku STS_REJECT 1 dan jejak komitenya.
func Sensus(os []BarisOSLama, riwayat []RiwayatLama, pega []KomitePega) SensusLama {
	perKlaim := map[string][]BarisOSLama{}
	for _, b := range os {
		perKlaim[IDDariKunci(b.CaseID)] = append(perKlaim[IDDariKunci(b.CaseID)], b)
	}
	riw := map[string]int{}
	komiteRiw := map[string]map[string]bool{}
	for _, r := range riwayat {
		id := IDDariKunci(r.IDPega)
		riw[id]++
		if strings.TrimSpace(r.IDKomite) == "" {
			continue
		}
		if komiteRiw[id] == nil {
			komiteRiw[id] = map[string]bool{}
		}
		komiteRiw[id][r.IDKomite] = true
	}
	blob := map[string]int{}
	for _, p := range pega {
		blob[IDDariKunci(p.Cover)]++
	}
	var s SensusLama
	s.KlaimLama = len(perKlaim)
	kunci := make([]string, 0, len(perKlaim))
	for k := range perKlaim {
		kunci = append(kunci, k)
	}
	sort.Strings(kunci)
	for _, k := range kunci {
		rows := append([]BarisOSLama(nil), perKlaim[k]...)
		UrutkanOSLama(rows)
		if strings.TrimSpace(rows[len(rows)-1].StsReject) != StsOSDitunda {
			continue
		}
		s.Ditunda++
		r := RingkasanLama{Klaim: k, BarisOS: len(rows), Riwayat: riw[k], KomiteRiwayat: len(komiteRiw[k]),
			KomitePega: blob[k]}
		switch {
		case r.KomitePega > 0:
			r.Sebab = SebabBlobBelumDiurai
			s.DenganBlob++
		case r.Riwayat > 0:
			r.Sebab = SebabTanpaTangga
			s.TanpaTangga++
		default:
			r.Sebab = SebabTanpaKomite
			s.TanpaKomite++
		}
		s.Baris = append(s.Baris, r)
	}
	return s
}

// IDDariKunci - pyID dari `pzInsKey` Pega ("<kelas> <pyID>"): kunci OS / riwayat / cover dicocokkan lewat pyID.
func IDDariKunci(k string) string {
	k = strings.TrimSpace(k)
	if i := strings.LastIndex(k, " "); i >= 0 {
		return k[i+1:]
	}
	return k
}

// UsulLama - keputusan 19-09-2026 tiket 13: penanda usul halaman lama benar / salah / kosong -> '1' / '0' / '0'.
func UsulLama(v string) string {
	if strings.EqualFold(strings.TrimSpace(v), "true") {
		return UsulYa
	}
	return UsulTidak
}

package models

// POLICY_HOLDER dan PLAN peserta dari DATA POLIS, bukan dari CSV (keputusan
// work owner 05-10-2026) - supaya peserta tidak berbeda dari polisnya.
//
// Untuk apa berkas ini: aturan MURNI yang menimpa kedua kolom setiap baris
// unggahan dengan nilai polis, untuk SEMUA Type. Pega memakai nilai CSV apa
// adanya; di sini nilai CSV DIABAIKAN (kolomnya tidak wajib di judul dan
// isinya tidak divalidasi).
//
//	POLICY_HOLDER peserta = T_PREMIUM_LIST.POLICY_HOLDER_NAME (nama, bukan ID)
//	PLAN peserta          = T_PREMIUM_LIST.BUSINESS_NAME (Class of Business)

import (
	"errors"
	"fmt"
	"strings"
)

// ErrDataPolisPesertaKosong - Policy Holder atau Class of Business polis
// kosong; seluruh unggahan ditolak sebelum baris diproses (409 berpesan).
var ErrDataPolisPesertaKosong = errors.New("models: data polis untuk peserta kosong")

// Kolom CSV yang diisi dari polis - NAMA KOLOM, bukan data orang. Dinamai menurut
// perannya (pemegang polis) supaya tidak menyerupai pemberian teks tetap ke medan
// bernama-orang (penjaga `TestNolNamaOrangDiKode`, 07-10-2026).
const (
	KolomPesertaPemegangPolis = "POLICY_HOLDER"
	KolomPesertaPlan          = "PLAN"
)

// DataPolisPeserta - nilai polis yang disalin ke setiap peserta.
type DataPolisPeserta struct {
	PolicyHolder string // POLICY_HOLDER_NAME
	Plan         string // BUSINESS_NAME
}

// Periksa menolak polis yang Policy Holder atau Class of Business-nya kosong.
func (d DataPolisPeserta) Periksa() error {
	var kosong []string
	if strings.TrimSpace(d.PolicyHolder) == "" {
		kosong = append(kosong, "Policy Holder")
	}
	if strings.TrimSpace(d.Plan) == "" {
		kosong = append(kosong, "Class of Business")
	}
	switch len(kosong) {
	case 0:
		return nil
	case 1:
		return fmt.Errorf("%w: %s is empty on the policy - fill it and press Save Data before uploading the CSV",
			ErrDataPolisPesertaKosong, kosong[0])
	default:
		return fmt.Errorf("%w: %s are empty on the policy - fill them and press Save Data before uploading the CSV",
			ErrDataPolisPesertaKosong, strings.Join(kosong, " and "))
	}
}

// TerapkanDataPolisKePeserta menimpa POLICY_HOLDER dan PLAN setiap baris
// dengan nilai polis (TrimSpace), juga bila kolomnya tidak ada di CSV.
// Galat ErrDataPolisPesertaKosong bila salah satu nilai polis kosong - baris
// tidak disentuh.
func TerapkanDataPolisKePeserta(policyHolder, plan string, baris []BarisUnggah) error {
	d := DataPolisPeserta{PolicyHolder: policyHolder, Plan: plan}
	if err := d.Periksa(); err != nil {
		return err
	}
	d.Terapkan(baris)
	return nil
}

// Terapkan menimpa kedua kolom tanpa memeriksa - pemanggil sudah memanggil
// Periksa (dipakai sebagai kait periksaBerkas, sebelum validasi per baris).
func (d DataPolisPeserta) Terapkan(baris []BarisUnggah) {
	ph, plan := strings.TrimSpace(d.PolicyHolder), strings.TrimSpace(d.Plan)
	for i := range baris {
		if baris[i].Nilai == nil {
			baris[i].Nilai = map[string]string{}
		}
		baris[i].Nilai[KolomPesertaPemegangPolis] = ph
		baris[i].Nilai[KolomPesertaPlan] = plan
	}
}

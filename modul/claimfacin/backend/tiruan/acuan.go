package tiruan

// Untuk apa berkas ini: ACUAN TIRUAN - bacaan acuan (`services.Acuan`) dari peta dalam memori yang diisi uji. Semua
// nilai fixture berawalan `UJI-` atau kode buatan; nol data DEV.

import (
	"context"
	"strings"

	"nusantarare/modul/claimfacin/backend/models"
	"nusantarare/modul/claimfacin/backend/repository"
)

// KunciPolis - kunci peta dokumen polis tiruan.
func KunciPolis(nopolis, prodke string) string { return nopolis + "|" + prodke }

// Acuan - acuan tiruan.
type Acuan struct {
	Polis      []models.BarisPolisCari
	Dokumen    map[string][]byte // KunciPolis -> DATA_JSON
	Riwayat    map[string][]models.RiwayatKlaimPolis
	Ditolak    map[string]bool
	Kurs       map[string]string
	NamaMU     map[string]string
	JenisReas  map[string]string
	WilayahPos map[string][]models.BarisWilayah
	Adjuster   map[string]string
	Agen       map[string]string
	Alamat     map[string][4]string
	Roster     []models.AnggotaKomite
	Tingkat    map[string]string
	Nama       map[string]string
	Email      map[string]string
	OSKasus    map[string]string
	OSPolis    map[string]string
	Batas      map[string]models.BatasTreaty // kodeBisnis|jenis -> batas
	QS         map[string][]models.BarisQuotaShare
	Reasuradur map[string][]models.ReasuradurTreaty
	MO         map[string][6]string
	Sebab      []repository.BarisSebab
	Katastrofe []repository.BarisKatastrofe
	Negara     []models.Pilihan
	Provinsi   []models.Pilihan
	Kota       []models.Pilihan
	Distrik    []models.Pilihan
	RW         []models.Pilihan
	// Rekening - klien|mata uang -> baris BANKACCOUNT; RekeningMU - mata uang -> baris (Payable 3).
	Rekening   map[string][]models.RekeningBank
	RekeningMU map[string][]models.RekeningBank
	// Kasir - NOAKSEPTASI -> KET DIRECTTOKASIR_LOG; Konversi - nomor akseptasi tanpa titik -> status konversi
	// (getStatusKonversi_Act, "1" = sudah di TRLOSS_DETAIL_T).
	Kasir    map[string]string
	Konversi map[string]string
	// LimitDirut - LIMIT_BOTTOM Direktur Utama; Saldo - polis tanpa titik|mata uang -> saldo premi; Proteksi - polis
	// ber-proteksi premi dibuka.
	LimitDirut string
	Saldo      map[string]string
	Proteksi   map[string]bool
	// Tangga - kasus komite -> anggota (tangga tertulis gudang tiruan dibaca lewat Gudang).
	Tangga func(komiteID string) []models.AnggotaKomite
}

// AcuanBaru membuat acuan tiruan kosong.
func AcuanBaru() *Acuan {
	return &Acuan{Dokumen: map[string][]byte{}, Riwayat: map[string][]models.RiwayatKlaimPolis{}, Ditolak: map[string]bool{},
		Kurs: map[string]string{}, NamaMU: map[string]string{}, JenisReas: map[string]string{},
		WilayahPos: map[string][]models.BarisWilayah{}, Adjuster: map[string]string{}, Agen: map[string]string{},
		Alamat: map[string][4]string{}, Tingkat: map[string]string{}, Nama: map[string]string{}, Email: map[string]string{},
		OSKasus: map[string]string{}, OSPolis: map[string]string{}, Batas: map[string]models.BatasTreaty{},
		QS: map[string][]models.BarisQuotaShare{}, Reasuradur: map[string][]models.ReasuradurTreaty{},
		MO: map[string][6]string{}, Rekening: map[string][]models.RekeningBank{},
		RekeningMU: map[string][]models.RekeningBank{}, Kasir: map[string]string{}, Konversi: map[string]string{},
		Saldo: map[string]string{}, Proteksi: map[string]bool{}}
}

// StatusKonversi - lihat `repository.Acuan.StatusKonversi` (tiruan: peta Konversi).
func (a *Acuan) StatusKonversi(_ context.Context, noAksep string) (string, error) {
	return a.Konversi[noAksep], nil
}

// CariPolis - baris polis berkolom cocok.
func (a *Acuan) CariPolis(_ context.Context, jenis, teks string) ([]models.BarisPolisCari, error) {
	out := []models.BarisPolisCari{}
	for _, p := range a.Polis {
		var v string
		switch jenis {
		case models.CariNoPolis:
			if p.PolicyNo == teks {
				out = append(out, p)
			}
			continue
		case models.CariCeding:
			v = p.CedingCoName
		case models.CariInsured:
			v = p.CustomerName
		case models.CariQQ:
			v = p.QQ
		}
		if teks != "" && strings.HasPrefix(v, teks) {
			out = append(out, p)
		}
	}
	return out, nil
}

// AdaPolis - pasangan nomor + prodke ada di fixture polis.
func (a *Acuan) AdaPolis(_ context.Context, nopolis, prodke string) (bool, error) {
	for _, p := range a.Polis {
		if p.PolicyNo == nopolis && p.Prodke == prodke {
			return true, nil
		}
	}
	return false, nil
}

// DokumenPolis - DATA_JSON fixture.
func (a *Acuan) DokumenPolis(_ context.Context, nopolis, prodke string) ([]byte, bool, error) {
	d, ada := a.Dokumen[KunciPolis(nopolis, prodke)]
	return d, ada, nil
}

// RiwayatKlaimPolis - riwayat fixture.
func (a *Acuan) RiwayatKlaimPolis(_ context.Context, nopolis string) ([]models.RiwayatKlaimPolis, error) {
	return a.Riwayat[nopolis], nil
}

// KlaimDitolak - fixture CLAIMREJECTED.
func (a *Acuan) KlaimDitolak(_ context.Context, kunci string) (bool, error) {
	return a.Ditolak[kunci], nil
}

// KursStandar - kurs fixture (kosong = "1").
func (a *Acuan) KursStandar(_ context.Context, cur string) (string, error) {
	if v, ada := a.Kurs[cur]; ada {
		return v, nil
	}
	return "1", nil
}

// NamaMataUang - nama mata uang fixture.
func (a *Acuan) NamaMataUang(_ context.Context, cur string) (string, error) {
	return a.NamaMU[cur], nil
}

// NamaJenisReas - REINSURANCETYPE.NOTE fixture.
func (a *Acuan) NamaJenisReas(_ context.Context, id string) (string, error) {
	return a.JenisReas[id], nil
}

// Wilayah - BrowseRW_SQL fixture.
func (a *Acuan) Wilayah(_ context.Context, pos string) ([]models.BarisWilayah, error) {
	return a.WilayahPos[pos], nil
}

// NamaAdjuster - nama adjuster fixture.
func (a *Acuan) NamaAdjuster(_ context.Context, id string) (string, error) {
	return a.Adjuster[id], nil
}

// AgenKlien - agen fixture menurut nama klien.
func (a *Acuan) AgenKlien(_ context.Context, nama string) (string, error) { return a.Agen[nama], nil }

// AlamatAgen - alamat fixture.
func (a *Acuan) AlamatAgen(_ context.Context, agen string) ([4]string, error) {
	return a.Alamat[agen], nil
}

// RosterKomite - roster fixture.
func (a *Acuan) RosterKomite(context.Context) ([]models.AnggotaKomite, error) { return a.Roster, nil }

// TingkatPelaku - jabatan fixture.
func (a *Acuan) TingkatPelaku(_ context.Context, akun string) (string, error) {
	return a.Tingkat[akun], nil
}

// NamaPelaku - nama fixture (kosong = akun).
func (a *Acuan) NamaPelaku(_ context.Context, akun string) (string, error) {
	if v := a.Nama[akun]; v != "" {
		return v, nil
	}
	return akun, nil
}

// EmailPelaku - surel fixture.
func (a *Acuan) EmailPelaku(_ context.Context, akun string) (string, error) {
	return a.Email[akun], nil
}

// EstimasiKasusTerbuka - Σ TRLOSS fixture.
func (a *Acuan) EstimasiKasusTerbuka(_ context.Context, id string) (string, error) {
	return a.OSKasus[id], nil
}

// EstimasiPolisTerbuka - Σ TRLOSS polis fixture.
func (a *Acuan) EstimasiPolisTerbuka(_ context.Context, nopolis string) (string, error) {
	return a.OSPolis[nopolis], nil
}

// BatasTreaty - batas fixture.
func (a *Acuan) BatasTreaty(_ context.Context, kode, jenis, _ string) (models.BatasTreaty, bool, error) {
	b, ada := a.Batas[kode+"|"+jenis]
	return b, ada, nil
}

// QuotaShare - QS fixture menurut induk.
func (a *Acuan) QuotaShare(_ context.Context, _, _, induk string) ([]models.BarisQuotaShare, error) {
	return a.QS[induk], nil
}

// ReasuradurTreaty - reasuradur fixture menurut jenis.
func (a *Acuan) ReasuradurTreaty(_ context.Context, jenis, _, _ string) ([]models.ReasuradurTreaty, error) {
	return a.Reasuradur[jenis], nil
}

// MarketingOfficer - MO fixture.
func (a *Acuan) MarketingOfficer(_ context.Context, id string) ([6]string, bool, error) {
	m, ada := a.MO[id]
	return m, ada, nil
}

// DaftarMataUang - mata uang fixture.
func (a *Acuan) DaftarMataUang(context.Context) ([]models.Pilihan, error) {
	out := []models.Pilihan{}
	for id, n := range a.NamaMU {
		out = append(out, models.Pilihan{Nilai: id, Label: n})
	}
	return out, nil
}

// DaftarJenisReas - jenis reasuransi fixture.
func (a *Acuan) DaftarJenisReas(context.Context) ([]models.Pilihan, error) {
	out := []models.Pilihan{}
	for id, n := range a.JenisReas {
		out = append(out, models.Pilihan{Nilai: id, Label: n})
	}
	return out, nil
}

// DaftarAdjuster - adjuster fixture.
func (a *Acuan) DaftarAdjuster(context.Context, string) ([]models.Pilihan, error) {
	out := []models.Pilihan{}
	for id, n := range a.Adjuster {
		out = append(out, models.Pilihan{Nilai: id, Label: n})
	}
	return out, nil
}

func saringPilihan(xs []models.Pilihan, cari string) []models.Pilihan {
	out := []models.Pilihan{}
	for _, x := range xs {
		if strings.Contains(strings.ToUpper(x.Nilai), strings.ToUpper(strings.TrimSpace(cari))) {
			out = append(out, x)
		}
	}
	return out
}

// DaftarNegara - negara fixture.
func (a *Acuan) DaftarNegara(_ context.Context, cari string) ([]models.Pilihan, error) {
	return saringPilihan(a.Negara, cari), nil
}

// DaftarProvinsi - provinsi fixture.
func (a *Acuan) DaftarProvinsi(_ context.Context, _, cari string) ([]models.Pilihan, error) {
	return saringPilihan(a.Provinsi, cari), nil
}

// DaftarKota - kota fixture.
func (a *Acuan) DaftarKota(_ context.Context, _, cari string) ([]models.Pilihan, error) {
	return saringPilihan(a.Kota, cari), nil
}

// DaftarDistrik - distrik fixture.
func (a *Acuan) DaftarDistrik(_ context.Context, _, cari string) ([]models.Pilihan, error) {
	return saringPilihan(a.Distrik, cari), nil
}

// DaftarRW - RW fixture.
func (a *Acuan) DaftarRW(_ context.Context, _, cari string) ([]models.Pilihan, error) {
	return saringPilihan(a.RW, cari), nil
}

// DaftarSebab - cause of loss fixture.
func (a *Acuan) DaftarSebab(context.Context, string) ([]repository.BarisSebab, error) {
	return a.Sebab, nil
}

// BarisSebabID - cause of loss menurut ID.
func (a *Acuan) BarisSebabID(_ context.Context, id string) (repository.BarisSebab, bool, error) {
	for _, s := range a.Sebab {
		if s.ID == id {
			return s, true, nil
		}
	}
	return repository.BarisSebab{}, false, nil
}

// DaftarKatastrofe - katastrofe fixture.
func (a *Acuan) DaftarKatastrofe(context.Context, string) ([]repository.BarisKatastrofe, error) {
	return a.Katastrofe, nil
}

// BarisKatastrofeID - katastrofe menurut ID.
func (a *Acuan) BarisKatastrofeID(_ context.Context, id string) (repository.BarisKatastrofe, bool, error) {
	for _, k := range a.Katastrofe {
		if k.ID == id {
			return k, true, nil
		}
	}
	return repository.BarisKatastrofe{}, false, nil
}

// ProgresKlaim - progres kosong (tiruan tidak membaca PROGRESSCLAIM; asersi lewat Gudang.Progres).
func (a *Acuan) ProgresKlaim(context.Context, string) ([]models.Baris, map[int][]models.Baris, error) {
	return nil, map[int][]models.Baris{}, nil
}

// RekeningBank - rekening fixture menurut klien + mata uang.
func (a *Acuan) RekeningBank(_ context.Context, klien, cur string) ([]models.RekeningBank, error) {
	return a.Rekening[klien+"|"+cur], nil
}

// RekeningBankMataUang - rekening fixture menurut mata uang.
func (a *Acuan) RekeningBankMataUang(_ context.Context, cur string) ([]models.RekeningBank, error) {
	return a.RekeningMU[cur], nil
}

// IDBankRekening - IDOFBANK fixture menurut nama bank + cabang + nomor rekening.
func (a *Acuan) IDBankRekening(_ context.Context, bank, cabang, akun string) (string, error) {
	for _, rs := range a.Rekening {
		for _, r := range rs {
			if r.NameOfBank == bank && r.BranchOfBank == cabang && r.AccountNo == akun {
				return r.IDOfBank, nil
			}
		}
	}
	return "", nil
}

// StatusKasir - KET fixture.
func (a *Acuan) StatusKasir(_ context.Context, noAksep string) (string, bool, error) {
	v, ada := a.Kasir[noAksep]
	return v, ada, nil
}

// EmailCeding - kosong (hanya produksi).
func (a *Acuan) EmailCeding(context.Context, string) (string, error) { return "", nil }

// LimitDirekturUtama - batas fixture.
func (a *Acuan) LimitDirekturUtama(context.Context) (string, bool, error) {
	return a.LimitDirut, a.LimitDirut != "", nil
}

// SaldoPremi - saldo fixture.
func (a *Acuan) SaldoPremi(_ context.Context, invoice, cur string) (string, error) {
	return a.Saldo[invoice+"|"+cur], nil
}

// AdaProteksiPremi - proteksi fixture.
func (a *Acuan) AdaProteksiPremi(_ context.Context, nopolis string) (bool, error) {
	return a.Proteksi[nopolis], nil
}

// TanggaKomite - tangga kasus komite (dari gudang tiruan bila dipasang).
func (a *Acuan) TanggaKomite(_ context.Context, komiteID string) ([]models.AnggotaKomite, error) {
	if a.Tangga == nil {
		return nil, nil
	}
	return a.Tangga(komiteID), nil
}

package tiruan

// Untuk apa berkas ini: ACUAN TIRUAN - bacaan acuan (`services.Acuan`) dari peta dalam memori yang diisi uji. Semua
// nilai fixture berawalan `UJI-` atau kode buatan; nol data DEV.

import (
	"context"
	"strings"

	"nusantarare/modul/claimnonprop/backend/models"
	"nusantarare/modul/claimnonprop/backend/repository"
)

// Rekening - satu baris BANKACCOUNT tiruan.
type Rekening struct {
	models.RekeningBank
	ClientID string
}

// Acuan - acuan tiruan.
type Acuan struct {
	NamaMU      map[string]string
	JenisXOL    map[string]models.BarisReinsType // nama layer -> REINSURANCETYPE type 4
	Master      map[string]models.MasterTreaty
	BarisMaster []models.BarisMaster
	PolisMaster map[string]bool
	Polis       map[string][]models.BarisPolis // IDMaster -> polis
	TreatyNama  []string
	WilayahPos  map[string][]models.BarisWilayah
	Alamat      map[string][4]string
	Adjuster    map[string]string
	MO          map[string][6]string
	Riwayat     map[string][]models.RiwayatKlaimPolis
	OldID       map[string]string // nama bisnis -> OLDID
	Rekening    []Rekening
	Roster      []models.AnggotaKomite // urut DEGREE; Degree "1" = LIMIT_BOTTOM <= 0
	Tangga      map[string][]models.AnggotaKomite
	Berkas      map[string]models.BerkasPolis
	Sebab       []repository.BarisSebab
	Katastrofe  []repository.BarisKatastrofe
	Kasir       map[string]string
	Konversi    map[string]string
	Tingkat     map[string]string
	Nama        map[string]string
	Email       map[string]string
	XOL2        []models.BarisXOL2
	LampiranB   []repository.BarisLampiranBayar
}

// AcuanBaru membuat acuan tiruan kosong.
func AcuanBaru() *Acuan {
	return &Acuan{NamaMU: map[string]string{}, JenisXOL: map[string]models.BarisReinsType{},
		Master: map[string]models.MasterTreaty{}, PolisMaster: map[string]bool{}, Polis: map[string][]models.BarisPolis{},
		WilayahPos: map[string][]models.BarisWilayah{}, Alamat: map[string][4]string{}, Adjuster: map[string]string{},
		MO: map[string][6]string{}, Riwayat: map[string][]models.RiwayatKlaimPolis{}, OldID: map[string]string{},
		Tangga: map[string][]models.AnggotaKomite{}, Berkas: map[string]models.BerkasPolis{}, Kasir: map[string]string{},
		Konversi: map[string]string{}, Tingkat: map[string]string{}, Nama: map[string]string{}, Email: map[string]string{}}
}

// MasterTreaty - master menurut ID.
func (a *Acuan) MasterTreaty(_ context.Context, id string) (models.MasterTreaty, bool, error) {
	m, ada := a.Master[id]
	return m, ada, nil
}

// JenisReasXOL - REINSURANCETYPE type 4 menurut nama.
func (a *Acuan) JenisReasXOL(_ context.Context, nama string) ([]models.BarisReinsType, error) {
	if j, ada := a.JenisXOL[nama]; ada {
		return []models.BarisReinsType{j}, nil
	}
	return nil, nil
}

// NamaMataUang - CURRENCY menurut ID.
func (a *Acuan) NamaMataUang(_ context.Context, cur string) (string, error) {
	return a.NamaMU[cur], nil
}

// DaftarPolis - polis master.
func (a *Acuan) DaftarPolis(_ context.Context, idMaster string) ([]models.BarisPolis, error) {
	return append([]models.BarisPolis{}, a.Polis[idMaster]...), nil
}

// AdaPolisMaster - ada polis realisasi master.
func (a *Acuan) AdaPolisMaster(_ context.Context, id string) (bool, error) {
	return a.PolisMaster[id], nil
}

// NamaTreatySpreading - nama treaty spreading.
func (a *Acuan) NamaTreatySpreading(context.Context, string, string) ([]string, error) {
	return append([]string{}, a.TreatyNama...), nil
}

// Wilayah - wilayah menurut kode pos.
func (a *Acuan) Wilayah(_ context.Context, pos string) ([]models.BarisWilayah, error) {
	return a.WilayahPos[pos], nil
}

// NamaAdjuster - nama adjuster / consultant.
func (a *Acuan) NamaAdjuster(_ context.Context, id string) (string, error) {
	return a.Adjuster[id], nil
}

// MarketingPolis - MO polis.
func (a *Acuan) MarketingPolis(_ context.Context, nopolis string) ([6]string, bool, error) {
	m, ada := a.MO[nopolis]
	return m, ada, nil
}

// RiwayatKlaimPolis - klaim lain pada polis.
func (a *Acuan) RiwayatKlaimPolis(_ context.Context, nopolis string) ([]models.RiwayatKlaimPolis, error) {
	return a.Riwayat[nopolis], nil
}

// KodeLamaBisnis - OLDID menurut nama bisnis.
func (a *Acuan) KodeLamaBisnis(_ context.Context, nama string) (string, error) {
	return a.OldID[nama], nil
}

func (a *Acuan) rekening(f func(r Rekening) bool) []models.RekeningBank {
	out := []models.RekeningBank{}
	for _, r := range a.Rekening {
		if f(r) {
			out = append(out, r.RekeningBank)
		}
	}
	return out
}

// RekeningBank - CLIENTID + CURRENCYID.
func (a *Acuan) RekeningBank(_ context.Context, klien, cur string) ([]models.RekeningBank, error) {
	return a.rekening(func(r Rekening) bool { return r.ClientID == klien && r.CurrencyID == cur }), nil
}

// RekeningBankKlien - CLIENTID.
func (a *Acuan) RekeningBankKlien(_ context.Context, klien string) ([]models.RekeningBank, error) {
	return a.rekening(func(r Rekening) bool { return r.ClientID == klien }), nil
}

// RekeningBankNama - CLIENTNAME [+ CURRENCYID].
func (a *Acuan) RekeningBankNama(_ context.Context, nama, cur string) ([]models.RekeningBank, error) {
	return a.rekening(func(r Rekening) bool { return r.ClientName == nama && (cur == "" || r.CurrencyID == cur) }), nil
}

// RosterKomite - roster NONPROP; hanyaTingkat1 = DEGREE 1.
func (a *Acuan) RosterKomite(_ context.Context, hanyaTingkat1 bool) ([]models.AnggotaKomite, error) {
	var out []models.AnggotaKomite
	for _, r := range a.Roster {
		if !hanyaTingkat1 || r.Degree == "1" {
			out = append(out, r)
		}
	}
	return out, nil
}

// AlamatAgen - alamat klien agen.
func (a *Acuan) AlamatAgen(_ context.Context, agen string) ([4]string, error) {
	return a.Alamat[agen], nil
}

// NamaPelaku - nama tampilan akun.
func (a *Acuan) NamaPelaku(_ context.Context, akun string) (string, error) {
	if n := a.Nama[akun]; n != "" {
		return n, nil
	}
	return akun, nil
}

// DaftarMaster - baris popup master tersaring.
func (a *Acuan) DaftarMaster(_ context.Context, s models.SaringanMaster) ([]models.BarisMaster, error) {
	out := []models.BarisMaster{}
	for _, b := range a.BarisMaster {
		if s.TreatyID != "" && !strings.Contains(strings.ToUpper(b.TreatyID), strings.ToUpper(s.TreatyID)) {
			continue
		}
		out = append(out, b)
	}
	return out, nil
}

// BarisMasterDari - satu baris popup master.
func (a *Acuan) BarisMasterDari(_ context.Context, _, id, cob, grup string) (models.BarisMaster, bool, error) {
	for _, b := range a.BarisMaster {
		if b.TreatyID == id && b.ClassOfBusinessID == cob && b.TreatyGroup == grup {
			return b, true, nil
		}
	}
	return models.BarisMaster{}, false, nil
}

// BerkasPolis - berkas NB / EDM Treaty In.
func (a *Acuan) BerkasPolis(_ context.Context, nopolis string) (models.BerkasPolis, bool, error) {
	b, ada := a.Berkas[nopolis]
	return b, ada, nil
}

// DaftarAdjuster - pilihan adjuster.
func (a *Acuan) DaftarAdjuster(context.Context, string) ([]models.Pilihan, error) {
	out := []models.Pilihan{}
	for id, n := range a.Adjuster {
		out = append(out, models.Pilihan{Nilai: id, Label: n})
	}
	return out, nil
}

// DaftarMataUang - pilihan mata uang.
func (a *Acuan) DaftarMataUang(context.Context) ([]models.Pilihan, error) {
	out := []models.Pilihan{}
	for id, n := range a.NamaMU {
		out = append(out, models.Pilihan{Nilai: id, Label: n})
	}
	return out, nil
}

// DaftarProvinsi - pilihan provinsi.
func (a *Acuan) DaftarProvinsi(context.Context, string) ([]models.Pilihan, error) {
	return []models.Pilihan{{Nilai: "UJI-PROVINSI", Label: "UJI-PROVINSI"}}, nil
}

// DaftarKlien - pilihan nama klien rekening.
func (a *Acuan) DaftarKlien(context.Context, string) ([]models.Pilihan, error) {
	out := []models.Pilihan{}
	for _, r := range a.Rekening {
		out = append(out, models.Pilihan{Nilai: r.ClientName, Label: r.ClientName})
	}
	return out, nil
}

// DaftarSebab - pilihan cause of loss.
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

// DaftarKatastrofe - pilihan katastrofe.
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

// TanggaKomite - tangga kasus komite.
func (a *Acuan) TanggaKomite(_ context.Context, id string) ([]models.AnggotaKomite, error) {
	return a.Tangga[id], nil
}

// StatusKasir - keterangan DIRECTTOKASIR_LOG.
func (a *Acuan) StatusKasir(_ context.Context, no string) (string, bool, error) {
	v, ada := a.Kasir[no]
	return v, ada, nil
}

// StatusKonversi - status konversi akseptasi.
func (a *Acuan) StatusKonversi(_ context.Context, no string) (string, error) {
	return a.Konversi[no], nil
}

// EmailCeding - email ceding (produksi saja di repository).
func (a *Acuan) EmailCeding(context.Context, string) (string, error) { return "", nil }

// EmailPelaku - email akun.
func (a *Acuan) EmailPelaku(_ context.Context, akun string) (string, error) {
	return a.Email[akun], nil
}

// IDBankRekening - ID bank rekening.
func (a *Acuan) IDBankRekening(_ context.Context, bank, cabang, akun string) (string, error) {
	for _, r := range a.Rekening {
		if r.NameOfBank == bank && r.BranchOfBank == cabang && r.AccountNo == akun {
			return r.IDOfBank, nil
		}
	}
	return "", nil
}

// TingkatPelaku - tingkat wewenang.
func (a *Acuan) TingkatPelaku(_ context.Context, akun string) (string, error) {
	return a.Tingkat[akun], nil
}

// RiwayatMaster - CLAIMXOL2.
func (a *Acuan) RiwayatMaster(context.Context, [3]string) ([]models.BarisXOL2, error) {
	return a.XOL2, nil
}

// LampiranBayar - lampiran invoice.
func (a *Acuan) LampiranBayar(context.Context, string) ([]repository.BarisLampiranBayar, error) {
	return a.LampiranB, nil
}

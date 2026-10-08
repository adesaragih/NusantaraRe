package tiruan

// Untuk apa berkas ini: ACUAN TIRUAN - bacaan acuan (`services.Acuan`) dari peta dalam memori yang diisi uji. Semua
// nilai fixture berawalan `UJI-` atau kode buatan; nol data DEV.

import (
	"context"
	"sort"
	"strings"

	"nusantarare/modul/claimprop/backend/models"
	"nusantarare/modul/claimprop/backend/repository"
)

// AnggotaRoster - satu baris EMAILKOMITE tiruan.
type AnggotaRoster struct {
	models.AnggotaKomite
	Batas string // LIMIT_BOTTOM
	Sts   string // STS_KLAIM
}

// Rekening - satu baris BANKACCOUNT tiruan.
type Rekening struct {
	models.RekeningBank
	ClientID string
}

// Acuan - acuan tiruan.
type Acuan struct {
	Kurs           map[string]string
	NamaMU         map[string]string
	JenisReas      map[string]string // ID -> NOTE
	TipeReas       map[string]string // ID -> TYPE
	Master         map[string]models.MasterTreaty
	BarisMaster    []models.BarisMaster
	PolisRealisasi map[string]bool
	PolisMaster    map[string]bool
	Polis          []models.BarisPolis
	TreatyGroup    map[string]string // bizCode|tahun
	YearQuartal    map[string]string // nopolis|prodke
	WilayahPos     map[string][]models.BarisWilayah
	Alamat         map[string][4]string
	Adjuster       map[string]string
	Berkas         map[string]models.BerkasPolis // nopolis -> berkas NB / EDM Treaty In
	MO             map[string][6]string
	Riwayat        map[string][]models.RiwayatKlaimPolis
	OldID          map[string]string
	Tahun          map[string]string // grup|ymd -> tahun ("*" = semua tanggal)
	LimitPLAMap    map[string]string // tahun|grup|reins
	Retro          map[string][]models.Retro
	Roster         []AnggotaRoster
	LimitDirut     string
	Rekening       []Rekening
	Saldo          map[string]string // invoice|cur
	Proteksi       map[string]bool
	Sebab          []repository.BarisSebab
	Katastrofe     []repository.BarisKatastrofe
	OS             map[string][]repository.BarisRingkasanOS
	Tangga         map[string][]models.AnggotaKomite
	Kasir          map[string]string
	Konversi       map[string]string
	Tingkat        map[string]string
}

// AcuanBaru membuat acuan tiruan kosong.
func AcuanBaru() *Acuan {
	return &Acuan{Kurs: map[string]string{}, NamaMU: map[string]string{}, JenisReas: map[string]string{},
		TipeReas: map[string]string{}, Master: map[string]models.MasterTreaty{}, PolisRealisasi: map[string]bool{},
		PolisMaster: map[string]bool{}, TreatyGroup: map[string]string{}, YearQuartal: map[string]string{},
		WilayahPos: map[string][]models.BarisWilayah{}, Alamat: map[string][4]string{}, Adjuster: map[string]string{},
		Berkas: map[string]models.BerkasPolis{},
		MO:     map[string][6]string{}, Riwayat: map[string][]models.RiwayatKlaimPolis{}, OldID: map[string]string{},
		Tahun: map[string]string{}, LimitPLAMap: map[string]string{}, Retro: map[string][]models.Retro{},
		Saldo: map[string]string{}, Proteksi: map[string]bool{}, OS: map[string][]repository.BarisRingkasanOS{},
		Tangga: map[string][]models.AnggotaKomite{}, Kasir: map[string]string{}, Konversi: map[string]string{},
		Tingkat: map[string]string{}}
}

func (a *Acuan) KursStandar(_ context.Context, c string) (string, error)  { return a.Kurs[c], nil }
func (a *Acuan) NamaMataUang(_ context.Context, c string) (string, error) { return a.NamaMU[c], nil }
func (a *Acuan) NamaJenisReasuransi(_ context.Context, id string) (string, error) {
	return a.JenisReas[id], nil
}

func (a *Acuan) IDJenisReasuransi(_ context.Context, nama, tipe string) (string, error) {
	for id, n := range a.JenisReas {
		if strings.EqualFold(n, nama) && a.TipeReas[id] == tipe {
			return id, nil
		}
	}
	return "", nil
}

func (a *Acuan) MasterTreaty(_ context.Context, id string) (models.MasterTreaty, bool, error) {
	m, ada := a.Master[id]
	return m, ada, nil
}

func (a *Acuan) AdaPolisMaster(_ context.Context, id string) (bool, error) {
	return a.PolisMaster[id], nil
}
func (a *Acuan) AdaPolisRealisasi(_ context.Context, p string) (bool, error) {
	return a.PolisRealisasi[p], nil
}
func (a *Acuan) TreatyGroupBisnis(_ context.Context, b, t string) (string, error) {
	return a.TreatyGroup[b+"|"+t], nil
}
func (a *Acuan) YearOfQuartal(_ context.Context, p, k string) (string, error) {
	return a.YearQuartal[p+"|"+k], nil
}
func (a *Acuan) Wilayah(_ context.Context, z string) ([]models.BarisWilayah, error) {
	return a.WilayahPos[z], nil
}
func (a *Acuan) AlamatKlien(_ context.Context, id string) ([4]string, bool, error) {
	v, ada := a.Alamat[id]
	return v, ada, nil
}
func (a *Acuan) NamaAdjuster(_ context.Context, id string) (string, error) {
	return a.Adjuster[id], nil
}
func (a *Acuan) MarketingPolis(_ context.Context, p string) ([6]string, bool, error) {
	v, ada := a.MO[p]
	return v, ada, nil
}
func (a *Acuan) RiwayatKlaimPolis(_ context.Context, p string) ([]models.RiwayatKlaimPolis, error) {
	return a.Riwayat[p], nil
}
func (a *Acuan) KodeLamaBisnis(_ context.Context, b string) (string, error) { return a.OldID[b], nil }
func (a *Acuan) TahunTreaty(_ context.Context, g, ymd string) (string, error) {
	if v, ada := a.Tahun[g+"|"+ymd]; ada {
		return v, nil
	}
	return a.Tahun[g+"|*"], nil
}
func (a *Acuan) LimitPLA(_ context.Context, t, g, r string) (string, bool, error) {
	v, ada := a.LimitPLAMap[t+"|"+g+"|"+r]
	return v, ada, nil
}
func (a *Acuan) DaftarRetro(_ context.Context, r, t, g string) ([]models.Retro, error) {
	return a.Retro[r+"|"+t+"|"+g], nil
}

func (a *Acuan) RosterKomite(_ context.Context, nilai, sts string) ([]models.AnggotaKomite, error) {
	n, err := models.AngkaTeks("nilai", nilai)
	if err != nil {
		return nil, err
	}
	var out []models.AnggotaKomite
	for _, r := range a.Roster {
		b, _ := models.AngkaTeks("batas", r.Batas)
		if r.Sts == sts && !models.Lebih(b, n) {
			out = append(out, r.AnggotaKomite)
		}
	}
	return out, nil
}

func (a *Acuan) LimitDirekturUtama(context.Context) (string, bool, error) {
	return a.LimitDirut, a.LimitDirut != "", nil
}

func (a *Acuan) rek(f func(r Rekening) bool) []models.RekeningBank {
	out := []models.RekeningBank{}
	for _, r := range a.Rekening {
		if f(r) {
			out = append(out, r.RekeningBank)
		}
	}
	return out
}

func (a *Acuan) RekeningBank(_ context.Context, k, c string) ([]models.RekeningBank, error) {
	return a.rek(func(r Rekening) bool { return r.ClientID == k && r.CurrencyID == c }), nil
}
func (a *Acuan) RekeningBankMataUang(_ context.Context, c string) ([]models.RekeningBank, error) {
	return a.rek(func(r Rekening) bool { return r.CurrencyID == c }), nil
}
func (a *Acuan) RekeningBankKlien(_ context.Context, k string) ([]models.RekeningBank, error) {
	return a.rek(func(r Rekening) bool { return r.ClientID == k }), nil
}
func (a *Acuan) SaldoPremi(_ context.Context, i, c string) (string, error) {
	return a.Saldo[i+"|"+c], nil
}
func (a *Acuan) AdaProteksiPremi(_ context.Context, p string) (bool, error) {
	return a.Proteksi[p], nil
}

// ---------------------------------------------------------------- pemilih layar

func (a *Acuan) DaftarMaster(_ context.Context, s models.SaringanMaster) ([]models.BarisMaster, error) {
	memuat := func(nilai, kata string) bool {
		return strings.Contains(strings.ToUpper(nilai), strings.ToUpper(strings.TrimSpace(kata)))
	}
	out := []models.BarisMaster{}
	for _, b := range a.BarisMaster {
		if b.ProportionType == models.ProporsiMaster && memuat(b.TreatyID, s.TreatyID) &&
			memuat(b.ClassOfBusiness, s.ClassOfBusiness) && memuat(b.TreatyContractName, s.ContractName) &&
			memuat(b.SOB, s.SOB) && memuat(b.Ceding, s.InsuredName) && memuat(b.TreatyType, s.TreatyType) &&
			memuat(b.TreatyGroup, s.TreatyGroup) && memuat(b.TreatyYear, s.TreatyYear) {
			out = append(out, b)
		}
	}
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].TreatyYear != out[j].TreatyYear {
			return out[i].TreatyYear > out[j].TreatyYear
		}
		return out[i].TreatyID > out[j].TreatyID
	})
	if len(out) > models.BatasMaster {
		out = out[:models.BatasMaster]
	}
	return out, nil
}

func (a *Acuan) BerkasPolis(_ context.Context, nopolis string) (models.BerkasPolis, bool, error) {
	b, ada := a.Berkas[nopolis]
	return b, ada, nil
}

func (a *Acuan) BarisMasterDari(_ context.Context, id, grup, cob string) (models.BarisMaster, bool, error) {
	for _, b := range a.BarisMaster {
		if b.TreatyID == id && b.TreatyGroupID == grup && b.ClassOfBusinessID == cob {
			return b, true, nil
		}
	}
	return models.BarisMaster{}, false, nil
}

func (a *Acuan) DaftarPolis(_ context.Context, noOffer, grup string) ([]models.BarisPolis, error) {
	out := []models.BarisPolis{}
	for _, p := range a.Polis {
		if p.NoOffer == noOffer && p.TreatyGroup == grup {
			out = append(out, p)
		}
	}
	return out, nil
}

func (a *Acuan) DaftarAdjuster(_ context.Context, cari string) ([]models.Pilihan, error) {
	out := []models.Pilihan{}
	for id, n := range a.Adjuster {
		out = append(out, models.Pilihan{Nilai: id, Label: n})
	}
	return out, nil
}

func (a *Acuan) DaftarMataUang(context.Context) ([]models.Pilihan, error) {
	out := []models.Pilihan{}
	for id, n := range a.NamaMU {
		out = append(out, models.Pilihan{Nilai: id, Label: n})
	}
	return out, nil
}

func (a *Acuan) DaftarJenisReas(_ context.Context, tipe string) ([]models.Pilihan, error) {
	out := []models.Pilihan{}
	for id, n := range a.JenisReas {
		if tipe == "" || a.TipeReas[id] == tipe {
			out = append(out, models.Pilihan{Nilai: id, Label: n})
		}
	}
	return out, nil
}

func (a *Acuan) DaftarProvinsi(context.Context, string) ([]models.Pilihan, error) {
	return []models.Pilihan{{Nilai: "UJI-PROVINSI", Label: "UJI-PROVINSI"}}, nil
}

func (a *Acuan) DaftarSebab(context.Context, string) ([]repository.BarisSebab, error) {
	return a.Sebab, nil
}

func (a *Acuan) BarisSebabID(_ context.Context, id string) (repository.BarisSebab, bool, error) {
	for _, s := range a.Sebab {
		if s.ID == id {
			return s, true, nil
		}
	}
	return repository.BarisSebab{}, false, nil
}

func (a *Acuan) DaftarKatastrofe(context.Context, string) ([]repository.BarisKatastrofe, error) {
	return a.Katastrofe, nil
}

func (a *Acuan) BarisKatastrofeID(_ context.Context, id string) (repository.BarisKatastrofe, bool, error) {
	for _, s := range a.Katastrofe {
		if s.ID == id {
			return s, true, nil
		}
	}
	return repository.BarisKatastrofe{}, false, nil
}

func (a *Acuan) RingkasanOS(_ context.Context, no string) ([]repository.BarisRingkasanOS, error) {
	return a.OS[no], nil
}

func (a *Acuan) TanggaKomite(_ context.Context, id string) ([]models.AnggotaKomite, error) {
	return a.Tangga[id], nil
}

func (a *Acuan) StatusKasir(_ context.Context, no string) (string, bool, error) {
	v, ada := a.Kasir[no]
	return v, ada, nil
}

func (a *Acuan) StatusKonversi(_ context.Context, no string) (string, error) {
	return a.Konversi[no], nil
}
func (a *Acuan) EmailCeding(context.Context, string) (string, error) { return "", nil }
func (a *Acuan) IDBankRekening(context.Context, string, string, string) (string, error) {
	return "UJI-BANK", nil
}
func (a *Acuan) TingkatPelaku(_ context.Context, op string) (string, error) {
	return a.Tingkat[op], nil
}
func (a *Acuan) NamaPelaku(_ context.Context, akun string) (string, error) { return akun, nil }

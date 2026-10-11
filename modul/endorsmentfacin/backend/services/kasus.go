package services

// Buat kasus Endorsment Fac In - irisan I3 (kontrak sesi utama 11-10-2026 + ralat bentuk jawaban; keputusan work
// owner W-EDM-08…12).
//
// Urutan: akun → isian wajib → pemeriksaan /periksa diulang di server → Detail EDM Type → OpenCase (bukakasus.go,
// kering, tanpa nomor) → SATU transaksi: nomor SEQ_WORK_POLIS_EDM, OpenCase dengan identitas nyata, T_WORK_POLIS,
// T_EDM_FACIN_KASUS. Penolakan = HASIL BISNIS (HasilKasus.Periksa), bukan galat; nol tulisan.
//
// Dibaca sesudah: layanan.go, bukakasus.go.
//
// Tidak di sini: before-image (Seam 4, langkah 14-15 SetValueToEDMWork), penyimpanan OfferFacIn, tautan
// EDMHandle/EDMHandle2 ke kasus portal (portal bukan entitas tersimpan di sistem baru - A43).

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/endorsmentfacin/backend/models"
	"nusantarare/modul/endorsmentfacin/backend/repository"
)

// JenisPengecualianBatal - `M_EDM_PENGECUALIAN_POLIS.JENIS` pengganti literal `SetValueToEDMWork` langkah 5.
const JenisPengecualianBatal = "BATAL"

var (
	// ErrDilarang - pelaku tanpa akun; pesannya untuk pengguna.
	ErrDilarang = errors.New("services: tidak berhak")
	// ErrBelumAda - tabel / sequence modul ini belum dibuat DBA (migrasi 260 / 262).
	ErrBelumAda = repository.ErrBelumAda
	// ErrTanpaTransaksi - Layanan dirakit tanpa DenganTransaksi: kesalahan perakitan.
	ErrTanpaTransaksi = errors.New("services: layanan tanpa pembuka transaksi")
)

// Batas lebar kolom T_EDM_FACIN_KASUS (migrasi 262), dalam BYTE - diperiksa sebelum menulis (bukan ORA-12899).
var batasKolomKasus = []struct {
	label string
	nilai func(models.MasukanKasus) string
	maks  int
}{
	{"Policy Number", func(m models.MasukanKasus) string { return m.PolicyNo }, 100},
	{"Endorsement Source", func(m models.MasukanKasus) string { return m.EndorsementSourceNote }, 1000},
	{"Endorsement Source Type", func(m models.MasukanKasus) string { return m.EndorsementSource }, 255},
	{"Endorsement Note", func(m models.MasukanKasus) string { return m.Note }, 4000},
	{"Endorsement Status", func(m models.MasukanKasus) string { return m.EndorsementStatus }, 255},
	{"EDM Type", func(m models.MasukanKasus) string { return m.EdmTypeNew }, 10},
	{"Detail EDM Type", func(m models.MasukanKasus) string { return m.DetailType }, 10},
	{"Endorsement Internal Retro", func(m models.MasukanKasus) string { return m.EndorsementInternalRetro }, 10},
}

// ValidasiKasus - isian wajib `pyRequired` Section/WorkPrimaryDetails.xml `[terverifikasi]`, urut layar; galat = medan
// pertama yang gagal. Detail EDM Type diperiksa sesudah pemeriksaan (validasiDetailType), karena syaratnya membaca
// jenis bisnis polis lama.
func ValidasiKasus(m models.MasukanKasus) error {
	kosong := func(s string) bool { return strings.TrimSpace(s) == "" }
	switch {
	case kosong(m.PolicyNo):
		return tolak("Policy Number is required")
	case kosong(m.EndorsementSourceNote):
		return tolak("Endorsement Source is required")
	case kosong(m.Note):
		return tolak("Endorsement Note is required")
	case kosong(m.EdmType):
		return tolak("Endorsement Type is required")
	case m.EdmType == string(models.EdmPerubahan) && kosong(m.EdmTypeNew):
		// `pyRequiredWhen .Quotation.EdmType = 4`
		return tolak("EDM Type is required")
	case kosong(m.EndorsementDate):
		return tolak("Endorsement Date is required")
	case kosong(m.EndorsementStatus):
		// A46 dijawab work owner 11-10-2026: wajib; Prompt List 0-3.
		return tolak("Endorsement Status is required")
	case kosong(m.EndorsementInternalRetro):
		return tolak("Endorsement Internal Retro is required")
	}
	if _, ada := models.PilihanEndorsementStatus[m.EndorsementStatus]; !ada {
		return tolak("Endorsement Status %q is not an option", m.EndorsementStatus)
	}
	if _, ada := models.PilihanInternalRetro[m.EndorsementInternalRetro]; !ada {
		// A47 dijawab work owner 11-10-2026: 1 Yes / 0 No.
		return tolak("Endorsement Internal Retro %q is not an option", m.EndorsementInternalRetro)
	}
	for _, b := range batasKolomKasus {
		if n := len(b.nilai(m)); n > b.maks {
			return tolak("%s is too long (%d bytes, maximum %d)", b.label, n, b.maks)
		}
	}
	return nil
}

// validasiDetailType - Section EdmType1/2/3 `pyRequiredWhen .Quotation.EdmType==4 && OutData.pxResults(1).CARI2!='Life'
// && .Quotation.EdmTypeNew==n`; pilihannya daftar SetEdmType (models.DetailTypeMenurutEdmTypeNew). Di luar itu Detail
// EDM Type tidak ditampilkan layar, jadi wajib kosong (keputusan agent A45).
func validasiDetailType(m models.MasukanKasus, jenisBisnis string) error {
	pilihan, berlaku := models.DetailTypeMenurutEdmTypeNew[m.EdmTypeNew]
	berlaku = berlaku && m.EdmType == string(models.EdmPerubahan) && jenisBisnis != "Life"
	if !berlaku {
		if strings.TrimSpace(m.DetailType) != "" {
			return tolak("Detail EDM Type applies only to Endorsement Type 4 with EDM Type 1, 2 or 3")
		}
		return nil
	}
	if strings.TrimSpace(m.DetailType) == "" {
		return tolak("Detail EDM Type is required")
	}
	for _, p := range pilihan {
		if string(p) == m.DetailType {
			return nil
		}
	}
	return tolak("Detail EDM Type %q is not an option for EDM Type %s", m.DetailType, m.EdmTypeNew)
}

// teksCentang - checkbox Pega "true"/"false".
func teksCentang(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

// Buat - `POST /kasus`. akun = pengenal akun pelaku (kosong → ErrDilarang).
func (l *Layanan) Buat(ctx context.Context, akun string, m models.MasukanKasus) (models.HasilKasus, error) {
	if strings.TrimSpace(akun) == "" {
		return models.HasilKasus{}, fmt.Errorf("%w: %s", ErrDilarang, "Log in to create an endorsement")
	}
	if l.tx == nil {
		return models.HasilKasus{}, ErrTanpaTransaksi
	}
	if err := ValidasiKasus(m); err != nil {
		return models.HasilKasus{}, err
	}
	d, err := UraiMasukanPeriksa(m.MasukanPeriksa)
	if err != nil {
		return models.HasilKasus{}, err
	}
	if err := l.bacaData(ctx, &d); err != nil {
		return models.HasilKasus{}, err
	}
	h, err := SusunHasilPeriksa(d)
	if err != nil {
		return models.HasilKasus{}, err
	}
	if h.Ditolak || !h.Lengkap {
		return models.HasilKasus{Pesan: []models.PesanPeriksa{}, Periksa: &h}, nil
	}
	if err := validasiDetailType(m, h.Bisnis.JenisBisnis); err != nil {
		return models.HasilKasus{}, err
	}

	// OpenCase - SetValueToEDMWork langkah 2 membaca GetEDMStatus_SQL juga untuk RI slip.
	if !d.StatusEDMDibaca {
		if d.StatusEDM, err = l.gudang.StatusEDM(ctx, d.PolicyNo); err != nil {
			return models.HasilKasus{}, err
		}
	}
	batal, err := l.gudang.Pengecualian(ctx, JenisPengecualianBatal)
	if err != nil && !errors.Is(err, repository.ErrBelumAda) { // tabel 261 belum ada = daftar kosong
		return models.HasilKasus{}, err
	}
	masuk := MasukanBukaKasus{
		Portal: KasusPortal{
			PolicyNo: m.PolicyNo, Note: m.Note, EndorsementSource: m.EndorsementSource,
			EndorsementSourceNote: m.EndorsementSourceNote, EndorsementDate: *d.TanggalEndorsement,
			EndorsementStatus: m.EndorsementStatus, AdminCharge: teksCentang(m.AdminCharge), Survey: teksCentang(m.Survey),
			EndorsementInternalRetro: m.EndorsementInternalRetro,
			Quotation: QuotationPortal{EdmType: d.EdmType, EdmTypeNew: m.EdmTypeNew,
				Type: models.JenisPenyesuaian(m.DetailType)},
		},
		StatusEDMPolis: d.StatusEDM, PengecualianBatal: batal,
		// Identitas sementara: OpenCase murni; jalan kering ini hanya untuk penolakan, sebelum nomor diambil.
		KasusBaru: IdentitasKasus{PxInsName: models.AwalanKasusEDM, PzInsKey: models.AwalanKasusEDM},
	}
	kering, err := OpenCase(masuk)
	if err != nil {
		return models.HasilKasus{}, err
	}
	if kering.Ditolak {
		h.Pesan = append(h.Pesan, models.PesanPeriksa{Medan: medanPolicyNo, Teks: kering.AlasanTolak})
		h.Ditolak = true
		return models.HasilKasus{Pesan: []models.PesanPeriksa{}, Periksa: &h}, nil
	}

	var id string
	err = l.tx(ctx, func(tx *db.Tx) error {
		var err error
		if id, err = l.gudang.NomorKasusBerikut(ctx, tx); err != nil {
			return err
		}
		// pzInsKey Pega tidak dibentuk (bentuknya tidak dibutuhkan dan tidak disimpan irisan ini - A43): = ID.
		masuk.KasusBaru = IdentitasKasus{PxInsName: id, PzInsKey: id}
		hb, err := OpenCase(masuk)
		if err != nil {
			return err
		}
		if err := l.gudang.SisipWorkPolis(ctx, tx, id, akun); err != nil {
			return err
		}
		return l.gudang.SisipKasus(ctx, tx, KepalaDari(hb.Kasus))
	})
	if err != nil {
		return models.HasilKasus{}, err
	}
	hasil := models.HasilKasus{CaseID: id, Pesan: []models.PesanPeriksa{}}
	for _, p := range kering.Pesan {
		hasil.Pesan = append(hasil.Pesan, models.PesanPeriksa{Medan: p.Medan, Teks: p.Teks})
	}
	return hasil, nil
}

// KepalaDari - baris T_EDM_FACIN_KASUS dari kasus keluaran OpenCase (`DataToEDM` + langkah 12-13 + konektor Start2).
func KepalaDari(k models.KasusEndorsement) models.KepalaKasus {
	o := k.OfferFacIn
	q := o.QuotationData
	return models.KepalaKasus{
		ID: k.PxInsName, OldPolicyNo: q.OldPolicyNo, EdmType: string(q.EdmType), EdmTypeNew: q.EdmTypeNew,
		DetailType: string(q.Type), EndorsementDate: q.EdmDate, EndorsementReason: o.EndorsmentReason,
		EndorsementSource: q.EdmSource, EndorsementSourceNote: q.EdmSourceNote, EndorsementStatus: q.EdmStatus,
		AdminCharge: q.EdmChargeFee, Survey: q.EdmSurvey, InternalRetro: q.EndorsementInternalRetro,
		NBStatus: k.NBStatus, IsCedingConfirm: k.IsCedingConfirm,
	}
}

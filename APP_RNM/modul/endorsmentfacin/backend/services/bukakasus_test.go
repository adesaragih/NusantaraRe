package services_test

// Seam 5 - alur masuk: kasus endorsement lahir dari kasus portal. Tiket E03.
//
// Dibaca sesudah: services/bukakasus.go. Nomor polis di berkas ini SINTETIS.

import (
	"testing"

	"nusantarare/modul/endorsmentfacin/backend/services"
)

func portalUji(t *testing.T) services.KasusPortal {
	t.Helper()
	return services.KasusPortal{
		PzInsKey:                 "UJI-PORTAL-KEY",
		PolicyNo:                 "UJI-POLIS-1",
		Note:                     "UJI-ALASAN",
		EndorsementSource:        "UJI-SRC",
		EndorsementSourceNote:    "UJI-SRCNOTE",
		EndorsementDate:          tanggal(t, "2026-01-02 00:00"),
		EndorsementStatus:        "UJI-ST",
		AdminCharge:              "UJI-FEE",
		Survey:                   "UJI-SURVEY",
		EndorsementInternalRetro: "0",
		Quotation:                services.QuotationPortal{EdmType: "4", EdmTypeNew: "1", Type: "1"},
	}
}

func bukaUji(t *testing.T, statusEDM string, pengecualian ...string) services.HasilBukaKasus {
	t.Helper()
	h, err := services.OpenCase(services.MasukanBukaKasus{
		Portal:            portalUji(t),
		StatusEDMPolis:    statusEDM,
		PengecualianBatal: services.DaftarPolis(pengecualian),
		KasusBaru:         services.IdentitasKasus{PxInsName: "UJI-EDM-1", PzInsKey: "UJI-EDM-KEY"},
	})
	if err != nil {
		t.Fatalf("OpenCase: %v", err)
	}
	return h
}

// TestTracerOpenCase - E03: kasus lahir berkelas endorsement, berprefiks EDM-,
// bertautan tiga arah, langsung berfase Policy di workbasket marketing.
func TestTracerOpenCase(t *testing.T) {
	h := bukaUji(t, "")
	if h.Ditolak {
		t.Fatalf("ditolak: %q", h.AlasanTolak)
	}
	k := h.Kasus
	if k.Kelas != services.KelasKasusEndorsement || k.PrefiksID != services.PrefiksIDEndorsement {
		t.Errorf("kelas %q prefiks %q", k.Kelas, k.PrefiksID)
	}
	// Tautan: EDM → portal, portal → EDM, label.
	if k.EndorsementID != "UJI-PORTAL-KEY" || h.PortalEDMHandle != "UJI-EDM-KEY" || k.PyLabel != "UJI-EDM-1" {
		t.Errorf("tautan: EndorsementID %q EDMHandle %q pyLabel %q", k.EndorsementID, h.PortalEDMHandle, k.PyLabel)
	}
	q := k.OfferFacIn.QuotationData
	if q.OldPolicyNo != "UJI-POLIS-1" || k.Policy.PolicyNo != "UJI-POLIS-1" || k.PolicyNumber != "UJI-POLIS-1" {
		t.Errorf("nomor polis lama tidak terisi: %+v", q)
	}
	if k.OfferFacIn.EndorsmentReason != "UJI-ALASAN" || q.EdmNote != "UJI-ALASAN" {
		t.Errorf("alasan endorsement: %q / %q", k.OfferFacIn.EndorsmentReason, q.EdmNote)
	}
	// DataToEDM.
	if q.EdmSource != "UJI-SRC" || q.EdmSourceNote != "UJI-SRCNOTE" || q.EdmType != "4" || q.EdmTypeNew != "1" ||
		q.Type != "1" || !q.EdmDate.Equal(tanggal(t, "2026-01-02 00:00")) || q.EdmStatus != "UJI-ST" ||
		q.EdmChargeFee != "UJI-FEE" || q.EdmSurvey != "UJI-SURVEY" || q.EndorsementInternalRetro != "0" {
		t.Errorf("DataToEDM: %+v", q)
	}
	if q.StatusBusiness != "3" || k.OfferFacIn.IsBanding != "false" {
		t.Errorf("StatusBusiness %q IsBanding %q", q.StatusBusiness, k.OfferFacIn.IsBanding)
	}
	// Konektor Start2 → Assignment7: fase Policy, tanpa penawaran/binding.
	if k.IsCedingConfirm != services.FasePolicy || k.FlagOnGoingPolicy != "1" || k.Position != "1" ||
		k.PositionNote != "ReasFacInMarketing" || k.NBStatus != "NEW EDM" || k.NBStatusNew != "NEW EDM" {
		t.Errorf("konektor masuk flow: %+v", k)
	}
	if h.Assignment.Workbasket != "ReasFacInMarketing" || h.Assignment.Tiket != "AdminPolicy" {
		t.Errorf("assignment: %+v", h.Assignment)
	}
}

// TestOpenCaseTidakPernahFasePenawaran - kasus uji negatif E03.
func TestOpenCaseTidakPernahFasePenawaran(t *testing.T) {
	k := bukaUji(t, "").Kasus
	if k.IsCedingConfirm == "Offer" || k.IsCedingConfirm == "Binding" || k.IsCedingConfirm == "" {
		t.Errorf("fase %q", k.IsCedingConfirm)
	}
}

// TestOpenCaseSudahBatalDitolak - langkah 4-6: polis yang sudah dibatalkan
// lewat endorsement ditolak (penolakan bisnis, bukan galat sistem), kecuali
// nomornya ada di daftar pengecualian konfigurasi.
func TestOpenCaseSudahBatalDitolak(t *testing.T) {
	h := bukaUji(t, "1")
	if !h.Ditolak || h.AlasanTolak != "Sudah Di endorsement Batal" {
		t.Errorf("ditolak %v alasan %q", h.Ditolak, h.AlasanTolak)
	}
	if h.Kasus.Kelas != "" {
		t.Errorf("kasus lahir padahal ditolak: %+v", h.Kasus)
	}

	// Pengecualian: kasus tetap lahir; pesan langkah 4 tetap terpasang.
	h = bukaUji(t, "1", "UJI-POLIS-1")
	if h.Ditolak || h.Kasus.Kelas == "" {
		t.Errorf("polis pengecualian ditolak: %+v", h)
	}
	if len(h.Pesan) != 1 {
		t.Errorf("pesan langkah 4: %v", h.Pesan)
	}

	// Status lain tidak memblokir.
	if h := bukaUji(t, "2"); h.Ditolak || len(h.Pesan) != 0 {
		t.Errorf("status 2: %+v", h)
	}
}

// TestOpenCaseBerlanjutKeBeforeImage - keluaran Seam 5 adalah masukan Seam 4:
// tanggal, jenis, dan penanda siklus sudah di tempat yang dibaca before-image.
func TestOpenCaseBerlanjutKeBeforeImage(t *testing.T) {
	h := bukaUji(t, "")
	lama := polisLamaFire(t)
	b, err := services.PrepareBeforeImage(services.MasukanBeforeImage{
		Kasus: h.Kasus, PolisLama: &lama, Predikat: services.Predikat{IsFire: true, IsEDM: true},
	})
	if err != nil || b.GalatPorsiPeriode != nil {
		t.Fatalf("%v / %v", err, b.GalatPorsiPeriode)
	}
	if b.Kasus.OfferFacIn.ProrateEDMEnd.Kosong() || b.Kasus.Kelas != services.KelasKasusEndorsement {
		t.Errorf("before-image atas kasus OpenCase: %+v", b.Kasus.OfferFacIn.ProrateEDMEnd)
	}
}

// TestOpenCaseIdentitasWajib - tanpa identitas kasus baru ini galat SISTEM,
// bukan penolakan bisnis.
func TestOpenCaseIdentitasWajib(t *testing.T) {
	_, err := services.OpenCase(services.MasukanBukaKasus{Portal: portalUji(t)})
	if err == nil {
		t.Fatal("tanpa galat")
	}
}

// TestTautanAssignment - langkah 22-25: tautan ketiga ke handle assignment.
func TestTautanAssignment(t *testing.T) {
	workbasket := []services.Assignment{
		{RefObjectKey: "UJI-LAIN", InsKey: "UJI-A0"},
		{RefObjectKey: "UJI-EDM-KEY", InsKey: "UJI-A1"},
	}
	if got := services.TautanAssignment("UJI-EDM-KEY", workbasket, nil); got != "UJI-A1" {
		t.Errorf("EDMHandle2 %q", got)
	}
	// Langkah 25 (worklist) menulis sesudah langkah 24 (workbasket).
	worklist := []services.Assignment{{RefObjectKey: "UJI-EDM-KEY", InsKey: "UJI-W1"}}
	if got := services.TautanAssignment("UJI-EDM-KEY", workbasket, worklist); got != "UJI-W1" {
		t.Errorf("EDMHandle2 %q, mau worklist", got)
	}
	if got := services.TautanAssignment("UJI-EDM-KEY", nil, nil); got != "" {
		t.Errorf("tanpa assignment: %q", got)
	}
}

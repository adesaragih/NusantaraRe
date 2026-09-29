package services

// Ambang batas hari sebuah klaim - butir ba, dibuka keputusan bh.
//
// Untuk apa berkas ini: menyambung tiga hal yang sampai hari ini berdiri
// sendiri-sendiri -
//
//	aturannya   `models.PenandaBatasHari`  (murni, sudah ada sejak kelompok Detail)
//	kuncinya    `PolicyDataLife.ProductNameID` (butir av, 28-09-2026)
//	ambangnya   view produk                 (butir bh, 28-09-2026)
//
// ⛔ ATURANNYA TIDAK PERNAH PUNYA PEMANGGIL sampai berkas ini. Ia ditulis
// murni dan diuji murni, lalu MENUNGGU - dan fungsi yang menunggu lama adalah
// fungsi yang mudah dianggap sudah jalan. Itu sebabnya penyambungannya
// dicatat, bukan sekadar dikerjakan.
//
// ⛔ AMBANGNYA DIBACA SERVER, bukan diterima dari klien. Ia menentukan sah
// atau tidaknya sebuah klaim; ambang yang dapat disebut pemanggil adalah
// ambang yang dapat dipilih pemanggil.
//
// Dibaca sesudah: models/validasitanggal.go (aturannya),
// repository/ambangproduk.go (pembacanya).

import (
	"context"
	"errors"
	"fmt"

	"nusantarare/internal/models"
	"nusantarare/internal/repository"
	"nusantarare/inti"
	"nusantarare/inti/db"
	intigalat "nusantarare/inti/galat"
)

// ErrAmbangProdukTakDitemukan dirujuk ulang supaya handler tidak perlu
// mengimpor repository.
var ErrAmbangProdukTakDitemukan = repository.ErrAmbangProdukTakDitemukan

// PenandaBatasKlaim adalah kedua penanda beserta ambang yang melahirkannya.
//
// ⚠️ AMBANGNYA IKUT DIKIRIM. Penanda tanpa ambangnya menyuruh orang percaya
// begitu saja; dengan ambangnya, ia dapat memeriksa sendiri kenapa klaimnya
// ditandai - dan menemukan produk yang ambangnya keliru.
type PenandaBatasKlaim struct {
	ProductNameID string `json:"productNameId"`
	ProductName   string `json:"productName"`
	// MaxExpiredClaim dan MaxDataReceive apa adanya dari produk.
	MaxExpiredClaim string `json:"maxExpiredClaim"`
	MaxDataReceive  string `json:"maxDataReceive"`
	// PenandaClaimReceived kosong berarti SAH; terisi berarti tanggal yang
	// melanggar - `dd/MM/yyyy`, persis yang Pega tulis ke medannya.
	PenandaClaimReceived string `json:"penandaClaimReceived"`
	// PenandaSTNC sama, untuk EFFECTIVE_DATE -> RECEIVED_DATE.
	PenandaSTNC string `json:"penandaStnc"`
}

// TanggalKlaimUntukAmbang adalah keempat tanggal yang kedua aturan pakai.
//
// ⚠️ Satu tipe, bukan empat parameter teks yang berjajar. Keempatnya
// `dd/MM/yyyy` dan tertukar di antaranya menghasilkan penanda yang tetap
// terlihat masuk akal - itu Data Clump yang memang perlu dibungkus.
type TanggalKlaimUntukAmbang struct {
	// DateOfLoss -> ClaimReceived, berambang MaxExpiredClaim.
	DateOfLoss    string
	ClaimReceived string
	// EffectiveDate -> Received, berambang MaxDataReceive.
	EffectiveDate string
	Received      string
}

// AmbangKlaim melayani perhitungan kedua penanda.
type AmbangKlaim struct{ svc *Service }

// AmbangKlaim menyusun layanannya.
func (s *Service) AmbangKlaim() *AmbangKlaim { return &AmbangKlaim{svc: s} }

// Hitung membaca ambang produk sebuah polis lalu menghitung kedua penanda.
//
// ⛔ PRODUKNYA DARI POLIS, bukan dari permintaan. Nomor polis -> versi
// berjalan -> `PRODUCT_NAME_ID` -> ambang. Rantai itu utuh di server.
func (a *AmbangKlaim) Hitung(ctx context.Context, pelaku inti.Pelaku,
	nomorPolis string, t TanggalKlaimUntukAmbang) (PenandaBatasKlaim, error) {

	if err := inti.WajibIdentitas(pelaku); err != nil {
		return PenandaBatasKlaim{}, err
	}
	if a == nil || a.svc == nil || !a.svc.PunyaDatabase() {
		return PenandaBatasKlaim{}, db.ErrTanpaOracle
	}
	if nomorPolis == "" {
		return PenandaBatasKlaim{}, fmt.Errorf("%w: nomor polis kosong", intigalat.ErrPermintaanTidakSah)
	}

	polis, err := a.svc.PembacaPolis().Ringkas(ctx, nomorPolis)
	if err != nil {
		return PenandaBatasKlaim{}, err
	}
	// ⛔ Polis TANPA produk gagal terang. Membiarkannya lewat berarti
	// menghitung ambang dari produk yang tidak disebut siapa pun.
	if polis.ProductNameID == "" {
		return PenandaBatasKlaim{}, fmt.Errorf(
			"%w: polis %q tidak menyebut produk", ErrAmbangProdukTakDitemukan, nomorPolis)
	}

	ambang, err := repository.NewProdukLife(a.svc.DB()).Ambang(ctx, polis.ProductNameID)
	if err != nil {
		return PenandaBatasKlaim{}, err
	}

	hasil := PenandaBatasKlaim{
		ProductNameID:   polis.ProductNameID,
		ProductName:     polis.ProductName,
		MaxExpiredClaim: ambang.MaxExpiredClaim,
		MaxDataReceive:  ambang.MaxDataReceive,
	}
	// ⛔ KEDUANYA DIHITUNG, walau yang pertama gagal. Menghentikan
	// perhitungan di penanda pertama membuat orang memperbaiki satu tanggal,
	// menyimpan, lalu ditolak lagi oleh yang kedua.
	var galat []error
	p1, err1 := models.PenandaBatasHari(t.DateOfLoss, t.ClaimReceived, ambang.MaxExpiredClaim)
	if err1 != nil {
		galat = append(galat, fmt.Errorf("penanda claim received: %w", err1))
	} else {
		hasil.PenandaClaimReceived = p1
	}
	p2, err2 := models.PenandaBatasHari(t.EffectiveDate, t.Received, ambang.MaxDataReceive)
	if err2 != nil {
		galat = append(galat, fmt.Errorf("penanda STNC: %w", err2))
	} else {
		hasil.PenandaSTNC = p2
	}
	if len(galat) > 0 {
		return hasil, fmt.Errorf("%w: %w", intigalat.ErrPermintaanTidakSah, errors.Join(galat...))
	}
	return hasil, nil
}

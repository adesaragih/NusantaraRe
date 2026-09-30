package services

// Penolak kontrak Claim Life yang belum disambung - refactor bentuk B.
//
// ⚠️ Metode di bawah bernama sama dengan penulis transisi Claim Life
// (`PerbaruiStatusBaris`, `CerminkanHeader`, ...), tetapi tidak menulis apa
// pun: setiap panggilan dijawab `ErrKlaimBelumDisambung`. Penjaga jejak dan
// penjaga penulis status mengenal berkas ini dengan namanya.

import (
	"context"
	"time"

	"nusantarare/inti/backend/db"
)

// klaimBelumDisambung menolak setiap panggilan dengan ErrKlaimBelumDisambung.
type klaimBelumDisambung struct{}

func (klaimBelumDisambung) TypeKlaim(context.Context, string) (string, error) {
	return "", ErrKlaimBelumDisambung
}

func (klaimBelumDisambung) KodeBisnisKlaim(context.Context, string) (string, error) {
	return "", ErrKlaimBelumDisambung
}

func (klaimBelumDisambung) PerbaruiStatusBaris(context.Context, *db.Tx, string, string, string,
	string, string, time.Time) error {
	return ErrKlaimBelumDisambung
}

func (klaimBelumDisambung) CerminkanHeader(context.Context, *db.Tx, string, string, string) error {
	return ErrKlaimBelumDisambung
}

func (klaimBelumDisambung) CabutPenandaDipilih(context.Context, *db.Tx, string) error {
	return ErrKlaimBelumDisambung
}

func (klaimBelumDisambung) NomorAkseptasiDipakai(context.Context, *db.Tx, string) (bool, error) {
	return false, ErrKlaimBelumDisambung
}

func (klaimBelumDisambung) PastikanKasusTerbuka(context.Context, string) error {
	return ErrKlaimBelumDisambung
}

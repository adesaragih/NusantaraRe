package repository

// T_WORK_POLIS.POSITION_NOTE untuk kasus EDM (keputusan work owner W-EDM-18, mengganti W-EDM-11; kolom dari migrasi
// nbfacin 254, VARCHAR2(255)). POSITION = tahap ('Policy'), POSITION_NOTE = workbasket ('ReasFacInMarketing') - pola NB
// tiket 77 (`modul/nbfacin/backend/repository/posisinote.go`, DISALIN polanya, tidak diimpor).
//
// Peralihan TOLERAN ("jangan ubah yang sudah berhasil"): selama 254 belum dijalankan, kolom tidak ada → SQL PERSIS seperti
// sebelum W-EDM-18 (POSITION = 'ReasFacInMarketing', tanpa POSITION_NOTE). Deteksi sekali per tabel; ketiadaan
// (ORA-00904) di-cache dan dicoba ulang tiap jedaUlangCatatan. Kasus yang lahir sebelum 254 dikonversi migrasi 263.
// Akibatnya biner aman dideploy sebelum maupun sesudah 254.

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"nusantarare/modul/endorsmentfacin/backend/models"
)

// KolomCatatanPosisi - kolom workbasket T_WORK_POLIS (migrasi nbfacin 254).
const KolomCatatanPosisi = "POSITION_NOTE"

// jedaUlangCatatan - ketiadaan kolom dicoba ulang sesudah jeda ini (254 dapat dijalankan selagi backend hidup).
const jedaUlangCatatan = time.Minute

var catatanPosisi = struct {
	sync.Mutex
	ada  map[string]bool
	coba map[string]time.Time
}{ada: map[string]bool{}, coba: map[string]time.Time{}}

func sqlUjiCatatanPosisi(work string) string {
	return "SELECT " + KolomCatatanPosisi + " FROM " + work + " WHERE ROWNUM < 1"
}

// kolomBelumAda - ORA-00904 (kolom tak dikenal) = migrasi belum dijalankan.
func kolomBelumAda(err error) bool { return err != nil && strings.Contains(err.Error(), "ORA-00904") }

// adaCatatanPosisi - true bila T_WORK_POLIS.POSITION_NOTE ada. Galat lain (koneksi) → false + log: bentuk lama tetap sah
// (kolom nullable) dan INSERT utama yang akan melaporkan galat sesungguhnya.
func (g *Gudang) adaCatatanPosisi(ctx context.Context, work string) bool {
	catatanPosisi.Lock()
	ada, tahu := catatanPosisi.ada[work]
	sejak := catatanPosisi.coba[work]
	catatanPosisi.Unlock()
	if ada || (tahu && time.Since(sejak) < jedaUlangCatatan) {
		return ada
	}
	baris, err := g.db.QueryContext(ctx, sqlUjiCatatanPosisi(work))
	if err == nil {
		_ = baris.Close()
	}
	ada = err == nil
	switch {
	case kolomBelumAda(err):
		log.Printf("endorsmentfacin: %s.%s belum ada (migrasi nbfacin 254); kasus EDM ditulis bentuk lama (POSITION = %s)",
			work, KolomCatatanPosisi, models.CatatanPosisiAwalEDM)
	case err != nil:
		log.Printf("endorsmentfacin: uji kolom %s.%s: %v", work, KolomCatatanPosisi, err)
	}
	catatanPosisi.Lock()
	catatanPosisi.ada[work], catatanPosisi.coba[work] = ada, time.Now()
	catatanPosisi.Unlock()
	return ada
}

// SqlSisipWorkPolisCatatan - bentuk W-EDM-18: POSITION = tahap, POSITION_NOTE = workbasket.
func SqlSisipWorkPolisCatatan(t string) string {
	return fmt.Sprintf(`INSERT INTO %s (ID, LINI, POSITION, POSITION_NOTE, FLAG_ONGOING_POLICY, STATUS_WORK, CREATE_OP,
	  CREATE_OP_NAME, TGL_CREATE, TGL_UPDATE) VALUES (:1, :2, :3, :4, :5, :6, :7, :8, SYSDATE, SYSDATE)`, t)
}

// SisipWorkPolisMenurut - teks SQL dan nilai bind baris T_WORK_POLIS kasus EDM: dengan kolom POSITION_NOTE bentuk
// W-EDM-18, tanpa kolom bentuk lama W-EDM-11 (PERSIS seperti sebelumnya, dikunci uji).
func SisipWorkPolisMenurut(t string, adaCatatan bool, id, akun string) (string, []any) {
	if adaCatatan {
		return SqlSisipWorkPolisCatatan(t), []any{id, models.LiniEDMFacIn, models.PosisiAwalEDM,
			models.CatatanPosisiAwalEDM, models.FlagAwalEDM, models.StatusAwalEDM, akun, akun}
	}
	return SqlSisipWorkPolis(t), []any{id, models.LiniEDMFacIn, models.CatatanPosisiAwalEDM, models.FlagAwalEDM,
		models.StatusAwalEDM, akun, akun}
}

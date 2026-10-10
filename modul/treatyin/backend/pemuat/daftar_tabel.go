//go:build ignore

// Mendaftar tabel skema beserta cacah barisnya — NOL TULIS.
//
//	go run modul/treatyin/backend/pemuat/daftar_tabel.go
//	go run modul/treatyin/backend/pemuat/daftar_tabel.go -pola T_TREATY_%
//
// ⛔ MENGAPA BERTANYA KE BASIS DATA, BUKAN MEMBACA MIGRASI
// Daftar `CREATE TABLE` di `backend/migrations/` adalah RIWAYAT, bukan
// keadaan: migrasi 436 mengganti nama sebelas tabel (`M_TREATYIN_*` ->
// `T_TREATY_*`), dan beberapa tabel lain dicabut sesudah dibuat
// (`MATA_UANG`, 4 Oktober 2026). Membacanya sebagai inventaris akan
// menyebut tabel yang sudah tidak ada dan melewatkan nama barunya.
//
// ⛔ NOL `DELETE`, `UPDATE`, `INSERT`, atau DDL di berkas ini. Ia hanya
// `SELECT` terhadap `ALL_TABLES` dan `COUNT(*)`.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"sort"
	"strings"

	"nusantarare/inti/backend/config"
	"nusantarare/inti/backend/db"
)

func main() {
	pola := flag.String("pola", "%", "saringan nama tabel (pola LIKE Oracle)")
	kolom := flag.String("kolom", "", "cetak kolom SATU tabel, bukan daftar")
	isi := flag.String("isi", "", "cetak 5 baris pertama SATU tabel (nol tulis)")
	migrasi := flag.Bool("migrasi", false, "cetak langkah migrasi yang SUDAH tercatat di T_MIGRASI")
	flag.Parse()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("konfigurasi: %v", err)
	}
	if !cfg.PunyaOracle() {
		log.Fatal("ORACLE_DSN kosong")
	}
	d, err := db.Open(cfg)
	if err != nil {
		log.Fatalf("membuka oracle: %v", err)
	}
	defer func() { _ = d.Close() }()

	ctx := context.Background()
	if err := d.Ping(ctx); err != nil {
		log.Fatalf("oracle tidak terjangkau: %v", err)
	}

	if *migrasi {
		cetakMigrasi(ctx, d)
		return
	}
	if *isi != "" {
		cetakIsi(ctx, d, strings.ToUpper(*isi))
		return
	}
	if *kolom != "" {
		cetakKolom(ctx, d, cfg.OracleSchema, strings.ToUpper(*kolom))
		cetakBatasan(ctx, d, cfg.OracleSchema, strings.ToUpper(*kolom))
		return
	}

	baris, err := d.QueryContext(ctx,
		`SELECT TABLE_NAME FROM ALL_TABLES WHERE OWNER = :1 AND TABLE_NAME LIKE :2 ORDER BY TABLE_NAME`,
		strings.ToUpper(cfg.OracleSchema), *pola)
	if err != nil {
		log.Fatalf("mendaftar tabel: %v", err)
	}
	var nama []string
	for baris.Next() {
		var n string
		if err := baris.Scan(&n); err != nil {
			log.Fatalf("membaca nama: %v", err)
		}
		nama = append(nama, n)
	}
	_ = baris.Close()
	sort.Strings(nama)

	fmt.Printf("skema %s · pola %s · %d tabel\n", cfg.OracleSchema, *pola, len(nama))
	for _, n := range nama {
		q, err := d.Qualify(n)
		if err != nil {
			fmt.Printf("  %-34s  (nama ditolak Qualify: %v)\n", n, err)
			continue
		}
		var c int64
		if err := d.QueryRowContext(ctx, "SELECT COUNT(*) FROM "+q).Scan(&c); err != nil {
			fmt.Printf("  %-34s  (gagal dicacah: %v)\n", n, err)
			continue
		}
		fmt.Printf("  %-34s %10d\n", n, c)
	}
}

// cetakKolom - bentuk satu tabel, untuk menulis jalur MUNDUR migrasi yang
// membuang tabel yang nol migrasi pernah membuatnya.
func cetakKolom(ctx context.Context, d *db.DB, skema, tabel string) {
	baris, err := d.QueryContext(ctx,
		`SELECT COLUMN_NAME, DATA_TYPE, CHAR_LENGTH, DATA_PRECISION, DATA_SCALE, NULLABLE, CHAR_USED
		   FROM ALL_TAB_COLUMNS WHERE OWNER = :1 AND TABLE_NAME = :2 ORDER BY COLUMN_ID`,
		strings.ToUpper(skema), tabel)
	if err != nil {
		log.Fatalf("membaca kolom: %v", err)
	}
	defer func() { _ = baris.Close() }()
	fmt.Printf("%s.%s\n", skema, tabel)
	for baris.Next() {
		var nama, jenis, boleh string
		var panjang int64
		var presisi, skala *int64
		var charUsed *string
		if err := baris.Scan(&nama, &jenis, &panjang, &presisi, &skala, &boleh, &charUsed); err != nil {
			log.Fatalf("scan kolom: %v", err)
		}
		cu := ""
		if charUsed != nil {
			cu = *charUsed
		}
		p, s2 := "-", "-"
		if presisi != nil {
			p = fmt.Sprint(*presisi)
		}
		if skala != nil {
			s2 = fmt.Sprint(*skala)
		}
		fmt.Printf("  %-28s %-12s charlen=%-5d presisi=%-4s skala=%-4s null=%s charused=%s\n",
			nama, jenis, panjang, p, s2, boleh, cu)
	}
}

// cetakBatasan - nama batasan dan kolomnya, plus sequence sejenis.
func cetakBatasan(ctx context.Context, d *db.DB, skema, tabel string) {
	b, err := d.QueryContext(ctx,
		`SELECT c.CONSTRAINT_NAME, c.CONSTRAINT_TYPE, cc.COLUMN_NAME, cc.POSITION
		   FROM ALL_CONSTRAINTS c JOIN ALL_CONS_COLUMNS cc
		     ON c.OWNER = cc.OWNER AND c.CONSTRAINT_NAME = cc.CONSTRAINT_NAME
		  WHERE c.OWNER = :1 AND c.TABLE_NAME = :2 AND c.CONSTRAINT_TYPE IN ('P','U','R')
		  ORDER BY c.CONSTRAINT_NAME, cc.POSITION`,
		strings.ToUpper(skema), tabel)
	if err != nil {
		log.Fatalf("membaca batasan: %v", err)
	}
	for b.Next() {
		var n, jenis, kol string
		var pos int64
		if err := b.Scan(&n, &jenis, &kol, &pos); err != nil {
			log.Fatalf("scan batasan: %v", err)
		}
		fmt.Printf("  BATASAN %-28s %s %s(%d)\n", n, jenis, kol, pos)
	}
	_ = b.Close()

	s2, err := d.QueryContext(ctx,
		`SELECT SEQUENCE_NAME FROM ALL_SEQUENCES WHERE SEQUENCE_OWNER = :1 AND SEQUENCE_NAME LIKE :2`,
		strings.ToUpper(skema), "%CURRENCY%")
	if err != nil {
		log.Fatalf("membaca sequence: %v", err)
	}
	defer func() { _ = s2.Close() }()
	for s2.Next() {
		var n string
		if err := s2.Scan(&n); err != nil {
			log.Fatalf("scan sequence: %v", err)
		}
		fmt.Printf("  SEQUENCE %s\n", n)
	}
}

// cetakMigrasi - nama langkah yang sudah tercatat. Dipakai untuk mengetahui
// apa yang TERTUNDA sebelum `-migrate` dijalankan: pelarinya nol punya
// penargetan, jadi ia menjalankan SELURUH yang tertunda lintas modul.
func cetakMigrasi(ctx context.Context, d *db.DB) {
	t, err := d.Qualify("T_MIGRASI")
	if err != nil {
		log.Fatalf("qualify: %v", err)
	}
	b, err := d.QueryContext(ctx, "SELECT NAMA FROM "+t+" ORDER BY NAMA")
	if err != nil {
		log.Fatalf("membaca T_MIGRASI: %v", err)
	}
	defer func() { _ = b.Close() }()
	n := 0
	for b.Next() {
		var nama string
		if err := b.Scan(&nama); err != nil {
			log.Fatalf("scan: %v", err)
		}
		fmt.Println(nama)
		n++
	}
	fmt.Printf("TOTAL %d langkah tercatat\n", n)
}

// cetakIsi - lima baris pertama satu tabel. NOL tulis; dipakai memastikan
// kolom yang layar tampilkan memang berisi di basis data.
func cetakIsi(ctx context.Context, d *db.DB, tabel string) {
	q, err := d.Qualify(tabel)
	if err != nil {
		log.Fatalf("qualify: %v", err)
	}
	baris, err := d.QueryContext(ctx, "SELECT * FROM "+q+" WHERE ROWNUM <= 5")
	if err != nil {
		log.Fatalf("membaca isi: %v", err)
	}
	defer func() { _ = baris.Close() }()
	kol, err := baris.Columns()
	if err != nil {
		log.Fatalf("kolom: %v", err)
	}
	fmt.Println(strings.Join(kol, " | "))
	for baris.Next() {
		sel := make([]any, len(kol))
		pt := make([]any, len(kol))
		for i := range sel {
			pt[i] = &sel[i]
		}
		if err := baris.Scan(pt...); err != nil {
			log.Fatalf("scan: %v", err)
		}
		var teks []string
		for _, v := range sel {
			if v == nil {
				teks = append(teks, "<NULL>")
				continue
			}
			if b, ok := v.([]byte); ok {
				teks = append(teks, string(b))
				continue
			}
			teks = append(teks, fmt.Sprint(v))
		}
		fmt.Println(strings.Join(teks, " | "))
	}
}

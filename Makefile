# Nusantara Re - Makefile aplikasi. Jalankan dari folder APP_RNM/ ini.
#
# Bila `make` belum terpasang, tiap target adalah satu-dua perintah biasa:
# lihat baris di bawah nama targetnya, dan PANDUAN-RANTAI-ALAT.md di folder induk.
#
# Nol alamat host di berkas ini. Seluruh koneksi datang dari env var
# (ADR-U-0004); daftar endpoint sesungguhnya ada di tabel Oracle
# M_LINK_SERVICE yang isinya belum ada di korpus (OQ-047).

GO      ?= go
NPM     ?= npm
API_PKG := ./cmd/api
BIN     := bin/api

# Skema uji yang dibuat oleh `make db-up`.
#
# TANPA BAWAAN, dan itu disengaja. Sampai 26-09-2026 baris ini berbunyi
# `?= POOLDATA` - nama skema warisan Pega. `make test-db` MENGHAPUS tabel di
# skema yang ditunjuk, termasuk OS_AKSEPTASI_KLAIM_LIFE, jadi bawaan itu cukup
# untuk menghapus tabel warisan sungguhan di instance pengembangan. Nilainya
# kini wajib datang dari DBA lewat env, dan skemauji menolak POOLDATA apa pun
# yang terjadi.
ORACLE_SCHEMA ?=
ORACLE_DEV_CONTAINER ?= nusantarare-oracle

.PHONY: help check generate typecheck run-api run-web build build-api build-web test test-db test-db-treatyin db-up db-down migrate

help:
	@echo "run-api    - jalankan backend Go"
	@echo "run-web    - jalankan frontend Vite"
	@echo "build      - build kedua sisi"
	@echo "generate   - tulis ulang daftar modul Go (inti/backend/daftar) dari folder modul/"
	@echo "test       - go vet + go test (tanpa Oracle)"
	@echo "typecheck  - periksa tipe TypeScript frontend (tsc --noEmit)"
	@echo "test-db    - test seam repository terhadap skema uji Oracle"
	@echo "             (MENGHAPUS tabel; perlu ORACLE_SKEMA_UJI=true)"
	@echo "db-up      - kontainer Oracle pengembangan [usulan]"
	@echo "db-down    - hentikan dan buang kontainer itu"
	@echo "migrate    - jalankan migrasi terhadap ORACLE_DSN"

run-api:
	$(GO) run $(API_PKG)

# Sejak 30-09-2026 (struktur tim satu folder per modul) npm dijalankan dari
# folder ini: package.json dan vite.config.ts tinggal di APP_RNM/, bukan frontend/.
run-web:
	$(NPM) run dev

build: build-api build-web

build-api:
	$(GO) build -o $(BIN) $(API_PKG)

build-web:
	$(NPM) ci && $(NPM) run build

# Daftar modul Go dibangkitkan dari folder modul/<nama>/backend/modul.go
# (struktur tim satu folder per modul, 30-09-2026). Hasilnya di-commit; uji
# inti/backend/daftar/bangkit menagihnya bila lupa dijalankan.
generate:
	$(GO) generate ./inti/backend/daftar

# Nol Oracle: hanya seam `services` murni dan unit lain yang tidak menyentuh DB.
test:
	$(GO) vet ./...
	$(GO) test ./...

# Frontend memakai TypeScript (.tsx). Vite tidak memeriksa tipe saat build,
# jadi pemeriksaannya dijalankan tersendiri di sini dan di awal `npm run build`.
typecheck:
	$(NPM) run typecheck

# Gabungan yang diminta tiket 01 Claim Life: build Go, test Go, tipe + test web.
# Nol Oracle - seam repository dan HTTP dijalankan terpisah lewat `make test-db`.
check:
	$(GO) build ./...
	$(GO) vet ./...
	$(GO) test ./...
	$(NPM) run typecheck
	$(NPM) test

# Seam `repository` terhadap skema uji Oracle nyata (brief bab 5).
# Ditandai build tag `db` supaya `make test` tetap lulus tanpa instance.
# ⛔ MENGHAPUS tabel di skema yang ditunjuk ORACLE_SCHEMA. Menuntut pengakuan
# sadar lewat ORACLE_SKEMA_UJI=true, dan menolak POOLDATA.
# Uji db yang MENUNTUT skema uji buangan (uji/skemauji memasang lalu MEMBONGKAR
# tabel tiruan). Pagarnya menolak POOLDATA, dan penolakan itu BENAR: daftar
# bongkarnya memuat OS_AKSEPTASI_KLAIM_LIFE dan empat tabel warisan lain.
test-db:
	$(GO) test -tags=db ./modul/claimlife/backend/repository/...

# Uji db Treaty In - berjalan di SKEMA APLIKASI (POOLDATA), sebab ke-37
# tabelnya berdiri di sana (keputusan pemilik proses 03-10-2026). Ia TIDAK
# memakai uji/skemauji: nol pemasangan, nol pembongkaran. Tiap test satu
# transaksi yang diakhiri Rollback, jadi nol baris menetap.
#
# Dipisah dari `test-db` dengan sengaja: menjalankan keduanya dalam satu
# ORACLE_SCHEMA tidak mungkin - yang satu menuntut skema buangan, yang lain
# menuntut skema aplikasi.
# ⛔ GAGAL KERAS bila ORACLE_DSN kosong, dan itu bukan kerewelan.
# `siapkan(t)` memanggil `t.Skip` ketika DSN kosong, dan `go test` mencetak
# `ok` untuk paket yang SELURUH testnya dilewati. Tanpa pagar di bawah,
# `make test-db-treatyin` berhasil tanpa menyentuh Oracle sama sekali - dan
# `ok ... 1.6s` tidak dapat dibedakan dari lulus sungguhan tanpa `-v`.
# Ditemukan 3 Oktober 2026: target ini memang berperilaku begitu.
#
# `config.Load()` membaca os.Getenv, BUKAN .env - jadi env-nya harus sudah
# termuat di shell pemanggil:  set -a && . ./.env && set +a && make test-db-treatyin
test-db-treatyin:
	@if [ -z "$$ORACLE_DSN" ]; then 		echo "ORACLE_DSN kosong - uji db akan DILEWATI seluruhnya dan mencetak 'ok'."; 		echo "Muat env lebih dulu:  set -a && . ./.env && set +a"; 		exit 1; 	fi
	$(GO) test -tags=db -count=1 -v ./modul/treatyin/backend/repository/... | grep -Ev '^(=== RUN|=== PAUSE|=== CONT)'

# [usulan] gvenzl/oracle-free. Kata sandi datang dari env, tidak pernah
# ditanam di berkas ini. Tidak pernah menunjuk instance produksi.
db-up:
	@test -n "$$ORACLE_DEV_PASSWORD" || (echo "setel ORACLE_DEV_PASSWORD dulu" && exit 1)
	@test -n "$(ORACLE_SCHEMA)" || (echo "setel ORACLE_SCHEMA dulu - bawaannya dicabut 26-09-2026" && exit 1)
	@test "$(ORACLE_SCHEMA)" != "POOLDATA" || (echo "ORACLE_SCHEMA tidak boleh POOLDATA" && exit 1)
	docker run -d --name $(ORACLE_DEV_CONTAINER) -p 1521:1521 \
	  -e ORACLE_PASSWORD="$$ORACLE_DEV_PASSWORD" \
	  -e APP_USER="$(ORACLE_SCHEMA)" \
	  -e APP_USER_PASSWORD="$$ORACLE_DEV_PASSWORD" \
	  gvenzl/oracle-free:latest

db-down:
	docker rm -f $(ORACLE_DEV_CONTAINER)

migrate:
	$(GO) run $(API_PKG) -migrate

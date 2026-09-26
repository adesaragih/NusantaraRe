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

.PHONY: help check typecheck run-api run-web build build-api build-web test test-db db-up db-down migrate

help:
	@echo "run-api    - jalankan backend Go"
	@echo "run-web    - jalankan frontend Vite"
	@echo "build      - build kedua sisi"
	@echo "test       - go vet + go test (tanpa Oracle)"
	@echo "typecheck  - periksa tipe TypeScript frontend (tsc --noEmit)"
	@echo "test-db    - test seam repository terhadap skema uji Oracle"
	@echo "             (MENGHAPUS tabel; perlu ORACLE_SKEMA_UJI=true)"
	@echo "db-up      - kontainer Oracle pengembangan [usulan]"
	@echo "db-down    - hentikan dan buang kontainer itu"
	@echo "migrate    - jalankan migrasi terhadap ORACLE_DSN"

run-api:
	$(GO) run $(API_PKG)

run-web:
	cd frontend && $(NPM) run dev

build: build-api build-web

build-api:
	$(GO) build -o $(BIN) $(API_PKG)

build-web:
	cd frontend && $(NPM) ci && $(NPM) run build

# Nol Oracle: hanya seam `services` murni dan unit lain yang tidak menyentuh DB.
test:
	$(GO) vet ./...
	$(GO) test ./...

# Frontend memakai TypeScript (.tsx). Vite tidak memeriksa tipe saat build,
# jadi pemeriksaannya dijalankan tersendiri di sini dan di awal `npm run build`.
typecheck:
	cd frontend && $(NPM) run typecheck

# Gabungan yang diminta tiket 01 Claim Life: build Go, test Go, tipe + test web.
# Nol Oracle - seam repository dan HTTP dijalankan terpisah lewat `make test-db`.
check:
	$(GO) build ./...
	$(GO) vet ./...
	$(GO) test ./...
	cd frontend && $(NPM) run typecheck
	cd frontend && $(NPM) test

# Seam `repository` terhadap skema uji Oracle nyata (brief bab 5).
# Ditandai build tag `db` supaya `make test` tetap lulus tanpa instance.
# ⛔ MENGHAPUS tabel di skema yang ditunjuk ORACLE_SCHEMA. Menuntut pengakuan
# sadar lewat ORACLE_SKEMA_UJI=true, dan menolak POOLDATA.
test-db:
	$(GO) test -tags=db ./internal/repository/...

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

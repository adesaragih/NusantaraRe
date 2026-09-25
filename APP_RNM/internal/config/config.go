// Package config membaca seluruh konfigurasi dari environment variable.
//
// ⛔ Nol literal host, endpoint, kata sandi, atau DSN di berkas ini maupun di
// berkas lain mana pun (ADR-U-0004, CLAUDE.md bab 10). Daftar endpoint
// sesungguhnya ada di tabel Oracle M_LINK_SERVICE yang isinya belum ada di
// korpus (OQ-047) - sampai ia ada, alamat datang dari env.
package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Prefix env var untuk alamat layanan luar: SERVICE_<NAMA>.
const prefixLayanan = "SERVICE_"

// Config adalah seluruh konfigurasi proses.
type Config struct {
	// HTTPAddr alamat dengar server. Hanya port secara bawaan; tidak ada
	// nama host yang ditanam.
	HTTPAddr string

	// OracleDSN kosong berarti backend berjalan tanpa Oracle. Fase 0
	// mengizinkannya supaya GET /healthz dapat dijawab tanpa instance; tiket
	// penyimpanan mana pun akan menuntutnya terisi.
	OracleDSN string

	// OracleSchema wajib terisi bila OracleDSN terisi: setiap query menyebut
	// nama skema secara eksplisit (ADR-U-0033).
	OracleSchema string

	// IsPegaProd menandai lingkungan yang menunjuk data produksi Pega
	// (ADR-U-0005). Dipakai sebagai gerbang perilaku, bukan sebagai hiasan.
	IsPegaProd bool

	// Layanan memetakan nama layanan luar ke alamatnya, seluruhnya dari env.
	Layanan map[string]string
}

// ErrKonfigurasi membungkus seluruh kegagalan pembacaan konfigurasi.
var ErrKonfigurasi = errors.New("konfigurasi")

// Load membaca konfigurasi dari environment.
func Load() (Config, error) {
	c := Config{
		HTTPAddr:     ambil("HTTP_ADDR", ":8080"),
		OracleDSN:    os.Getenv("ORACLE_DSN"),
		OracleSchema: strings.TrimSpace(os.Getenv("ORACLE_SCHEMA")),
		Layanan:      map[string]string{},
	}

	raw := strings.TrimSpace(os.Getenv("IS_PEGA_PROD"))
	if raw != "" {
		b, err := strconv.ParseBool(raw)
		if err != nil {
			return Config{}, fmt.Errorf("%w: IS_PEGA_PROD bukan boolean: %q", ErrKonfigurasi, raw)
		}
		c.IsPegaProd = b
	}

	if c.OracleDSN != "" && c.OracleSchema == "" {
		return Config{}, fmt.Errorf(
			"%w: ORACLE_SCHEMA wajib terisi bila ORACLE_DSN terisi (ADR-U-0033)", ErrKonfigurasi)
	}

	for _, kv := range os.Environ() {
		nama, nilai, ada := strings.Cut(kv, "=")
		if !ada || !strings.HasPrefix(nama, prefixLayanan) || nilai == "" {
			continue
		}
		c.Layanan[strings.TrimPrefix(nama, prefixLayanan)] = nilai
	}

	return c, nil
}

// PunyaOracle menyatakan apakah proses dikonfigurasi menyentuh Oracle.
func (c Config) PunyaOracle() bool { return c.OracleDSN != "" }

func ambil(nama, bawaan string) string {
	if v := strings.TrimSpace(os.Getenv(nama)); v != "" {
		return v
	}
	return bawaan
}

package handlers_test

// Status HTTP galat penyimpanan nyata (keputusan work owner 03-10-2026): layanan gagal 502, belum siap 503, berkas
// tidak ada di penyimpanan 409 - kalimatnya tanpa sebab mentah.

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"testing/iotest"

	"nusantarare/modul/masterproductnamelife/backend/handlers"
	"nusantarare/modul/masterproductnamelife/backend/models"
	"nusantarare/modul/masterproductnamelife/backend/services"
	"nusantarare/modul/masterproductnamelife/backend/tiruan"
)

// storagePalsu - antrean di memori; Buka dan Hapus menjawab `galat`; `putus` = isi gagal dibaca di tengah.
type storagePalsu struct {
	antre map[string][]byte
	galat error
	putus bool
}

func (s *storagePalsu) SimpanAntrean(_ context.Context, id string, isi io.Reader) error {
	d, err := io.ReadAll(isi)
	s.antre[id] = d
	return err
}
func (s *storagePalsu) BuangAntrean(_ context.Context, id string) { delete(s.antre, id) }
func (s *storagePalsu) Kirim(_ context.Context, o models.ObjekPenyimpanan, _, _ string) (models.ObjekPenyimpanan, error) {
	o.URLPublic = "UJI-URL"
	return o, nil
}
func (s *storagePalsu) Buka(context.Context, models.ObjekPenyimpanan) (io.ReadCloser, *models.ObjekPenyimpanan, error) {
	if s.galat != nil {
		return nil, nil, s.galat
	}
	if s.putus {
		return io.NopCloser(io.MultiReader(bytes.NewReader([]byte("isi")), iotest.ErrReader(errors.New("UJI putus")))), nil, nil
	}
	return io.NopCloser(bytes.NewReader([]byte("isi"))), nil, nil
}
func (s *storagePalsu) Hapus(context.Context, string, *models.ObjekPenyimpanan) error { return s.galat }

func TestRuteGalatPenyimpananNyata(t *testing.T) {
	g := tiruan.Baru()
	g.IsiJSON("100007", `{"ID":"100007"}`, "")
	g.AppName = "UJI-APP"
	st := &storagePalsu{antre: map[string][]byte{}}
	l := services.BaruLayanan(g, g.Transaksi, nil).DenganPenyimpanan(st)
	srv := httptest.NewServer(handlers.RouterDengan(l, true, true))
	t.Cleanup(srv.Close)
	dasar := srv.URL + handlers.Prefix + "/produk/100007/lampiran"
	kode, badan := kirimBerkasUji(t, dasar, "UJI.pdf", "isi")
	if kode != http.StatusOK {
		t.Fatalf("unggah: %d %s", kode, badan)
	}
	var a models.Lampiran
	if err := json.Unmarshal([]byte(badan), &a); err != nil || a.Status != models.StatusTerunggah {
		t.Fatalf("lampiran: %+v %v", a, err)
	}
	for _, k := range []struct {
		galat error
		kode  int
	}{
		{fmt.Errorf("%w: UJI-RINCI", services.ErrStorageGagal), http.StatusBadGateway},
		{fmt.Errorf("%w: UJI-RINCI", services.ErrStorageBelumSiap), http.StatusServiceUnavailable},
		{fmt.Errorf("%w: UJI-RINCI", services.ErrBerkasTidakDiStorage), http.StatusConflict},
	} {
		st.galat = k.galat
		kode, b, _ := ambil(t, "GET", dasar+"/"+a.ID+"/unduh")
		if kode != k.kode || !strings.Contains(string(b), strings.TrimPrefix(errors.Unwrap(k.galat).Error(), "services: ")) {
			t.Errorf("unduh %v: %d %s", k.galat, kode, b)
		}
		kode, b, _ = ambil(t, "DELETE", dasar+"/"+a.ID)
		if kode != k.kode {
			t.Errorf("hapus %v: %d %s", k.galat, kode, b)
		}
	}
}

// Unduhan yang putus di tengah TIDAK berakhir sebagai 200 berisi berkas terpotong: sambungan diputus.
func TestRuteUnduhPutusTidakTerpotongDiamDiam(t *testing.T) {
	g := tiruan.Baru()
	g.IsiJSON("100007", `{"ID":"100007"}`, "")
	g.AppName = "UJI-APP"
	st := &storagePalsu{antre: map[string][]byte{}}
	l := services.BaruLayanan(g, g.Transaksi, nil).DenganPenyimpanan(st)
	srv := httptest.NewServer(handlers.RouterDengan(l, true, true))
	t.Cleanup(srv.Close)
	dasar := srv.URL + handlers.Prefix + "/produk/100007/lampiran"
	_, badan := kirimBerkasUji(t, dasar, "UJI.pdf", "isi")
	var a models.Lampiran
	_ = json.Unmarshal([]byte(badan), &a)
	st.putus = true
	req, _ := http.NewRequest("GET", dasar+"/"+a.ID+"/unduh", nil)
	req.Header.Set("X-Pelaku", "UJI-PELAKU")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return // sambungan diputus sebelum kepala - juga bukan berkas terpotong
	}
	defer func() { _ = res.Body.Close() }()
	if _, err := io.ReadAll(res.Body); err == nil {
		t.Errorf("status %d dengan badan utuh - berkas terpotong diterima sebagai lengkap", res.StatusCode)
	}
}

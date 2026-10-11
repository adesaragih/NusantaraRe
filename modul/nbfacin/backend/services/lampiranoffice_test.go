package services

import (
	"context"
	"errors"
	"net/url"
	"strings"
	"testing"
	"time"

	"nusantarare/modul/nbfacin/backend/models"
)

// urlUjiOffice - URL penyimpanan tiruan (disusun lewat url.URL: ADR-U-0013 melarang literal alamat di kode).
func urlUjiOffice(kueri string) string {
	return (&url.URL{Scheme: "https", Host: "penyimpanan.uji", Path: "/objek", RawQuery: kueri}).String()
}

// storagePenautTiruan - klien penyimpanan "nyata" tiruan: Tautan menjawab URL tetap (atau galat).
type storagePenautTiruan struct {
	KlienStorageLampiran
	url   string
	galat error
	link  []models.LinkStorage
}

func (s *storagePenautTiruan) Tautan(_ context.Context, link models.LinkStorage) (string, error) {
	s.link = append(s.link, link)
	return s.url, s.galat
}

// TestTautanOfficeLampiran - GET …/office: ekstensi kantor saja (400), case lain (404), belum tersimpan (409), stub lokal (409),
// penyimpanan nyata -> URL; objek hilang -> ErrLampiranTanpaBerkas.
func TestTautanOfficeLampiran(t *testing.T) {
	ctx := context.Background()
	r := rangkaUji(t, "0")
	xlsx, _ := r.svc.UnggahLampiran(ctx, pelakuPemegang, kasusLampiran, berkasUji("Hitungan.XLSX", "OFFER", "ISI-XLSX-UJI"))
	pdf, _ := r.svc.UnggahLampiran(ctx, pelakuPemegang, kasusLampiran, berkasUji("a.pdf", "OFFER", "ISI-PDF-UJI"))
	if _, err := r.svc.TautanOfficeLampiran(ctx, pelakuPemegang, kasusLampiran, pdf); !errors.Is(err, ErrLampiranBukanOffice) ||
		!errors.Is(err, ErrMasukanLampiran) {
		t.Errorf("pdf: %v", err)
	}
	// stub lokal: URLPUBLIC kosong dan klien tanpa Tautan.
	if _, err := r.svc.TautanOfficeLampiran(ctx, pelakuLain, kasusLampiran, xlsx); !errors.Is(err, ErrLampiranOfficeTanpaStorage) {
		t.Errorf("stub: %v", err)
	}
	r.svc.DenganTahap(tahapTiruan{kasusLampiran: {Position: "ReasFacInMarketing", FlagOnGoingPolicy: "0"}, "UJI-NB-2": {}})
	if _, err := r.svc.TautanOfficeLampiran(ctx, pelakuPemegang, "UJI-NB-2", xlsx); !errors.Is(err, ErrLampiranTidakAda) {
		t.Errorf("case lain: %v", err)
	}
	// penyimpanan nyata (tiruan): catatan T_STORAGE_IMAGE ber-URLPUBLIC.
	p := &storagePenautTiruan{KlienStorageLampiran: r.storage, url: urlUjiOffice("tanda=X")}
	r.svc.storageLampiran = p
	// klien nyata, tetapi catatan objek tanpa URLPUBLIC (diunggah ke stub sebelum pindah ke penyimpanan nyata) -> 409.
	if _, err := r.svc.TautanOfficeLampiran(ctx, pelakuPemegang, kasusLampiran, xlsx); !errors.Is(err, ErrLampiranOfficeTanpaStorage) ||
		len(p.link) != 0 {
		t.Errorf("nyata tanpa URLPUBLIC: %v %d", err, len(p.link))
	}
	img := r.repo.dok[xlsx].TStorageID
	o := r.repo.objek[img]
	o.URLPublic = urlUjiOffice("")
	r.repo.objek[img] = o
	u, err := r.svc.TautanOfficeLampiran(ctx, pelakuPemegang, kasusLampiran, xlsx)
	if err != nil || u != p.url || len(p.link) != 1 || p.link[0].ImageID != img {
		t.Errorf("nyata: %v %q %+v", err, u, p.link)
	}
	p.galat = ErrBerkasTidakAdaDiPenyimpanan
	if _, err := r.svc.TautanOfficeLampiran(ctx, pelakuPemegang, kasusLampiran, xlsx); !errors.Is(err, ErrLampiranTanpaBerkas) {
		t.Errorf("hilang: %v", err)
	}
	// belum tersimpan (unggah gagal -> tanpa catatan T_STORAGE_IMAGE).
	r2 := rangkaUji(t, "0")
	r2.storage.galatSimpan = ErrStorageTakTerjangkau
	doc, _ := r2.svc.UnggahLampiran(ctx, pelakuPemegang, kasusLampiran, berkasUji("b.docx", "OFFER", "ISI-DOCX-UJI"))
	if _, err := r2.svc.TautanOfficeLampiran(ctx, pelakuPemegang, kasusLampiran, doc); !errors.Is(err, ErrLampiranBelumTersimpan) {
		t.Errorf("belum tersimpan: %v", err)
	}
}

// TestStorageJarakJauhTautan - URL tersimpan berlaku: tanpa geturl; lewat: geturl + penyegaran; EXPDATE kosong: tanpa objek; URL
// bukan https ditolak dan TIDAK muncul di pesan galat.
func TestStorageJarakJauhTautan(t *testing.T) {
	ctx := context.Background()
	l := layananStorageUji(t)
	kini := time.Date(2026, 10, 8, 6, 0, 0, 0, time.UTC)
	klien, _, pc := storageJarakJauhUji(t, l, kini)
	st, nyata := klien.(penautStorage) // klien penyimpanan nyata wajib dapat memberi tautan (View Office Online)
	if !nyata {
		t.Fatal("storageJarakJauh tanpa Tautan")
	}
	link := models.LinkStorage{ImageID: lokUji.ImageID, URLPublic: l.srv.URL + "/objek/a", ExpUTC: "20261008063000",
		AppFolder: models.AppFolderUji("UJIAPP", "Policy/Doc/2026/10/"+lokUji.Namafile), App: "UJIAPP", Namafile: lokUji.Namafile}
	if u, err := st.Tautan(ctx, link); err != nil || u != l.srv.URL+"/objek/a" || len(l.minta["geturl"]) != 0 {
		t.Errorf("tersimpan: %v %q %d", err, u, len(l.minta["geturl"]))
	}
	link.ExpUTC = "20261008060030"
	l.jawaban["geturl"] = map[string]string{"URLImage": l.srv.URL + "/objek/b", "exp": "20261008T070000"}
	if u, err := st.Tautan(ctx, link); err != nil || u != l.srv.URL+"/objek/b" || len(l.minta["geturl"]) != 1 || len(pc.o) != 1 {
		t.Errorf("geturl: %v %q %d %d", err, u, len(l.minta["geturl"]), len(pc.o))
	}
	polos := (&url.URL{Scheme: "http", Host: "polos", Path: "/objek/rahasia"}).String()
	l.jawaban["geturl"] = map[string]string{"URLImage": polos}
	if _, err := st.Tautan(ctx, link); !errors.Is(err, ErrStorageJawabanRusak) || strings.Contains(err.Error(), "rahasia") {
		t.Errorf("http: %v", err)
	}
	link.ExpUTC = ""
	if _, err := st.Tautan(ctx, link); !errors.Is(err, ErrBerkasTidakAdaDiPenyimpanan) {
		t.Errorf("tanpa EXPDATE: %v", err)
	}
}

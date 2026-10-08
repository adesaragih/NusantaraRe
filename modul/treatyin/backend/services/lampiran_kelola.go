package services

// Modal `View File` panel Attachment (`Section/ShowAttachmentTreaty.xml`) —
// unduh, hapus, dan ganti kategori (permintaan pemakai 8 Oktober 2026).
//
//	nama berkas (tautan)  DownloadAttachmentTreaty(ImageID=.type, ViewOffice=false)
//	View Office Online    idem ViewOffice=true — hanya xls/xlsx/doc/docx/ppt/pptx
//	                      → `https://view.officeapps.live.com/op/view.aspx?src=` +
//	                      @encodeURL(url)
//	Delete                Delete_act(ID, ImageID) — `ViewState !='1' ||
//	                      RevisionState='1'`: DeleteGoogleStorage_Act (Google/
//	                      delete, Namafile = APPFOLDER tanpa `gs://<App>/`;
//	                      [10] DeleteStorage_SQL) → [3] DeleteAttachment2_Sql
//	Change Category/Save  ChangeDokument_Act("Change"/"Save") — status bukan
//	                      Resolve Complete/Decline: ChangeKateAttachment2_Sql
//	                      per baris
//
// GetUrlGoogleStorage_Act: URL tersimpan dipakai selama EXPDATE belum lewat;
// selain itu Google/geturl (Folder = APPFOLDER tanpa Namafile dan tanpa
// `gs://<App>/`, Durasi 1800) lalu Update_T_Storage_SQL.
//
// ⚠️ PENYIMPANGAN yang dinyatakan:
//   - Hapus di penyimpanan GAGAL → baris TETAP (galat dilaporkan, dapat
//     diulang) — Pega melanjutkan menghapus baris walau objeknya tertinggal.
//   - Jawaban geturl tanpa `appfolder` tidak mengosongkan APPFOLDER tersimpan
//     (Pega menimpanya kosong, dan jalur objek hilang untuk selamanya).
//   - Nama kategori (`CATEGORY`) diambil dari katalog yang SUDAH dirapikan dan
//     mengikuti nama Non-Prop untuk `00007` — sama dengan unggah; SQL Pega
//     menyalin `NOTE` mentah (ekor CR LF `00004`).

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"strings"
	"time"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/unggah"
	"nusantarare/modul/treatyin/backend/models"
)

// penampilOffice - `DownloadAttachmentTreaty` [3].
const penampilOffice = "https://view.officeapps.live.com/op/view.aspx?src="

// jenisOffice - syarat tampil `View Office Online`.
var jenisOffice = map[string]bool{"xls": true, "xlsx": true, "doc": true, "docx": true, "ppt": true, "pptx": true}

// JenisOffice - berkas ini punya tombol `View Office Online`.
func JenisOffice(ext string) bool { return jenisOffice[strings.ToLower(strings.TrimSpace(ext))] }

// MasukanUbahKategori - satu baris `Change Category` → `Save`.
type MasukanUbahKategori struct {
	IDLampiran   string `json:"id"`
	KodeKategori string `json:"kategori"`
}

// namaKategoriTampil - nama yang panel TAMPILKAN untuk satu kode
// (`StatusDoc.CARI41`): katalog, lalu nama Non-Prop untuk kode `00007`.
func namaKategoriTampil(nama string, proporsional bool) string {
	if proporsional {
		return nama
	}
	for nonProp, prop := range namaKategoriSetara {
		if nama == prop {
			return nonProp
		}
	}
	return nama
}

// lampiranKontrak - baris lampiran `idLampiran` milik kontrak `idKontrak`.
func (l *Layanan) lampiranKontrak(ctx context.Context, idKontrak, idLampiran string) (models.BarisLampiranWarisan, bool, error) {
	daftar, err := l.gudang.BacaLampiranKontrak(ctx, idKontrak)
	if err != nil {
		return models.BarisLampiranWarisan{}, false, err
	}
	for _, b := range daftar {
		if b.ID == idLampiran {
			return b, true, nil
		}
	}
	return models.BarisLampiranWarisan{}, false, nil
}

// kepalaKontrakLampiran - kepala `TREATY_IN` kontrak yang lampirannya disentuh.
func (l *Layanan) kepalaKontrakLampiran(ctx context.Context, id string) (map[string]any, error) {
	kepala, ada, err := l.gudang.BacaKepalaTreatyIn(ctx, id)
	if err != nil {
		return nil, err
	}
	if !ada {
		return nil, fmt.Errorf("%w: %s", ErrKontrakTidakAda, id)
	}
	return kepala, nil
}

// TautanLampiran - `DownloadAttachmentTreaty`: URL bertanda tangan berkas,
// dibungkus penampil Office bila `office`.
func (l *Layanan) TautanLampiran(ctx context.Context, p inti.Pelaku, idKontrak, idLampiran string, office bool) (string, error) {
	if err := inti.WajibIdentitas(p); err != nil {
		return "", err
	}
	b, tautan, err := l.tautanBertanda(ctx, idKontrak, idLampiran)
	if err != nil {
		return "", err
	}
	if office {
		if !JenisOffice(b.JenisMime) {
			return "", ditolak("View Office Online hanya untuk berkas xls/xlsx/doc/docx/ppt/pptx.")
		}
		return penampilOffice + url.QueryEscape(tautan), nil
	}
	return tautan, nil
}

// BerkasUnduhan - isi satu lampiran untuk dialirkan ke layar.
type BerkasUnduhan struct {
	Nama string
	Mime string
	Isi  io.ReadCloser
}

// IsiLampiran - tautan nama berkas (`DownloadAttachmentTreaty`): isi berkas
// dari URL bertanda tangan, DIALIRKAN backend. Layar mengunduhnya lewat
// `fetch` beridentitas (`unduhBerkasBeridentitas`), bukan membuka URL —
// penjaga lintas modul `unduhdokumen.test.ts`.
func (l *Layanan) IsiLampiran(ctx context.Context, p inti.Pelaku, idKontrak, idLampiran string) (BerkasUnduhan, error) {
	if err := inti.WajibIdentitas(p); err != nil {
		return BerkasUnduhan{}, err
	}
	b, tautan, err := l.tautanBertanda(ctx, idKontrak, idLampiran)
	if err != nil {
		return BerkasUnduhan{}, err
	}
	if l.simpanan == nil {
		return BerkasUnduhan{}, fmt.Errorf("%w: the storage sender is not wired", ErrSimpananBelumSiap)
	}
	isi, err := l.simpanan.Ambil(ctx, tautan)
	if err != nil {
		return BerkasUnduhan{}, err
	}
	return BerkasUnduhan{Nama: b.NamaBerkas, Mime: unggah.MimeDariNamaFile(b.NamaBerkas), Isi: isi}, nil
}

// tautanBertanda - `GetUrlGoogleStorage_Act`: URL tersimpan selama EXPDATE
// berlaku; selain itu Google/geturl lalu `Update_T_Storage_SQL`.
func (l *Layanan) tautanBertanda(ctx context.Context, idKontrak, idLampiran string) (models.BarisLampiranWarisan, string, error) {
	b, ada, err := l.lampiranKontrak(ctx, strings.TrimSpace(idKontrak), strings.TrimSpace(idLampiran))
	if err != nil {
		return b, "", err
	}
	if !ada {
		return b, "", fmt.Errorf("%w: lampiran %s", ErrKontrakTidakAda, idLampiran)
	}
	o, ada, err := l.gudang.BacaObjekSimpanan(ctx, b.IDSimpanan)
	if err != nil {
		return b, "", err
	}
	if !ada || strings.TrimSpace(o.URLPublik) == "" {
		return b, "", ditolak(fmt.Sprintf("Berkas %s tidak punya objek di penyimpanan (T_STORAGE_IMAGE).", b.NamaBerkas))
	}
	tautan := o.URLPublik
	// [6] JIKA EXPDATE SUDAH EXPIRED — URL tersimpan dipakai selama berlaku.
	if exp, err := time.ParseInLocation("02/01/2006 15:04:05", strings.TrimSpace(o.Exp), zonaLampiran); err != nil || !exp.After(time.Now()) {
		if l.simpanan == nil {
			return b, "", fmt.Errorf("%w: the storage sender is not wired", ErrSimpananBelumSiap)
		}
		// [6.2] Folder = APPFOLDER tanpa Namafile, tanpa `gs://<App>/`.
		folder := strings.ReplaceAll(strings.ReplaceAll(o.AppFolder, o.NamaObjek, ""), "gs://"+o.App+"/", "")
		j, err := l.simpanan.URLBaru(ctx, PermintaanSimpanan{App: o.App, Durasi: durasiLampiran, Folder: folder, Namafile: o.NamaObjek})
		if err != nil {
			return b, "", err
		}
		// [6.6] dilewati bila URLImage kosong → URL tersimpan.
		if strings.TrimSpace(j.URLImage) != "" {
			o.URLPublik, o.Exp = j.URLImage, expSimpanan(j.Exp)
			if strings.TrimSpace(j.AppFolder) != "" {
				o.AppFolder = j.AppFolder
			}
			if err := l.gudang.PerbaruiObjekSimpanan(ctx, o, j.DateTime); err != nil {
				return b, "", err
			}
			tautan = j.URLImage
		}
	}
	return b, tautan, nil
}

// HapusLampiran - `Delete_act`: objek di penyimpanan lebih dulu, lalu baris.
func (l *Layanan) HapusLampiran(ctx context.Context, p inti.Pelaku, idKontrak, idLampiran string) (PanelLampiran, error) {
	if err := inti.WajibIdentitas(p); err != nil {
		return PanelLampiran{}, err
	}
	idKontrak, idLampiran = strings.TrimSpace(idKontrak), strings.TrimSpace(idLampiran)
	kepala, err := l.kepalaKontrakLampiran(ctx, idKontrak)
	if err != nil {
		return PanelLampiran{}, err
	}
	b, ada, err := l.lampiranKontrak(ctx, idKontrak, idLampiran)
	if err != nil {
		return PanelLampiran{}, err
	}
	if !ada {
		return PanelLampiran{}, fmt.Errorf("%w: lampiran %s", ErrKontrakTidakAda, idLampiran)
	}
	if b.IDSimpanan != "" {
		o, adaObjek, err := l.gudang.BacaObjekSimpanan(ctx, b.IDSimpanan)
		if err != nil {
			return PanelLampiran{}, err
		}
		if adaObjek && strings.TrimSpace(o.URLPublik) != "" {
			if l.simpanan == nil {
				return PanelLampiran{}, fmt.Errorf("%w: the storage sender is not wired", ErrSimpananBelumSiap)
			}
			// [6] Namafile = APPFOLDER tanpa `gs://<App>/`.
			jalur := strings.ReplaceAll(o.AppFolder, "gs://"+o.App+"/", "")
			if strings.TrimSpace(jalur) == "" {
				return PanelLampiran{}, fmt.Errorf("%w: APPFOLDER is empty; the storage object cannot be named", ErrSimpananGagal)
			}
			if err := l.simpanan.Hapus(ctx, PermintaanSimpanan{App: o.App, Namafile: jalur}); err != nil {
				return PanelLampiran{}, err
			}
		}
	}
	if err := l.gudang.HapusLampiran(ctx, idKontrak, idLampiran, b.IDSimpanan); err != nil {
		return PanelLampiran{}, err
	}
	return l.susunPanelLampiran(ctx, idKontrak, SifatProporsional(teksDok(kepala, "ProportionType")))
}

// UbahKategoriLampiran - `ChangeDokument_Act("Save")`: kategori baru tiap baris.
func (l *Layanan) UbahKategoriLampiran(ctx context.Context, p inti.Pelaku, idKontrak string, ubah []MasukanUbahKategori) (PanelLampiran, error) {
	if err := inti.WajibIdentitas(p); err != nil {
		return PanelLampiran{}, err
	}
	idKontrak = strings.TrimSpace(idKontrak)
	kepala, err := l.kepalaKontrakLampiran(ctx, idKontrak)
	if err != nil {
		return PanelLampiran{}, err
	}
	// Syarat tampil wadahnya: `StatusAkseptasi != 'Resolve Complete' && != 'Decline'`.
	if s := teksDok(kepala, "StatusAkseptasi"); s == "Resolve Complete" || s == "Decline" {
		return PanelLampiran{}, ditolak("Change Category tidak berlaku untuk kontrak berstatus " + s + ".")
	}
	proporsional := SifatProporsional(teksDok(kepala, "ProportionType"))
	katalog, err := l.gudang.BacaKatalogKategoriLampiran(ctx)
	if err != nil {
		return PanelLampiran{}, err
	}
	daftar, err := l.gudang.BacaLampiranKontrak(ctx, idKontrak)
	if err != nil {
		return PanelLampiran{}, err
	}
	milik := map[string]bool{}
	for _, b := range daftar {
		milik[b.ID] = true
	}
	var tulis []models.PerubahanKategori
	for _, u := range ubah {
		kode := strings.TrimSpace(u.KodeKategori)
		nama, ada := katalog[kode]
		if !ada {
			return PanelLampiran{}, ditolak(fmt.Sprintf("Kategori %q tidak ada di M_KATEGORIMASTERTREATY.", kode))
		}
		if !milik[strings.TrimSpace(u.IDLampiran)] {
			return PanelLampiran{}, ditolak(fmt.Sprintf("Lampiran %q bukan milik kontrak %s.", u.IDLampiran, idKontrak))
		}
		tulis = append(tulis, models.PerubahanKategori{
			IDLampiran: strings.TrimSpace(u.IDLampiran), KodeKategori: kode, NamaKategori: namaKategoriTampil(nama, proporsional),
		})
	}
	if len(tulis) > 0 {
		if err := l.gudang.UbahKategoriLampiran(ctx, idKontrak, tulis); err != nil {
			return PanelLampiran{}, err
		}
	}
	return l.susunPanelLampiran(ctx, idKontrak, proporsional)
}

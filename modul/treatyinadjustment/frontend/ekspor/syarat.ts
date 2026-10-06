// Penilai SYARAT TAMPIL kerangka bangkitan — ditulis tangan, satu per teks.
//
// ⛔ Satu entri per TEKS SYARAT persis seperti di ekspor, bukan pengurai
// ekspresi Pega. Pengurai umum akan "mengerti" syarat yang belum pernah
// dibaca siapa pun; peta ini menolak yang tidak dikenalnya — `layar.test`
// gagal bila kerangka memuat teks yang tidak ada di sini.
//
// ⚠️ Semua syarat dinilai atas halaman AKAR (`TreatyIn.*`) — juga di panel
// Old: Section Old menulis `TreatyIn.FacultativeShare >0`, bukan
// `TreatyIn.OLDDATA.…`. `ViewState` adalah keadaan SESI (Edit `0`, View `1`),
// bukan nilai tersimpan.

type Halaman = Readonly<Record<string, string>>

const lihat = (m: Halaman) => m.ViewState === '1'
/** `When/TreatyMasterInEDM.xml`: `EDMState` 1, 2, atau 3. */
const edm = (m: Halaman) => m.EDMState === '1' || m.EDMState === '2' || m.EDMState === '3'
const angka = (v: string | undefined) => {
  const n = Number(v ?? '')
  return Number.isFinite(n) ? n : 0
}
const belumSelesai = (m: Halaman) => m.StatusAkseptasi !== 'Resolve Complete'

export const PENILAI_SYARAT: Readonly<Record<string, (m: Halaman) => boolean>> = {
  'TreatyIn.FacultativeShare >0': (m) => angka(m.FacultativeShare) > 0,
  'TreatyIn.FacultativeShare>0': (m) => angka(m.FacultativeShare) > 0,
  "TreatyIn.ViewState !='1'": (m) => !lihat(m),
  "TreatyIn.ViewState != '1'": (m) => !lihat(m),
  'TreatyIn.ViewState != 1': (m) => !lihat(m),
  "TreatyIn.ViewState !='1' && TreatyMasterInEDM": (m) => !lihat(m) && edm(m),
  // Properti telanjang sebagai syarat = "benar"; nilainya TEKS di dokumen.
  'TreatyIn.IsMultipleRetro': (m) => m.IsMultipleRetro === 'true',
  "TreatyIn.IsMultipleRetro = 'true'": (m) => m.IsMultipleRetro === 'true',
  "TreatyIn.ReportingPeriod='other'": (m) => m.ReportingPeriod === 'other',
  "TreatyIn.IsEditData !='1'": (m) => m.IsEditData !== '1',
  // ⛔ `facsharedisp` halaman SESI (tombol "Show Facultative Share"), tidak
  // tersimpan di dokumen — syaratnya tidak pernah terpenuhi di layar ini.
  "TreatyIn.FacultativeShare>0 && facsharedisp.CARI1 = '1'": () => false,
  'facsharedisp.CARI1 = 1': () => false,
  'TreatyIn.IsProRate = True': (m) => m.IsProRate === 'true',
  'TreatyIn.IsProRate = true': (m) => m.IsProRate === 'true',
  'TreatyIn.IsProRate == true': (m) => m.IsProRate === 'true',
  'TreatyIn.IsProRate != true': (m) => m.IsProRate !== 'true',
  // `!=` Pega atas kunci yang tidak ada bernilai BENAR — kosong bukan 3.
  'TreatyIn.EDMState != 3 && TreatyIn.EDMMaterialType == 1': (m) => m.EDMState !== '3' && m.EDMMaterialType === '1',
  'TreatyIn.ViewState != 1 ||TreatyIn.RevisionState=1': (m) => !lihat(m) || m.RevisionState === '1',
  "TreatyIn.RevisionState='1'": (m) => m.RevisionState === '1',
  TreatyMasterInEDM: edm,
  '!TreatyMasterInEDM': (m) => !edm(m),
  "(TreatyIn.ViewState !='1' && TreatyIn.StatusAkseptasi != 'Resolve Complete')": (m) => !lihat(m) && belumSelesai(m),
  "TreatyIn.ViewState !='1' && TreatyIn.StatusAkseptasi != 'Resolve Complete'": (m) => !lihat(m) && belumSelesai(m),
}

/**
 * Seluruh syarat sebuah butir terpenuhi?
 *
 * ⛔ Syarat yang TIDAK dikenal → butirnya TIDAK dirender. Menampilkannya
 * berarti menebak bahwa syaratnya benar. Uji menjamin peta ini lengkap
 * untuk kerangka hari ini, jadi cabang ini hanya menangkap bangkitan baru.
 */
export function syaratTerpenuhi(syarat: readonly string[], m: Halaman): boolean {
  return syarat.every((s) => PENILAI_SYARAT[s]?.(m) ?? false)
}

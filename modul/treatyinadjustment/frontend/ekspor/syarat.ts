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
//
// ⭐ Section RINCIAN BARIS (`KERANGKA_RINCIAN`) menulis syarat RELATIF —
// `.SpreadingTypeXOL != ''`, `.TreatyType = 'QUOTA SHARE'` — atas halaman
// baris yang dibuka. Perender menaruh skalar baris itu di halaman yang sama
// dengan kunci BERTITIK (`.TreatyType`), jadi satu peta tetap cukup.

type Halaman = Readonly<Record<string, string>>

const lihat = (m: Halaman) => m.ViewState === '1'
/** `When/TreatyMasterInEDM.xml`: `EDMState` 1, 2, atau 3. */
const edm = (m: Halaman) => m.EDMState === '1' || m.EDMState === '2' || m.EDMState === '3'
const angka = (v: string | undefined) => {
  const n = Number(v ?? '')
  return Number.isFinite(n) ? n : 0
}
const belumSelesai = (m: Halaman) => m.StatusAkseptasi !== 'Resolve Complete'
/** `TreatyIn.IsEditData` — `1` berarti data kontrak dikunci dari suntingan. */
const dataTerkunci = (m: Halaman) => m.IsEditData === '1'
/** `EDMMaterialType` — 1 / 2 (teks pilihannya tidak ada di ekspor). */
const material = (n: '1' | '2') => (m: Halaman) => m.EDMMaterialType === n
/** Properti `= ''` Pega: kunci yang TIDAK ADA juga kosong. */
const kosong = (k: string) => (m: Halaman) => (m[k] ?? '') === ''
/** Skalar halaman BARIS (`.X`) — kunci yang tidak ada = kosong. */
const baris = (k: string) => (m: Halaman) => m[`.${k}`] ?? ''
const salahSatu = (k: string, nilai: readonly string[]) => (m: Halaman) => nilai.includes(baris(k)(m))
const SURPLUS3 = ['SURPLUS', '2ND SURPLUS', '3RD SURPLUS'] as const
const SURPLUS4 = [...SURPLUS3, 'SPECIAL SURPLUS'] as const
const spreadingXOL = (ada: boolean) => (m: Halaman) => (baris('SpreadingTypeXOL')(m) !== '') === ada
const spreadingID = (ada: boolean) => (m: Halaman) => (baris('SpreadingTypeID')(m) !== '') === ada

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

  // ---------------------------------------------------------------------
  // ⭐ 7 Oktober 2026 — syarat BACA-SAJA dan KUNCI panel New
  // (`pyReadOnlyCondition`, `pyDisabledWhen`) serta syarat tampil tombol
  // Add/Delete yang kini dibangkitkan. Semuanya dinilai atas halaman akar.
  // ---------------------------------------------------------------------
  'TreatyIn.ViewState = 1': lihat,
  "TreatyIn.ViewState ='1'": lihat,
  "TreatyIn.IsEditData!='1'": (m) => !dataTerkunci(m),
  'TreatyIn.IsEditData= 1': dataTerkunci,
  "TreatyIn.IsEditData='1'": dataTerkunci,
  "TreatyIn.IsEditData ='1'": dataTerkunci,
  'TreatyIn.EDMMaterialType = 2': material('2'),
  'TreatyIn.ViewState = 1 || TreatyIn.EDMMaterialType = 2': (m) => lihat(m) || material('2')(m),
  "TreatyIn.ViewState =='1' || TreatyIn.EDMMaterialType = 2": (m) => lihat(m) || material('2')(m),
  'TreatyIn.IsEditData= 1 || TreatyIn.EDMMaterialType = 1': (m) => dataTerkunci(m) || material('1')(m),
  'TreatyIn.IsEditData= 1 || TreatyIn.EDMMaterialType = 2': (m) => dataTerkunci(m) || material('2')(m),
  'TreatyIn.ProportionType = "Proportional"': (m) => m.ProportionType === 'Proportional',
  "TreatyIn.ID = ''": kosong('ID'),
  "TreatyIn.EDMEffective = ''": kosong('EDMEffective'),
  // Penjaga yang selalu benar — ditulis di ekspor apa adanya.
  '1=1': () => true,

  // ---------------------------------------------------------------------
  // ⭐ 7 Oktober 2026 — Section RINCIAN BARIS. Ejaan akar yang baru muncul
  // di sana, lalu syarat RELATIF atas halaman baris.
  // ---------------------------------------------------------------------
  "TreatyIn.ViewState = '1'": lihat,
  'TreatyIn.ViewState !=1': (m) => !lihat(m),
  "TreatyIn.ViewState ='1' || TreatyIn.EDMMaterialType = 2": (m) => lihat(m) || material('2')(m),
  "TreatyIn.ViewState = '1' || TreatyIn.EDMMaterialType = 2": (m) => lihat(m) || material('2')(m),
  'TreatyIn.IsEditData = 1': dataTerkunci,
  'TreatyIn.IsEditData = 1 || TreatyIn.EDMMaterialType = 2': (m) => dataTerkunci(m) || material('2')(m),
  // ⭐ `FlagExcel` halaman SESI — diisi `GetAchievement` (rute Treaty In)
  // dan disimpan di akar panel New berkunci utuh (7 Oktober 2026).
  "FlagExcel.CARI1=='1'": (m) => m['FlagExcel.CARI1'] === '1',
  ".TreatyType = 'QUOTA SHARE'": (m) => baris('TreatyType')(m) === 'QUOTA SHARE',
  ".TreatyType = 'SURPLUS' || .TreatyType = '2ND SURPLUS' || .TreatyType = '3RD SURPLUS'": salahSatu('TreatyType', SURPLUS3),
  ".TreatyType = 'SURPLUS' || .TreatyType = '2ND SURPLUS' || .TreatyType = '3RD SURPLUS' || .TreatyType = 'SPECIAL SURPLUS'":
    salahSatu('TreatyType', SURPLUS4),
  // Ditulis ekspor atas `.Note`, bukan `.TreatyType` — dibaca apa adanya.
  ".Note = 'QUOTA SHARE'": (m) => baris('Note')(m) === 'QUOTA SHARE',
  ".Note = 'SURPLUS' || .Note = '2ND SURPLUS' || .Note = '3RD SURPLUS' || .Note = 'SPECIAL SURPLUS'": salahSatu('Note', SURPLUS4),
  ".SpreadingTypeID != ''": spreadingID(true),
  ".SpreadingTypeID == ''": spreadingID(false),
  ".SpreadingTypeID = ''": spreadingID(false),
  ".SpreadingTypeXOL!=''": spreadingXOL(true),
  ".SpreadingTypeXOL !=''": spreadingXOL(true),
  ".SpreadingTypeXOL != ''": spreadingXOL(true),
  ".SpreadingTypeXOL= ''": spreadingXOL(false),
  ".SpreadingTypeXOL =''": spreadingXOL(false),
  ".SpreadingTypeXOL==''": spreadingXOL(false),
  ".SpreadingTypeXOL = ''": spreadingXOL(false),
  // Decimal kosong = 0 di aritmetika Pega (sama dengan `angka`).
  '.Limit2 != 0': (m) => angka(baris('Limit2')(m)) !== 0,

  // ---------------------------------------------------------------------
  // ⭐ 7 Oktober 2026 — SYARAT AKSI (`pyActionConditions`), teks
  // `syarat_aksi` pembangkit. Aksi yang syaratnya tidak terpenuhi
  // DILEWATI, langkah lain rantainya tetap berjalan. Baris sel grid
  // (`.Note`) dinilai atas BARIS sel pemicu.
  // ---------------------------------------------------------------------
  // Kotak centang Share Across The Board — `Associated Property = true`.
  "TreatyIn.RNMShareAcrossTheBoard = 'true'": (m) => m.RNMShareAcrossTheBoard === 'true',
  'TreatyIn.BrokeragePercent > 0': (m) => angka(m.BrokeragePercent) > 0,
  "TreatyIn.FacultativeShare = ''": kosong('FacultativeShare'),
  // ⛔ Tidak ada Material Type 7897987 — langkah EDM Update Summary tidak
  // pernah berjalan, persis Pega (dan rute Treaty In `summary`).
  'TreatyIn.EDMMaterialType = 7897987 && TreatyIn.EDMState != 3': (m) => m.EDMMaterialType === '7897987' && m.EDMState !== '3',
  'TreatyIn.EDMState = 3': (m) => m.EDMState === '3',
  ".TreatyType = 'SURPLUS'": (m) => baris('TreatyType')(m) === 'SURPLUS',
  ".Note = 'SURPLUS' || .Note = '2ND SURPLUS' || .Note = '3RD SURPLUS'": salahSatu('Note', SURPLUS3),
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

// Layar Adjustment — `Section/InputTreatyInAdjustment.xml`.
//
// ⛔ DUA mode, keduanya dari ekspor:
//
//   daftar  `OutputParam.DATASHOW != 1` @481804 — tombol Add Revision /
//           Add Adjustment Premium, penyaring, grid dua belas kolom.
//   detail  `OutputParam.DATASHOW = 1` — kepala (@40422), lalu Old Data dan
//           New Data BERDAMPINGAN (@181268), Attachment (@276129), deret
//           tombol (@336222), History (@380295).
//
// ⛔ Datanya dari `TREATY_IN_EDM` dan tabel pendaratan `T_TREATY_*`. Bukan
// `VERSI_KONTRAK`: diukur nol baris.
//
// ⭐ TOMBOL TULIS HIDUP sejak 7 Oktober 2026 — Save, Submit, Actions,
// Decline offer menulis ke `TREATY_IN_EDM` + `T_TREATY_*` lewat rute
// `/api/treaty-in/penyesuaian/*` (penulisnya di modul Treaty In). Isian
// masuk basis data HANYA dari tombol-tombol itu.
//
// ⚠️ Panel Warning (`TreatyWarning.CARI1 != ''` @317776) TIDAK dirender:
// isinya diisi Activity saat jalan, dan halaman `TreatyWarning` tidak
// tersimpan di dokumen — syaratnya tidak pernah terpenuhi di sini.

import { useCallback, useEffect, useRef, useState } from 'react'

import { Area, Gagal, Halaman, Kosong, Memuat, Modal, Panel } from '../../../../inti/frontend/components/ui/dasar'
import { ambilSesiSaya } from '../../../../inti/frontend/klien'
import {
  ambilDaftarPenyesuaian,
  ambilPenyesuaian,
  hapusPenyesuaian,
  kirimPenyesuaian,
  simpanPenyesuaian,
  type BarisPenyesuaian,
  type BarisRiwayatWarisan,
  type HasilSimpanPenyesuaian,
  type JenisPicker,
  type MasukanSimpanPenyesuaian,
  type Penyesuaian,
  type SisiKiriman,
  type SisiPenyesuaian,
} from '../api'
import type { JenisTulis } from '../komponen/aksiTombol'
import { gabungPohon, samaNilai, teks } from '../komponen/baris'
import { persenLebar } from '../komponen/lebar'
import { gabungPesanSalinan } from '../komponen/salinanLampiran'
import PilihMaster from '../komponen/PilihMaster'
import SisiForm, { selNilai, type ModeLayar } from '../komponen/SisiPenyesuaian'
import {
  JENIS_DAFTAR,
  KOLOM_DAFTAR,
  LEBAR_DAFTAR,
  LEBAR_TOMBOL_BARIS,
  MEDAN_KANAN_BARU,
  MEDAN_KANAN_LAMA,
  MEDAN_KIRI_BARU,
  MEDAN_KIRI_LAMA,
  PENYESUAIAN,
  PILIHAN_AKSEPTASI,
  TAB_BARU_NP,
  TAB_BARU_NP_ADJ,
  TAB_BARU_P,
  TAB_LAMA_NP,
  TAB_LAMA_P,
  TOMBOL,
} from '../labelsPenyesuaian'
import { PanelRiwayat } from './LampiranKontrak'
import PanelLampiranPenyesuaian from '../komponen/PanelLampiranPenyesuaian'
import { teksPromptEDM } from '../labelsPromptEDM'
import PanelPolisMaster, { idMasterPolis } from '../komponen/PanelPolisMaster'

/** Baris per halaman grid daftar — `pyGridPaginator` @651683, ukuran bawaan Pega 10. */
const UKURAN_HALAMAN = 10

/**
 * Nilai grid daftar, urut sesuai `KOLOM_DAFTAR`.
 *
 * Type / Material Type: sel Pega @726743 / @733188 adalah `pxDropdown` yang
 * menampilkan TEKS (`.CARI2` daftar `EDMStates` / `EDMMaterial`), bukan kode.
 * `EDMState 3` tidak ada di rule Property; teksnya mengikuti label kepala
 * "Adjustment Premium" @82379 yang menggantikan dropdown Type untuk kode itu.
 */
export function selDaftar(b: BarisPenyesuaian): string[] {
  const jenis = b.jenisPenyesuaian === '3' ? PENYESUAIAN.premiPenyesuaian : teksPromptEDM('EDMState', b.jenisPenyesuaian)
  return [
    b.id, b.idAsal, jenis,
    teksPromptEDM('EDMMaterialType', b.jenisMaterial), b.namaKontrak, b.sifatProporsi,
    b.asalBisnis, b.cedant, b.tanggalMulai, b.tanggalBerakhir, b.posisi, b.statusAkseptasi,
  ].map((v, i) => selNilai(JENIS_DAFTAR[i] ?? 'teks', v))
}

/**
 * ⭐ Pencarian daftar (8 Oktober 2026). Dicocokkan dengan TEKS YANG TAMPIL
 * (`selDaftar` — tanggal `08-10-2026`, bukan bentuk simpanannya), tanpa
 * membedakan huruf besar/kecil. Beberapa kata = SEMUANYA harus ada, boleh di
 * kolom berbeda: `marsh 2024` menemukan kontrak broker MARSH tahun 2024.
 */
export function cocokCari(b: BarisPenyesuaian, cari: string): boolean {
  const kata = kataCari(cari)
  if (kata.length === 0) return true
  const isi = selDaftar(b).join(' ').toLowerCase()
  return kata.every((k) => isi.includes(k))
}

/** Kata ketikan Search — huruf kecil, tanpa spasi kosong. */
export function kataCari(cari: string): string[] {
  return cari.toLowerCase().split(/\s+/).filter((k) => k !== '')
}

/**
 * Potong teks sel menjadi bagian biasa dan bagian COCOK (disorot `<mark>`),
 * supaya terlihat MENGAPA sebuah baris muncul di hasil pencarian.
 */
export function potongSorot(teks: string, kata: readonly string[]): { isi: string; cocok: boolean }[] {
  if (kata.length === 0 || teks === '') return [{ isi: teks, cocok: false }]
  const kecil = teks.toLowerCase()
  const tanda = new Array<boolean>(teks.length).fill(false)
  for (const k of kata) {
    for (let i = kecil.indexOf(k); i >= 0; i = kecil.indexOf(k, i + 1)) {
      for (let j = i; j < i + k.length; j++) tanda[j] = true
    }
  }
  const hasil: { isi: string; cocok: boolean }[] = []
  for (let i = 0; i < teks.length; i++) {
    const c = tanda[i] === true
    const akhir = hasil[hasil.length - 1]
    if (akhir !== undefined && akhir.cocok === c) akhir.isi += teks[i]
    else hasil.push({ isi: teks[i] ?? '', cocok: c })
  }
  return hasil
}

/**
 * Baris `CommentList` → baris panel History.
 *
 * ⛔ Diukur: `T_VIEW_COMMENT` NOL baris berpengenal penyesuaian, sedang
 * `CommentList` berisi di 280 dari 280 dokumen. Itulah yang harness ikat
 * (`TreatyIn.CommentList` @407001).
 */
export function riwayatDari(sisi: SisiPenyesuaian): BarisRiwayatWarisan[] {
  return (sisi.larik.CommentList ?? []).map((b) => ({
    tanggal: teks(b.Date),
    operator: teks(b.OperatorName),
    disetujui: teks(b.IsApproved),
    catatan: teks(b.Suggest),
  }))
}

/**
 * Syarat TAMPIL tombol Save — @119719:
 * `TreatyIn.ViewState !='1' && TreatyIn.StatusAkseptasi != 'Resolve Complete'`.
 *
 * ⚠️ Kunci `StatusAkseptasi` yang TIDAK ADA bukan `Resolve Complete` — `!=`
 * Pega bernilai benar, jadi tombolnya tampil.
 */
/**
 * Mode layar EFEKTIF panel New — `ViewState` sesudah DT `TreatyInSetEdit`.
 *
 * ⛔ YANG MENENTUKAN ADALAH TOMBOLNYA, BUKAN NILAI TERSIMPAN — ralat atas
 * salah tafsir yang dilaporkan pemilik proses 8 Oktober 2026 (*"kenapa
 * setelah saya coba tidak bisa edit juga"*).
 *
 * `TreatyInSetEdit[2]` berbunyi `WHEN TreatyIn.RevisionState==1 → ViewState
 * = 1`, dan dahulu kami membacanya sebagai kolom tersimpan. Ia bukan:
 * `SetTreatyIn_Act[7]`, satu-satunya yang menyetel properti itu,
 * berprasyarat `param.revisionstate==1` — PARAMETER TOMBOL.
 * `Section/InputTreatyInOffer` memasangkannya:
 *
 *	Edit  → viewstate=<kosong>  revisionstate=<kosong>
 *	View  → viewstate=1         revisionstate=1
 *
 * Keduanya bergerak bersama, jadi `RevisionState` tidak pernah menambahkan
 * apa pun di atas `viewstate`: tombol Edit membuka SUNTING, tombol View
 * membuka BACA. Itu yang `modeAwal` sudah bawa.
 *
 * ⚠️ Kolom `REVISIONSTATE` (migrasi 448) menyimpan hal LAIN — *"kontrak ini
 * masuk jalur revisi tangga akseptasi"*, ditulis tombol Revision layar
 * Treaty In. Memakainya untuk mengunci layar ini membuat setiap penyesuaian
 * turunan kontrak itu lahir terkunci, dan Add Revision tidak menghasilkan
 * apa-apa.
 *
 * ⭐ Dipertahankan sebagai fungsi — bukan diratakan menjadi `modeAwal` —
 * supaya tempat keputusannya tetap satu dan beralasan, dan pagar di
 * `terkunci-pendaratan.test.ts` tetap menunjuk ke sini.
 *
 * ⛔ BERKAS TUNTAS = BACA (9 Oktober 2026). Laporan pemakai: sesudah
 * approve, status `Resolve Complete` tetapi seluruh isian di halaman masih
 * dapat diubah. Di Pega halaman tertutup sesudah akseptasi, dan berkas
 * `Resolve Complete`/`Decline` tidak dapat dibuka sunting lagi: tombol Edit
 * daftar tersembunyi (@782051) dan Save tersembunyi (@119719). Jadi status
 * tuntas/ditolak mengunci layar apa pun tombol pembukanya. (Revision
 * mengosongkan `StatusAkseptasi`, jadi jalurnya tidak terkunci.)
 */
export function modeEfektif(
  modeAwal: ModeLayar,
  _adaDraf: boolean,
  medan: Readonly<Record<string, string>>,
): ModeLayar {
  const status = medan.StatusAkseptasi ?? ''
  if (status === 'Resolve Complete' || status === 'Decline') return '1'
  return modeAwal
}

export function simpanTampil(mode: ModeLayar, medan: Readonly<Record<string, string>>): boolean {
  return mode !== '1' && medan.StatusAkseptasi !== 'Resolve Complete'
}

/**
 * Kepala mode detail — kelima sel ber-`pyReadOnly=true`.
 *
 * ⭐ BENTUK PEGA (9 Oktober 2026, gambar pemakai `1001802/R01`): SATU wadah
 * tanpa judul — kiri lima pasangan label tebal + nilai TEKS bertumpuk (ID
 * Original, ID Revision, Reinsurance Type, Adjustment Type, Material Type),
 * kanan panel `Existing Policy for Master ID` (@111283). Bukan kotak isian,
 * bukan radio: semuanya baca-saja di Pega. Teks pilihan = PROMPT VALUE
 * (`labelsPromptEDM.ts`); `NonProportional` tampil `Non Proportional`.
 */
function Kepala({ p }: { p: Penyesuaian }) {
  const m = p.baru.medan
  const butir = (label: string, nilai: string) => (
    <div className="tria__kepala-butir">
      <dt>{label}</dt>
      <dd>{nilai === '' ? '\u00a0' : nilai}</dd>
    </div>
  )
  const jenis = m.ProportionType === PENYESUAIAN.nonProporsional ? PENYESUAIAN.nonProporsionalTeks : (m.ProportionType ?? '')
  return (
    <section className="panel tria__kepala-pega" aria-label={PENYESUAIAN.judul}>
      <dl className="tria__kepala-data">
        {butir(PENYESUAIAN.idAsal, m.OLDID ?? p.idAsal)}
        {butir(PENYESUAIAN.id, m.ID ?? p.id)}
        {butir(PENYESUAIAN.jenisReasuransi, jenis)}
        {/* `EDMState = 3` → LABEL "Adjustment Premium" @82379 menggantikan
            dropdown @76962 (`EDMState != 3`). */}
        {m.EDMState === '3' ? (
          <div className="tria__kepala-butir">
            <dt>{PENYESUAIAN.premiPenyesuaian}</dt>
          </div>
        ) : (
          butir(PENYESUAIAN.jenisPenyesuaian, teksPromptEDM('EDMState', m.EDMState ?? ''))
        )}
        {butir(PENYESUAIAN.jenisMaterial, teksPromptEDM('EDMMaterialType', m.EDMMaterialType ?? ''))}
      </dl>
      <PanelPolisMaster ringkas idMaster={idMasterPolis(m, p.id, p.idAsal)} />
    </section>
  )
}

/** Posisi penyetuju — `Position` ∈ {SecHead, DeptHead, Director} (@155994). */
const POSISI_PENYETUJU: readonly string[] = ['ReasTreatyInSecHead', 'ReasTreatyInDeptHead', 'ReasTreatyInDirector']

/**
 * Syarat tampil Actions.
 *
 * ⛔ MENYIMPANG DARI EKSPOR atas permintaan pemakai 9 Oktober 2026
 * ("munculkan saja"). Ekspor @155994: hanya pemegang workbasket
 * `TreatyIn.Position` penyetuju DAN `StatusAkseptasi` Accept/Reject — berkas
 * di posisi Admin (belum diajukan) tidak punya Actions sama sekali.
 *
 * Kini Actions SELALU tampil selama berkas belum tuntas (`Resolve Complete`)
 * atau ditolak (`Decline`). Penjaga tetap di server: akun yang tidak
 * memegang workbasket posisi berkas DITOLAK dengan pesan (`KirimPenyesuaian`),
 * dan dari posisi Admin hanya `Accept` yang berlaku (`pilihanAksi`).
 * `POSISI_PENYETUJU` dipertahankan untuk `pilihanAksi`.
 */
export function actionsTampil(medan: Readonly<Record<string, string>>, _workbasket: readonly string[]): boolean {
  const status = medan.StatusAkseptasi ?? ''
  return status !== 'Resolve Complete' && status !== 'Decline'
}

/**
 * Pilihan modal Actions menurut posisi berkas — tangga `Akseptasi_DT`:
 * posisi Admin (atau kosong) HANYA punya cabang `Accept` (= mengajukan ke
 * Sec Head); posisi penyetuju Accept / Reject / Decline.
 */
export function pilihanAksi(posisi: string | undefined): readonly string[] {
  return POSISI_PENYETUJU.includes(posisi ?? '') ? PILIHAN_AKSEPTASI : ['Accept']
}

/**
 * Halaman layar → kiriman tombol tulis.
 *
 * ⛔ Grid Rate of Exchange (`CurrencyList`) cermin `TREATYEXCHANGEYEARLY`,
 * bukan larik dokumen, dan prosedur EDM tidak menulisnya. Ia dikirim HANYA
 * bila berubah — server lalu melaporkannya tak tersimpan, bukan menelannya.
 */
export function sisiKiriman(kini: SisiPenyesuaian, asal: SisiPenyesuaian): SisiKiriman {
  const larik = { ...kini.larik }
  if (samaNilai(larik.CurrencyList ?? [], gabungPohon(asal).larik.CurrencyList ?? [])) delete larik.CurrencyList
  return { medan: kini.medan, larik }
}

/**
 * Deret tombol — blok `TreatyMasterInEDM` @103200 (EDMState 1/2/3; terukur
 * 280 dari 280 penyesuaian memenuhinya).
 *
 * ⭐ HIDUP sejak 7 Oktober 2026:
 *
 *   Save     @119719 → pra-DT `TreatyInAddNew(Status=1)` → `SaveTreatyIn_EDM_Act`
 *   Close    ALWAYS
 *   Actions  @155994 → modal `TreatyInActionEDM` → `TreatyInAkseptasiEDM_Act`
 *
 * ⛔ Tombol yang syarat tampilnya tidak terpenuhi TIDAK dirender — persis
 * Pega. Blok dev (`OperatorID.pyOrgDivision = 'IT'` @167976) tidak dibangun.
 */
function DeretTombol({
  mode,
  medan,
  onTutup,
  onSimpan,
  aksiTampil,
  onAksi,
  sibuk,
}: {
  mode: ModeLayar
  medan: Readonly<Record<string, string>>
  onTutup: () => void
  onSimpan: () => void
  aksiTampil: boolean
  onAksi: () => void
  sibuk: boolean
}) {
  return (
    <div className="tria__aksi" role="group" aria-label={PENYESUAIAN.judul}>
      {simpanTampil(mode, medan) && (
        <button type="button" className="btn btn--primary" disabled={sibuk} onClick={onSimpan}>
          {sibuk ? TOMBOL.menyimpan : TOMBOL.simpan}
        </button>
      )}
      <button type="button" className="btn" onClick={onTutup}>
        {TOMBOL.tutup}
      </button>
      {aksiTampil && (
        <button type="button" className="btn" disabled={sibuk} onClick={onAksi}>
          {TOMBOL.aksi}
        </button>
      )}
    </div>
  )
}

/**
 * Mode detail — satu penyesuaian.
 *
 * `draf` = penyesuaian baru dari `Choose` di picker tombol Add: SUDAH
 * tersusun, BELUM tersimpan, jadi tidak dibaca ulang dari server.
 */
function Detail({
  id,
  draf,
  mode: modeAwal,
  onTutup,
  onTersimpan,
}: {
  id: string
  draf?: Penyesuaian
  mode: ModeLayar
  onTutup: () => void
  /** Draf yang baru TERSIMPAN — dibuka ulang sebagai penyesuaian tersimpan. */
  onTersimpan: (id: string) => void
}) {
  const [p, setP] = useState<Penyesuaian | null>(draf ?? null)
  const [galat, setGalat] = useState<unknown>(null)
  // ⭐ Sesudah tombol tulis berhasil, penyesuaian dibaca ULANG dari tabel —
  // layar memperlihatkan yang TERSIMPAN, bukan yang diketik.
  const [muatUlang, setMuatUlang] = useState(0)
  // Keadaan TERKINI panel New — dilaporkan `SisiForm` (`onKini`).
  const kini = useRef<SisiPenyesuaian | null>(null)
  const catatKini = useCallback((s: SisiPenyesuaian) => {
    kini.current = s
  }, [])
  const [sibuk, setSibuk] = useState(false)
  const [hasil, setHasil] = useState<{ galat: boolean; pesan: string; takTersimpan: string[] } | null>(null)
  // Workbasket pemakai — syarat tampil Actions.
  const [workbasket, setWorkbasket] = useState<string[]>([])
  useEffect(() => {
    let dibuang = false
    ambilSesiSaya()
      .then((s) => {
        if (!dibuang) setWorkbasket(s.peran)
      })
      .catch(() => undefined)
    return () => {
      dibuang = true
    }
  }, [])
  // Modal `TreatyInActionEDM` dan `TreatyInDeclineConfirmationEDM`.
  const [aksiBuka, setAksiBuka] = useState(false)
  const [pilihan, setPilihan] = useState<string>('Accept')
  const [komentar, setKomentar] = useState('')
  const [tolakBuka, setTolakBuka] = useState(false)

  useEffect(() => {
    if (draf !== undefined) {
      setP(draf)
      return
    }
    let dibuang = false
    setP(null)
    setGalat(null)
    // ⭐ Tombol `Edit` daftar @782051 → `SetTreatyInEDM_Act(viewstate=0)` →
    // DT `TreatyInSetEdit` atas halaman yang dimuat — SEKALI, saat dibuka.
    // Muat ulang sesudah Save/Submit TIDAK mereset lagi (di Pega Submit
    // menutup halaman). Belum tersimpan sampai Save.
    const setelEdit = modeAwal === '0' && muatUlang === 0
    Promise.all([ambilPenyesuaian(id), setelEdit ? ambilSesiSaya().catch(() => null) : Promise.resolve(null)])
      .then(([x, sesi]) => {
        if (!dibuang) setP(setelEdit ? terapkanSetEdit(x, sesi?.akunId ?? '') : x)
      })
      .catch((e: unknown) => {
        if (!dibuang) setGalat(e)
      })
    return () => {
      dibuang = true
    }
  }, [id, draf, muatUlang])

  if (galat !== null) return <Gagal galat={galat} />
  if (p === null) return <Memuat />

  // `TreatyInSetEdit`: mode ditentukan TOMBOLNYA. Lihat `modeEfektif`.
  const mode: ModeLayar = modeEfektif(modeAwal, draf !== undefined, p.baru.medan)
  const adalahDraf = draf !== undefined
  /** Isi layar → kiriman. Panel Old ikut HANYA untuk draf. */
  const masukan = (tambahan: Record<string, string> = {}): MasukanSimpanPenyesuaian => {
    const k = kini.current ?? gabungPohon(p.baru)
    return {
      id: p.id,
      draf: adalahDraf,
      baru: sisiKiriman({ ...k, medan: { ...k.medan, ...tambahan } }, p.baru),
      lama: adalahDraf ? { medan: p.lama.medan, larik: gabungPohon(p.lama).larik } : undefined,
    }
  }
  /** Satu penekanan tombol tulis; sesudah berhasil, dibaca ulang dari tabel. */
  const tekan = (jalan: () => Promise<HasilSimpanPenyesuaian>) => {
    setSibuk(true)
    setHasil(null)
    jalan()
      .then((h) => {
        setHasil({ galat: false, pesan: gabungPesanSalinan(h.pesan, h.salinanLampiran), takTersimpan: h.kunciTakTersimpan })
        kini.current = null
        setMuatUlang((n) => n + 1)
        if (adalahDraf) onTersimpan(h.id)
      })
      .catch((e: unknown) => {
        setHasil({ galat: true, pesan: e instanceof Error ? e.message : String(e), takTersimpan: [] })
      })
      .finally(() => {
        setSibuk(false)
      })
  }
  /** Submit / Decline offer tab Information & Submit (`SisiForm` → kerangka). */
  const tulis = (jenis: JenisTulis) => {
    if (jenis === 'submit') {
      tekan(() => kirimPenyesuaian({ ...masukan(), aksi: 'submit' }))
      return
    }
    setTolakBuka(true)
  }
  /** `TreatyInDeclineConfirmation_postactEDM` — baris EDM dihapus fisik. */
  const tolak = () => {
    setTolakBuka(false)
    // Draf belum pernah tersimpan: tidak ada baris untuk dihapus.
    if (adalahDraf) {
      onTutup()
      return
    }
    setSibuk(true)
    setHasil(null)
    hapusPenyesuaian(p.id)
      .then(() => {
        onTutup()
      })
      .catch((e: unknown) => {
        setHasil({ galat: true, pesan: e instanceof Error ? e.message : String(e), takTersimpan: [] })
        setSibuk(false)
      })
  }

  // ⛔ Cabang dibaca SEKALI dari halaman akar, dan KEDUA panel memakainya:
  // `TreatyInNONProportionalOldData.xml` memilih tab lewat
  // `TreatyIn.ProportionType` @278100 @292774 — akar, bukan `OLDDATA`.
  const cabang = p.baru.medan.ProportionType ?? ''
  const np = cabang === 'NonProportional'
  // ⭐ Cabang ADJUST PREMIUM panel New — `EDMState=3` (@566980): Actual GNPI,
  // Actual Limits, Actual Share, Premium Adjustment (halaman `ActualValue`).
  const adjustPremi = np && (p.baru.medan.EDMState ?? '') === '3'
  // ⭐ Isi tab TIDAK ada di pendaratan → panel New dan deret tombolnya dalam
  // mode LIHAT: grid kosong tidak dapat ditambah lalu disimpan sebagai data
  // separuh (audit 8 Oktober 2026). Panel Old tetap menurut `mode`.
  const terkunciPendaratan = p.terdarat === false
  const modeBaru: ModeLayar = terkunciPendaratan ? '1' : mode

  return (
    <>
      <Kepala p={p} />
      {draf !== undefined && (
        <p className="tria__redup" role="note">
          {PENYESUAIAN.drafBelumTersimpan}
        </p>
      )}
      {terkunciPendaratan && mode === '0' && (
        <div className="alert alert--warn" role="note">
          {PENYESUAIAN.belumTerdarat}
        </div>
      )}
      {/* `Existing Policy for Master ID` (@111283) kini DI DALAM `Kepala`,
          di kanan — bentuk Pega (9 Oktober 2026). */}
      {/* ⭐ BERDAMPINGAN: Old di kiri, New di kanan — @181268. Wadahnya
          bersyarat `(NonProportional && DATASHOW=1) || (Proportional &&
          DATASHOW=1)`: kontrak tanpa cabang yang dikenal tidak membukanya. */}
      {cabang === 'NonProportional' || cabang === 'Proportional' ? (
        <div className="tria__bandingan">
          <SisiForm
            key={`lama-${p.id}-${String(muatUlang)}`}
            judul={PENYESUAIAN.panelLama}
            sisi={p.lama}
            akar={p.baru}
            bacaSaja
            mode={mode}
            cabang={cabang}
            medanKiri={MEDAN_KIRI_LAMA}
            medanKanan={MEDAN_KANAN_LAMA}
            bagian={np ? 'TreatyInTabsNonProportionalOldData' : 'TreatyInTabsProportionalOldData'}
            tab={np ? TAB_LAMA_NP : TAB_LAMA_P}
          />
          <SisiForm
            key={`baru-${p.id}-${modeBaru}-${String(muatUlang)}`}
            judul={PENYESUAIAN.panelBaru}
            sisi={p.baru}
            akar={p.baru}
            lama={p.lama}
            bacaSaja={false}
            mode={modeBaru}
            cabang={cabang}
            medanKiri={MEDAN_KIRI_BARU}
            medanKanan={MEDAN_KANAN_BARU}
            bagian={adjustPremi ? 'TreatyInTabsNonProportionalAdjustPremi' : np ? 'TreatyInTabsNonProportional' : 'TreatyInTabsProportional'}
            tab={adjustPremi ? TAB_BARU_NP_ADJ : np ? TAB_BARU_NP : TAB_BARU_P}
            onKini={catatKini}
            tulis={tulis}
            sibukTulis={sibuk}
          />
        </div>
      ) : null}

      {/* Draf: lampiran ASALNYA — `TreatyInEDMSetValue` [8]
          `TreatyRevisionCopyAttachment` menyalin lampiran itu ke pengenal
          baru, dan salinannya menunggu Save. */}
      {/* ⭐ WADAH KEDUA (8 Oktober 2026) — Attachment · deret tombol ·
          History di wadah putih TERSENDIRI di bawah Old/New, bentuk
          Treaty In. Panel Attachment HIDUP: tombol Upload/Delete
          `TreatyIn.ViewState !='1' || TreatyIn.RevisionState='1'`
          (`WorkAttachments`/`ShowAttachmentTreaty`); draf belum ber-ID →
          nol tombol tulis sampai Save. */}
      <div className="tria__inbox-lampiran">
      <PanelLampiranPenyesuaian
        idKontrak={draf !== undefined ? p.idAsal : p.id}
        draf={draf !== undefined}
        jenis={cabang}
        bisaUnggah={mode !== '1' || p.baru.medan.RevisionState === '1'}
        statusAkseptasi={p.baru.medan.StatusAkseptasi ?? ''}
      />
      <DeretTombol
        mode={modeBaru}
        medan={p.baru.medan}
        onTutup={onTutup}
        onSimpan={() => {
          tekan(() => simpanPenyesuaian(masukan()))
        }}
        aksiTampil={actionsTampil(p.baru.medan, workbasket)}
        onAksi={() => {
          setPilihan('Accept')
          setKomentar('')
          setAksiBuka(true)
        }}
        sibuk={sibuk}
      />
      {hasil !== null && (
        <div className={hasil.galat ? 'alert alert--error' : 'alert alert--info'} role={hasil.galat ? 'alert' : 'status'}>
          {hasil.pesan}
          {hasil.takTersimpan.length > 0 && (
            <div className="tria__redup">
              {TOMBOL.takTersimpan} {hasil.takTersimpan.join(', ')}
            </div>
          )}
        </div>
      )}
      {aksiBuka && (
        <Modal
          judul={TOMBOL.judulAksi}
          onTutup={() => {
            setAksiBuka(false)
          }}
          labelBatal={TOMBOL.batal}
          onKirim={() => {
            setAksiBuka(false)
            tekan(() => kirimPenyesuaian({ ...masukan({ Comment: komentar }), aksi: 'akseptasi', pilihan }))
          }}
          aksi={
            <button type="submit" className="btn btn--primary" disabled={sibuk}>
              {TOMBOL.kirim}
            </button>
          }
        >
          <p>
            {TOMBOL.dariID} <strong>{p.id}</strong>
          </p>
          <div className="tria__radio" role="radiogroup" aria-label={TOMBOL.pilihan}>
            {pilihanAksi(p.baru.medan.Position).map((v) => (
              <label key={v}>
                <input
                  type="radio"
                  name="tria-pilihan-akseptasi"
                  value={v}
                  checked={pilihan === v}
                  onChange={() => {
                    setPilihan(v)
                  }}
                />
                {v}
              </label>
            ))}
          </div>
          <Area label={TOMBOL.komentar} value={komentar} onChange={setKomentar} baris={4} />
        </Modal>
      )}
      {tolakBuka && (
        <Modal
          judul={TOMBOL.judulTolak}
          onTutup={() => {
            setTolakBuka(false)
          }}
          labelBatal={TOMBOL.batal}
          onKirim={tolak}
          aksi={
            <button type="submit" className="btn btn--primary">
              {TOMBOL.tolak}
            </button>
          }
        >
          <p>{TOMBOL.tanyaTolak}</p>
        </Modal>
      )}
      <PanelRiwayat riwayat={riwayatDari(p.baru)} petunjuk={PENYESUAIAN.petunjukHistory} />
      </div>
    </>
  )
}

/** Mode daftar — grid `TREATY_IN_EDM`. */
/**
 * Syarat tampil tombol `Edit` daftar — `Section/InputTreatyInAdjustment.xml`
 * @782051: `OperatorID.pyWorkBasketList(2).pyWorkBasketName =
 * 'ReasTreatyInAdmin' && (.Position = 'ReasTreatyInAdmin' || .Position = '')
 * && (.StatusAkseptasi != 'Decline' && .StatusAkseptasi != 'Resolve
 * Complete')`. Workbasket = peran sesi (seperti syarat Actions).
 */
export function editTampil(b: Pick<BarisPenyesuaian, 'kodePosisi' | 'statusAkseptasi'>, workbasket: readonly string[]): boolean {
  const posisi = b.kodePosisi ?? ''
  return (
    workbasket.includes('ReasTreatyInAdmin') &&
    (posisi === 'ReasTreatyInAdmin' || posisi === '') &&
    b.statusAkseptasi !== 'Decline' &&
    b.statusAkseptasi !== 'Resolve Complete'
  )
}

/**
 * Syarat tampil tombol `Add Revision` (@500688) dan `Add Adjustment Premium`
 * (@515563): `OperatorID.pyWorkBasketList(2).pyWorkBasketName =
 * 'ReasTreatyInAdmin'`. Workbasket = peran sesi, seperti `editTampil`.
 */
export function tambahTampil(workbasket: readonly string[]): boolean {
  return workbasket.includes('ReasTreatyInAdmin')
}

/**
 * DT `TreatyInSetEdit` — tombol `Edit` daftar, atas halaman yang DIMUAT
 * (clipboard; tersimpan baru oleh Save):
 *   ViewState = 0 · Position = ReasTreatyInAdmin · IsEditData = 0 ·
 *   PositionUsername = operator · StatusAkseptasi = "" · Comment = "".
 *
 * ⛔ Cabang `RevisionState == 1` TIDAK berlaku di sini: ia menyala dari
 * `param.revisionstate`, dan tombol `Edit` mengirim parameter itu KOSONG
 * (`Section/InputTreatyInOffer`). Lihat `modeEfektif` untuk uraiannya.
 */
export function terapkanSetEdit(p: Penyesuaian, operator: string): Penyesuaian {
  const m = { ...p.baru.medan }
  m.ViewState = '0'
  m.Position = 'ReasTreatyInAdmin'
  m.IsEditData = '0'
  m.PositionUsername = operator
  m.StatusAkseptasi = ''
  m.Comment = ''
  return { ...p, baru: { ...p.baru, medan: m } }
}

function Daftar({ onBuka, onDraf }: { onBuka: (id: string, mode: ModeLayar) => void; onDraf: (p: Penyesuaian) => void }) {
  // Workbasket pemakai — syarat tampil `Edit` (@782051).
  const [workbasket, setWorkbasket] = useState<string[]>([])
  useEffect(() => {
    let dibuang = false
    ambilSesiSaya()
      .then((s) => {
        if (!dibuang) setWorkbasket(s.peran)
      })
      .catch(() => undefined)
    return () => {
      dibuang = true
    }
  }, [])
  const [baris, setBaris] = useState<BarisPenyesuaian[] | null>(null)
  const [galat, setGalat] = useState<unknown>(null)
  const [halaman, setHalaman] = useState(1)
  // Picker yang sedang terbuka — `null` = tertutup.
  const [picker, setPicker] = useState<JenisPicker | null>(null)
  // Ketikan kotak Search — menyaring SELURUH baris (semuanya sudah dimuat),
  // bukan hanya halaman yang tampil.
  const [cari, setCari] = useState('')
  const kotakCari = useRef<HTMLInputElement>(null)

  // Pintasan `/` — langsung mengetik di kotak Search dari mana saja di
  // halaman, kecuali sedang mengetik di isian lain.
  useEffect(() => {
    const tekan = (e: KeyboardEvent) => {
      if (e.key !== '/' || e.ctrlKey || e.metaKey || e.altKey) return
      const t = e.target
      if (t instanceof HTMLElement && (t.isContentEditable || /^(INPUT|TEXTAREA|SELECT)$/.test(t.tagName))) return
      e.preventDefault()
      kotakCari.current?.focus()
    }
    document.addEventListener('keydown', tekan)
    return () => {
      document.removeEventListener('keydown', tekan)
    }
  }, [])

  useEffect(() => {
    let dibuang = false
    ambilDaftarPenyesuaian()
      .then((d) => {
        if (!dibuang) setBaris(d)
      })
      .catch((e: unknown) => {
        if (!dibuang) setGalat(e)
      })
    return () => {
      dibuang = true
    }
  }, [])

  const seluruh = baris ?? []
  const semua = seluruh.filter((b) => cocokCari(b, cari))
  const kata = kataCari(cari)
  const tampil = semua.slice((halaman - 1) * UKURAN_HALAMAN, halaman * UKURAN_HALAMAN)
  const lebar = [...LEBAR_DAFTAR, LEBAR_TOMBOL_BARIS, LEBAR_TOMBOL_BARIS]

  return (
    <Panel judul={PENYESUAIAN.judul}>
      {/* ⭐ Kedua tombol kepala HIDUP sejak 7 Oktober 2026 — keduanya
          membuka picker (`komponen/PilihMaster.tsx`). `Choose` di picker
          menyusun DRAF penyesuaian baru; di Pega ia juga menyimpan, di sini
          simpanannya menunggu tombol Save (keputusan pemilik proses).

          ⭐ Kedua tombol bersyarat tampil `OperatorID.pyWorkBasketList(2).
          pyWorkBasketName = 'ReasTreatyInAdmin'` (@500688 Add Revision,
          @515563 Add Adjustment Premium) — DITERAPKAN 9 Oktober 2026 lewat
          `tambahTampil`, workbasket = peran sesi (seperti tombol Edit
          @782051).

          ⛔ `Show/Hide filter` (sel 112) TIDAK dirender — dicabut 7 Oktober
          2026. Ekspor memberinya `pyVisible = NEVER`, dan ia satu-satunya
          pemanggil `TreatyInShowHide`, DataTransform yang mengisi
          `SearchFilter.CARI1`. Panel penyaring (sel 117-126) hanya tampil bila
          `CARI1 = 1`, jadi di layar Adjustment Pega penyaringnya tidak pernah
          muncul. Tombol mati yang Pega sendiri sembunyikan bukan bagian
          layarnya. */}
      <div className="tria__aksi" role="group" aria-label={PENYESUAIAN.judul}>
        {tambahTampil(workbasket) && (
          <>
            <button
              type="button"
              className="btn btn--primary"
              onClick={() => {
                setPicker('revisi')
              }}
            >
              {PENYESUAIAN.tambahRevisi}
            </button>
            <button
              type="button"
              className="btn btn--primary"
              onClick={() => {
                setPicker('premi')
              }}
            >
              {PENYESUAIAN.tambahPremi}
            </button>
          </>
        )}
        {/* ⭐ Pencarian (8 Oktober 2026) — di KANAN deret tombol. Chip jumlah
            hasil duduk DI KIRI bar, supaya bar tidak bergeser saat chip
            muncul. */}
        {baris !== null && kataCari(cari).length > 0 && (
          <span className={'tria__cari-jumlah' + (semua.length === 0 ? ' tria__cari-jumlah--nol' : '')} aria-live="polite">
            {PENYESUAIAN.cariJumlah(semua.length, seluruh.length)}
          </span>
        )}
        <div className="tria__cari" role="search">
          <svg className="tria__cari-ikon" viewBox="0 0 20 20" aria-hidden="true">
            <circle cx="8.5" cy="8.5" r="5.5" />
            <path d="M12.5 12.5 17 17" />
          </svg>
          <input
            ref={kotakCari}
            type="search"
            className="tria__cari-kotak"
            aria-label={PENYESUAIAN.cari}
            placeholder={PENYESUAIAN.cariPetunjuk}
            value={cari}
            onChange={(e) => {
              setCari(e.target.value)
              setHalaman(1)
            }}
            onKeyDown={(e) => {
              if (e.key === 'Escape' && cari !== '') {
                e.preventDefault()
                setCari('')
                setHalaman(1)
              }
            }}
          />
          {cari !== '' ? (
            <button
              type="button"
              className="tria__cari-bersih"
              aria-label={PENYESUAIAN.cariBersihkan}
              title={PENYESUAIAN.cariBersihkan}
              onClick={() => {
                setCari('')
                setHalaman(1)
                kotakCari.current?.focus()
              }}
            >
              ✕
            </button>
          ) : (
            <kbd className="tria__cari-pintas" title={PENYESUAIAN.cariPintasan}>
              /
            </kbd>
          )}
        </div>
      </div>
      {picker !== null && (
        <PilihMaster
          jenis={picker}
          onTutup={() => {
            setPicker(null)
          }}
          onDraf={(p) => {
            setPicker(null)
            onDraf(p)
          }}
        />
      )}

      {galat !== null && <Gagal galat={galat} />}
      {baris === null && galat === null && <Memuat />}
      {baris !== null && (
        <>
          <div className="table-wrap">
            {/* ⛔ `--daftar`: grid INI punya dua belas kolom data. Lebar
                piksel ekspor menjadi PERBANDINGAN `<col>`, dan tata letak
                `auto` (8 Oktober 2026) menjamin tiap kolom minimal selebar
                kata terpanjangnya — "SAHABAT INSURANCE", `Complete`,
                `NonProportional` tidak pernah terpenggal. */}
            <table className="tria__tabel tria__tabel--daftar">
              <colgroup>
                {lebar.map((_, i) => (
                  <col key={i} style={{ width: persenLebar(lebar, i) }} />
                ))}
              </colgroup>
              <thead>
                <tr>
                  {KOLOM_DAFTAR.map((k) => (
                    <th key={k} scope="col">
                      {k}
                    </th>
                  ))}
                  <th scope="col" aria-label={PENYESUAIAN.ubah} />
                  <th scope="col" aria-label={PENYESUAIAN.lihat} />
                </tr>
              </thead>
              <tbody>
                {semua.length === 0 && seluruh.length === 0 && (
                  <tr>
                    <td colSpan={KOLOM_DAFTAR.length + 2}>
                      <Kosong pesan={PENYESUAIAN.tanpaBaris} petunjuk={PENYESUAIAN.petunjukDaftarKosong} />
                    </td>
                  </tr>
                )}
                {/* Ada data, tetapi nol yang cocok dengan pencarian. */}
                {semua.length === 0 && seluruh.length > 0 && (
                  <tr>
                    <td colSpan={KOLOM_DAFTAR.length + 2} className="tria__cari-kosong">
                      <span>{PENYESUAIAN.cariTanpaHasil(cari.trim())}</span>
                      <button
                        type="button"
                        className="tria__cari-tautan"
                        onClick={() => {
                          setCari('')
                          setHalaman(1)
                          kotakCari.current?.focus()
                        }}
                      >
                        {PENYESUAIAN.cariBersihkan}
                      </button>
                    </td>
                  </tr>
                )}
                {tampil.map((b) => (
                  <tr key={b.id}>
                    {selDaftar(b).map((v, i) => (
                      // ⛔ Tanggal diberi kelasnya di SEL, bukan di `<col>`:
                      // `white-space` tidak berlaku pada `<col>` — elemen itu
                      // hanya menghormati `width`, `background`, `border`, dan
                      // `visibility`. Aturan yang ditaruh di sana diam-diam
                      // tidak berlaku, dan `08-10-2026` tetap pecah dua baris.
                      <td key={i} className={JENIS_DAFTAR[i] === 'tanggal' ? 'tria__sel-tanggal' : undefined}>
                        {/* ⭐ Bagian yang cocok dengan pencarian disorot. */}
                        {potongSorot(v, kata).map((p, j) =>
                          p.cocok ? (
                            <mark key={j} className="tria__sorot">
                              {p.isi}
                            </mark>
                          ) : (
                            p.isi
                          ),
                        )}
                      </td>
                    ))}
                    {/* `Edit` @782051 dan `View` @798870 sama-sama MEMBUKA —
                        bedanya `ViewState` 0 lawan 1. Membuka bukan menulis;
                        yang menulis tombol Save, dan ia mati. */}
                    <td>
                      {/* `Edit` @782051 — syarat tampil `editTampil`. */}
                      {editTampil(b, workbasket) && (
                      <button
                        type="button"
                        className="btn btn--ghost btn--sm"
                        onClick={() => {
                          onBuka(b.id, '0')
                        }}
                      >
                        {PENYESUAIAN.ubah}
                      </button>
                      )}
                    </td>
                    <td>
                      <button
                        type="button"
                        className="btn btn--ghost btn--sm"
                        onClick={() => {
                          onBuka(b.id, '1')
                        }}
                      >
                        {PENYESUAIAN.lihat}
                      </button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
          <Halaman halaman={halaman} ukuran={UKURAN_HALAMAN} total={semua.length} onPindah={setHalaman} />
        </>
      )}
    </Panel>
  )
}

/**
 * Mode layar draf dari `Choose` — `ViewState` yang `TreatyInSetEdit` setel:
 * `0` (Edit), atau `1` bila dokumennya ber-`RevisionState = 1`.
 */
export function modeDraf(p: Penyesuaian): ModeLayar {
  return p.baru.medan.ViewState === '1' ? '1' : '0'
}

export default function PenyesuaianKontrak() {
  const [buka, setBuka] = useState<{ id: string; mode: ModeLayar; draf?: Penyesuaian } | null>(null)
  return (
    <div className="inbox">
      <header className="inbox__kepala">
        <h2 className="inbox__judul">{PENYESUAIAN.judulMenu}</h2>
        {buka !== null && (
          <button
            type="button"
            className="btn btn--ghost"
            onClick={() => {
              setBuka(null)
            }}
          >
            {PENYESUAIAN.kembali}
          </button>
        )}
      </header>
      {buka === null ? (
        <Daftar
          onBuka={(id, mode) => {
            setBuka({ id, mode })
          }}
          onDraf={(p) => {
            setBuka({ id: p.id, mode: modeDraf(p), draf: p })
          }}
        />
      ) : (
        <Detail
          id={buka.id}
          draf={buka.draf}
          mode={buka.mode}
          onTutup={() => {
            setBuka(null)
          }}
          onTersimpan={(idTersimpan) => {
            setBuka({ id: idTersimpan, mode: buka.mode })
          }}
        />
      )}
    </div>
  )
}

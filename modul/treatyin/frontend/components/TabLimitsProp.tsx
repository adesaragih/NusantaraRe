// Tab **Limits** cabang PROPORSIONAL — tiga tingkat yang dapat dibuka,
// dibaca dari ekspor (`labelsLimitsProp.ts`).
//
//   Kind of Treaty  (grid, `Limits[]`)          expandPane → LimitProportional
//   └ Treaty Group  (grid, `Limits[].Detail[]`) expandPane → DetailLimits
//     └ medan · grid · sebelas tab
//
// ⛔ MENGGANTIKAN `PohonLimits` untuk cabang ini saja. Pohon lama memasang
// medan LAYER non-prop (Cover, MDP, ROL …) di bawah `Kind of Treaty`; Section
// prop tidak punya satu pun, sehingga layar berisi deretan "—". Cabang
// non-prop tetap memakai `PohonLimits`.
//
// ⭐ Mode `ubah`: seluruh medan dapat disunting dan tiap grid ber-`Add`
// punya `Add`/`Delete` yang bekerja — atas salinan pohon di layar ini.
// Mode `lihat`: baca-saja, tombol tidak dirender (`TreatyIn.ViewState !='1'`).

import { createContext, useContext, useEffect, useState, type ReactNode } from 'react'

import { Field, FieldAngka, Kosong } from '../../../../inti/frontend/components/ui/dasar'
import { PilihCari as Pilih } from '../../../../inti/frontend/components/ui/pilihSaring'
import { formatNumber } from '../../../../inti/frontend/lib/format'
import {
  ambilKelasBisnis,
  catatLogAchievement,
  ambilOpsiLimits,
  hitungCadangan,
  hitungDeduksi,
  hitungAchievement,
  hitungLimit,
  type OpsiLimits,
  type PilihanWarisan,
  type SimpulLimit,
} from '../api'
import type { GolonganAngka } from '../labels'
import { unduhXlsx } from '../../../../inti/frontend/lib/exportXlsx'
import {
  ACHIEVEMENT,
  GRID_DETAIL,
  KEPALA_KOSONG,
  KOLOM_ACH_PARAMETER,
  KOLOM_KIND_OF_TREATY,
  KOLOM_TREATY_GROUP,
  LIMITS_PROP,
  TAB_DETAIL_LIMIT,
  TAB_KUNCI_MATERIAL,
  syaratQS,
  syaratSurplus,
  type KolomLimit,
  type MedanLimit,
} from '../labelsLimitsProp'
import { useProperti } from '../halaman'
import type { ModeForm } from '../mode'
import { DropdownDaftar } from './IsianAuto'
import { TombolHapus, TombolTambah } from './limitsUI'
import { BlokPega, GridPega } from './gridPega'
import { selAngka } from './angka'
import { StripTabNavigasi } from './navigasi'
import { PemicuUbah } from './pemicuUbah'
import { saringAngka } from './saringAngka'
import { SelKosongPega, TataPegaBlok } from './tataPega'

/** Teks satu medan simpul — larik dan kunci yang tidak ada → kosong. */
export function teksDari(s: SimpulLimit, kunci: string): string {
  const v = s[kunci]
  return typeof v === 'string' ? v : ''
}

/** Larik di simpul — yang tidak ada → kosong. */
function larikDari(s: SimpulLimit, kunci: string): SimpulLimit[] {
  const v = s[kunci]
  return Array.isArray(v) ? v : []
}

/**
 * Format satu nilai menurut golongan + desimal ekspor.
 *
 * ⛔ Pemformatnya `selAngka` modul (memanggil `format.ts` inti) — bukan
 * pemformat kedua. `-1` = `pyDecimalPlaces = -999`, tak dibatasi.
 */
export function formatLimit(golongan: GolonganAngka, desimal: number | null, nilai: string): string {
  if (desimal === -1) return formatNumber(nilai, -1)
  return selAngka(desimal === null ? golongan : [golongan, desimal], nilai)
}

/** Isi dropdown tab Limits — diberikan sekali di akar, dibaca tiap tingkat. */
const OPSI_KOSONG: OpsiLimits = {
  jenisTreaty: [],
  kelompokTreaty: [],
  mataUang: [],
}
const OpsiLimitsCtx = createContext<OpsiLimits>(OPSI_KOSONG)

/**
 * Pohon `Limits[]` TERKINI — diberikan di akar, dibaca `DetailLimits`.
 *
 * ⛔ Perlu sebab langkah 9 `LimitCalculation` menapaki SELURUH
 * `TreatyIn.Limits` untuk menemukan Quota Share bergrup sama. Simpul
 * `Detail` sendirian tidak tahu saudaranya.
 */
const PohonLimitsCtx = createContext<readonly SimpulLimit[]>([])

/**
 * `EDMMaterialType = 2` — kunci materialitas. Sel yang ber-`pyDisabledWhen
 * TreatyIn.EDMMaterialType = 2` dimatikan walau mode Edit.
 */
const KunciMaterialCtx = createContext(false)

/** Keadaan sub-tab Achievement — `SearchData` dan `TempQuarter*` Pega bersifat GLOBAL. */
interface KeadaanAchievement {
  asAt: string
  tahun: string
  kuartal: readonly string[]
  tahunKuartal: readonly string[]
  flagExcel: boolean
  pesan: string
  /** As At berubah → DT `Reset_DT` (Quarter Year dikosongkan). */
  pilihAsAt: (v: string) => void
  /** Quarter Year berubah → `GetAchievement(search)`. */
  pilihTahun: (v: string) => void
  /** Refresh → `GetAchievement`. */
  segarkan: () => void
  /** `TreatyIn.ID` — `InsertToLogAchievement` [1] `InputParam.CARI1`. */
  idKontrak: string
}
const AchievementCtx = createContext<KeadaanAchievement | null>(null)

/** Sub-tab Achievement satu Treaty Group — `DetailLimits` @2372416. */
function PanelAchievement({ d, grid }: { d: SimpulLimit; grid: (larik: string) => ReactNode }) {
  const a = useContext(AchievementCtx)
  const daftar = larikDari(d, 'AchievementLists')
  // `InsertToLogAchievement` → `LOG_ACHIEVEMENT` (keputusan pemakai
  // 8 Oktober 2026): satu baris log per baris AchievementLists.
  const [log, setLog] = useState<{ sibuk: boolean; pesan: string; galat: boolean }>({ sibuk: false, pesan: '', galat: false })
  const kirim = () => {
    const kolom = ['Quarter', 'QUARTERYEAR', 'CurrencyID', 'Currency', 'PREMIUM', 'RICOMM', 'BROKERAGE', 'NETPREMIUM',
      'PaidClaim', 'CASHCALL', 'OutstandingClaim', 'IncuredClaim', 'Total', 'LossRatio']
    setLog({ sibuk: true, pesan: '', galat: false })
    catatLogAchievement({
      idKontrak: a?.idKontrak ?? '',
      baris: daftar.map((b) => Object.fromEntries(kolom.map((k) => [k, teksDari(b, k)]))),
    })
      .then((h) => {
        setLog({ sibuk: false, pesan: ACHIEVEMENT.tercatat(h.disisipkan, h.dilewati), galat: false })
      })
      .catch((e: unknown) => {
        setLog({ sibuk: false, pesan: e instanceof Error ? e.message : String(e), galat: true })
      })
  }
  const unduh = () => {
    // `GenerateCSVTreaty`: 13 kolom AchievementLists, baris tanpa Quarter
    // (baris total) dilewati.
    const kolom = GRID_DETAIL.AchievementLists?.kolom ?? []
    unduhXlsx(
      ACHIEVEMENT.berkasExcel,
      kolom.map((k) => ({ kunci: k.kunci, label: k.label })),
      daftar.filter((b) => teksDari(b, 'Quarter') !== '').map((b) => Object.fromEntries(kolom.map((k) => [k.kunci, teksDari(b, k.kunci)]))),
    )
  }
  return (
    <>
      {grid('AchievementLists')}
      <GridLimit kolom={KOLOM_ACH_PARAMETER} baris={larikDari(d, 'CurrencyList')} bisaUbah={false} tambah={false} bacaSaja onUbah={() => undefined} />
      {a !== null && (
        // ⭐ TATA LETAK PEGA (8 Oktober 2026) — `Section/DetailLimits.xml`
        // `Inline labels left` L98928: As At Quarter dan Quarter Year SEBARIS,
        // label masing-masing di kiri kotaknya.
        <TataPegaBlok tata="alir">
          <TataPegaBlok tata="kiri">
            <Pilih
              label={ACHIEVEMENT.asAt}
              value={a.asAt}
              kosong={LIMITS_PROP.pilihKosong}
              opsi={a.kuartal.map((q) => ({ value: q, label: q }))}
              onChange={a.pilihAsAt}
            />
          </TataPegaBlok>
          <TataPegaBlok tata="kiri">
            <Pilih
              label={ACHIEVEMENT.tahun}
              value={a.tahun}
              kosong={LIMITS_PROP.pilihKosong}
              opsi={a.tahunKuartal.map((q) => ({ value: q, label: q }))}
              onChange={a.pilihTahun}
            />
          </TataPegaBlok>
        </TataPegaBlok>
      )}
      {/* `Inline` L100440 — Refresh · Excel · Insert Log sebaris. */}
      <div className="trin__aksi">
        <button type="button" className="btn" onClick={a?.segarkan}>
          {ACHIEVEMENT.refresh}
        </button>
        {a?.flagExcel === true && (
          <>
            <button type="button" className="btn" onClick={unduh}>
              {ACHIEVEMENT.excel}
            </button>
            {/* `InsertToLogAchievement` → `LOG_ACHIEVEMENT` — HIDUP sejak
                8 Oktober 2026 (keputusan pemakai). */}
            <button type="button" className="btn" disabled={log.sibuk} onClick={kirim}>
              {ACHIEVEMENT.kirim}
            </button>
          </>
        )}
      </div>
      {a !== null && a.pesan !== '' && (
        <p className="trin__galat" role="alert">
          {a.pesan}
        </p>
      )}
      {log.pesan !== '' && (
        <p className={log.galat ? 'trin__galat' : 'trin__redup'} role={log.galat ? 'alert' : 'status'}>
          {log.pesan}
        </p>
      )}
    </>
  )
}

/**
 * Medan yang `LimitCalculation` TULIS, per mode — hanya ini yang
 * digabungkan kembali ke simpul.
 *
 * ⛔ Menggabungkan seluruh jawaban akan melahirkan kunci yang dokumennya
 * tidak punya (`IOOPct` kosong pada Quota Share, misalnya), dan layar
 * membedakan "kunci tidak ada" dari "kunci kosong".
 */
export const DITULIS_LIMIT_CALCULATION = {
  // [2]–[4]
  qs: ['RetentionPct', 'CessionPct', 'RetentionList', 'CessionList'],
  // [6]–[14]
  surplus: ['IOOPct', 'CessionPct', 'IOOLimitList', 'RetentionList', 'CessionList', 'COBList'],
} as const

/** Pilihan untuk `Pilih`; nama kembar diberi kodenya. */
function opsiDari(daftar: readonly PilihanWarisan[], nilai: 'id' | 'nama') {
  const lihat = new Set<string>()
  const keluar: { value: string; label: string }[] = []
  for (const o of daftar) {
    const v = nilai === 'id' ? o.id : o.nama
    if (lihat.has(v)) continue
    lihat.add(v)
    keluar.push({
      value: v,
      label: nilai === 'id' && o.kembar ? `${o.nama} (${o.id} — ${LIMITS_PROP.namaKembar})` : o.nama,
    })
  }
  return keluar
}

/**
 * Dropdown terikat kode (`.TreatyTypeID`, `.TreatyGroupID`, `.CurrencyID`):
 * mode ubah `Pilih`; mode lihat nama baca-saja, seperti dropdown read-only Pega.
 * Memilih mengisi nama pasangannya — `Set…Name_Act` di Pega.
 */
function PilihKode({
  label,
  daftar,
  simpul,
  kunciID,
  kunciNama,
  bisaUbah,
  onUbah,
  namaDari = (o) => o.nama,
  ketik = false,
}: {
  label: string
  daftar: readonly PilihanWarisan[]
  simpul: SimpulLimit
  kunciID: string
  kunciNama: string
  bisaUbah: boolean
  onUbah: (s: SimpulLimit) => void
  /** Nilai yang ditulis ke `kunciNama` dari pilihan — bawaan: Name-nya. */
  namaDari?: (o: PilihanWarisan) => string
  /**
   * Dropdown yang DAPAT DIKETIK untuk menyaring — komponen yang SAMA dengan
   * pemilih Ceding (`DropdownDaftar` → `DropdownWarisan`): ketik-saring,
   * papan tik, nama kembar diberi pengenalnya. Nilainya tetap HANYA dari
   * daftar, sama seperti `pxDropdown`.
   */
  ketik?: boolean
}) {
  // ⭐ PENCARIAN BALIK: pengenal dari NAMA, ketika pengenalnya tidak tersimpan.
  //
  // Terukur: `TreatyTypeID` terisi pada 19 dari 1.360 elemen `Limits[]`,
  // sementara `TreatyType` (namanya) hampir selalu ada. Dropdown mengikat
  // pengenal, jadi tanpa pencarian ini 98,6% layer berbunyi `Choose` padahal
  // namanya tertulis di kepala simpul tepat di atasnya.
  //
  // ⛔ BUKAN KARANGAN — Pega memakai pola yang sama, dan rulenya ada di
  // korpus: `Claim Non Prop/RDBList/GetReinsuranceTypeBYName_SQL.xml`
  // mencari `ID` dari `NOTE` pada `REINSURANCETYPE`.
  //
  // ⚠️ Hanya untuk nama TUNGGAL. Nama kembar membuat pencarian memilih salah
  // satu, dan layer akan menunjuk jenis treaty yang KELIRU tanpa ada yang
  // tahu — `Kembar` dihitung di `repository.tandaiKembar`.
  //
  // ⛔ Hasilnya TIDAK DITULIS ke dokumen. Ia hanya memilihkan di layar; yang
  // menulis tetap perbuatan pemakai, lewat `onChange` di bawah. Menulis
  // diam-diam akan mengubah 1.341 layer tanpa seorang pun memintanya.
  const kodeTersimpan = teksDari(simpul, kunciID)
  const namaTersimpan = teksDari(simpul, kunciNama)
  const kode =
    kodeTersimpan !== ''
      ? kodeTersimpan
      : (cariKodeDariNama(daftar, namaTersimpan) ?? '')
  if (!bisaUbah) {
    const nama = daftar.find((o) => o.id === kode)?.nama ?? teksDari(simpul, kunciNama)
    return <Field label={label} value={nama !== '' ? nama : kode} readOnly onChange={() => undefined} />
  }
  if (ketik) {
    return (
      <DropdownDaftar
        label={label}
        // Kotaknya memperlihatkan label PILIHAN (Name) — yang daftar
        // tampilkan — bukan `kunciNama`, yang bisa nama turunan (SOA Name).
        nilai={daftar.find((o) => o.id === kode)?.nama ?? ''}
        pilihan={daftar}
        bisaUbah
        onPilih={(_, id) => {
          const o = daftar.find((x) => x.id === id)
          onUbah({ ...simpul, [kunciID]: id, [kunciNama]: o === undefined ? '' : namaDari(o) })
        }}
      />
    )
  }
  return (
    <Pilih
      label={label}
      value={kode}
      kosong={LIMITS_PROP.pilihKosong}
      opsi={opsiDari(daftar, 'id')}
      onChange={(v) => {
        onUbah({
          ...simpul,
          [kunciID]: v,
          [kunciNama]: ((o) => (o === undefined ? '' : namaDari(o)))(daftar.find((o) => o.id === v)),
        })
      }}
    />
  )
}

/**
 * Pengenal satu pilihan dari NAMANYA — hanya bila namanya TUNGGAL.
 *
 * ⛔ Mengembalikan `null` pada nama kembar, nama kosong, dan nama yang tidak
 * ada di daftar. Ketiganya berarti hal yang sama bagi layar: jangan menebak.
 */
export function cariKodeDariNama(
  daftar: readonly PilihanWarisan[],
  nama: string,
): string | null {
  const n = nama.trim()
  if (n === '') return null
  const cocok = daftar.filter((o) => o.nama.trim() === n)
  if (cocok.length !== 1) return null
  const satu = cocok[0]
  if (satu === undefined || satu.kembar) return null
  return satu.id
}

/**
 * Kind of Treaty dari baris menu Reinsurance Type yang dipilih sebagai
 * Treaty Type — kolom **SOA Name** (`SOANOTE`).
 *
 * ⭐ KEPUTUSAN PEMAKAI 6 Oktober 2026: "treaty type diambil dari kolom name
 * dari menu itu dan kind of treaty nya itu SOA name nya". Menyimpang dari
 * `SetTreatyTypeName_Act` (yang menulis `.Note`).
 *
 * ⚠️ SOA Name KOSONG (20 dari 84 jenis aktif, termasuk `QUOTA SHARE` 10035
 * dan `SURPLUS` 10042) → kolom **Name**, supaya Kind of Treaty tidak
 * kosong. Lihat `PERTANYAAN-TERBUKA-LAYAR-PEGA.md` §20.
 */
export function kindOfTreatyDari(o: PilihanWarisan): string {
  const soa = (o.namaSoa ?? '').trim()
  return soa !== '' ? soa : o.nama
}

/** Ganti satu simpul di larik, tanpa mengubah larik aslinya. */
function ganti(larik: readonly SimpulLimit[], i: number, baru: SimpulLimit): SimpulLimit[] {
  return larik.map((x, j) => (j === i ? baru : x))
}

/** Satu medan simpul: teks baca-saja (lihat) atau kotak isian (ubah). */
/**
 * Pembungkus `Stacked with labels left` satu sel — label di KIRI medan
 * (`.trin__tata--kiri > .field`). Tanpa `aktif`, anaknya dirender apa adanya.
 */
function BungkusKiri({ aktif, children }: { aktif: boolean; children: ReactNode }) {
  return aktif ? <div className="trin__tata trin__tata--kiri">{children}</div> : <>{children}</>
}

function MedanSimpul({
  m,
  simpul,
  bisaUbah,
  onUbah,
  tataKiri = false,
}: {
  m: MedanLimit
  simpul: SimpulLimit
  bisaUbah: boolean
  onUbah: (s: SimpulLimit) => void
  /**
   * ⭐ TATA LETAK PEGA (8 Oktober 2026) — medan ini duduk di layout `Stacked
   * with labels left` (`Section/DetailLimits.xml`): label di KIRI kotaknya.
   * Medan bermata uang (Event Limits) menjadi pasangan `Inline 30 70 table`
   * L27339 — [mata uang berlabel, 30%] | [nilai, 70%].
   */
  tataKiri?: boolean
}) {
  const opsi = useContext(OpsiLimitsCtx)
  const ada = m.kunci in simpul
  const nilai = teksDari(simpul, m.kunci)
  const kelas = !tataKiri ? 'trin__limit-medan' : m.mataUang !== undefined ? 'trin__tata trin__tata--t3070' : 'trin__tata trin__tata--kiri'
  return (
    <div className={kelas}>
      {/* `.Currency…` pxDropdown — `BrowseCurrencyTreatyIn_RD`, nilai `.Currency`. */}
      {m.mataUang !== undefined && (
        <BungkusKiri aktif={tataKiri}>
          {bisaUbah ? (
            <Pilih
              label={m.label}
              value={teksDari(simpul, m.mataUang)}
              kosong={LIMITS_PROP.pilihKosong}
              opsi={opsiDari(opsi.mataUang, 'nama')}
              onChange={(v) => {
                onUbah({ ...simpul, [m.mataUang ?? '']: v })
              }}
            />
          ) : (
            <Field label={m.label} value={teksDari(simpul, m.mataUang)} readOnly onChange={() => undefined} />
          )}
        </BungkusKiri>
      )}
      {/* ⭐ PEMISAH RIBUAN SAAT MENGETIK — hanya untuk kolom yang BENAR
          angka DAN presisinya terbaca di ekspor (`pyDecimalPlaces`).

          ⛔ DUA GOLONGAN SENGAJA DILEWATI, dan keduanya pernah menggigit:
            • `teks` — `>=30% up to < 50%` kebetulan berangka; menguraikannya
              sebagai bilangan merusaknya.
            • `desimal === null` — presisinya TIDAK dinyatakan ekspor. Memaksa
              dua desimal membuat nomor urut tampil `1,00`; pemilik proses
              melaporkannya tepat begitu di kartu Layers tab Limits Non-Prop. */}
      {m.golongan === 'teks' || m.desimal === null ? (
        <Field
          label={m.mataUang !== undefined ? '' : m.label}
          // ⛔ Mode ubah: nilai TERSIMPAN apa adanya di kotak; mode lihat:
          // terjemahan tampil. Terjemahan untuk membaca, bukan untuk mengisi.
          value={bisaUbah ? nilai : formatLimit(m.golongan, m.desimal, nilai)}
          readOnly={!bisaUbah}
          onChange={(v) => {
            // Kontrol Number ditolak hurufnya; `teks` (mis. `>=30% up to
            // < 50%`) lewat apa adanya. `Periode` (Period (Month)) bergolongan
            // teks di sini tetapi kontrol Number di `Section/DetailLimits.xml` @2170296.
            onUbah({ ...simpul, [m.kunci]: m.golongan !== 'teks' || m.kunci === 'Periode' ? saringAngka(v) : v })
          }}
        />
      ) : (
        <FieldAngka
          label={m.mataUang !== undefined ? '' : m.label}
          value={nilai}
          desimal={m.desimal}
          persen={m.golongan !== 'uang'}
          readOnly={!bisaUbah}
          onChange={(v) => {
            onUbah({ ...simpul, [m.kunci]: v })
          }}
        />
      )}
      {!ada && !bisaUbah && <span className="trin__redup">{LIMITS_PROP.tidakAda}</span>}
    </div>
  )
}

/** Peristiwa sel grid yang memicu Activity di Pega. */
export interface PeristiwaGrid {
  jenis: 'mata-uang' | 'nilai' | 'hapus'
  /** Larik SESUDAH perubahan. */
  baris: SimpulLimit[]
  /** Indeks baris (mulai 0). */
  r: number
  /** Baris SEBELUM perubahan — `.Note` dan `.Layer` dibaca darinya. */
  lama: SimpulLimit
  /** Kunci sel yang berubah. */
  kunci: string
}

/**
 * Grid satu larik — kolom dari ekspor; `Add`/`Remove` hidup di mode ubah.
 *
 * ⛔ Kolom yang syarat SEL-nya tidak terpenuhi untuk sebuah baris tetap
 * berdiri — selnya yang kosong. Begitu Pega merender syarat sel.
 */
export function GridLimit({
  kolom,
  baris,
  bisaUbah,
  tambah,
  onUbah,
  judul,
  bacaSaja = false,
  barisBaru,
  kelasBisnis = [],
  onPeristiwa,
}: {
  kolom: readonly KolomLimit[]
  baris: readonly SimpulLimit[]
  bisaUbah: boolean
  tambah: boolean
  onUbah: (baris: SimpulLimit[]) => void
  judul?: ReactNode
  /** Grid baca-saja di ekspor (`readOnly`) — tak dapat disunting di mode mana pun. */
  bacaSaja?: boolean
  /** Nilai awal baris baru — DT/Activity tombol `Add` di Pega. */
  barisBaru?: SimpulLimit
  /** Pilihan autocomplete `Class of Business` grup ini. */
  kelasBisnis?: readonly PilihanWarisan[]
  /** Pemicu Activity per sel (mata uang, nilai, hapus baris). */
  onPeristiwa?: (p: PeristiwaGrid) => void
}) {
  const opsi = useContext(OpsiLimitsCtx)
  const total = kolom.reduce((a, k) => a + k.lebar, 0) || 1
  const sunting = bisaUbah && !bacaSaja
  const pakaiTombol = sunting && tambah
  const lapor = (p: PeristiwaGrid) => {
    if (onPeristiwa !== undefined) onPeristiwa(p)
  }
  return (
    <div className="trin__limit-grid">
      {(judul !== undefined || pakaiTombol) && (
        <div className="tl-grid-kepala">
          <span className="trin__limit-label">{judul ?? ''}</span>
          {pakaiTombol && (
            <TombolTambah
              label={LIMITS_PROP.tambah}
              onClick={() => {
                onUbah([...baris, { ...(barisBaru ?? {}) }])
              }}
            />
          )}
        </div>
      )}
      <div className="table-wrap">
        <table className="trin__tabel trin__tabel--pega">
          <colgroup>
            {kolom.map((k, i) => (
              <col key={i} style={{ width: `${((k.lebar / total) * 100).toFixed(2)}%` }} />
            ))}
            {pakaiTombol && <col className="tl-kolom-aksi" />}
          </colgroup>
          <thead>
            <tr>
              {kolom.map((k, i) => (
                <th key={i} scope="col">
                  {/* Kepala kosong di ekspor diberi nama kolomnya — grid tanpa
                      kepala tidak terbaca, mata maupun pembaca layar. */}
                  {k.label !== '' ? k.label : (KEPALA_KOSONG[k.kunci] ?? '')}
                </th>
              ))}
              {pakaiTombol && <th scope="col" aria-label={LIMITS_PROP.hapusBaris} />}
            </tr>
          </thead>
          <tbody>
            {baris.length === 0 && (
              <tr>
                <td colSpan={kolom.length + (pakaiTombol ? 1 : 0)}>
                  <Kosong pesan={LIMITS_PROP.tanpaBaris} />
                </td>
              </tr>
            )}
            {baris.map((b, r) => (
              <tr key={r}>
                {kolom.map((k, i) => {
                  if (k.kunci === '' || (k.syarat !== undefined && !k.syarat(b))) return <td key={i} />
                  const v = teksDari(b, k.kunci)
                  const ubahSel = (baru: SimpulLimit, jenis?: PeristiwaGrid['jenis']) => {
                    const semua = ganti(baris, r, baru)
                    onUbah(semua)
                    if (jenis !== undefined) lapor({ jenis, baris: semua, r, lama: b, kunci: k.kunci })
                  }
                  if (k.cek === true) {
                    // `.Layer` — `pxCheckbox` (Layer = saklar hitung otomatis),
                    // keterangan di sampingnya (`pyCheckboxCaption`).
                    const kotak = (
                      <input
                        type="checkbox"
                        aria-label={k.keterangan ?? (k.label || k.kunci)}
                        checked={v === 'true'}
                        disabled={!sunting}
                        onChange={(e) => {
                          ubahSel({ ...b, [k.kunci]: e.target.checked ? 'true' : 'false' })
                        }}
                      />
                    )
                    return (
                      <td key={i}>
                        {k.keterangan === undefined ? (
                          kotak
                        ) : (
                          <label className="trin__cek">
                            {kotak}
                            {k.keterangan}
                          </label>
                        )}
                      </td>
                    )
                  }
                  return (
                    <td key={i} className={k.golongan === 'teks' ? undefined : 'trin__angka'}>
                      {sunting && k.auto === 'kelasBisnis' ? (
                        <DropdownDaftar
                          label=""
                          nilai={v}
                          pilihan={kelasBisnis}
                          bisaUbah
                          onPilih={(nama, id) => {
                            // DT `SetCoBID` — nama dan kode Class of Business.
                            ubahSel({ ...b, [k.kunci]: nama, ClassOfBusinessID: id })
                          }}
                        />
                      ) : sunting && k.mataUang === 'id' ? (
                        <PilihKode
                          label=""
                          daftar={opsi.mataUang}
                          simpul={b}
                          kunciID="CurrencyID"
                          kunciNama={k.kunci}
                          bisaUbah
                          onUbah={(x) => {
                            ubahSel(x, 'mata-uang')
                          }}
                        />
                      ) : sunting && k.mataUang === 'nama' ? (
                        <Pilih
                          label=""
                          value={v}
                          kosong={LIMITS_PROP.pilihKosong}
                          opsi={opsiDari(opsi.mataUang, 'nama')}
                          onChange={(x) => {
                            ubahSel({ ...b, [k.kunci]: x }, 'mata-uang')
                          }}
                        />
                      ) : sunting && k.golongan !== 'teks' && k.desimal !== null ? (
                        /* ⭐ SEL ANGKA BERPEMISAH RIBUAN — laporan pemilik
                           proses 7 Oktober 2026: grid `100% Limit`,
                           `Retention`, dan `Cession to R/I` menampilkan
                           `1500000000` apa adanya.

                           ⛔ Pembungkusnya menangkap `onBlur`: `FieldAngka`
                           memakai `onBlur`-nya sendiri untuk keluar dari
                           mode ketik, sementara peristiwa `change` Pega
                           harus tetap dilaporkan. `onBlur` React
                           MENGGELEMBUNG, jadi keduanya hidup berdampingan.

                           ⚠️ `k.desimal === null` TIDAK ikut ke sini —
                           presisinya tidak dinyatakan ekspor, dan memaksa
                           dua desimal membuat nomor urut tampil `1,00`.
                           Pelajaran dari kartu Layers tab Limits Non-Prop. */
                        <PemicuUbah
                          nilai={v}
                          aksi={() => {
                            lapor({ jenis: 'nilai', baris: [...baris], r, lama: b, kunci: k.kunci })
                          }}
                        >
                          <FieldAngka
                            label=""
                            value={v}
                            desimal={k.desimal}
                            persen={k.golongan !== 'uang'}
                            onChange={(x) => {
                              ubahSel({ ...b, [k.kunci]: x })
                            }}
                          />
                        </PemicuUbah>
                      ) : sunting ? (
                        // Peristiwa `change` Pega — sekali, dan hanya bila nilainya berubah.
                        <PemicuUbah
                          nilai={v}
                          aksi={() => {
                            lapor({ jenis: 'nilai', baris: [...baris], r, lama: b, kunci: k.kunci })
                          }}
                        >
                          <input
                            className="field__input"
                            type="text"
                            value={v}
                            aria-label={k.label || k.kunci}
                            onChange={(e) => {
                              ubahSel({ ...b, [k.kunci]: k.golongan !== 'teks' ? saringAngka(e.target.value) : e.target.value })
                            }}
                          />
                        </PemicuUbah>
                      ) : (
                        formatLimit(k.golongan, k.desimal, v)
                      )}
                    </td>
                  )
                })}
                {pakaiTombol && (
                  <td>
                    <TombolHapus
                      label={LIMITS_PROP.hapusBaris}
                      labelAkses={`${LIMITS_PROP.hapusBaris} ${judul ?? ''} ${r + 1}`.replace(/\s+/g, ' ')}
                      onClick={() => {
                        const sisa = baris.filter((_, j) => j !== r)
                        onUbah(sisa)
                        lapor({ jenis: 'hapus', baris: sisa, r, lama: b, kunci: '' })
                      }}
                    />
                  </td>
                )}
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  )
}

/** Judul grid nilai + persennya (baca-saja) sebagai lencana. */
/**
 * Persen di samping judul grid — bentuk Pega: NILAI, spasi, lalu LABEL `%`
 * tersendiri (`.RetentionPct` @474451 + `%`, `.CessionPct` @656819 + `%`),
 * terbaca `0 %` / `100 %` di tangkapan layar pemakai 7 Oktober 2026. Angka
 * apa adanya, tanpa pembulatan.
 */
export function persenSamping(nilai: string): string {
  return nilai.trim() === '' ? LIMITS_PROP.persen : `${formatNumber(nilai, -1)} ${LIMITS_PROP.persen}`
}

/**
 * Retention % dan Cession % Quota Share — `LimitCalculation` [3], satu-satunya
 * penulis keduanya untuk QUOTA SHARE:
 *
 *   `.RetentionPct = 100 - .QSPct`   `.CessionPct = 100 - .RetentionPct`
 *
 * ⭐ Dihitung DARI QS % yang sedang tampil, bukan menunggu jawaban server —
 * keluhan pemakai 7 Oktober 2026: "yang lainnya malah hilang angkanya".
 * Treaty Group baru belum pernah menjalankan rumusnya, jadi kedua properti
 * itu kosong. Hasilnya sama persis dengan nilai yang `LimitCalculation`
 * tulis sesudah QS % berubah, dan diukur 100% cocok atas 1.468 detail QS
 * tersimpan (`hitung_limit.go`).
 *
 * QS % kosong → nilai tersimpan apa adanya (Pega pun belum menghitungnya).
 * Desimalnya mengikuti QS %; `toFixed` menjinakkan galat float
 * (`100 - 33.33`).
 */
export function persenQuotaShare(d: SimpulLimit): { retensi: string; cession: string } {
  const qs = teksDari(d, 'QSPct').trim()
  const n = Number(qs)
  if (qs === '' || !Number.isFinite(n)) {
    return { retensi: teksDari(d, 'RetentionPct'), cession: teksDari(d, 'CessionPct') }
  }
  const desimal = (qs.split('.')[1] ?? '').length
  const retensi = (100 - n).toFixed(desimal)
  return { retensi, cession: (100 - Number(retensi)).toFixed(desimal) }
}

function JudulNilai({ judul, persen }: { judul: string; persen: string }) {
  return (
    <span className="tl-grid-judul">
      {judul}
      {persen !== '' && <span className="tl-persen">{persen}</span>}
    </span>
  )
}

/** `.Note` baris retensi yang memicu `LimitCalculation(surplus)` saat MATA UANG berubah. */
const SURPLUS_MATA_UANG = ['SURPLUS', '2ND SURPLUS', '3RD SURPLUS']

/**
 * Sub-tab DetailLimits yang medannya duduk LANGSUNG di layout `Stacked with
 * labels left` — label di kiri kotak (`Section/DetailLimits.xml`: Event
 * Limits L26873 → pasangan `Inline 30 70 table` L27339, Deduction In A
 * L35922, Experience Premium Refund L51491). `% Premium Reserve` (sel
 * `Inline grid double` L44861) dan medan LPC (`Inline` L72437) TIDAK.
 */
const SUB_LABEL_KIRI: ReadonlySet<string> = new Set(['Event Limits', 'Deduction In A', 'Experience Premium Refund'])

/**
 * ⭐ TATA LETAK PEGA (8 Oktober 2026) — susunan isi satu sub-tab DetailLimits
 * menurut `pyLayoutOtherFormat` `Section/DetailLimits.xml`. Urutan `butir`
 * = urutan `TAB_DETAIL_LIMIT[...].isi` (`labelsLimitsProp.ts`).
 */
function susunSubTabLimit(judul: string, butir: ReactNode[]): ReactNode {
  switch (judul) {
    case 'Reserve':
      // `Stacked with labels left` L44563: [`Inline grid double` L44861 —
      // % Premium Reserve (label di ATAS) | kosong], lalu `Inline grid 30 70`
      // L45874 — [label "Premium Reserve" 30% | grid 70%] (gambar Pega 08).
      return (
        <TataPegaBlok tata="kiri">
          <TataPegaBlok tata="g2">
            {butir[0]}
            <SelKosongPega />
          </TataPegaBlok>
          <TataPegaBlok tata="t3070">{butir.slice(1)}</TataPegaBlok>
        </TataPegaBlok>
      )
    case 'LPC':
      // Layout bertumpuk L71989: label "Loss Participation Clause (LPC)",
      // lalu `Inline` L72437 — empat medan sebaris.
      return (
        <TataPegaBlok tata="tumpuk">
          {butir[0]}
          <TataPegaBlok tata="alir">{butir.slice(1)}</TataPegaBlok>
        </TataPegaBlok>
      )
    case 'Deduction':
      // Layout bertumpuk L37172: grid Deduction Details lalu grid total.
      return <TataPegaBlok tata="tumpuk">{butir}</TataPegaBlok>
    case 'PLA':
    case 'Cash Loss Limit':
    case 'Claim Cooperation':
    case 'EPI':
      // `Stacked with labels left` (L53773, L59848, L65923, L74566) →
      // `Inline grid 30 70` (L54239, L60314, L66389, L75032): [label 30% |
      // grid 70%], bentuk yang sama dengan Premium Reserve di gambar Pega 08.
      return (
        <TataPegaBlok tata="kiri">
          <TataPegaBlok tata="t3070">{butir}</TataPegaBlok>
        </TataPegaBlok>
      )
    default:
      // Event Limits L26873, Deduction In A L35922, Experience Premium
      // Refund L51491, Achievement L80637 — `Stacked with labels left`.
      return <TataPegaBlok tata="kiri">{butir}</TataPegaBlok>
  }
}

/** Tingkat 3 — `Section/DetailLimits.xml`. */
/**
 * ⭐ PERBAIKAN 7 Oktober 2026 — "input tiba-tiba hilang". Rumus berjalan
 * ASINKRON (rute services). Dahulu jawabannya ditimpakan sebagai simpul UTUH
 * yang disalin saat rumus dipicu — isian yang diketik selama rumus berjalan
 * ikut tertimpa nilai lamanya. Kini jawaban diterapkan lewat `terapkan`
 * sebagai FUNGSI atas simpul TERKINI di penampung halaman, dan hanya medan
 * yang Activity tulis yang diganti — seperti refresh Pega atas clipboard.
 */
type Terapkan = (f: (kini: SimpulLimit) => SimpulLimit) => void

function DetailLimits({
  d,
  bisaUbah,
  onUbah,
  terapkan,
}: {
  d: SimpulLimit
  bisaUbah: boolean
  onUbah: (d: SimpulLimit) => void
  /** Terapkan ubahan atas Detail TERKINI (bukan salinan render). */
  terapkan?: Terapkan
}) {
  // Tanpa jalur terkini (pemakaian lama/uji): jatuh ke `onUbah` atas `d`.
  const terapkanKini: Terapkan = terapkan ?? ((f) => onUbah(f(d)))
  const [tab, setTab] = useState<string>(TAB_DETAIL_LIMIT[0]?.judul ?? '')
  const [gagalHitung, setGagalHitung] = useState('')
  const [pesan, setPesan] = useState<string[]>([])
  const [kelasBisnis, setKelasBisnis] = useState<PilihanWarisan[]>([])
  const opsiLimits = useContext(OpsiLimitsCtx)
  const pohon = useContext(PohonLimitsCtx)
  // `pyDisabledWhen TreatyIn.EDMMaterialType = 2` — kunci materialitas.
  const kunciMaterial = useContext(KunciMaterialCtx)
  const bisaUbahM = bisaUbah && !kunciMaterial
  const jenis = teksDari(d, 'TreatyType')
  const qs = syaratQS(jenis)
  const idGrup = teksDari(d, 'TreatyGroupID')

  // Autocomplete `Class of Business` — `BrowseTreatyBusinessWOType_RD`
  // berparameter `pTreatyGroupId = .TreatyGroupID`.
  useEffect(() => {
    if (!bisaUbah || idGrup === '') return
    let dibuang = false
    ambilKelasBisnis(idGrup)
      .then((k) => {
        if (!dibuang) setKelasBisnis(k)
      })
      .catch(() => undefined)
    return () => {
      dibuang = true
    }
  }, [bisaUbah, idGrup])

  const gagal = (e: unknown) => {
    setGagalHitung(e instanceof Error ? e.message : String(e))
  }

  /**
   * ⭐ `LimitCalculation` atas simpul `dBaru`, dengan parameter persis sel
   * pemicunya. Hanya medan yang Activity TULIS yang digabung kembali.
   *
   * ⛔ `autocalculate = false` — langkah 1 Activity KELUAR tanpa mengubah
   * apa pun, jadi tidak ada yang dikirim.
   *
   * ⛔ Hasilnya TAMPIL saja: ia mengganti simpul di salinan pohon layar ini,
   * dan jalur Save menunggu keputusan pemilik proses.
   */
  const jalankanLC = (dBaru: SimpulLimit, mode: 'qs' | 'surplus', tambah: '' | 'man', otomatis: boolean) => {
    if (!bisaUbah || !otomatis) return
    setGagalHitung('')
    const semua = pohon.flatMap((l) => larikDari(l, 'Detail'))
    hitungLimit({ jenis: mode, tambah, otomatis: true, detail: dBaru, pohon: semua })
      .then((h) => {
        // Hanya medan yang `LimitCalculation` tulis — atas Detail TERKINI.
        terapkanKini((kini) => {
          const baru: SimpulLimit = { ...kini }
          for (const k of DITULIS_LIMIT_CALCULATION[mode]) {
            const v = h[k]
            if (v !== undefined) baru[k] = v
          }
          return baru
        })
      })
      .catch(gagal)
  }

  /**
   * Sel 37 `.QSPct` / sel 42 `.Surplus` — `LimitCalculation(qs|surplus, '',
   * true)`, dipicu saat medan DILEPAS (peristiwa `change`, bukan tiap ketukan).
   */
  const hitung = (mode: 'qs' | 'surplus') => {
    jalankanLC(d, mode, '', true)
  }

  /**
   * Pemicu grid 100% Limit (sel 75–77) dan Retention (sel 117–119):
   * `autocalculate = .Layer` baris itu — kotak centang `Layer` adalah
   * saklar hitung otomatisnya.
   *
   * ⚠️ Sel 75 mengirim `kindoftreaty = 'QS'` HURUF BESAR ke perbandingan
   * `=="qs"` — salah alamat yang menjatuhkan Quota Share ke cabang
   * surplus dan menggandakan 100% Limit-nya. Di sini ia dibaca `qs`.
   */
  const peristiwaIOO = (p: PeristiwaGrid) => {
    if (teksDari(p.lama, 'Note') !== 'QUOTA SHARE') return
    jalankanLC({ ...d, IOOLimitList: p.baris }, 'qs', '', teksDari(p.lama, 'Layer') === 'true')
  }
  const peristiwaRetensi = (p: PeristiwaGrid) => {
    const note = teksDari(p.lama, 'Note')
    const syarat =
      p.jenis === 'hapus' ||
      (p.jenis === 'mata-uang' && SURPLUS_MATA_UANG.includes(note)) ||
      (p.jenis === 'nilai' && syaratSurplus(note))
    if (!syarat) return
    jalankanLC({ ...d, RetentionList: p.baris }, 'surplus', 'man', teksDari(p.lama, 'Layer') === 'true')
  }

  /** `CalculateDeduction(sts, index)` — mata uang/hapus: '', nilai: val, persen: pct. */
  const peristiwaDeduksi = (p: PeristiwaGrid) => {
    if (!bisaUbah) return
    const sts = p.jenis === 'nilai' ? (p.kunci === 'DeductionPct' ? 'pct' : p.kunci === 'Deduction' ? 'val' : null) : ''
    if (sts === null) return
    setGagalHitung('')
    hitungDeduksi({ sts, indeks: p.r, DeductionList: p.baris, GrossPremiumList: larikDari(d, 'GrossPremiumList') })
      .then((h) => {
        setPesan(h.pesan)
        terapkanKini((kini) => ({ ...kini, DeductionList: h.DeductionList, DeductionTotalList: h.DeductionTotalList }))
      })
      .catch(gagal)
  }

  /** `PremiumReserveCalculate` — saat `% Premium Reserve` dilepas. */
  const hitungCadanganPremi = () => {
    if (!bisaUbahM) return
    setGagalHitung('')
    hitungCadangan({ PremiumReservePct: teksDari(d, 'PremiumReservePct'), CessionList: larikDari(d, 'CessionList') })
      .then((h) => {
        setPesan(h.pesan)
        terapkanKini((kini) => ({ ...kini, ReserveList: h.ReserveList }))
      })
      .catch(gagal)
  }

  /** Nilai awal baris `Add` per larik — DT/Activity tombolnya di Pega. */
  const barisBaru: Record<string, SimpulLimit> = {
    // DT `AddLimitRetentionCession`: `.Note = .TreatyType`.
    IOOLimitList: { Note: jenis },
    RetentionList: { Note: jenis },
    CessionList: { Note: jenis },
    // `AddDeduction`: `Currency = .CurrencyIOOLimit`.
    DeductionList: { Currency: teksDari(d, 'CurrencyIOOLimit') },
    // `AddValue(type)`: `ID = ""`.
    ReserveList: { ID: '' },
    PLAList: { ID: '' },
    CashLossList: { ID: '' },
    ClaimCoopList: { ID: '' },
    EPIList: { ID: '' },
  }
  const pemicu: Record<string, (p: PeristiwaGrid) => void> = {
    IOOLimitList: peristiwaIOO,
    RetentionList: peristiwaRetensi,
    DeductionList: peristiwaDeduksi,
  }

  const grid = (larik: string, judul: ReactNode, kunci: boolean) => {
    const g = GRID_DETAIL[larik]
    if (g === undefined) return null
    return (
      <GridLimit
        key={larik}
        judul={judul}
        kolom={g.kolom}
        baris={larikDari(d, larik)}
        bisaUbah={kunci ? bisaUbahM : bisaUbah}
        tambah={g.tambah}
        bacaSaja={g.bacaSaja === true}
        barisBaru={barisBaru[larik]}
        kelasBisnis={kelasBisnis}
        onPeristiwa={pemicu[larik]}
        onUbah={(b) => {
          onUbah({ ...d, [larik]: b })
        }}
      />
    )
  }
  const medan = (m: MedanLimit, ubah = bisaUbah, kiri = false) => (
    <MedanSimpul key={m.kunci} m={m} simpul={d} bisaUbah={ubah} onUbah={onUbah} tataKiri={kiri} />
  )
  const aktif = TAB_DETAIL_LIMIT.find((t) => t.judul === tab) ?? TAB_DETAIL_LIMIT[0]
  const tabTerkunci = aktif !== undefined && TAB_KUNCI_MATERIAL.includes(aktif.judul)
  return (
    <>
      {/* `.TreatyGroupID` pxDropdown — `BrowseTreatyGroup_RD` (TREATYGROUP):
          nilai `ID`, label `.TreatyGroupName`; berubah → `SetTreatyGroupName_Act`
          lalu `LimitCalculation(surplus, '', true)` bila jenisnya SURPLUS.
          ⚠️ `FetchQSfromMaster` (bila QUOTA SHARE) mengisi larik spreading
          yang TAMPIL di tab Share — dibangun bersama tab Share, bukan di sini. */}
      {/* ⭐ TATA LETAK PEGA (8 Oktober 2026) — `Section/DetailLimits.xml`
          `Inline grid double` L644, tiga sel: Treaty Group (`Stacked with
          labels left` L946) | sel `Spacer` `1=2` (L1879) — TETAP memesan
          slotnya | grid Class of Business (L2823). Jadi CoB turun ke baris
          kedua, separuh KIRI — persis gambar Pega 05/06/08. */}
      <TataPegaBlok tata="g2">
        <div className="trin__tata trin__tata--kiri">
          <PilihKode
            label={LIMITS_PROP.treatyGroup}
            daftar={opsiLimits.kelompokTreaty}
            simpul={d}
            kunciID="TreatyGroupID"
            kunciNama="TreatyGroup"
            bisaUbah={bisaUbahM}
            onUbah={(x) => {
              onUbah(x)
              if (jenis === 'SURPLUS') jalankanLC(x, 'surplus', '', true)
            }}
          />
        </div>
        <SelKosongPega />
        {/* Class of Business — autocomplete; Add CoB/Delete MATI di ekspor. */}
        {grid('COBList', undefined, true)}
      </TataPegaBlok>

      {/* ⭐ TATA LETAK PEGA (8 Oktober 2026) — `Section/DetailLimits.xml`
          L4956, `Stacked with labels left` L5178: QS % dan Lines masing-masing
          `Inline grid double` (L5665, L7103) — medannya di separuh KIRI,
          label di kiri kotak; lalu `Inline grid 30 70` L8544. */}
      <TataPegaBlok tata="kiri">
      {(qs || syaratSurplus(jenis)) && (
        <>
          {qs && (
            // Sel 37 `.QSPct` — `LimitCalculation(qs, '', true)` pada `change`.
            <TataPegaBlok tata="g2">
            <PemicuUbah
              nilai={teksDari(d, 'QSPct')}
              aksi={() => {
                hitung('qs')
              }}
            >
              {medan(
                {
                  label: LIMITS_PROP.qs,
                  kunci: 'QSPct',
                  golongan: 'persen',
                  desimal: 2,
                },
                bisaUbah,
                true,
              )}
            </PemicuUbah>
            <SelKosongPega />
            </TataPegaBlok>
          )}
          {syaratSurplus(jenis) && (
            // Sel 42 `.Surplus` — `LimitCalculation(surplus, '', true)` pada `change`.
            <TataPegaBlok tata="g2">
            <PemicuUbah
              nilai={teksDari(d, 'Surplus')}
              aksi={() => {
                hitung('surplus')
              }}
            >
              {medan(
                {
                  label: LIMITS_PROP.lines,
                  kunci: 'Surplus',
                  golongan: 'uang',
                  desimal: 2,
                },
                bisaUbah,
                true,
              )}
            </PemicuUbah>
            <SelKosongPega />
            </TataPegaBlok>
          )}
        </>
      )}
      {gagalHitung !== '' && (
        <p className="tl-pesan" role="alert">
          {gagalHitung}
        </p>
      )}
      {pesan.length > 0 && (
        <ul className="tl-pesan" role="alert">
          {pesan.map((p, i) => (
            <li key={i}>{p}</li>
          ))}
        </ul>
      )}

        {/* Tiap grid: judul + persennya BERDAMPINGAN sebagai teks baca-saja
            (`100 %`, `0 %`, `100 %` — permintaan pemakai 7 Oktober 2026,
            persis tangkapan layar Pega).
            `.RetentionPct` / `.CessionPct` Read-only SELALU di ekspor —
            ditulis `LimitCalculation`. Tampil bila QUOTA SHARE (@297544,
            @474451, @656819).

            ⭐ TATA LETAK PEGA (8 Oktober 2026) — `Inline grid 30 70` L8544
            = TIGA BARIS [judul 30% | grid 70%], BUKAN tiga grid berdampingan:
            gambar Pega 05/06/08 menaruh "100% Limit", "Retention", "Cession
            to R/I" di kiri dan gridnya mulai ±30% lebar. Sel judul =
            layout L8565/L14519/L20620, sel grid = L10815/L16834/L22935.
            Tombol Add tetap di kepala gridnya. */}
        <TataPegaBlok tata="t3070">
          <JudulNilai judul={LIMITS_PROP.limit100} persen={qs ? `${LIMITS_PROP.seratus} ${LIMITS_PROP.persen}` : ''} />
          {grid('IOOLimitList', undefined, false)}
          <JudulNilai judul={LIMITS_PROP.retensi} persen={qs ? persenSamping(persenQuotaShare(d).retensi) : ''} />
          {grid('RetentionList', undefined, false)}
          <JudulNilai judul={LIMITS_PROP.cession} persen={qs ? persenSamping(persenQuotaShare(d).cession) : ''} />
          {grid('CessionList', undefined, false)}
        </TataPegaBlok>
      </TataPegaBlok>

      <StripTabNavigasi tab={TAB_DETAIL_LIMIT.map((t) => t.judul)} aktif={aktif?.judul ?? ''} onPilih={setTab} />
      {/* `[BLOK tab] X` → `[BLOK] X`: judul blok = judul sub-tab, seperti ekspor. */}
      <BlokPega judul={aktif?.judul ?? ''}>
        {susunSubTabLimit(
          aktif?.judul ?? '',
          (aktif?.isi ?? []).map((b, i) =>
          b.t === 'medan' ? (
            b.m.kunci === 'PremiumReservePct' ? (
              <PemicuUbah key={b.m.kunci} nilai={teksDari(d, 'PremiumReservePct')} aksi={hitungCadanganPremi}>
                {medan(b.m, bisaUbahM)}
              </PemicuUbah>
            ) : (
              medan(b.m, tabTerkunci ? bisaUbahM : bisaUbah, SUB_LABEL_KIRI.has(aktif?.judul ?? ''))
            )
          ) : b.t === 'grid' && b.larik === 'AchievementLists' ? (
            <PanelAchievement key={b.larik} d={d} grid={(larik) => grid(larik, undefined, false)} />
          ) : b.t === 'grid' ? (
            grid(b.larik, b.judul, tabTerkunci)
          ) : (
            <span key={i} className="trin__limit-label">
              {b.teks}
            </span>
          ),
          ),
        )}
      </BlokPega>
    </>
  )
}

/**
 * DT `TreatyTypeSetIndex` — dijalankan `SetTreatyTypeName_Act` sesudah
 * Treaty Type berubah (`.TreatyTypeID` `change` → refresh activity):
 *
 *   1    `.ID = .pxListSubscript`
 *   2    FOR_EACH_PAGE_IN `.Detail`
 *   2.1    `.TreatyType = .TreatyType` — jenis INDUK ke tiap Treaty Group
 *
 * ⭐ Tanpa langkah 2, Treaty Group yang ditambahkan SEBELUM Treaty Type
 * dipilih (atau sebelum ia diganti) tetap membawa jenis lamanya — dan
 * medan QS % (`.TreatyType = 'QUOTA SHARE'` di simpul Treaty Group) tidak
 * pernah tampil. `AddClassofBusiness` menyalin jenis yang sama saat
 * menambah baris; langkah ini menjaganya tetap sama sesudahnya.
 */
export function terapkanJenisTreaty(l: SimpulLimit, indeks: number): SimpulLimit {
  const jenis = teksDari(l, 'TreatyType')
  const baru: SimpulLimit = { ...l, ID: String(indeks + 1) }
  if ('Detail' in l) baru.Detail = larikDari(l, 'Detail').map((d) => ({ ...d, TreatyType: jenis }))
  return baru
}

/** Tingkat 2 — `Section/LimitProportional.xml`. */
function LimitProportional({
  l,
  indeks,
  bisaUbah,
  onUbah,
  terapkan,
}: {
  l: SimpulLimit
  /** `.pxListSubscript` Kind of Treaty (mulai 0) — `ParentID` Treaty Group. */
  indeks: number
  bisaUbah: boolean
  onUbah: (l: SimpulLimit) => void
  /** Terapkan ubahan atas Kind of Treaty TERKINI — lihat `DetailLimits`. */
  terapkan?: Terapkan
}) {
  const detail = larikDari(l, 'Detail')
  const { jenisTreaty } = useContext(OpsiLimitsCtx)
  // ⭐ Bentuk Pega (8 Oktober 2026) — `Section/LimitProportional.xml`:
  // Treaty Type, lalu grid `masterDetail` Treaty Group (lebar 1368, kolom
  // tombol 83) yang barisnya dibuka menjadi `DetailLimits`. Kartu berchip
  // dan kotak `Kind of Treaty` tambahan DICABUT — Kind of Treaty sudah
  // tampil di baris grid induknya.
  return (
    <>
      {/* ⭐ Treaty Type — pilihan dari menu Reinsurance Type (`REINSURANCETYPE`,
          Flag `active`, urut ID menurun — RD `BrowseReinsuranceType_RD`):
          nilai `ID`, label kolom **Name**. Memilih mengisi Kind of Treaty
          dengan kolom **SOA Name** baris yang sama (`kindOfTreatyDari`).
          Baca-saja hanya bila `TreatyIn.ViewState = 1`. */}
      <PilihKode
        label={LIMITS_PROP.treatyType}
        daftar={jenisTreaty}
        simpul={l}
        kunciID="TreatyTypeID"
        kunciNama="TreatyType"
        bisaUbah={bisaUbah}
        onUbah={(x) => {
          // `SetTreatyTypeName_Act` → DT `TreatyTypeSetIndex`.
          onUbah(terapkanJenisTreaty(x, indeks))
        }}
        namaDari={kindOfTreatyDari}
        // ⭐ Dapat diketik seperti Ceding — permintaan pemakai 7 Oktober 2026.
        ketik
      />
      <GridPega
        label={KOLOM_TREATY_GROUP[0]?.label}
        kolom={[
          {
            judul: KOLOM_TREATY_GROUP[0]?.label ?? '',
            lebar: 1368,
            isi: (d) => (teksDari(d, 'TreatyGroup') === '' ? LIMITS_PROP.belumDipilih : teksDari(d, 'TreatyGroup')),
          },
        ]}
        baris={detail}
        tombol={
          bisaUbah
            ? {
                lebar: 83,
                tambah: {
                  label: LIMITS_PROP.tambah,
                  onKlik: () => {
                    // `AddClassofBusiness`: ID kosong, ParentID = subscript induk,
                    // TreatyType = jenis induk.
                    onUbah({
                      ...l,
                      Detail: [...detail, { ID: '', ParentID: String(indeks + 1), TreatyType: teksDari(l, 'TreatyType') }],
                    })
                  },
                },
                hapus: {
                  label: LIMITS_PROP.hapus,
                  akses: (_, i) => `${LIMITS_PROP.hapus} ${KOLOM_TREATY_GROUP[0]?.label ?? ''} ${String(i + 1)}`,
                  onKlik: (i) => {
                    onUbah({ ...l, Detail: detail.filter((_, j) => j !== i) })
                  },
                },
              }
            : undefined
        }
        rincian={(d, i) => (
          <DetailLimits
            d={d}
            bisaUbah={bisaUbah}
            onUbah={(baru) => {
              onUbah({ ...l, Detail: ganti(detail, i, baru) })
            }}
            terapkan={
              terapkan === undefined
                ? undefined
                : (f) => {
                    terapkan((lKini) => {
                      const ds = larikDari(lKini, 'Detail')
                      return { ...lKini, Detail: ganti(ds, i, f(ds[i] ?? {})) }
                    })
                  }
            }
          />
        )}
      />
    </>
  )
}

/**
 * Tab Limits prop. Pohonnya disalin ke keadaan layar sekali (pemanggil
 * memberi `key` per kontrak); suntingan hidup di sini.
 */
export default function TabLimitsProp({
  pohon,
  mode,
  idKontrak = '',
  edmJenisMaterial = '',
  opsi: opsiAwal,
}: {
  pohon: readonly SimpulLimit[]
  mode: ModeForm
  /** Pengenal kontrak — `InputParam.CARI1` RDB `GetAchievement`. */
  idKontrak?: string
  /** `EDMMaterialType` — `2` mengunci sel ber-`pyDisabledWhen` materialitas. */
  edmJenisMaterial?: string
  /** Isi dropdown tab ini; bila tidak diberikan, diminta sendiri. */
  opsi?: OpsiLimits
}) {
  // ⭐ PENAMPUNG HALAMAN (`../halaman.tsx`) — `TreatyIn.Limits` yang SAMA
  // dibaca dan ditulis tab Share (grid Kind of Treaty, `TreatyInPropshare`)
  // dan Achievement In IDR. Isian bertahan saat pindah tab.
  const [limits, setLimits] = useProperti<SimpulLimit[]>('Limits', () => [...pohon])
  const [opsi, setOpsi] = useState<OpsiLimits>(opsiAwal ?? OPSI_KOSONG)
  const bisaUbah = mode === 'ubah'
  const kunciMaterial = edmJenisMaterial.trim() === '2'
  // Halaman SESI Pega (`SearchData`, `TempQuarter*`, `FlagExcel`) — ikut
  // penampung supaya pilihan As At / Quarter Year tidak hilang saat pindah tab.
  const [asAt, setAsAt] = useProperti('SearchData.CARI1', '')
  const [tahun, setTahun] = useProperti('SearchData.CARI2', '')
  const [kuartal, setKuartal] = useProperti<string[]>('TempQuarter', [])
  const [tahunKuartal, setTahunKuartal] = useProperti<string[]>('TempQuarterYear', [])
  const [flagExcel, setFlagExcel] = useProperti('FlagExcel.CARI1', false)
  const [pesanAch, setPesanAch] = useState('')

  /**
   * `GetAchievement` atas SELURUH pohon — hanya medan yang Activity tulis
   * yang digabung ke tiap Treaty Group.
   */
  const jalankanAchievement = (cari: boolean, tahunPilih: string) => {
    setPesanAch('')
    hitungAchievement({
      idKontrak,
      cari,
      asAt,
      tahun: tahunPilih,
      rnmShareP: '',
      limits: limits.map((l) => ({
        TreatyType: teksDari(l, 'TreatyType'),
        Detail: larikDari(l, 'Detail').map((d) => ({ TreatyGroup: teksDari(d, 'TreatyGroup'), EPIList: larikDari(d, 'EPIList') })),
      })),
    })
      .then((h) => {
        setKuartal(h.kuartal)
        setTahunKuartal(h.tahunKuartal)
        setFlagExcel(h.flagExcel)
        setLimits((kini) =>
          kini.map((l, i) => ({
            ...l,
            Detail: larikDari(l, 'Detail').map((d, k) => ({ ...d, ...(h.limits[i]?.[k] ?? {}) })),
          })),
        )
      })
      .catch((e: unknown) => {
        setPesanAch(e instanceof Error ? e.message : String(e))
      })
  }
  const keadaanAch: KeadaanAchievement = {
    idKontrak,
    asAt,
    tahun,
    kuartal,
    tahunKuartal,
    flagExcel,
    pesan: pesanAch,
    pilihAsAt: (v) => {
      setAsAt(v)
      setTahun('') // DT `Reset_DT`: SearchData.CARI2 = ""
    },
    pilihTahun: (v) => {
      setTahun(v)
      jalankanAchievement(true, v)
    },
    segarkan: () => {
      jalankanAchievement(false, tahun)
    },
  }
  useEffect(() => {
    if (opsiAwal !== undefined) return
    let dibuang = false
    ambilOpsiLimits()
      .then((o) => {
        if (!dibuang) setOpsi(o)
      })
      .catch(() => undefined)
    return () => {
      dibuang = true
    }
  }, [opsiAwal])
  return (
    <OpsiLimitsCtx.Provider value={opsi}>
      <KunciMaterialCtx.Provider value={kunciMaterial}>
      <AchievementCtx.Provider value={keadaanAch}>
      <PohonLimitsCtx.Provider value={limits}>
      {/* ⭐ Bentuk Pega (8 Oktober 2026) — blok `Limits`: grid `masterDetail`
          Kind of Treaty (lebar 1315, kolom tombol 80), rincian
          `LimitProportional`. Kartu berchip DIGANTI grid. */}
      <div className="trin__blok trin__tab" role="group" aria-label={LIMITS_PROP.judul}>
        <BlokPega judul={LIMITS_PROP.judul}>
          <GridPega
            label={KOLOM_KIND_OF_TREATY[0]?.label}
            kolom={[
              {
                judul: KOLOM_KIND_OF_TREATY[0]?.label ?? '',
                lebar: 1315,
                isi: (l) => (teksDari(l, 'TreatyType') === '' ? LIMITS_PROP.kindBelumDipilih : teksDari(l, 'TreatyType')),
              },
            ]}
            baris={limits}
            tombol={
              bisaUbah
                ? {
                    lebar: 80,
                    tambah: {
                      label: LIMITS_PROP.tambah,
                      mati: kunciMaterial,
                      onKlik: () => {
                        // `TreatyInPropAdd(limits)`: `Limits(APPEND).ID = ""`.
                        setLimits([...limits, { ID: '', Detail: [] }])
                      },
                    },
                    hapus: {
                      label: LIMITS_PROP.hapus,
                      mati: kunciMaterial,
                      akses: (l) => `${LIMITS_PROP.hapus} ${teksDari(l, 'TreatyType') || (KOLOM_KIND_OF_TREATY[0]?.label ?? '')}`,
                      onKlik: (i) => {
                        setLimits(limits.filter((_, j) => j !== i))
                      },
                    },
                  }
                : undefined
            }
            rincian={(l, i) => (
              // Kind of Treaty = `.TreatyType`, diisi `SetTreatyTypeName_Act`
              // dari pilihan Treaty Type (`.Note`, kolom Name menu Reinsurance Type).
              <LimitProportional
                l={l}
                indeks={i}
                bisaUbah={bisaUbah}
                onUbah={(baru) => {
                  setLimits((ls) => ganti(ls, i, baru))
                }}
                terapkan={(f) => {
                  setLimits((ls) => ganti(ls, i, f(ls[i] ?? {})))
                }}
              />
            )}
          />
        </BlokPega>
      </div>
      </PohonLimitsCtx.Provider>
      </AchievementCtx.Provider>
      </KunciMaterialCtx.Provider>
    </OpsiLimitsCtx.Provider>
  )
}

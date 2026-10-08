// Perender kerangka bangkitan — isi tab Old dan New, apa adanya dari ekspor.
//
// ⛔ Komponen ini TIDAK memutuskan tata letak; `kerangka.gen.ts` yang
// memutuskannya (urutan butir, judul blok, kolom grid, desimal per sel,
// letak tombol). Yang diputuskan di sini hanya cara MENAMPILKAN satu jenis
// butir.
//
// ⭐ 7 Oktober 2026 — MODE SUNTING panel New. Di mode Edit (`ViewState 0`)
// medan dan sel grid panel New dapat disunting kecuali syarat baca-sajanya
// (`baca`, dari `pyReadOnlyCondition`/`pyDisabledWhen`) benar, dan tombol
// Add/Delete grid bekerja. Panel Old TETAP baca-saja tanpa pengecualian.
// Suntingan tinggal di keadaan layar — tidak ada yang menulis ke basis data
// sebelum Save, dan Save masih mati.
//
// ⭐ 7 Oktober 2026 — RINCIAN BARIS. Grid `expandPane` membuka Section
// rincian (`KERANGKA_RINCIAN`) di bawah barisnya, dengan konteks BARIS
// (`konteksBaris`). Medan/sel yang punya perilaku `change` di ekspor
// (`aksiUbah`) menjalankan rumusnya saat isian ditinggalkan atau Enter.
//
// ⭐ 7 Oktober 2026 — JALUR. Tiap konteks tahu jalurnya dari akar panel
// (`Limits(2).Detail(1)`); rumus berjalan atas akar + jalur itu, dan
// hasilnya diterapkan GABUNG TIGA ARAH di akar (`SisiPenyesuaian`). Sel
// pemicu (`.pxListSubscript`, kuncinya) dan skalar barisnya (`.Note`,
// `.Layer`) ikut dikirim — parameter dan syarat aksi Pega membacanya.

import { Fragment, useEffect, useId, useRef, useState, type ReactNode } from 'react'

import { Area, Field, FieldAngka, Kosong } from '../../../../inti/frontend/components/ui/dasar'
import { keInputTanggal } from '../../../../inti/frontend/lib/tanggalInput'
import type { BarisBersarang, OpsiLimitsTreatyIn, SisiPenyesuaian } from '../api'
import { golongan } from '../ekspor/golongan'
import type { AksiTombol, ButirKerangka, GridKerangka, MedanKerangka, SumberPilihan, TombolKerangka } from '../ekspor/jenis'
import { KERANGKA_INCLUDE, KERANGKA_RINCIAN } from '../ekspor/kerangka.gen'
import { syaratTerpenuhi } from '../ekspor/syarat'
import { PENYESUAIAN } from '../labelsPenyesuaian'
import {
  hapusBarisGrid,
  labelTombol,
  rantaiSesudahHapus,
  rantaiSesudahTambah,
  tambahBarisGrid,
  tambahDari,
  terkunci,
  tombolMati,
  tombolTampil,
  tulisanDari,
  unduhanDari,
  type JenisTulis,
} from './aksiTombol'
import { unduhAchievement } from './unduhAchievement'
import { AreaBacaSaja, BelumDibangun, Centang, MedanTakAda, TanggalBacaSaja, selNilai } from './medan'
import { rantaiRumus, type HasilTerapan, type LingkupRumus, type SelRumus } from './rumus'
import { denganNilaiKini, type Opsi } from './pilihan'
import { kunciBerjalur, nilaiJalur, type LangkahJalur } from './baris'
import { FieldTanggalKetik, KotakTanggalKetik } from './TanggalKetik'

/**
 * Include yang ISINYA tidak dibangun — Retro, keputusan §17
 * (`treatyin/docs/KEPUTUSAN-PENYELARASAN-REPO.md`): jarang dipakai. Wadah
 * tab dan syaratnya tetap dibangun; isinya kotak "jarang".
 */
export const INCLUDE_RETRO: ReadonlySet<string> = new Set([
  'TreatyInFacultativeShareCalculation',
  'TreatyInFacultativeShareCalculationOldData',
  'TreatyInFacultativeRetro',
])

type Baris = BarisBersarang

export interface KonteksKerangka {
  /** Halaman panel ini: `OLDDATA` (Old) atau akar (New) — keadaan TERKINI. */
  sisi: SisiPenyesuaian
  /** Halaman akar `TreatyIn` — sebagian ikatan Old menunjuknya. */
  akar: SisiPenyesuaian
  /** Halaman akar + `ViewState` sesi — tempat SEMUA syarat dinilai. */
  halaman: Readonly<Record<string, string>>
  /** Panel New di mode Edit. Tanpa ini semua butir baca-saja. */
  ubah?: boolean
  ubahMedan?: (kunci: string, nilai: string) => void
  ubahLarik?: (larik: string, baris: Baris[]) => void
  /**
   * Timpakan hasil satu rumus ke keadaan panel New. Hasil RANTAI membawa
   * akar sebelum/sesudah (`awal`/`akhir`) — diterapkan di akar.
   */
  terapkan?: (h: HasilTerapan) => void
  /** Daftar dropdown/autocomplete untuk satu sumber — `undefined` = belum ada. */
  opsi?: (sp: SumberPilihan, kunci: string) => Opsi[] | undefined
  /** `TreatyIn.OLDDATA` — sisi Old; dibaca rumus (Installment langkah 10). */
  lama?: SisiPenyesuaian
  /** Akar panel — di Section rincian `sisi` adalah BARIS, akarnya di sini. */
  panel?: SisiPenyesuaian
  /**
   * Jalur halaman ini dari akar panel — `[]` di akar. Tanpa jalur, hasil
   * rumus di rincian diterapkan ke barisnya lewat `ganti` saja.
   */
  jalur?: readonly LangkahJalur[]
  /** Daftar master yang dimuat panel ini — pasangan id ↔ nama (`Set…Name_Act`). */
  master?: OpsiLimitsTreatyIn
  /** Tulis satu medan AKAR panel — medan halaman SESI (`dari = 'sesi'`) di Section rincian. */
  ubahMedanAkar?: (kunci: string, nilai: string) => void
  /**
   * ⭐ Tombol TULIS (Submit / Decline offer EDM) — form yang menjalankannya.
   * Tanpa ini tombolnya mati (panel Old, uji satu tab).
   */
  tulis?: (jenis: JenisTulis) => void
  /** Tombol tulis sedang berjalan. */
  sibukTulis?: boolean
}

const tanpaAksi = () => undefined

/** Lingkup rumus konteks ini — akar panel, jalur, sisi Old, dan pemicunya. */
const lingkup = (k: KonteksKerangka, x: { sel?: SelRumus; halaman?: Readonly<Record<string, string>> } = {}): LingkupRumus => ({
  panel: k.panel ?? k.sisi,
  lama: k.lama,
  jalur: k.jalur,
  master: k.master,
  sel: x.sel,
  halaman: x.halaman ?? k.halaman,
})

/** `.pxListSubscript` halaman konteks ini (mulai 1) — 0 di akar. */
const nomorHalaman = (k: KonteksKerangka) => {
  const j = k.jalur ?? []
  const t = j[j.length - 1]
  return t === undefined ? 0 : t.indeks + 1
}

/** Skalar satu baris sebagai kunci BERTITIK (`.Note`) — halaman syarat/parameter sel itu. */
const bertitik = (b: Baris): Record<string, string> => {
  const out: Record<string, string> = {}
  for (const [k, v] of Object.entries(b)) if (typeof v === 'string') out[`.${k}`] = v
  return out
}

/**
 * Konteks SATU baris untuk Section rincian (`expandPane`).
 *
 * Ikatan `.X` menunjuk baris itu; syarat relatif membaca skalarnya lewat
 * kunci BERTITIK di halaman syarat (`.SpreadingTypeXOL`), dan kunci bertitik
 * baris induk (rincian di dalam rincian) tidak terbawa. Suntingan dan hasil
 * rumus MENGGANTIKAN baris itu di larik induknya lewat `ganti`.
 */
export function konteksBaris(k: KonteksKerangka, b: Baris, ganti: (nb: Baris) => void, langkah?: LangkahJalur): KonteksKerangka {
  const medan: Record<string, string> = {}
  const larik: Record<string, Baris[]> = {}
  for (const [kunci, v] of Object.entries(b)) {
    if (typeof v === 'string') medan[kunci] = v
    else if (Array.isArray(v)) larik[kunci] = v
  }
  const halaman: Record<string, string> = {}
  for (const [kunci, v] of Object.entries(k.halaman)) if (!kunci.startsWith('.')) halaman[kunci] = v
  for (const [kunci, v] of Object.entries(medan)) halaman[`.${kunci}`] = v
  const jalur = langkah === undefined || k.jalur === undefined ? undefined : [...k.jalur, langkah]
  return {
    ...k,
    sisi: { medan, larik },
    halaman,
    panel: k.panel ?? k.sisi,
    jalur,
    ubahMedan: (kunci, v) => {
      ganti({ ...b, [kunci]: v })
    },
    ubahLarik: (l, baris) => {
      ganti({ ...b, [l]: baris })
    },
    terapkan: (h) => {
      // ⭐ Hasil RANTAI berjalur diteruskan ke akar (gabung tiga arah atas
      // keadaan TERKINI); tanpa jalur, baris ini diganti seperti dahulu.
      if (jalur !== undefined && h.akhir !== undefined && k.terapkan !== undefined) {
        k.terapkan(h)
        return
      }
      ganti({ ...b, ...h.medan, ...h.larik })
    },
  }
}

/**
 * Rantai rumus perilaku `change` sebuah medan/sel — hanya di mode sunting,
 * dan hanya bila SELURUH langkahnya punya rumus (`rantaiRumus`).
 */
function rantaiUbah(
  k: KonteksKerangka,
  aksi: readonly AksiTombol[] | null | undefined,
  pemicu: { sel?: SelRumus; halaman?: Readonly<Record<string, string>> } = {},
) {
  if (k.ubah !== true || aksi === null || aksi === undefined) return undefined
  const r = rantaiRumus({ aksi })
  if (r === undefined) return undefined
  return (setPesan: (p: string[]) => void) => {
    setPesan([])
    r(k.sisi, lingkup(k, pemicu))
      .then((h) => {
        k.terapkan?.(h)
        setPesan(h.pesan)
      })
      .catch((e: unknown) => {
        setPesan([e instanceof Error ? e.message : String(e)])
      })
  }
}

/**
 * Pemicu perilaku `change` Pega: dijalankan saat isian DITINGGALKAN
 * (focusout) atau Enter — hanya bila nilainya berbeda dari saat terakhir
 * dipicu, sehingga Enter lalu keluar tidak menjalankannya dua kali.
 * Pembungkus `display: contents`: tata letak kontrolnya tidak bergeser.
 *
 * ⚠️ Dropdown di Pega memicu SEGERA saat dipilih; di sini saat ditinggalkan.
 *
 * ⛔ PERBAIKAN 7 Oktober 2026 — kotak yang MENYIMPAN isiannya saat
 * ditinggalkan (tanggal ketik, `TanggalKetik.tsx`) menulis nilai barunya di
 * blur yang SAMA; pembungkus ini masih melihat nilai lama, sehingga sel
 * Initial Date / Reporting Date tidak pernah memicu rumusnya. Pemeriksaan
 * kini diulang SEKALI sesudah render berikutnya — dan hanya sesudah blur
 * itu, supaya nilai yang diubah RUMUS tidak memicu perilaku `change`.
 */
function PemicuUbah({ nilai, onPicu, children }: { nilai: string; onPicu: () => void; children: ReactNode }) {
  const terakhir = useRef(nilai)
  const menyusul = useRef(false)
  const picu = () => {
    if (nilai === terakhir.current) return
    terakhir.current = nilai
    onPicu()
  }
  useEffect(() => {
    if (!menyusul.current) return
    menyusul.current = false
    picu()
  })
  return (
    <div
      className="tria__pemicu"
      onFocus={() => {
        terakhir.current = nilai
      }}
      onBlur={() => {
        picu()
        // Render dari blur ini (event diskret) selesai sebelum tenggat ini.
        menyusul.current = true
        setTimeout(() => {
          menyusul.current = false
        }, 0)
      }}
      onKeyDown={(e) => {
        if (e.key === 'Enter') picu()
      }}
    >
      {children}
    </div>
  )
}

/** Pesan rumus (`Property-Set-Messages`), apa adanya. */
function PesanRumus({ pesan }: { pesan: readonly string[] }) {
  if (pesan.length === 0) return null
  return (
    <span className="tria__galat" role="alert">
      {pesan.join(' · ')}
    </span>
  )
}

function halamanDari(dari: GridKerangka['dari'], k: KonteksKerangka): SisiPenyesuaian {
  if (dari === 'sesi') return k.panel ?? k.sisi
  return dari === 'akar' ? k.akar : k.sisi
}

/** Butir ini dapat disunting di layar sekarang? Ikatan ke AKAR (panel Old) tidak pernah. */
const bisaSunting = (k: KonteksKerangka, dari: GridKerangka['dari']) => k.ubah === true && dari !== 'akar'

/** Satu tombol ekspor di dalam grid (kepala atau baris). */
function TombolGrid({ t, mati, onKlik, akses }: { t: TombolKerangka; mati: boolean; onKlik: () => void; akses?: string }) {
  return (
    <button type="button" className="btn btn--ghost btn--sm" disabled={mati} aria-label={akses} onClick={onKlik}>
      {labelTombol(t)}
    </button>
  )
}

/**
 * Kontrol berdaftar — `pxDropdown` menjadi `<select>`, `pxAutoComplete`
 * menjadi isian ber-`datalist` (ketik bebas, pilihan dari daftar).
 */
function Berdaftar({ format, nilai, label, opsi, onUbah }: { format: string; nilai: string; label: string; opsi: readonly Opsi[]; onUbah: (v: string) => void }) {
  const idDaftar = useId()
  if (format === 'pxDropdown') {
    return (
      <select
        className="field__input"
        value={nilai}
        aria-label={label}
        onChange={(e) => {
          onUbah(e.target.value)
        }}
      >
        <option value="" />
        {denganNilaiKini(opsi, nilai, PENYESUAIAN.diLuarDaftar).map((o) => (
          <option key={o.value} value={o.value}>
            {o.label}
          </option>
        ))}
      </select>
    )
  }
  return (
    <>
      <input
        className="field__input"
        type="text"
        list={idDaftar}
        value={nilai}
        aria-label={label}
        onChange={(e) => {
          onUbah(e.target.value)
        }}
      />
      <datalist id={idDaftar}>
        {opsi.map((o) => (
          <option key={o.value} value={o.value}>
            {o.label}
          </option>
        ))}
      </datalist>
    </>
  )
}

/** Satu sel grid yang dapat disunting — kontrolnya menurut `pyFormat`. */
function SelSunting({
  format,
  nilai,
  label,
  opsi,
  onUbah,
}: {
  format: string
  nilai: string
  label: string
  opsi?: readonly Opsi[]
  onUbah: (v: string) => void
}) {
  if (opsi !== undefined && (format === 'pxDropdown' || format === 'pxAutoComplete')) {
    return <Berdaftar format={format} nilai={nilai} label={label} opsi={opsi} onUbah={onUbah} />
  }
  if (format === 'pxCheckbox') {
    return (
      <input
        type="checkbox"
        checked={nilai === 'true'}
        aria-label={label}
        onChange={(e) => {
          onUbah(e.target.checked ? 'true' : 'false')
        }}
      />
    )
  }
  if (format === 'pxDateTime') {
    // ⭐ 7 Oktober 2026 — tanggal DAPAT DIKETIK (`DD/MM/YYYY`, `18012024`,
    // …) atau dipilih dari ikon kalender (`TanggalKetik.tsx`). Nilai
    // TERSIMPAN (`20180101`) masuk apa adanya; keluarannya tetap
    // `YYYY-MM-DD`, bentuk yang kotak `type="date"` sebelumnya kirim.
    return (
      <KotakTanggalKetik
        label={label}
        value={nilai}
        onChange={(kabel) => {
          onUbah(keInputTanggal(kabel))
        }}
      />
    )
  }
  // ⭐ ANGKA BERPEMISAH RIBUAN — permintaan pemilik proses 7 Oktober 2026.
  //
  // ⛔ RALAT ATAS CATATAN SEBELUMNYA. Berkas ini dahulu berbunyi: *"Nilai
  // MENTAH ke kotak isian, tanpa format angka: memformat di tiap ketukan
  // membuat koma desimal mustahil diketik."* Pengamatannya BENAR, tetapi
  // kesimpulannya terlalu jauh: yang merusak pengetikan adalah
  // `formatNumber` — ia MEMBULATKAN dan MEMBUANG NOL EKOR, sehingga `12,`
  // menjadi `12` dan `1,0` menjadi `1`.
  //
  // `FieldAngka` (inti) memformat BAGIAN BULAT saja dan membiarkan ekor
  // desimal apa adanya selama kotaknya dipegang. `inti/frontend/lib/
  // angkaKetik.test.ts` memaku keempat keadaan tengah-pengetikan itu.
  if (format === 'pxNumber') {
    return (
      <FieldAngka
        label={label}
        value={nilai}
        desimal={2}
        onChange={onUbah}
      />
    )
  }
  // Sisanya TEKS — nilainya mentah, sebab bukan bilangan.
  return (
    <input
      className="field__input"
      type="text"
      value={nilai}
      aria-label={label}
      onChange={(e) => {
        onUbah(e.target.value)
      }}
    />
  )
}

/** Grid bangkitan; kolom bersyarat sel (`Auto Calculate` @313578) dinilai per kolom. */
export function GridEkspor({ g, k }: { g: GridKerangka; k: KonteksKerangka }) {
  const hal = halamanDari(g.dari, k)
  const ada = Object.prototype.hasOwnProperty.call(hal.larik, g.larik)
  const baris = hal.larik[g.larik] ?? []
  const sunting = bisaSunting(k, g.dari)
  // ⭐ Rincian baris (`expandPane`) — tertutup semula, seperti Pega.
  const rincian = g.rincian !== undefined ? KERANGKA_RINCIAN[g.rincian] : undefined
  const [terbuka, setTerbuka] = useState<ReadonlySet<number>>(() => new Set())
  const [pesanSel, setPesanSel] = useState<string[]>([])
  const tombolBaris = (i: number) => {
    const t = g.tombol[i] ?? null
    return sunting && tombolTampil(t, k.halaman) ? t : null
  }
  const tombolKepala = (i: number) => {
    const t = g.tombolKepala[i] ?? null
    return sunting && tombolTampil(t, k.halaman) ? t : null
  }
  // ⛔ Kolom tombol tampil hanya bila salah satu tombolnya tampil — di panel
  // Old dan di mode View kolom itu tidak ada, seperti sebelumnya.
  const tampil = g.kunci
    .map((_, i) => i)
    .filter((i) => {
      if (g.tombol[i] !== null && g.tombol[i] !== undefined) return tombolBaris(i) !== null || tombolKepala(i) !== null
      const s = g.syaratSel[i]
      return s === null || s === undefined || syaratTerpenuhi([s], k.halaman)
    })
  const total = tampil.reduce((a, i) => a + (g.lebar[i] ?? 0), 0) || 1
  const ganti = (baru: Baris[]) => k.ubahLarik?.(g.larik, baru)
  const lebarKolom = tampil.length + (rincian !== undefined ? 1 : 0)
  return (
    <div className="tria__grid">
      <div className="table-wrap">
        <table className="tria__tabel">
          {/* Lebar DARI ekspor sebagai perbandingan — tata letak responsif. */}
          <colgroup>
            {rincian !== undefined && <col className="tria__buka-kolom" />}
            {tampil.map((i) => (
              <col key={i} style={{ width: `${(((g.lebar[i] ?? 0) / total) * 100).toFixed(2)}%` }} />
            ))}
          </colgroup>
          <thead>
            <tr>
              {rincian !== undefined && <th scope="col" aria-label={PENYESUAIAN.rincianBuka} />}
              {tampil.map((i) => {
                const t = tombolKepala(i)
                return (
                  <th key={i} scope="col">
                    {t !== null ? (
                      // Tombol di sel KEPALA — `addRow` grid ini, atau Activity
                      // tambah baris (`TreatyInNonAddItem(Type=…)`).
                      <TombolGrid
                        t={t}
                        mati={
                          tombolMati(t, k.halaman) ||
                          (tambahBarisGrid(t) ? rantaiSesudahTambah(t) === undefined : tambahDari(t) === undefined)
                        }
                        onKlik={() => {
                          const tb = tambahDari(t)
                          if (tb !== undefined && tb.larik === g.larik) {
                            ganti([...baris, tb.baris(k.sisi, nomorHalaman(k))])
                          } else if (tambahBarisGrid(t)) {
                            const tambahan = [...baris, {}]
                            ganti(tambahan)
                            // `addRow` lalu Activity/DT sesudahnya (`Add Treaty Group`
                            // → DT `TreatyTypeSetIndex`), atas halaman yang barisnya
                            // SUDAH bertambah.
                            const sesudah = rantaiSesudahTambah(t)
                            if (sesudah === null || sesudah === undefined) return
                            sesudah({ ...k.sisi, larik: { ...k.sisi.larik, [g.larik]: tambahan } }, lingkup(k))
                              .then((h) => {
                                k.terapkan?.(h)
                                setPesanSel(h.pesan)
                              })
                              .catch((e: unknown) => {
                                setPesanSel([e instanceof Error ? e.message : String(e)])
                              })
                          }
                        }}
                      />
                    ) : (
                      g.kolom[i]
                    )}
                  </th>
                )
              })}
            </tr>
          </thead>
          <tbody>
            {baris.length === 0 && (
              <tr>
                <td colSpan={lebarKolom || 1}>
                  <Kosong
                    pesan={PENYESUAIAN.tanpaBaris}
                    petunjuk={ada ? PENYESUAIAN.petunjukGridKosong : PENYESUAIAN.takAdaDiWarisan}
                  />
                </td>
              </tr>
            )}
            {baris.map((b, r) => (
              <Fragment key={r}>
              <tr>
                {rincian !== undefined && (
                  <td className="tria__buka">
                    <button
                      type="button"
                      className="tria__buka-tombol"
                      aria-expanded={terbuka.has(r)}
                      aria-label={`${PENYESUAIAN.rincianBuka} ${String(r + 1)}`}
                      onClick={() => {
                        setTerbuka((x) => {
                          const y = new Set(x)
                          if (y.has(r)) y.delete(r)
                          else y.add(r)
                          return y
                        })
                      }}
                    >
                      {terbuka.has(r) ? '▾' : '▸'}
                    </button>
                  </td>
                )}
                {tampil.map((i) => {
                  const t = tombolBaris(i)
                  if (g.tombol[i] !== null && g.tombol[i] !== undefined) {
                    const sesudah = t === null ? null : rantaiSesudahHapus(t)
                    return (
                      <td key={i}>
                        {t !== null && hapusBarisGrid(t) && (
                          <TombolGrid
                            t={t}
                            mati={tombolMati(t, k.halaman) || sesudah === undefined}
                            akses={`${labelTombol(t)} ${String(r + 1)}`}
                            onKlik={() => {
                              const sisa = baris.filter((_, x) => x !== r)
                              // Indeks bergeser — rincian yang terbuka ditutup.
                              setTerbuka(new Set())
                              ganti(sisa)
                              // `deleteRow` lalu Activity sesudahnya, atas halaman yang
                              // barisnya SUDAH terhapus.
                              if (sesudah === null || sesudah === undefined) return
                              // Sel pemicu = baris yang DIHAPUS (`.pxListSubscript`,
                              // `.Note`, `.Layer` barisnya).
                              sesudah(
                                { ...k.sisi, larik: { ...k.sisi.larik, [g.larik]: sisa } },
                                lingkup(k, { sel: { larik: g.larik, indeks: r, kunci: '' }, halaman: { ...k.halaman, ...bertitik(b) } }),
                              )
                                .then((h) => {
                                  k.terapkan?.(h)
                                  setPesanSel(h.pesan)
                                })
                                .catch((e: unknown) => {
                                  setPesanSel([e instanceof Error ? e.message : String(e)])
                                })
                            }}
                          />
                        )}
                      </td>
                    )
                  }
                  const kunci = g.kunci[i] ?? ''
                  const fmt = g.format[i] ?? ''
                  const v = nilaiJalur(b, kunci)
                  // ⛔ Jalur berindeks (`RnmLimitListDisplay(1).Value`) hasil
                  // rumus — tidak disunting langsung.
                  if (sunting && !terkunci(g.baca[i], k.halaman) && !kunciBerjalur(kunci)) {
                    const sel = (
                      <SelSunting
                        format={fmt}
                        nilai={v}
                        label={g.kolom[i] ?? kunci}
                        opsi={pilihanSel(g.pilihan[i] ?? null, kunci, k)}
                        onUbah={(nv) => {
                          ganti(baris.map((x, y) => (y === r ? { ...x, [kunci]: nv } : x)))
                        }}
                      />
                    )
                    // ⭐ Perilaku `change` sel (`% Installment` →
                    // SetTotalInstallment) berjalan di halaman GRID ini.
                    const jalankan = rantaiUbah(k, g.aksiUbah[i], {
                      sel: { larik: g.larik, indeks: r, kunci },
                      halaman: { ...k.halaman, ...bertitik(b) },
                    })
                    return (
                      <td key={i}>
                        {jalankan === undefined ? (
                          sel
                        ) : (
                          <PemicuUbah
                            nilai={v}
                            onPicu={() => {
                              jalankan(setPesanSel)
                            }}
                          >
                            {sel}
                          </PemicuUbah>
                        )}
                      </td>
                    )
                  }
                  if (fmt === 'pxCheckbox') {
                    return (
                      <td key={i}>
                        <input type="checkbox" checked={v === 'true'} disabled readOnly aria-label={g.kolom[i]} />
                      </td>
                    )
                  }
                  const jenis = golongan(kunci, fmt)
                  return (
                    <td key={i} className={jenis === 'teks' || jenis === 'tanggal' ? undefined : 'tria__angka'}>
                      {selNilai(jenis, v, g.desimal[i] ?? null)}
                    </td>
                  )
                })}
              </tr>
              {rincian !== undefined && terbuka.has(r) && (
                <tr className="tria__rincian">
                  <td colSpan={lebarKolom}>
                    <RenderKerangka
                      isi={rincian}
                      k={konteksBaris(
                        k,
                        b,
                        (nb) => {
                          ganti(baris.map((x, y) => (y === r ? nb : x)))
                        },
                        { larik: g.larik, indeks: r },
                      )}
                    />
                  </td>
                </tr>
              )}
              </Fragment>
            ))}
          </tbody>
        </table>
      </div>
      <PesanRumus pesan={pesanSel} />
      {/* ⛔ 8 Oktober 2026 — `rincianHilang` DIKOSONGKAN atas permintaan
          pemilik proses. Tanpa penjaga ini, sisanya yang tampil adalah
          tanda kurung berisi nama aturan Pega: `(DetailInstallment)`. */}
      {g.rincianHilang !== undefined && PENYESUAIAN.rincianHilang !== '' && (
        <span className="tria__redup">
          {PENYESUAIAN.rincianHilang} ({g.rincianHilang})
        </span>
      )}
    </div>
  )
}

/** Daftar satu sel/medan berdaftar, atau `undefined` (kotak teks). */
function pilihanSel(sp: SumberPilihan | null | undefined, kunci: string, k: KonteksKerangka): Opsi[] | undefined {
  if (sp === null || sp === undefined) return undefined
  return k.opsi?.(sp, kunci)
}

/** Medan bangkitan — kontrolnya menurut `pyFormat`; dapat disunting di mode Edit panel New. */
export function MedanEkspor({ m, k }: { m: MedanKerangka; k: KonteksKerangka }) {
  const [pesan, setPesan] = useState<string[]>([])
  const hal = halamanDari(m.dari, k)
  const label = m.label !== '' ? m.label : (m.caption ?? '')
  const sunting = bisaSunting(k, m.dari) && !terkunci(m.baca, k.halaman)
  if (!Object.prototype.hasOwnProperty.call(hal.medan, m.kunci) && !sunting) return <MedanTakAda label={label} />
  const v = hal.medan[m.kunci] ?? ''
  if (sunting) {
    const ubah = (nv: string) => (m.dari === 'sesi' ? k.ubahMedanAkar : k.ubahMedan)?.(m.kunci, nv)
    const opsi = pilihanSel(m.pilihan, m.kunci, k)
    const kontrol = (() => {
      if (opsi !== undefined && (m.format === 'pxDropdown' || m.format === 'pxAutoComplete')) {
        return (
          <div className="field">
            <label className="field__label">{label}</label>
            <Berdaftar format={m.format} nilai={v} label={label} opsi={opsi} onUbah={ubah} />
          </div>
        )
      }
      switch (m.format) {
        case 'pxTextArea':
          return <Area label={label} value={v} onChange={ubah} />
        case 'pxDateTime':
          // Dapat diketik atau dipilih dari kalender — `TanggalKetik.tsx`.
          return <FieldTanggalKetik label={label} value={v} onChange={ubah} />
        case 'pxCheckbox':
          return <Centang label={m.caption ?? label} nilai={v} bacaSaja={false} onUbah={ubah} />
        default:
          // ⛔ Nilai MENTAH, tanpa format angka — lihat `SelSunting`.
          return <Field label={label} value={v} onChange={ubah} />
      }
    })()
    // ⭐ Perilaku `change` ekspor (`Installment` → TreatyInSetValueInstallment).
    const jalankan = rantaiUbah(k, m.aksiUbah)
    if (jalankan === undefined) return kontrol
    return (
      <>
        <PemicuUbah
          nilai={v}
          onPicu={() => {
            jalankan(setPesan)
          }}
        >
          {kontrol}
        </PemicuUbah>
        <PesanRumus pesan={pesan} />
      </>
    )
  }
  switch (m.format) {
    case 'pxTextArea':
      return <AreaBacaSaja label={label} nilai={v} />
    case 'pxDateTime':
      // ⛔ Nilai TERSIMPAN ke kotak tanggal — bukan terjemahan tampil (§5).
      return <TanggalBacaSaja label={label} nilai={v} />
    case 'pxCheckbox':
      return <Centang label={m.caption ?? label} nilai={v} bacaSaja onUbah={tanpaAksi} />
    default:
      return <Field label={label} value={selNilai(golongan(m.kunci, m.format), v, m.desimal)} onChange={tanpaAksi} readOnly />
  }
}

/**
 * Tombol bangkitan di luar grid.
 *
 * Panel New mode Edit: tombol tambah baris (`TreatyInNonAddItem`,
 * `TreatyInPropAdd`, …) dan tombol yang SELURUH rantai Activity-nya punya
 * rumus (`komponen/rumus.ts`) bekerja. Yang lain DIMATIKAN, bukan
 * dihilangkan. Panel Old: mati.
 */
function TombolEkspor({ t, k }: { t: TombolKerangka; k: KonteksKerangka }) {
  const [sibuk, setSibuk] = useState(false)
  const [pesan, setPesan] = useState<string[]>([])
  if (unduhanDari(t) === 'achievement') {
    // `GenerateCSVTreaty` — membaca halaman tombol ini; tidak mengubah apa pun.
    return (
      <button
        type="button"
        className="btn btn--ghost btn--sm"
        onClick={() => {
          unduhAchievement(k.sisi)
        }}
      >
        {labelTombol(t)}
      </button>
    )
  }
  // ⭐ Submit / Decline offer EDM — form yang menulis (`k.tulis`).
  const tulisan = tulisanDari(t)
  if (tulisan !== undefined && k.tulis !== undefined) {
    const tulis = k.tulis
    return (
      <button
        type="button"
        className={tulisan === 'submit' ? 'btn btn--primary btn--sm' : 'btn btn--ghost btn--sm'}
        disabled={k.sibukTulis === true || tombolMati(t, k.halaman)}
        onClick={() => {
          tulis(tulisan)
        }}
      >
        {labelTombol(t)}
      </button>
    )
  }
  const tb = k.ubah === true ? tambahDari(t) : undefined
  const rumus = k.ubah === true && tb === undefined ? rantaiRumus(t) : undefined
  if (tb === undefined && rumus === undefined) {
    return (
      <button
        type="button"
        className="btn btn--ghost btn--sm"
        disabled
        title={
          t.aksi.some((a) => a.aktivitas === 'InsertToLogAchievement')
            ? PENYESUAIAN.achievementMenunggu
            : k.ubah === true
              ? PENYESUAIAN.rumusBelum
              : PENYESUAIAN.tombolTulisMati
        }
      >
        {labelTombol(t)}
      </button>
    )
  }
  return (
    <>
      <button
        type="button"
        className="btn btn--ghost btn--sm"
        disabled={sibuk || tombolMati(t, k.halaman)}
        onClick={() => {
          if (tb !== undefined) {
            k.ubahLarik?.(tb.larik, [...(k.sisi.larik[tb.larik] ?? []), tb.baris(k.sisi, nomorHalaman(k))])
            return
          }
          if (rumus === undefined) return
          setSibuk(true)
          setPesan([])
          rumus(k.sisi, lingkup(k))
            .then((h) => {
              k.terapkan?.(h)
              setPesan(h.pesan)
            })
            .catch((e: unknown) => {
              setPesan([e instanceof Error ? e.message : String(e)])
            })
            .finally(() => {
              setSibuk(false)
            })
        }}
      >
        {labelTombol(t)}
      </button>
      {pesan.length > 0 && (
        <span className="tria__galat" role="alert">
          {pesan.join(' · ')}
        </span>
      )}
    </>
  )
}

/** Satu daftar butir kerangka. */
export function RenderKerangka({ isi, k }: { isi: readonly ButirKerangka[]; k: KonteksKerangka }) {
  return (
    <>
      {isi.map((b) => {
        if (!syaratTerpenuhi(b.syarat, k.halaman)) return null
        switch (b.t) {
          case 'blok':
            return (
              <div key={b.at} className="tria__blok">
                {b.judul !== '' && <h5 className="tria__subjudul">{b.judul}</h5>}
                <RenderKerangka isi={b.anak} k={k} />
              </div>
            )
          case 'grid':
            return <GridEkspor key={b.at} g={b} k={k} />
          case 'medan':
            return <MedanEkspor key={b.at} m={b} k={k} />
          case 'teks':
            return (
              <span key={b.at} className="tria__teks-sel">
                {b.teks}
              </span>
            )
          case 'tombol':
            return <TombolEkspor key={b.at} t={b} k={k} />
          case 'include': {
            if (INCLUDE_RETRO.has(b.nama)) {
              return <BelumDibangun key={b.at} judul={PENYESUAIAN.retroJarang} petunjuk={PENYESUAIAN.retroJarangPetunjuk} />
            }
            const anak = KERANGKA_INCLUDE[b.nama]
            if (anak === undefined) {
              // ⛔ Tanpa penjaga, petunjuknya menjadi `NamaSection: ` — nama
              // aturan Pega berikut titik dua yang menggantung.
              return (
                <BelumDibangun
                  key={b.at}
                  judul={PENYESUAIAN.belumDibangun}
                  petunjuk={PENYESUAIAN.includeTakAda === '' ? '' : `${b.nama}: ${PENYESUAIAN.includeTakAda}`}
                />
              )
            }
            return <RenderKerangka key={b.at} isi={anak} k={k} />
          }
          default:
            return null
        }
      })}
    </>
  )
}

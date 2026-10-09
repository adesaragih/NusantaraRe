// Renderer pohon tata Claim Prop. Satu komponen untuk semua section: kondisi tampil / hanya-baca / nonaktif sudah
// dievaluasi server; di sini hanya pemetaan kendali Pega ke elemen HTML. Medan berkendali aksi (refresh ber-activity)
// memanggil server saat nilainya berubah; medan tanpa aksi (postValue) hanya mengubah nilai lokal yang ikut terkirim
// pada aksi berikutnya.
//
// Rupa = layout lama Pega dirapikan (work owner 08-10-2026 "berikut layout lama, ikuti dan rapihkan") dengan kulit
// Kelola User: medan "Stacked with labels left" (label di kiri, nilai di kanan; hanya-baca = teks), `dua` = Inline grid
// double, `sebaris` = Inline (Quarter/Year, RNM Share, baris tombol), `tab` = layout group Tab, `judul` = kepala tengah,
// tombol ikon dari XML (pi-plus / pi-trash / pi-pencil / pi-check), grid ber-paging (pyGridPaginator). Pengelompokan di
// `susun.ts`.

import { createContext, Fragment, useContext, useEffect, useRef, useState } from 'react'

import { PilihSaring } from '../../../../inti/frontend/components/ui/pilihSaring'
import type { Halaman, Pilihan, Tata } from '../api'
import { tampilTanggal } from '../ketikTanggal'
import { ambil, hanyaAngka, jalurBaris, tampilAngka } from '../nilai'
import InputAngka from './InputAngka'
import InputTanggal from './InputTanggal'
import { barisTerbuka, bukaAwal, rincianUntuk, type RincianGrid } from './rincian'
import { jumlahHalaman, potongHalaman, ratakan, susunIsi, susunLayar, type Butir } from './susun'
import { butirSaring, teksTerpilih } from './tampilanNama'

export interface KonteksTata {
  h: Halaman
  /** postValue: ubah nilai lokal saja. */
  ubah: (jalur: string, nilai: string) => void
  /** Ubah nilai lalu jalankan aksi server (refresh ber-activity / tombol). */
  aksi: (aksi: string, indeks?: number, ubahan?: Record<string, string>) => void
  opsi: (sumber: string, indeks: number) => Pilihan[]
  /** Saran autocomplete (adjuster, provinsi, rekening, allocation). */
  saran?: (sumber: string, indeks: number, cari: string) => void
  pesanMedan: Record<string, string[]>
  sibuk: boolean
  /** Panel rinci baris grid (expand pane AdjustmentDetail) - hanya untuk grid `rincian.daftar`. */
  rincian?: RincianGrid
}

/** Indeks baris (1..n) dari jalur `daftar(n).prop`; 0 bila jalur halaman. */
function indeksDari(jalur: string): number {
  const m = /\((\d+)\)\.[^()]*$/.exec(jalur)
  return m ? Number(m[1]) : 0
}

const idMedan = (jalur: string) => `claimprop-${jalur.replace(/[^A-Za-z0-9]/g, '-')}`

/** Teks tampilan nilai hanya-baca menurut kendalinya. */
function teksTampil(t: Tata, v: string, opsi: () => Pilihan[]): string {
  switch (t.kendali) {
    case 'angka':
      return tampilAngka(v)
    case 'tanggal':
    case 'tanggal-waktu':
      return tampilTanggal(v, t.kendali)
    case 'pilih':
    case 'radio':
      return opsi().find((o) => o.nilai === v)?.label ?? v
    default:
      return v
  }
}

/** Ikon tombol dari XML (pi-plus, pi-trash, pi-pencil, pi-check). */
function Ikon({ jenis }: { jenis: NonNullable<Tata['ikon']> }) {
  const d = {
    tambah: 'M12 5v14M5 12h14',
    hapus: 'M4 7h16M10 11v6M14 11v6M6 7l1 13h10l1-13M9 7V4h6v3',
    ubah: 'M4 20h4L19 9l-4-4L4 16v4zM14 6l4 4',
    simpan: 'M5 12l5 5 9-10',
  }[jenis]
  return (
    <svg viewBox="0 0 24 24" width="15" height="15" aria-hidden="true" focusable="false">
      <path d={d} fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" />
    </svg>
  )
}

/**
 * Kendali satu medan. `sel` = sel grid (isian ringkas); hanya-baca = teks nilai (layout lama), isian = kotak
 * `field__input` kulit Kelola User.
 */
function Medan({
  t,
  k,
  indeks = 0,
  jalur,
  id,
}: {
  t: Tata
  k: KonteksTata
  indeks?: number
  jalur: string
  id?: string
}) {
  const v = ambil(k.h, jalur)
  // Nilai saat fokus: textarea beraksi memanggil server hanya bila isinya berubah (event `change` XML).
  const awal = useRef<string | null>(null)
  const sel = indeks > 0
  const kunci = !!(t.hanyaBaca || t.nonaktif)
  const n = indeks || indeksDari(jalur)
  const opsi = () => k.opsi(t.sumber ?? '', n)
  const ganti = (baru: string, segera: boolean) => {
    if (t.aksi && segera) k.aksi(t.aksi, n, { [jalur]: baru })
    else k.ubah(jalur, baru)
  }
  const angka = t.kendali === 'angka'
  const kelas = ['field__input', sel && 'claimprop__input--sel', angka && 'claimprop__input--angka']
    .filter(Boolean)
    .join(' ')

  if (t.kendali === 'centang') {
    return (
      <span className="claimprop__centang">
        <input
          id={id}
          type="checkbox"
          checked={v === 'true'}
          disabled={kunci || k.sibuk}
          onChange={(e) => ganti(e.target.checked ? 'true' : 'false', true)}
        />
      </span>
    )
  }
  if (kunci) {
    const teks = t.tampilan ? teksTerpilih(t, v, (j) => ambil(k.h, j), opsi()) : teksTampil(t, v, opsi)
    return (
      <span id={id} className={angka ? 'claimprop__nilai claimprop__angka' : 'claimprop__nilai'}>
        {teks}
      </span>
    )
  }
  if (t.kendali === 'radio' && !sel) {
    return (
      <span id={id} className="claimprop__radio" role="radiogroup">
        {opsi().map((o) => (
          <label key={o.nilai}>
            <input
              type="radio"
              name={idMedan(jalur)}
              checked={v === o.nilai}
              disabled={k.sibuk}
              onChange={() => ganti(o.nilai, true)}
            />
            {o.label}
          </label>
        ))}
      </span>
    )
  }
  switch (t.kendali) {
    case 'angka':
      // hanya angka, separator Indonesia, maks 4 desimal (work owner 08-10-2026); nilai halaman tetap mentah
      return (
        <InputAngka
          id={id}
          className={kelas}
          value={v}
          disabled={k.sibuk}
          onChange={(mentah) => k.ubah(jalur, mentah)}
          onBlur={() => t.aksi && ganti(v, true)}
        />
      )
    case 'telepon':
      return (
        <input
          id={id}
          type="tel"
          inputMode="numeric"
          className={kelas}
          value={v}
          disabled={k.sibuk}
          onChange={(e) => k.ubah(jalur, hanyaAngka(e.target.value))}
          onBlur={(e) => t.aksi && ganti(hanyaAngka(e.target.value), true)}
        />
      )
    case 'tanggal':
    case 'tanggal-waktu':
      // diketik dd-mm-yyyy (+ hh:mm), pemisah otomatis, tombol kalender (work owner 08-10-2026); aksi server hanya saat
      // isian lengkap dan sah atau dikosongkan
      return (
        <InputTanggal
          id={id}
          className={kelas}
          jenis={t.kendali}
          value={v}
          disabled={k.sibuk}
          onChange={(nilai, final) => (final ? ganti(nilai, true) : k.ubah(jalur, nilai))}
        />
      )
    case 'area':
      return (
        <textarea
          id={id}
          className="field__input claimprop__area"
          value={v}
          disabled={k.sibuk}
          onChange={(e) => k.ubah(jalur, e.target.value)}
          onFocus={(e) => (awal.current = e.target.value)}
          onBlur={(e) => t.aksi && e.target.value !== awal.current && ganti(e.target.value, true)}
        />
      )
    case 'pilih':
    case 'radio': {
      const daftar = opsi()
      return (
        <select id={id} className={kelas} value={v} disabled={k.sibuk} onChange={(e) => ganti(e.target.value, true)}>
          <option value="">Choose</option>
          {daftar.map((o) => (
            <option key={o.nilai} value={o.nilai}>
              {o.label}
            </option>
          ))}
          {v !== '' && !daftar.some((o) => o.nilai === v) && <option value={v}>{v}</option>}
        </select>
      )
    }
    case 'otomatis': {
      // Dropdown yang dapat dicari (keputusan work owner 08-10-2026: Consultant ID / Adjuster ID "model dropdown yang
      // bisa di search"): PilihSaring inti, saringan ke server lewat `saran`; nilai berubah hanya saat butir dipilih.
      // Medan ber-`tampilan` (Consultant / Adjuster) memilih dan menampilkan nama saja, nilai tetap ID (tampilanNama.ts).
      if (!sel) {
        return (
          <span className="claimprop__saring">
            <PilihSaring
              label=""
              value={v}
              teksTerpilih={teksTerpilih(t, v, (j) => ambil(k.h, j), opsi())}
              opsi={butirSaring(t, opsi())}
              onCari={(kata) => k.saran?.(t.sumber ?? '', n, kata)}
              onPilih={(o) => ganti(o.value, true)}
            />
          </span>
        )
      }
      const idDaftar = `${idMedan(jalur)}-saran`
      return (
        <>
          <input
            id={id}
            className={kelas}
            list={idDaftar}
            value={v}
            disabled={k.sibuk}
            onFocus={() => k.saran?.(t.sumber ?? '', n, '')}
            onChange={(e) => {
              k.ubah(jalur, e.target.value)
              k.saran?.(t.sumber ?? '', n, e.target.value)
            }}
            onBlur={(e) => t.aksi && ganti(e.target.value, true)}
          />
          <datalist id={idDaftar}>
            {opsi().map((o) => (
              <option key={o.nilai + o.label} value={o.nilai}>
                {o.label}
              </option>
            ))}
          </datalist>
        </>
      )
    }
    default:
      return (
        <input
          id={id}
          className={kelas}
          value={v}
          disabled={k.sibuk}
          onChange={(e) => k.ubah(jalur, e.target.value)}
          onBlur={(e) => t.aksi && ganti(e.target.value, true)}
        />
      )
  }
}

/**
 * Nomor baris adjustment (1..n) pemilik panel rinci yang sedang dirender; 0 di luar panel. Tombol panel (Send to
 * Committe, Acceptation, Generate DLA, View Komite No) adalah aksi baris: server menyaring tata baris lewat Indeks
 * (`services/aksi.go` aksiBarisAdj), jadi tanpa nomor ini aksinya ditolak "tidak tersedia" (laporan work owner
 * 09-10-2026).
 */
export const BarisAdjustment = createContext(0)

/** Indeks yang dikirim tombol: indeks baris grid-nya sendiri, atau nomor baris panel rinci yang memuatnya. */
export function indeksTombol(indeks: number, baris: number): number {
  return indeks > 0 ? indeks : baris
}

/**
 * Tombol: ikon dari XML = tombol ikon (label jadi keterangan); Save / Submit = tombol utama (Kelola User: Simpan);
 * selainnya tombol bertepi. Tanpa label dan tanpa ikon = pengaturan. `utama` = tombol kaki modal.
 */
export function Tombol({ t, k, indeks = 0, utama }: { t: Tata; k: KonteksTata; indeks?: number; utama?: boolean }) {
  const baris = useContext(BarisAdjustment)
  const primer = utama ?? (indeks === 0 && !t.ikon && (t.label === 'Save' || t.label === 'Submit'))
  const kelas = [
    'btn',
    primer ? 'btn--primary' : 'btn--ghost',
    (indeks > 0 || t.ikon) && 'btn--sm',
    t.ikon && 'claimprop__ikon',
    t.ikon === 'hapus' && 'claimprop__ikon--hapus',
  ]
    .filter(Boolean)
    .join(' ')
  const nama = t.label || t.id
  return (
    <button
      type="button"
      className={kelas}
      disabled={t.nonaktif}
      aria-busy={k.sibuk || undefined}
      title={t.catatan ?? nama}
      aria-label={t.ikon || !t.label ? nama : undefined}
      onClick={(e) => {
        e.stopPropagation()
        k.aksi(t.aksi ?? t.id ?? '', indeksTombol(indeks, baris))
      }}
    >
      {t.ikon ? <Ikon jenis={t.ikon} /> : t.label ? t.label : '⚙'}
    </button>
  )
}

/** Satu baris "Stacked with labels left": label di kiri, nilai + satuan + tombol menempel di kanan. */
function BarisMedan({ t, tombol, k }: { t: Butir; tombol: readonly Tata[]; k: KonteksTata }) {
  const jalur = t.jalur ?? ''
  const id = idMedan(jalur)
  const pesan = k.pesanMedan[jalur] ?? []
  return (
    <div className={'claimprop__baris' + (t.label ? '' : ' claimprop__baris--tanpa-label')} title={t.catatan}>
      <label className="claimprop__label-medan" htmlFor={id}>
        {t.label}
        {t.wajib && <span className="field__req">*</span>}
      </label>
      <div className="claimprop__isi-medan">
        <Medan t={t} k={k} jalur={jalur} id={id} />
        {t.satuan && <span className="claimprop__satuan">{t.satuan}</span>}
        {tombol.map((b, i) => (
          <Tombol key={i} t={b} k={k} />
        ))}
        {pesan.length > 0 && <span className="field__error">{pesan.join('; ')}</span>}
      </div>
    </div>
  )
}

/** Layout Inline: anak sebaris; label bagian = label baris (Quarter/Year). */
function Sebaris({ t, k }: { t: Tata; k: KonteksTata }) {
  const isi = ratakan(t.anak ?? [])
  if (isi.length === 0) return null
  const baris = (
    <div className="claimprop__sebaris">
      {isi.map((u, i) => {
        if (u.jenis === 'medan') {
          const jalur = u.jalur ?? ''
          return (
            <span key={i} className="claimprop__sebaris-medan" title={u.catatan}>
              {u.label && (
                <label className="claimprop__label-sebaris" htmlFor={idMedan(jalur)}>
                  {u.label}
                  {u.wajib && <span className="field__req">*</span>}
                </label>
              )}
              <Medan t={u} k={k} jalur={jalur} id={idMedan(jalur)} />
              {u.satuan && <span className="claimprop__satuan">{u.satuan}</span>}
            </span>
          )
        }
        if (u.jenis === 'tombol') return <Tombol key={i} t={u} k={k} />
        if (u.jenis === 'label') {
          return (
            <span key={i} className="claimprop__satuan">
              {u.label}
            </span>
          )
        }
        return <TataView key={i} tata={[u]} k={k} />
      })}
    </div>
  )
  if (!t.label) return baris
  return (
    <div className="claimprop__baris">
      <span className="claimprop__label-medan">{t.label}</span>
      <div className="claimprop__isi-medan">{baris}</div>
    </div>
  )
}

/** Layout Inline grid double: setiap anak satu sel; bagian tanpa judul = satu kolom. */
function Dua({ t, k }: { t: Tata; k: KonteksTata }) {
  return (
    <div className="claimprop__dua">
      {(t.anak ?? []).map((c, i) => (
        <div key={i} className="claimprop__sel">
          <TataView tata={c.jenis === 'bagian' && !c.label && !c.letak ? (c.anak ?? []) : [c]} k={k} />
        </div>
      ))}
    </div>
  )
}

/** Layout group Tab: satu tab per bagian berjudul. */
function Tab({ t, k }: { t: Tata; k: KonteksTata }) {
  const tab = (t.anak ?? []).filter((a) => a.jenis === 'bagian')
  const [aktif, setAktif] = useState(0)
  const i = Math.min(aktif, Math.max(0, tab.length - 1))
  if (tab.length === 0) return null
  return (
    <div className="claimprop__tab">
      <div className="tabs" role="tablist">
        {tab.map((a, j) => (
          <button
            key={j}
            type="button"
            role="tab"
            aria-selected={j === i}
            className={'tabs__item' + (j === i ? ' tabs__item--aktif' : '')}
            onClick={() => setAktif(j)}
          >
            {a.label}
          </button>
        ))}
      </div>
      <div role="tabpanel" className="claimprop__tab-isi">
        <TataView tata={tab[i]?.anak ?? []} k={k} />
      </div>
    </div>
  )
}

/**
 * Layout bebas berkolom XML (`LetakTabel`, mis. Deductible / Treaty Gross / Total di AdjustmentDetail): baris pertama =
 * judul kolom, sel medan tanpa label; judul kolom berisi angka rata kanan seperti nilainya.
 */
function TabelTetap({ t, k }: { t: Tata; k: KonteksTata }) {
  const [kepala, ...isi] = (t.anak ?? []).filter((b) => b.jenis === 'bagian')
  if (!kepala) return null
  const angka = (j: number) => isi.some((b) => b.anak?.[j]?.kendali === 'angka')
  const sel = (c: Tata) => {
    if (c.jenis === 'medan') {
      const jalur = c.jalur ?? ''
      return <Medan t={c} k={k} indeks={indeksDari(jalur)} jalur={jalur} />
    }
    return c.jenis === 'label' ? c.label : null
  }
  return (
    <table className="claimprop__tabel claimprop__tabel--tetap">
      <thead>
        <tr>
          {(kepala.anak ?? []).map((c, j) => (
            <th key={j}>{angka(j) ? <span className="claimprop__angka">{c.label}</span> : c.label}</th>
          ))}
        </tr>
      </thead>
      <tbody>
        {isi.map((b, i) => (
          <tr key={i}>
            {(b.anak ?? []).map((c, j) => (
              <td key={j}>{sel(c)}</td>
            ))}
          </tr>
        ))}
      </tbody>
    </table>
  )
}

/** Kepala layar: judul tahap di tengah, Claim No di bawahnya. */
function Kepala({ t, k }: { t: Tata; k: KonteksTata }) {
  const isi = ratakan(t.anak ?? [])
  return (
    <div className="claimprop__kepala">
      {isi.map((u, i) => {
        if (u.jenis === 'label') {
          return i === 0 ? (
            <h3 key={i} className="claimprop__kepala-judul">
              {u.label}
            </h3>
          ) : (
            <div key={i} className="claimprop__kepala-teks">
              {u.label}
            </div>
          )
        }
        if (u.jenis === 'medan') {
          return (
            <div key={i} className="claimprop__kepala-teks">
              <span className="claimprop__label-sebaris">{u.label}</span> <Medan t={u} k={k} jalur={u.jalur ?? ''} />
            </div>
          )
        }
        return null
      })}
    </div>
  )
}

function Grid({ t, k }: { t: Tata; k: KonteksTata }) {
  const rinci = rincianUntuk(k.rincian, t.jalur)
  const kolom = t.kolom ?? []
  const semua = t.baris ?? []
  // Panel rinci baris terbuka tanpa klik (work owner 08-10-2026 "tidak harus klik angka sebelah kiri"): baris terbaru
  // terbuka sejak awal dan sesudah Add; klik di mana pun pada baris membuka / menutup.
  const nBaris = semua.length
  const [buka, setBuka] = useState<number | null>(() => bukaAwal(nBaris))
  const [hal, setHal] = useState(1)
  const adaRinci = rinci !== undefined
  const nLalu = useRef(nBaris)
  useEffect(() => {
    const lalu = nLalu.current
    nLalu.current = nBaris
    if (nBaris === lalu) return
    setBuka((b) => barisTerbuka(b, lalu, nBaris))
    if (adaRinci && nBaris > lalu && t.perHalaman) setHal(jumlahHalaman(nBaris, t.perHalaman))
  }, [nBaris, adaRinci, t.perHalaman])
  const nHal = jumlahHalaman(semua.length, t.perHalaman)
  const halIni = Math.min(hal, nHal)
  const awal = t.perHalaman ? (halIni - 1) * t.perHalaman : 0
  const baris = potongHalaman(semua, t.perHalaman, halIni)
  // Add di sel kepala kolom tombol terakhir (pzPegaDefaultGridIcons); tanpa kolom tombol = sel kepala tambahan
  const kolomAkhirTombol = kolom[kolom.length - 1]?.jenis === 'tombol'
  const selKepalaTambah = !!t.tambah && !kolomAkhirTombol
  const lebar = kolom.length + (t.bernomor ? 1 : 0) + (selKepalaTambah ? 1 : 0)
  return (
    <div className="claimprop__grid" title={t.catatan}>
      {t.perHalaman && semua.length > t.perHalaman ? (
        <div className="claimprop__pager">
          <button
            type="button"
            className="btn btn--ghost btn--sm"
            disabled={halIni <= 1}
            onClick={() => setHal(halIni - 1)}
            aria-label="Previous"
          >
            ‹
          </button>
          <span>
            Page {halIni} of {nHal}
          </span>
          <button
            type="button"
            className="btn btn--ghost btn--sm"
            disabled={halIni >= nHal}
            onClick={() => setHal(halIni + 1)}
            aria-label="Next"
          >
            ›
          </button>
        </div>
      ) : null}
      <table className="claimprop__tabel">
        <thead>
          <tr>
            {t.bernomor && <th>#</th>}
            {kolom.map((c, i) => (
              <th key={i} className={c.jenis === 'tombol' ? 'claimprop__th-aksi' : undefined}>
                {c.jenis !== 'tombol' ? (
                  // kolom angka: judul rata kanan seperti nilainya (work owner 08-10-2026 "ga sejajar header sama nilai")
                  c.kendali === 'angka' ? (
                    <span className="claimprop__angka">{c.label}</span>
                  ) : (
                    c.label
                  )
                ) : i === kolom.length - 1 && t.tambah ? (
                  <Tombol t={t.tambah} k={k} />
                ) : null}
              </th>
            ))}
            {selKepalaTambah && t.tambah && (
              <th className="claimprop__th-aksi">
                <Tombol t={t.tambah} k={k} />
              </th>
            )}
          </tr>
        </thead>
        <tbody>
          {baris.length === 0 && (
            <tr>
              <td className="muted" colSpan={lebar}>
                —
              </td>
            </tr>
          )}
          {baris.map((sel, i) => {
            const n = awal + i + 1
            return (
              <Fragment key={n}>
                <tr
                  className={rinci ? 'inbox__baris' : undefined}
                  aria-expanded={rinci ? buka === n : undefined}
                  onClick={rinci ? () => setBuka(buka === n ? null : n) : undefined}
                >
                  {t.bernomor && (
                    <td>
                      {rinci && (
                        <span className="claimprop__panah" aria-hidden="true">
                          {buka === n ? '▾' : '▸'}
                        </span>
                      )}
                      {n}
                    </td>
                  )}
                  {kolom.map((c, j) => {
                    const s = sel[j]
                    if (!s || !s.tampil) return <td key={j} />
                    const cel: Tata = { ...c, hanyaBaca: s.hanyaBaca, nonaktif: s.nonaktif }
                    return (
                      <td
                        key={j}
                        className={c.jenis === 'tombol' ? 'claimprop__td-aksi' : undefined}
                        onClick={(e) => {
                          // hanya isian yang dapat diubah menahan klik; sel hanya-baca ikut membuka panel baris
                          if (c.jenis === 'medan' && !cel.hanyaBaca && !cel.nonaktif && c.kendali !== 'tampil') {
                            e.stopPropagation()
                          }
                        }}
                      >
                        {c.jenis === 'tombol' ? (
                          <Tombol t={cel} k={k} indeks={n} />
                        ) : (
                          <Medan t={cel} k={k} indeks={n} jalur={jalurBaris(t.jalur ?? '', n, c.jalur ?? '')} />
                        )}
                      </td>
                    )
                  })}
                  {selKepalaTambah && <td />}
                </tr>
                {rinci && buka === n && (
                  <tr className="claimprop__baris-rinci">
                    <td colSpan={lebar}>
                      <BarisAdjustment.Provider value={n}>{rinci.isi(n)}</BarisAdjustment.Provider>
                    </td>
                  </tr>
                )}
              </Fragment>
            )
          })}
        </tbody>
      </table>
      {(t.kaki ?? []).length > 0 && (
        <div className="claimprop__grid-kaki">
          {ratakan(t.kaki ?? []).map((m, i) => {
            if (m.jenis === 'medan') {
              return (
                <span key={i} className="claimprop__sebaris-medan">
                  {m.label && <span className="claimprop__label-sebaris">{m.label}</span>}
                  <Medan t={m} k={k} jalur={m.jalur ?? ''} id={idMedan(m.jalur ?? '')} />
                </span>
              )
            }
            if (m.jenis === 'tombol') return <Tombol key={i} t={m} k={k} />
            if (m.jenis === 'label') {
              return (
                <span key={i} className="claimprop__label-sebaris">
                  {m.label}
                </span>
              )
            }
            return null
          })}
        </div>
      )}
    </div>
  )
}

/** Isi satu kartu / kolom / tab / modal: "Stacked with labels left". */
export default function TataView({ tata, k }: { tata: readonly Tata[]; k: KonteksTata }) {
  return (
    <div className="claimprop__tumpuk">
      {susunIsi(tata).map((u, i) => {
        switch (u.jenis) {
          case 'medan':
            return <BarisMedan key={i} t={u.t} tombol={u.tombol} k={k} />
          case 'tombol':
            return (
              <div key={i} className={'claimprop__tombol' + (u.akhir ? ' claimprop__tombol--akhir' : '')}>
                {u.tombol.map((b, j) => (
                  <Tombol key={j} t={b} k={k} />
                ))}
              </div>
            )
          case 'judul':
            return (
              <h4 key={i} className="claimprop__subjudul">
                {u.label}
              </h4>
            )
          case 'grid':
            return <Grid key={i} t={u.t} k={k} />
          case 'letak':
            if (u.t.letak === 'dua') return <Dua key={i} t={u.t} k={k} />
            if (u.t.letak === 'tab') return <Tab key={i} t={u.t} k={k} />
            if (u.t.letak === 'judul') return <Kepala key={i} t={u.t} k={k} />
            if (u.t.letak === 'tabel') return <TabelTetap key={i} t={u.t} k={k} />
            return <Sebaris key={i} t={u.t} k={k} />
          default:
            return (
              <section key={i} className="claimprop__sub">
                <h4 className="claimprop__subjudul">{u.t.label}</h4>
                <TataView tata={u.t.anak ?? []} k={k} />
              </section>
            )
        }
      })}
    </div>
  )
}

/** Tingkat layar kasus: kepala tengah, kartu `panel` bertitel, baris aksi. */
export function LayarTata({ tata, k }: { tata: readonly Tata[]; k: KonteksTata }) {
  return (
    <>
      {susunLayar(tata).map((b, i) => {
        if (b.jenis === 'kepala') return <Kepala key={i} t={b.t} k={k} />
        if (b.jenis === 'aksi') {
          return (
            <div key={i} className="claimprop__aksi">
              {b.tombol.map((t, j) => (
                <Tombol key={j} t={t} k={k} />
              ))}
            </div>
          )
        }
        return (
          <section key={i} className="panel">
            {b.judul && <h3 className="panel__title">{b.judul}</h3>}
            <TataView tata={b.isi} k={k} />
          </section>
        )
      })}
    </>
  )
}

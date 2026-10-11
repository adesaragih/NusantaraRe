// Disalin dari `modul/claimnonprop/frontend/components/TataView.tsx` (pola, bukan impor; asal Claim Prop): keputusan
// work owner yang disebut di bawah adalah keputusan layar Claim Prop yang ditiru Claim Fac In (prompt Claim Fac In tahap 1).
// Renderer pohon tata Claim Fac In. Satu komponen untuk semua section: kondisi tampil / hanya-baca / nonaktif sudah
// dievaluasi server; di sini hanya pemetaan kendali Pega ke elemen HTML. Medan berkendali aksi (refresh ber-activity)
// memanggil server saat nilainya berubah; medan tanpa aksi (postValue) hanya mengubah nilai lokal yang ikut terkirim
// pada aksi berikutnya.
//
// Rupa = layout lama Pega dirapikan (work owner 08-10-2026 "berikut layout lama, ikuti dan rapihkan") dengan kulit
// Kelola User: medan "Stacked with labels left" (label di kiri, nilai di kanan; hanya-baca = teks), `dua` = Inline grid
// double, `sebaris` = Inline, `tab` = layout group Tab, `judul` = kepala tengah, tombol ikon dari XML (pi-plus /
// pi-trash / pi-pencil / pi-check), grid ber-paging (pyGridPaginator). Pengelompokan di `susun.ts`.
//
// Claim Fac In: panel baris grid masterDetail UMUM (`rincian.ts`) dirender REKURSIF di bawah barisnya. Setiap unsur
// mengirim aksinya beralamat `konteks` = panel / modal tempatnya dirender ('' = layar utama) dan `indeks` = baris grid
// tempat sel / tombolnya berada di dalam konteks itu (0 = bukan sel grid, juga tombol Add kepala grid) -
// `services.PermintaanAksi`.

import { Fragment, useEffect, useRef, useState } from 'react'

import { PilihSaring } from '../../../../inti/frontend/components/ui/pilihSaring'
import type { Halaman, Pilihan, Tata } from '../api'
import { tampilTanggal } from '../ketikTanggal'
import { ambil, hanyaAngka, jalurBaris, sumberServer, tampilAngka } from '../nilai'
import InputAngka from './InputAngka'
import InputTanggal from './InputTanggal'
import { barisTerbuka, panelBaris } from './rincian'
import { jumlahHalaman, potongHalaman, ratakan, susunIsi, susunLayar, type Butir } from './susun'
import { butirSaring, teksTerpilih } from './tampilanNama'

export interface KonteksTata {
  h: Halaman
  /** postValue: ubah nilai lokal saja. */
  ubah: (jalur: string, nilai: string) => void
  /**
   * Ubah nilai lalu jalankan aksi server (refresh ber-activity / tombol). `indeks` = baris grid tempat sel / tombol
   * berada di dalam konteksnya (0 = bukan sel grid); `konteks` = kunci panel baris / modal ('' = layar utama); `param`
   * = pilihan autocomplete (dibaca ulang server).
   */
  aksi: (aksi: string, indeks: number, ubahan: Record<string, string>, konteks: string, param?: string) => void
  opsi: (sumber: string, konteks: string, indeks: number) => Pilihan[]
  /** Muat daftar pilihan server (autocomplete ber-saringan, dropdown bersumber server). */
  saran?: (sumber: string, konteks: string, indeks: number, cari: string) => void
  pesanMedan: Record<string, string[]>
  sibuk: boolean
  /** Konteks tempat unsur dirender: '' layar utama, kunci panel baris (`models.KunciPanel`), atau kunci modal. */
  konteks: string
  /** Panel baris grid masterDetail (`Layar.panel`); tanpa = tidak ada baris yang dapat dibuka. */
  panel?: Readonly<Record<string, Tata[] | undefined>>
  /** Penutup konteks ini (modal / panel baris): tombol tanpa aksi = Cancel. */
  tutup?: () => void
  /** Bertambah setiap layar baru diterima - dropdown bersumber server membaca ulang daftarnya. */
  versi?: number
}

/** Konteks unsur di dalam panel baris / modal `kunci`: aksinya beralamat ke sana; `tutup` = penutup tombol tanpa aksi. */
export function konteksDi(k: KonteksTata, kunci: string, tutup?: () => void): KonteksTata {
  return { ...k, konteks: kunci, tutup }
}

/** Klik tombol: tanpa aksi = penutup konteksnya (Cancel modal / panel); selainnya aksi beralamat konteks + baris grid. */
export function klikTombol(t: Pick<Tata, 'aksi'>, k: KonteksTata, indeks: number): void {
  if (!t.aksi) {
    k.tutup?.()
    return
  }
  k.aksi(t.aksi, indeks, {}, k.konteks)
}

/**
 * Ubah nilai medan: medan beraksi (refresh ber-activity) memanggil server bila `segera`; selainnya nilai lokal saja
 * (postValue). `param` = nilai pilihan autocomplete (server membaca ulang barisnya).
 */
export function gantiMedan(
  t: Pick<Tata, 'aksi'>,
  k: KonteksTata,
  jalur: string,
  indeks: number,
  baru: string,
  segera: boolean,
  param = '',
): void {
  if (t.aksi && segera) k.aksi(t.aksi, indeks, { [jalur]: baru }, k.konteks, param)
  else k.ubah(jalur, baru)
}

const idMedan = (jalur: string) => `claimfacin-${jalur.replace(/[^A-Za-z0-9]/g, '-')}`

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
 * Kendali satu medan. `indeks` = baris grid tempat sel berada (0 = bukan sel grid); `sel` = isian ringkas sel tabel.
 * Hanya-baca = teks nilai (layout lama), isian = kotak `field__input` kulit Kelola User.
 */
function Medan({
  t,
  k,
  indeks = 0,
  sel = false,
  jalur,
  id,
}: {
  t: Tata
  k: KonteksTata
  indeks?: number
  sel?: boolean
  jalur: string
  id?: string
}) {
  const v = ambil(k.h, jalur)
  // Nilai saat fokus: textarea beraksi memanggil server hanya bila isinya berubah (event `change` XML).
  const awal = useRef<string | null>(null)
  // `tampil` (pxDisplayText) selalu teks - server pun tidak pernah menerimanya (`models.MedanTerbuka`).
  const kunci = !!(t.hanyaBaca || t.nonaktif) || t.kendali === 'tampil'
  const sumber = t.sumber ?? ''
  const { konteks, saran, versi } = k
  // Dropdown bersumber server (Choose Currency panel adjustment): daftar dibaca saat medan tampil dan sesudah setiap
  // layar baru (daftar bergantung pada baris panelnya).
  const muat = (t.kendali === 'pilih' || t.kendali === 'radio') && sumberServer(sumber)
  useEffect(() => {
    if (muat) saran?.(sumber, konteks, indeks, '')
  }, [muat, sumber, konteks, indeks, saran, versi])
  const opsi = () => k.opsi(sumber, konteks, indeks)
  const ganti = (baru: string, segera: boolean, param = '') => gantiMedan(t, k, jalur, indeks, baru, segera, param)
  const angka = t.kendali === 'angka'
  const kelas = ['field__input', sel && 'claimfacin__input--sel', angka && 'claimfacin__input--angka']
    .filter(Boolean)
    .join(' ')

  if (t.kendali === 'centang') {
    return (
      <span className="claimfacin__centang">
        <input
          id={id}
          type="checkbox"
          aria-label={t.label || undefined}
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
      <span id={id} className={angka ? 'claimfacin__nilai claimfacin__angka' : 'claimfacin__nilai'}>
        {teks}
      </span>
    )
  }
  if (t.kendali === 'radio' && !sel) {
    return (
      <span id={id} className="claimfacin__radio" role="radiogroup">
        {opsi().map((o) => (
          <label key={o.nilai}>
            <input
              type="radio"
              name={idMedan(`${konteks}|${jalur}`)}
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
          onBlur={(berubah) => t.aksi && berubah && ganti(v, true)}
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
          onFocus={(e) => (awal.current = e.target.value)}
          onBlur={(e) => t.aksi && e.target.value !== awal.current && ganti(hanyaAngka(e.target.value), true)}
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
          className="field__input claimfacin__area"
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
        <select
          id={id}
          className={kelas}
          aria-label={sel ? t.label || undefined : undefined}
          value={v}
          disabled={k.sibuk}
          onChange={(e) => ganti(e.target.value, true)}
        >
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
    case 'otomatis':
      // Dropdown yang dapat dicari (keputusan work owner 08-10-2026: "model dropdown yang bisa di search"): PilihSaring
      // inti, saringan ke server lewat `saran`; nilai berubah hanya saat butir dipilih - juga di sel grid (item objek,
      // coverage, okupasi). Pilihan dikirim sebagai nilai medan DAN `param` aksi (server membaca ulang barisnya). Medan
      // ber-`tampilan` (Consultant / Adjuster) memilih dan menampilkan nama saja, nilai tetap ID (tampilanNama.ts).
      return (
        <span className="claimfacin__saring">
          <PilihSaring
            label={t.label ?? ''}
            sembunyikanLabel
            value={v}
            teksTerpilih={teksTerpilih(t, v, (j) => ambil(k.h, j), opsi())}
            opsi={butirSaring(t, opsi())}
            onCari={(kata) => saran?.(sumber, konteks, indeks, kata)}
            onPilih={(o) => ganti(o.value, true, o.value)}
          />
        </span>
      )
    default:
      return (
        <input
          id={id}
          className={kelas}
          value={v}
          disabled={k.sibuk}
          onChange={(e) => k.ubah(jalur, e.target.value)}
          onFocus={(e) => (awal.current = e.target.value)}
          onBlur={(e) => t.aksi && e.target.value !== awal.current && ganti(e.target.value, true)}
        />
      )
  }
}

/**
 * Tombol: ikon dari XML = tombol ikon (label jadi keterangan); Save / Submit = tombol utama (Kelola User: Simpan);
 * selainnya tombol bertepi. Tanpa label dan tanpa ikon = pengaturan. `utama` = tombol kaki modal. `indeks` = baris grid
 * tempat tombol berada (0 = bukan sel grid / tombol Add kepala grid).
 */
export function Tombol({ t, k, indeks = 0, utama }: { t: Tata; k: KonteksTata; indeks?: number; utama?: boolean }) {
  const primer = utama ?? (indeks === 0 && !t.ikon && (t.label === 'Save' || t.label === 'Submit'))
  const kelas = [
    'btn',
    primer ? 'btn--primary' : 'btn--ghost',
    (indeks > 0 || t.ikon) && 'btn--sm',
    t.ikon && 'claimfacin__ikon',
    t.ikon === 'hapus' && 'claimfacin__ikon--hapus',
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
        klikTombol(t, k, indeks)
      }}
    >
      {t.ikon ? <Ikon jenis={t.ikon} /> : t.label ? t.label : '⚙'}
    </button>
  )
}

/** Satu baris "Stacked with labels left": label di kiri, nilai + satuan + tombol menempel di kanan. */
function BarisMedan({ t, tombol, k }: { t: Butir; tombol: readonly Tata[]; k: KonteksTata }) {
  const jalur = t.jalur ?? ''
  const id = idMedan(`${k.konteks}|${jalur}`)
  const pesan = k.pesanMedan[jalur] ?? []
  return (
    <div className={'claimfacin__baris' + (t.label ? '' : ' claimfacin__baris--tanpa-label')} title={t.catatan}>
      <label className="claimfacin__label-medan" htmlFor={id}>
        {t.label}
        {t.wajib && <span className="field__req">*</span>}
      </label>
      <div className="claimfacin__isi-medan">
        <Medan t={t} k={k} jalur={jalur} id={id} />
        {t.satuan && <span className="claimfacin__satuan">{t.satuan}</span>}
        {tombol.map((b, i) => (
          <Tombol key={i} t={b} k={k} />
        ))}
        {pesan.length > 0 && <span className="field__error">{pesan.join('; ')}</span>}
      </div>
    </div>
  )
}

/** Layout Inline: anak sebaris; label bagian = label baris. */
function Sebaris({ t, k }: { t: Tata; k: KonteksTata }) {
  const isi = ratakan(t.anak ?? [])
  if (isi.length === 0) return null
  const baris = (
    <div className="claimfacin__sebaris">
      {isi.map((u, i) => {
        if (u.jenis === 'medan') {
          const jalur = u.jalur ?? ''
          const id = idMedan(`${k.konteks}|${jalur}`)
          const pesan = k.pesanMedan[jalur] ?? []
          return (
            <span key={i} className="claimfacin__sebaris-medan" title={u.catatan}>
              {u.label && (
                <label className="claimfacin__label-sebaris" htmlFor={id}>
                  {u.label}
                  {u.wajib && <span className="field__req">*</span>}
                </label>
              )}
              <Medan t={u} k={k} jalur={jalur} id={id} />
              {u.satuan && <span className="claimfacin__satuan">{u.satuan}</span>}
              {pesan.length > 0 && <span className="field__error">{pesan.join('; ')}</span>}
            </span>
          )
        }
        if (u.jenis === 'tombol') return <Tombol key={i} t={u} k={k} />
        if (u.jenis === 'label') {
          return (
            <span key={i} className="claimfacin__satuan">
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
    <div className="claimfacin__baris">
      <span className="claimfacin__label-medan">{t.label}</span>
      <div className="claimfacin__isi-medan">{baris}</div>
    </div>
  )
}

/** Layout Inline grid double: setiap anak satu sel; bagian tanpa judul = satu kolom. */
function Dua({ t, k }: { t: Tata; k: KonteksTata }) {
  return (
    <div className="claimfacin__dua">
      {(t.anak ?? []).map((c, i) => (
        <div key={i} className="claimfacin__sel">
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
    <div className="claimfacin__tab">
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
      <div role="tabpanel" className="claimfacin__tab-isi">
        {/* key per tab: grid tab lain (Estimation vs Adjustment, keduanya ObjectList) tidak mewarisi baris terbukanya */}
        <TataView key={i} tata={tab[i]?.anak ?? []} k={k} />
      </div>
    </div>
  )
}

/**
 * Layout bebas berkolom XML (`LetakTabel`): baris pertama = judul kolom, sel medan tanpa label; judul kolom berisi angka
 * rata kanan seperti nilainya.
 */
function TabelTetap({ t, k }: { t: Tata; k: KonteksTata }) {
  const [kepala, ...isi] = (t.anak ?? []).filter((b) => b.jenis === 'bagian')
  if (!kepala) return null
  const angka = (j: number) => isi.some((b) => b.anak?.[j]?.kendali === 'angka')
  const sel = (c: Tata) => {
    if (c.jenis === 'medan') return <Medan t={c} k={k} sel jalur={c.jalur ?? ''} />
    return c.jenis === 'label' ? c.label : null
  }
  return (
    <table className="claimfacin__tabel claimfacin__tabel--tetap">
      <thead>
        <tr>
          {(kepala.anak ?? []).map((c, j) => (
            <th key={j}>{angka(j) ? <span className="claimfacin__angka">{c.label}</span> : c.label}</th>
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

/** Kepala layar: judul tahap di tengah, teks di bawahnya. */
function Kepala({ t, k }: { t: Tata; k: KonteksTata }) {
  const isi = ratakan(t.anak ?? [])
  return (
    <div className="claimfacin__kepala">
      {isi.map((u, i) => {
        if (u.jenis === 'label') {
          return i === 0 ? (
            <h3 key={i} className="claimfacin__kepala-judul">
              {u.label}
            </h3>
          ) : (
            <div key={i} className="claimfacin__kepala-teks">
              {u.label}
            </div>
          )
        }
        if (u.jenis === 'medan') {
          return (
            <div key={i} className="claimfacin__kepala-teks">
              <span className="claimfacin__label-sebaris">{u.label}</span> <Medan t={u} k={k} jalur={u.jalur ?? ''} />
            </div>
          )
        }
        return null
      })}
    </div>
  )
}

/**
 * Isi panel baris `kunci` di bawah barisnya - REKURSIF: grid di dalamnya ber-`rincian` membuka panelnya sendiri. Unsur
 * di dalamnya beralamat `konteks` = `kunci`; tombol tanpa aksi menutup panel.
 */
export function PanelBaris({ kunci, k, onTutup }: { kunci: string; k: KonteksTata; onTutup?: () => void }) {
  return (
    <div className="claimfacin__rinci">
      <TataView tata={k.panel?.[kunci] ?? []} k={konteksDi(k, kunci, onTutup)} />
    </div>
  )
}

function Grid({ t, k }: { t: Tata; k: KonteksTata }) {
  const kolom = t.kolom ?? []
  const semua = t.baris ?? []
  // Panel rinci: satu baris terbuka; klik di mana pun pada baris membuka / menutup. Baris baru (Add) langsung terbuka.
  const nBaris = semua.length
  const adaRinci = !!t.rincian
  const [buka, setBuka] = useState<number | null>(null)
  const [hal, setHal] = useState(1)
  const nLalu = useRef(nBaris)
  useEffect(() => {
    const lalu = nLalu.current
    nLalu.current = nBaris
    if (nBaris === lalu || !adaRinci) return
    setBuka((b) => barisTerbuka(b, lalu, nBaris))
    if (nBaris > lalu && t.perHalaman) setHal(jumlahHalaman(nBaris, t.perHalaman))
  }, [nBaris, adaRinci, t.perHalaman])
  const nHal = jumlahHalaman(semua.length, t.perHalaman)
  const halIni = Math.min(hal, nHal)
  const awal = t.perHalaman ? (halIni - 1) * t.perHalaman : 0
  const baris = potongHalaman(semua, t.perHalaman, halIni)
  // Add di sel kepala kolom tombol terakhir (pzPegaDefaultGridIcons); tanpa kolom tombol = sel kepala tambahan
  const kolomAkhirTombol = kolom[kolom.length - 1]?.jenis === 'tombol'
  const selKepalaTambah = !!t.tambah && !kolomAkhirTombol
  // Grid ber-panel tanpa nomor baris: kolom panah di depan (penanda baris yang dapat dibuka).
  const kolomPanah = adaRinci && !t.bernomor
  const kolomDepan = !!t.bernomor || kolomPanah
  const lebar = kolom.length + (kolomDepan ? 1 : 0) + (selKepalaTambah ? 1 : 0)
  return (
    <div className="claimfacin__grid" title={t.catatan}>
      {t.label && <h4 className="claimfacin__subjudul">{t.label}</h4>}
      {t.perHalaman && semua.length > t.perHalaman ? (
        <div className="claimfacin__pager">
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
      <table className="claimfacin__tabel">
        <thead>
          <tr>
            {kolomDepan && <th>{t.bernomor ? '#' : ''}</th>}
            {kolom.map((c, i) => (
              <th key={i} className={c.jenis === 'tombol' ? 'claimfacin__th-aksi' : undefined}>
                {c.jenis !== 'tombol' ? (
                  // kolom angka: judul rata kanan seperti nilainya (work owner 08-10-2026 "ga sejajar header sama nilai")
                  c.kendali === 'angka' ? (
                    <span className="claimfacin__angka">{c.label}</span>
                  ) : (
                    c.label
                  )
                ) : i === kolom.length - 1 && t.tambah ? (
                  <Tombol t={t.tambah} k={k} />
                ) : null}
              </th>
            ))}
            {selKepalaTambah && t.tambah && (
              <th className="claimfacin__th-aksi">
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
            const kunci = panelBaris(t, n, k.panel)
            const terbuka = kunci !== null && buka === n
            return (
              <Fragment key={n}>
                <tr
                  className={kunci ? 'inbox__baris' : undefined}
                  aria-expanded={kunci ? terbuka : undefined}
                  onClick={kunci ? () => setBuka(terbuka ? null : n) : undefined}
                >
                  {kolomDepan && (
                    <td>
                      {kunci && (
                        <span className="claimfacin__panah" aria-hidden="true">
                          {terbuka ? '▾' : '▸'}
                        </span>
                      )}
                      {t.bernomor && n}
                    </td>
                  )}
                  {kolom.map((c, j) => {
                    const s = sel[j]
                    if (!s || !s.tampil) return <td key={j} />
                    const cel: Tata = { ...c, hanyaBaca: s.hanyaBaca, nonaktif: s.nonaktif }
                    return (
                      <td
                        key={j}
                        className={c.jenis === 'tombol' ? 'claimfacin__td-aksi' : undefined}
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
                          <Medan t={cel} k={k} indeks={n} sel jalur={jalurBaris(t.jalur ?? '', n, c.jalur ?? '')} />
                        )}
                      </td>
                    )
                  })}
                  {selKepalaTambah && <td />}
                </tr>
                {kunci !== null && terbuka && (
                  <tr className="claimfacin__baris-rinci">
                    <td colSpan={lebar}>
                      <PanelBaris kunci={kunci} k={k} onTutup={() => setBuka(null)} />
                    </td>
                  </tr>
                )}
              </Fragment>
            )
          })}
        </tbody>
      </table>
      {(t.kaki ?? []).length > 0 && (
        <div className="claimfacin__grid-kaki">
          {ratakan(t.kaki ?? []).map((m, i) => {
            if (m.jenis === 'medan') {
              const jalur = m.jalur ?? ''
              return (
                <span key={i} className="claimfacin__sebaris-medan">
                  {m.label && <span className="claimfacin__label-sebaris">{m.label}</span>}
                  <Medan t={m} k={k} jalur={jalur} id={idMedan(`${k.konteks}|${jalur}`)} />
                </span>
              )
            }
            if (m.jenis === 'tombol') return <Tombol key={i} t={m} k={k} />
            if (m.jenis === 'label') {
              return (
                <span key={i} className="claimfacin__label-sebaris">
                  {m.label}
                </span>
              )
            }
            return <TataView key={i} tata={[m]} k={k} />
          })}
        </div>
      )}
    </div>
  )
}

/** Key React grid: posisi + jalur + prefiks panel (`rincian`). */
export function kunciGrid(i: number, t: Pick<Tata, 'jalur' | 'rincian'>): string {
  return `${i}|${t.jalur ?? ''}|${t.rincian ?? ''}`
}

/** Isi satu kartu / kolom / tab / modal / panel: "Stacked with labels left". */
export default function TataView({ tata, k }: { tata: readonly Tata[]; k: KonteksTata }) {
  return (
    <div className="claimfacin__tumpuk">
      {susunIsi(tata).map((u, i) => {
        switch (u.jenis) {
          case 'medan':
            return <BarisMedan key={i} t={u.t} tombol={u.tombol} k={k} />
          case 'tombol':
            return (
              <div key={i} className={'claimfacin__tombol' + (u.akhir ? ' claimfacin__tombol--akhir' : '')}>
                {u.tombol.map((b, j) => (
                  <Tombol key={j} t={b} k={k} />
                ))}
              </div>
            )
          case 'judul':
            return (
              <h4 key={i} className="claimfacin__subjudul">
                {u.label}
              </h4>
            )
          case 'grid':
            // key berjalur + prefiks panel: grid lain di posisi yang sama tidak mewarisi baris terbuka / halamannya
            return <Grid key={kunciGrid(i, u.t)} t={u.t} k={k} />
          case 'letak':
            if (u.t.letak === 'dua') return <Dua key={i} t={u.t} k={k} />
            if (u.t.letak === 'tab') return <Tab key={i} t={u.t} k={k} />
            if (u.t.letak === 'judul') return <Kepala key={i} t={u.t} k={k} />
            if (u.t.letak === 'tabel') return <TabelTetap key={i} t={u.t} k={k} />
            return <Sebaris key={i} t={u.t} k={k} />
          default:
            return (
              <section key={i} className="claimfacin__sub">
                <h4 className="claimfacin__subjudul">{u.t.label}</h4>
                <TataView tata={u.t.anak ?? []} k={k} />
              </section>
            )
        }
      })}
    </div>
  )
}

/** Tingkat layar: kepala tengah, kartu `panel` bertitel, baris aksi. */
export function LayarTata({ tata, k }: { tata: readonly Tata[]; k: KonteksTata }) {
  return (
    <>
      {susunLayar(tata).map((b, i) => {
        if (b.jenis === 'kepala') return <Kepala key={i} t={b.t} k={k} />
        if (b.jenis === 'aksi') {
          return (
            <div key={i} className="claimfacin__aksi">
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

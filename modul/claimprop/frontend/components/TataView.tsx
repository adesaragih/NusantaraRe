// Renderer pohon tata Claim Prop. Satu komponen untuk semua section: kondisi tampil / hanya-baca / nonaktif sudah
// dievaluasi server; di sini hanya pemetaan kendali Pega ke elemen HTML. Medan berkendali aksi (refresh ber-activity)
// memanggil server saat nilainya berubah; medan tanpa aksi (postValue) hanya mengubah nilai lokal yang ikut terkirim
// pada aksi berikutnya.
//
// Rupa = Kelola User (keputusan work owner 08-10-2026): kartu `panel` + `panel__title`, isi `form-grid`, medan `field`
// berlabel di atas kotak `field__input`; medan hanya-baca tetap berkotak (`field__input--readonly`), sel grid ringkas.
// Pengelompokannya di `susun.ts`.

import { Fragment, useState, type ReactNode } from 'react'

import type { Halaman, Pilihan, Tata } from '../api'
import { ambil, dariInputWaktu, jalurBaris, keInputWaktu, tampilAngka } from '../nilai'
import { susunIsi, susunLayar, type Butir } from './susun'

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
  /** Panel rinci baris grid (expand pane AdjustmentDetail). */
  rincian?: (jalurDaftar: string, n: number) => ReactNode
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
      return v.slice(0, 10)
    case 'tanggal-waktu':
      return keInputWaktu(v).replace('T', ' ')
    case 'pilih':
    case 'radio':
      return opsi().find((o) => o.nilai === v)?.label ?? v
    default:
      return v
  }
}

/**
 * Kendali satu medan. `sel` = sel grid (isian ringkas, nilai hanya-baca berupa teks); selainnya kotak penuh Kelola
 * User, hanya-baca = kotak `field__input--readonly`.
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
  const sel = indeks > 0
  const kunci = !!(t.hanyaBaca || t.nonaktif)
  const n = indeks || indeksDari(jalur)
  const opsi = () => k.opsi(t.sumber ?? '', n)
  const ganti = (baru: string, segera: boolean) => {
    if (t.aksi && segera) k.aksi(t.aksi, n, { [jalur]: baru })
    else k.ubah(jalur, baru)
  }
  const angka = t.kendali === 'angka'
  const kelas = (tambahan = '') =>
    ['field__input', sel && 'claimprop__input--sel', angka && 'claimprop__input--angka', tambahan]
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
    const teks = teksTampil(t, v, opsi)
    if (sel) return <span className={angka ? 'claimprop__angka' : undefined}>{teks}</span>
    if (t.kendali === 'area') {
      return <textarea id={id} className="field__input field__input--readonly claimprop__area" readOnly value={teks} />
    }
    return <input id={id} className={kelas('field__input--readonly')} readOnly tabIndex={-1} value={teks} />
  }
  switch (t.kendali) {
    case 'angka':
      return (
        <input
          id={id}
          className={kelas()}
          inputMode="decimal"
          value={v}
          disabled={k.sibuk}
          onChange={(e) => k.ubah(jalur, e.target.value.replace(/,/g, ''))}
          onBlur={(e) => t.aksi && ganti(e.target.value.replace(/,/g, ''), true)}
        />
      )
    case 'tanggal':
      return (
        <input
          id={id}
          type="date"
          className={kelas()}
          value={v.slice(0, 10)}
          disabled={k.sibuk}
          onChange={(e) => ganti(e.target.value, true)}
        />
      )
    case 'tanggal-waktu':
      return (
        <input
          id={id}
          type="datetime-local"
          className={kelas()}
          value={keInputWaktu(v)}
          disabled={k.sibuk}
          onChange={(e) => ganti(dariInputWaktu(e.target.value), true)}
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
          onBlur={(e) => t.aksi && ganti(e.target.value, true)}
        />
      )
    case 'pilih':
    case 'radio': {
      const daftar = opsi()
      return (
        <select id={id} className={kelas()} value={v} disabled={k.sibuk} onChange={(e) => ganti(e.target.value, true)}>
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
      const idDaftar = `${idMedan(jalur)}-saran`
      return (
        <>
          <input
            id={id}
            className={kelas()}
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
          className={kelas()}
          value={v}
          disabled={k.sibuk}
          onChange={(e) => k.ubah(jalur, e.target.value)}
          onBlur={(e) => t.aksi && ganti(e.target.value, true)}
        />
      )
  }
}

/** Tombol utama (Kelola User: Simpan) = Save / Submit; selainnya tombol bertepi. Ikon tanpa label = pengaturan. */
function Tombol({ t, k, indeks = 0 }: { t: Tata; k: KonteksTata; indeks?: number }) {
  const utama = indeks === 0 && (t.label === 'Save' || t.label === 'Submit')
  const kelas = ['btn', utama ? 'btn--primary' : 'btn--ghost', indeks > 0 && 'btn--sm'].filter(Boolean).join(' ')
  return (
    <button
      type="button"
      className={kelas}
      disabled={t.nonaktif || k.sibuk}
      title={t.catatan ?? t.id}
      aria-label={t.label ? undefined : t.id}
      onClick={() => k.aksi(t.aksi ?? t.id ?? '', indeks)}
    >
      {t.label ? t.label : '⚙'}
    </button>
  )
}

/** Satu sel form-grid: label di atas, kotak isian + satuan + tombol menempel di kanan. */
function SelMedan({ t, tombol, k }: { t: Butir; tombol: readonly Tata[]; k: KonteksTata }) {
  const jalur = t.jalur ?? ''
  const id = idMedan(jalur)
  const pesan = k.pesanMedan[jalur] ?? []
  return (
    <div className="field" title={t.catatan}>
      <label className="field__label" htmlFor={id}>
        {t.label ? t.label : <>&nbsp;</>}
        {t.wajib && <span className="field__req">*</span>}
      </label>
      <div className="claimprop__baris-isian">
        <Medan t={t} k={k} jalur={jalur} id={id} />
        {t.satuan && <span className="claimprop__satuan">{t.satuan}</span>}
        {tombol.map((b, i) => (
          <Tombol key={i} t={b} k={k} />
        ))}
      </div>
      {pesan.length > 0 && <p className="field__error">{pesan.join('; ')}</p>}
    </div>
  )
}

function Grid({ t, k }: { t: Tata; k: KonteksTata }) {
  const [buka, setBuka] = useState<number | null>(null)
  const kolom = t.kolom ?? []
  const baris = t.baris ?? []
  return (
    <div className="claimprop__grid" title={t.catatan}>
      {t.tambah && (
        <div className="claimprop__grid-alat">
          <Tombol t={t.tambah} k={k} />
        </div>
      )}
      <table>
        <thead>
          <tr>
            {t.bernomor && <th>#</th>}
            {kolom.map((c, i) => (
              <th key={i}>{c.jenis === 'tombol' ? '' : c.label}</th>
            ))}
          </tr>
        </thead>
        <tbody>
          {baris.length === 0 && (
            <tr>
              <td className="muted" colSpan={kolom.length + (t.bernomor ? 1 : 0)}>
                —
              </td>
            </tr>
          )}
          {baris.map((sel, i) => (
            <Fragment key={i}>
              <tr
                className={k.rincian ? 'inbox__baris' : undefined}
                onClick={k.rincian ? () => setBuka(buka === i + 1 ? null : i + 1) : undefined}
              >
                {t.bernomor && <td>{i + 1}</td>}
                {kolom.map((c, j) => {
                  const s = sel[j]
                  if (!s || !s.tampil) return <td key={j} />
                  const cel: Tata = { ...c, hanyaBaca: s.hanyaBaca, nonaktif: s.nonaktif }
                  return (
                    <td
                      key={j}
                      onClick={(e) => {
                        if (c.jenis === 'medan') e.stopPropagation()
                      }}
                    >
                      {c.jenis === 'tombol' ? (
                        <Tombol t={cel} k={k} indeks={i + 1} />
                      ) : (
                        <Medan t={cel} k={k} indeks={i + 1} jalur={jalurBaris(t.jalur ?? '', i + 1, c.jalur ?? '')} />
                      )}
                    </td>
                  )
                })}
              </tr>
              {k.rincian && buka === i + 1 && (
                <tr>
                  <td colSpan={kolom.length + (t.bernomor ? 1 : 0)}>{k.rincian(t.jalur ?? '', i + 1)}</td>
                </tr>
              )}
            </Fragment>
          ))}
        </tbody>
      </table>
      {(t.kaki ?? []).length > 0 && (
        <div className="claimprop__grid-kaki">
          {(t.kaki ?? []).map((m, i) =>
            m.jenis === 'medan' ? (
              <span key={i} className="claimprop__kaki-nilai">
                {m.label && <span className="field__label">{m.label}</span>}
                <Medan t={m} k={k} jalur={m.jalur ?? ''} id={idMedan(m.jalur ?? '')} />
              </span>
            ) : m.jenis === 'tombol' ? (
              <Tombol key={i} t={m} k={k} />
            ) : null,
          )}
        </div>
      )}
    </div>
  )
}

/** Isi satu kartu / sub-bagian / modal: form-grid Kelola User. */
export default function TataView({ tata, k }: { tata: readonly Tata[]; k: KonteksTata }) {
  return (
    <div className="form-grid">
      {susunIsi(tata).map((u, i) => {
        switch (u.jenis) {
          case 'medan':
            return <SelMedan key={i} t={u.t} tombol={u.tombol} k={k} />
          case 'tombol':
            return (
              <div key={i} className={'field--lebar claimprop__tombol' + (u.akhir ? ' claimprop__tombol--akhir' : '')}>
                {u.tombol.map((b, j) => (
                  <Tombol key={j} t={b} k={k} />
                ))}
              </div>
            )
          case 'judul':
            return (
              <h4 key={i} className="field--lebar claimprop__subjudul">
                {u.label}
              </h4>
            )
          case 'grid':
            return (
              <div key={i} className="field--lebar">
                <Grid t={u.t} k={k} />
              </div>
            )
          default:
            return (
              <section key={i} className="field--lebar claimprop__sub">
                <h4 className="claimprop__subjudul">{u.t.label}</h4>
                <TataView tata={u.t.anak ?? []} k={k} />
              </section>
            )
        }
      })}
    </div>
  )
}

/** Tingkat layar kasus: kartu `panel` bertitel + baris aksi. */
export function LayarTata({ tata, k }: { tata: readonly Tata[]; k: KonteksTata }) {
  return (
    <>
      {susunLayar(tata).map((b, i) =>
        b.jenis === 'aksi' ? (
          <div key={i} className="claimprop__aksi">
            {b.tombol.map((t, j) => (
              <Tombol key={j} t={t} k={k} />
            ))}
          </div>
        ) : (
          <section key={i} className="panel">
            {b.judul && <h3 className="panel__title">{b.judul}</h3>}
            <TataView tata={b.isi} k={k} />
          </section>
        ),
      )}
    </>
  )
}

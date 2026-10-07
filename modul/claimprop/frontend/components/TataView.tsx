// Renderer pohon tata Claim Prop. Satu komponen untuk semua section: kondisi tampil / hanya-baca / nonaktif sudah
// dievaluasi server; di sini hanya pemetaan kendali Pega ke elemen HTML. Medan berkendali aksi (refresh ber-activity)
// memanggil server saat nilainya berubah; medan tanpa aksi (postValue) hanya mengubah nilai lokal yang ikut terkirim
// pada aksi berikutnya.

import { Fragment, useState, type ReactNode } from 'react'

import type { Halaman, Pilihan, Tata } from '../api'
import { ambil, dariInputWaktu, jalurBaris, keInputWaktu, tampilAngka } from '../nilai'

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

function Medan({ t, k, indeks = 0, jalur }: { t: Tata; k: KonteksTata; indeks?: number; jalur: string }) {
  const v = ambil(k.h, jalur)
  const kunci = t.hanyaBaca || t.nonaktif || k.sibuk
  const ganti = (baru: string, segera: boolean) => {
    if (t.aksi && segera) k.aksi(t.aksi, indeks || indeksDari(jalur), { [jalur]: baru })
    else k.ubah(jalur, baru)
  }
  const pesan = k.pesanMedan[jalur]
  // sel grid memakai isian ringkas; medan bagian memakai isian penuh Kelola User
  const kelasIsian = indeks > 0 ? 'field__input claimprop__input--sel' : 'field__input'
  let isi: ReactNode
  switch (t.kendali) {
    case 'tampil':
      isi = <span className="claimprop__tampil">{v}</span>
      break
    case 'angka':
      isi = kunci ? (
        <span className="claimprop__angka">{tampilAngka(v)}</span>
      ) : (
        <input
          className={'field__input claimprop__input--angka' + (indeks > 0 ? ' claimprop__input--sel' : '')}
          inputMode="decimal"
          value={v}
          onChange={(e) => k.ubah(jalur, e.target.value.replace(/,/g, ''))}
          onBlur={(e) => t.aksi && ganti(e.target.value.replace(/,/g, ''), true)}
        />
      )
      break
    case 'tanggal':
      isi = (
        <input
          type="date"
          className={kelasIsian}
          value={v.slice(0, 10)}
          disabled={kunci}
          onChange={(e) => ganti(e.target.value, true)}
        />
      )
      break
    case 'tanggal-waktu':
      isi = (
        <input
          type="datetime-local"
          className={kelasIsian}
          value={keInputWaktu(v)}
          disabled={kunci}
          onChange={(e) => ganti(dariInputWaktu(e.target.value), true)}
        />
      )
      break
    case 'area':
      isi = (
        <textarea
          className="field__input claimprop__area"
          value={v}
          readOnly={kunci}
          onChange={(e) => k.ubah(jalur, e.target.value)}
          onBlur={(e) => t.aksi && ganti(e.target.value, true)}
        />
      )
      break
    case 'centang':
      isi = (
        <input
          type="checkbox"
          checked={v === 'true'}
          disabled={kunci}
          onChange={(e) => ganti(e.target.checked ? 'true' : 'false', true)}
        />
      )
      break
    case 'pilih':
    case 'radio': {
      const opsi = k.opsi(t.sumber ?? '', indeks || indeksDari(jalur))
      const label = opsi.find((o) => o.nilai === v)?.label ?? v
      isi = kunci ? (
        <span className="claimprop__tampil">{label}</span>
      ) : (
        <select className={kelasIsian} value={v} onChange={(e) => ganti(e.target.value, true)}>
          <option value="">Choose</option>
          {opsi.map((o) => (
            <option key={o.nilai} value={o.nilai}>
              {o.label}
            </option>
          ))}
          {v !== '' && !opsi.some((o) => o.nilai === v) && <option value={v}>{v}</option>}
        </select>
      )
      break
    }
    case 'otomatis': {
      const idDaftar = `claimprop-${jalur.replace(/[^A-Za-z0-9]/g, '-')}`
      const opsi = k.opsi(t.sumber ?? '', indeks || indeksDari(jalur))
      isi = kunci ? (
        <span className="claimprop__tampil">{v}</span>
      ) : (
        <>
          <input
            className={kelasIsian}
            list={idDaftar}
            value={v}
            onFocus={() => k.saran?.(t.sumber ?? '', indeks || indeksDari(jalur), '')}
            onChange={(e) => {
              k.ubah(jalur, e.target.value)
              k.saran?.(t.sumber ?? '', indeks || indeksDari(jalur), e.target.value)
            }}
            onBlur={(e) => t.aksi && ganti(e.target.value, true)}
          />
          <datalist id={idDaftar}>
            {opsi.map((o) => (
              <option key={o.nilai + o.label} value={o.nilai}>
                {o.label}
              </option>
            ))}
          </datalist>
        </>
      )
      break
    }
    default:
      isi = kunci ? (
        <span className="claimprop__tampil">{v}</span>
      ) : (
        <input
          className={kelasIsian}
          value={v}
          onChange={(e) => k.ubah(jalur, e.target.value)}
          onBlur={(e) => t.aksi && ganti(e.target.value, true)}
        />
      )
  }
  return (
    <span className="claimprop__nilai" title={t.catatan}>
      {isi}
      {pesan && pesan.length > 0 && <span className="claimprop__pesan-medan">{pesan.join('; ')}</span>}
    </span>
  )
}

function Tombol({ t, k, indeks = 0 }: { t: Tata; k: KonteksTata; indeks?: number }) {
  return (
    <button
      type="button"
      className={indeks > 0 || !t.label ? 'btn btn--ghost btn--sm' : 'btn btn--sm'}
      disabled={t.nonaktif || k.sibuk}
      title={t.catatan ?? t.id}
      onClick={() => k.aksi(t.aksi ?? t.id ?? '', indeks)}
    >
      {t.label ? t.label : '⚙'}
    </button>
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
          <TataView tata={t.kaki ?? []} k={k} />
        </div>
      )}
    </div>
  )
}

/** Merender satu daftar unsur tata. */
export default function TataView({ tata, k }: { tata: readonly Tata[]; k: KonteksTata }) {
  return (
    <>
      {tata.map((t, i) => {
        switch (t.jenis) {
          case 'bagian':
            return t.label ? (
              <fieldset key={i} className="claimprop__bagian">
                <legend className="claimprop__judul-bagian">{t.label}</legend>
                <div className="claimprop__isi-bagian">
                  <TataView tata={t.anak ?? []} k={k} />
                </div>
              </fieldset>
            ) : (
              <div key={i} className="claimprop__isi-bagian">
                <TataView tata={t.anak ?? []} k={k} />
              </div>
            )
          case 'label':
            return (
              <div key={i} className="claimprop__label">
                {t.label}
              </div>
            )
          case 'tombol':
            return <Tombol key={i} t={t} k={k} />
          case 'grid':
            return (
              <Grid
                key={i}
                t={t}
                k={k.rincian && t.jalur === 'ClaimData.AdjustmentList' ? k : { ...k, rincian: undefined }}
              />
            )
          default:
            return (
              <label key={i} className="field claimprop__medan">
                <span className="field__label">
                  {t.label}
                  {t.wajib && <span className="field__req">*</span>}
                </span>
                <Medan t={t} k={k} jalur={t.jalur ?? ''} />
              </label>
            )
        }
      })}
    </>
  )
}

// Form Tambah / Ubah satu baris master (inti):
// - isian = kolom bukan turunan; ID diisi pengguna saat Tambah, baca-saja saat Ubah (ID dari rute); master ber-ID
//   otomatis (Accumulation) tidak menampilkan isian ID saat Tambah dan meminta `negara` (awalan ID);
// - kolom rujukan (meta `rujukan`) = `PilihRujukan` (data aktif saja);
// - kolom turunan (nama rujukan, jejak ubah) tampil baca-saja saat Ubah;
// - wajib dan lebar (byte) diperiksa sebelum dikirim; galat backend (400 / 409) tampil di kaki form.

import { useState } from 'react'

import { Field, Gagal, Modal } from '../components/ui/dasar'
import type { BarisMaster, KlienMaster, KolomMaster, MetaMaster } from './api'
import { labelKolomMaster, TEKS_MASTER as T } from './labels'
import PilihRujukan from './PilihRujukan'

/** Kunci masukan tambahan master ber-ID otomatis (bukan kolom). */
export const KUNCI_NEGARA = 'negara'
/** Lebar `negara` (= NATION.ID). */
export const LEBAR_NEGARA = 10

/** Panjang teks dalam byte UTF-8 (lebar kolom Oracle BYTE). */
export const panjangByte = (v: string) => new TextEncoder().encode(v).length

/** Medan yang menjadi ISIAN form. */
export function medanForm(meta: MetaMaster, baru: boolean): KolomMaster[] {
  const kolom = meta.kolom.filter((k) => !k.turunan && !(k.kunci === 'id' && (!baru || meta.idOtomatis)))
  const tambahan: KolomMaster[] =
    baru && meta.idOtomatis ? [{ kunci: KUNCI_NEGARA, kolom: KUNCI_NEGARA, lebar: LEBAR_NEGARA, wajib: true, turunan: false }] : []
  return [...tambahan, ...kolom]
}

/** Galat per medan (kosong = sah). */
export function periksa(meta: MetaMaster, baru: boolean, isi: Record<string, string>): Record<string, string> {
  const g: Record<string, string> = {}
  for (const k of medanForm(meta, baru)) {
    const v = (isi[k.kunci] ?? '').trim()
    if (k.wajib && v === '') g[k.kunci] = T.wajib
    else if (k.lebar > 0 && panjangByte(v) > k.lebar) g[k.kunci] = T.terlaluPanjang(k.lebar)
  }
  return g
}

/** Badan permintaan: medan form saja (turunan tidak dikirim), dipangkas. */
export function badan(meta: MetaMaster, baru: boolean, isi: Record<string, string>): Record<string, string> {
  return Object.fromEntries(medanForm(meta, baru).map((k) => [k.kunci, (isi[k.kunci] ?? '').trim()]))
}

export default function FormMaster({
  klien,
  meta,
  baris,
  onTutup,
  onTersimpan,
}: {
  klien: KlienMaster
  meta: MetaMaster
  /** null = Tambah. */
  baris: BarisMaster | null
  onTutup: () => void
  onTersimpan: (id: string) => void
}) {
  const baru = baris === null
  const [isi, setIsi] = useState<Record<string, string>>(() =>
    Object.fromEntries(meta.kolom.map((k) => [k.kunci, baris ? String(baris[k.kunci] ?? '') : ''])),
  )
  const [coba, setCoba] = useState(false)
  const [menyimpan, setMenyimpan] = useState(false)
  const [galat, setGalat] = useState<unknown>(null)
  const set = (k: string) => (v: string) => setIsi((x) => ({ ...x, [k]: v }))
  const salah = coba ? periksa(meta, baru, isi) : {}
  const idBaris = String(baris?.id ?? '')

  async function simpan() {
    setCoba(true)
    if (Object.keys(periksa(meta, baru, isi)).length > 0) return
    setMenyimpan(true)
    setGalat(null)
    try {
      const h = baru ? await klien.tambah(badan(meta, baru, isi)) : await klien.ubah(idBaris, badan(meta, baru, isi))
      onTersimpan(h.id)
    } catch (err) {
      setGalat(err)
    } finally {
      setMenyimpan(false)
    }
  }

  return (
    <Modal
      judul={baru ? T.judulTambah(meta.judul) : T.judulUbah(meta.judul)}
      onTutup={onTutup}
      onKirim={() => void simpan()}
      lebar
      aksi={
        <button type="submit" className="btn btn--primary" disabled={menyimpan}>
          {menyimpan ? T.menyimpan : T.simpan}
        </button>
      }
    >
      <div className="form-grid">
        {(!baru || meta.idOtomatis) && (
          <Field label={labelKolomMaster('id', 'ID')} value={baru ? T.idOtomatis : idBaris} onChange={() => {}} readOnly />
        )}
        {medanForm(meta, baru).map((k) => {
          const label = labelKolomMaster(k.kunci, k.kolom)
          const r = meta.rujukan.find((x) => x.kunci === k.kunci)
          if (r)
            return (
              <PilihRujukan
                key={k.kunci}
                klien={klien}
                rujukan={r}
                label={label}
                value={isi[k.kunci] ?? ''}
                onChange={set(k.kunci)}
                required={k.wajib}
                error={salah[k.kunci]}
              />
            )
          return (
            <Field key={k.kunci} label={label} value={isi[k.kunci] ?? ''} onChange={set(k.kunci)} required={k.wajib} error={salah[k.kunci]} />
          )
        })}
        {!baru &&
          meta.kolom
            .filter((k) => k.turunan && k.kunci !== 'id')
            .map((k) => (
              <Field key={k.kunci} label={`${labelKolomMaster(k.kunci, k.kolom)} (${T.turunan})`} value={isi[k.kunci] ?? ''} onChange={() => {}} readOnly />
            ))}
      </div>
      <Gagal galat={galat} />
    </Modal>
  )
}

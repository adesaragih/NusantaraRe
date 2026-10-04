// Form Tambah / Ubah satu baris master (rencana M-3 … M-5, kontrak 02):
// - medan = kolom master yang bukan turunan; kolom turunan tampil baca-saja saat Ubah (diisi backend);
// - ID: diisi pengguna saat Tambah, baca-saja saat Ubah (MD-9); master ber-ID otomatis (Accumulation) tidak menampilkan
//   ID saat Tambah dan meminta `negara` (awalan ID, MD-3);
// - kolom rujukan = `PilihRujukan` (hanya data aktif);
// - wajib dan lebar (byte) diperiksa sebelum dikirim; galat backend (400 / 409) tampil di kaki form.

import { useState } from 'react'

import { Field, Gagal, Modal } from '../../../../inti/frontend/components/ui/dasar'
import { tambahMaster, ubahMaster, type BarisMaster, type MetaMaster } from '../api'
import { labelKolom, RUJUKAN, TEKS } from '../labels'
import PilihRujukan from './PilihRujukan'

/** Kunci masukan tambahan Accumulation (bukan kolom). */
export const KUNCI_NEGARA = 'negara'
/** Lebar `negara` (= NATION.ID). */
export const LEBAR_NEGARA = 10

/** Panjang teks dalam byte UTF-8 (lebar kolom Oracle BYTE). */
export const panjangByte = (v: string) => new TextEncoder().encode(v).length

/**
 * Medan yang tampil sebagai ISIAN: kolom bukan turunan; ID hanya saat Tambah dan bila tidak otomatis (saat Ubah ID
 * tampil baca-saja dan diambil dari rute, MD-9).
 */
export function medanForm(meta: MetaMaster, baru: boolean) {
  const kolom = meta.kolom.filter((k) => !k.turunan && !(k.kunci === 'id' && (!baru || meta.idOtomatis)))
  const tambahan = baru && meta.idOtomatis ? [{ kunci: KUNCI_NEGARA, kolom: KUNCI_NEGARA, lebar: LEBAR_NEGARA, wajib: true, turunan: false }] : []
  return [...tambahan, ...kolom]
}

/** Galat per medan (kosong = sah). */
export function periksa(meta: MetaMaster, baru: boolean, isi: Record<string, string>): Record<string, string> {
  const g: Record<string, string> = {}
  for (const k of medanForm(meta, baru)) {
    const v = (isi[k.kunci] ?? '').trim()
    if (k.wajib && v === '') g[k.kunci] = TEKS.wajib
    else if (k.lebar > 0 && panjangByte(v) > k.lebar) g[k.kunci] = TEKS.terlaluPanjang(k.lebar)
  }
  return g
}

/** Badan permintaan: medan form saja (turunan tidak dikirim), dipangkas. */
export function badan(meta: MetaMaster, baru: boolean, isi: Record<string, string>): Record<string, string> {
  return Object.fromEntries(medanForm(meta, baru).map((k) => [k.kunci, (isi[k.kunci] ?? '').trim()]))
}

export default function FormMaster({
  meta,
  baris,
  onTutup,
  onTersimpan,
}: {
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
  const rujukan = RUJUKAN[meta.kunci] ?? {}

  async function simpan() {
    setCoba(true)
    if (Object.keys(periksa(meta, baru, isi)).length > 0) return
    setMenyimpan(true)
    setGalat(null)
    try {
      const h = baru
        ? await tambahMaster(meta.kunci, badan(meta, baru, isi))
        : await ubahMaster(meta.kunci, String(baris?.id ?? ''), badan(meta, baru, isi))
      onTersimpan(h.id)
    } catch (err) {
      setGalat(err)
    } finally {
      setMenyimpan(false)
    }
  }

  return (
    <Modal
      judul={baru ? TEKS.judulTambah(meta.judul) : TEKS.judulUbah(meta.judul)}
      onTutup={onTutup}
      onKirim={() => void simpan()}
      lebar
      aksi={
        <button type="submit" className="btn btn--primary" disabled={menyimpan}>
          {menyimpan ? TEKS.menyimpan : TEKS.simpan}
        </button>
      }
    >
      <div className="form-grid">
        {(!baru || meta.idOtomatis) && (
          <Field label={labelKolom('id', 'ID')} value={baru ? TEKS.idOtomatis : String(baris?.id ?? '')} onChange={() => {}} readOnly />
        )}
        {medanForm(meta, baru).map((k) => {
          const label = labelKolom(k.kunci, k.kolom)
          const r = rujukan[k.kunci]
          if (r)
            return (
              <PilihRujukan
                key={k.kunci}
                label={label}
                value={isi[k.kunci] ?? ''}
                onChange={set(k.kunci)}
                master={r.master}
                nilai={r.nilai}
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
              <Field key={k.kunci} label={`${labelKolom(k.kunci, k.kolom)} (${TEKS.turunan})`} value={isi[k.kunci] ?? ''} onChange={() => {}} readOnly />
            ))}
      </div>
      <Gagal galat={galat} />
    </Modal>
  )
}

// Daftar satu master (inti; template Kelola User: toolbar cari + Tambah, tabel, aksi per baris):
// - cari (ID / nama, jeda 400 ms) dan saring status (Semua / Aktif / Nonaktif), 50 baris per halaman;
// - Tambah / Ubah -> `FormMaster`; Aktifkan / Nonaktifkan -> konfirmasi -> `PUT P/{id}/status` (tanpa hapus).

import { useEffect, useState } from 'react'

import { Gagal, Kosong, Memuat, Modal } from '../components/ui/dasar'
import type { BarisMaster, HalamanMaster, KlienMaster, MetaMaster, StatusSaring } from './api'
import FormMaster from './FormMaster'
import { labelKolomMaster as labelKolom, TEKS_MASTER as TEKS } from './labels'

/** Jeda cari. */
export const JEDA_CARI_MS = 400

/** Jumlah halaman dari total dan ukuran (minimal 1). */
export const jumlahHalaman = (total: number, ukuran: number) => Math.max(1, Math.ceil(total / Math.max(1, ukuran)))

/**
 * Nomor halaman yang ditampilkan: selalu halaman 1 dan terakhir, ditambah dua halaman di kiri-kanan halaman aktif;
 * `null` = jeda (…). Contoh: kini 6 dari 20 -> 1 … 4 5 6 7 8 … 20.
 */
export function nomorHalaman(kini: number, dari: number): (number | null)[] {
  const tampil = new Set<number>([1, dari])
  for (let n = kini - 2; n <= kini + 2; n++) if (n >= 1 && n <= dari) tampil.add(n)
  const urut = [...tampil].sort((a, b) => a - b)
  const hasil: (number | null)[] = []
  urut.forEach((n, i) => {
    if (i > 0 && n - urut[i - 1]! > 1) hasil.push(null)
    hasil.push(n)
  })
  return hasil
}

export default function DaftarMaster({ klien, meta }: { klien: KlienMaster; meta: MetaMaster }) {
  const [kueri, setKueri] = useState('')
  const [q, setQ] = useState('')
  const [status, setStatus] = useState<StatusSaring>('')
  const [halaman, setHalaman] = useState(1)
  const [data, setData] = useState<HalamanMaster | null>(null)
  const [galat, setGalat] = useState<unknown>(null)
  const [muat, setMuat] = useState(0)
  /** undefined = form tertutup; null = Tambah; baris = Ubah. */
  const [form, setForm] = useState<BarisMaster | null | undefined>(undefined)
  const [konfirmasi, setKonfirmasi] = useState<BarisMaster | null>(null)
  const [sibuk, setSibuk] = useState(false)
  const [galatAksi, setGalatAksi] = useState<unknown>(null)
  const [pesan, setPesan] = useState<string | null>(null)

  // Kotak cari -> kueri sesudah jeda; halaman kembali ke 1.
  useEffect(() => {
    const j = window.setTimeout(() => {
      setQ(kueri)
      setHalaman(1)
    }, JEDA_CARI_MS)
    return () => window.clearTimeout(j)
  }, [kueri])

  useEffect(() => {
    let batal = false
    setData(null)
    setGalat(null)
    klien.cari(q, status, halaman).then(
      (h) => {
        if (!batal) setData(h)
      },
      (err: unknown) => {
        if (!batal) setGalat(err)
      },
    )
    return () => {
      batal = true
    }
  }, [klien, q, status, halaman, muat])

  async function ubahStatus(b: BarisMaster) {
    const id = String(b.id ?? '')
    setSibuk(true)
    setGalatAksi(null)
    try {
      const h = await klien.ubahStatus(id, !b.aktif)
      setPesan(TEKS.statusBerubah(id, h.aktif))
      setKonfirmasi(null)
      setMuat((x) => x + 1)
    } catch (err) {
      setGalatAksi(err)
    } finally {
      setSibuk(false)
    }
  }

  const kolom = meta.kolom
  const dari = data ? jumlahHalaman(data.total, data.ukuran) : 1

  return (
    <>
      <div className="toolbar">
        <input
          className="field__input"
          type="search"
          aria-label={TEKS.cari}
          placeholder={TEKS.cari}
          value={kueri}
          onChange={(e) => setKueri(e.target.value)}
        />
        <select
          className="field__input"
          aria-label={TEKS.kolomStatus}
          value={status}
          onChange={(e) => {
            setStatus(e.target.value as StatusSaring)
            setHalaman(1)
          }}
        >
          <option value="">{TEKS.semuaStatus}</option>
          <option value="aktif">{TEKS.aktif}</option>
          <option value="nonaktif">{TEKS.nonaktif}</option>
        </select>
        <span className="toolbar__spacer" />
        <button
          type="button"
          className="btn btn--primary"
          onClick={() => {
            setPesan(null)
            setForm(null)
          }}
        >
          {TEKS.tambah}
        </button>
      </div>

      {pesan !== null && (
        <div className="alert alert--ok" role="status">
          {pesan}
        </div>
      )}
      {galat !== null && <Gagal galat={galat} />}
      {data === null && galat === null && <Memuat pesan={TEKS.memuat} />}
      {data !== null && data.baris.length === 0 && <Kosong pesan={q === '' && status === '' ? TEKS.kosong : TEKS.tidakCocok} />}
      {data !== null && data.baris.length > 0 && (
        <>
          <div className="table-wrap">
            <table className="inbox__tabel">
              <thead>
                <tr>
                  {kolom.map((k) => (
                    <th key={k.kunci} scope="col">
                      {labelKolom(k.kunci, k.kolom)}
                    </th>
                  ))}
                  <th scope="col">{TEKS.kolomStatus}</th>
                  <th scope="col" className="table__actions">
                    {TEKS.kolomAksi}
                  </th>
                </tr>
              </thead>
              <tbody>
                {data.baris.map((b) => (
                  <tr key={String(b.id)} className="inbox__baris">
                    {kolom.map((k) => (
                      <td key={k.kunci}>{String(b[k.kunci] ?? '')}</td>
                    ))}
                    <td>
                      <span className="badge">{b.aktif ? TEKS.aktif : TEKS.nonaktif}</span>
                    </td>
                    <td className="table__actions">
                      <button
                        type="button"
                        className="btn btn--ghost btn--sm"
                        onClick={() => {
                          setPesan(null)
                          setForm(b)
                        }}
                      >
                        {TEKS.ubah}
                      </button>
                      <button
                        type="button"
                        className="btn btn--ghost btn--sm"
                        onClick={() => {
                          setPesan(null)
                          setGalatAksi(null)
                          setKonfirmasi(b)
                        }}
                      >
                        {b.aktif ? TEKS.nonaktifkan : TEKS.aktifkan}
                      </button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
          <div className="pager">
            <span className="muted">{TEKS.menampilkan(data.halaman, data.ukuran, data.total)}</span>
            <nav className="pager__aksi" aria-label={TEKS.navigasiHalaman}>
              <button type="button" className="btn btn--ghost btn--sm" disabled={halaman <= 1} onClick={() => setHalaman(1)}>
                {TEKS.pertama}
              </button>
              <button type="button" className="btn btn--ghost btn--sm" disabled={halaman <= 1} onClick={() => setHalaman((h) => h - 1)}>
                {TEKS.sebelumnya}
              </button>
              {nomorHalaman(data.halaman, dari).map((n, i) =>
                n === null ? (
                  <span key={`jeda-${i}`} className="muted">
                    …
                  </span>
                ) : (
                  <button
                    key={n}
                    type="button"
                    className={'btn btn--sm ' + (n === data.halaman ? 'btn--primary' : 'btn--ghost')}
                    aria-current={n === data.halaman ? 'page' : undefined}
                    onClick={() => setHalaman(n)}
                  >
                    {n}
                  </button>
                ),
              )}
              <button type="button" className="btn btn--ghost btn--sm" disabled={halaman >= dari} onClick={() => setHalaman((h) => h + 1)}>
                {TEKS.berikutnya}
              </button>
              <button type="button" className="btn btn--ghost btn--sm" disabled={halaman >= dari} onClick={() => setHalaman(dari)}>
                {TEKS.terakhir}
              </button>
            </nav>
          </div>
        </>
      )}

      {form !== undefined && (
        <FormMaster
          klien={klien}
          meta={meta}
          baris={form}
          onTutup={() => setForm(undefined)}
          onTersimpan={(id) => {
            setForm(undefined)
            setPesan(TEKS.tersimpan(id))
            setMuat((x) => x + 1)
          }}
        />
      )}
      {konfirmasi !== null && (
        <Modal
          judul={TEKS.judulKonfirmasi(!konfirmasi.aktif)}
          onTutup={() => setKonfirmasi(null)}
          aksi={
            <button type="button" className="btn btn--primary" disabled={sibuk} onClick={() => void ubahStatus(konfirmasi)}>
              {konfirmasi.aktif ? TEKS.nonaktifkan : TEKS.aktifkan}
            </button>
          }
        >
          <p>{TEKS.konfirmasiStatus(String(konfirmasi.id ?? ''), !konfirmasi.aktif)}</p>
          {galatAksi !== null && <Gagal galat={galatAksi} />}
        </Modal>
      )}
    </>
  )
}

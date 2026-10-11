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

/** Kunci jejak ubah (MD-7) - digabung dua kolom "Dibuat" / "Diubah" (nama, tanggal di bawahnya), pola Kelola User. */
export const KUNCI_JEJAK = ['createOp', 'tglCreate', 'updateOp', 'tglUpdate'] as const

/**
 * Susunan kolom tampil (work owner 10-10-2026 "perbaiki biar user friendly, template nya ikutin menu kelola user"):
 * - Old ID di bawah ID (pola Contact ID di bawah username), bukan kolom sendiri;
 * - empat kolom jejak menjadi "Dibuat" / "Diubah" bila meta memuatnya;
 * - kolom data lain tetap urutan meta.
 */
export function susunKolom(kunci: readonly string[]) {
  const ada = new Set(kunci)
  const oldDiBawahId = ada.has('id') && ada.has('oldId')
  const data = kunci.filter((k) => !(KUNCI_JEJAK as readonly string[]).includes(k) && !(oldDiBawahId && k === 'oldId'))
  return {
    data,
    oldDiBawahId,
    dibuat: ada.has('createOp') || ada.has('tglCreate'),
    diubah: ada.has('updateOp') || ada.has('tglUpdate'),
  }
}

/** Sel dua baris: nilai utama, nilai kedua redup di bawahnya; keduanya kosong = "—". */
function SelDua({ atas, bawah }: { atas: string; bawah: string }) {
  if (atas === '' && bawah === '') return <span className="muted">—</span>
  return (
    <>
      {atas === '' ? <span className="muted">—</span> : atas}
      {bawah !== '' && <div className="muted">{bawah}</div>}
    </>
  )
}

/** Tabel satu halaman master (template Kelola User) - murni, tanpa muat data. */
export function TabelMaster({
  kolom,
  baris,
  onUbah,
  onStatus,
}: {
  kolom: MetaMaster['kolom']
  baris: BarisMaster[]
  onUbah: (b: BarisMaster) => void
  onStatus: (b: BarisMaster) => void
}) {
  const susun = susunKolom(kolom.map((k) => k.kunci))
  const kolomData = kolom.filter((k) => susun.data.includes(k.kunci))
  const teks = (b: BarisMaster, k: string) => String(b[k] ?? '')
  return (
    <table className="inbox__tabel kelola-user__tabel">
      <thead>
        <tr>
          {kolomData.map((k) => (
            <th key={k.kunci} scope="col">
              {labelKolom(k.kunci, k.kolom)}
            </th>
          ))}
          {susun.dibuat && <th scope="col">{TEKS.kolomDibuat}</th>}
          {susun.diubah && <th scope="col">{TEKS.kolomDiubah}</th>}
          <th scope="col">{TEKS.kolomStatus}</th>
          <th scope="col" className="table__actions">
            {TEKS.kolomAksi}
          </th>
        </tr>
      </thead>
      <tbody>
        {baris.map((b) => (
          <tr key={String(b.id)} className="inbox__baris">
            {kolomData.map((k) => (
              <td key={k.kunci}>
                {teks(b, k.kunci)}
                {k.kunci === 'id' && susun.oldDiBawahId && teks(b, 'oldId') !== '' && (
                  <div className="muted">
                    {labelKolom('oldId', 'oldId')} {teks(b, 'oldId')}
                  </div>
                )}
              </td>
            ))}
            {susun.dibuat && (
              <td>
                <SelDua atas={teks(b, 'createOp')} bawah={teks(b, 'tglCreate')} />
              </td>
            )}
            {susun.diubah && (
              <td>
                <SelDua atas={teks(b, 'updateOp')} bawah={teks(b, 'tglUpdate')} />
              </td>
            )}
            <td>
              <span className="kelola-user__status">
                <span className={`badge ${b.aktif ? 'kelola-user__badge--aktif' : 'kelola-user__badge--nonaktif'}`}>
                  {b.aktif ? TEKS.aktif : TEKS.nonaktif}
                </span>
              </span>
            </td>
            <td className="table__actions">
              <span className="kelola-user__aksi">
                <button
                  type="button"
                  className="btn btn--ghost btn--sm"
                  onClick={() => onUbah(b)}
                >
                  {TEKS.ubah}
                </button>
                <button
                  type="button"
                  className="btn btn--ghost btn--sm"
                  onClick={() => onStatus(b)}
                >
                  {b.aktif ? TEKS.nonaktifkan : TEKS.aktifkan}
                </button>
              </span>
            </td>
          </tr>
        ))}
      </tbody>
    </table>
  )
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
          className="field__input kelola-user__cari"
          type="search"
          aria-label={TEKS.cari}
          placeholder={TEKS.cari}
          value={kueri}
          onChange={(e) => setKueri(e.target.value)}
        />
        <select
          className="field__input master-data__saring"
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
            <TabelMaster
              kolom={kolom}
              baris={data.baris}
              onUbah={(b) => {
                setPesan(null)
                setForm(b)
              }}
              onStatus={(b) => {
                setPesan(null)
                setGalatAksi(null)
                setKonfirmasi(b)
              }}
            />
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

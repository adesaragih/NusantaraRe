// Lampiran klaim induk di layar komite - perintah work owner 11-10-2026 ("TAMBAHKAN ATTACHMENT ... KLAIM FACIN KE KOMITE
// KLAIM FACIN. YANG MUNCUL HANYA KATEGORI YANG ADA ISI NYA"; "TOMBOL ATTACHNYA DI KOMITE DI DISABLE AJA"). Susunan sama
// dengan tab lampiran Claim Fac In (grid Category - Count Attach - Upload File - View File, 5 baris lalu "Show All"),
// salinan hanya-baca - modul tidak saling impor: server hanya mengirim kategori yang berisi; Add attachment dan Upload
// File tampil NONAKTIF; jendela View File tanpa Delete / Change Category. View pdf / gambar dan View Office Online
// membuka popup penampil (pola NB Treaty In, `inti/frontend/lampiran/penampil.ts`).

import { useCallback, useEffect, useRef, useState } from 'react'

import { ApiFailure } from '../../../../inti/frontend/klien'
import { Gagal, Memuat, Modal } from '../../../../inti/frontend/components/ui/dasar'
import {
  BINGKAI_PENAMPIL,
  PARAM_PENAMPIL,
  PENAMPIL_OFFICE,
  bisaViewOffice,
  jenisViewOnline,
  mimeViewOnline,
} from '../../../../inti/frontend/lampiran/penampil'
import {
  ambilIsiLampiran,
  ambilLampiran,
  tautanOfficeLampiran,
  unduhLampiran,
  type KategoriLampiran,
  type Lampiran,
  type LampiranKomite,
} from '../api'
import { LAMPIRAN_KCFI } from '../labels'

/** Baris grid sebelum "Show All" (sama dengan tab lampiran Claim Fac In). */
const BARIS_AWAL = 5

function IkonUnggah() {
  return (
    <svg
      width="20"
      height="20"
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      strokeWidth="2"
      aria-hidden="true"
    >
      <path d="M12 16V4M7 9l5-5 5 5M4 15v4h16v-4" />
    </svg>
  )
}

function IkonLihat() {
  return (
    <svg
      width="20"
      height="20"
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      strokeWidth="2"
      aria-hidden="true"
    >
      <path d="M12 4v12M7 11l5 5 5-5M4 15v4h16v-4" />
    </svg>
  )
}

export default function PanelLampiran({ id }: { id: string }) {
  const [data, setData] = useState<LampiranKomite | null>(null)
  const [galat, setGalat] = useState<unknown>(null)
  const [semua, setSemua] = useState(false)
  const [lihat, setLihat] = useState<KategoriLampiran | null>(null)

  const muat = useCallback(() => {
    ambilLampiran(id).then(
      (d) => {
        setData(d)
        setGalat(null)
      },
      (g: unknown) => setGalat(g),
    )
  }, [id])

  useEffect(() => {
    muat()
  }, [muat])

  const daftarKategori = data?.kategori ?? []
  const tampil = semua ? daftarKategori : daftarKategori.slice(0, BARIS_AWAL)
  const pesanGalat = galat instanceof ApiFailure ? (galat.detail.message ?? null) : null
  return (
    <section className="panel komiteclaimfacin__kartu komiteclaimfacin__lampiran" aria-label={LAMPIRAN_KCFI.judul}>
      <h3 className="panel__title komiteclaimfacin__kartu-judul">{LAMPIRAN_KCFI.judul}</h3>
      {galat !== null &&
        (pesanGalat ? <div className="alert alert--error">{pesanGalat}</div> : <Gagal galat={galat} />)}
      <div className="komiteclaimfacin__lampiran-aksi">
        <button type="button" className="btn btn--ghost btn--sm" disabled title={LAMPIRAN_KCFI.hanyaDiKlaim}>
          {LAMPIRAN_KCFI.tambah}
        </button>
        <button type="button" className="btn btn--ghost btn--sm" onClick={muat}>
          {LAMPIRAN_KCFI.muatUlang}
        </button>
      </div>
      {data === null && galat === null ? (
        <Memuat />
      ) : daftarKategori.length === 0 ? (
        <p className="komiteclaimfacin__kosong">{LAMPIRAN_KCFI.kosong}</p>
      ) : (
        <>
          <div className="komiteclaimfacin__grid-alat">
            <span>
              {LAMPIRAN_KCFI.menampilkan} {tampil.length} {LAMPIRAN_KCFI.dari} {daftarKategori.length}
            </span>
            {!semua && daftarKategori.length > BARIS_AWAL && (
              <button type="button" className="btn btn--ghost btn--sm" onClick={() => setSemua(true)}>
                {LAMPIRAN_KCFI.semua}
              </button>
            )}
          </div>
          <GridKategori kategori={tampil} onLihat={setLihat} />
        </>
      )}

      {lihat && data && (
        <DialogLihat
          id={id}
          awal={lihat.id}
          kategori={daftarKategori}
          lampiran={data.lampiran}
          onTutup={() => setLihat(null)}
        />
      )}
    </section>
  )
}

/** Grid kategori (hanya yang berisi): Upload File selalu NONAKTIF di komite, View File membuka daftar berkas. */
export function GridKategori({
  kategori,
  onLihat,
}: {
  kategori: KategoriLampiran[]
  onLihat: (k: KategoriLampiran) => void
}) {
  return (
    <div className="komiteclaimfacin__gulir">
      <table className="komiteclaimfacin__tabel">
        <thead>
          <tr>
            <th>{LAMPIRAN_KCFI.kategori}</th>
            <th>{LAMPIRAN_KCFI.cacah}</th>
            <th className="komiteclaimfacin__kolom-ikon">{LAMPIRAN_KCFI.unggah}</th>
            <th className="komiteclaimfacin__kolom-ikon">{LAMPIRAN_KCFI.lihat}</th>
          </tr>
        </thead>
        <tbody>
          {kategori.map((k) => (
            <tr key={k.id}>
              <td>{k.label}</td>
              <td>{k.countAttach}</td>
              <td className="komiteclaimfacin__kolom-ikon">
                <button
                  type="button"
                  className="btn btn--sm komiteclaimfacin__ikon-lampiran komiteclaimfacin__ikon-lampiran--unggah"
                  aria-label={`${LAMPIRAN_KCFI.unggah} ${k.label}`}
                  title={LAMPIRAN_KCFI.hanyaDiKlaim}
                  disabled
                >
                  <IkonUnggah />
                </button>
              </td>
              <td className="komiteclaimfacin__kolom-ikon">
                <button
                  type="button"
                  className="btn btn--sm komiteclaimfacin__ikon-lampiran komiteclaimfacin__ikon-lampiran--lihat"
                  aria-label={`${LAMPIRAN_KCFI.lihat} ${k.label}`}
                  onClick={() => onLihat(k)}
                >
                  <IkonLihat />
                </button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}

/**
 * View File hanya-baca: pilih Category ("Label - (jumlah)", hanya kategori yang berisi), tabel File Name / Category /
 * Create Date / Attached By / View / centang Action, tombol Download Selected. View dan View Office Online membuka
 * POPUP PENAMPIL layar penuh yang menggantikan daftar.
 */
function DialogLihat({
  id,
  awal,
  kategori,
  lampiran,
  onTutup,
}: {
  id: string
  awal: string
  kategori: KategoriLampiran[]
  lampiran: Lampiran[]
  onTutup: () => void
}) {
  const [kat, setKat] = useState(awal)
  const [pilih, setPilih] = useState<ReadonlySet<string>>(new Set())
  const [galat, setGalat] = useState<unknown>(null)
  const [sibuk, setSibuk] = useState(false)
  const [penampil, setPenampil] = useState<{ nama: string; jenis: 'office' | 'pdf' | 'gambar'; siap: boolean } | null>(
    null,
  )
  const formOffice = useRef<HTMLFormElement>(null)
  const urlOffice = useRef<HTMLInputElement>(null)
  const bingkai = useRef<HTMLIFrameElement>(null)
  const gambar = useRef<HTMLImageElement>(null)
  // Objek URL yang sedang tampil (dicabut saat penampil ditutup) dan giliran buka (jawaban penampil tertutup dibuang).
  const objekAktif = useRef<string | null>(null)
  const giliran = useRef(0)

  useEffect(
    () => () => {
      if (objekAktif.current !== null) URL.revokeObjectURL(objekAktif.current)
    },
    [],
  )

  const daftar = lampiran.filter((a) => a.kategori === kat)
  const terpilih = daftar.filter((a) => pilih.has(a.id))
  const labelKat = (k: string) => kategori.find((x) => x.id === k)?.label ?? k
  const semuaTerpilih = daftar.length > 0 && terpilih.length === daftar.length

  const centang = (lid: string, ya: boolean) =>
    setPilih((s) => {
      const n = new Set(s)
      if (ya) n.add(lid)
      else n.delete(lid)
      return n
    })

  const jalankan = (kerja: () => Promise<void>) => {
    if (sibuk) return
    setSibuk(true)
    setGalat(null)
    kerja()
      .catch((g: unknown) => setGalat(g))
      .finally(() => setSibuk(false))
  }

  const tutupPenampil = () => {
    giliran.current++
    if (objekAktif.current !== null) URL.revokeObjectURL(objekAktif.current)
    objekAktif.current = null
    setPenampil(null)
  }

  /** View Office Online: URL bertanda tangan -> form GET ke bingkai popup penampil. */
  const bukaOffice = (a: Lampiran) => {
    if (sibuk) return
    const ke = ++giliran.current
    setPenampil({ nama: a.namaFile, jenis: 'office', siap: false })
    jalankan(async () => {
      try {
        const { url } = await tautanOfficeLampiran(id, a)
        if (ke !== giliran.current) return
        if (formOffice.current !== null && urlOffice.current !== null) {
          urlOffice.current.value = url
          formOffice.current.submit()
        }
        setPenampil((p) => (p === null ? p : { ...p, siap: true }))
      } catch (g: unknown) {
        if (ke === giliran.current) tutupPenampil()
        throw g
      }
    })
  }

  /** View pdf / gambar: isi dari rute unduh beridentitas -> objek URL lokal di popup penampil. */
  const bukaOnline = (a: Lampiran) => {
    const jenis = jenisViewOnline(a.mime)
    if (sibuk || jenis === null) return
    const ke = ++giliran.current
    setPenampil({ nama: a.namaFile, jenis, siap: false })
    jalankan(async () => {
      try {
        const isi = await ambilIsiLampiran(id, a)
        if (ke !== giliran.current) return
        const objekURL = URL.createObjectURL(new Blob([isi], { type: mimeViewOnline(a.mime) }))
        objekAktif.current = objekURL
        if (jenis === 'gambar' && gambar.current !== null) gambar.current.src = objekURL
        if (jenis === 'pdf' && bingkai.current !== null) bingkai.current.src = objekURL
        setPenampil((p) => (p === null ? p : { ...p, siap: true }))
      } catch (g: unknown) {
        if (ke === giliran.current) tutupPenampil()
        throw g
      }
    })
  }

  // ⛔ `key` BERBEDA untuk popup penampil dan popup daftar (bug layar beku lampiran Bordereaux 08-10-2026).
  if (penampil !== null) {
    return (
      <Modal key="penampil" judul={penampil.nama} onTutup={tutupPenampil} labelBatal={LAMPIRAN_KCFI.tutup} penuh>
        {!penampil.siap && <Memuat />}
        <form ref={formOffice} method="get" action={PENAMPIL_OFFICE} target={BINGKAI_PENAMPIL} hidden>
          <input ref={urlOffice} type="hidden" name={PARAM_PENAMPIL} />
        </form>
        {penampil.jenis === 'gambar' ? (
          <img ref={gambar} alt={penampil.nama} className="komiteclaimfacin__penampil-gambar" />
        ) : (
          <iframe
            ref={bingkai}
            name={BINGKAI_PENAMPIL}
            title={penampil.nama}
            className="komiteclaimfacin__penampil-bingkai"
          />
        )}
      </Modal>
    )
  }

  const pesanGalat = galat instanceof ApiFailure ? (galat.detail.message ?? null) : null
  return (
    <Modal key="daftar" judul={LAMPIRAN_KCFI.lihat} onTutup={onTutup} labelBatal={LAMPIRAN_KCFI.tutup} lebar>
      <label>
        {LAMPIRAN_KCFI.kategori}{' '}
        <select
          className="field__input"
          value={kat}
          disabled={sibuk}
          onChange={(e) => {
            setKat(e.target.value)
            setPilih(new Set())
          }}
        >
          {kategori.map((k) => (
            <option key={k.id} value={k.id}>
              {k.label} - ({k.countAttach})
            </option>
          ))}
        </select>
      </label>
      {galat !== null &&
        (pesanGalat ? <div className="alert alert--error">{pesanGalat}</div> : <Gagal galat={galat} />)}
      {daftar.length === 0 ? (
        <p className="komiteclaimfacin__kosong">{LAMPIRAN_KCFI.kosong}</p>
      ) : (
        <div className="komiteclaimfacin__gulir">
          <table className="komiteclaimfacin__tabel">
            <thead>
              <tr>
                <th>{LAMPIRAN_KCFI.nama}</th>
                <th>{LAMPIRAN_KCFI.kategori}</th>
                <th>{LAMPIRAN_KCFI.tanggal}</th>
                <th>{LAMPIRAN_KCFI.oleh}</th>
                <th aria-label={LAMPIRAN_KCFI.view} />
                <th>
                  <label>
                    <input
                      type="checkbox"
                      checked={semuaTerpilih}
                      disabled={sibuk}
                      onChange={(e) => setPilih(e.target.checked ? new Set(daftar.map((a) => a.id)) : new Set())}
                    />{' '}
                    {LAMPIRAN_KCFI.aksi}
                  </label>
                </th>
              </tr>
            </thead>
            <tbody>
              {daftar.map((a) => (
                <tr key={a.id}>
                  <td>
                    <button
                      type="button"
                      className="btn btn--ghost btn--sm"
                      disabled={sibuk}
                      onClick={() => jalankan(() => unduhLampiran(id, a))}
                    >
                      {a.namaFile}
                    </button>
                  </td>
                  <td>{labelKat(a.kategori)}</td>
                  <td>{a.tanggal}</td>
                  <td>{a.operator}</td>
                  <td>
                    <div className="komiteclaimfacin__lampiran-aksi">
                      {a.adaObjek && jenisViewOnline(a.mime) !== null && (
                        <button
                          type="button"
                          className="btn btn--ghost btn--sm"
                          disabled={sibuk}
                          onClick={() => bukaOnline(a)}
                        >
                          {LAMPIRAN_KCFI.view}
                        </button>
                      )}
                      {a.adaObjek && bisaViewOffice(a.mime) && (
                        <button
                          type="button"
                          className="btn btn--ghost btn--sm"
                          disabled={sibuk}
                          onClick={() => bukaOffice(a)}
                        >
                          {LAMPIRAN_KCFI.viewOffice}
                        </button>
                      )}
                    </div>
                  </td>
                  <td>
                    <input
                      type="checkbox"
                      aria-label={`${LAMPIRAN_KCFI.aksi} ${a.namaFile}`}
                      checked={pilih.has(a.id)}
                      disabled={sibuk}
                      onChange={(e) => centang(a.id, e.target.checked)}
                    />
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
      <div className="komiteclaimfacin__grid-alat">
        <button
          type="button"
          className="btn btn--ghost btn--sm"
          disabled={sibuk || terpilih.length === 0}
          onClick={() =>
            jalankan(async () => {
              for (const a of terpilih) await unduhLampiran(id, a)
            })
          }
        >
          {LAMPIRAN_KCFI.unduhPilih}
        </button>
      </div>
    </Modal>
  )
}

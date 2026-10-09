// Lampiran "Reas" kasus - grid `AttachmentGridReas` korpus NB FacIn (dipinjam NB / EDM Treaty In lewat
// `SetCategoryAttach`), di atas `inti/backend/dokumenpolis` (tabel `DOCUMENT_POLIS`) dan penyimpanan bersama
// `inti/backend/penyimpanan`. Panel di bawah layar kasus (keputusan work owner 08-10-2026).
//
//  grid kategori `CATEGORY_ATTACH_REAS`: Category b1340, Count b1493, Upload File b1643, View File b1791 - berbentuk
//  contoh work owner 08-10-2026: pager inti 5 baris di atas tabel, Count lencana, Upload / View File tombol ikon
//  (nama kolomnya jadi `title` / `aria-label`). Tanpa Refresh (tidak ada di XML): grid dimuat ulang sesudah Upload /
//  Delete.
//  `Upload File` = modal `AttachContentGIS` (Attach b24 / Cancel b22): seret-lepas atau pilih, banyak berkas diunggah
//  satu per satu (pola lampiran Bordereaux).
//  `View File` = popup `ReasViewAttachment`: File b1363 (unduh), Note b1618, Upload Date b1762, View Office Online b2670
//  (xls…pptx dan objek tercatat, b2951), Delete b3434; `View` pdf / gambar (seperti Product Name Life). View dan View
//  Office Online membuka popup penampil layar penuh.
//  Upload File dan Delete hanya bila `bolehUbah` = kasus belum Resolve (keputusan work owner 08-10-2026: "semua bisa
//  asal belum resolve") - diputuskan backend, panel hanya menyembunyikan tombolnya.
//
// ⛔ Unduh dan View lewat fetch beridentitas (penjaga `unduhdokumen.test.ts`); View Office Online = form GET
// tersembunyi ke penampil di bingkai popup penampil (`penampil.ts`).

import { useCallback, useEffect, useRef, useState } from 'react'

import { BahasaUI } from '../components/ui/bahasaUI'
import { Gagal, Halaman, Memuat, Modal } from '../components/ui/dasar'
import { pesanGalat } from '../klien'
import {
  ambilDokumen,
  ambilGrid,
  ambilIsiDokumen,
  hapusDokumen,
  tautanOffice,
  unduhDokumen,
  unggahDokumen,
  type DokumenReas,
  type KategoriReas,
} from './api'
import { LAMPIRAN_REAS } from './labels'
import {
  BINGKAI_PENAMPIL,
  PARAM_PENAMPIL,
  PENAMPIL_OFFICE,
  bisaViewOffice,
  gabungBerkas,
  jenisViewOnline,
  mimeViewOnline,
} from './penampil'

/** Grid kategori: 5 baris per halaman, pager inti `Halaman` di atas tabel (permintaan work owner 08-10-2026). */
const UKURAN_GRID = 5

/** Upload File: panah naik di atas garis dasar (`currentColor`, pola ikon `inti/frontend/components/ui/dasar`). */
function IkonUnggah() {
  return (
    <svg
      width={16}
      height={16}
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      strokeWidth={2.2}
      strokeLinecap="round"
      strokeLinejoin="round"
      aria-hidden
    >
      <path d="M12 16V5M7 10l5-5 5 5M5 20h14" />
    </svg>
  )
}

/** View File: panah turun di atas garis dasar. */
function IkonUnduh() {
  return (
    <svg
      width={16}
      height={16}
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      strokeWidth={2.2}
      strokeLinecap="round"
      strokeLinejoin="round"
      aria-hidden
    >
      <path d="M12 4v11M7 10l5 5 5-5M5 20h14" />
    </svg>
  )
}

/** `dasar` = jalur rute lampiran kasus, mis. `${PREFIX}/kasus/${encodeURIComponent(id)}/lampiran`. */
export default function PanelLampiranReas({ dasar }: { dasar: string }) {
  const [grid, setGrid] = useState<{ daftar: KategoriReas[]; bolehUbah: boolean } | null>(null)
  const [galat, setGalat] = useState<unknown>(null)
  const [hal, setHal] = useState(1)
  const [unggahKe, setUnggahKe] = useState<KategoriReas | null>(null)
  const [lihatKe, setLihatKe] = useState<KategoriReas | null>(null)

  const muat = useCallback(() => {
    setGalat(null)
    ambilGrid(dasar).then(setGrid, (g: unknown) => setGalat(g))
  }, [dasar])

  useEffect(() => {
    muat()
  }, [muat])

  const bolehUbah = grid?.bolehUbah ?? false

  // Teks bawaan komponen inti (pager `Halaman`, modal, galat) berbahasa Inggris seperti label panel ini (permintaan
  // work owner 08-10-2026 "ubah pake bahasa inggris").
  return (
    <BahasaUI.Provider value="en">
      <section className="panel lampiran-reas">
        <h3 className="lampiran-reas__judul">{LAMPIRAN_REAS.judul}</h3>
        {galat !== null && <Gagal galat={galat} />}
        {grid === null && galat === null && <Memuat pesan={LAMPIRAN_REAS.memuat} />}
        {grid !== null && (
          <GridKategori
            daftar={grid.daftar}
            bolehUbah={bolehUbah}
            halaman={hal}
            onHalaman={setHal}
            onUnggah={setUnggahKe}
            onLihat={setLihatKe}
          />
        )}

        {unggahKe !== null && (
          <DialogUnggah
            dasar={dasar}
            kategori={unggahKe}
            onTutup={(berubah) => {
              setUnggahKe(null)
              if (berubah) muat()
            }}
          />
        )}
        {lihatKe !== null && (
          <DialogLihat
            dasar={dasar}
            kategori={lihatKe}
            bolehUbah={bolehUbah}
            onTutup={(berubah) => {
              setLihatKe(null)
              if (berubah) muat()
            }}
          />
        )}
      </section>
    </BahasaUI.Provider>
  )
}

/**
 * Grid kategori `AttachmentGridReas` - satu halaman `UKURAN_GRID` baris. Upload File hanya bila `bolehUbah`; View File
 * selalu. Halaman dijepit seperti `Halaman` menjepitnya: daftar yang memendek tidak meninggalkan halaman kosong.
 */
export function GridKategori({
  daftar,
  bolehUbah,
  halaman,
  onHalaman,
  onUnggah,
  onLihat,
}: {
  daftar: readonly KategoriReas[]
  bolehUbah: boolean
  halaman: number
  onHalaman: (h: number) => void
  onUnggah: (k: KategoriReas) => void
  onLihat: (k: KategoriReas) => void
}) {
  if (daftar.length === 0) return <p className="muted">{LAMPIRAN_REAS.kosong}</p>
  const kini = Math.min(halaman, Math.max(1, Math.ceil(daftar.length / UKURAN_GRID)))
  const tampak = daftar.slice((kini - 1) * UKURAN_GRID, kini * UKURAN_GRID)
  return (
    <>
      <Halaman halaman={kini} ukuran={UKURAN_GRID} total={daftar.length} onPindah={onHalaman} />
      <div className="table-wrap lampiran-reas__grid">
        <table>
          <thead>
            <tr>
              <th scope="col">{LAMPIRAN_REAS.category}</th>
              <th scope="col" className="lampiran-reas__tengah">
                {LAMPIRAN_REAS.count}
              </th>
              <th scope="col" className="lampiran-reas__tengah">
                {LAMPIRAN_REAS.uploadFile}
              </th>
              <th scope="col" className="lampiran-reas__tengah">
                {LAMPIRAN_REAS.viewFile}
              </th>
            </tr>
          </thead>
          <tbody>
            {tampak.map((k) => (
              <tr key={k.nama}>
                <td>{k.nama}</td>
                <td className="lampiran-reas__tengah">
                  <span className={'lampiran-reas__cacah' + (k.cacah > 0 ? ' lampiran-reas__cacah--ada' : '')}>
                    {k.cacah}
                  </span>
                </td>
                <td className="lampiran-reas__tengah">
                  {bolehUbah && (
                    <button
                      type="button"
                      className="lampiran-reas__ikon lampiran-reas__ikon--unggah"
                      title={LAMPIRAN_REAS.uploadFile}
                      aria-label={`${LAMPIRAN_REAS.uploadFile} ${k.nama}`}
                      onClick={() => onUnggah(k)}
                    >
                      <IkonUnggah />
                    </button>
                  )}
                </td>
                <td className="lampiran-reas__tengah">
                  <button
                    type="button"
                    className="lampiran-reas__ikon lampiran-reas__ikon--lihat"
                    title={LAMPIRAN_REAS.viewFile}
                    aria-label={`${LAMPIRAN_REAS.viewFile} ${k.nama}`}
                    onClick={() => onLihat(k)}
                  >
                    <IkonUnduh />
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </>
  )
}

/**
 * `AttachContentGIS` - seret-lepas atau klik untuk memilih, berkas digabung ke daftar; Attach mengunggah satu per satu,
 * yang gagal tetap terpilih.
 */
function DialogUnggah({
  dasar,
  kategori,
  onTutup,
}: {
  dasar: string
  kategori: KategoriReas
  onTutup: (berubah: boolean) => void
}) {
  const [berkas, setBerkas] = useState<File[]>([])
  const [proses, setProses] = useState<string | null>(null)
  const [gagal, setGagal] = useState<{ nama: string; galat: unknown }[]>([])
  const [pesan, setPesan] = useState<string | null>(null)
  const [masuk, setMasuk] = useState(false)
  const [seret, setSeret] = useState(false)

  const tambah = (baru: readonly File[]) => {
    if (proses !== null || baru.length === 0) return
    setBerkas((b) => gabungBerkas(b, baru))
    setGagal([])
    setPesan(null)
  }

  const kirim = async () => {
    if (proses !== null) return
    if (berkas.length === 0) {
      setPesan(LAMPIRAN_REAS.tanpaBerkas)
      return
    }
    setPesan(null)
    const tertolak: { berkas: File; galat: unknown }[] = []
    let ada = masuk
    for (const [i, f] of berkas.entries()) {
      setProses(LAMPIRAN_REAS.mengunggah(i + 1, berkas.length, f.name))
      try {
        await unggahDokumen(dasar, kategori.nama, f)
        ada = true
      } catch (g: unknown) {
        tertolak.push({ berkas: f, galat: g })
      }
    }
    setProses(null)
    setMasuk(ada)
    if (tertolak.length === 0) {
      onTutup(true)
      return
    }
    setBerkas(tertolak.map((t) => t.berkas))
    setGagal(tertolak.map((t) => ({ nama: t.berkas.name, galat: t.galat })))
  }

  return (
    <Modal
      judul={LAMPIRAN_REAS.judulUnggah(kategori.nama)}
      onTutup={() => onTutup(masuk)}
      onKirim={() => void kirim()}
      labelBatal={LAMPIRAN_REAS.cancel}
      aksi={
        <button type="submit" className="btn btn--primary" disabled={proses !== null}>
          {LAMPIRAN_REAS.attach}
        </button>
      }
    >
      <label
        className={'lampiran-reas__seret' + (seret ? ' lampiran-reas__seret--aktif' : '')}
        onDragOver={(e) => {
          e.preventDefault()
          e.dataTransfer.dropEffect = proses === null ? 'copy' : 'none'
          setSeret(true)
        }}
        onDragLeave={(e) => {
          // Pindah ke anak kotak (teks, input) bukan meninggalkan kotak.
          if (!e.currentTarget.contains(e.relatedTarget as Node | null)) setSeret(false)
        }}
        onDrop={(e) => {
          e.preventDefault()
          setSeret(false)
          tambah(Array.from(e.dataTransfer.files))
        }}
      >
        <span>{LAMPIRAN_REAS.seretBerkas}</span>
        <input
          type="file"
          multiple
          aria-label={LAMPIRAN_REAS.pilihBerkas}
          disabled={proses !== null}
          onChange={(e) => {
            const dipilih = Array.from(e.target.files ?? [])
            // Dikosongkan: memilih berkas yang sama lagi tetap memicu onChange.
            e.target.value = ''
            tambah(dipilih)
          }}
        />
      </label>
      {berkas.length > 0 && (
        <ul className="lampiran-reas__pilihan">
          <li className="muted">{LAMPIRAN_REAS.berkasTerpilih(berkas.length)}</li>
          {berkas.map((f) => (
            <li key={f.name + f.size}>
              <span>{f.name}</span>
              <button
                type="button"
                className="btn btn--ghost btn--sm"
                disabled={proses !== null}
                onClick={() => setBerkas((b) => b.filter((x) => x !== f))}
              >
                {LAMPIRAN_REAS.buangPilihan}
              </button>
            </li>
          ))}
        </ul>
      )}
      {proses !== null && <p className="muted">{proses}</p>}
      {pesan !== null && <p className="lampiran-reas__galat">{pesan}</p>}
      {gagal.length > 0 && (
        <div className="lampiran-reas__galat">
          <p>{LAMPIRAN_REAS.gagalUnggah}</p>
          <ul>
            {gagal.map((g) => (
              <li key={g.nama}>
                {g.nama}: {pesanGalat(g.galat) ?? String(g.galat)}
              </li>
            ))}
          </ul>
        </div>
      )}
    </Modal>
  )
}

/**
 * `ReasViewAttachment` - dokumen satu kategori. File mengunduh; `View` (pdf / gambar) dan View Office Online membuka
 * POPUP PENAMPIL layar penuh. Selama penampil terbuka, popup daftar diganti olehnya lalu kembali saat ditutup - modal
 * bertumpuk tertutup bersama oleh Escape.
 */
function DialogLihat({
  dasar,
  kategori,
  bolehUbah,
  onTutup,
}: {
  dasar: string
  kategori: KategoriReas
  bolehUbah: boolean
  onTutup: (berubah: boolean) => void
}) {
  const [daftar, setDaftar] = useState<DokumenReas[] | null>(null)
  const [galat, setGalat] = useState<unknown>(null)
  const [sibuk, setSibuk] = useState(false)
  const [berubah, setBerubah] = useState(false)
  // Popup penampil: berkas, jenis isi, dan apakah isinya sudah terpasang.
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

  const muat = useCallback(() => {
    ambilDokumen(dasar, kategori.nama).then(
      (r) => setDaftar(r.daftar),
      (g: unknown) => setGalat(g),
    )
  }, [dasar, kategori.nama])

  useEffect(() => {
    muat()
  }, [muat])

  useEffect(
    () => () => {
      if (objekAktif.current !== null) URL.revokeObjectURL(objekAktif.current)
    },
    [],
  )

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
  const bukaOffice = (d: DokumenReas) => {
    if (sibuk) return
    const ke = ++giliran.current
    setPenampil({ nama: d.namaFile, jenis: 'office', siap: false })
    jalankan(async () => {
      try {
        const { url } = await tautanOffice(dasar, d)
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
  const bukaOnline = (d: DokumenReas) => {
    const jenis = jenisViewOnline(d.ekstensi)
    if (sibuk || jenis === null) return
    const ke = ++giliran.current
    setPenampil({ nama: d.namaFile, jenis, siap: false })
    jalankan(async () => {
      try {
        const isi = await ambilIsiDokumen(dasar, d)
        if (ke !== giliran.current) return
        // Tipe dari ekstensi, bukan dari jawaban: hanya pdf / gambar raster yang dirender.
        const objekURL = URL.createObjectURL(new Blob([isi], { type: mimeViewOnline(d.ekstensi) }))
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

  // ⛔ `key` BERBEDA untuk popup penampil dan popup daftar: keduanya `<Modal>` di posisi yang sama, sehingga tanpa
  // `key` React memakai ulang instans yang sama - popup daftar mewarisi status "sedang menutup" milik penampil
  // (backdrop keluar menutupi layar, Close / X / Escape mati). Bug lampiran Bordereaux 08-10-2026.
  if (penampil !== null) {
    return (
      <Modal key="penampil" judul={penampil.nama} onTutup={tutupPenampil} labelBatal={LAMPIRAN_REAS.close} penuh>
        {!penampil.siap && <Memuat />}
        {/* `View Office Online`: `<penampil>?src=<URL bertanda tangan>` di bingkai popup ini. */}
        <form ref={formOffice} method="get" action={PENAMPIL_OFFICE} target={BINGKAI_PENAMPIL} hidden>
          <input ref={urlOffice} type="hidden" name={PARAM_PENAMPIL} />
        </form>
        {penampil.jenis === 'gambar' ? (
          <img ref={gambar} alt={penampil.nama} className="lampiran-reas__penampil-gambar" />
        ) : (
          <iframe
            ref={bingkai}
            name={BINGKAI_PENAMPIL}
            title={penampil.nama}
            className="lampiran-reas__penampil-bingkai"
          />
        )}
      </Modal>
    )
  }

  return (
    <Modal
      key="daftar"
      judul={LAMPIRAN_REAS.judulLihat(kategori.nama)}
      onTutup={() => onTutup(berubah)}
      labelBatal={LAMPIRAN_REAS.close}
      lebar
    >
      {galat !== null && <Gagal galat={galat} />}
      {daftar === null && galat === null && <Memuat pesan={LAMPIRAN_REAS.memuat} />}
      {daftar !== null &&
        (daftar.length === 0 ? (
          <p className="muted">{LAMPIRAN_REAS.kosong}</p>
        ) : (
          <div className="table-wrap">
            <table>
              <thead>
                <tr>
                  <th scope="col">{LAMPIRAN_REAS.file}</th>
                  <th scope="col">{LAMPIRAN_REAS.note}</th>
                  <th scope="col">{LAMPIRAN_REAS.uploadDate}</th>
                  <th scope="col" aria-label={LAMPIRAN_REAS.view} />
                  <th scope="col" aria-label={LAMPIRAN_REAS.delete} />
                </tr>
              </thead>
              <tbody>
                {daftar.map((d) => (
                  <tr key={d.id}>
                    <td>
                      <button
                        type="button"
                        className="lampiran-reas__nama"
                        disabled={sibuk}
                        onClick={() => jalankan(() => unduhDokumen(dasar, d))}
                      >
                        {d.namaFile}
                      </button>
                    </td>
                    <td>{d.kategori}</td>
                    <td>{d.tanggal}</td>
                    <td>
                      <div className="lampiran-reas__lihat">
                        {d.adaObjek && jenisViewOnline(d.ekstensi) !== null && (
                          <button
                            type="button"
                            className="btn btn--ghost btn--sm"
                            disabled={sibuk}
                            onClick={() => bukaOnline(d)}
                          >
                            {LAMPIRAN_REAS.view}
                          </button>
                        )}
                        {/* b2951: `.MIME` xls…pptx && `.T_STORAGE_ID != ''` */}
                        {d.adaObjek && bisaViewOffice(d.ekstensi) && (
                          <button
                            type="button"
                            className="btn btn--ghost btn--sm"
                            disabled={sibuk}
                            onClick={() => bukaOffice(d)}
                          >
                            {LAMPIRAN_REAS.viewOffice}
                          </button>
                        )}
                      </div>
                    </td>
                    <td>
                      {bolehUbah && (
                        <button
                          type="button"
                          className="btn btn--danger btn--sm"
                          disabled={sibuk}
                          onClick={() =>
                            jalankan(async () => {
                              await hapusDokumen(dasar, d)
                              setBerubah(true)
                              muat()
                            })
                          }
                        >
                          {LAMPIRAN_REAS.delete}
                        </button>
                      )}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        ))}
    </Modal>
  )
}

import { useState, type FormEvent } from 'react'

import { PanelTotalKlaim } from '../components/PanelTotalKlaim'

import {
  ambilKlaimLife,
  kodeStatusGalat,
  pesanGalat,
  tampilUang,
  tolakBarisAdjustment,
  serahkanKeKomite,
  tambahPutaran,
  simpanAdjustment,
  bolehSimpanAdjustment,
  bolehPutaranBaru,
  bolehSerahkanDiLayar,
  dampakHapusKlaim,
  hapusKlaim,
  STATUS_OUTSTANDING,
  type DampakHapus,
  type Klaim,
} from '../services/api'

// ============================================================================
// pages/KlaimLife.tsx — halaman "buka satu klaim Life" (tiket 01 AC-2).
//
// Cara membaca berkas ini bila baru mengenal React:
//   - Komponen = fungsi yang MENGEMBALIKAN tampilan. Bagian yang mirip HTML di
//     dalam `return (...)` disebut JSX; karena berkas ini TypeScript, namanya TSX.
//   - `useState` = kotak penyimpan nilai. Setiap kali nilainya diganti lewat
//     fungsi set-nya (misalnya `setKlaim`), React menggambar ulang tampilan.
//   - `<string>` di belakang useState memberi tahu TypeScript isi kotak itu apa.
//   - `{ ... }` di dalam JSX = "sisipkan nilai JavaScript di sini".
//
// ⛔ Nilai uang tidak pernah disentuh aritmetika dan tidak pernah dibungkus
// Number(). Ia datang sebagai teks desimal dan ditampilkan apa adanya
// (ADR-U-0003, ADR-U-0016). Penjelasannya di services/api.ts bab 3.
// ============================================================================
export default function KlaimLife() {
  // Empat kotak keadaan halaman ini.
  const [id, setId] = useState<string>('') // yang sedang diketik pengguna
  const [klaim, setKlaim] = useState<Klaim | null>(null) // klaim yang berhasil dibaca
  const [galat, setGalat] = useState<string | null>(null) // pesan bila gagal
  const [memuat, setMemuat] = useState<boolean>(false) // sedang menunggu server?
  const [menolak, setMenolak] = useState<string | null>(null) // baris yang sedang ditolak
  const [menyerahkan, setMenyerahkan] = useState<string | null>(null) // baris ke Komite
  const [memutar, setMemutar] = useState<string | null>(null) // peserta yang dibuka putarannya
  const [mengaksep, setMengaksep] = useState<string | null>(null) // peserta yang diaksep
  const [dampak, setDampak] = useState<DampakHapus | null>(null) // isi popup konfirmasi
  const [menghapus, setMenghapus] = useState<boolean>(false) // permintaan hapus berjalan

  // Membuka popup: MEMBACA dampak, tidak menghapus apa pun.
  async function bukaKonfirmasiHapus() {
    if (!klaim) return
    setGalat(null)
    try {
      setDampak(await dampakHapusKlaim(klaim.id))
    } catch {
      setGalat('Gagal menghitung dampak penghapusan.')
    }
  }

  // ⛔ Batal hanya menutup popup. Tidak ada yang perlu dibatalkan, sebab
  // membuka popup tidak menulis apa pun.
  function batalHapus() {
    setDampak(null)
    setGalat(null)
  }

  async function lanjutkanHapus() {
    if (!klaim) return
    setGalat(null)
    setMenghapus(true)
    try {
      await hapusKlaim(klaim.id)
      setKlaim(null)
      setDampak(null)
    } catch (err: unknown) {
      const kode = kodeStatusGalat(err)
      setGalat(
        kode === 501
          ? 'Penghapusan klaim adalah penanda, bukan hapus fisik (ADR-U-0031). Kolom penandanya belum diputuskan work owner.'
          : kode === 409
            ? 'Klaim sudah diserahkan ke Komite dan tidak dapat dihapus.'
            : kode === 403
              ? 'Hanya ReasLifeAdmin yang dapat menghapus klaim.'
              : 'Gagal menghapus klaim.',
      )
      setDampak(null)
    } finally {
      setMenghapus(false)
    }
  }

  // Menolak satu baris, lalu MEMBACA ULANG klaimnya dari server.
  //
  // Dibaca ulang, bukan diubah di layar: status klaim adalah TURUNAN dari
  // seluruh barisnya, dan menebaknya di sisi klien berarti dua sumber
  // kebenaran yang dapat berbeda.
  async function aksep(pesertaID: string) {
    if (!klaim) return
    setGalat(null)
    setMengaksep(pesertaID)
    try {
      const nomor = await simpanAdjustment(klaim.id, pesertaID)
      setKlaim(await ambilKlaimLife(klaim.id))
      setGalat(`Baris diaksep dengan nomor ${nomor}.`)
    } catch (err: unknown) {
      const kode = kodeStatusGalat(err)
      setGalat(
        kode === 403
          ? 'Hanya pemegang tahap klaim ini yang dapat mengaksep barisnya.'
          : kode === 409
            ? (pesanGalat(err) ?? 'Baris tidak dalam keadaan yang dapat diaksep.')
            : kode === 501
              ? (pesanGalat(err) ?? 'Akseptasi belum dapat disimpan.')
              : 'Gagal mengaksep baris.',
      )
    } finally {
      setMengaksep(null)
    }
  }

  async function putaranBaru(pesertaID: string) {
    if (!klaim) return
    setGalat(null)
    setMemutar(pesertaID)
    try {
      await tambahPutaran(klaim.id, pesertaID)
      setKlaim(await ambilKlaimLife(klaim.id))
    } catch (err: unknown) {
      const kode = kodeStatusGalat(err)
      setGalat(
        kode === 403
          ? 'Hanya ReasLifeSPV yang dapat membuka putaran berikutnya.'
          : kode === 409
            ? 'Putaran berikutnya hanya lahir sesudah baris terakhir ditolak.'
            : kode === 501
              ? 'Jejak audit belum dapat direkam; tempatnya belum diputuskan.'
              : 'Gagal membuka putaran berikutnya.',
      )
    } finally {
      setMemutar(null)
    }
  }

  async function serahkan(pesertaID: string, adjID: string) {
    if (!klaim) return
    setGalat(null)
    setMenyerahkan(adjID)
    try {
      await serahkanKeKomite(klaim.id, pesertaID, adjID)
      setKlaim(await ambilKlaimLife(klaim.id))
    } catch (err: unknown) {
      const kode = kodeStatusGalat(err)
      setGalat(
        kode === 403
          ? 'Peran Anda tidak berwenang menyerahkan baris bertipe ini ke Komite.'
          : kode === 409
            ? 'Baris sudah diserahkan, atau bukan lagi Outstanding.'
            : kode === 422
              ? // ⛔ Pesan 422 dicetak APA ADANYA. Salah satunya berbunyi
                // "Name of bank cannot be empty" - kalimat sistem lama, yang
                // dikenali pengguna lama. Menggantinya dengan kalimat kita
                // sendiri memutus pengenalan itu.
                (pesanGalat(err) ?? 'Penyerahan ditolak.')
              : kode === 501
                ? 'Penyerahan belum dapat disimpan; tempatnya belum diputuskan.'
                : 'Gagal menyerahkan baris ke Komite.',
      )
    } finally {
      setMenyerahkan(null)
    }
  }

  async function tolak(adjID: string) {
    if (!klaim) return
    setGalat(null)
    setMenolak(adjID)
    try {
      await tolakBarisAdjustment(klaim.id, adjID)
      setKlaim(await ambilKlaimLife(klaim.id))
    } catch (err: unknown) {
      const kode = kodeStatusGalat(err)
      setGalat(
        kode === 403
          ? 'Hanya ReasLifeAdmin yang dapat menolak baris.'
          : kode === 409
            ? 'Baris sudah diputus dan tidak dapat ditolak lagi.'
            : kode === 422
              ? 'Klaim belum bernomor.'
              : kode === 501
                ? 'Jejak audit belum dapat direkam; tempatnya belum diputuskan.'
                : 'Gagal menolak baris.',
      )
    } finally {
      setMenolak(null)
    }
  }

  // Dipanggil saat tombol "Buka" ditekan (form dikirim).
  async function cari(e: FormEvent<HTMLFormElement>) {
    e.preventDefault() // jangan biarkan browser memuat ulang halaman
    setGalat(null)
    setKlaim(null)
    setMemuat(true)
    try {
      setKlaim(await ambilKlaimLife(id.trim()))
    } catch (err: unknown) {
      setGalat(kodeStatusGalat(err) === 404 ? 'Klaim tidak ada.' : 'Gagal membaca klaim.')
    } finally {
      setMemuat(false) // apa pun hasilnya, berhenti "memuat"
    }
  }

  return (
    <section>
      <h2>Klaim Life</h2>

      <form onSubmit={cari}>
        <label htmlFor="idKlaim">Pengenal klaim</label>{' '}
        <input
          id="idKlaim"
          value={id}
          onChange={(e) => setId(e.target.value)}
          placeholder="mis. UJI-KLAIM-1"
        />{' '}
        <button type="submit" disabled={memuat || id.trim() === ''}>
          {memuat ? 'Memuat…' : 'Buka'}
        </button>
      </form>

      {/* `galat && (...)` = tampilkan hanya bila galat terisi. */}
      {galat && <p role="alert">{galat}</p>}

      {klaim && (
        <article>
          <h3>{klaim.nomorKlaim || klaim.id}</h3>
          <dl>
            <dt>Nomor polis</dt>
            <dd>{klaim.nomorPolis || '—'}</dd>
            <dt>Nama bisnis</dt>
            <dd>{klaim.namaBisnis || '—'}</dd>
            <dt>Jumlah baris</dt>
            <dd>{klaim.cacahBaris}</dd>
            {/* Status klaim adalah TURUNAN dari baris-barisnya, bukan kolom
                tersimpan. Yang ditampilkan kata, bukan angka - dan bukan pula
                terjemahan `kodeStatus`, yang hanya cerminan baris terakhir. */}
            <dt>Status klaim</dt>
            <dd>{klaim.statusTurunan}</dd>
          </dl>

          {/* `.map` = ulangi blok di bawah untuk SETIAP peserta.
              `key` wajib ada supaya React tahu baris mana yang berubah. */}
          <p>
            <button type="button" onClick={() => void bukaKonfirmasiHapus()}>
              Hapus klaim…
            </button>
          </p>

          {/* Popup konfirmasi: angka dulu, keputusan kemudian. Tiap jenis
              disebut sendiri - satu angka total menyembunyikan tingkat mana
              yang ternyata lebih besar dari dugaan. */}
          {dampak && (
            <aside>
              <h4>Hapus klaim {klaim.nomorKlaim || klaim.id}?</h4>
              <p>Yang akan ikut terhapus:</p>
              <ul>
                <li>Header klaim: {dampak.header}</li>
                <li>Peserta: {dampak.peserta}</li>
                <li>Baris adjustment: {dampak.adjustment}</li>
                <li>Spreading: {dampak.spreading}</li>
                <li>Spreading retro: {dampak.spreadingRetro}</li>
                <li>Dokumen: {dampak.dokumen}</li>
                <li>Baris work: {dampak.barisWork}</li>
              </ul>
              <p>Total {dampak.total} baris.</p>
              {/* ⛔ Baris datar warisan BUKAN bagian daftar di atas dan tidak
                  dijumlahkan ke totalnya. Nasibnya saat klaim dihapus belum
                  diputuskan work owner; menaruhnya di bawah judul "yang akan
                  ikut terhapus" berarti menjawab pertanyaan itu diam-diam. */}
              <p>
                ⚠️ Selain itu terdapat {dampak.barisDatarWarisan} baris datar warisan
                ber-CASEID sama. <strong>Nasibnya belum diputuskan</strong>: apakah ikut
                terhapus atau ditinggal karena hilir sudah membacanya.
              </p>
              <button
                type="button"
                disabled={menghapus}
                onClick={() => void lanjutkanHapus()}
              >
                {menghapus ? 'Menghapus…' : 'Ya, hapus'}
              </button>{' '}
              <button type="button" onClick={batalHapus}>
                Batal
              </button>
            </aside>
          )}

          {klaim.peserta.map((p) => (
            <section key={p.id}>
              <h4>
                Peserta {p.nomorSertifikat || p.id}{' '}
                <small>({p.baris.length} baris)</small>
              </h4>

              {/* ⛔ Penolakan bukan akhir: klaim TIDAK terminal (ADR-U-0011).
                  Yang terminal adalah baris, dan baris berikutnya memulai
                  putaran baru dengan angka yang diperbaiki. */}
              {/* ⭐ Jalur akseptasi Claim Life sendiri - SaveAdjustment_Act.
                  Prasyaratnya dari XML: peserta DIPILIH, baris terakhir belum
                  bernomor, dan masih Outstanding. Perannya TIDAK diperiksa di
                  layar: pohon XML membuktikan tombolnya tidak bergerbang
                  peran, dan yang menggerbanginya pemegang tahap. */}
              {bolehSimpanAdjustment(p, p.baris) && (
                <p>
                  <button
                    type="button"
                    disabled={mengaksep === p.id}
                    onClick={() => void aksep(p.id)}
                  >
                    {mengaksep === p.id ? 'Menyimpan…' : 'Save Adjustment'}
                  </button>
                </p>
              )}

              {bolehPutaranBaru(p.baris) && (
                <p>
                  <button
                    type="button"
                    disabled={memutar === p.id}
                    onClick={() => void putaranBaru(p.id)}
                  >
                    {memutar === p.id ? 'Membuka…' : 'Putaran berikutnya'}
                  </button>
                </p>
              )}

              {p.baris.length === 0 ? (
                <p>Belum ada baris adjustment.</p>
              ) : (
                <table>
                  <thead>
                    <tr>
                      <th>Baris</th>
                      <th>Status</th>
                      <th>Jumlah klaim</th>
                      <th>Nomor akseptasi</th>
                      <th>Tindakan</th>
                      <th>Komite</th>
                    </tr>
                  </thead>
                  <tbody>
                    {p.baris.map((b) => (
                      <tr key={b.id}>
                        <td>{b.id}</td>
                        <td>
                          {/* Tiket 01 AC-3: status sebagai KATA — bukan nama field,
                              bukan angka. Kode mentahnya (b.kodeStatus) tetap ada di
                              kontrak API bagi yang menelusuri, tetapi tidak dicetak. */}
                          {b.status}
                        </td>
                        <td>{tampilUang(b.jumlahKlaim) || '—'}</td>
                        <td>{b.nomorAkseptasi || '—'}</td>
                        <td>
                          {/* Gerbang XML: `pyPosition=='ReasLifeAdmin' &&
                              CLAIM_NO != '' && .STS_REJECT == 0`. Dua syarat
                              terakhir dapat diperiksa di sini; PERANNYA tidak —
                              layar tidak tahu peran siapa pun, dan yang
                              menegakkannya services (403). Menyembunyikan
                              tombol di sini adalah kenyamanan, bukan pagar. */}
                          {b.status === STATUS_OUTSTANDING && klaim.nomorKlaim ? (
                            <button
                              type="button"
                              disabled={menolak === b.id}
                              onClick={() => void tolak(b.id)}
                            >
                              {menolak === b.id ? 'Menolak…' : 'Reject Outstanding'}
                            </button>
                          ) : (
                            '—'
                          )}
                        </td>
                        <td>
                          {/* Wewenangnya bergantung Type (QP/QR hanya SPV), dan
                              layar TIDAK tahu peran siapa pun. Yang diperiksa di
                              sini hanya keadaan baris; perannya ditegakkan
                              services (403). Menyembunyikan tombol adalah
                              kenyamanan, bukan pagar. */}
                          {bolehSerahkanDiLayar(b) ? (
                            <button
                              type="button"
                              disabled={menyerahkan === b.id}
                              onClick={() => void serahkan(p.id, b.id)}
                            >
                              {menyerahkan === b.id ? 'Menyerahkan…' : 'Send ke Komite'}
                            </button>
                          ) : b.komiteId ? (
                            <span title={b.komiteId}>sudah diserahkan</span>
                          ) : (
                            '—'
                          )}
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              )}
            </section>
          ))}
        </article>
      )}

      {/* ⛔ Kelima total layar Detail DINYATAKAN belum bersumber, bukan
          dihilangkan dan bukan dijumlahkan sendiri. Rule penghitungnya
          (`CheckTotalAdjustmentClaim`) dirujuk sepuluh kali di
          `ClaimLifeDetailGCNM.xml` tetapi NOL berkasnya ada di ekspor, jadi
          baris mana yang ikut dihitung belum terjawab - dan ini angka uang
          (ADR-U-0003). Lihat OQ-H. */}
      <PanelTotalKlaim />
    </section>
  )
}
